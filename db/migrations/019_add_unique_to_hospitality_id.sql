-- +goose Up
-- +goose StatementBegin
ALTER TABLE student
ADD CONSTRAINT unique_hospitality_id UNIQUE (hospitality_id);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE hostel_metadata
RENAME COLUMN price TO amrita_dayscholar_price;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE hostel_metadata
ADD COLUMN non_amrita_price INTEGER DEFAULT 0 NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE student
DROP CONSTRAINT unique_hospitality_id;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE hostel_metadata
RENAME COLUMN amrita_dayscholar_price TO price;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE hostel_metadata
DROP COLUMN non_amrita_price;
-- +goose StatementEnd
