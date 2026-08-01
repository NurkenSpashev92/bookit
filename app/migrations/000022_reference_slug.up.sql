-- Add a unique slug to the reference tables (types, categories, conveniences),
-- backfilled from the english name (name/name_en), de-duplicated by id suffix.

-- types
ALTER TABLE types ADD COLUMN IF NOT EXISTS slug VARCHAR(255);
UPDATE types
SET slug = NULLIF(trim(BOTH '-' FROM regexp_replace(lower(COALESCE(name_en, '')), '[^a-z0-9]+', '-', 'g')), '')
WHERE slug IS NULL OR slug = '';
UPDATE types SET slug = 'type-' || id WHERE slug IS NULL;
WITH d AS (
    SELECT id, row_number() OVER (PARTITION BY slug ORDER BY id) AS rn FROM types
)
UPDATE types t SET slug = t.slug || '-' || t.id FROM d WHERE d.id = t.id AND d.rn > 1;
ALTER TABLE types ALTER COLUMN slug SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_types_slug ON types(slug);

-- categories
ALTER TABLE categories ADD COLUMN IF NOT EXISTS slug VARCHAR(255);
UPDATE categories
SET slug = NULLIF(trim(BOTH '-' FROM regexp_replace(lower(COALESCE(name_en, '')), '[^a-z0-9]+', '-', 'g')), '')
WHERE slug IS NULL OR slug = '';
UPDATE categories SET slug = 'category-' || id WHERE slug IS NULL;
WITH d AS (
    SELECT id, row_number() OVER (PARTITION BY slug ORDER BY id) AS rn FROM categories
)
UPDATE categories c SET slug = c.slug || '-' || c.id FROM d WHERE d.id = c.id AND d.rn > 1;
ALTER TABLE categories ALTER COLUMN slug SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_categories_slug ON categories(slug);

-- conveniences (single name column)
ALTER TABLE conveniences ADD COLUMN IF NOT EXISTS slug VARCHAR(255);
UPDATE conveniences
SET slug = NULLIF(trim(BOTH '-' FROM regexp_replace(lower(COALESCE(name, '')), '[^a-z0-9]+', '-', 'g')), '')
WHERE slug IS NULL OR slug = '';
UPDATE conveniences SET slug = 'convenience-' || id WHERE slug IS NULL;
WITH d AS (
    SELECT id, row_number() OVER (PARTITION BY slug ORDER BY id) AS rn FROM conveniences
)
UPDATE conveniences c SET slug = c.slug || '-' || c.id FROM d WHERE d.id = c.id AND d.rn > 1;
ALTER TABLE conveniences ALTER COLUMN slug SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_conveniences_slug ON conveniences(slug);
