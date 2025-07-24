-- name: ListOrganizersQuery :many

SELECT
  name as organizer_name,
  abbr as abbreviation,
  org_type as organizer_type,
  student_head,
  student_co_head,
  faculty_head
FROM organizer;
