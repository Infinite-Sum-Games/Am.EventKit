-- name: FetchAllPeopleQuery :many
SELECT * FROM people;

-- name: FetchPeopleByDepartmentQuery :many
SELECT * FROM people 
INNER JOIN people_to_event_mapping ON people.id = people_to_event_mapping.person_id
INNER JOIN event_to_organizer_mapping ON people_to_event_mapping.event_id = event_to_organizer_mapping.event_id
INNER JOIN organizer ON event_to_organizer_mapping.organizer_id = organizer.id
WHERE organizer.name = $1;
-- name: FetchPeopleByEventQuery :many
SELECT * FROM people
INNER JOIN people_to_event_mapping ON people.id = people_to_event_mapping.person_id
INNER JOIN event ON people_to_event_mapping.event_id = event.id
WHERE event.name = $1;

-- name: FetchPeopleByDayQuery :many
SELECT * FROM people WHERE ARRAY_CONTAINS(event_day, $1);

-- name: AddNewPersonQuery :one
INSERT INTO people (
  name, 
  phone_number, 
  profession, 
  email) VALUES ($1, $2, $3, $4)
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
  event_day) VALUES ($1, $2, $3)
RETURNING
id,
event_id,
person_id,
event_day;

-- name: UpdatePersonDetailsQuery :one
UPDATE people SET 
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
DELETE FROM people WHERE id = $1
RETURNING id;
