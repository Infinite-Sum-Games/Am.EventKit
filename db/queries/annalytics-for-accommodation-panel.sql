-- name: GetInsideCampusAnalyticsQuery :many
SELECT
    logged_at::date AS date,
    jsonb_build_object(
        'IN',  COUNT(*) FILTER (WHERE direction = 'IN'),
        'OUT', COUNT(*) FILTER (WHERE direction = 'OUT')
    ) AS counts
FROM gate_management
GROUP BY date
ORDER BY date;

-- name: GetLiveBedsAnalyticsQuery :many
SELECT
    room_filled
FROM hostel_metadata;
