-- +goose Up
-- +goose StatementBegin
-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pg_cron;
-- +goose StatementEnd

-- +goose StatementBegin
-- Enums
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

CREATE TYPE dispute_status_enum AS ENUM (
  'OPEN',
  'CLOSED_AS_TRUE',
  'CLOSED_AS_FALSE'
);

CREATE TYPE gate_log_direction_enum AS ENUM (
  'IN',
  'OUT'
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Tables
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
  hospitality_id TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "student_pkey" PRIMARY KEY (id),
  CONSTRAINT unique_hospitality_id UNIQUE (hospitality_id)
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

  CONSTRAINT "password_reset_pkey" PRIMARY KEY (id),
  CONSTRAINT unique_password_reset_email UNIQUE (email)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS admin (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  refresh_token TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT admin_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizer (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  org_type organizer_type_enum  NOT NULL,
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
  cover_image_url TEXT,
  price INTEGER NOT NULL,
  is_per_head BOOLEAN NOT NULL,
  rules TEXT NOT NULL,
  event_type event_type_enum NOT NULL,
  is_group BOOLEAN NOT NULL,
  is_technical BOOLEAN DEFAULT false,
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
      ON UPDATE CASCADE,
  CONSTRAINT "favourites_event_id_fkey"
  FOREIGN KEY (event_id)
  REFERENCES event (id)
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
  event_day INTEGER[],

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
  registration_fee INTEGER NOT NULL,
  registration_fee_without_gst INTEGER DEFAULT 0,
  product_info TEXT NOT NULL,
  seats_released  INTEGER NOT NULL DEFAULT 0,
  txn_status TEXT NOT NULL,
  team_details JSONB,
  metadata JSONB,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "bookings_pkey" PRIMARY KEY (id),
  CONSTRAINT "bookings_txn_id_unique" UNIQUE (txn_id),
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
  metadata JSONB,

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

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dispute (
  id UUID DEFAULT gen_random_uuid(),
  txn_id TEXT NOT NULL,
  student_email TEXT,
  description TEXT,
  event_id UUID NOT NULL,
  dispute_status dispute_status_enum DEFAULT 'OPEN',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "dispute_pkey" PRIMARY KEY (id),
  CONSTRAINT "dispute_to_student_mapping_fkey"
    FOREIGN KEY (student_email)
    REFERENCES student (email)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT "dispute_to_event_mapping_fkey"
    FOREIGN KEY (event_id)
    REFERENCES event (id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT "dispute_to_txn_mapping_fkey"
    FOREIGN KEY (txn_id)
    REFERENCES bookings (txn_id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE
);
-- +goose StatementEnd

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
  hostel_name TEXT NOT NULL,
  amrita_dayscholar_price INTEGER DEFAULT 0 NOT NULL,
  non_amrita_price INTEGER DEFAULT 0 NOT NULL,
  room_filled INTEGER NOT NULL DEFAULT 0,
  CONSTRAINT room_filled_lte_room_count
  CHECK (room_filled <= room_count),

  CONSTRAINT "hostel_metadata_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accomodation_personell (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  refresh_token TEXT,
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
  payment_status TEXT DEFAULT '' NOT NULL,
  payment_expires TIMESTAMP,
  day_count INTEGER DEFAULT 0 NOT NULL,
  check_in TIMESTAMP NOT NULL,
  check_out TIMESTAMP NOT NULL,
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
  direction gate_log_direction_enum NOT NULL,
  logged_at TIMESTAMP NOT NULL DEFAULT NOW(),
  personell_id UUID NOT NULL,
  student_id UUID NOT NULL,

  CONSTRAINT "gate_management_pkey" PRIMARY KEY (id),
  CONSTRAINT "gate_management_student_id_fkey"
  FOREIGN KEY (student_id)
  REFERENCES student(id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT "gate_management_personell_id_fkey" 
    FOREIGN KEY (personell_id)
    REFERENCES accomodation_personell(id)
      ON DELETE RESTRICT
      ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Materialized Views
CREATE MATERIALIZED VIEW revenue_analytics AS
SELECT
    1 AS id,
    -- TOTAL REVENUE
    (
        SELECT COALESCE(SUM(registration_fee), 0)
        FROM bookings
        WHERE txn_status = 'SUCCESS'
    ) AS total_revenue,
    -- REVENUE GROUPED BY EVENTS
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                e.id AS event_id,
                e.name AS event_name,
                SUM(b.registration_fee) AS revenue
            FROM bookings b
            JOIN event e ON e.id = b.event_id
            WHERE b.txn_status = 'SUCCESS'
            GROUP BY e.id
            ORDER BY e.name
        ) x
    ) AS revenue_per_event,
    -- REVENUE GROUPED BY DATE
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                DATE(b.created_at) AS revenue_date,
                SUM(b.registration_fee) AS revenue
            FROM bookings b
            WHERE b.txn_status = 'SUCCESS'
            GROUP BY revenue_date
            ORDER BY revenue_date
        ) x
    ) AS revenue_per_date,
    -- REVENUE GROUPED BY ORGANIZERS
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                o.id AS organizer_id,
                o.name AS organizer_name,
                SUM(b.registration_fee) AS revenue
            FROM bookings b
            JOIN event_to_organizer_mapping m ON m.event_id = b.event_id
            JOIN organizer o ON o.id = m.organizer_id
            WHERE b.txn_status = 'SUCCESS'
            GROUP BY o.id
            ORDER BY o.name
        ) x
    ) AS revenue_per_organizer;

CREATE UNIQUE INDEX revenue_analytics_id_index 
ON revenue_analytics(id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW event_registration_analytics AS
WITH participants AS (
    SELECT student_id FROM solo_event_participant
    UNION
    SELECT student_id FROM team_members
)
SELECT
    1 AS id,
    -- TOTAL EVENT REGISTRATIONS
    (SELECT COUNT(*) FROM participants) AS total_event_registrations,
    -- PARTICIPANTS VS NON-PARTICIPANTS
    (
        SELECT jsonb_build_object(
            'participants', (SELECT COUNT(*) FROM participants),
            'non_participants', 
                (SELECT COUNT(*) FROM student s 
                 WHERE s.id NOT IN (SELECT student_id FROM participants))
        )
    ) AS participant_split,
    -- (REGISTRATIONS VS TOTAL SEATS) PER EVENT
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                e.id AS event_id,
                e.name AS event_name,
                e.total_seats,
                COUNT(DISTINCT sp.student_id) 
                    + COUNT(DISTINCT tm.student_id) AS registered_count
            FROM event e
            LEFT JOIN solo_event_participant sp ON sp.event_id = e.id
            LEFT JOIN teams t ON t.event_id = e.id
            LEFT JOIN team_members tm ON tm.team_id = t.id
            GROUP BY e.id
            ORDER BY e.name
        ) x
    ) AS event_registration_stats;

CREATE UNIQUE INDEX event_registration_analytics_id_index 
ON event_registration_analytics(id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW people_registration_analytics AS
SELECT
    1 AS id,
    -- INTERNAL VS EXTERNAL WEBSITE REGISTRATIONS
    (
        SELECT jsonb_build_object(
            'internal', COUNT(*) FILTER (WHERE is_amrita_student = TRUE),
            'external', COUNT(*) FILTER (WHERE is_amrita_student = FALSE)
        )
        FROM student
    ) AS website_registration_split,
    -- TOTAL WEBSITE REGISTRATIONS
    (
        SELECT COUNT(*) FROM student
    ) AS total_website_registrations;

CREATE UNIQUE INDEX people_registration_analytics_id_index 
ON people_registration_analytics(id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW transaction_analytics AS
SELECT
    1 AS id,
    -- FAILED VS SUCCESSFUL VS PENDING TRANSACTIONS
    (
        SELECT jsonb_build_object(
            'success', COUNT(*) FILTER (WHERE txn_status='SUCCESS'),
            'failed', COUNT(*) FILTER (WHERE txn_status='FAILED'),
            'pending', COUNT(*) FILTER (WHERE txn_status='PENDING'),
            'total', COUNT(*)
        )
        FROM bookings
    ) AS transaction_summary;

CREATE UNIQUE INDEX transaction_analytics_id_index 
ON transaction_analytics(id);
-- +goose StatementEnd

-- +goose StatementBegin
-- Schedule analytics refresh
SELECT cron.schedule(
    'refresh_analytics_every_30min',
    '*/30 * * * *',
$$
REFRESH MATERIALIZED VIEW CONCURRENTLY revenue_analytics;
REFRESH MATERIALIZED VIEW CONCURRENTLY event_registration_analytics;
REFRESH MATERIALIZED VIEW CONCURRENTLY people_registration_analytics;
REFRESH MATERIALIZED VIEW CONCURRENTLY transaction_analytics;
$$
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Unschedule analytics refresh
SELECT cron.unschedule('refresh_analytics_every_30min');

-- Drop materialized views
DROP MATERIALIZED VIEW IF EXISTS revenue_analytics;
DROP MATERIALIZED VIEW IF EXISTS event_registration_analytics;
DROP MATERIALIZED VIEW IF EXISTS people_registration_analytics;
DROP MATERIALIZED VIEW IF EXISTS transaction_analytics;

-- Drop tables
DROP TABLE IF EXISTS gate_management;
DROP TABLE IF EXISTS hostel_check_in;
DROP TABLE IF EXISTS accomodation_details;
DROP TABLE IF EXISTS accomodation_personell;
DROP TABLE IF EXISTS hostel_metadata;
DROP TABLE IF EXISTS dispute;
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
DROP TABLE IF EXISTS admin;
DROP TABLE IF EXISTS password_reset;
DROP TABLE IF EXISTS student_onboarding;
DROP TABLE IF EXISTS student;

-- Drop types
DROP TYPE IF EXISTS gate_log_direction_enum;
DROP TYPE IF EXISTS dispute_status_enum;
DROP TYPE IF EXISTS attendance_mode_enum;
DROP TYPE IF EXISTS event_mode_enum;
DROP TYPE IF EXISTS event_status_enum;
DROP TYPE IF EXISTS event_type_enum;
DROP TYPE IF EXISTS organizer_type_enum;
DROP TYPE IF EXISTS account_status_enum;

-- Drop extensions
DROP EXTENSION IF EXISTS pg_cron;
DROP EXTENSION IF EXISTS citext;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd