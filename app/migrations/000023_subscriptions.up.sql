CREATE TYPE subscription_type AS ENUM ('basic', 'pro', 'max');
CREATE TYPE subscription_status AS ENUM ('active', 'in_active');

CREATE TABLE IF NOT EXISTS subscriptions (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       subscription_type NOT NULL DEFAULT 'basic',
    status     subscription_status NOT NULL DEFAULT 'active',
    start_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_date   TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_subscriptions_user_id ON subscriptions(user_id);

ALTER TABLE users ADD COLUMN IF NOT EXISTS subscription_type subscription_type NOT NULL DEFAULT 'basic';
