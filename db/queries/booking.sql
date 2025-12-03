-- name: CreateBooking :one
INSERT INTO bookings (
  event_id,
  student_id,
  txn_id,
  registration_fee,
  txn_status,
  product_info,
  seats_released
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: CreateTeam :one
INSERT INTO teams (
  team_name,
  event_id,
  leader_name,
  booking_id
) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CreateTeamMember :one
INSERT INTO team_members (
  team_id,
  student_id,
  student_role,
  student_name,
  student_email
) VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetEventForBooking :one
SELECT
  id,
  price,
  is_group,
  is_per_head,
  max_teamsize,
  min_teamsize,
  total_seats,
  seats_filled,
  event_status
FROM 
  event
WHERE 
  id = $1;

-- name: GetAnyBookingByUserAndEvent :one
SELECT 
  b.id,
  b.txn_status
FROM bookings b
LEFT JOIN teams t 
  ON b.id = t.booking_id
LEFT JOIN team_members tm 
  ON t.id = tm.team_id
WHERE 
  (
    b.student_id = $1
    OR tm.student_id = $1
  )
  AND b.event_id = $2
  AND b.txn_status != 'failed'
LIMIT 1;

-- Hopefully, we can use this for removing decrement too (should try)
-- name: UpdateEventSeats :exec
UPDATE event
SET seats_filled = seats_filled + $1
WHERE id = $2;