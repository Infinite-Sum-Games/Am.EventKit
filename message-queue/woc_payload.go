package messagequeue

import (
	"encoding/json"
	"strings"
)

// Payload builder for Winter of Code 25' event.
type WoCPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func CreateWoCPayload(email, fullName string) ([]byte, error) {
	nameParts := strings.Split(fullName, " ")
	firstName := ""
	lastName := ""
	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}

	payload := WoCPayload{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return jsonPayload, nil
}
