-- +goose Up
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS revenue_analytics;
DROP MATERIALIZED VIEW IF EXISTS transaction_analytics;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE bookings
ALTER COLUMN registration_fee
TYPE INTEGER
USING TRUNC(registration_fee);
-- +goose StatementEnd

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