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

-- name: DeleteHostelQuery :execrows
DELETE FROM hostel_metadata
WHERE id = $1;

-- name: AllotHostelQuery :execrows
UPDATE accomodation_details
SET 
  hostel_id = $2,
  payment_expires = NOW() + INTERVAL '30 minutes',
  payment_status = 'PENDING',
  updated_at = NOW()
  where accomodation_details.id = $1
  AND hostel_id IS NULL
  AND is_male = (
    SELECT is_male FROM hostel_metadata WHERE hostel_metadata.id = $2
  );

-- name: AffirmAccommodationPaymentQuery :execrows
UPDATE accomodation_details
SET 
  is_paid = TRUE,
  payment_status = 'COMPLETED',
  updated_at = NOW()
  where id = $1;

-- name: GetHostelQuery :one
SELECT
  id,
  room_count,
  is_male,
  warden_email,
  latitude,
  longtitude,
  map_url,
  hostel_name
FROM hostel_metadata
WHERE id = $1;

-- name: DecrementHostelRoomCountQuery :execrows
UPDATE hostel_metadata
SET 
  room_count = room_count - 1
  where id = $1
  AND room_count>0;

-- name: GetAccommodationByIdQuery :one
SELECT
  id AS accommodation_id,
  student_id,
  hostel_id,
  name,
  email,
  phone_number,
  is_male,
  room_preference,
  college_name,
  college_roll_number,
  is_hosteller,
  is_amrita_campus,
  check_in::date as check_in_date,
  to_char(check_in, 'HH12:MI AM') as check_in_time,
  check_out::date as check_out_date,
  to_char(check_out, 'HH12:MI AM') as check_out_time
FROM accomodation_details
WHERE id = $1;

-- name: GetAllHostelsQuery :many
SELECT
  id AS hostel_id,
  room_count,
  is_male,
  hostel_name
  FROM hostel_metadata;

-- name: UpdateAccommodationByIdQuery :execrows
UPDATE accomodation_details
SET
  is_male = $2,
  is_hosteller = $3,
  college_roll_number = $4,
  college_name = $5,
  room_preference = $6,
  is_amrita_campus = $7,
  check_in = $8,
  check_out = $9,
  updated_at = NOW()
WHERE id = $1;

-- name: MapQrStudentIdQuery :execrows
UPDATE student
SET hospitality_id = $2
WHERE id = $1;
