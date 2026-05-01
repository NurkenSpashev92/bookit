-- ============================================================
-- Migration 000016: Trigram index for ILIKE '%text%' name search
-- ============================================================
-- Plain B-tree indexes can't accelerate substring/wildcard search
-- (ILIKE '%foo%'). pg_trgm + GIN make it fast.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS ix_houses_name_en_trgm
    ON houses USING GIN (name_en gin_trgm_ops);

CREATE INDEX IF NOT EXISTS ix_houses_name_kz_trgm
    ON houses USING GIN (name_kz gin_trgm_ops);

CREATE INDEX IF NOT EXISTS ix_houses_name_ru_trgm
    ON houses USING GIN (name_ru gin_trgm_ops);
