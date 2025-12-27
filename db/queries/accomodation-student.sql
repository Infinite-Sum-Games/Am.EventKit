-- name: CheckUserAccomodationExistsQuery :one
SELECT 
    CASE 
        WHEN EXISTS (
          SELECT 1 
          FROM accomodation_form_resp
          WHERE student_id = $1
        ) 
        THEN TRUE 
        ELSE FALSE 
    END AS has_accomodation;
