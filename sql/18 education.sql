CREATE TABLE IF NOT EXISTS
    educations (
        id         bigserial PRIMARY KEY,
        start_date time,
        end_date   time,
        note       text
    );