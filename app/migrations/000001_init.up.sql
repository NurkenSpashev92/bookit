CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TYPE subscription_type AS ENUM ('basic', 'pro', 'max');
CREATE TYPE subscription_status AS ENUM ('active', 'in_active');

CREATE TABLE users (
    id                SERIAL PRIMARY KEY,
    email             VARCHAR(255),
    first_name        VARCHAR(255),
    last_name         VARCHAR(255),
    middle_name       VARCHAR(255),
    password          VARCHAR(255),
    date_of_birth     DATE,
    phone_number      VARCHAR(128),
    avatar            VARCHAR(255),
    payment_qr        VARCHAR(255),
    payment_phone     VARCHAR(20),
    is_superuser      BOOLEAN NOT NULL DEFAULT FALSE,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    subscription_type subscription_type NOT NULL DEFAULT 'basic',
    date_joined       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX users_email_unique ON users (email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX users_phone_number_unique ON users (phone_number) WHERE phone_number IS NOT NULL;

CREATE TABLE countries (
    id         SERIAL PRIMARY KEY,
    name_kz    VARCHAR(255) NOT NULL,
    name_en    VARCHAR(255) NOT NULL,
    name_ru    VARCHAR(255) NOT NULL,
    code       VARCHAR(10),
    slug       VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX countries_slug_key ON countries (slug);

CREATE TABLE cities (
    id           SERIAL PRIMARY KEY,
    name_ru      VARCHAR(255) NOT NULL,
    name_en      VARCHAR(255) NOT NULL,
    name_kz      VARCHAR(255) NOT NULL,
    postall_code VARCHAR(20),
    slug         VARCHAR(255),
    country_id   INTEGER NOT NULL REFERENCES countries (id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ix_cities_country_id ON cities (country_id);
CREATE UNIQUE INDEX cities_slug_key ON cities (slug);

CREATE TABLE types (
    id         SERIAL PRIMARY KEY,
    name_kz    VARCHAR(255) NOT NULL,
    name_ru    VARCHAR(255) NOT NULL,
    name_en    VARCHAR(255) NOT NULL,
    slug       VARCHAR(255) NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX uq_types_slug ON types (slug);

CREATE TABLE categories (
    id         SERIAL PRIMARY KEY,
    name_kz    VARCHAR(255) NOT NULL,
    name_ru    VARCHAR(255) NOT NULL,
    name_en    VARCHAR(255) NOT NULL,
    slug       VARCHAR(255) NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX uq_categories_slug ON categories (slug);

CREATE TABLE conveniences (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    slug       VARCHAR(255) NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX uq_conveniences_slug ON conveniences (slug);

CREATE TABLE houses (
    id                 SERIAL PRIMARY KEY,
    name_en            VARCHAR(255) NOT NULL,
    name_kz            VARCHAR(255) NOT NULL,
    name_ru            VARCHAR(255) NOT NULL,
    slug               VARCHAR(255) NOT NULL UNIQUE,
    price              INTEGER NOT NULL DEFAULT 0,
    rooms_qty          INTEGER NOT NULL DEFAULT 0,
    guest_qty          INTEGER NOT NULL DEFAULT 0,
    bedroom_qty        INTEGER NOT NULL DEFAULT 0,
    bath_qty           INTEGER DEFAULT 0,
    description_en     TEXT NOT NULL,
    description_kz     TEXT NOT NULL,
    description_ru     TEXT NOT NULL,
    address_en         VARCHAR(255) NOT NULL,
    address_kz         VARCHAR(255) NOT NULL,
    address_ru         VARCHAR(255) NOT NULL,
    lng                NUMERIC,
    lat                NUMERIC,
    is_active          BOOLEAN NOT NULL DEFAULT TRUE,
    priority           INTEGER NOT NULL DEFAULT 0,
    comments_ru        TEXT,
    comments_en        TEXT,
    comments_kz        TEXT,
    owner_id           INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type_id            INTEGER NOT NULL REFERENCES types (id) ON DELETE RESTRICT,
    city_id            INTEGER REFERENCES cities (id) ON DELETE SET NULL,
    country_id         INTEGER REFERENCES countries (id) ON DELETE SET NULL,
    guests_with_pets   BOOLEAN NOT NULL DEFAULT FALSE,
    guests_with_babies BOOLEAN NOT NULL DEFAULT FALSE,
    best_house         BOOLEAN NOT NULL DEFAULT FALSE,
    promotion          BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified        BOOLEAN NOT NULL DEFAULT FALSE,
    is_sale            BOOLEAN NOT NULL DEFAULT FALSE,
    is_newest          BOOLEAN NOT NULL DEFAULT TRUE,
    is_hot             BOOLEAN NOT NULL DEFAULT FALSE,
    is_featured        BOOLEAN NOT NULL DEFAULT FALSE,
    is_discount        BOOLEAN NOT NULL DEFAULT FALSE,
    district_en        VARCHAR(255),
    district_kz        VARCHAR(255),
    district_ru        VARCHAR(255),
    phone_number       VARCHAR(20),
    like_count         INTEGER NOT NULL DEFAULT 0,
    view_count         INTEGER NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_price CHECK (price >= 0),
    CONSTRAINT chk_rooms_qty CHECK (rooms_qty >= 0),
    CONSTRAINT chk_guest_qty CHECK (guest_qty >= 0),
    CONSTRAINT chk_bedroom_qty CHECK (bedroom_qty >= 0),
    CONSTRAINT chk_bath_qty CHECK (bath_qty >= 0),
    CONSTRAINT chk_like_count CHECK (like_count >= 0),
    CONSTRAINT chk_view_count CHECK (view_count >= 0)
);
CREATE INDEX ix_houses_owner_id ON houses (owner_id);
CREATE INDEX ix_houses_type_id ON houses (type_id);
CREATE INDEX ix_houses_city_id ON houses (city_id);
CREATE INDEX ix_houses_country_id ON houses (country_id);
CREATE INDEX ix_houses_price ON houses (price);
CREATE INDEX ix_houses_guest_qty ON houses (guest_qty);
CREATE INDEX ix_houses_rooms_qty ON houses (rooms_qty);
CREATE INDEX ix_houses_active ON houses (id DESC) WHERE is_active = TRUE;
CREATE INDEX ix_houses_inactive ON houses (id DESC) WHERE is_active = FALSE;
CREATE INDEX ix_houses_active_city ON houses (city_id, id DESC) WHERE is_active = TRUE;
CREATE INDEX ix_houses_best ON houses (id DESC) WHERE best_house = TRUE;
CREATE INDEX ix_houses_promo ON houses (id DESC) WHERE promotion = TRUE;
CREATE INDEX ix_houses_is_verified ON houses (id DESC) WHERE is_verified;
CREATE INDEX ix_houses_is_hot ON houses (id DESC) WHERE is_hot;
CREATE INDEX ix_houses_is_featured ON houses (id DESC) WHERE is_featured;
CREATE INDEX ix_houses_name_en_trgm ON houses USING gin (name_en gin_trgm_ops);
CREATE INDEX ix_houses_name_kz_trgm ON houses USING gin (name_kz gin_trgm_ops);
CREATE INDEX ix_houses_name_ru_trgm ON houses USING gin (name_ru gin_trgm_ops);

CREATE TABLE images (
    id         SERIAL PRIMARY KEY,
    original   VARCHAR(255),
    thumbnail  VARCHAR(255),
    width      INTEGER,
    height     INTEGER,
    mimetype   VARCHAR(100),
    size       INTEGER,
    is_label   BOOLEAN,
    house_id   INTEGER REFERENCES houses (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ix_images_house_id_id ON images (house_id, id);

CREATE TABLE house_category (
    id          SERIAL PRIMARY KEY,
    house_id    INTEGER NOT NULL REFERENCES houses (id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    UNIQUE (house_id, category_id)
);
CREATE INDEX ix_house_category_category_id ON house_category (category_id);

CREATE TABLE house_convenience (
    id             SERIAL PRIMARY KEY,
    house_id       INTEGER NOT NULL REFERENCES houses (id) ON DELETE CASCADE,
    convenience_id INTEGER NOT NULL REFERENCES conveniences (id) ON DELETE CASCADE,
    UNIQUE (house_id, convenience_id)
);
CREATE INDEX ix_house_convenience_convenience_id ON house_convenience (convenience_id);

CREATE TABLE house_likes (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    house_id   INTEGER NOT NULL REFERENCES houses (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_house_likes_user_house UNIQUE (user_id, house_id)
);
CREATE INDEX ix_house_likes_house_user ON house_likes (house_id, user_id);
CREATE INDEX ix_house_likes_user_created ON house_likes (user_id, created_at DESC);

CREATE TABLE house_views (
    id         SERIAL PRIMARY KEY,
    house_id   INTEGER NOT NULL REFERENCES houses (id) ON DELETE CASCADE,
    user_id    INTEGER REFERENCES users (id) ON DELETE SET NULL,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ix_house_views_house_created ON house_views (house_id, created_at DESC);
CREATE INDEX ix_house_views_user_id ON house_views (user_id) WHERE user_id IS NOT NULL;

CREATE TABLE bookings (
    id          SERIAL PRIMARY KEY,
    house_id    INTEGER NOT NULL REFERENCES houses (id) ON DELETE CASCADE,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    start_date  DATE NOT NULL,
    end_date    DATE NOT NULL,
    guest_count INTEGER NOT NULL DEFAULT 1,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_price INTEGER NOT NULL DEFAULT 0,
    message     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_bookings_dates CHECK (end_date >= start_date),
    CONSTRAINT chk_bookings_guests CHECK (guest_count >= 1),
    CONSTRAINT chk_bookings_price CHECK (total_price >= 0)
);
CREATE INDEX ix_bookings_house_status ON bookings (house_id, status);
CREATE INDEX ix_bookings_user_id ON bookings (user_id);

CREATE TABLE faq (
    id          SERIAL PRIMARY KEY,
    question_kz VARCHAR(500),
    answer_kz   TEXT,
    question_ru VARCHAR(500),
    answer_ru   TEXT,
    question_en VARCHAR(500),
    answer_en   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE inquiries (
    id           SERIAL PRIMARY KEY,
    email        VARCHAR(255) NOT NULL,
    phone_number VARCHAR(12),
    text         TEXT NOT NULL,
    is_approved  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE subscriptions (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type       subscription_type NOT NULL DEFAULT 'basic',
    status     subscription_status NOT NULL DEFAULT 'active',
    start_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_date   TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_subscriptions_dates CHECK (end_date IS NULL OR end_date >= start_date)
);
CREATE INDEX ix_subscriptions_user_id ON subscriptions (user_id, id DESC);
CREATE UNIQUE INDEX uq_subscriptions_active_user ON subscriptions (user_id) WHERE status = 'active';

CREATE FUNCTION trg_house_likes_inc() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    UPDATE houses SET like_count = like_count + 1 WHERE id = NEW.house_id;
    RETURN NEW;
END;
$$;
CREATE FUNCTION trg_house_likes_dec() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    UPDATE houses SET like_count = GREATEST(like_count - 1, 0) WHERE id = OLD.house_id;
    RETURN OLD;
END;
$$;
CREATE TRIGGER house_likes_after_insert AFTER INSERT ON house_likes
    FOR EACH ROW EXECUTE FUNCTION trg_house_likes_inc();
CREATE TRIGGER house_likes_after_delete AFTER DELETE ON house_likes
    FOR EACH ROW EXECUTE FUNCTION trg_house_likes_dec();
