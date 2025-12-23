-- name: GetEventsByOrganizerIdQuery :many
SELECT
    e.id,
    e.name,
    CASE
        WHEN e.is_group THEN 'GROUP'
        ELSE 'SOLO'
    END AS is_group
FROM event e
JOIN event_to_organizer_mapping AS etom ON e.id = etom.event_id
WHERE
    etom.organizer_id = $1;

-- name: CheckIfGroupEventQuery :one
SELECT is_group
FROM event
WHERE id = $1;

-- name: GetOrganizerSoloEventParticipantListQuery :many
SELECT
    s.name as student_name,
    s.college_name as college,
    s.college_city as city,
    s.email,
    s.is_amrita_student
FROM student s
LEFT JOIN 
    solo_event_participant sep ON sep.student_id = s.id
JOIN
    bookings b ON b.id = sep.booking_id
WHERE 
    sep.event_id = $1 
    AND b.txn_status = 'SUCCESS'
GROUP BY 
  s.name,
  s.college_name,
  s.college_city,
  s.email,
  s.is_amrita_student;


-- name: GetOrganizerGroupEventParticipantListQuery :many
SELECT
    t.team_name,
    s.name AS student_name,
    s.college_name as college,
    s.college_city as city,
    s.email,
    s.is_amrita_student 
FROM
    teams AS t
JOIN
    team_members AS tm ON t.id = tm.team_id
JOIN
    student AS s ON tm.student_id = s.id
JOIN
    bookings AS b on b.id = t.booking_id
WHERE
    t.event_id = $1
    AND b.txn_status = 'SUCCESS';
