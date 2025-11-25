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
SET 
  refresh_token = $1,
  updated_at = NOW()
WHERE 
  email = $2
  AND account_status = 'VERIFIED'
RETURNING
  refresh_token;

-- name: UpdateOrganizerRefreshTokenQuery :one
UPDATE organizer
SET 
  refresh_token = $1,
  updated_at = NOW()
WHERE
  email = $2
RETURNING
  refresh_token;

-- name: RevokeRefreshTokenQuery :one
UPDATE
	student
SET
	refresh_token = NULL,
  updated_at = NOW()
WHERE
	email = $1
	AND status = 'active'
RETURNING
	refresh_token;

-- name: GetStudentOtpQuery :one
SELECT 
  so.id,
  so.name,
  so.email,
  so.password,
  so.phone_number,
  so.is_amrita_student,
  so.amrita_roll_number,
  so.college_name,
  so.college_city,
  so.expiry_at
FROM student_onboarding so
WHERE 
  so.email = $1
  AND so.otp = $2
  AND so.expiry_at > NOW()
  AND NOT EXISTS (
    SELECT 1 FROM student s WHERE s.email = $1
  );

-- name: ResendStudentOtpQuery :one
SELECT
  so.name,
  so.email,
  so.otp,
  so.expiry_at
FROM student_onboarding so
WHERE
  so.email = $1
  AND so.expiry_at > NOW()
  AND NOT EXISTS (
    SELECT 1 FROM student s WHERE s.email = $1
);

-- name: UpdateStudentPasswordQuery :exec
UPDATE student
SET password = $2
WHERE email = $1;

-- name: OnboardStudentQuery :one
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: LoginUserQuery :one
SELECT
  id,
  name,
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

-- name: PasswordChangeOtpQuery :one
INSERT INTO password_reset (
  name,
  email,
  password,
  otp,
  expiry_at
) 
SELECT
  s.name, 
  s.email, 
  $2, 
  $3, 
  $4
FROM
  student s
WHERE 
  s.email = $1
  AND s.account_status = 'VERIFIED'
RETURNING 
  email, name;

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
  password = $1,
  updated_at = NOW()
WHERE
  email = $2
  AND account_status = 'VERIFIED'
RETURNING
  email;

-- name: UpsertStudentOnboardingQuery :one
INSERT INTO student_onboarding (
  name,
  email,
  password,
  phone_number,
  is_amrita_student,
  amrita_roll_number,
  college_name,
  college_city,
  otp,
  expiry_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (email)
DO UPDATE SET
  name = EXCLUDED.name,
  password = EXCLUDED.password,
  phone_number = EXCLUDED.phone_number,
  is_amrita_student = EXCLUDED.is_amrita_student,
  amrita_roll_number = EXCLUDED.amrita_roll_number,
  college_name = EXCLUDED.college_name,
  college_city = EXCLUDED.college_city,
  otp = EXCLUDED.otp,
  expiry_at = EXCLUDED.expiry_at
RETURNING
  id, 
  email;
