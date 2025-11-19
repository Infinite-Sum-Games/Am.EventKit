-- name: FetchUserProfileQuery :one
SELECT 
  name, 
  department_name,
  email,
  phone_number,
  is_amrita_student,
  amrita_roll_number,
  college_name,
  college_city,
  academic_year
FROM student 
WHERE account_status = 'VERIFIED' and email = $1;

-- name: EditUserProfileQuery :execrows
UPDATE student
SET 
  name = $2,
  phone_number = $3,
  college_name = $4,
  college_city = $5
WHERE email = $1 AND account_status = 'VERIFIED';
