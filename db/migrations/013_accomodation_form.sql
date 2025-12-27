-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accomodation_form_resp (
  id UUID DEFAULT gen_random_uuid(),
  student_id UUID NOT NULL,
  name TEXT NOT NULL,
  is_male BOOLEAN DEFAULT TRUE,
  college_roll_number TEXT NOT NULL,
  college_name_mentioned TEXT NOT NULL,
  is_amrita BOOLEAN DEFAULT FALSE,
  check_in TIMESTAMP,
  check_out TIMESTAMP,

  CONSTRAINT "accomodation_form_resp_pkey" PRIMARY KEY (id),

  CONSTRAINT "accomodation_form_resp_student_id"
    FOREIGN KEY (student_id)
    REFERENCES student(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accomodation_form_resp;
-- +goose StatementEnd
