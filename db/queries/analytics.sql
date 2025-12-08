-- name: GetRevenueQuery :many
SELECT booking_fee,
event_id,
event_name,
event_type,
event_date,
start_time,
end_time,
organizer_id,
organizer_name
FROM revenue_analytics;

-- name: GetRevenueSummaryQuery :one
SELECT total_revenue,
revenue_per_event,
revenue_per_date,
revenue_per_event_type,
revenue_per_organizer
FROM revenue_analytics
LIMIT 1;


-- name: GetParticipantQuery :many
SELECT student_id,
student_name,
student_email,
event_id,
event_name,
event_type
FROM participant_analytics;

-- name: GetParticipantSummaryQuery :one
SELECT global_total_bookings,
participants_per_event,
bookings_per_organizer
FROM participant_analytics
LIMIT 1;

-- name: GetRegistrationsQuery :many
SELECT student_id,
student_name,
student_email,
is_amrita_student,
event_id,
event_name,
event_type,
organizer_id,
organizer_name
FROM registrations_analytics;

-- name: GetRegistrationSummaryQuery :one
SELECT total_registrations,
registrations_by_student_type
FROM registrations_analytics;

-- name: ListPeopleQuery :many
SELECT person_id,
person_name
FROM people_analytics;

-- name: GetPeopleCountQuery :one
SELECT total_people
FROM people_analytics;
