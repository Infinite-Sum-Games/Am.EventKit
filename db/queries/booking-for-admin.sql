-- name: FetchAdminTransactionsQuery :many
SELECT
  b.txn_id as transaction_id,
  e.id as event_id,
  e.name as event_name,
  s.name as student_name,
  s.email as email,
  s.phone_number as student_phone_number,
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

