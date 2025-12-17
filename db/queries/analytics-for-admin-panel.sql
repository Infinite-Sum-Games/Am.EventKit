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