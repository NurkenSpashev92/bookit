CREATE TABLE IF NOT EXISTS faq (
    id       SERIAL PRIMARY KEY,
    question_kz VARCHAR(500),
    answer_kz   TEXT,
    question_ru VARCHAR(500),
    answer_ru   TEXT,
    question_en VARCHAR(500),
    answer_en   TEXT,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_faq_id ON faq(id);

CREATE TABLE IF NOT EXISTS inquiries (
    id           SERIAL PRIMARY KEY,
    email        VARCHAR(255) NOT NULL,
    phone_number VARCHAR(12),
    text         TEXT NOT NULL,
    is_approved  BOOLEAN NOT NULL DEFAULT FALSE,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_inquiry_id ON inquiries(id);
