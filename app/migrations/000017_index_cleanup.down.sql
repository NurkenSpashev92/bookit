-- Возврат индексов в состояние до миграции 000016.
CREATE INDEX IF NOT EXISTS house_house_owner_id ON houses(owner_id);
CREATE INDEX IF NOT EXISTS house_slug ON houses(slug);
CREATE INDEX IF NOT EXISTS houses_name_en_idx ON houses(name_en);
CREATE INDEX IF NOT EXISTS houses_name_kz_idx ON houses(name_kz);
CREATE INDEX IF NOT EXISTS houses_name_ru_idx ON houses(name_ru);
