-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS revenue_analytics AS
SELECT b.registration_fee, 
e.name, 
e.event_type, 
es.event_date,
es.start_time,
es.end_time,
o.name
FROM bookings AS b
INNER JOIN event AS e on b.event_id = e.id
INNER JOIN event_schedule AS es ON e.id = es.event_id
INNER JOIN event_to_organizer_mapping AS etom ON e.id = etom.event_id
INNER JOIN organizer AS o ON etom.organizer_id = o.id;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS participant_analytics AS
SELECT s.name 
FROM bookings AS b
INNER JOIN student AS s ON b.student_id = s.id

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
