-- name: ListEventsQuery :many
SELECT
    id,
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
    created_at,
    updated_at
FROM event;

-- name: GetEventByIdQuery :one
SELECT
    id,
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
    created_at,
    updated_at
FROM event
WHERE id = $1;
