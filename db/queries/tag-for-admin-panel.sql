-- name: ListTagsQuery :many
SELECT
    t.id,
    t.name,
    t.abbreviation,
    COALESCE(
      jsonb_agg(json_build_object(
          'id', e.id,
          'name', e.name
      )) FILTER (WHERE e.id IS NOT NULL),
    '[]'::jsonb
    ) AS events
FROM tags AS t
LEFT JOIN event_tag_mapping etm ON t.id = etm.tag_id
LEFT JOIN event e ON e.id = etm.event_id
WHERE
  t.abbreviation NOT LIKE '!%'
GROUP BY
  t.id, 
  t.name, 
  t.abbreviation;

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
