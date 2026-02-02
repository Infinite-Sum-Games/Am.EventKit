package common

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"github.com/Infinite-Sum-Games/Am.EventKit/configs"
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

func GenerateSHA512Hash(txnId, email, amount, productInfo, name string) string {
	var key, salt string

	if configs.Env.App.Env == "PRODUCTION" {
		key = configs.Env.Payment.PayUProdKey
		salt = configs.Env.Payment.PayUProdSalt
	} else {
		key = configs.Env.Payment.PayUTestKey
		salt = configs.Env.Payment.PayUTestSalt
	}

	fields := []string{
		key,         // 1
		txnId,       // 2
		amount,      // 3
		productInfo, // 4
		name,        // 5
		email,       // 6
		"",          // 7 udf1
		"",          // 8 udf2
		"",          // 9 udf3
		"",          // 10 udf4
		"",          // 11 udf5
		"",          // 12 udf6
		"",          // 13 udf7
		"",          // 14 udf8
		"",          // 15 udf9
		"",          // 16 udf10
		salt,        // 18
	}
	data := strings.Join(fields, "|")

	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}

func GenerateVerifyPayUHash(txnID string) string {
	var key, salt string

	if configs.Env.App.Env == "PRODUCTION" {
		key = configs.Env.Payment.PayUProdKey
		salt = configs.Env.Payment.PayUProdSalt
	} else {
		key = configs.Env.Payment.PayUTestKey
		salt = configs.Env.Payment.PayUTestSalt
	}

	data := fmt.Sprintf(
		"%s|verify_payment|%s|%s",
		key,
		txnID,
		salt,
	)

	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}

func BuildVerifyPayUForm(txnID string) string {
	var key string

	if configs.Env.App.Env == "PRODUCTION" {
		key = configs.Env.Payment.PayUProdKey
	} else {
		key = configs.Env.Payment.PayUTestKey
	}

	values := url.Values{}
	values.Set("key", key)
	values.Set("command", "verify_payment")
	values.Set("hash", GenerateVerifyPayUHash(txnID))
	values.Set("var1", txnID)

	return values.Encode()
}
