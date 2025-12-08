-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS revenue_analytics AS
SELECT b.registration_fee AS booking_fee, 
e.id AS event_id,
e.name AS event_name, 
e.event_type AS event_type,
es.event_date AS event_date,
es.start_time AS start_time,
es.end_time AS end_time,
o.id AS organizer_id,
o.name AS organizer_name,
SUM(b.registration_fee) OVER () AS total_revenue,
SUM(b.registration_fee) OVER (PARTITION BY e.id) AS revenue_per_event,
SUM(b.registration_fee) OVER (PARTITION BY es.event_date) AS revenue_per_date,
SUM(b.registration_fee) OVER (PARTITION BY e.event_type) AS revenue_per_event_type,
SUM(b.registration_fee) OVER (PARTITION BY o.id) AS revenue_per_organizer
FROM bookings AS b
INNER JOIN event AS e on b.event_id = e.id
INNER JOIN event_schedule AS es ON e.id = es.event_id
INNER JOIN event_to_organizer_mapping AS etom ON e.id = etom.event_id
INNER JOIN organizer AS o ON etom.organizer_id = o.id;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS participant_analytics AS
SELECT s.id AS student_id,
s.name AS student_name,
s.email AS student_email,
e.id AS event_id,
e.name AS event_name,
e.event_type AS event_type,
COUNT(b.id) OVER () AS global_total_bookings,
COUNT(b.id) OVER (PARTITION BY e.id) AS participants_per_event,
COUNT(b.id) OVER (PARTITION BY o.id) AS bookings_per_organizer
FROM bookings AS b
INNER JOIN student AS s ON b.student_id = s.id
INNER JOIN event AS e ON b.event_id = e.id
INNER JOIN event_to_organizer_mapping AS etom ON e.id = etom.event_id
INNER JOIN organizer AS o ON etom.organizer_id = o.id;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS registrations_analytics AS
SELECT UNIQUE b.student_id AS student_id,
s.name AS student_name,
s.email AS student_email,
s.is_amrita_student AS is_amrita_student,
e.id AS event_id,
e.name AS event_name,
e.event_type AS event_type,
o.id AS organizer_id,
o.name AS organizer_name,
COUNT(UNIQUE b.student_id) OVER () AS total_registrations
COUNT(UNIQUE b.student_id) OVER (PARTITION BY s.is_amrita_student) 
AS registrations_by_student_type,
FROM bookings AS b
INNER JOIN student AS s ON b.student_id = s.id
INNER JOIN event AS e ON b.event_id = e.id
INNER JOIN event_to_organizer_mapping AS etom ON e.id = etom.event_id
INNER JOIN organizer AS o ON etom.organizer_id = o.id;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS people_analytics AS
SELECT p.id AS person_id,
p.name AS person_name,
COUNT(p.id) OVER () AS total_people,
FROM person AS p;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
