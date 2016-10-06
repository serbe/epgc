CREATE TABLE IF NOT EXISTS
    phones (
        id         bigserial PRIMARY KEY,
        people_id  bigint,
        company_id bigint,
        phone      bigint,
        fax        bool NOT NULL DEFAULT false,
        note       text
    );