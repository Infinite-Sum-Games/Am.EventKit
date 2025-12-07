-- name: GetRevenueAnalytics :many
SELECT *
FROM revenue_analytics;


-- name: RefreshRevenueAnalytics :exec
SELECT *
FROM cron.job
WHERE jobname = 'refresh_revenue_every_30m'
UNION ALL
SELECT cron.schedule(
    'refresh_revenue_every_30m',
    '*/30 * * * *',
    $$REFRESH MATERIALIZED VIEW CONCURRENTLY revenue_analytics;$$
)
WHERE NOT EXISTS (
    SELECT 1 FROM cron.job WHERE jobname = 'refresh_revenue_every_30m'
);

