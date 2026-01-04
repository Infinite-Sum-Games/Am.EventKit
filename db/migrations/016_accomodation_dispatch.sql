-- +goose Up

-- +goose StatementBegin
ALTER TABLE students
ADD COLUMN hospitality_id TEXT;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE gate_log_direction_enum AS ENUM (
  'IN',
  'OUT'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hostel_check_in (
  id UUID DEFAULT gen_random_uuid(),
  accomodation_id UUID NOT NULL UNIQUE,
  checked_in_at TIMESTAMP NOT NULL DEFAULT NOW(),
  checked_out_at TIMESTAMP,
  checked_in_by UUID NOT NULL,
  checked_out_by UUID,

  CONSTRAINT "hostel_check_in_pkey" PRIMARY KEY (id),

  CONSTRAINT "hostel_check_in_accomodation_id_fkey" 
    FOREIGN KEY (accomodation_id)
    REFERENCES accomodation_details(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "hostel_check_in_checked_in_by_fkey" 
    FOREIGN KEY (checked_in_by)
    REFERENCES accomodation_personell(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "hostel_check_in_checked_out_by_fkey" 
    FOREIGN KEY (checked_out_by)
    REFERENCES accomodation_personell(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gate_management (
  id UUID DEFAULT gen_random_uuid(),
  accomodation_id UUID NOT NULL,
  direction gate_log_direction_enum NOT NULL,
  logged_at TIMESTAMP NOT NULL DEFAULT NOW(),
  personell_id UUID NOT NULL,

  CONSTRAINT "gate_management_pkey" PRIMARY KEY (id),

  CONSTRAINT "gate_management_accomodation_id_fkey" 
    FOREIGN KEY (accomodation_id)
    REFERENCES accomodation_details(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "gate_management_personell_id_fkey" 
    FOREIGN KEY (personell_id)
    REFERENCES accomodation_personell(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,
);
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS gate_management;
DROP TABLE IF EXISTS hostel_check_in;
DROP TYPE IF EXISTS gate_log_direction_enum;
-- +goose StatementEnd
