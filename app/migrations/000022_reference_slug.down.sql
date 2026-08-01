DROP INDEX IF EXISTS uq_types_slug;
ALTER TABLE types DROP COLUMN IF EXISTS slug;

DROP INDEX IF EXISTS uq_categories_slug;
ALTER TABLE categories DROP COLUMN IF EXISTS slug;

DROP INDEX IF EXISTS uq_conveniences_slug;
ALTER TABLE conveniences DROP COLUMN IF EXISTS slug;
