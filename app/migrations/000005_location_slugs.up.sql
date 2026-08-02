ALTER TABLE countries ADD COLUMN slug VARCHAR(255);
ALTER TABLE cities ADD COLUMN slug VARCHAR(255);
UPDATE countries SET slug = lower(regexp_replace(trim(name_en), '[^a-zA-Z0-9]+', '-', 'g')) WHERE slug IS NULL OR slug = '';
UPDATE cities SET slug = lower(regexp_replace(trim(name_en), '[^a-zA-Z0-9]+', '-', 'g')) WHERE slug IS NULL OR slug = '';
CREATE UNIQUE INDEX IF NOT EXISTS countries_slug_key ON countries (slug);
CREATE UNIQUE INDEX IF NOT EXISTS cities_slug_key ON cities (slug);
