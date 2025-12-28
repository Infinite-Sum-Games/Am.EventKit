-- +goose Up
-- +goose StatementBegin
ALTER TABLE bookings
ADD COLUMN registration_fee_without_gst INTEGER DEFAULT 0;
-- +goose StatementEnd
