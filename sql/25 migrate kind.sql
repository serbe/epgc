ALTER TABLE kinds ADD CONSTRAINT kinds_name_key UNIQUE (name);
ALTER TABLE kinds ADD COLUMN created_at TIMESTAMP without time zone;
ALTER TABLE kinds ADD COLUMN updated_at TIMESTAMP without time zone;