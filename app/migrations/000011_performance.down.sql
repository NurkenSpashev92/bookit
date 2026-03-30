-- Rollback performance optimizations

DROP TRIGGER IF EXISTS house_likes_after_insert ON house_likes;
DROP TRIGGER IF EXISTS house_likes_after_delete ON house_likes;
DROP FUNCTION IF EXISTS trg_house_likes_inc();
DROP FUNCTION IF EXISTS trg_house_likes_dec();

ALTER TABLE houses DROP COLUMN IF EXISTS like_count;

DROP INDEX IF EXISTS ix_images_house_id_id;
DROP INDEX IF EXISTS ix_house_likes_house_user;
DROP INDEX IF EXISTS ix_houses_active;
DROP INDEX IF EXISTS ix_houses_best;
DROP INDEX IF EXISTS ix_houses_promo;
DROP INDEX IF EXISTS ix_house_likes_user_created;

-- Restore original indexes
CREATE INDEX IF NOT EXISTS ix_house_likes_house_id ON house_likes(house_id);
CREATE INDEX IF NOT EXISTS ix_house_likes_user_id ON house_likes(user_id);
