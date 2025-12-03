-- name: FetchAllPeopleQuery :many
SELECT * FROM people;

-- name: FetchPeopleByDepartmentQuery :many
SELECT 
  p.name, 
  p.phone_number, 
  p.profession, 
  p.email
FROM people AS p
  INNER JOIN people_to_event_mapping AS ptem 
    ON p.id = ptem.person_id
  INNER JOIN event_to_organizer_mapping AS etom 
    ON ptem.event_id = etom.event_id
  INNER JOIN organizer AS o 
    ON etom.organizer_id = o.id
WHERE 
  o.id = $1;

-- name: FetchPeopleByEventQuery :many
SELECT
  p.name, 
  p.phone_number, 
  p.profession, 
  p.email
FROM people AS p
  INNER JOIN people_to_event_mapping AS ptem 
    ON p.id = ptem.person_id
  INNER JOIN event AS e
    ON ptem.event_id = e.id
WHERE 
  e.id = $1;

-- name: FetchPeopleByDayQuery :many
SELECT
  p.name, 
  p.phone_number, 
  p.profession, 
  p.email
FROM people AS p
  INNER JOIN people_to_event_mapping AS ptem
    ON p.id = ptem.person_id
WHERE 
  ptem.event_day && $1;

-- name: AddNewPersonQuery :one
INSERT INTO people (
  name, 
  phone_number, 
  profession, 
  email
) VALUES ($1, $2, $3, $4)
RETURNING
  id,
  name,
  phone_number,
  profession,
  email;

-- name: MapPersonToEventQuery :one
INSERT INTO people_to_event_mapping (
  event_id, 
  person_id, 
  event_day
) VALUES ($1, $2, $3)
RETURNING
  id,
  event_id,
  person_id,
  event_day;

-- name: UpdatePersonDetailsQuery :one
UPDATE people 
SET 
  name = $2, 
  phone_number = $3, 
  profession = $4,
  email = $5 
WHERE id = $1
RETURNING
  id,
  name,
  phone_number,
  profession,
  email;

-- name: DeletePersonQuery :one
DELETE FROM people 
WHERE id = $1
RETURNING id;
