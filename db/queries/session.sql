-- name: FetchUserSessionQuery :one
SELECT 
  name,
  email
FROM
  student
WHERE
  email = $1
  AND account_status = 'VERIFIED';
