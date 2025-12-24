SELECT
  e.name AS event_name,
  COUNT(*) FILTER (WHERE b.txn_status = 'SUCCESS') AS tickets_sold,
  COUNT(*) FILTER (WHERE b.txn_status = 'PENDING') AS txn_in_progress,
  COUNT(*) FILTER (WHERE b.txn_status = 'FAILED')  AS failed_txn
FROM bookings b
JOIN event e ON e.id = b.event_id
GROUP BY b.event_id, e.name;
