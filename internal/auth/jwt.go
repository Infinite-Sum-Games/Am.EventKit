package auth

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
)

// SetRSAKeys sets the keys into memory
func SetRSAKeys(priv *rsa.PrivateKey, pub *rsa.PublicKey) {
	rsaPrivateKey = priv
	rsaPublicKey = pub
}

// GenerateJWT creates a new JWT signed with the RSA private key
func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(rsaPrivateKey)
}

// ValidateJWT parses and verifies JWT token using RSA public key
func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return rsaPublicKey, nil
	})
}
