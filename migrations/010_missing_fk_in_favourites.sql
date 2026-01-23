-- +goose Up
-- +goose StatementBegin
ALTER TABLE favourites
ADD CONSTRAINT "favourites_event_id_fkey"
FOREIGN KEY (event_id)
REFERENCES event (id);
-- +goose StatementEnd
