-- name: LoginAdminQuery :one
SELECT
  id,
  name,
  email,
  password,
  refresh_token
FROM admin
WHERE email = $1;
