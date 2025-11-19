package pkg

import (
	"math/big"
	"strconv"
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
			otpStr, otpSlice, err := GenerateOTP()
			assert.NoError(t, err)

			otpInt, convErr := strconv.Atoi(otpStr)
			assert.NoError(t, convErr)
			assert.True(t, otpInt >= 0 && otpInt < 10000, "OTP %d out of expected range [0, 9999]", otpInt)

			// Check slice matches string characters
			assert.Equal(t, len(otpStr), len(otpSlice))
			for j, ch := range otpStr {
				assert.Equal(t, string(ch), otpSlice[j])
			}
		}
	})

	t.Run("Generates::different_otps", func(t *testing.T) {
		otpMax.Set(big.NewInt(10000))
		otps := make(map[string]struct{})
		uniqueCount := 0

		for range 1000 {
			otpStr, _, err := GenerateOTP()
			assert.NoError(t, err)
			if _, exists := otps[otpStr]; !exists {
				otps[otpStr] = struct{}{}
				uniqueCount++
			}
		}
		assert.Greater(t, uniqueCount, 900, "Expected a high number of unique OTPs for 1000 generations")
	})

	t.Run("Handles::large_max", func(t *testing.T) {
		otpMax.Set(big.NewInt(1000000000))

		for range 100 {
			otpStr, _, err := GenerateOTP()
			assert.NoError(t, err)

			otpInt, convErr := strconv.Atoi(otpStr)
			assert.NoError(t, convErr)
			assert.True(t, otpInt >= 0 && otpInt < 1000000000, "OTP %d out of expected range [0, 999999999]", otpInt)
		}
	})

	t.Run("Handles::zero_max", func(t *testing.T) {
		otpMax.Set(big.NewInt(0))

		otpStr, otpSlice, err := GenerateOTP()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max must be > 0 (got 0)")
		assert.Equal(t, "", otpStr)
		assert.Empty(t, otpSlice)
	})

	t.Run("Handles::negative_max", func(t *testing.T) {
		otpMax.Set(big.NewInt(-5))

		otpStr, otpSlice, err := GenerateOTP()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max must be > 0 (got -5)")
		assert.Equal(t, "", otpStr)
		assert.Empty(t, otpSlice)
	})

	t.Run("Handles::one_as_max", func(t *testing.T) {
		otpMax.Set(big.NewInt(1))

		for range 100 {
			otpStr, otpSlice, err := GenerateOTP()
			assert.NoError(t, err)
			assert.Equal(t, "0", otpStr)
			assert.Equal(t, []string{"0"}, otpSlice)
		}
	})
}
