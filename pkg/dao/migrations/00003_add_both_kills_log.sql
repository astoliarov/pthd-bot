-- +goose Up
CREATE TABLE IF NOT EXISTS bot_kill_log
(
    id          INTEGER  not null
        constraint id
            primary key autoincrement,
    victim      varchar  not null,
    source      varchar  not null,
    happened_at datetime not null
);

-- +goose Down
ALTER TABLE bot_kill_log DROP COLUMN source;





