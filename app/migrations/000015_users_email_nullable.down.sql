DROP INDEX IF EXISTS users_email_unique;
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
CREATE UNIQUE INDEX ix_users_email ON users(email);
