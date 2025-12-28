-- name: CheckStudentHasTicketQuery :one
SELECT EXISTS (
    SELECT 1
    FROM bookings b 
    WHERE b.student_id = $1
      AND b.txn_status = 'SUCCESS'
      
    UNION ALL
    
    SELECT 1
    FROM team_members tm
    JOIN teams t ON t.id = tm.team_id
    JOIN bookings b ON b.id = t.booking_id
    WHERE tm.student_id = $1
      AND b.txn_status = 'SUCCESS'
);

-- name: CheckUserAccomodationExistsQuery :one
SELECT 
    CASE 
        WHEN EXISTS (
          SELECT 1 
          FROM accomodation_details
          WHERE student_id = $1
        ) 
        THEN TRUE 
        ELSE FALSE 
    END AS has_accomodation;

-- name: FetchStudentMetadata :one
SELECT
  id,
  name,
  phone_number
FROM student
WHERE
  email = $1;

-- name: InsertAccomodationFormEntryQuery :one
INSERT INTO accomodation_details (
  student_id,
  name,
  is_male,
  email,
  phone_number,
  is_hosteller,
  college_roll_number,
  college_name,
  room_preference,
  is_amrita_campus,
  check_in,
  check_out
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;
