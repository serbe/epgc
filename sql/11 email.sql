CREATE TABLE IF NOT EXISTS
    emails (
        id         bigserial PRIMARY KEY,
        company_id bigint,
        people_id  bigint,
        email      text,
        note       text
    );