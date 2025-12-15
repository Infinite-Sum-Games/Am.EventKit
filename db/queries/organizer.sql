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
  name = $2,
  email = $3,
  password = $4,
  org_type = $5,
  student_head = $6,
  student_co_head = $7,
  faculty_head = $8
WHERE
  id = $1;

-- name: DeleteOrganizerByIDQuery :execrows
DELETE FROM organizer
WHERE id = $1;
