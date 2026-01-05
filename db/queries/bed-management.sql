-- name: FetchUnclaimedBedsQuery :many
SELECT
  hm.id,
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
  AND hm.id IS NOT NULL
  AND ad.updated_at > NOW() - INTERVAL '30 minutes';

-- name: DeleteUnclaimedBedQuery :execrows
WITH updated_accommodation AS (
    UPDATE accomodation_details AS ad
    SET 
      ad.payment_status = 'FAILED'
    WHERE ad.id = $1
    RETURNING ad.hostel_id
) UPDATE hostel_metadata hm
SET 
  room_filled = room_filled - 1
WHERE 
  hm.id = (SELECT hostel_id FROM updated_accommodation);
