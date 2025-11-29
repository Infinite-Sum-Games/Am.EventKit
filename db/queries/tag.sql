-- name: ListTagsQuery :many
SELECT
    id,
    name,
    abbreviation
FROM tags;

-- name: CreateTagQuery :exec
INSERT INTO tags (name, abbreviation)
VALUES ($1, $2);

-- name: UpdateTagByIDQuery :execrows
UPDATE tags
SET
    name = $2,
    abbreviation = $3
WHERE
    id = $1;

-- name: DeleteTagByIDQuery :execrows
DELETE FROM tags
WHERE id = $1;