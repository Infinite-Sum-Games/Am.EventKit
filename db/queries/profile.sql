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
