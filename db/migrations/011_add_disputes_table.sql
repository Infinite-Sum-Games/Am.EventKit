-- +goose up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dispute (
  id UUID DEFAULT gen_random_uuid(),
  email TEXT,
  description TEXT,
  event_id UUID NOT NULL,
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
