ALTER TABLE types
    ADD COLUMN IF NOT EXISTS name_kz VARCHAR(255),
    ADD COLUMN IF NOT EXISTS name_ru VARCHAR(255),
    ADD COLUMN IF NOT EXISTS name_en VARCHAR(255);

UPDATE types
SET name_kz = COALESCE(name_kz, name),
    name_ru = COALESCE(name_ru, name),
    name_en = COALESCE(name_en, name);

ALTER TABLE types
    ALTER COLUMN name_kz SET NOT NULL,
    ALTER COLUMN name_ru SET NOT NULL,
    ALTER COLUMN name_en SET NOT NULL;

ALTER TABLE types DROP COLUMN IF EXISTS name;
