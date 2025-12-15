-- name: GetAllAdminEventsQuery :many
SELECT
  id,
  cover_image_url as poster_url,
  name,
  blurb,
  event_type,
  event_status,
  price,
  is_per_head,
  is_group,
  is_technical,
  seats_filled,
  total_seats,
  updated_at
FROM
  event;

-- name: GetAdminEventByEventIdQuery :one
SELECT
  e.id,
  e.name,
  e.blurb,
  e.cover_image_url as poster_url,
  e.event_type,
  e.event_mode,
  e.event_status,
  e.description,
  e.attendance_mode,
  e.is_group,
  e.min_teamsize,
  e.max_teamsize,
  e.is_per_head,
  e.price,
  e.rules,
  e.seats_filled,
  e.total_seats,
  e.is_technical,
  e.updated_at,

  COALESCE(
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', o.id,
        'name', o.name,
        'org_type', o.org_type
    )) FILTER (WHERE o.id IS NOT NULL),
  '[]'::jsonb
  ) AS organizers,

  COALESCE(
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', es.id,
        'event_date', es.event_date,
        'start_time', es.start_time,
        'end_time', es.end_time,
        'venue', es.venue,
        'updated_at', es.updated_at
    )) FILTER (WHERE es.id IS NOT NULL),
  '[]'::jsonb
  ) AS schedules,

  COALESCE(
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', p.id,
        'name', p.name,
        'profession', p.profession
    )) FILTER (WHERE o.id IS NOT NULL),
  '[]'::jsonb
  ) AS people,

  COALESCE(
    JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', t.id,
        'name', t.name,
        'abbreviation', t.abbr
    )) FILTER (WHERE t.id IS NOT NULL),
  '[]'::jsonb
  ) AS tags

FROM event e

LEFT JOIN event_to_organizer_mapping m ON e.id = m.event_id
LEFT JOIN organizer o ON m.organizer_id = o.id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON e.id = m.event_id
LEFT JOIN people_to_event_mapping pem ON e.id = pem.event_id
LEFT JOIN people p ON pem.person_id = p.id

WHERE
  e.id = $1
GROUP BY e.id;

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
  event_mode = $2,
  attendance_mode = $3,
  is_technical = $4,
  event_status = $5,
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
  venue,
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
  venue,
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
