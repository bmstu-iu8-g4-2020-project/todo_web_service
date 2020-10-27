drop table if exists tg_user cascade;

CREATE TABLE if not exists tg_user
(
    user_id  integer,
    username varchar(255)
);
