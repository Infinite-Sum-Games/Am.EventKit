-- name: LoginOrganizerQuery :one
SELECT
  id,
  email,
  password,
  refresh_token
FROM
  organizer
WHERE
  email = $1;

-- name: UpdateOrganizerRefreshTokenQuery :one
UPDATE organizer
SET 
  refresh_token = $1,
  updated_at = NOW()
WHERE
  email = $2
RETURNING
  refresh_token;

-- name: RevokeOrganizerRefreshTokenQuery :one
UPDATE organizer
SET
  refresh_token = NULL,
  updated_at = NOW()
WHERE
  email = $1
RETURNING
  email;

-- name: CheckOrganizerRefreshTokenQuery :one
SELECT refresh_token 
FROM organizer
WHERE email = $1;
