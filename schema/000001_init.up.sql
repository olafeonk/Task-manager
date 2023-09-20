CREATE TABLE users
(
    id            serial       not null unique,
    username      varchar(255) not null,
    password_hash varchar(255) not null
);

CREATE TABLE tasks
(
    id          serial                                      not null unique,
    created_at  timestamp                                   not null default CURRENT_TIMESTAMP,
    updated_at  timestamp                                   not null default CURRENT_TIMESTAMP,
    user_id     int references users (id) on delete cascade not null,
    name        varchar(200)                                not null,
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

CREATE UNIQUE INDEX users_idx ON users (username, password_hash);
