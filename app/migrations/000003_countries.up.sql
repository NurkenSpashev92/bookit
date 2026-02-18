CREATE TABLE IF NOT EXISTS  countries (
    id         SERIAL PRIMARY KEY,

    name_kz    VARCHAR(255) NOT NULL,
    name_en    VARCHAR(255) NOT NULL,
    name_ru    VARCHAR(255) NOT NULL,

    code       VARCHAR(10),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_country_id ON countries(id);