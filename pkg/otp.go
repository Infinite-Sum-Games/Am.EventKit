package pkg

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (int, error) {
	max := big.NewInt(10000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, fmt.Errorf("failed to generate OTP: %w", err)
	}
	return int(n.Int64()), nil
}
