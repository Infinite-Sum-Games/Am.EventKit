-- +goose Up
-- +goose StatementBegin
ALTER TABLE people_to_event_mapping
    ADD COLUMN event_day INTEGER[];
-- +goose StatementEnd
