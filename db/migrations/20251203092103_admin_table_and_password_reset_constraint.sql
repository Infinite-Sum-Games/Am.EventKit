-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS admin (
  id UUID DEFAULT gen_random_uuid(),
  name TEXT,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  refresh_token TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT admin_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE password_reset
ADD CONSTRAINT unique_password_reset_email UNIQUE (email);
-- +goose StatementBegin

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS admin;
-- +goose StatementEnd
