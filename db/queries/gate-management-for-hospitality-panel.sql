-- name: GetLiveGateLogsQuery :many
SELECT
  s.name AS student_name,
  s.email AS student_email,
  s.college_name,
  gm.direction,
  gm.logged_at,
  per.name AS personell_name
FROM gate_management gm
LEFT JOIN student s ON gm.student_id = s.id
LEFT JOIN accomodation_personell per ON gm.personell_id = per.id
ORDER BY gm.logged_at DESC;

-- name: GetLiveHostelCheckInQuery :many
SELECT
  s.name AS student_name,
  s.email AS student_email,
  s.college_name,
  hm.hostel_name,
  h.checked_in_at AS logged_at,
  per.name AS personell_name
FROM hostel_check_in h
LEFT JOIN accomodation_personell per ON h.checked_in_by = per.id
LEFT JOIN accomodation_details ad ON ad.id = h.accomodation_id
LEFT JOIN hostel_metadata hm ON ad.hostel_id = hm.id
LEFT JOIN student s ON ad.student_id = s.id
ORDER BY h.checked_in_at DESC;
