-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW revenue_analytics AS
SELECT
    1 AS id,

    -- TOTAL REVENUE
    (
        SELECT COALESCE(SUM(registration_fee), 0)
        FROM bookings
        WHERE txn_status = 'SUCCESS'
    ) AS total_revenue,

    -- REVENUE GROUPED BY EVENTS
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                e.id AS event_id,
                e.name AS event_name,
                SUM(b.registration_fee) AS revenue
            FROM bookings b
            JOIN event e ON e.id = b.event_id
            WHERE b.txn_status = 'SUCCESS'
            GROUP BY e.id
            ORDER BY e.name
        ) x
    ) AS revenue_per_event,

    -- REVENUE GROUPED BY DATE
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                DATE(b.created_at) AS revenue_date,
                SUM(b.registration_fee) AS revenue
            FROM bookings b
            WHERE b.txn_status = 'SUCCESS'
            GROUP BY revenue_date
            ORDER BY revenue_date
        ) x
    ) AS revenue_per_date,

    -- REVENUE GROUPED BY ORGANIZERS
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                o.id AS organizer_id,
                o.name AS organizer_name,
                SUM(b.registration_fee) AS revenue
            FROM bookings b
            JOIN event_to_organizer_mapping m ON m.event_id = b.event_id
            JOIN organizer o ON o.id = m.organizer_id
            WHERE b.txn_status = 'SUCCESS'
            GROUP BY o.id
            ORDER BY o.name
        ) x
    ) AS revenue_per_organizer;

-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX revenue_analytics_id_index 
ON revenue_analytics(id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW event_registration_analytics AS
SELECT
    1 AS id,

    -- TOTAL EVENT REGISTRATIONS
    (
        WITH all_participants AS (
            SELECT student_id FROM solo_event_participant
            UNION
            SELECT student_id FROM team_members
        )
        SELECT COUNT(*) FROM all_participants
    ) AS total_event_registrations,

    -- PARTICIPANTS VS NON-PARTICIPANTS
    (
        WITH participants AS (
            SELECT student_id FROM solo_event_participant
            UNION
            SELECT student_id FROM team_members
        )
        SELECT jsonb_build_object(
            'participants', (SELECT COUNT(*) FROM participants),
            'non_participants', 
                (SELECT COUNT(*) FROM student s 
                 WHERE s.id NOT IN (SELECT student_id FROM participants))
        )
    ) AS participant_split,

    -- (REGISTRATIONS VS TOTAL SEATS) PER EVENT
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                e.id AS event_id,
                e.name AS event_name,
                e.total_seats,
                COUNT(DISTINCT sp.student_id) 
                    + COUNT(DISTINCT tm.student_id) AS registered_count
            FROM event e
            LEFT JOIN solo_event_participant sp ON sp.event_id = e.id
            LEFT JOIN teams t ON t.event_id = e.id
            LEFT JOIN team_members tm ON tm.team_id = t.id
            GROUP BY e.id
            ORDER BY e.name
        ) x
    ) AS event_registration_stats;


-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX event_registration_analytics_id_index 
ON event_registration_analytics(id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW people_registration_analytics AS
SELECT
    1 AS id,

    -- INTERNAL VS EXTERNAL WEBSITE REGISTRATIONS
    (
        SELECT jsonb_build_object(
            'internal', COUNT(*) FILTER (WHERE is_amrita_student = TRUE),
            'external', COUNT(*) FILTER (WHERE is_amrita_student = FALSE)
        )
        FROM student
    ) AS website_registration_split,

    -- TOTAL WEBSITE REGISTRATIONS
    (
        SELECT COUNT(*) FROM student
    ) AS total_website_registrations;

-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX people_registration_analytics_id_index 
ON people_registration_analytics(id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE MATERIALIZED VIEW transaction_analytics AS
SELECT
    1 AS id,

    -- FAILED VS SUCCESSFUL VS PENDING TRANSACTIONS
    (
        SELECT jsonb_build_object(
            'success', COUNT(*) FILTER (WHERE txn_status='SUCCESS'),
            'failed', COUNT(*) FILTER (WHERE txn_status='FAILED'),
            'pending', COUNT(*) FILTER (WHERE txn_status='PENDING'),
            'total', COUNT(*)
        )
        FROM bookings
    ) AS transaction_summary;

-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX transaction_analytics_id_index 
ON transaction_analytics(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS revenue_analytics;
DROP MATERIALIZED VIEW IF EXISTS event_registration_analytics;
DROP MATERIALIZED VIEW IF EXISTS people_registration_analytics;
DROP MATERIALIZED VIEW IF EXISTS transaction_analytics;
SELECT cron.unschedule('refresh_analytics_every_30min');
-- +goose StatementEnd