ALTER TABLE houses ADD COLUMN view_count INTEGER NOT NULL DEFAULT 0;
CREATE INDEX ix_houses_owner_id ON houses(owner_id);
