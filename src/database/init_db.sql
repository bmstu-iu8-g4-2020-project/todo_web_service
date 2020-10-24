drop table if exists tg_user cascade;

CREATE TABLE if not exists task_table
(
    id       SERIAL PRIMARY KEY UNIQUE,
    username varchar(255),
    user_id  integer
);
