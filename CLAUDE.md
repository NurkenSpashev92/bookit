# Project: Bookit Go API

## Code Style & Engineering Standards

Always respond and write code as a Senior Go Engineer.

### Principles

* Strictly follow SOLID principles in every package, struct, interface, and method
* Apply Clean Architecture, DDD, OOP-style composition, and separation of concerns
* Follow DRY, KISS, Clean Code — small focused functions, meaningful names, predictable behavior
* Prefer composition over inheritance
* Keep business logic framework-independent
* Use interfaces only when there is more than one implementation or when mocking is required in tests
* Avoid god structs and large service layers
* Keep packages cohesive and focused on a single responsibility
* Every public method and struct must have a clear purpose

---

## Architecture

```text
HTTP Request
    → Handler
    → Schema / DTO
    → Service
        → Factory
        → Mapper
        → Repository
    → Response DTO
```

### Folder Structure

```text
app/
├── cmd/                # entrypoints
│   ├── apiserver
│   ├── router
│   └── seed            # demo data generator
├── configs/            # config loading and env handling
├── docs/               # swagger documentation
├── internal/
│   ├── shared/         # cross-domain DTOs, validation, error responses
│   ├── initializers/   # database, redis, storage bootstrapping
│   └── <domain>/       # analytics, booking, content, identity,
│       ├── handler/    # interaction, location, property
│       │               # handlers + routes.go (domain route registrar)
│       ├── model/      # domain entities + domain errors
│       ├── port/       # interfaces to other domains (dependency inversion)
│       ├── repository/ # persistence layer
│       ├── schema/     # request / response DTOs
│       └── service/    # business logic
├── migrations/         # SQL migrations
├── pkg/
│   ├── aws
│   ├── cache
│   ├── imageproc
│   ├── middleware
│   ├── store
│   └── utils
```

Code is split by domain, not by layer: everything about a domain lives under
`internal/<domain>/`. Cross-domain calls go through `port/` interfaces so the
import graph stays acyclic.

### Routing and composition root

* Every domain owns its routes in `internal/<domain>/handler/routes.go` and
  exposes `RegisterRoutes(api fiber.Router, deps Deps, guards shared.Guards)`;
  the registrar constructs its handlers from the services it is given
* `cmd/router/router.go` only holds the middleware chain and calls the domain
  registrars — no paths, no handler construction
* `cmd/apiserver/container.go` is the only place that constructs repositories
  and services; `main.go` owns just the process lifecycle
* Dictionary resources share `shared.RegisterCRUD` — public reads,
  authenticated writes
* Inside a group, literal paths are registered before `/:param` ones; where two
  domains share a prefix (`/houses`), the registration order in `router.go`
  keeps that invariant and `test/routes_test.go` locks it

---

## Handlers

Handlers are responsible only for:

* Parsing request
* Validation
* Calling service layer
* Returning JSON response

Handlers must not:

* Contain business logic
* Access database directly
* Perform SQL queries
* Instantiate repositories manually
* Contain complex conditional flows

Handlers depend on an interface declared in the handler package, never on a
concrete `*service.XService` — that keeps the direction of dependencies inward
and makes the handler testable with a stub.

Every error response goes through `shared.Fail(c, err)`. The status comes from
the error itself, so handlers never map errors to status codes and never hand a
client `err.Error()` directly. 5xx details are logged and answered with a
neutral message.

```go
type HouseService interface {
    Create(ctx context.Context, req schema.CreateHouseRequest, ownerID int) (model.House, error)
}

func (h *HouseHandler) Create(c fiber.Ctx) error {
    user, err := middleware.CurrentUser(c)
    if err != nil {
        return shared.Fail(c, err)
    }

    request, err := shared.Bind[schema.CreateHouseRequest](c)
    if err != nil {
        return shared.Fail(c, err)
    }

    house, err := h.houseService.Create(c.Context(), request, user.ID)
    if err != nil {
        return shared.Fail(c, err)
    }

    return shared.Created(c, house)
}
```

Shared handler kit (`internal/shared`):

* `Bind[T](c)` — decodes the body and runs `Validate()` when the DTO has one
* `ParamInt(c, name)` / `ParamString(c, name)` — required path parameters
* `QueryIntInRange(c, name, fallback, min, max)` — optional numeric query
* `Fail(c, err)`, `OK(c, msg)`, `Created(c, payload)`, `List(c, items)`
* `middleware.CurrentUser(c)` / `middleware.CurrentUserID(c)` — authenticated user

---

## Errors

Domain errors are declared once per domain in `internal/<domain>/model/errors.go`
as typed values, and carry the kind the transport layer needs:

```go
var (
    ErrHouseNotFound = shared.NotFound("house not found")
    ErrSlugExists    = shared.Conflict("slug already exists")
)
```

Rules:

* Constructors: `shared.Invalid`, `Unauthorized`, `Forbidden`, `NotFound`, `Conflict`
* Kinds map to statuses in one place (`shared.StatusOf`); handlers never switch on errors
* Repositories translate storage failures with `store.MapNoRows(err, model.ErrX)`
* Wrap with `fmt.Errorf("...: %w", ErrX)` to add context — `errors.Is` and the
  status mapping keep working, and the client still sees the clean message
* Anything untyped is treated as an internal failure: 500 with a neutral body

---

## Schemas / DTOs

Schemas belong in `internal/<domain>/schema/`.

Rules:

* Separate request and response DTOs
* DTOs must not contain business logic
* DTOs should be small and explicit
* Use validation tags
* Never expose internal models directly in handlers

Example:

```go
type CreateHouseRequest struct {
    Title       string `json:"title" binding:"required,min=3,max=255"`
    Description string `json:"description" binding:"required"`
    CityID      uint   `json:"city_id" binding:"required"`
}
```

---

## Services

Services contain business logic only.

Services must:

* Work through interfaces
* Not know about HTTP layer
* Not return fiber.Ctx or HTTP responses
* Use repositories for persistence
* Use factories for entity creation
* Use mappers for DTO to model conversion
* Use transactions when multiple repositories are involved

Example:

```go
type HouseService struct {
    repository HouseRepository
    factory    HouseFactory
    mapper     HouseMapper
}

func (s *HouseService) Create(ctx context.Context, dto schemas.CreateHouseRequest) (*schemas.HouseResponse, error) {
    house := s.factory.Create()

    s.mapper.MapCreateRequest(house, dto)

    if err := s.repository.Create(ctx, house); err != nil {
        return nil, err
    }

    return s.mapper.ToResponse(house), nil
}
```

---

## Repositories

Repositories are responsible only for persistence and queries.

Rules:

* Declare repository interfaces close to service layer or in repositories/contracts
* Use SQL (pgx) only inside repositories
* Do not place business logic in repositories
* Keep repository methods small and explicit
* Always pass context.Context

Example:

```go
type HouseRepository interface {
    Create(ctx context.Context, house *models.House) error
    FindByID(ctx context.Context, id uint) (*models.House, error)
    FindAll(ctx context.Context) ([]models.House, error)
    Delete(ctx context.Context, id uint) error
}
```

---

## Factories

Factories create new domain entities.

Rules:

* Never instantiate models directly in services
* All new entities must be created through factories
* Factory methods should initialize defaults when needed

Example:

```go
type HouseFactory struct{}

func (f *HouseFactory) Create() *models.House {
    return &models.House{
        IsActive: true,
    }
}
```

---

## Mappers

Mappers transform DTOs into models and models into responses.

Rules:

* No persistence logic
* No business logic
* One responsibility only: mapping

Example:

```go
type HouseMapper struct{}

func (m *HouseMapper) MapCreateRequest(model *models.House, dto schemas.CreateHouseRequest) {
    model.Title = dto.Title
    model.Description = dto.Description
    model.CityID = dto.CityID
}

func (m *HouseMapper) ToResponse(model *models.House) *schemas.HouseResponse {
    return &schemas.HouseResponse{
        ID:          model.ID,
        Title:       model.Title,
        Description: model.Description,
    }
}
```

---

## Models

Rules:

* Models represent domain entities
* Keep model methods minimal
* Do not place HTTP logic in models
* Use explicit field names
* Keep struct tags minimal — pgx scans by column position, not by tag
* Add created_at and updated_at timestamps where needed

Example:

```go
type House struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CityID      *int      `json:"city_id,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## Database

* Always use migrations for schema changes
* Never change production tables manually
* Avoid SELECT *
* Always specify needed columns
* Use indexes for frequently queried fields
* Avoid N+1 queries
* Fetch related data with JOIN / LATERAL in the same query, never in a loop
* Use transactions for atomic operations
* Keep SQL readable and optimized

---

## Testing

Follow TDD whenever possible.

Rules:

* Write unit tests for services
* Mock repositories in service tests
* Test handlers separately
* Cover happy path and edge cases
* Keep tests isolated and deterministic
* Use table-driven tests

Example:

```go
func TestHouseService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   schemas.CreateHouseRequest
        wantErr bool
    }{
        {
            name: "success",
            input: schemas.CreateHouseRequest{
                Title: "House",
                CityID: 1,
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

---

## Naming Rules

* Use singular names for models: House, User, Country
* Use plural names for repository collections: FindHouses, ListCities
* Interfaces should end with -Repository, -Service, -Store when appropriate
* File names should be lowercase with underscores only when necessary
* Keep package names short and meaningful

---

## Forbidden

```text
❌ Business logic inside handlers
❌ SQL queries inside handlers
❌ Direct DB access inside services
❌ Large god services with 1000+ lines
❌ Shared utils for unrelated logic
❌ Repeating validation rules everywhere
❌ Returning raw domain models directly from handlers
❌ Panic for expected business errors
❌ Global mutable state
❌ Circular dependencies
```

---

## Preferred Stack

* Fiber v3 for HTTP layer
* pgx/v5 for database access (no ORM — plain SQL in repositories)
* PostgreSQL as primary database
* Redis for cache
* Swagger for API docs
* Docker + Docker Compose for local development
* Testify for assertions and mocks
* golang-migrate for migrations

---

## Senior Engineering Expectations

* Always propose scalable solutions
* Always think about performance and maintainability
* Always suggest indexes when needed
* Always consider caching for expensive queries
* Always think about concurrency and race conditions
* Always think about context cancellation and timeouts
* Always prefer explicit code over magic
* Always keep code easy to read and easy to test

---

## Architecture Map — where things live (AI navigation)

Read this first when locating or adding code. All paths are under `app/`.

### Composition root & process
* `cmd/apiserver/main.go` — process lifecycle only (start/stop/signals)
* `cmd/apiserver/container.go` — **the only place** that constructs repositories
  and services and wires them into `router.Services`
* `cmd/router/router.go` — middleware chain + calls each domain's
  `RegisterRoutes`. Holds `cachedResponseRoutes` (path→namespace/TTL) and
  `Guards{Required, Optional, Admin}`. No paths or handler construction here.
* `cmd/seed/` — demo-data generator (`data.go` seed lists, `refs.go` upserts,
  `relations.go` house↔category/convenience links, `images.go` bounded S3 upload)

### Domains — `internal/<domain>/`
`analytics · booking · content · identity · interaction · location · property`
Each domain owns the full vertical stack:
* `handler/` — HTTP parse/validate/call-service/respond + `routes.go` registrar
  (`RegisterRoutes(api, deps, guards)`); handlers depend on a **local interface**,
  never on `*service.X`
* `schema/` — request/response DTOs (validation tags + `Validate()`)
* `service/` — business logic; works through interfaces; no HTTP
* `repository/` — plain SQL (pgx); the only place with SQL
* `model/` — entities + `errors.go` (typed domain errors)
* `port/` — interfaces to OTHER domains (dependency inversion; keeps import
  graph acyclic). Example: `property/port/{booking_checker,like_checker}.go`
* `factory/` / `mapper/` where present — entity creation / DTO↔model mapping

### Cross-cutting
* `internal/shared/` — handler kit reused by every domain:
  `Bind[T]`, `ParamInt/ParamString`, `QueryIntInRange`, `Fail/OK/Created/List/Items[T]`,
  `RegisterCRUD`, `Guards`, error kinds (`Invalid/Unauthorized/Forbidden/NotFound/Conflict`,
  `StatusOf`, `MessageOf`), pagination (`WantsPagination`, `ParsePageParams`,
  `PageEnvelope`, `Paginated`, `ListMaybePaginated[T]`)
* `internal/initializers/` — `cache.go` (Redis+local), `fiber.go` (app bootstrap)
* `internal/platform/healthcheck/` — liveness/readiness
* `pkg/` — framework-agnostic infra: `aws` (S3), `cache` (two-tier), `imageproc`,
  `logger` (buffered + rotating), `middleware` (auth/cors/recover/response_cache/
  upload/request_logger), `store` (pg pool, `MapNoRows`), `upload` (per-resource
  image uploader), `utils` (`slug`)
* `configs/` — env loading (`NewDBConfig`, `NewRedisConfig`, `NewCacheConfig`,
  `NewLogConfig`, `NewAwsConfig`)
* `migrations/` — golang-migrate SQL (`NNNNNN_name.{up,down}.sql`)
* `docs/` — generated Swagger (do not hand-edit)

### Where to add a new feature
* New endpoint on an existing resource → that domain's `handler/` + its `routes.go`
* New reference/dictionary resource → `service`+`repository`+`handler`, register
  reads public / writes admin via `shared.RegisterCRUD`
* New cross-domain dependency → declare a `port/` interface, inject the concrete
  service in `container.go`
* New table/column → a migration pair, never manual DDL

---

## Patterns used in this codebase (reuse these, don't reinvent)

* **Optional pagination**: list handlers return a plain array by default and a
  `PaginatedResponse` envelope when `?page=` is present. Use the single generic
  `shared.ListMaybePaginated(c, svc.GetAll, svc.GetAllPaginated)` — do **not**
  re-inline the `if WantsPagination { … }` block per handler.
* **Errors**: return typed `shared.*` errors from services; handlers call only
  `shared.Fail(c, err)`. Repositories translate storage errors with
  `store.MapNoRows(err, model.ErrX)`. Never map errors to status codes in handlers.
* **CRUD dictionaries**: `shared.RegisterCRUD(group, guards, handler)` — public
  reads, `guards.Admin` writes.
* **Two-tier cache** (`pkg/cache`): `Cache` (Redis L2) fronted by `LocalCache`
  (in-process L1, TTL capped ~2s). Read via `GetBytes`; write via
  `SetBytes`/`SetBytesTTL`; invalidate via `InvalidateNamespace(ns)` (clears both
  `ns:` and `resp:ns:`). Cache-heavy reads use `singleflight` to collapse stampedes
  (see `HouseService.GetAllPaginated`).
* **HTTP response cache**: declarative — add the route to `cachedResponseRoutes`
  in `router.go` with a namespace + TTL; middleware caches anonymous GET 200s only.
  Every cached namespace MUST have write-path `InvalidateNamespace`.
* **Concurrency**: independent I/O in a request → `errgroup` (see house list
  count∥select, analytics charts). Bound CPU/upload fan-out with `group.SetLimit(n)`.
  Fire-and-forget background work → the `recordViewAsync` model: bounded semaphore
  channel + **detached** `context.Background()` + timeout + non-blocking send that
  drops under load. Never reuse the request `ctx` in a goroutine that outlives the response.
* **Nested reads**: fold child collections into one SQL query via LATERAL +
  `json_agg` (see house images) — never per-row loops.

## Performance conventions

* Add `WHERE is_active` (and matching partial/composite indexes) to public list
  queries; suggest `(filter_col, id DESC)` composites for filter+ordering paths.
* Prefer `= ANY($ids::int[])` scoped to the current page over fetching a full set
  to annotate a page (e.g. liked-house membership).
* Avoid per-request `UPDATE` on wide hot tables (counter churn → MVCC bloat);
  derive counts from event tables or a narrow `*_counters` table.
* Cache expensive auth-only aggregates (stats) at the service level by owner key
  with a short TTL; TTL-based expiry is fine when exact freshness isn't required.
* See **`OPTIMIZATION.md`** (repo root) for the full, ranked, `file:line`-grounded
  performance audit and the recommended low-risk quick-win batch.
