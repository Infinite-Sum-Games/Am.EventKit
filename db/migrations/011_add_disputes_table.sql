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
  student_email TEXT,
  description TEXT,
  event_id UUID NOT NULL,
  dispute_status dispute_status_enum DEFAULT 'OPEN',
  team_member_datails JSONB,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "dispute_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose down
-- +goose StatementBegin
DROP TABLE IF EXISTS dispute;
-- +goose StatementEnd
