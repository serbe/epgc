CREATE TABLE IF NOT EXISTS
    posts (
        id   bigserial PRIMARY KEY,
        name text,
        go   bool NOT NULL DEFAULT FALSE,
        note text,
        UNIQUE (name, go)
    );