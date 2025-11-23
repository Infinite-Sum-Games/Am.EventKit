package pkg

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

var otpMax = big.NewInt(10000)

// Totally over-engineered OTP generation function :(
func GenerateOTP() (string, []string, error) {
	if otpMax.Cmp(big.NewInt(0)) <= 0 {
		return "", []string{}, fmt.Errorf("failed to generate OTP: max must be > 0 (got %s)", otpMax.String())
	}

	n, err := rand.Int(rand.Reader, otpMax)
	if err != nil {
		return "", []string{}, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Convert OTP from Int64 to Int
	otp := int(n.Int64())

	// Get a string & []string from Int
	otpStr := strconv.Itoa(otp)
	otpSlice := make([]string, len(otpStr))
	for i, ch := range otpStr {
		otpSlice[i] = string(ch)
	}

	return otpStr, otpSlice, nil
}
