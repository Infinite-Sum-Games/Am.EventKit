-- name: GetTeamDetailsByEventID :many
WITH team_leader AS (
    SELECT
        tm.team_id,
        s.name AS leader_name,
        s.email AS leader_email,
        s.phone_number AS leader_phone_number,
        s.college_name AS leader_college_name
    FROM team_members tm
    JOIN student s ON tm.student_id = s.id
    WHERE tm.student_role = 'leader'
),
team_members_list AS (
    SELECT
        tm.team_id,
        json_agg(json_build_object(
            'name', s.name,
            'email', s.email,
            'phone_number', s.phone_number,
            'college_name', s.college_name
        )) AS members
    FROM team_members tm
    JOIN student s ON tm.student_id = s.id
    WHERE tm.student_role != 'leader'
    GROUP BY tm.team_id
)
SELECT
    t.team_name,
    tl.leader_name,
    tl.leader_email,
    tl.leader_phone_number,
    tl.leader_college_name,
    t.metadata->'problem_stmt' AS problem_statements,
    COALESCE(tml.members, '[]'::json) AS team_members
FROM teams AS t
JOIN bookings AS b ON t.booking_id = b.id
JOIN team_leader AS tl ON t.id = tl.team_id
LEFT JOIN team_members_list AS tml ON t.id = tml.team_id
WHERE t.event_id = $1;
