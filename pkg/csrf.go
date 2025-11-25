package pkg

import (
	"fmt"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/gin-gonic/gin"
)

var CsrfRoutes = map[string]string{
	"/api/v1/auth/user/login":                      "user.login",
	"/api/v1/auth/user/register":                   "user.register",
	"/api/v1/auth/user/register/otp/verify":        "user.otp_verify",
	"/api/v1/auth/organizer/login":                 "organizer.login",
	"/api/v1/auth/user/forgot-password":            "user.forgot_password",
	"/api/v1/auth/user/forgot-password/otp/verify": "user.forgot_password_otp",
	"/api/v1/user/profile":                         "user.edit_profile",
}

func CreateCsrfToken(email string, c *gin.Context) (string, error) {
	csrfPurpose := CsrfRoutes[c.FullPath()]
	if csrfPurpose == "" || email == "" {
		return "", fmt.Errorf("[CSRF-ERROR]: bad params in CreateCsrfToken()")
	}

	token := paseto.NewToken()

	token.SetJti(email)
	token.SetSubject("csrf_token")
	token.SetIssuer("Anokha-25: AUTH-SERVICE")
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(CsrfTokenValidTime))
	token.SetAudience(csrfPurpose)

	signed := token.V4Sign(SignKey, nil)
	return signed, nil

}
