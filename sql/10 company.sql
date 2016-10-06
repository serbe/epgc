CREATE TABLE IF NOT EXISTS 
    companies (
        id       bigserial PRIMARY KEY,
        name     text,
        address  text,
        scope_id bigint,
        note     text,
        UNIQUE(name, scope_id)
    );