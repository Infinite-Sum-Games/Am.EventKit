package pkg

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unable to generate hash %w", err)
	}
	hashedPassword := string(hashed)
	return hashedPassword, nil
}
