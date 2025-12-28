-- name: GetRevenueAnalytics :one
SELECT 
    id,
    total_revenue,
    revenue_per_event,
    revenue_per_date,
    revenue_per_organizer
FROM revenue_analytics
WHERE id = 1;

-- name: GetEventRegistrationAnalytics :one
SELECT
    id,
    total_event_registrations,
    participant_split,
    event_registration_stats
FROM event_registration_analytics
WHERE id = 1;

-- name: GetPeopleRegistrationAnalytics :one
SELECT
    id,
    website_registration_split,
    total_website_registrations
FROM people_registration_analytics
WHERE id = 1;

-- name: GetTransactionAnalytics :one
SELECT
    id,
    transaction_summary
FROM transaction_analytics
WHERE id = 1;

-- name: GetQuickDashboardQuery :many
WITH successful_bookings AS (
-- CTE to calculate revenue and seats filled from successful txns
  SELECT
    event_id,
    SUM(registration_fee) AS revenue,
    SUM(COALESCE(registration_fee_without_gst, 0)) AS revenue_without_gst,
    COUNT(id) AS seats_filled
  FROM
    bookings
  WHERE
    txn_status = 'SUCCESS'
  GROUP BY
    event_id
),
-- CTE to count the total number of members in teams for each event
team_participant_count AS (
  SELECT
    t.event_id,
    COUNT(tm.id) AS participant_count
  FROM
    team_members tm
  JOIN teams t ON tm.team_id = t.id
  JOIN bookings b ON t.booking_id = b.id
  WHERE b.txn_status = 'SUCCESS'
  GROUP BY t.event_id
)
-- Data combination step
SELECT
  e.id AS event_id,
  e.name AS event_name,
  COALESCE(sb.revenue, 0) AS revenue,
  COALESCE(sb.revenue_without_gst, 0) AS revenue_without_gst,
  COALESCE(sb.seats_filled, 0) AS seats_filled,
  e.total_seats,
  e.is_group,
  e.event_type,
CASE
  WHEN e.is_group = TRUE THEN
    COALESCE(tpc.participant_count, 0)
  ELSE
    COALESCE(sb.seats_filled, 0)
  END AS actual_participant_count
FROM
  event e
LEFT JOIN successful_bookings sb ON e.id = sb.event_id
LEFT JOIN team_participant_count tpc ON e.id = tpc.event_id
ORDER BY e.name;
