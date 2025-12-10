-- name: LoginAdminQuery :one
SELECT
  id,
  name,
  email,
  password,
  refresh_token
FROM admin
WHERE email = $1;

-- name: UpdateAdminRefreshTokenQuery :one
UPDATE admin
SET
  refresh_token = $1,
  updated_at = NOW()
WHERE
  email = $2
RETURNING
  refresh_token;

-- name: RevokeAdminRefreshTokenQuery :one
UPDATE admin
SET
  refresh_token = NULL,
  updated_at = NOW()
WHERE
  email = $1
RETURNING
  email;

-- name: CheckAdminRefreshTokenQuery :one
SELECT refresh_token 
FROM admin
WHERE email = $1;
