-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_cron;
-- +goose StatementEnd

-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS revenue_analytics;

CREATE MATERIALIZED VIEW revenue_analytics AS
SELECT
    1 AS id,

    -- TOTAL REVENUE
    (
        SELECT jsonb_build_object(
            'revenue', COALESCE(SUM(registration_fee), 0),
            'revenue_without_gst', COALESCE(SUM(registration_fee_without_gst), 0)
        )
        FROM bookings
        WHERE txn_status = 'SUCCESS'
    ) AS total_revenue,

    -- REVENUE GROUPED BY EVENTS (FULLY ALIGNED)
    (
        WITH successful_bookings AS (
            SELECT
                event_id,
                SUM(registration_fee) AS revenue,
                SUM(COALESCE(registration_fee_without_gst, 0)) AS revenue_without_gst,
                COUNT(id) AS seats_filled
            FROM bookings
            WHERE txn_status = 'SUCCESS'
            GROUP BY event_id
        ),
        team_participant_count AS (
            SELECT
                t.event_id,
                COUNT(tm.id) AS participant_count
            FROM team_members tm
            JOIN teams t ON tm.team_id = t.id
            JOIN bookings b ON t.booking_id = b.id
            WHERE b.txn_status = 'SUCCESS'
            GROUP BY t.event_id
        )
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT
                e.id AS event_id,
                e.name AS event_name,

                COALESCE(sb.seats_filled, 0) AS seats_filled,
                COALESCE(sb.revenue, 0) AS revenue,
                COALESCE(sb.revenue_without_gst, 0) AS revenue_without_gst,

                e.total_seats,
                e.is_group,
                e.event_type,

                CASE
                    WHEN e.is_group = TRUE THEN
                        COALESCE(tpc.participant_count, 0)
                    ELSE
                        COALESCE(sb.seats_filled, 0)
                END 
                AS actual_participant_count

            FROM event e
            LEFT JOIN successful_bookings sb ON e.id = sb.event_id
            LEFT JOIN team_participant_count tpc ON e.id = tpc.event_id
            ORDER BY e.name
        ) x
    ) AS revenue_per_event,

    -- REVENUE GROUPED BY DATE
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                DATE(b.created_at) AS revenue_date,
                SUM(b.registration_fee) AS revenue,
                SUM(b.registration_fee_without_gst) AS revenue_without_gst
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
                SUM(b.registration_fee) AS revenue,
                SUM(b.registration_fee_without_gst) AS revenue_without_gst
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
DROP MATERIALIZED VIEW IF EXISTS people_registration_analytics;

CREATE MATERIALIZED VIEW people_registration_analytics AS
SELECT
    1 AS id,

    -- Internal vs External website registrations
    (
        SELECT jsonb_build_object(
            'internal', COUNT(*) FILTER (WHERE is_amrita_student = TRUE),
            'external', COUNT(*) FILTER (WHERE is_amrita_student = FALSE)
        )
        FROM student
    ) AS website_registration_split,

    -- Total website registrations
    (
        SELECT COUNT(*) FROM student
    ) AS total_website_registrations,

    -- Website registrations grouped by year (Amrita students only)
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT
                year_label,
                COUNT(*) AS student_count
            FROM (
                SELECT
                    CASE
                        WHEN batch_year = 25 THEN '1st_year'
                        WHEN batch_year = 24 THEN '2nd_year'
                        WHEN batch_year = 23 THEN '3rd_year'
                        WHEN batch_year = 22 THEN '4th_year'
                        WHEN batch_year = 21 THEN '5th_year'
                        ELSE 'unknown'
                    END AS year_label
                FROM (
                    SELECT
                        substring(
                            split_part(amrita_roll_number, '.', 4),
                            '\d{2}'
                        )::int AS batch_year
                    FROM student
                    WHERE is_amrita_student = TRUE
                      AND amrita_roll_number IS NOT NULL
                ) y
            ) labeled
            GROUP BY year_label
            ORDER BY year_label
        ) x
    ) AS website_registration_by_year,

    -- Event registrations grouped by year (unique students)
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT
                year_label,
                COUNT(*) AS student_count
            FROM (
                SELECT DISTINCT
                    s.id,
                    CASE
                        WHEN batch_year = 25 THEN '1st_year'
                        WHEN batch_year = 24 THEN '2nd_year'
                        WHEN batch_year = 23 THEN '3rd_year'
                        WHEN batch_year = 22 THEN '4th_year'
                        WHEN batch_year = 21 THEN '5th_year'
                        ELSE 'unknown'
                    END AS year_label
                FROM bookings b
                JOIN student s ON s.id = b.student_id
                CROSS JOIN LATERAL (
                    SELECT
                        substring(
                            split_part(s.amrita_roll_number, '.', 4),
                            '\d{2}'
                        )::int AS batch_year
                ) r
                WHERE b.txn_status = 'SUCCESS'
                  AND s.is_amrita_student = TRUE
                  AND s.amrita_roll_number IS NOT NULL
            ) dedup
            GROUP BY year_label
            ORDER BY year_label
        ) x
    ) AS event_registration_by_year;

-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX people_registration_analytics_id_index 
ON people_registration_analytics(id);
-- +goose StatementEnd