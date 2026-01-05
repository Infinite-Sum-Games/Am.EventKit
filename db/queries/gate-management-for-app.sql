-- name: GateCheckInQuery :one

-- name: GateCheckOutQuery :one

-- name: FetchStudentGateLogs :many
SELECT
  gm.direction,
  gm.logged_at
FROM gate_management gm
LEFT JOIN student s ON s.id = gm.student_id
WHERE
  s.hospitality_id = $1;

-- name: HostelCheckInQuery :execrows
UPDATE hostel_check_in
SET
  checked_in_at = NOW(),
  checked_in_by = $2
WHERE
  accomodation_id = $1;
