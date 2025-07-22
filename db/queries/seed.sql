-- name: InsertEventQuery :exec
INSERT INTO event(name, blurb, description, price, is_per_head, rules, 
  event_type, is_group, total_seats, seats_filled, event_status, event_mode, 
  attendance_mode) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: InsertOrganizerQuery :exec
INSERT INTO organizer(name, abbr, org_type, student_head, student_co_head, faculty_head)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: InsertEventToOrganizerMappingQuery :exec
INSERT INTO event_to_organizer_mapping(event_id, organizer_id)
VALUES ($1, $2);

-- name: InsertEventScheduleQuery :exec
INSERT INTO event_schedule(event_id, event_date, start_time, end_time, venue)
VALUES ($1, $2, $3, $4, $5);

-- name: InsertPeopleQuery :exec
INSERT INTO people(name, phone_number)
VALUES ($1, $2);

-- name: InsertPeopleToEventMappingQuery :exec
INSERT INTO people_to_event_mapping(event_id, person_id)
VALUES ($1, $2);

-- name: InsertTagsQuery :exec
INSERT INTO tags(name, abbreviation)
VALUES ($1, $2);

-- name: InsertEventTagMappingQuery :exec
INSERT INTO event_tag_mapping(tag_id, event_id)
VALUES ($1, $2);

-- name: ListEventsQuery :many
SELECT id, name, blurb, description, price, is_per_head, rules, event_type, is_group, 
       total_seats, seats_filled, event_status, event_mode, attendance_mode
FROM event;

-- name: ListOrganizersQuery :many
SELECT id, name, abbr, org_type, student_head, student_co_head, faculty_head
FROM organizer;

-- name: ListPeopleQuery :many
SELECT id, name, phone_number, profession, email
FROM people;

-- name: ListEventToOrganizerMappingQuery :many
SELECT id, event_id, organizer_id
FROM event_to_organizer_mapping;

-- name: ListEventScheduleQuery :many
SELECT id, event_id, event_date, start_time, end_time, venue
FROM event_schedule;

-- name: ListPeopleToEventMappingQuery :many
SELECT id, event_id, person_id
FROM people_to_event_mapping;

-- name: ListEventTagMappingQuery :many
SELECT id, tag_id, event_id
FROM event_tag_mapping;

-- name: TruncateAllTablesQuery :exec
TRUNCATE TABLE organizer, people, tags, event, event_to_organizer_mapping, event_schedule, people_to_event_mapping, event_tag_mapping CASCADE;
