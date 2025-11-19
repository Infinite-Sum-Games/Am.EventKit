package pkg

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/gin-gonic/gin"
)

/*
	JTI stores the email
	Audience stores the username
*/

const (
	RefreshTokenValidTime = time.Hour * 24 * 90
	AuthTokenValidTime    = time.Hour
	TempTokenValidTime    = time.Minute * 5
	CsrfTokenValidTime    = time.Minute * 5
	privateKeyPath        = "app.rsa"
	publicKeyPath         = "app.rsa.pub"
)

var (
	VerifyKey paseto.V4AsymmetricPublicKey
	SignKey   paseto.V4AsymmetricSecretKey
)

func InitPaseto() error {
	privateKeyBinary, err := os.ReadFile("app.rsa")
	if err != nil {
		return err
	}
	privateKeyHex := hex.EncodeToString(privateKeyBinary)

	publicKeyBinary, err := os.ReadFile("app.pub.rsa")
	if err != nil {
		return err
	}
	publicKeyHex := hex.EncodeToString(publicKeyBinary)

	// Verify using public key
	VerifyKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publicKeyHex)
	if err != nil {
		return fmt.Errorf("Error in public-paseto: %w", err)
	}
	// Sign using private key
	SignKey, err = paseto.NewV4AsymmetricSecretKeyFromHex(privateKeyHex)
	if err != nil {
		return fmt.Errorf("Error in private-paseto: %w", err)
	}
	return nil
}

func CreateAuthToken(id, username, email string, user, host, organizer bool) string {
	token := paseto.NewToken()
	token.SetJti(email)
	token.SetAudience(username)
	token.SetIssuer("Anokha-25: AUTH-SERVICE")
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(AuthTokenValidTime))
	token.SetSubject("access_token")
	if err := token.Set("USER-ID", id); err != nil {
		Log.Error("[AUTH-ERROR]: Failed to set USER-ID claim", err)
	}
	if err := token.Set("STUDENT-ROLE", user); err != nil {
		Log.Error("[AUTH-ERROR]: Failed to set STUDENT-ROLE claim", err)
	}
	if err := token.Set("ORGANIZER-ROLE", organizer); err != nil {
		Log.Error("[AUTH-ERROR]: Failed to set ORGANIZER-ROLE claim", err)
	}

	signed := token.V4Sign(SignKey, nil)
	return signed
}

func CreateRefreshToken(id, username, email string, user, host, organizer bool) string {
	token := paseto.NewToken()
	token.SetJti(email)
	token.SetAudience(username)
	token.SetIssuer("Anokha-25: AUTH-SERVICE")
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(RefreshTokenValidTime))
	token.SetSubject("refresh_token")
	if err := token.Set("USER-ID", id); err != nil {
		Log.Error("[AUTH-ERROR]: Failed to set USER-ID claim", err)
	}
	if err := token.Set("STUDENT-ROLE", user); err != nil {
		Log.Error("[AUTH-ERROR]: Failed to set STUDENT-ROLE claim", err)
	}
	if err := token.Set("ORGANIZER-ROLE", organizer); err != nil {
		Log.Error("[AUTH-ERROR]: Failed to set ORGANIZER-ROLE claim", err)
	}

	signed := token.V4Sign(SignKey, nil)
	return signed
}

func CreateTempToken(username, email string) string {
	token := paseto.NewToken()
	token.SetJti(email)
	token.SetAudience(username)
	token.SetIssuer("Anokha-25: AUTH-SERVICE")
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(TempTokenValidTime))
	token.SetSubject("temp_token")

	signed := token.V4Sign(SignKey, nil)
	return signed
}

func CreateCsrfToken(email string, purpose string) string {
	token := paseto.NewToken()
	token.SetJti(email)
	token.SetIssuer("Anokha-25: AUTH-SERVICE")
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(CsrfTokenValidTime))
	token.SetSubject(purpose)

	signed := token.V4Sign(SignKey, nil)
	return signed
}

func ParseToken(token, tokeType string) (bool, *paseto.Token) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy("Anokha-25: AUTH-SERVICE"))
	parser.AddRule(paseto.Subject(tokeType))
	parser.AddRule(paseto.ValidAt(time.Now()))
	parser.AddRule(paseto.NotExpired())

	parsedToken, err := parser.ParseV4Public(VerifyKey, token, nil)
	if err != nil {
		return false, nil
	}
	return true, parsedToken
}

func VerifyTokens(c *gin.Context, authToken, refreshToken string) bool {
	ok, parsedAuthToken := ParseToken(authToken, "access_token")
	if !ok {
		if _, err := VerifyRefreshToken(c, refreshToken); err != nil {
			return false
		}
		return false
	}
	ok, parsedRefToken := ParseToken(refreshToken, "refresh_token")
	if !ok {
		return false
	}

	authData := parsedAuthToken.Claims()
	refData := parsedRefToken.Claims()

	// Verification conditions
	c1 := authData["audience"] != refData["audience"]
	c2 := authData["jti"] != refData["jti"]
	c3 := authData["USER-ROLE"] != refData["USER-ROLE"]
	c4 := authData["ORGANIZER-ROLE"] != refData["ORGANIZER-ROLE"]
	c5 := authData["USER-ID"] != refData["USER-ID"]

	if c1 || c2 || c3 || c4 || c5 {
		return false
	}

	// Setting up variables in *gin.Context for passing around in handlers
	c.Set("userId", authData["USER-ID"])
	c.Set("username", authData["audience"])
	c.Set("email", authData["jti"])
	c.Set("USER-ROLE", authData["USER-ROLE"])
	c.Set("ORGANIZER-ROLE", authData["ORGANIZER-ROLE"])

	return true
}

func VerifyTempToken(c *gin.Context, tempToken string) bool {
	ok, parsedTempToken := ParseToken(tempToken, "temp_token")
	if !ok {
		return false
	}

	tempData := parsedTempToken.Claims()
	c.Set("username", tempData["audience"])
	c.Set("email", tempData["jti"])

	return true
}

func VerifyRefreshToken(c *gin.Context, refreshToken string) (*paseto.Token, error) {
	email, ok := c.Get("email")
	if !ok {
		Log.WarnCtx(c, "[GIN-ERROR]: Email not passed down in context")
		return nil, fmt.Errorf("could not fetch email from gin.Context")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := cmd.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	q := db.New()
	token, err := q.CheckRefreshTokenQuery(ctx, conn, fmt.Sprintf("%v", email))

	// Possible scenarios
	// 1. RefreshToken does not exist
	// 2. RefreshToken has become invalid
	// 3. RefreshToken is perfect and it can generate AuthToken
	if err != nil {
		Log.FatalCtx(c, "[AUTH-ERROR] Failed to fetch refresh token from DB", err)
		return nil, err
	}
	if token.String == "" {
		return nil, fmt.Errorf("[AUTH-ERROR] Refresh token not available in DB")
	}

	ok, validToken := ParseToken(token.String, "refresh_token")
	if !ok {
		return nil, fmt.Errorf("[AUTH-ERROR]: Failed to parse refresh token")
	}

	return validToken, nil
}
