CREATE TABLE IF NOT EXISTS  cities (
    id            SERIAL PRIMARY KEY,

    name_ru       VARCHAR(255) NOT NULL,
    name_en       VARCHAR(255) NOT NULL,
    name_kz       VARCHAR(255) NOT NULL,

    postall_code  VARCHAR(255),

    country_id    INTEGER NOT NULL REFERENCES countries(id),

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX product_city_country_id ON cities(country_id);
CREATE INDEX ix_city_id ON cities(id);

