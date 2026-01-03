-- +goose Up
-- +goose StatementBegin
CREATE TYPE payment_status_enum AS ENUM (
  'PENDING', 
  'COMPLETED', 
  'FAILED');
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE accomodation_details
ADD COLUMN payment_status TEXT DEFAULT 'PENDING' NOT NULL,
ADD COLUMN payment_expires TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accomodation_details
DROP COLUMN payment_status,
DROP COLUMN payment_expires;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TYPE payment_status_enum;
-- +goose StatementEnd
