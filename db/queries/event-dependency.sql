-- name: AddEventDependency :exec
INSERT INTO event_dependency_mapping (
  start_event, 
  end_event)
VALUES ($1, $2);

-- name: RemoveEventDependency :execrows
DELETE FROM event_dependency_mapping
WHERE start_event = $1 AND end_event = $2;

-- name: ListAllEventDependencies :many
SELECT start_event, end_event
FROM event_dependency_mapping;

-- name: ListDependenciesForEvent :many
SELECT start_event
FROM event_dependency_mapping
WHERE end_event = $1;

-- name: ListDependentsForEvent :many
SELECT end_event
FROM event_dependency_mapping
WHERE start_event = $1;
  