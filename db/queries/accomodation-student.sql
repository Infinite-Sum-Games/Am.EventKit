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
  is_amrita_campus
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;
