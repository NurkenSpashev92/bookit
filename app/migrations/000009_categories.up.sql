CREATE TABLE IF NOT EXISTS categories (
    id        SERIAL PRIMARY KEY,
    name_kz      VARCHAR(255),
    name_ru      VARCHAR(255),
    name_en      VARCHAR(255),
    is_active BOOLEAN,
    icon      VARCHAR(255),

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_categories_id ON categories(id);

CREATE TABLE IF NOT EXISTS house_category (
    id          SERIAL PRIMARY KEY,
    house_id  INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX ix_house_category_id ON house_category(id);
CREATE INDEX ix_house_category_house_id ON house_category(house_id);
CREATE INDEX ix_house_category_category_id ON house_category(category_id);

CREATE UNIQUE INDEX uq_house_category_house_id_category_id
ON house_category(house_id, category_id);
