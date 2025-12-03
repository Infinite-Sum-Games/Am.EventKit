package pkg

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashedPassword := string(hashed)
	return hashedPassword, nil
}

func CompareHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateSHA512Hash(txnId, email, amount, productInfo, eventId, salt string) string {
	data := strings.Join([]string{txnId, email, amount, productInfo, eventId, salt}, "|")

	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}
