CREATE TABLE IF NOT EXISTS
    sirenTypes (
        id         bigserial primary key,
        name       text,
        radius     bigint,
        note       text,
        created_at TIMESTAMP without time zone,
        updated_at TIMESTAMP without time zone,
        UNIQUE(name, radius)
    );