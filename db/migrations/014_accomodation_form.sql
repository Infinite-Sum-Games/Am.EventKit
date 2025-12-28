-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hostel_metadata (
  id UUID DEFAULT gen_random_uuid(),
  room_count INTEGER DEFAULT 0 NOT NULL,
  is_male BOOLEAN DEFAULT TRUE NOT NULL,
  warden_email TEXT,
  warden_password TEXT,
  warden_refresh_token TEXT,
  latitude TEXT,
  longtitude TEXT,
  map_url TEXT,

  CONSTRAINT "hostel_metadata_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accomodation_personell (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "accomodation_personell_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accomodation_details (
  id UUID DEFAULT gen_random_uuid(),
  student_id UUID NOT NULL,
  hostel_id UUID,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  is_male BOOLEAN DEFAULT TRUE NOT NULL,
  is_hosteller BOOLEAN DEFAULT TRUE NOT NULL,
  college_roll_number TEXT NOT NULL,
  college_name TEXT NOT NULL,
  room_preference TEXT NOT NULL,
  is_amrita_campus BOOLEAN DEFAULT FALSE NOT NULL,
  is_paid BOOLEAN DEFAULT FALSE NOT NULL,
  check_in TIMESTAMP,
  check_out TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "accomodation_details_pkey" PRIMARY KEY (id),

  CONSTRAINT "accomodation_details_student_id_fkey"
    FOREIGN KEY (student_id)
    REFERENCES student(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "accomodation_details_hostel_id_fkey"
    FOREIGN KEY (hostel_id)
    REFERENCES hostel_metadata(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accomodation_form_resp;
DROP TABLE IF EXISTS accomodation_personell;
-- +goose StatementEnd
