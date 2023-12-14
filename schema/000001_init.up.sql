CREATE TABLE tasks
(
    id          serial                                      not null unique,
    created_at  timestamp                                   not null default CURRENT_TIMESTAMP,
    updated_at  timestamp                                   not null default CURRENT_TIMESTAMP,
    telegram_id   varchar(20)  not null,
    text        text not null,
    start_time_at timestamp not null,
    status_end  text check ( status_end in ('START', 'END') ),
    end_task_at timestamp
);

CREATE FUNCTION update_updated_on_user_task()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at
= now();
RETURN NEW;
END;
$$
language 'plpgsql';

CREATE FUNCTION update_end_task_at_task()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.end_task_at
= now();
RETURN NEW;
END;
$$
language 'plpgsql';

CREATE TRIGGER update_user_task_updated_at
    BEFORE UPDATE
    ON
        tasks
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_on_user_task();

CREATE TRIGGER update_end_task_at
    BEFORE UPDATE
    ON
        tasks
    FOR EACH ROW
    WHEN (NEW.status_end = 'END')
    EXECUTE PROCEDURE update_end_task_at_task();
