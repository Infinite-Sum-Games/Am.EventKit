-- +goose Up

-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS event_registration_analytics;
DROP MATERIALIZED VIEW IF EXISTS participants_analytics;

CREATE MATERIALIZED VIEW participants_analytics AS
WITH participants AS (
    SELECT student_id FROM solo_event_participant
    UNION
    SELECT student_id FROM team_members
    INNER JOIN teams t ON t.id = team_members.team_id
    INNER JOIN bookings b ON b.id = t.booking_id
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

    -- Amrita vs Non-Amrita Participants
    (
        SELECT jsonb_build_object(
            'amrita_participants', 
                (SELECT COUNT(*) FROM participants p
                 JOIN student s ON s.id = p.student_id
                 WHERE s.is_amrita_student = TRUE),
            'non_amrita_participants', 
                (SELECT COUNT(*) FROM participants p
                 JOIN student s ON s.id = p.student_id
                 WHERE s.is_amrita_student = FALSE)
        )
    ) AS amrita_non_amrita_split,

    -- (REGISTRATIONS VS TOTAL SEATS) PER EVENT
    (
        SELECT jsonb_agg(row_to_json(x))
        FROM (
            SELECT 
                e.id AS event_id,
                e.name AS event_name,
                e.total_seats,
                COUNT(DISTINCT CASE WHEN b_sp.txn_status = 'SUCCESS' THEN sp.student_id END) 
                    + COUNT(DISTINCT CASE WHEN b_t.txn_status = 'SUCCESS' THEN tm.student_id END) AS registered_count
            FROM event e
            LEFT JOIN solo_event_participant sp ON sp.event_id = e.id
   66       LEFT JOIN bookings b_sp ON b_sp.id = sp.booking_id
   67       LEFT JOIN teams t ON t.event_id = e.id
   68       LEFT JOIN bookings b_t ON b_t.id = t.booking_id
   69       LEFT JOIN team_members tm ON tm.team_id = t.id
   70       GROUP BY e.id, e.name, e.total_seats            
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
