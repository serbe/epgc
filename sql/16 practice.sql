CREATE TABLE IF NOT EXISTS
    practices (
        id               bigserial PRIMARY KEY,
        company_id       bigint,
        kind_id          bigint,
        topic            text,
        date_of_practice time,
        note             text
    );