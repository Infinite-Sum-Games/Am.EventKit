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
ADD COLUMN payment_expires TIMESTAMP,
ADD COLUMN day_count INTEGER DEFAULT 0 NOT NULL,
DROP COLUMN is_paid;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE hostel_metadata
ADD COLUMN price INTEGER DEFAULT 0 NOT NULL,
ADD COLUMN room_filled INTEGER NOT NULL DEFAULT 0,
ADD CONSTRAINT room_filled_lte_room_count
CHECK (room_filled <= room_count);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accomodation_details
DROP COLUMN payment_status,
DROP COLUMN day_count,
DROP COLUMN payment_expires;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE hostel_metadata
DROP COLUMN price,
DROP COLUMN room_filled;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TYPE payment_status_enum;
-- +goose StatementEnd
