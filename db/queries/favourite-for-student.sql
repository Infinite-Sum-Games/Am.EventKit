-- name: MarkFavouriteEventQuery :one
INSERT INTO favourites (
  email,
  event_id
) VALUES ($1, $2)
RETURNING 
  event_id;

-- name: UnmarkFavouriteEventQuery :one
DELETE FROM favourites
WHERE
  email = $1
  AND event_id = $2 
RETURNING
  event_id;
