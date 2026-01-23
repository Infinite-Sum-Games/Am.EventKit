-- +goose Up

-- +goose StatementBegin
ALTER TABLE accomodation_personell
ADD COLUMN refresh_token TEXT;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE accomodation_details
ALTER COLUMN check_in SET NOT NULL,
ALTER COLUMN check_out SET NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE hostel_metadata
ADD COLUMN hostel_name TEXT NOT NULL;
-- +goose StatementEnd
