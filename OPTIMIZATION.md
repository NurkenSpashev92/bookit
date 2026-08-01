# Bookit — Performance & Optimization Report

> Аудит горячих путей: БД/N+1, конкурентность/воркеры, кэш.
> Все ссылки — `file:line` относительно `app/`. Статусы: ✅ уже сделано, ⚠️ возможность, 🐘 нужна миграция.

## TL;DR

Кодовая база **уже сильно оптимизирована**. Классического N+1 нет; есть worker-pool
на загрузку картинок, фоновый воркер просмотров, `singleflight` + двухуровневый кэш
на список домов, `errgroup` на count+select и аналитику, LATERAL + `json_agg` вместо
per-row подзапросов. «Насыпать горутин» не нужно — выигрыши точечные и перечислены ниже.

---

## ✅ Что уже реализовано (использовать как образец при доработках)

| Паттерн | Где | Модель для |
|---|---|---|
| Фоновый воркер (семафор + detached ctx + дроп под нагрузкой) | `internal/property/service/house.go:51` `recordViewAsync` (`maxConcurrentViewWrites=32`) | любой новый async (нотификации, cache warming) |
| Bounded upload-пул | `internal/property/service/image.go:96` `uploadGroup.SetLimit(6)` | параллельная загрузка в S3 |
| Bounded пул в seed | `cmd/seed/images.go:92` `SetLimit(uploadConcurrency=16)` | оффлайн-массовые операции |
| `singleflight` + двухуровневый кэш на список | `internal/property/service/house.go:68` | защита от thundering herd |
| `errgroup` count ∥ select | `internal/property/repository/house.go:97` | пагинированные списки |
| `errgroup` на 4 агрегата аналитики | `internal/analytics/service/stats.go:42` | параллельные независимые запросы |
| LATERAL + `json_agg` (нет N+1 по картинкам) | `internal/property/repository/house.go` | вложенные коллекции одним запросом |
| Двухуровневый кэш Redis(L2)+local(L1) | `pkg/cache/cache.go`, `pkg/cache/memory.go` | все read-through кэши |
| Инвалидация namespace (оба уровня) | `internal/property/service/house.go:176/183/191` `InvalidateNamespace("houses")` | сброс кэша при мутации |

---

## ⚠️ Приоритет 1 — БД, горячий путь

### P1.1 — `GetUserLikedHouseIDs` тянет ВСЕ лайки юзера на каждый запрос списка
`internal/interaction/repository/house_like.go:99` (вызов: `internal/property/service/house.go:105`)

**Проблема:** `SELECT house_id FROM house_likes WHERE user_id=$1 ORDER BY created_at DESC`
загружает все лайки (сотни строк), чтобы разметить 20 на странице; `ORDER BY` бесполезен —
результат используется как множество (`applyLiked`, `house.go:120`).

**Фикс:** передавать ID текущей страницы →
`SELECT house_id FROM house_likes WHERE user_id=$1 AND house_id = ANY($2::int[])`, снять `ORDER BY`.
Обслуживается индексом `uq_house_likes_user_house` как index-only scan. Ограничено размером страницы.

**Риск:** низкий, локальный (страница уже доступна в `applyLiked`). Без миграции.

### P1.2 — `view_count` UPDATE на каждый просмотр раздувает самую читаемую таблицу
`internal/property/repository/house.go:326` `RecordView` (из `GetBySlug`, `house.go:152`)

**Проблема:** `UPDATE houses SET view_count = view_count+1` на каждый просмотр → новые MVCC-версии
широкой строки `houses`, блокировки на популярном листинге, dead-tuple bloat → тормозятся все
list/detail чтения + autovacuum. То же для `like_count` через триггеры (`migrations/000011:16`).

**Фикс:** каждый просмотр уже пишется в `house_views` (000014). Считать счётчик оттуда
(кэш/периодический rollup), либо вынести счётчики в отдельную `house_counters` (1 строка/дом),
чтобы churn не раздувал горячую таблицу.

**Риск:** средний — решение по схеме. 🐘 миграция.

### P1.3 — Публичный список не фильтрует `is_active`
`internal/property/repository/house.go:126` (list) + `:99` (count)

**Проблема:** нет `WHERE h.is_active = TRUE` → (а) неактивные дома утекают в публичный список
(корректность); (б) partial-индексы `ix_houses_active/best/promo` (`migrations/000011:56`) **мертвы** —
только оверхед на запись.

**Фикс:** добавить `WHERE h.is_active = TRUE` в builder list/count. Оживит `ix_houses_active(id DESC) WHERE is_active`
для `ORDER BY h.id DESC` как ordered index scan.

**Риск:** низкий по коду, но меняет выдачу (корректность — желаемо). Индекс уже есть.

### P1.4 — Нет композитных индексов под «фильтр + ORDER BY id DESC»
`internal/property/repository/house.go:54` (фильтры), `:159` (`ORDER BY h.id DESC`)

**Проблема:** фильтры по `city_id`/`country_id`/`price`/… покрыты одиночными btree; Postgres
использует один и делает отдельную сортировку. Под самый частый «город X, свежие первыми» нет
`(city_id, id DESC)`.

**Фикс:** 🐘 `CREATE INDEX ON houses(city_id, id DESC) WHERE is_active;` (+ `(country_id, id DESC)` при
использовании). Даёт выдачу уже в порядке — без сортировки.

**Риск:** только запись/размер индекса. 🐘 миграция.

### P1.5 — `COUNT(*)` + `SELECT` двумя проходами (мягкое)
`house.go:99` vs `:111`; аналогично во всех `*Paginated` репозиториях

**Проблема:** два независимых скана на страницу; count — unbounded, растёт с таблицей; глубокий OFFSET дорожает.

**Фикс (если появится в нагрузке):** `count(*) OVER () AS total` одним сканом, либо keyset-пагинация
(`WHERE id < $lastId ORDER BY id DESC LIMIT n`) — идеально ложится на индексы из P1.4.

**Риск:** низкий приоритет — уже параллельно (`errgroup`) и кэшируется.

---

## ⚠️ Приоритет 2 — Конкурентность

### P2.1 — `GetBySlug`: 3 запроса подряд на самом горячем эндпоинте
`internal/property/service/house.go:143`

**Проблема:** `GetBySlug` → `StatusWithCount` → `getActiveBooking` строго последовательно.
`StatusWithCount` зависит только от `slug`+`userID` (не от результата первого запроса).

**Фикс:** `errgroup` — гнать `GetBySlug` ∥ `StatusWithCount`, `getActiveBooking` после (нужен `house.ID`).
Экономит ~1 RTT на каждой детальной. **Гонка:** писать в локальные переменные, присваивать в `house`
после `Wait()`.

**Риск:** низкий, если аккуратно с общим `house`. Без миграции.

### P2.2 — `procGroup` без лимита (декодирование картинок)
`internal/property/service/image.go:82`

**Проблема:** группа обработки (decode+resize+encode, CPU-bound) без `SetLimit`, хотя upload-группа
ограничена 6. При `maxHouseImages=15` — до 15 одновременных декодов bitmap = скачок памяти/CPU без выигрыша.

**Фикс:** `procGroup.SetLimit(runtime.GOMAXPROCS(0))` — одна строка.

**Риск:** минимальный. Без миграции.

### P2.3 — Инициализация зависимостей при старте последовательна (низкий приоритет)
`cmd/apiserver/container.go:46`

**Проблема:** pgx pool connect, Redis ping (`internal/initializers/cache.go:23`) — независимые сетевые RTT
подряд. Влияет только на время старта.

**Фикс:** `errgroup` для pool ∥ ping (S3-клиент сети не делает). Сохранить порядок cleanup.

**Риск:** низкий, разовый выигрыш.

### P2.4 — Удаление старого аватара из S3 инлайн (низкий приоритет)
`internal/identity/service/avatar.go:55` (и `:78`)

**Проблема:** после апдейта юзер ждёт сетевой RTT на удаление старого объекта без функциональной нужды.

**Фикс:** fire-and-forget по модели `recordViewAsync` — фоновая горутина с `context.Background()` + таймаут
(request-ctx мёртв после ответа). Осиротевшие удаления допустимо дропать.

**Риск:** низкий.

---

## ⚠️ Приоритет 3 — Кэш

### P3.1 — Детальная страница `GetBySlug` не кэшируется вообще
`internal/property/service/house.go:142`

**Проблема:** роут `/houses/` не матчит `/houses/{slug}` (`cmd/router/router.go:37`); запрос на каждый вызов.
Самый частый detail-эндпоинт.

**Фикс:** кэшировать анонимную базу под `houses:detail:<slug>` (существующий `InvalidateNamespace("houses")`
уже покроет инвалидацию при update/delete), персональные `IsLiked/IsBooked` накладывать после — как в
`GetAllPaginated`. TTL 5м.

### P3.2 — Stats не кэшируется (auth-only, не ловится response-cache)
`internal/analytics/service/stats.go`

**Проблема:** `GetCharts` = 4 тяжёлых агрегата на вызов; эндпоинты auth-only → response-cache пропускает
(`pkg/middleware/response_cache.go:40`). Меняются медленно (просмотры/лайки).

**Фикс:** service-level кэш по владельцу: `stats:dashboard:<ownerID>`, `stats:charts:<ownerID>:<days>` и т.д.
TTL 30–60с, без инвалидации (по TTL).

### P3.3 — Reference-кэш без инвалидации при записи
`internal/property/service/{category,type}.go`, `internal/location/service/{city,country}.go`

**Проблема:** categories/types/cities/countries теперь кэшируются 15с (`cmd/router/router.go:35` `referenceTTL`),
но Create/Update/Delete **не сбрасывают** кэш → админ-правка видна до 15с.

**Фикс:** внедрить `*cache.Cache` в 4 сервиса, вызывать `InvalidateNamespace(ns)` в Create/Update/Delete.

### P3.4 — `like_count` устаревает в кэшированном списке
`internal/interaction/service/house_like.go:27` `Like`/`Unlike`

**Проблема:** список домов кэшируется 5м, `HouseListItem.LikeCount` в нём. Like/Unlike не инвалидирует `houses`
→ счётчик устаревает до 5 мин.

**Фикс:** либо инвалидировать `houses` из `Like/Unlike` (но убьёт эффективность 5м-кэша под нагрузкой),
либо убрать волатильные счётчики из кэшируемого payload, либо принять staleness (локальный L1 бьёт её до ~2с).

### Заметки по TTL
- Local L1 TTL = 2с (`configs/redis.go:10`), `Set` кэпит любой TTL до 2с (`memory.go:54`) — это shield от
  thundering herd; ограничивает staleness неинвалидируемых путей до ~2с внутри процесса. Оставить.
- Reference response TTL = 15с — ок при отсутствии инвалидации.
- Houses TTL = 5м (дефолт конфига, т.к. route TTL=0) — корректно (мутации инвалидируют). Явный TTL на роут
  `houses` был бы чище, чем полагаться на глобальный fallback 300с.

---

## Рекомендованная быстрая пачка (низкий риск, без миграций)

1. **P1.1** — `ANY($ids)` на лайки (каждый авторизованный список).
2. **P2.1** — параллель `GetBySlug` ∥ `StatusWithCount` (−1 RTT на detail).
3. **P2.2** — `procGroup.SetLimit(GOMAXPROCS)` (одна строка).

Отдельными шагами: **P3.1/P3.2** (кэш detail/stats), затем **P1.3 + P1.4** (`is_active` + композитный индекс, 🐘 миграция),
затем решение по **P1.2** (схема счётчиков).

## Что НЕ трогать (уже хорошо)
`recordViewAsync`, upload-пул, `singleflight`, LATERAL-джоины, локальный L1-кэш, `errgroup` в аналитике/списках.
