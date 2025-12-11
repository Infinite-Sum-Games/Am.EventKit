-- name: GetAllDiscounts :many
SELECT *
FROM discount
ORDER BY created_at DESC;

-- name: CreateDiscount :one
INSERT INTO discount (
  discount_type,
  discounted_solo_seats,
  discounted_team_seats,
  start_time,
  end_time
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: EditDiscount :one
UPDATE discount
SET
  discount_type = COALESCE($2, discount_type),
  discounted_solo_seats = COALESCE($3, discounted_solo_seats),
  discounted_team_seats = COALESCE($4, discounted_team_seats),
  start_time = COALESCE($5, start_time),
  end_time = COALESCE($6, end_time),
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDiscount :exec
DELETE FROM discount
WHERE id = $1;
