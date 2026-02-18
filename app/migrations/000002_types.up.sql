CREATE TABLE IF NOT EXISTS  types (
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(255),
    icon      VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_types_id ON types(id);
