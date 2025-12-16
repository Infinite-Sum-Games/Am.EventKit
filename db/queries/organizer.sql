-- name: ListOrganizersQuery :many
SELECT
  id,
  name,
  email,
  org_type,
  student_head,
  student_co_head,
  faculty_head
FROM organizer;

-- name: CreateOrganizerQuery :exec
INSERT INTO organizer (
  name,
  email,
  password,
  org_type,
  student_head,
  student_co_head,
  faculty_head
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7
);

-- name: UpdateOrganizerByIDQuery :execrows
UPDATE organizer
SET
  name = $1,
  email = $2,
  org_type = $3,
  student_head = $4,
  student_co_head = $5,
  faculty_head = $6
WHERE
  id = $7;

-- name: DeleteOrganizerByIDQuery :execrows
DELETE FROM organizer
WHERE id = $1;

-- name: ChangeOrganizerPasswordQuery :one
UPDATE organizer
SET 
  password = $1
WHERE
  id = $2
RETURNING id;
