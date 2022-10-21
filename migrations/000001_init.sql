-- +goose Up
CREATE SCHEMA IF NOT EXISTS log;

CREATE TABLE IF NOT EXISTS log.logs (
    id serial primary key,
    raw text  NOT NULL,
    created_at timestamp NOT NULL DEFAULT timezone('utc', now())
);

-- +goose Down
DROP SCHEMA IF EXISTS log CASCADE;
