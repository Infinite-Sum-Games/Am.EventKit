-- name: FetchUserSessionQuery :one
SELECT 
  name,
  email
FROM
  student
WHERE
  email = $1
  AND account_status = 'VERIFIED';

-- name: FetchAdminSessionQuery :one
SELECT
  name,
  email
FROM
  admin
WHERE
  email = $1;
