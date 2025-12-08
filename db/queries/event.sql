-- name: GetEventsQuery :many
SELECT
    e.id AS event_id,
    e.cover_image_url AS event_image_url,
    e.name AS event_name,
    e.event_status,
    e.blurb AS event_description,
    MIN(es.event_date) AS event_date,
    e.is_group,
    e.event_type,
    e.is_technical,

    COALESCE(
        JSONB_AGG(DISTINCT t.abbreviation) FILTER (WHERE t.id IS NOT NULL),
        '[]'::jsonb
    ) AS tags,

    e.price AS event_price,
    e.total_seats AS max_seats,
    e.seats_filled

FROM event e

LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id

GROUP BY e.id;

-- name: GetEventByIdQuery :one
SELECT
    e.id,
    e.name AS event_name,
    e.blurb,
    e.description AS event_description,
    e.cover_image_url,
    e.price,
    e.is_per_head,
    e.rules,
    e.event_type,
    e.is_group,
    e.max_teamsize,
    e.min_teamsize,
    e.total_seats,
    e.seats_filled,
    e.event_status,
    e.event_mode,
    e.is_technical,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'organizer_name', o.name,
        'org_abbreviation', LOWER(SUBSTRING(o.email FROM 1 FOR 3)),
        'org_type', o.org_type
      )) FILTER (WHERE o.id IS NOT NULL),
      '[]'::jsonb
    ) AS organizers,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'event_date', es.event_date,
        'start_time', es.start_time::time,
        'end_time', es.end_time::time,
        'venue', es.venue
      )) FILTER (WHERE es.id IS NOT NULL),
      '[]'::jsonb
    ) AS schedules,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'tag_name', t.name,
        'tag_abbreviation', t.abbreviation
      )) FILTER (WHERE t.id IS NOT NULL),
      '[]'::jsonb
    ) AS tags,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'person_name', p.name,
        'profession', p.profession,
        'phone_number', p.phone_number,
        'email', p.email
      )) FILTER (WHERE p.id IS NOT NULL),
      '[]'::jsonb
    ) AS people

FROM event e

LEFT JOIN event_to_organizer_mapping m ON e.id = m.event_id
LEFT JOIN organizer o ON m.organizer_id = o.id
LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id
LEFT JOIN people_to_event_mapping pem ON e.id = pem.event_id
LEFT JOIN people p ON pem.person_id = p.id

WHERE e.id = $1
GROUP BY e.id;

-- name: GetEventsWithAuthQuery :many
SELECT
    e.id AS event_id,
    e.cover_image_url AS event_image_url,
    e.name AS event_name,
    e.event_status,
    e.blurb AS event_description,
    MIN(es.event_date) AS event_date,
    e.is_group,
    e.event_type,
    e.is_technical,

    COALESCE(
        JSONB_AGG(DISTINCT t.abbreviation) FILTER (WHERE t.id IS NOT NULL),
        '[]'::jsonb
    ) AS tags,

    e.price AS event_price,
    e.total_seats AS max_seats,
    e.seats_filled,

    /* registration and favourite status for the given student */
    (COUNT(DISTINCT b.id) > 0 OR COUNT(DISTINCT tm.id) > 0) AS is_registered,
    (COUNT(DISTINCT f.id) > 0) AS is_starred

FROM event e

LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id
LEFT JOIN bookings b ON e.id = b.event_id AND b.student_id = $1
LEFT JOIN teams te ON te.event_id = e.id
LEFT JOIN team_members tm ON tm.team_id = te.id AND tm.student_id = $1
LEFT JOIN favourites f ON e.id = f.event_id AND f.email = $2

GROUP BY e.id;

-- name: GetEventByIdWithAuthQuery :one
SELECT
    e.id,
    e.name AS event_name,
    e.blurb,
    e.description AS event_description,
    e.cover_image_url,
    e.price,
    e.is_per_head,
    e.rules,
    e.event_type,
    e.is_group,
    e.max_teamsize,
    e.min_teamsize,
    e.total_seats,
    e.seats_filled,
    e.event_status,
    e.event_mode,
    e.is_technical,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'organizer_name', o.name,
        'org_abbreviation', LOWER(SUBSTRING(o.email FROM 1 FOR 3)),
        'org_type', o.org_type
      )) FILTER (WHERE o.id IS NOT NULL),
      '[]'::jsonb
    ) AS organizers,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'event_date', es.event_date,
        'start_time', es.start_time::time,
        'end_time', es.end_time::time,
        'venue', es.venue
      )) FILTER (WHERE es.id IS NOT NULL),
      '[]'::jsonb
    ) AS schedules,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'tag_name', t.name,
        'tag_abbreviation', t.abbreviation
      )) FILTER (WHERE t.id IS NOT NULL),
      '[]'::jsonb
    ) AS tags,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'person_name', p.name,
        'profession', p.profession,
        'phone_number', p.phone_number,
        'email', p.email
      )) FILTER (WHERE p.id IS NOT NULL),
      '[]'::jsonb
    ) AS people,

    (COUNT(DISTINCT b.id) > 0 OR COUNT(DISTINCT tm.id) > 0) AS is_registered,
    (COUNT(DISTINCT f.id) > 0) AS is_starred

FROM event e

LEFT JOIN event_to_organizer_mapping m ON e.id = m.event_id
LEFT JOIN organizer o ON m.organizer_id = o.id
LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id
LEFT JOIN people_to_event_mapping pem ON e.id = pem.event_id
LEFT JOIN people p ON pem.person_id = p.id
LEFT JOIN bookings b ON e.id = b.event_id AND b.student_id = $2
LEFT JOIN teams te ON te.event_id = e.id
LEFT JOIN team_members tm ON tm.team_id = te.id AND tm.student_id = $2
LEFT JOIN favourites f ON e.id = f.event_id AND f.email = $3

WHERE e.id = $1
GROUP BY e.id;

-- name: CreateEventQuery :one
INSERT INTO event (
  name,
  blurb,
  description,
  cover_image_url,
  price,
  is_per_head,
  rules,
  event_type,
  is_group,
  max_teamsize,
  min_teamsize,
  total_seats,
  seats_filled,
  event_status,
  event_mode,
  attendance_mode,
  is_technical
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
)
RETURNING id;

-- name: InsertEventScheduleQuery :exec
INSERT INTO event_schedule (
  event_id,
  event_date,
  start_time,
  end_time,
  venue
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: InsertEventTagMappingQuery :exec
INSERT INTO event_tag_mapping (
  tag_id,
  event_id
) VALUES (
  $1, $2
);

-- name: InsertEventOrganizerMappingQuery :exec
INSERT INTO event_to_organizer_mapping (
  event_id,
  organizer_id
) VALUES (
  $1, $2
);

-- name: InsertPeopleToEventMappingQuery :exec
INSERT INTO people_to_event_mapping (
  event_id,
  person_id
) VALUES (
  $1, $2
);

-- name: UpdateEventQuery :execrows
UPDATE event SET
  name = $2,
  blurb = $3,
  description = $4,
  cover_image_url = $5,
  price = $6,
  is_per_head = $7,
  rules = $8,
  event_type = $9,
  is_group = $10,
  max_teamsize = $11,
  min_teamsize = $12,
  total_seats = $13,
  seats_filled = $14,
  event_status = $15,
  event_mode = $16,
  attendance_mode = $17,
  is_technical = $18
WHERE id = $1;

-- name: DeleteEventSchedulesByEventIDQuery :exec
DELETE FROM event_schedule WHERE event_id = $1;

-- name: DeleteEventTagMappingsByEventIDQuery :exec
DELETE FROM event_tag_mapping WHERE event_id = $1;

-- name: DeleteEventOrganizerMappingsByEventIDQuery :exec
DELETE FROM event_to_organizer_mapping WHERE event_id = $1;

-- name: DeletePeopleToEventMappingsByEventIDQuery :exec
DELETE FROM people_to_event_mapping WHERE event_id = $1;

-- name: DeleteEventQuery :execrows
DELETE FROM event WHERE id = $1;

-- name: ToggleEventStatusQuery :execrows
UPDATE event
SET event_status = (
    CASE
        WHEN event_status = 'ACTIVE'::event_status_enum THEN 'CLOSED'::event_status_enum
        ELSE 'ACTIVE'::event_status_enum
    END
)
WHERE id = $1;
