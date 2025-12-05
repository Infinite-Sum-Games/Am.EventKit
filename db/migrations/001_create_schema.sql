-- +goose up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS "public";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE account_status_enum AS ENUM (
  'VERIFIED',
  'DISABLED'
);

CREATE TYPE organizer_type_enum AS ENUM (
  'DEPARTMENT',
  'CLUB'
);

CREATE TYPE event_type_enum AS ENUM (
  'EVENT',
  'WORKSHOP'
);

CREATE TYPE event_status_enum AS ENUM (
  'CLOSED',
  'ACTIVE',
  'COMPLETED'
);

CREATE TYPE event_mode_enum AS ENUM (
  'ONLINE',
  'OFFLINE'
);

CREATE TYPE attendance_mode_enum AS ENUM (
  'SOLO',
  'DUO'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS student (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  is_amrita_student BOOLEAN DEFAULT FALSE,
  amrita_roll_number TEXT,
  college_name TEXT DEFAULT 'Amrita Vishwa Vidyapeetham' NOT NULL,
  college_city TEXT DEFAULT 'Coimbatore' NOT NULL,
  account_status account_status_enum DEFAULT 'VERIFIED',
  refresh_token TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "student_pkey" PRIMARY KEY (id)
);
CREATE UNIQUE INDEX student_unique_roll_number
ON student(amrita_roll_number)
WHERE amrita_roll_number IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS student_onboarding (
  id SERIAL NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  is_amrita_student BOOLEAN NOT NULL DEFAULT TRUE,
  amrita_roll_number TEXT,
  college_name TEXT DEFAULT 'Amrita Vishwa Vidyapeetham' NOT NULL,
  college_city TEXT DEFAULT 'Coimbatore' NOT NULL,
  otp TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  expiry_at TIMESTAMP NOT NULL,

  CONSTRAINT "student_onboarding_pkey" PRIMARY KEY (id)
);

CREATE UNIQUE INDEX student_onboarding_unique_roll_number
ON student_onboarding(amrita_roll_number)
WHERE amrita_roll_number IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS password_reset (
  id SERIAL NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  otp TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  expiry_at TIMESTAMP NOT NULL,

  CONSTRAINT "password_reset_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizer (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE, -- eg: Computer Science and Engineering
  email TEXT NOT NULL UNIQUE, -- eg: cse@cb.amrita.edu
  password TEXT NOT NULL, -- eg: Single-time hash; need: For attendance
  org_type organizer_type_enum  NOT NULL, -- eg: DEPARTMENT | CLUB
  student_head TEXT NOT NULL,
  student_co_head TEXT,
  faculty_head TEXT NOT NULL,
  refresh_token TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "organizer_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  blurb TEXT NOT NULL,
  description TEXT NOT NULL,
  cover_image_url TEXT UNIQUE,
  price NUMERIC NOT NULL,
  is_per_head BOOLEAN NOT NULL,
  rules TEXT NOT NULL,
  event_type event_type_enum NOT NULL,
  is_group BOOLEAN NOT NULL,
  max_teamsize INTEGER,
  min_teamsize INTEGER,
  total_seats INTEGER NOT NULL,
  seats_filled INTEGER NOT NULL DEFAULT 0,
  event_status event_status_enum NOT NULL,
  event_mode event_mode_enum NOT NULL,
  attendance_mode attendance_mode_enum NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "event_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS favourites (
  id SERIAL NOT NULL,
  email TEXT NOT NULL,
  event_id UUID NOT NULL,

  CONSTRAINT "favourites_pkey" PRIMARY KEY (id),

  CONSTRAINT "favourites_email_event_id_unique" UNIQUE (email, event_id),

  CONSTRAINT "favourites_email" 
    FOREIGN KEY (email)
    REFERENCES student(email)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_schedule (
  id UUID DEFAULT gen_random_uuid(),
  event_id UUID NOT NULL,
  event_date DATE NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  venue TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "event_schedule_pkey" PRIMARY KEY (id),

  CONSTRAINT "event_schedule_event_id_fkey"
    FOREIGN KEY (event_id)
    REFERENCES event(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS people (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  phone_number TEXT NOT NULL UNIQUE,
  profession TEXT,
  email TEXT UNIQUE,

  CONSTRAINT "people_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE people_to_event_mapping (
  id SERIAL NOT NULL,
  event_id UUID NOT NULL,
  person_id UUID NOT NULL,

  CONSTRAINT "people_to_event_mapping_pkey" PRIMARY KEY (id),

  CONSTRAINT "people_to_event_mapping_event_id_fkey"
    FOREIGN KEY (event_id)
    REFERENCES event(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "people_to_event_mapping_person_id_fkey"
    FOREIGN KEY (person_id)
    REFERENCES people(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_to_organizer_mapping (
  id Serial NOT NULL,
  event_id UUID NOT NULL,
  organizer_id UUID NOT NULL,

  CONSTRAINT "event_to_organizer_mapping_pkey" PRIMARY KEY (id),

  CONSTRAINT "event_to_organizer_mapping_event_id_fkey"
    FOREIGN KEY (event_id)
    REFERENCES event(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "event_to_organizer_mapping_organizer_id_fkey"
    FOREIGN KEY (organizer_id)
    REFERENCES organizer(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  abbreviation TEXT NOT NULL UNIQUE,

  CONSTRAINT "tags_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_tag_mapping (
  id Serial NOT NULL,
  tag_id UUID NOT NULL,
  event_id UUID NOT NULL,

  CONSTRAINT "event_tag_mapping_pkey" PRIMARY KEY (id),

  CONSTRAINT "event_tag_mapping_tag_id_fkey"
    FOREIGN KEY (tag_id)
    REFERENCES tags(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "event_tag_mapping_event_id_fkey"
    FOREIGN KEY (event_id)
    REFERENCES event(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE bookings (
  id UUID DEFAULT gen_random_uuid(),
  txn_id TEXT NOT NULL,
  student_id UUID NOT NULL,
  event_id UUID NOT NULL ,
  registration_fee NUMERIC NOT NULL,
  product_info TEXT NOT NULL,
  seats_released  INTEGER NOT NULL DEFAULT 0,
  txn_status TEXT NOT NULL,
  team_details JSONB,
  metadata JSONB,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "bookings_pkey" PRIMARY KEY (id),

  CONSTRAINT "bookings_student_id_fkey"
  FOREIGN KEY (student_id)
  REFERENCES student(id)
  ON DELETE RESTRICT
  ON UPDATE CASCADE,

  CONSTRAINT "bookings_event_id_fkey"
  FOREIGN KEY (event_id)
  REFERENCES event(id)
  ON DELETE RESTRICT
  ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
  id UUID DEFAULT gen_random_uuid(),
  team_name CITEXT NOT NULL,
  event_id UUID NOT NULL,
  leader_name TEXT NOT NULL,
  booking_id UUID NOT NULL,

  CONSTRAINT "teams_pkey" PRIMARY KEY (id),

  CONSTRAINT "team_name_event_id_unique" UNIQUE (team_name, event_id),

  CONSTRAINT "teams_event_id_fkey"
    FOREIGN KEY (event_id)
    REFERENCES event(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "teams_booking_id_fkey"
  FOREIGN KEY (booking_id)
  REFERENCES bookings(id)
  ON DELETE RESTRICT
  ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS team_members (
  id UUID DEFAULT gen_random_uuid(),
  team_id UUID NOT NULL,
  student_id UUID NOT NULL,
  student_role TEXT NOT NULL,
  student_name TEXT NOT NULL,
  student_email TEXT NOT NULL,

  CONSTRAINT "team_members_pkey" PRIMARY KEY (id),

  CONSTRAINT "team_members_team_id_fkey"
    FOREIGN KEY(team_id)
    REFERENCES teams(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "team_members_student_id_fkey"
    FOREIGN KEY(student_id)
    REFERENCES student(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS team_events_attendance (
  id UUID DEFAULT gen_random_uuid(),
  student_id UUID NOT NULL,
  event_schedule_id UUID NOT NULL,
  check_in TIMESTAMP,
  check_out TIMESTAMP,

  CONSTRAINT "team_events_attendance_pkey" PRIMARY KEY (id),

  CONSTRAINT "team_events_attendance_student_id_fkey"
    FOREIGN KEY(student_id)
    REFERENCES student(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "team_events_attendance_event_schedule_id_fkey"
    FOREIGN KEY(event_schedule_id)
    REFERENCES event_schedule(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd
--
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS solo_event_participant (
  id SERIAL NOT NULL,
  student_id UUID NOT NULL,
  event_id UUID NOT NULL,
  event_schedule_id UUID NOT NULL,
  booking_id UUID NOT NULL,
  student_name TEXT NOT NULL,
  student_email TEXT NOT NULL,
  check_in TIMESTAMP,
  check_out TIMESTAMP,

  CONSTRAINT "solo_event_participant_pkey" PRIMARY KEY (id),

  CONSTRAINT "solo_event_participant_student_id_fkey"
    FOREIGN KEY(student_id)
    REFERENCES student(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,
  
  CONSTRAINT "solo_event_participant_event_id_fkey"
    FOREIGN KEY(event_id)
    REFERENCES event(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "solo_event_participant_event_schedule_id_fkey"
    FOREIGN KEY(event_schedule_id)
    REFERENCES event_schedule(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE,

  CONSTRAINT "solo_event_participant_booking_id_fkey"
    FOREIGN KEY(booking_id)
    REFERENCES bookings(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose down
-- +goose StatementBegin
DROP TABLE IF EXISTS solo_event_participant;
DROP TABLE IF EXISTS team_events_attendance;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS event_tag_mapping;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS event_to_organizer_mapping;
DROP TABLE IF EXISTS people_to_event_mapping;
DROP TABLE IF EXISTS people;
DROP TABLE IF EXISTS event_schedule;
DROP TABLE IF EXISTS favourites;
DROP TABLE IF EXISTS event;
DROP TABLE IF EXISTS organizer;
DROP TABLE IF EXISTS student_onboarding;
DROP TABLE IF EXISTS password_reset;
DROP TABLE IF EXISTS student;

DROP TYPE IF EXISTS attendance_mode_enum;
DROP TYPE IF EXISTS event_mode_enum;
DROP TYPE IF EXISTS event_status_enum;
DROP TYPE IF EXISTS event_type_enum;
DROP TYPE IF EXISTS organizer_type_enum;
DROP TYPE IF EXISTS account_status_enum;
-- +goose StatementEnd
