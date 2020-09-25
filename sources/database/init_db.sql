CREATE TABLE if not exists task_table
(
    id                 SERIAL PRIMARY KEY UNIQUE,
    title              varchar(64),
    description        varchar(128),
    notifications_freq int,
    is_accomplished    boolean,
    creation_time      timestamp without time zone
);

CREATE TABLE if not exists user_table
(
    id           SERIAL PRIMARY KEY UNIQUE,
    name         varchar(32),
    login        varchar(32) UNIQUE,
    password     varchar(128),
    phone_number varchar(32) UNIQUE
);

CREATE TABLE if not exists shedule_table
(
--?
);