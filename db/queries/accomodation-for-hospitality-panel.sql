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
  payment_status,
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
  map_url = $6,
  is_male = $7
  where id = $1;

-- name: DeleteHostelQuery :execrows
DELETE FROM hostel_metadata
WHERE id = $1;

-- name: AllotHostelQuery :execrows
UPDATE accomodation_details
SET 
  hostel_id = $2,
  day_count = $3,
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

-- name: MapQrStudentIdQuery :one
UPDATE student
SET hospitality_id = $2
WHERE student.id = $1
RETURNING
  (SELECT id FROM accomodation_details WHERE student_id = $1 LIMIT 1) AS accommodation_id,
  EXISTS (
    SELECT 1
    FROM accomodation_details
    WHERE student_id = $1
  ) AS has_opted_accommodation;

-- name: GetFinanceDetailsByHospitalityIdQuery :one
SELECT
  ad.id AS accommodation_id,
  ad.name AS name,
  ad.email AS email,
  ad.day_count AS day_count,
  ad.payment_status AS payment_status,
  ad.is_amrita_campus AS is_amrita_campus,
  ad.is_hosteller AS is_hosteller,
  hm.hostel_name AS hostel_name
FROM accomodation_details ad
INNER JOIN student s ON s.id = ad.student_id
INNER JOIN hostel_metadata hm ON hm.id = ad.hostel_id
WHERE s.hospitality_id = $1;

-- name: GetStudentDetailsForSecurityQuery :one
SELECT 
    s.name AS student_name,
    s.email AS student_email,
    s.phone_number AS student_phone_number,
    s.college_name AS college_name,

    CASE
        WHEN ad.student_id IS NOT NULL THEN ad.college_roll_number
        ELSE NULL
    END AS college_roll_number,

    CASE 
        WHEN ad.student_id IS NOT NULL THEN ad.check_in::date
        ELSE NULL
    END AS check_in_date,

    CASE 
        WHEN ad.student_id IS NOT NULL THEN to_char(ad.check_in, 'HH12:MI AM')
        ELSE NULL
    END AS check_in_time,

    CASE 
        WHEN ad.student_id IS NOT NULL THEN ad.check_out::date
        ELSE NULL
    END AS check_out_date,

    CASE 
        WHEN ad.student_id IS NOT NULL THEN to_char(ad.check_out, 'HH12:MI AM')
        ELSE NULL
    END AS check_out_time,

    CASE 
        WHEN ad.student_id IS NOT NULL THEN hm.hostel_name
        ELSE NULL
    END AS hostel_name

FROM student s
LEFT JOIN accomodation_details ad 
    ON ad.student_id = s.id
LEFT JOIN hostel_metadata hm 
    ON hm.id = ad.hostel_id
WHERE s.hospitality_id = $1;

-- name: GetAllHostelDetailsQuery :many
SELECT
  hm.id AS hostel_id,
  hm.hostel_name AS hostel_name,
  hm.room_count AS room_count,
  hm.is_male AS is_male,
  hm.latitude AS latitude,
  hm.longtitude AS longtitude,
  hm.map_url AS map_url,
  hm.warden_email AS warden_email,
  hm.room_filled AS room_filled
FROM hostel_metadata hm;
  
