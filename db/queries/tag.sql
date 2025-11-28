-- name: ListTagsQuery :many
SELECT
    id,
    name,
    abbreviation
FROM tags;

-- name: GetTagByIDQuery :one
SELECT 
    id,
    name,
    abbreviation
FROM tags
WHERE id = $1;

-- name: CreateTagQuery :exec
INSERT INTO tags (name, abbreviation)
VALUES ($1, $2);