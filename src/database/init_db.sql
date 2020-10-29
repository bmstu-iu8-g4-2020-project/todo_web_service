drop table if exists tg_user cascade;
drop table if exists fast_task cascade;

CREATE TABLE if not exists tg_user (
    user_id  integer,
    username varchar(255)
);

CREATE TABLE if not exists fast_task (
    id serial primary key,
    assignee_id integer references tg_user(user_id),
    chat_id bigint,
    task_name varchar(255),
    notify_interval bigint,
    deadline timestamp
);
