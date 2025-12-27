-- name: GetStudentByEmail :one
SELECT * FROM student
WHERE email = $1;

-- name: GetStudentsByEmails :many
SELECT * FROM student
WHERE email = ANY($1::text[]);

-- name: FetchEventsByOrganizerQuery :many
SELECT
  e.id AS event_id,
  e.name AS event_name,
  e.is_group,
  COALESCE(
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
      'id', es.id,
      'event_date', es.event_date,
      'start_time', es.start_time,
      'end_time', es.end_time,
      'venue', es.venue
    )) FILTER (WHERE es.id IS NOT NULL),
    '[]'::jsonb
  ) AS schedules
FROM event e
JOIN event_to_organizer_mapping etom
  ON e.id = etom.event_id
LEFT JOIN event_schedule es
  ON e.id = es.event_id
WHERE
  etom.organizer_id = $1
GROUP BY
  e.id, e.name, e.is_group;

-- name: FetchParticipantsBySoloEventQuery :many
SELECT
  sep.id AS attendance_id,
  sep.student_id AS student_id,
  sep.student_name AS student_name,
  sep.student_email AS student_email,
  sep.check_in AS check_in,
  sep.check_out AS check_out
FROM solo_event_participant sep
WHERE sep.event_schedule_id = $1;

-- name: FetchParticipantsByTeamEventQuery :many
SELECT
  tea.id AS attendance_id,
  tea.student_id AS student_id,
  s.name AS student_name,
  s.email AS student_email,
  tea.check_in AS check_in,
  tea.check_out AS check_out
FROM team_events_attendance tea
INNER JOIN student s
  ON tea.student_id = s.id
WHERE tea.event_schedule_id = $1;

-- name: MarkSoloCheckInQuery :execrows
UPDATE solo_event_participant
SET check_in = NOW()
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: MarkSoloCheckOutQuery :execrows
UPDATE solo_event_participant
SET check_out = NOW()
WHERE student_id = $1
  AND event_schedule_id = $2
  AND check_in IS NOT NULL;

-- name: MarkSoloBothQuery :execrows
UPDATE solo_event_participant
SET check_in = NOW(), 
  check_out = NOW()
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: MarkTeamCheckInQuery :execrows
UPDATE team_events_attendance
SET check_in = NOW()
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: MarkTeamCheckOutQuery :execrows
UPDATE team_events_attendance
SET check_out = NOW()
WHERE student_id = $1
  AND event_schedule_id = $2
  AND check_in IS NOT NULL;

-- name: MarkTeamBothQuery :execrows
UPDATE team_events_attendance
SET check_in = NOW(), 
  check_out = NOW()
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: UnMarkSoloCheckInQuery :execrows
UPDATE solo_event_participant
SET check_in = NULL
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: UnMarkSoloCheckOutQuery :execrows
UPDATE solo_event_participant
SET check_out = NULL
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: UnMarkSoloBothQuery :execrows
UPDATE solo_event_participant
SET check_in = NULL, 
  check_out = NULL
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: UnMarkTeamCheckInQuery :execrows
UPDATE team_events_attendance
SET check_in = NULL
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: UnMarkTeamCheckOutQuery :execrows
UPDATE team_events_attendance
SET check_out = NULL
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: UnMarkTeamBothQuery :execrows
UPDATE team_events_attendance
SET check_in = NULL, 
  check_out = NULL
WHERE student_id = $1
  AND event_schedule_id = $2;

-- name: CheckStudentRegisteredForEvent :one
SELECT 
    id,
    txn_id,
    student_id,
    event_id,
    registration_fee,
    product_info,
    seats_released,
    txn_status,
    team_details,
    metadata,
    created_at,
    updated_at
FROM bookings
WHERE 
  student_id = $1
  AND event_id = $2
  AND txn_status = 'SUCCESS'
LIMIT 1;

-- name: GetScheduleById :one
SELECT 
    id,
    event_id,
    event_date,
    start_time,
    end_time,
    venue,
    created_at,
    updated_at
FROM event_schedule
WHERE id = $1;

-- name: GetAttendanceRecord :one
SELECT
    id,
    student_id,
    event_id,
    event_schedule_id,
    booking_id,
    student_name,
    student_email,
    check_in,
    check_out
FROM solo_event_participant
WHERE student_id = $1
  AND event_schedule_id = $2
LIMIT 1;

-- name: InsertCheckIn :one
INSERT INTO solo_event_participant (
    student_id,
    event_id,
    event_schedule_id,
    booking_id,
    student_name,
    student_email,
    check_in
) VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING
    id,
    student_id,
    event_id,
    event_schedule_id,
    booking_id,
    student_name,
    student_email,
    check_in,
    check_out;

-- name: UpdateCheckOut :one
UPDATE solo_event_participant
SET check_out = NOW()
WHERE id = $1
RETURNING
    id,
    student_id,
    event_id,
    event_schedule_id,
    booking_id,
    student_name,
    student_email,
    check_in,
    check_out;

-- name: GetSchedulesByEventID :many
SELECT id FROM event_schedule WHERE event_id = $1;

-- name: CreateSoloEventParticipant :one
INSERT INTO solo_event_participant (
  student_id, 
  event_id, 
  event_schedule_id, 
  booking_id,
  student_name,
  student_email
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: CreateTeamAttendance :one
INSERT INTO team_events_attendance (
  student_id, 
  event_schedule_id
) VALUES ($1, $2)
RETURNING id;
