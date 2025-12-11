-- name: AddEventDependency :exec
INSERT INTO event_dependency_mapping (
  start_event, 
  end_event)
VALUES ($1, $2);

-- name: RemoveEventDependency :execrows
DELETE FROM event_dependency_mapping
WHERE start_event = $1 AND end_event = $2;