-- +goose Up
-- +goose StatementBegin
ALTER TABLE bookings
ADD CONSTRAINT "bookings_txn_id_unique"
UNIQUE (txn_id);
-- +goose StatementEnd
