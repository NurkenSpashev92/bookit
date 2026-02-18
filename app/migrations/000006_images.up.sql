CREATE TABLE IF NOT EXISTS images (
    id          SERIAL PRIMARY KEY,

    original    VARCHAR(255),
    thumbnail   VARCHAR(255),
    width       INTEGER,
    height      INTEGER,
    mimetype    VARCHAR(100),
    size        INTEGER,
    is_label    BOOLEAN,

    house_id  INTEGER REFERENCES houses(id) ON DELETE CASCADE,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_images_id ON images(id);
CREATE INDEX ix_images_house_id ON images(house_id);
