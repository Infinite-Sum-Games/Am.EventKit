-- name: FindEmailQuery :one
SELECT EXISTS (
  SELECT 1 
  FROM student
  WHERE email = $1
);

-- name: CheckRefreshTokenQuery :one
SELECT refresh_token 
FROM student 
WHERE email = $1;

-- name: UpdateRefreshTokenQuery :one
UPDATE student 
SET refresh_token = $1 
WHERE 
  email = $2
  AND account_status = 'VERIFIED'
RETURNING
  refresh_token;

-- name: UpdateOrganizerRefreshTokenQuery :one
UPDATE organizer
SET refresh_token = $1
WHERE
  email = $2
RETURNING
  refresh_token;

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

-- name: GetStudentOtpQuery :one
SELECT 
  id,
  name,
  otp
FROM student_onboarding 
WHERE 
  email = $1
  AND expiry_at > NOW();

-- name: UpdateStudentPasswordQuery :exec
UPDATE student
SET password = $2
WHERE email = $1;

-- name: FinalizeStudentSignUpQuery :one
INSERT INTO student (
  name,
  email,
  password,
  phone_number,
  is_amrita_student,
  amrita_roll_number,
  college_name,
  college_city,
  account_status
)
SELECT
  name,
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
WHERE student_onboarding.email = $1
RETURNING id;

-- name: LoginUserQuery :one
SELECT
  id,
  email,
  refresh_token
FROM
  student
WHERE
  email = $1
  AND password = $2
  AND account_status = 'VERIFIED';

-- name: LoginOrganizerQuery :one
SELECT
  id,
  email,
  refresh_token
FROM
  organizer
WHERE
  email = $1
  AND password = $2;

-- INSERT INTO password_reset (
--   email,
--   password,
--   otp,
--   expiry_at
-- ) VALUES ($1, $2, $3, $4)
-- RETURNING
--   email;

-- name: PasswordChangeOtpQuery :one
INSERT INTO password_reset (
  email,
  password,
  otp,
  expiry_at
) SELECT
    $1, $2, $3, $4
WHERE EXISTS (
    SELECT 1
    FROM 
      student s
    WHERE 
      s.email = $1
      AND s.account_status = 'VERIFIED'
) RETURNING email;

-- name: PasswordChangeVerifyOtpQuery :one
SELECT 
  email, 
  password
FROM
  password_reset
WHERE
  email = $1
  AND otp = $2
  AND expiry_at > NOW();

-- name: UpdateUserPasswordQuery :one
UPDATE student
SET
  password = $1
WHERE
  email = $2
  AND account_status = 'VERIFIED'
RETURNING
  email;
