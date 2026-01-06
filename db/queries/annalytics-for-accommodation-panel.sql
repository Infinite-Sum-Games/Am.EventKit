-- name: GetInsideCampusAnalyticsQuery :many
SELECT
    logged_at::date AS date,
    jsonb_build_object(
        'IN',  COUNT(*) FILTER (WHERE direction = 'IN'),
        'OUT', COUNT(*) FILTER (WHERE direction = 'OUT'),
        'CURRENTLY_INSIDE', COUNT(*) FILTER (WHERE direction = 'IN') - COUNT(*) FILTER (WHERE direction = 'OUT')
    ) AS counts
FROM gate_management
GROUP BY date
ORDER BY date;

-- name: GetLiveBedsAnalyticsQuery :many
SELECT jsonb_build_object(
    'hostels', jsonb_agg(
        jsonb_build_object(
            'id', id,
            'hostel_name', hostel_name,
            'room_count', room_count,
            'room_filled', room_filled
        )
    ),
    'total_beds_filled', SUM(room_filled)
)
FROM hostel_metadata;

