-- name: SeedBookingsQuery :exec
INSERT INTO bookings(
  txn_id ,
  student_id,
  event_id,
  registration_fee,
  product_info,
  seats_released,
  txn_status
) VALUES($1, $2, $3, $4, $5, $6, $7);

-- name: ViewBookingsSeedQuery :many
SELECT 
  id, 
  txn_id, 
  student_id, 
  event_id, 
  registration_fee, 
  product_info, 
  seats_released, 
  txn_status
FROM bookings;

-- name: SeedAdminQuery :exec
INSERT INTO admin(
  name, 
  email,
  password
) VALUES ($1, $2, $3);

-- name: SeedAmritaStudentQuery :exec
INSERT INTO student(
  name, 
  email,
  password,
  phone_number,
  is_amrita_student,
  amrita_roll_number
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: SeedNonAmritaStudentQuery :exec
INSERT INTO student(
  name,
  email,
  password,
  phone_number
) VALUES ($1, $2, $3, $4);

-- name: SeedEventQuery :exec
INSERT INTO event(
  name, 
  blurb, 
  description, 
  price, 
  is_per_head, 
  rules, 
  event_type, 
  is_group, 
  total_seats, 
  seats_filled, 
  event_status, 
  event_mode, 
  attendance_mode,
  cover_image_url
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);

-- name: SeedOrganizerQuery :exec
INSERT INTO organizer(
  name, 
  email, 
  password,
  org_type, 
  student_head, 
  student_co_head, 
  faculty_head
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: SeedEventToOrganizerMappingQuery :exec
INSERT INTO event_to_organizer_mapping(
  event_id, 
  organizer_id
) VALUES ($1, $2);

-- name: SeedEventScheduleQuery :exec
INSERT INTO event_schedule(
  event_id, 
  event_date, 
  start_time, 
  end_time, 
  venue
) VALUES ($1, $2, $3, $4, $5);

-- name: SeedPeopleQuery :exec
INSERT INTO people(
  name, 
  phone_number
) VALUES ($1, $2);

-- name: SeedPeopleToEventMappingQuery :exec
INSERT INTO people_to_event_mapping(
  event_id, 
  person_id,
  event_day
) VALUES ($1, $2, $3);

-- name: SeedTagsQuery :exec
INSERT INTO tags(
  name, 
  abbreviation
) VALUES ($1, $2);

-- name: SeedEventTagMappingQuery :exec
INSERT INTO event_tag_mapping(
  tag_id, 
  event_id
) VALUES ($1, $2);

-- name: ViewAdminSeedQuery :many
SELECT * FROM admin;

-- name: ViewStudentSeedQuery :many
SELECT
  id, 
  name, 
  email, 
  phone_number, 
  is_amrita_student, 
  amrita_roll_number
FROM student;

-- name: ViewEventSeedQuery :many
SELECT 
  id, 
  name, 
  blurb, 
  description, 
  price, 
  is_per_head, 
  rules, 
  event_type, 
  is_group, 
  total_seats, 
  seats_filled, 
  event_status, 
  event_mode, 
  attendance_mode
FROM event;

-- name: ViewPeopleSeedQuery :many
SELECT 
  id, 
  name, 
  phone_number, 
  profession, 
  email
FROM people;

-- name: ViewEventToOrganizerMappingSeedQuery :many
SELECT 
  id, 
  event_id, 
  organizer_id
FROM event_to_organizer_mapping;

-- name: ViewEventScheduleSeedQuery :many
SELECT 
  id, 
  event_id, 
  event_date, 
  start_time, 
  end_time, 
  venue
FROM event_schedule;

-- name: ViewPeopleToEventMappingSeedQuery :many
SELECT 
  id, 
  event_id, 
  person_id
FROM people_to_event_mapping;

-- name: ViewEventTagMappingSeedQuery :many
SELECT 
  id, 
  tag_id, 
  event_id
FROM event_tag_mapping;

-- name: ViewOrganizerSeedQuery :many
SELECT
  id,
  name as organizer_name,
  email as organizer_email,
  org_type as organizer_type,
  student_head,
  student_co_head,
  faculty_head
FROM organizer;

-- name: ViewTagSeedQuery :many
SELECT
    id,
    name,
    abbreviation
FROM tags;

-- name: TruncateAllTablesQuery :exec
TRUNCATE TABLE 
  student,
  organizer, 
  people, 
  tags, 
  event, 
  event_to_organizer_mapping, 
  event_schedule, 
  people_to_event_mapping, 
  event_tag_mapping 
CASCADE;

