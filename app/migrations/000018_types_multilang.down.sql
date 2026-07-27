ALTER TABLE types ADD COLUMN IF NOT EXISTS name VARCHAR(255);

UPDATE types
SET name = COALESCE(name, name_ru, name_en, name_kz);

ALTER TABLE types ALTER COLUMN name SET NOT NULL;

ALTER TABLE types
    DROP COLUMN IF EXISTS name_kz,
    DROP COLUMN IF EXISTS name_ru,
    DROP COLUMN IF EXISTS name_en;
