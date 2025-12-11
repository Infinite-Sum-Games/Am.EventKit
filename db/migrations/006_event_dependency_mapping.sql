-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_dependency_mapping (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  start_event UUID NOT NULL,
  end_event UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT fk_start_event
    FOREIGN KEY (start_event) REFERENCES event(id) ON DELETE CASCADE,

  CONSTRAINT fk_end_event
    FOREIGN KEY (end_event) REFERENCES event(id) ON DELETE CASCADE,

  CONSTRAINT unique_dependency UNIQUE (start_event, end_event),
  CONSTRAINT no_self_dependency CHECK (start_event <> end_event)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_dependency_mapping;
-- +goose StatementEnd
