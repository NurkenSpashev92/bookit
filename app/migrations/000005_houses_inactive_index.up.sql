CREATE INDEX ix_houses_inactive ON houses (id DESC) WHERE is_active = FALSE;
