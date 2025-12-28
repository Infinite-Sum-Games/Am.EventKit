-- name: FetchUserSessionQuery :one
SELECT 
  id,
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

-- name: FetchOrganizerSessionQuery :one
SELECT
  id,
  name,
  email
FROM
  organizer
WHERE
  email = $1;

-- name: FetchHospitalitySessionQuery :one
SELECT
  id,
  name,
  email
FROM
  organizer
WHERE
  email = $1;
