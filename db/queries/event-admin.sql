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
