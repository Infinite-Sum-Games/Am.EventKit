package profile

import "github.com/google/uuid"

type Profile struct {
	ID       uuid.UUID `json:"user_id"`
	Email    string    `json:"user_email"`
	FullName string    `json:"user_full_name"`
}
