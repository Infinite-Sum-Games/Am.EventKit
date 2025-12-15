-- +goose Up
-- +goose StatementBegin
ALTER TABLE event
DROP CONSTRAINT event_cover_image_url_key; 

ALTER TABLE teams
ADD COLUMN metadata JSONB;
-- +goose StatementEnd


