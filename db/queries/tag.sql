-- name: ListTagsQuery :many
SELECT
    id,
    name,
    abbreviation
FROM tags;
