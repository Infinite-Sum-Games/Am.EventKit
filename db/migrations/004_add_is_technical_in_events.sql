-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
ADD COLUMN is_technical BOOLEAN DEFAULT false;
-- +goose StatementEnd
