package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CheckEmailRequest struct {
	Email string `json:"email"`
}

func (s CheckEmailRequest) Validate() error {
	return v.ValidateStruct(&s,
		v.Field(&s.Email, v.Required, is.Email, is.LowerCase),
	)
}

type StudentOnboardingRequest struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	PhoneNumber      string `json:"phone_number"`
	IsAmritaStudent  bool   `json:"is_amrita_student"`
	AmritaRollNumber string `json:"amrita_roll_number"`
	CollegeName      string `json:"college_name"`
	CollegeCity      string `json:"college_city"`
}

func (s StudentOnboardingRequest) Validate() error {
	return v.ValidateStruct(&s,
		v.Field(&s.Name, v.Required, v.Length(3, 50)),
		v.Field(&s.Email, v.Required, is.Email, is.LowerCase),
		v.Field(&s.Password, v.Required, v.Length(8, 0)),
		v.Field(&s.PhoneNumber, v.Required, v.Length(10, 10)),
		v.Field(&s.CollegeName, v.Required, v.Length(3, 128)),
		v.Field(&s.CollegeCity, v.Required, v.Length(3, 50)),
	)
}

type LoginRequest struct {
	Email          string `json:"email"`
	HashedPassword string `json:"password"`
}

func (l LoginRequest) Validate() error {
	return v.ValidateStruct(&l,
		v.Field(&l.Email, v.Required, is.Email, is.LowerCase),
		v.Field(&l.HashedPassword, v.Required))
}

type OtpRequest struct {
	Otp string `json:"otp"`
}

func (s OtpRequest) Validate() error {
	return v.ValidateStruct(&s,
		v.Field(&s.Otp, v.Required, v.Length(4, 6), is.Digit),
	)
}

type ForgetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

func (f ForgetPasswordRequest) Validate() error {
	return v.ValidateStruct(&f,
		v.Field(&f.Email, v.Required, is.Email, is.LowerCase),
		v.Field(&f.NewPassword, v.Required, v.Length(8, 0)))
}
