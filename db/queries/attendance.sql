-- name: GetEventsByDateAndOrganizer :many
SELECT
    e.id,
    e.name AS event_name,
    e.blurb,
    e.description AS event_description,
    e.cover_image_url,
    e.price,
    e.is_per_head,
    e.rules,
    e.event_type,
    e.is_group,
    e.max_teamsize,
    e.min_teamsize,
    e.total_seats,
    e.seats_filled,
    e.event_status,
    e.event_mode,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'organizer_name', o.name,
        'org_abbreviation', o.abbr,
        'org_type', o.org_type
      )) FILTER (WHERE o.id IS NOT NULL),
      '[]'::jsonb
    ) AS organizers,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'event_date', es.event_date,
        'start_time', es.start_time,
        'end_time', es.end_time,
        'venue', es.venue
      )) FILTER (WHERE es.id IS NOT NULL),
      '[]'::jsonb
    ) AS schedules,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'tag_name', t.name,
        'tag_abbreviation', t.abbreviation
      )) FILTER (WHERE t.id IS NOT NULL),
      '[]'::jsonb
    ) AS tags

FROM event e

LEFT JOIN event_schedule es 
  ON e.id = es.event_id
LEFT JOIN event_to_organizer_mapping m 
  ON e.id = m.event_id
LEFT JOIN organizer o 
  ON m.organizer_id = o.id
LEFT JOIN event_tag_mapping etm 
  ON e.id = etm.event_id
LEFT JOIN tags t 
  ON etm.tag_id = t.id

WHERE 
  es.event_date = $1
  AND o.email = $2

GROUP BY e.id
ORDER BY es.start_time ASC;

-- name: GetStudentByEmail :one
SELECT * FROM student
WHERE email = $1;

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

