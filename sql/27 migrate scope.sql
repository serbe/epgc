ALTER TABLE scopes ADD CONSTRAINT scopes_name_key UNIQUE (name);
ALTER TABLE scopes ADD COLUMN created_at TIMESTAMP without time zone;
ALTER TABLE scopes ADD COLUMN updated_at TIMESTAMP without time zone;