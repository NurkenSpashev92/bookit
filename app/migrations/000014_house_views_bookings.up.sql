-- Track individual house views
CREATE TABLE house_views (
    id          SERIAL PRIMARY KEY,
    house_id    INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    user_id     INTEGER REFERENCES users(id) ON DELETE SET NULL,
    ip_address  VARCHAR(45),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_house_views_house_id ON house_views(house_id);
CREATE INDEX ix_house_views_user_id ON house_views(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX ix_house_views_created_at ON house_views(house_id, created_at DESC);

-- Bookings table
CREATE TABLE bookings (
    id          SERIAL PRIMARY KEY,
    house_id    INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date  DATE NOT NULL,
    end_date    DATE NOT NULL,
    guest_count INTEGER NOT NULL DEFAULT 1,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_price INTEGER NOT NULL DEFAULT 0,
    message     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_bookings_house_id ON bookings(house_id);
CREATE INDEX ix_bookings_user_id ON bookings(user_id);
CREATE INDEX ix_bookings_status ON bookings(house_id, status);
