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

-- name: AddHostelQuery :one
INSERT INTO hostel_metadata (
  room_count,
  is_male,
  warden_email,
  latitude,
  longtitude,
  map_url,
  hostel_name)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdateHostelQuery :execrows
UPDATE hostel_metadata
SET 
  room_count = $2,
  warden_email = $3,
  latitude = $4,
  longtitude = $5,
  map_url = $6
  where id = $1;
