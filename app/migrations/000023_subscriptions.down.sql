ALTER TABLE users DROP COLUMN IF EXISTS subscription_type;

DROP TABLE IF EXISTS subscriptions;

DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS subscription_type;
