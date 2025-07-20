-- name: InsertEvent :exec
INSERT INTO event(name, blurb, description, price, is_per_head, rules, 
  event_type, is_group, total_seats, seats_filled, event_status, event_mode, 
  attendance_mode) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: InsertOrganizer :exec
INSERT INTO organizer(name, abbr, org_type, student_head, student_co_head, faculty_head)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: InsertEventToOrganizerMapping :exec
INSERT INTO event_to_organizer_mapping(event_id, organizer_id)
VALUES ($1, $2);

-- name: InsertEventSchedule :exec
INSERT INTO event_schedule(event_id, event_date, start_time, end_time, venue)
VALUES ($1, $2, $3, $4, $5);

-- name: InsertPeople :exec
INSERT INTO people(name, phone_number)
VALUES ($1, $2);

-- name: InsertPeopleToEventMapping :exec
INSERT INTO people_to_event_mapping(event_id, person_id)
VALUES ($1, $2);

-- name: InsertTags :exec
INSERT INTO tags(name, abbrevation)
VALUES ($1, $2);

-- name: InsertEventTagMapping :exec
INSERT INTO event_tag_mapping(tag_id, event_id)
VALUES ($1, $2);

-- name: ListEvents :many
SELECT id, name, blurb, description, price, is_per_head, rules, event_type, is_group, 
       total_seats, seats_filled, event_status, event_mode, attendance_mode
FROM event;

-- name: ListOrganizers :many
SELECT id, name, abbr, org_type, student_head, student_co_head, faculty_head
FROM organizer;

-- name: ListPeople :many
SELECT id, name, phone_number, profession, email
FROM people;

-- name: ListTags :many
SELECT id, name, abbrevation
FROM tags;
