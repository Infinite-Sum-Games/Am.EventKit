-- +goose Up
-- +goose StatementBegin
ALTER TABLE event
ADD COLUMN is_technical BOOLEAN DEFAULT false;
-- +goose StatementEnd
