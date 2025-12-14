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

-- name: AddEventDetailsQuery :one
UPDATE event
SET
  name = $1,
  blurb = $2,
  description = $3,
  rules = $4, 
  price = $5,
  is_per_head = $6,
  updated_at = NOW()
WHERE
  id = $7
RETURNING
  id,
  name,
  blurb,
  description,
  rules,
  price,
  is_per_head
  updated_at;

-- name: AddEventPosterQuery :one
UPDATE event
SET
  cover_image_url = $1,
  updated_at = NOW()
WHERE
  id = $2
RETURNING
  id,
  cover_image_url,
  updated_at;

-- name: DeleteEventPosterQuery :one
UPDATE event
SET 
  cover_image_url = NULL
WHERE 
  id = $1
RETURNING id;

-- name: AddEventDimensionQuery :one
UPDATE event
SET
  is_group = $1,
  total_seats = $2,
  min_teamsize = $3,
  max_teamsize = $4,
  updated_at = NOW()
WHERE
  id = $5
RETURNING
  is_group,
  total_seats,
  min_teamsize,
  max_teamsize,
  updated_at;

-- name: AddEventTogglesQuery :one
UPDATE event
SET
  event_type = $1,
  is_group = $2,
  event_status = $3,
  event_mode = $4,
  attendance_mode = $5,
  updated_at = NOW()
WHERE
  id = $6
RETURNING
  event_type,
  is_group,
  is_technical,
  event_status,
  event_mode,
  attendance_mode,
  updated_at;

-- name: ConnectEventAndOrganizerQuery :one
INSERT INTO event_to_organizer_mapping (
  event_id,
  organizer_id
) VALUES ($1, $2)
RETURNING
  event_id,
  organizer_id;

-- name: DisconnectEventAndOrganizerQuery :one
DELETE FROM event_to_organizer_mapping
WHERE 
  event_id = $1
  AND organizer_id = $2
RETURNING id;

-- name: ConnectEventAndTagsQuery :one
INSERT INTO event_tag_mapping (
  event_id,
  tag_id
) VALUES ($1, $2)
RETURNING
  event_id,
  tag_id;

-- name: DisconnectEventAndTagsQuery :one
DELETE FROM event_tag_mapping
WHERE 
  event_id = $1
  AND tag_id = $2
RETURNING
  event_id;

-- name: AddEventScheduleQuery :one
INSERT INTO event_schedule (
  event_id,
  event_date,
  start_time,
  end_time,
  venue
) VALUES ($1, $2, $3, $4, $5)
RETURNING
  id,
  event_id,
  event_date,
  start_time,
  end_time;

-- name: EditEventScheduleQuery :one
UPDATE event_schedule
SET
  event_date = $1,
  start_time = $2,
  end_time = $3,
  venue = $4,
  updated_at = NOW()
WHERE
  id = $5
RETURNING
  id,
  event_id,
  event_date,
  start_time,
  end_time,
  updated_at;

-- name: DeleteEventScheduleByIdQuery :one
DELETE FROM event_schedule
WHERE 
  id = $1
RETURNING id;

-- name: ConnectEventAndPeopleQuery :one
INSERT INTO people_to_event_mapping (
  event_id,
  person_id
) VALUES ($1, $2)
RETURNING
  event_id,
  person_id;

-- name: DisconnectEventAndPeopleQuery :one
DELETE FROM people_to_event_mapping
WHERE 
  event_id = $1
  AND person_id = $2
RETURNING id;

-- name: PublishEventQuery :one
UPDATE event
SET
  event_status = 'ACTIVE',
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  event_status,
  updated_at;

-- name: UnpublishEventQuery :one
UPDATE event
SET
  event_status = 'CLOSED',
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  event_status,
  updated_at;

-- name: MarkEventAsCompletedQuery :one
UPDATE event
SET
  event_status = 'COMPLETED',
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  event_status,
  updated_at;

-- name: UnmarkEventsAsCompletedQuery :one
UPDATE event
SET
  event_status = 'ACTIVE',
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  event_status,
  updated_at;
