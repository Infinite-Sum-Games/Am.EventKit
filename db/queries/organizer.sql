-- name: ListOrganizersQuery :many

SELECT
  name,
  abbr,
  org_type,
  student_head,
  student_co_head,
  faculty_head
FROM organizer;
