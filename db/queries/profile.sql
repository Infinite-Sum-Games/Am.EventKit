-- name: FetchUserProfileQuery :one
SELECT 
  name, 
  email,
  phone_number,
  is_amrita_student,
  amrita_roll_number,
  college_name,
  college_city
FROM student 
WHERE account_status = 'VERIFIED' and email = $1;

-- name: EditUserProfileQuery :execrows
UPDATE student
SET 
  name = $2,
  phone_number = $3,
  college_name = $4,
  college_city = $5,
  updated_at = NOW()
WHERE 
  email = $1 
  AND account_status = 'VERIFIED';

-- name: GetAllTransactionsOfUserQuery :many
SELECT 
  b.id,
  b.txn_id,
  e.name AS event_name,
  b.registration_fee,
  b.txn_status,
  b.created_at
FROM bookings b
LEFT JOIN
  event e
ON b.event_id = e.id
WHERE b.student_id = $1;

-- name: GetMySoloEventTickets :many
SELECT
  e.id AS event_id,
  e.name AS event_name,
  e.price,
  e.is_technical,
  e.event_mode,
  e.event_type AS event_type,

  COALESCE(
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'schedule_id', es.id,
        'event_date', es.event_date,
        'start_time', es.start_time,
        'end_time', es.end_time,
        'venue', es.venue
    )) FILTER (WHERE e.id IS NOT NULL),
    '[]'::jsonb
  ) AS schedules

FROM event e
LEFT JOIN solo_event_participant sep 
  ON e.id = sep.event_id
LEFT JOIN bookings b 
  ON sep.booking_id = b.id
LEFT JOIN student s 
  ON sep.student_id = s.id
LEFT JOIN event_schedule es
  ON e.id = es.event_id
WHERE
  s.email = $1
  AND s.id = $2
  AND b.txn_status = 'SUCCESS'
GROUP BY
  e.id,
  e.name,
  e.price,
  e.is_technical,
  e.event_mode;

-- name: GetMyTeamEventTickets :many
SELECT
  e.id AS event_id,
  e.name AS event_name,
  e.price,
  e.is_technical,
  e.event_mode,
  t.team_name,
  e.event_type AS event_type,

  COALESCE(
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'schedule_id', es.id,
        'event_date', es.event_date,
        'start_time', es.start_time,
        'end_time', es.end_time,
        'venue', es.venue
    )) FILTER (WHERE es.id IS NOT NULL),
    '[]'::jsonb
  ) AS schedules

FROM event e
LEFT JOIN teams t ON e.id = t.event_id
LEFT JOIN team_members tm ON t.id = tm.team_id
LEFT JOIN bookings b ON t.booking_id = b.id
LEFT JOIN student s ON tm.student_id = s.id
LEFT JOIN event_schedule es ON e.id = es.event_id

WHERE
  s.email = $1
  AND b.txn_status = 'SUCCESS'
GROUP BY
  e.id,
  t.team_name;
