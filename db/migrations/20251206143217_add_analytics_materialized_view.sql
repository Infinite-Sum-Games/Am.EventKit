-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS revenue_analytics AS
SELECT b.id AS booking_id,
b.registration_fee AS booking_fee, 
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

-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX revenue_analytics_event_id_id
ON revenue_analytics (booking_id);
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

-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX participant_analytics_student_id_event_id_idx
ON participant_analytics (student_id, event_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS registrations_analytics AS
WITH unique_student AS (
    SELECT DISTINCT b.student_id
    FROM bookings b
),

total_unique AS (
    SELECT COUNT(*) AS total_registrations
    FROM unique_student
),

unique_by_type AS (
    SELECT
        s.is_amrita_student,
        COUNT(*) AS registrations_by_student_type
    FROM unique_student us
    JOIN student s ON s.id = us.student_id
    GROUP BY s.is_amrita_student
)

SELECT
    s.id AS student_id,
    s.name AS student_name,
    s.email AS student_email,
    s.is_amrita_student AS is_amrita_student,

    e.id AS event_id,
    e.name AS event_name,
    e.event_type AS event_type,
    o.id AS organizer_id,
    o.name AS organizer_name,

    tu.total_registrations,
    ubt.registrations_by_student_type

FROM bookings b
JOIN student s ON s.id = b.student_id
JOIN event e ON e.id = b.event_id
JOIN event_to_organizer_mapping etom ON e.id = etom.event_id
JOIN organizer o ON etom.organizer_id = o.id
CROSS JOIN total_unique tu
LEFT JOIN unique_by_type ubt
    ON ubt.is_amrita_student = s.is_amrita_student;

-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX registrations_analytics_student_id_event_id_idx
ON registrations_analytics (student_id, event_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS people_analytics AS
SELECT p.id AS person_id,
p.name AS person_name,
COUNT(p.id) OVER () AS total_people
FROM people AS p;

-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX people_analytics_person_id_idx 
ON people_analytics (person_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS revenue_analytics;
DROP MATERIALIZED VIEW IF EXISTS participant_analytics;
DROP MATERIALIZED VIEW IF EXISTS registrations_analytics;
DROP MATERIALIZED VIEW IF EXISTS people_analytics;
SELECT cron.unschedule('refresh_revenue_every_30m');
SELECT cron.unschedule('refresh_participant_every_30m');
SELECT cron.unschedule('refresh_registrations_every_30m');
SELECT cron.unschedule('refresh_people_every_30m');
-- +goose StatementEnd
