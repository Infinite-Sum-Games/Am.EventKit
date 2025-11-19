-- name: GetEventsQuery :many
SELECT
    e.id AS event_id,
    e.cover_image_url AS event_image_url,
    e.name AS event_name,
    e.event_status,
    e.blurb AS event_description,
    MIN(es.event_date) AS event_date,
    e.is_group,
    COALESCE(
        JSONB_AGG(DISTINCT t.name) FILTER (WHERE t.id IS NOT NULL),
        '[]'::jsonb
    ) AS tags,
    e.price AS event_price,
    FALSE AS is_registered,
    FALSE AS is_starred,
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
    COALESCE(
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'organizer_name', o.name,
        'org_abbreviation', o.abbr,
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
      JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'tag_name', t.name,
        'tag_abbreviation', t.abbreviation
      )) FILTER (WHERE t.id IS NOT NULL),
      '[]'::jsonb
    ) AS tags
FROM event e
LEFT JOIN event_to_organizer_mapping m ON e.id = m.event_id
LEFT JOIN organizer o ON m.organizer_id = o.id
LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id
WHERE e.id = $1
GROUP BY e.id;
