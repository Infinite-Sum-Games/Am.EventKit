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

-- name: HostelCheckInQuery :one
WITH student_lookup AS (
  SELECT id
  FROM student
  WHERE student.hospitality_id = $1
),
accommodation_lookup AS (
  SELECT id, hostel_id, student_id
  FROM accomodation_details
  WHERE 
    student_id = (SELECT id FROM student_lookup)
    AND payment_status = 'COMPLETED'
),
new_check_in AS (
  INSERT INTO hostel_check_in (accomodation_id, checked_in_by)
  SELECT id, $2
  FROM accommodation_lookup
  RETURNING accomodation_id
)
SELECT
  s.name,
  hm.hostel_name
FROM new_check_in
JOIN accomodation_details ad ON ad.id = new_check_in.accomodation_id
JOIN hostel_metadata hm ON hm.id = ad.hostel_id
JOIN student s ON s.id = ad.student_id;

-- name: HostelGateStatusQuery :one
WITH student_info AS (
    SELECT
        s.id as student_id,
        s.name as student_name,
        s.email as student_email,
        EXISTS (
            SELECT 1
            FROM accomodation_details ad
            WHERE ad.student_id = s.id AND ad.payment_status = 'COMPLETED'
        ) as has_accommodation
    FROM student s
    WHERE s.hospitality_id = $1
),
gate_logs AS (
    SELECT
        direction,
        logged_at
    FROM gate_management
    WHERE student_id = (SELECT student_id FROM student_info)
    ORDER BY logged_at DESC
),
last_check_in AS (
    SELECT logged_at FROM gate_logs WHERE direction = 'IN' LIMIT 1
),
last_check_out AS (
    SELECT logged_at FROM gate_logs WHERE direction = 'OUT' LIMIT 1
)
SELECT
    si.student_name,
    si.student_email,
    si.has_accommodation AS single_check_in,
    (SELECT logged_at FROM last_check_in) AS last_check_in,
    (SELECT logged_at FROM last_check_out) AS last_check_out
FROM student_info si;

-- name: GateCheckStatusQuery :one  
SELECT
  s.name,
  s.college_name,
  s.email,
  ad.payment_status AS accomodation_status,
  (
    SELECT logged_at
    FROM gate_management
    WHERE 
      student_id = s.id 
      AND direction = 'IN'
    ORDER BY logged_at DESC
    LIMIT 1
  ) AS last_check_in,
  (
    SELECT logged_at
    FROM gate_management
    WHERE
      student_id = s.id 
      AND direction = 'OUT'
    ORDER BY logged_at DESC
    LIMIT 1
  ) AS last_check_out
FROM student s
LEFT JOIN accomodation_details ad ON s.id = ad.student_id
WHERE
  s.hospitality_id = $1;
