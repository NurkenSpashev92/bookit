-- ============================================================
-- Migration 000010: Schema normalization & house_likes
-- ============================================================

-- 1. Create house_likes table (replaces like_count in houses)
CREATE TABLE IF NOT EXISTS house_likes (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    house_id   INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_house_likes_user_house ON house_likes(user_id, house_id);
CREATE INDEX ix_house_likes_house_id ON house_likes(house_id);
CREATE INDEX ix_house_likes_user_id ON house_likes(user_id);

-- 2. Drop like_count from houses (replaced by house_likes)
ALTER TABLE houses DROP CONSTRAINT IF EXISTS chk_like_count;
ALTER TABLE houses DROP COLUMN IF EXISTS like_count;

-- 3. Fix lng/lat types: VARCHAR -> NUMERIC
ALTER TABLE houses ALTER COLUMN lng TYPE NUMERIC USING NULLIF(lng, '')::NUMERIC;
ALTER TABLE houses ALTER COLUMN lat TYPE NUMERIC USING NULLIF(lat, '')::NUMERIC;

-- 4. Fix priority: TEXT -> INTEGER
ALTER TABLE houses ALTER COLUMN priority TYPE INTEGER USING COALESCE(NULLIF(priority, '')::INTEGER, 0);
ALTER TABLE houses ALTER COLUMN priority SET DEFAULT 0;

-- 5. Fix users.is_active default to TRUE (new users should be active)
ALTER TABLE users ALTER COLUMN is_active SET DEFAULT TRUE;

-- 6. Add missing NOT NULL constraints
ALTER TABLE types ALTER COLUMN name SET NOT NULL;
ALTER TABLE categories ALTER COLUMN name_kz SET NOT NULL;
ALTER TABLE categories ALTER COLUMN name_ru SET NOT NULL;
ALTER TABLE categories ALTER COLUMN name_en SET NOT NULL;
ALTER TABLE categories ALTER COLUMN is_active SET NOT NULL;
ALTER TABLE categories ALTER COLUMN is_active SET DEFAULT TRUE;
ALTER TABLE conveniences ALTER COLUMN name SET NOT NULL;
ALTER TABLE conveniences ALTER COLUMN is_active SET NOT NULL;
ALTER TABLE conveniences ALTER COLUMN is_active SET DEFAULT TRUE;

-- 7. Fix phone_number length (international numbers)
ALTER TABLE houses ALTER COLUMN phone_number TYPE VARCHAR(20);

-- 8. Fix postall_code length
ALTER TABLE cities ALTER COLUMN postall_code TYPE VARCHAR(20);

-- 9. Add ON DELETE SET NULL for houses FK (prevent cascade delete of houses when type/city/country deleted)
-- Drop old FK and re-create with proper behavior
ALTER TABLE houses DROP CONSTRAINT IF EXISTS houses_type_id_fkey;
ALTER TABLE houses ADD CONSTRAINT houses_type_id_fkey FOREIGN KEY (type_id) REFERENCES types(id) ON DELETE RESTRICT;

ALTER TABLE houses DROP CONSTRAINT IF EXISTS houses_city_id_fkey;
ALTER TABLE houses ADD CONSTRAINT houses_city_id_fkey FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE SET NULL;

ALTER TABLE houses DROP CONSTRAINT IF EXISTS houses_country_id_fkey;
ALTER TABLE houses ADD CONSTRAINT houses_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE SET NULL;

ALTER TABLE houses DROP CONSTRAINT IF EXISTS houses_owner_id_fkey;
ALTER TABLE houses ADD CONSTRAINT houses_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE;

-- 10. Remove redundant indexes on primary keys (PG already indexes PKs)
DROP INDEX IF EXISTS ix_users_id;
DROP INDEX IF EXISTS ix_types_id;
DROP INDEX IF EXISTS ix_country_id;
DROP INDEX IF EXISTS ix_city_id;
DROP INDEX IF EXISTS ix_images_id;
DROP INDEX IF EXISTS ix_conveniences_id;
DROP INDEX IF EXISTS ix_faq_id;
DROP INDEX IF EXISTS ix_inquiry_id;
DROP INDEX IF EXISTS ix_categories_id;
DROP INDEX IF EXISTS ix_house_convenience_id;
DROP INDEX IF EXISTS ix_house_category_id;

-- 11. Remove duplicate index on houses(type_id)
DROP INDEX IF EXISTS houses_type_id_idx;

-- 12. Add missing NOT NULL + default for boolean fields in houses
ALTER TABLE houses ALTER COLUMN is_active SET DEFAULT TRUE;
ALTER TABLE houses ALTER COLUMN guests_with_pets SET DEFAULT FALSE;
ALTER TABLE houses ALTER COLUMN best_house SET DEFAULT FALSE;
ALTER TABLE houses ALTER COLUMN promotion SET DEFAULT FALSE;
