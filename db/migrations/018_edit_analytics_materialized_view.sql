-- +goose Up

-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS event_registration_analytics;

CREATE MATERIALIZED VIEW event_registration_analytics AS
WITH participants AS (
    SELECT student_id FROM solo_event_participant
    UNION
    SELECT team_members.student_id FROM team_members
    INNER JOIN teams t ON team_members.team_id = t.id
    INNER JOIN bookings b ON t.booking_id = b.id
    WHERE b.txn_status = 'SUCCESS'
)
SELECT
    1 AS id,

    -- TOTAL EVENT REGISTRATIONS
    (SELECT COUNT(*) FROM participants) AS total_event_registrations,

    -- PARTICIPANTS VS NON-PARTICIPANTS
    (
        SELECT jsonb_build_object(
            'participants', (SELECT COUNT(*) FROM participants),
            'non_participants', 
                (SELECT COUNT(*) FROM student s 
                 WHERE s.id NOT IN (SELECT student_id FROM participants))
        )
    ) AS participant_split,

    -- (REGISTRATIONS VS TOTAL SEATS) PER EVENT
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                e.id AS event_id,
                e.name AS event_name,
                e.total_seats,
                COUNT(DISTINCT sp.student_id) 
                    + (SELECT (COUNT(DISTINCT tm.student_id))
                    FROM team_members tm
                    INNER JOIN teams t ON tm.team_id = t.id
                    INNER JOIN bookings b ON t.booking_id = b.id
                    WHERE t.event_id = e.id AND b.txn_status = 'SUCCESS')
                    AS registered_count
            FROM event e
            LEFT JOIN solo_event_participant sp 
                ON sp.event_id = e.id
            GROUP BY e.id, e.name
            ORDER BY e.name
        ) x
    ) AS event_registration_stats;


-- Create unique index for CONCURRENT refresh support
CREATE UNIQUE INDEX event_registration_analytics_id_index 
ON event_registration_analytics(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS event_registration_analytics;
-- +goose StatementEnd
