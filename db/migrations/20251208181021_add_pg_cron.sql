-- +goose Up
-- +goose StatementBegin
SELECT cron.schedule(
    'refresh_revenue_every_30m',
    '*/30 * * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY revenue_analytics'
);
-- +goose StatementEnd

-- +goose StatementBegin
SELECT cron.schedule(
    'refresh_participant_every_30m',
    '*/30 * * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY participant_analytics'
);
-- +goose StatementEnd

-- +goose StatementBegin
SELECT cron.schedule(
    'refresh_registrations_every_30m',
    '*/30 * * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY registrations_analytics'
);
-- +goose StatementEnd

-- +goose StatementBegin
SELECT cron.schedule(
    'refresh_people_every_30m',
    '*/30 * * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY people_analytics'
);
-- +goose StatementEnd



