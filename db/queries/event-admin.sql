-- name: NewUntitledEventQuery :one
INSERT INTO event (
  name,
  blurb,
  description,
  price,
  is_per_head,
  rules,
  event_type,
  is_group,
  max_teamsize,
  min_teamsize,
  total_seats,
  event_status,
  event_mode,
  attendance_mode,
  is_technical
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING 
  id,
  name,
  blurb,
  description,
  cover_image_url,
  price,
  is_per_head,
  rules,
  is_group,
  min_teamsize,
  max_teamsize,
  total_seats,
  event_type,
  is_technical,
  event_mode,
  attendance_mode,
  event_status;

-- name: PublishEventQuery :one
UPDATE event
SET
  event_status = 'ACTIVE'
WHERE
  id = $1
RETURNING
  event_status;

-- name: DeleteEventPosterQuery :one
UPDATE event
SET 
  cover_image_url = NULL
WHERE id = $1
RETURNING id;

-- name: DisconnectEventAndTags :one
DELETE FROM event_tag_mapping
WHERE 
  event_id = $1
  AND tag_id = $2
RETURNING
  event_id;

-- name: DisconnectEventAndOrganizerQuery :exec
DELETE FROM event_to_organizer_mapping
WHERE 
  event_id = $1
  AND organizer_id = $2
RETURNING id;

-- name: DisconnectEventAndPeopleQuery :one
DELETE FROM people_to_event_mapping
WHERE 
  event_id = $1
  AND person_id = $2
RETURNING id;

-- name: DeleteEventScheduleByIdQuery :one
DELETE FROM event_schedule
WHERE id = $1
RETURNING id;

-- name: UnpublishEventQuery :one
UPDATE event
SET
  event_status = 'CLOSED'
WHERE
  id = $1
RETURNING
  event_status;

-- name: MarkEventAsCompletedQuery :one
UPDATE event
SET
  event_status = 'COMPLETED'
WHERE
  id = $1
RETURNING
  event_status;

-- name: UnmarkEventsAsCompletedQuery :one
UPDATE event
SET
  event_status = 'ACTIVE'
WHERE
  id = $1
RETURNING
  event_status;
