-- name: GateCheckInOutQuery :one
WITH student_lookup AS (
    SELECT id AS found_student_id
    FROM student
    WHERE hospitality_id = $1
)
INSERT INTO gate_management (
    student_id,
    direction,
    logged_at,
    personell_id
) 
SELECT 
    found_student_id, 
    $2,
    NOW(), 
    $3
FROM student_lookup
RETURNING direction;

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
