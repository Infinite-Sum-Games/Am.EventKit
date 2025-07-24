-- name: ListEventsQuery :many
SELECT
    e.id,
    e.name,
    e.blurb,
    e.description,
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
    e.attendance_mode,
    e.created_at,
    e.updated_at,

    -- Organizer details
    JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', o.id,
        'name', o.name,
        'abbr', o.abbr,
        'org_type', o.org_type,
        'student_head', o.student_head,
        'student_co_head', o.student_co_head,
        'faculty_head', o.faculty_head
    )) FILTER (WHERE o.id IS NOT NULL) AS organizers,

    -- Event schedule
    JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', es.id,
        'event_id', es.event_id,
        'event_date', es.event_date,
        'start_time', es.start_time,
        'end_time', es.end_time,
        'venue', es.venue
    )) FILTER (WHERE es.id IS NOT NULL) AS schedules,

    -- Tags 
    JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', t.id,
        'name', t.name,
        'abbreviation', t.abbreviation
    )) FILTER (WHERE t.id IS NOT NULL) AS tags

FROM event e
LEFT JOIN event_to_organizer_mapping m ON e.id = m.event_id
LEFT JOIN organizer o ON m.organizer_id = o.id
LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id
WHERE ($1 IS NULL OR m.organizer_id = $1)
GROUP BY e.id;

-- name: GetEventByIdQuery :one
SELECT
    e.id,
    e.name,
    e.blurb,
    e.description,
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
    e.attendance_mode,
    e.created_at,
    e.updated_at,

    -- Organizer details
    JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', o.id,
        'name', o.name,
        'abbr', o.abbr,
        'org_type', o.org_type,
        'student_head', o.student_head,
        'student_co_head', o.student_co_head,
        'faculty_head', o.faculty_head
    )) FILTER (WHERE o.id IS NOT NULL) AS organizers,

    -- Event schedule
    JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', es.id,
        'event_date', es.event_date,
        'start_time', es.start_time,
        'end_time', es.end_time,
        'venue', es.venue
    )) FILTER (WHERE es.id IS NOT NULL) AS schedules,

    -- Tags 
    JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
        'id', t.id,
        'name', t.name
    )) FILTER (WHERE t.id IS NOT NULL) AS tags

FROM event e
LEFT JOIN event_to_organizer_mapping m ON e.id = m.event_id
LEFT JOIN organizer o ON m.organizer_id = o.id
LEFT JOIN event_schedule es ON e.id = es.event_id
LEFT JOIN event_tag_mapping etm ON e.id = etm.event_id
LEFT JOIN tags t ON etm.tag_id = t.id
WHERE e.id = $1
GROUP BY e.id;
