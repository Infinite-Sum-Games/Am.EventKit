package auth

type IAuth interface {
	LoginUser(email, password string) error
	Logout(email string) error
	LoginAdmin(email, password string) error
	ForgetUserPassword(email string) error
	ConfirmPasswordchange() error
	ResendPasswordChangeOtp() error
}
