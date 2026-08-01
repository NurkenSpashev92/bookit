DROP TABLE IF EXISTS
    subscriptions,
    inquiries,
    faq,
    bookings,
    house_views,
    house_likes,
    house_convenience,
    house_category,
    images,
    houses,
    conveniences,
    categories,
    types,
    cities,
    countries,
    users
CASCADE;

DROP FUNCTION IF EXISTS trg_house_likes_inc();
DROP FUNCTION IF EXISTS trg_house_likes_dec();

DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS subscription_type;
