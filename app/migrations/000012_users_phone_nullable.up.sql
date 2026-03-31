-- Drop the unique constraint on phone_number and allow NULLs.
-- Re-add a unique index that permits multiple NULLs.
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_phone_number_key;
CREATE UNIQUE INDEX IF NOT EXISTS users_phone_number_unique ON users(phone_number) WHERE phone_number IS NOT NULL;
