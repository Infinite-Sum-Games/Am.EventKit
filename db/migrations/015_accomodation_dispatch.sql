-- +goose Up
-- +goose StatementBegin
ALTER TABLE accomodation_personell
ADD COLUMN refresh_token TEXT;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hostel_check_in (

);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gate_management (
);
-- +goose StatementEnd
