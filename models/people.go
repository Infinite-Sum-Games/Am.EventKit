package models

import (
	"regexp"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateNewPerson struct {
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Profession  *string `json:"profession"`
	Email       *string `json:"email"`
}

func (p CreateNewPerson) Validate() error {
	return v.ValidateStruct(&p,
		v.Field(&p.Name, v.Required, v.Length(2, 100)),
		v.Field(&p.PhoneNumber, v.Required,
			v.Match(regexp.MustCompile(`^[0-9]{10}$`))),
		v.Field(&p.Profession, v.When(p.Profession != nil,
			v.Length(1, 100),
		)),
		v.Field(&p.Email, v.When(p.Email != nil,
			v.Length(5, 200),
			is.Email,
		)),
	)
}

type UpdatePersonEventRequest struct {
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Profession  *string `json:"profession"`
	Email       *string `json:"email"`
}

func (p UpdatePersonEventRequest) Validate() error {
	return v.ValidateStruct(&p,
		v.Field(&p.Name, v.When(
			p.Name != "",
			v.Length(2, 100),
		)),
		v.Field(&p.PhoneNumber, v.When(
			p.PhoneNumber != "",
			v.Length(10, 10),
			v.Match(regexp.MustCompile(`^[0-9]{10}$`)),
		)),
		v.Field(&p.Profession, v.When(
			p.Profession != nil,
			v.Length(1, 100),
		)),
		v.Field(&p.Email, v.When(p.Email != nil,
			v.Length(5, 200),
			is.Email,
		)),
	)
}
