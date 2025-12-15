-- +goose Up
-- +goose StatementBegin
ALTER TABLE event
DROP CONSTRAINT event_cover_image_url_key; 
-- +goose StatementEnd
