-- name: FetchAdminTransactionsQuery :many
SELECT
  b.txn_id AS transaction_id,
  e.id AS event_id,
  e.name AS event_name,
  s.name AS student_name,
  e.is_group AS is_group,
  s.email AS email,
  s.phone_number AS student_phone_number,
  s.college_name,
  s.college_city,
  s.is_amrita_student,
  b.txn_status,
  b.registration_fee

FROM bookings b

LEFT JOIN event e ON b.event_id = e.id
LEFT JOIN student s ON b.student_id = s.id

WHERE
  b.txn_status = $1
ORDER BY b.updated_at DESC;

