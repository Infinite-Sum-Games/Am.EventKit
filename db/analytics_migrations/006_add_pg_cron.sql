-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pg_cron;

SELECT cron.schedule(
    'refresh_analytics_every_30min',
    '*/30 * * * *',
$$
REFRESH MATERIALIZED VIEW CONCURRENTLY revenue_analytics;
REFRESH MATERIALIZED VIEW CONCURRENTLY event_registration_analytics;
REFRESH MATERIALIZED VIEW CONCURRENTLY people_registration_analytics;
REFRESH MATERIALIZED VIEW CONCURRENTLY transaction_analytics;
$$
);
-- +goose StatementEnd