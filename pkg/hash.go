package pkg

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	hashedPassword := string(hashed)
	return hashedPassword, nil
}

func CompareHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateSHA512Hash(
	txnId,
	email,
	amount,
	productInfo,
	name string) string {

	fields := []string{
		cmd.Env.PayUKey,  // 1
		txnId,            // 2
		amount,           // 3
		productInfo,      // 4
		name,             // 5
		email,            // 6
		"",               // 7 udf1
		"",               // 8 udf2
		"",               // 9 udf3
		"",               // 10 udf4
		"",               // 11 udf5
		"",               // 12 udf6
		"",               // 13 udf7
		"",               // 14 udf8
		"",               // 15 udf9
		"",               // 16 udf10
		cmd.Env.PayUSalt, // 25 salt (mandatory)
	}
	data := strings.Join(fields, "|")

	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}
