ALTER TABLE contacts ADD CONSTRAINT contacts_name_birthday_key UNIQUE (name, birthday);
ALTER TABLE contacts ADD COLUMN created_at TIMESTAMP without time zone;
ALTER TABLE contacts ADD COLUMN updated_at TIMESTAMP without time zone;