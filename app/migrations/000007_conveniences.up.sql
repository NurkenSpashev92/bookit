CREATE TABLE IF NOT EXISTS conveniences (
    id         SERIAL PRIMARY KEY,

    name       VARCHAR(255),
    is_active  BOOLEAN,
    icon       VARCHAR(100),

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_conveniences_id ON conveniences(id);


CREATE TABLE IF NOT EXISTS house_convenience (
    id              SERIAL PRIMARY KEY,
    house_id      INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    convenience_id  INTEGER NOT NULL REFERENCES conveniences(id) ON DELETE CASCADE
);

CREATE INDEX ix_house_convenience_id ON house_convenience(id);
CREATE INDEX ix_house_convenience_house_id ON house_convenience(house_id);
CREATE INDEX ix_house_convenience_convenience_id ON house_convenience(convenience_id);

CREATE UNIQUE INDEX uq_house_convenience ON house_convenience(house_id, convenience_id);
