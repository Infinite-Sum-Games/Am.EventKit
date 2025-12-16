-- +goose Up
-- +goose StatementBegin
ALTER TABLE event
ALTER COLUMN price
TYPE INTEGER
USING TRUNC(price);
-- +goose StatementEnd
