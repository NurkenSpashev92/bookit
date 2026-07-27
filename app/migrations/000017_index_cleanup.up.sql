-- ============================================================
-- Migration 000016: чистка избыточных индексов на houses
--
-- Каждый лишний индекс — это дополнительная запись при каждом INSERT/UPDATE
-- дома и лишние страницы в кеше БД, при нулевой пользе на чтении.
-- ============================================================

-- 1. Дубликат: owner_id проиндексирован дважды.
--    Оставляем ix_houses_owner_id из миграции 000010.
DROP INDEX IF EXISTS house_house_owner_id;

-- 2. Дубликат: slug уже покрыт уникальным индексом houses_slug_key,
--    отдельный btree ничего не добавляет.
DROP INDEX IF EXISTS house_slug;

-- 3. btree по name_* не используется: поиск идёт через ILIKE '%...%',
--    который обслуживают gin-индексы ix_houses_name_*_trgm (миграция 000010),
--    а сортировки по имени в запросах нет.
DROP INDEX IF EXISTS houses_name_en_idx;
DROP INDEX IF EXISTS houses_name_kz_idx;
DROP INDEX IF EXISTS houses_name_ru_idx;
