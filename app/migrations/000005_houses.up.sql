CREATE TABLE IF NOT EXISTS houses (
    id              SERIAL PRIMARY KEY,

    name_en         VARCHAR(255) NOT NULL,
    name_kz         VARCHAR(255) NOT NULL,
    name_ru         VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL,

    price INTEGER NOT NULL DEFAULT 0,

    rooms_qty       INTEGER NOT NULL DEFAULT 0,
    guest_qty       INTEGER NOT NULL DEFAULT 0,
    bedroom_qty     INTEGER NOT NULL DEFAULT 0,
    bath_qty        INTEGER DEFAULT 0,

    description_en  TEXT NOT NULL,
    description_kz  TEXT NOT NULL,
    description_ru  TEXT NOT NULL,

    address_en      VARCHAR(255) NOT NULL,
    address_kz      VARCHAR(255) NOT NULL,
    address_ru      VARCHAR(255) NOT NULL,

    lng             VARCHAR(255),
    lat             VARCHAR(255),

    is_active       BOOLEAN NOT NULL,
    priority        TEXT NOT NULL,
    like_count      INTEGER NOT NULL DEFAULT 0,

    comments_ru     TEXT,
    comments_en     TEXT,
    comments_kz     TEXT,

    owner_id        INTEGER NOT NULL REFERENCES users(id),
    type_id         INTEGER NOT NULL REFERENCES types(id),
    city_id         INTEGER REFERENCES cities(id),
    country_id      INTEGER REFERENCES countries(id),

    guests_with_pets   BOOLEAN NOT NULL,
    best_house       BOOLEAN NOT NULL,
    promotion          BOOLEAN NOT NULL,

    district_en     VARCHAR(255),
    district_kz     VARCHAR(255),
    district_ru     VARCHAR(255),

    phone_number    VARCHAR(12),

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- UNIQUE
    CONSTRAINT houses_slug_key UNIQUE (slug),

    -- CHECK CONSTRAINTS
    CONSTRAINT chk_bath_qty          CHECK (bath_qty >= 0),
    CONSTRAINT chk_bedroom_qty       CHECK (bedroom_qty >= 0),
    CONSTRAINT chk_guest_qty         CHECK (guest_qty >= 0),
    CONSTRAINT chk_rooms_qty         CHECK (rooms_qty >= 0),
    CONSTRAINT chk_like_count        CHECK (like_count >= 0),
    CONSTRAINT chk_price             CHECK (price >= 0)
);

CREATE INDEX house_house_owner_id ON houses(owner_id);
CREATE INDEX house_slug ON houses(slug);
CREATE INDEX house_house_type_id ON houses(type_id);
CREATE INDEX houses_city_id ON houses(city_id);
CREATE INDEX houses_country_id ON houses(country_id);
CREATE INDEX houses_guest_q_idx ON houses(guest_qty);
CREATE INDEX houses_name_en_idx ON houses(name_en);
CREATE INDEX houses_name_kz_idx ON houses(name_kz);
CREATE INDEX houses_name_ru_idx ON houses(name_ru);
CREATE INDEX houses_price_p_idx ON houses(price);
CREATE INDEX houses_rooms_q_idx ON houses(rooms_qty);
CREATE INDEX houses_type_id_idx ON houses(type_id);
