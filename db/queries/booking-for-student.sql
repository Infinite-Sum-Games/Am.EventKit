-- name: CreateBooking :one
INSERT INTO bookings (
  event_id,
  student_id,
  txn_id,
  registration_fee,
  txn_status,
  product_info,
  seats_released,
  metadata
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: CreateTeam :one
INSERT INTO teams (
  team_name,
  event_id,
  leader_name,
  booking_id,
  metadata
) VALUES ($1, $2, $3, $4, $5)
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
  e.id,
  e.price,
  e.is_group,
  e.is_per_head,
  e.max_teamsize,
  e.min_teamsize,
  e.total_seats,
  e.seats_filled,
  e.event_status,
  COALESCE(
    to_jsonb(ARRAY_AGG(t.abbreviation) FILTER (WHERE t.abbreviation LIKE '!%')),
    '[]'::jsonb
  ) AS special_tags
FROM event e
LEFT JOIN event_tag_mapping etm
  ON e.id = etm.event_id
LEFT JOIN tags t
  ON t.id = etm.tag_id
WHERE e.id = $1
GROUP BY e.id;


-- name: GetAnyBookingByUsersAndEvent :many
(
SELECT b.student_id
FROM bookings b
WHERE b.student_id = ANY($1::uuid[])
  AND b.event_id = $2
  AND b.txn_status != 'FAILED'
)
UNION
(
SELECT tm.student_id
FROM bookings b
JOIN teams t 
  ON t.booking_id = b.id
JOIN team_members tm 
  ON tm.team_id = t.id
WHERE tm.student_id = ANY($1::uuid[])
  AND b.event_id = $2
  AND b.txn_status != 'FAILED'
);

-- name: GetAnyPendingBookingByUser :many
SELECT id
FROM bookings
WHERE student_id = $1
AND txn_status = 'PENDING';

-- Hopefully, we can use this for removing decrement too (should try)
-- name: UpdateEventSeats :exec
UPDATE event
SET seats_filled = seats_filled + $1
WHERE id = $2;

-- name: GetBookingByTxnID :one
SELECT * FROM bookings WHERE txn_id = $1;

-- name: GetTeamMembersByTeamID :many
SELECT * FROM team_members
INNER JOIN teams t 
  ON team_members.team_id = t.id
INNER JOIN bookings b 
  ON t.booking_id = b.id
  AND b.txn_status = 'SUCCESS'
WHERE team_id = $1;

-- name: GetTeamIDByBooking :one
SELECT teams.id FROM teams 
INNER JOIN bookings b
  ON teams.booking_id = b.id 
  AND b.txn_status = 'SUCCESS'
WHERE booking_id = $1;

-- name: DeleteTeam :exec
DELETE 
FROM teams 
WHERE booking_id = $1;

-- name: DeleteTeamDetailsOfTeam :exec
DELETE 
FROM team_members 
WHERE team_id = $1;

-- name: UpdateBookingStatus :exec
UPDATE bookings
SET txn_status = $2
WHERE id = $1;

-- name: GetEmailByTxnId :one
SELECT s.email
FROM bookings b 
INNER JOIN student s 
  ON b.student_id = s.id
WHERE b.txn_id = $1;
