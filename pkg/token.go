package pkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ApplicationClaims struct {
	role string `json: "role"`
	name string `json: "name"`
	jwt.RegisteredClaims
}

func SignToken(email string, name string, userId, role string, tokenType string) (string, error) {
	var expiryTime time.Time
	switch tokenType {
	case "temp":
		expiryTime = time.Now().Add(30 * time.Minute)
	case "access":
		expiryTime = time.Now().Add(4 * time.Hour)
	case "refresh":
		expiryTime = time.Now().Add(90 * 24 * time.Hour)
	default:
		return "", fmt.Errorf("invalid tokenType provided. valid types: %s, %s or %s",
			"temp", "access", "refresh")
	}
	claims := ApplicationClaims{
		role,
		name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "anokha.backend.official",
			Audience:  []string{email},
			ID:        userId,
			Subject:   tokenType,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//TODO: use the jwt secret from env once merged...for now using random string
	signedToken, err := token.SignedString("Random String")
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}
	return signedToken, nil
}
