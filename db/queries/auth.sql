-- name: CheckRefreshTokenQuery :one
SELECT refresh_token 
FROM student 
WHERE email = $1;

-- name: RevokeRefreshTokenQuery :one
UPDATE
	student
SET
	refresh_token = NULL
WHERE
	email = $1
	AND status = 'active'
RETURNING
	refresh_token;

-- name: UpsertStudentOnboarding :exec
INSERT INTO student_onboarding (
  name,
  department_name,
  email,
  password,
  phone_number,
  is_amrita_student,
  amrita_roll_number,
  college_name,
  college_city,
  academic_year,
  otp,
  expiry_at
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
ON CONFLICT (email)
DO UPDATE SET
  name = EXCLUDED.name,
  department_name = EXCLUDED.department_name,
  password = EXCLUDED.password,
  phone_number = EXCLUDED.phone_number,
  is_amrita_student = EXCLUDED.is_amrita_student,
  amrita_roll_number = EXCLUDED.amrita_roll_number,
  college_name = EXCLUDED.college_name,
  college_city = EXCLUDED.college_city,
  academic_year = EXCLUDED.academic_year,
  otp = EXCLUDED.otp,
  expiry_at = EXCLUDED.expiry_at;

-- name: UpdateStudentPassword :exec
UPDATE student
SET password = $2
WHERE email = $1;

-- name: UpdateStudentOTP :exec
UPDATE student_onboarding
SET otp = $2,
    expiry_at = $3
WHERE email = $1;

-- name: FinalizeStudentSignUp :exec
INSERT INTO student (
  name,
  department_name,
  email,
  password,
  phone_number,
  is_amrita_student,
  amrita_roll_number,
  college_name,
  college_city,
  academic_year,
  account_status
)
SELECT
  name,
  department_name,
  email,
  password,
  phone_number,
  is_amrita_student,
  amrita_roll_number,
  college_name,
  college_city,
  academic_year,
  'VERIFIED'
FROM student_onboarding
WHERE email = $1;
