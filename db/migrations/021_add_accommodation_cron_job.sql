-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION cleanup_expired_hostel_allotments_func() 
RETURNS void AS $$
BEGIN
    WITH expired_rows AS (
      SELECT id, hostel_id
      FROM accomodation_details
      WHERE payment_status = 'PENDING'
        AND payment_expires < NOW()
        AND hostel_id IS NOT NULL
      FOR UPDATE
    ),
    update_hostels AS (
      UPDATE hostel_metadata h
      SET room_count = h.room_count + expired_counts.total
      FROM (
        SELECT hostel_id, COUNT(*) as total
        FROM expired_rows
        GROUP BY hostel_id
      ) AS expired_counts
      WHERE h.id = expired_counts.hostel_id
      RETURNING 1
    )
    UPDATE accomodation_details
    SET
      hostel_id = NULL,
      payment_status = 'FAILED',
      updated_at = NOW()
    WHERE id IN (SELECT id FROM expired_rows);
END;
$$ LANGUAGE plpgsql;

SELECT cron.schedule(
  'cleanup_expired_hostel_allotments',
  '*/10 * * * *',
  'SELECT cleanup_expired_hostel_allotments_func()'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT cron.unschedule('cleanup_expired_hostel_allotments');
DROP FUNCTION IF EXISTS cleanup_expired_hostel_allotments_func();
-- +goose StatementEnd
