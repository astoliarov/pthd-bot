-- +goose Up
CREATE TABLE IF NOT EXISTS team_kill_log
(
    id          INTEGER  not null
        constraint id
            primary key autoincrement,
    killer      varchar  not null,
    victim      varchar  not null,
    happened_at datetime not null
);

-- +goose Down
DROP TABLE team_kill_log;
