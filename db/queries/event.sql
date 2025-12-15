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
      array_agg(DISTINCT t.abbreviation) FILTER (WHERE t.id IS NOT NULL)
    ) AS tags,

    e.price AS event_price,
    (e.seats_filled = e.total_seats) AS is_full

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
    (e.seats_filled = e.total_seats) AS is_full,
    e.is_per_head,
    e.rules,
    e.event_type,
    e.is_group,
    e.max_teamsize,
    e.min_teamsize,
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
        'start_time', es.start_time,
        'end_time', es.end_time,
        'venue', es.venue
      )) FILTER (WHERE es.id IS NOT NULL),
      '[]'::jsonb
    ) AS schedules,

    COALESCE(
      array_agg(DISTINCT t.name) FILTER (WHERE t.id IS NOT NULL)
    ) AS tags,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'person_name', p.name,
        'profession', p.profession
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

WHERE 
  e.id = $1
  AND e.event_status = 'ACTIVE'
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
    (e.seats_filled = e.total_seats) AS is_full,

    /* registration and favourite status for the given student */
    (COUNT(DISTINCT b.id) > 0 OR COUNT(DISTINCT tm.id) > 0) AS is_registered,
    (COUNT(DISTINCT f.id) > 0) AS is_starred

FROM event e

LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id
LEFT JOIN bookings b ON e.id = b.event_id AND b.student_id = $1 AND b.txn_status = 'SUCCESS'
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
    (e.seats_filled = e.total_seats) AS is_full,
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
        'start_time', es.start_time,
        'end_time', es.end_time,
        'venue', es.venue
      )) FILTER (WHERE es.id IS NOT NULL),
      '[]'::jsonb
    ) AS schedules,

    COALESCE(
      array_agg(DISTINCT t.name) FILTER (WHERE t.id IS NOT NULL)
    ) AS tags,

    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'person_name', p.name,
        'profession', p.profession
      )) FILTER (WHERE p.id IS NOT NULL),
      '[]'::jsonb
    ) AS people,

    (COUNT(DISTINCT b.id) > 0 OR COUNT(DISTINCT tm_user.id) > 0) AS is_registered,
    (COUNT(DISTINCT f.id) > 0) AS is_starred

FROM event e

LEFT JOIN event_to_organizer_mapping m ON e.id = m.event_id
LEFT JOIN organizer o ON m.organizer_id = o.id
LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id
LEFT JOIN people_to_event_mapping pem ON e.id = pem.event_id
LEFT JOIN people p ON pem.person_id = p.id
LEFT JOIN bookings b ON e.id = b.event_id AND b.student_id = $2 AND b.txn_status = 'SUCCESS'
LEFT JOIN teams te ON te.event_id = e.id
LEFT JOIN team_members tm_user ON tm_user.team_id = te.id AND tm_user.student_id = $2
LEFT JOIN favourites f ON e.id = f.event_id AND f.email = $3
WHERE e.id = $1
GROUP BY e.id;
