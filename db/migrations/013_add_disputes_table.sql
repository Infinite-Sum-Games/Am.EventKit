-- +goose up
-- +goose StatementBegin
CREATE TYPE dispute_status_enum AS ENUM (
  'OPEN',
  'CLOSED_AS_TRUE',
  'CLOSED_AS_FALSE'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dispute (
  id UUID DEFAULT gen_random_uuid(),
  txn_id TEXT NOT NULL,
  student_email TEXT,
  description TEXT,
  event_id UUID NOT NULL,
  dispute_status dispute_status_enum DEFAULT 'OPEN',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "dispute_pkey" PRIMARY KEY (id),

  CONSTRAINT "dispute_to_student_mapping_fkey"
    FOREIGN KEY (student_email)
    REFERENCES student (email)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,

  CONSTRAINT "dispute_to_event_mapping_fkey"
    FOREIGN KEY (event_id)
    REFERENCES event (id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,

  CONSTRAINT "dispute_to_txn_mapping_fkey"
    FOREIGN KEY (txn_id)
    REFERENCES bookings (txn_id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE
  );
-- +goose StatementEnd

-- +goose down
-- +goose StatementBegin
DROP TABLE IF EXISTS dispute;
DROP TYPE IF EXISTS dispute_status_enum;
-- +goose StatementEnd
