-- Drop the unique constraint and index on email, allow NULLs.
-- Re-add a partial unique index that permits multiple NULLs.
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
DROP INDEX IF EXISTS ix_users_email;
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON users(email) WHERE email IS NOT NULL;
