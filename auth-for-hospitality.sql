-- name: LoginHospitalityQuery :one
SELECT
  id,
  name,
  email,
  password,
  refresh_token
FROM hospitality_personell
WHERE email = $1;

-- name: UpdateHospitalityRefreshTokenQuery :one
UPDATE hospitality_personell
SET
  refresh_token = $1,
  updated_at = NOW()
WHERE
  email = $2
RETURNING
  refresh_token;

-- name: RevokeHospitalityRefreshTokenQuery :one
UPDATE hospitality_personell
SET
  refresh_token = NULL,
  updated_at = NOW()
WHERE
  email = $1
RETURNING
  email;

-- name: CheckHospitalityRefreshTokenQuery :one
SELECT refresh_token
FROM hospitality_personell
WHERE email = $1;
