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
│   └── router
├── configs/            # config loading and env handling
├── docs/               # swagger documentation
├── internal/
│   ├── handlers/       # HTTP handlers only
│   ├── initializers/   # database, redis, storage bootstrapping
│   ├── models/         # domain models and entities
│   ├── repositories/   # persistence layer
│   ├── schemas/        # request / response DTOs
│   └── services/       # business logic
├── migrations/         # SQL migrations
├── pkg/
│   ├── aws
│   ├── middleware
│   ├── store
│   └── utils
```

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

Example:

```go
func (h *HouseHandler) Create(c *gin.Context) {
    var request schemas.CreateHouseRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.houseService.Create(c.Request.Context(), request)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, response)
}
```

---

## Schemas / DTOs

Schemas belong in `internal/schemas/`.

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
* Not return gin.Context or HTTP responses
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
* Use GORM or SQL only inside repositories
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
* Use GORM tags carefully
* Add created_at and updated_at timestamps where needed

Example:

```go
type House struct {
    ID          uint      `gorm:"primaryKey"`
    Title       string    `gorm:"size:255;not null"`
    Description string    `gorm:"type:text"`
    CityID      uint      `gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
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
* Use Preload() when related data is needed
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
❌ Returning raw GORM models directly from handlers
❌ Panic for expected business errors
❌ Global mutable state
❌ Circular dependencies
```

---

## Preferred Stack

* Gin for HTTP layer
* GORM for ORM
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
