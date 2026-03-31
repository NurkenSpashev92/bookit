DROP INDEX IF EXISTS users_phone_number_unique;
ALTER TABLE users ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);
