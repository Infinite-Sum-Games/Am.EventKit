package models

import "github.com/google/uuid"

type Event struct {
	ID   uuid.UUID `json:"event_id"`
	Name string    `json:"event_name"`
}
