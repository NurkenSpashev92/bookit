CREATE TABLE IF NOT EXISTS  users (
    id              SERIAL PRIMARY KEY,

    email           VARCHAR(255) UNIQUE,
    first_name      VARCHAR(255),
    last_name       VARCHAR(255),
    middle_name     VARCHAR(255),

    password        VARCHAR(255),

    date_of_birth   DATE,

    phone_number    VARCHAR(128) UNIQUE,
    avatar          VARCHAR(255),

    is_superuser    BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT FALSE,

    date_joined     TIMESTAMP NOT NULL DEFAULT NOW(),

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_users_id ON users(id);
CREATE UNIQUE INDEX ix_users_email ON users(email);
