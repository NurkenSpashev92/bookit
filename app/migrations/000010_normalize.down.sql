-- Revert migration 000010

-- Restore like_count
ALTER TABLE houses ADD COLUMN IF NOT EXISTS like_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE houses ADD CONSTRAINT chk_like_count CHECK (like_count >= 0);

-- Drop house_likes
DROP TABLE IF EXISTS house_likes;

-- Revert lng/lat to VARCHAR
ALTER TABLE houses ALTER COLUMN lng TYPE VARCHAR(255) USING lng::VARCHAR;
ALTER TABLE houses ALTER COLUMN lat TYPE VARCHAR(255) USING lat::VARCHAR;

-- Revert priority to TEXT
ALTER TABLE houses ALTER COLUMN priority TYPE TEXT USING priority::TEXT;
ALTER TABLE houses ALTER COLUMN priority DROP DEFAULT;

-- Revert users.is_active default
ALTER TABLE users ALTER COLUMN is_active SET DEFAULT FALSE;

-- Revert NOT NULL constraints
ALTER TABLE types ALTER COLUMN name DROP NOT NULL;
ALTER TABLE categories ALTER COLUMN name_kz DROP NOT NULL;
ALTER TABLE categories ALTER COLUMN name_ru DROP NOT NULL;
ALTER TABLE categories ALTER COLUMN name_en DROP NOT NULL;
ALTER TABLE categories ALTER COLUMN is_active DROP NOT NULL;
ALTER TABLE categories ALTER COLUMN is_active DROP DEFAULT;
ALTER TABLE conveniences ALTER COLUMN name DROP NOT NULL;
ALTER TABLE conveniences ALTER COLUMN is_active DROP NOT NULL;
ALTER TABLE conveniences ALTER COLUMN is_active DROP DEFAULT;

-- Revert phone_number / postall_code
ALTER TABLE houses ALTER COLUMN phone_number TYPE VARCHAR(12);
ALTER TABLE cities ALTER COLUMN postall_code TYPE VARCHAR(255);

-- Revert FK constraints
ALTER TABLE houses DROP CONSTRAINT IF EXISTS houses_type_id_fkey;
ALTER TABLE houses ADD CONSTRAINT houses_type_id_fkey FOREIGN KEY (type_id) REFERENCES types(id);

ALTER TABLE houses DROP CONSTRAINT IF EXISTS houses_city_id_fkey;
ALTER TABLE houses ADD CONSTRAINT houses_city_id_fkey FOREIGN KEY (city_id) REFERENCES cities(id);

ALTER TABLE houses DROP CONSTRAINT IF EXISTS houses_country_id_fkey;
ALTER TABLE houses ADD CONSTRAINT houses_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries(id);

ALTER TABLE houses DROP CONSTRAINT IF EXISTS houses_owner_id_fkey;
ALTER TABLE houses ADD CONSTRAINT houses_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id);

-- Re-create redundant PK indexes
CREATE INDEX IF NOT EXISTS ix_users_id ON users(id);
CREATE INDEX IF NOT EXISTS ix_types_id ON types(id);
CREATE INDEX IF NOT EXISTS ix_country_id ON countries(id);
CREATE INDEX IF NOT EXISTS ix_city_id ON cities(id);
CREATE INDEX IF NOT EXISTS ix_images_id ON images(id);
CREATE INDEX IF NOT EXISTS ix_conveniences_id ON conveniences(id);
CREATE INDEX IF NOT EXISTS ix_faq_id ON faq(id);
CREATE INDEX IF NOT EXISTS ix_inquiry_id ON inquiries(id);
CREATE INDEX IF NOT EXISTS ix_categories_id ON categories(id);
CREATE INDEX IF NOT EXISTS ix_house_convenience_id ON house_convenience(id);
CREATE INDEX IF NOT EXISTS ix_house_category_id ON house_category(id);

-- Re-create duplicate index
CREATE INDEX IF NOT EXISTS houses_type_id_idx ON houses(type_id);

-- Revert boolean defaults
ALTER TABLE houses ALTER COLUMN is_active DROP DEFAULT;
ALTER TABLE houses ALTER COLUMN guests_with_pets DROP DEFAULT;
ALTER TABLE houses ALTER COLUMN best_house DROP DEFAULT;
ALTER TABLE houses ALTER COLUMN promotion DROP DEFAULT;
