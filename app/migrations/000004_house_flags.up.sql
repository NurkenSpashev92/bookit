ALTER TABLE houses
    ADD COLUMN is_verified BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN is_sale     BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN is_newest    BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN is_hot       BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN is_featured  BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN is_discount  BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX ix_houses_is_verified ON houses (id DESC) WHERE is_verified;
CREATE INDEX ix_houses_is_hot ON houses (id DESC) WHERE is_hot;
CREATE INDEX ix_houses_is_featured ON houses (id DESC) WHERE is_featured;
