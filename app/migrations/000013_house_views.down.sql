ALTER TABLE houses DROP COLUMN IF EXISTS view_count;
DROP INDEX IF EXISTS ix_houses_owner_id;
