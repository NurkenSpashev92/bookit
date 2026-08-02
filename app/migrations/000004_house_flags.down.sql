DROP INDEX IF EXISTS ix_houses_is_featured;
DROP INDEX IF EXISTS ix_houses_is_hot;
DROP INDEX IF EXISTS ix_houses_is_verified;

ALTER TABLE houses
    DROP COLUMN IF EXISTS is_discount,
    DROP COLUMN IF EXISTS is_featured,
    DROP COLUMN IF EXISTS is_hot,
    DROP COLUMN IF EXISTS is_newest,
    DROP COLUMN IF EXISTS is_sale,
    DROP COLUMN IF EXISTS is_verified;
