package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHash(t *testing.T) {
	t.Run("ValidPassword::generates_valid_hash", func(t *testing.T) {
		password := "myStrongPassword123!"
		hashedPassword, err := Hash(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)

		compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		assert.NoError(t, compareErr)
	})

	t.Run("EmptyPassword::generates_valid_hash", func(t *testing.T) {
		password := ""
		hashedPassword, err := Hash(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)

		compareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		assert.NoError(t, compareErr)
	})

	t.Run("DifferentPasswords::produce_different_hashes", func(t *testing.T) {
		password1A := "uniquePassword1"
		password1B := "uniquePassword1"
		password2 := "anotherUniquePass2"

		hash1A, err1A := Hash(password1A)
		assert.NoError(t, err1A)
		hash1B, err1B := Hash(password1B)
		assert.NoError(t, err1B)
		hash2, err2 := Hash(password2)
		assert.NoError(t, err2)

		assert.NotEqual(t, hash1A, hash1B)
		assert.NotEqual(t, hash1A, hash2)
		assert.NotEqual(t, hash1B, hash2)

		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(hash1A), []byte(password1A)))
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(hash1B), []byte(password1B)))
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(hash2), []byte(password2)))
	})

	t.Run("InvalidHash::comparison_fail_expected", func(t *testing.T) {
		password := "testpassword"
		invalidHash := "$2a$10$dS1d2S3d4F5g6H7j8K9l.eX/R/Y/Z0a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r"

		compareErr := bcrypt.CompareHashAndPassword([]byte(invalidHash), []byte(password))
		assert.Error(t, compareErr)
		assert.Contains(t, compareErr.Error(), "hashedPassword is not the hash of the given password")
	})
}
