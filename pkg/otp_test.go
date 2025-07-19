package pkg

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOTP(t *testing.T) {
	originalOTPMax := new(big.Int).Set(otpMax)

	t.Cleanup(func() {
		otpMax.Set(originalOTPMax)
	})

	t.Run("Generate::valid_otp_within_range", func(t *testing.T) {
		otpMax.Set(big.NewInt(10000))

		for range 1000 {
			otp, err := GenerateOTP()
			assert.NoError(t, err)
			assert.True(t, otp >= 0 && otp < 10000, "OTP %d out of expected range [0, 9999]", otp)
		}
	})

	t.Run("Generates::different_otps", func(t *testing.T) {
		otpMax.Set(big.NewInt(10000))
		otps := make(map[int]struct{})
		uniqueCount := 0

		for range 1000 {
			otp, err := GenerateOTP()
			assert.NoError(t, err)
			if _, exists := otps[otp]; !exists {
				otps[otp] = struct{}{}
				uniqueCount++
			}
		}
		assert.Greater(t, uniqueCount, 900, "Expected a high number of unique OTPs for 1000 generations")
	})

	t.Run("Handles::large_max", func(t *testing.T) {
		otpMax.Set(big.NewInt(1000000000))

		for range 100 {
			otp, err := GenerateOTP()
			assert.NoError(t, err)
			assert.True(t, otp >= 0 && otp < 1000000000, "OTP %d out of large expected range [0, 999999999]", otp)
		}
	})

	t.Run("Handles::zero_max", func(t *testing.T) {
		otpMax.Set(big.NewInt(0))

		otp, err := GenerateOTP()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max must be > 0 (got 0)")
		assert.Equal(t, 0, otp)
	})

	t.Run("Handles::negative_max", func(t *testing.T) {
		otpMax.Set(big.NewInt(-5))

		otp, err := GenerateOTP()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max must be > 0 (got -5)")
		assert.Equal(t, 0, otp)
	})

	t.Run("Handles::one_as_max", func(t *testing.T) {
		otpMax.Set(big.NewInt(1))

		for range 100 {
			otp, err := GenerateOTP()
			assert.NoError(t, err)
			assert.Equal(t, 0, otp)
		}
	})
}
