-- ============================================================
-- Migration 000011: Performance optimizations
-- ============================================================

-- 1. Denormalize like_count back to houses for fast reads
--    Maintained by triggers on house_likes INSERT/DELETE
ALTER TABLE houses ADD COLUMN IF NOT EXISTS like_count INTEGER NOT NULL DEFAULT 0;

-- Backfill existing like counts
UPDATE houses h
SET like_count = (SELECT COUNT(*) FROM house_likes WHERE house_id = h.id);

-- 2. Trigger to increment like_count on INSERT
CREATE OR REPLACE FUNCTION trg_house_likes_inc() RETURNS TRIGGER AS $$
BEGIN
    UPDATE houses SET like_count = like_count + 1 WHERE id = NEW.house_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER house_likes_after_insert
    AFTER INSERT ON house_likes
    FOR EACH ROW EXECUTE FUNCTION trg_house_likes_inc();

-- 3. Trigger to decrement like_count on DELETE
CREATE OR REPLACE FUNCTION trg_house_likes_dec() RETURNS TRIGGER AS $$
BEGIN
    UPDATE houses SET like_count = GREATEST(like_count - 1, 0) WHERE id = OLD.house_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER house_likes_after_delete
    AFTER DELETE ON house_likes
    FOR EACH ROW EXECUTE FUNCTION trg_house_likes_dec();

-- 4. Composite index for images by house (used in LATERAL subquery)
CREATE INDEX IF NOT EXISTS ix_images_house_id_id ON images(house_id, id);

-- 5. Composite index for house_likes count queries
-- (house_id already indexed, but covering index with user_id is useful)
DROP INDEX IF EXISTS ix_house_likes_house_id;
CREATE INDEX IF NOT EXISTS ix_house_likes_house_user ON house_likes(house_id, user_id);

-- 6. Partial index for active houses (most list queries filter active)
CREATE INDEX IF NOT EXISTS ix_houses_active ON houses(id DESC) WHERE is_active = TRUE;

-- 7. Index for best_house / promotion filters
CREATE INDEX IF NOT EXISTS ix_houses_best ON houses(id DESC) WHERE best_house = TRUE;
CREATE INDEX IF NOT EXISTS ix_houses_promo ON houses(id DESC) WHERE promotion = TRUE;

-- 8. Composite index for liked houses query (user's liked houses ordered by time)
DROP INDEX IF EXISTS ix_house_likes_user_id;
CREATE INDEX IF NOT EXISTS ix_house_likes_user_created ON house_likes(user_id, created_at DESC);
