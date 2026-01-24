-- name: FetchUnclaimedBedsQuery :many
SELECT
  ad.id,
  hm.hostel_name,
  s.name as student_name,
  s.email as student_email,
  ad.college_name,
  s.phone_number
FROM accomodation_details ad
LEFT JOIN hostel_metadata hm ON ad.hostel_id = hm.id
LEFT JOIN student s ON ad.student_id = s.id
WHERE
  ad.payment_status = 'PENDING'
  AND hm.id IS NOT NULL;

-- name: DeleteUnclaimedBedQuery :execrows
WITH updated_accommodation AS (
    UPDATE accomodation_details
    SET 
      payment_status = 'FAILED',
      hostel_id = NULL
    WHERE accomodation_details.id = $1
    RETURNING (SELECT hostel_id FROM accomodation_details WHERE id = $1) AS old_hostel_id
) UPDATE hostel_metadata hm
SET 
  room_filled = room_filled - 1
FROM updated_accommodation ua
WHERE 
  hm.id = ua.old_hostel_id;
