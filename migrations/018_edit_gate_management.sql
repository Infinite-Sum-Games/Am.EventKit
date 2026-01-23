-- +goose Up

-- +goose StatementBegin
ALTER TABLE gate_management
DROP COLUMN accomodation_id;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE gate_management
ADD COLUMN student_id UUID NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE gate_management
ADD CONSTRAINT "gate_management_student_id_fkey"
  FOREIGN KEY (student_id)
  REFERENCES student(id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;
-- +goose StatementEnd
