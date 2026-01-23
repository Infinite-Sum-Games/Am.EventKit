-- +goose Up
-- +goose StatementBegin
ALTER TABLE accomodation_details 
ALTER COLUMN payment_status SET DEFAULT '';
-- +goose StatementEnd

-- The old data will not migrate so need a manual query on prod data.
-- 
-- UPDATE accomodation_details 
-- SET payment_status = '' 
-- WHERE payment_status = 'PENDING';
