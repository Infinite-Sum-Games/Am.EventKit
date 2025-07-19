package pkg

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var otpMax = big.NewInt(10000)

func GenerateOTP() (int, error) {
	if otpMax.Cmp(big.NewInt(0)) <= 0 {
		return 0, fmt.Errorf("failed to generate OTP: max must be > 0 (got %s)", otpMax.String())
	}
	n, err := rand.Int(rand.Reader, otpMax)
	if err != nil {
		return 0, fmt.Errorf("failed to generate OTP: %w", err)
	}
	return int(n.Int64()), nil
}
