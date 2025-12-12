CREATE TABLE revenue_analytics (
    id integer,
    total_revenue numeric,
    revenue_per_event jsonb,
    revenue_per_date jsonb,
    revenue_per_organizer jsonb
);

CREATE TABLE event_registration_analytics (
    id integer,
    total_event_registrations integer,
    participant_split jsonb,
    event_registration_stats jsonb
);

CREATE TABLE people_registration_analytics (
    id integer,
    website_registration_split jsonb,
    total_website_registrations integer
);

CREATE TABLE transaction_analytics (
    id integer,
    transaction_summary jsonb
);