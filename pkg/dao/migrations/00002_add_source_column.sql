-- +goose Up
ALTER TABLE team_kill_log
    ADD COLUMN source VARCHAR;

-- +goose Down
ALTER TABLE team_kill_log DROP COLUMN source;





