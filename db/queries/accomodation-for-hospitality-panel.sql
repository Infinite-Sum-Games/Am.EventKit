-- name: GetAllAccommodationRequestsQuery :many
SELECT
  id,
  name,
  email,
  phone_number,
  college_name,
  college_roll_number,
  is_male,
  room_preference,
  is_paid,
  'RESERVED' AS check_in_status,
  'No Hostel Allotted' AS hostel,
  check_in::date as check_in_date,
  to_char(check_in, 'HH12:MI AM') as check_in_time,
  check_out::date as check_out_date,
  to_char(check_out, 'HH12:MI AM') as check_out_time
FROM accomodation_details
ORDER BY created_at DESC;
