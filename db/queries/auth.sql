-- name: RevokeRefreshTokenQuery :one
UPDATE
	student
SET
	refresh_token = NULL
WHERE
	email = $1
	AND status = 'active'
RETURNING
	refresh_token;
