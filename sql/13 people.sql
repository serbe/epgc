CREATE TABLE IF NOT EXISTS
    peoples (
        id         bigserial PRIMARY KEY,
        name       text,
        company_id bigint,
        post_id    bigint,
        post_go_id bigint,
        rank_id    bigint,
        birthday   date,
        note       text
    );