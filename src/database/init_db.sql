drop table if exists tg_user cascade;
drop table if exists fast_task cascade;
drop table if exists schedule cascade;
drop table if exists schedule_task cascade;

CREATE TABLE if not exists tg_user (
    user_id  integer unique,
    username varchar(255),
    state_code integer,
    state_request text
);

CREATE index on tg_user(user_id);

CREATE TABLE if not exists fast_task (
    id serial primary key,
    assignee_id integer references tg_user(user_id),
    chat_id bigint,
    task_name varchar(255),
    notify_interval bigint,
    deadline timestamptz
);

CREATE TABLE if not exists schedule (
    id serial primary key,
    assignee_id integer references tg_user(user_id),
    week_day varchar(10),
    title varchar(255),
    place varchar(50),
    speaker varchar(50),
    start_time time,
    end_time time
);
