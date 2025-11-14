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
  department_name = $3,
  phone_number = $4,
  college_name = $5,
  college_city = $6,
  academic_year = $7
WHERE email = $1 AND account_status = 'VERIFIED';
