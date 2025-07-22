-- name: ListOrganizersQuery :many

SELECT
  name,
  abbr,
  type,
  student_head,
  student_co_head,
  faculty_head
FROM organizer;
