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
