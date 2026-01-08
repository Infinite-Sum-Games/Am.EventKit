package models

import (
	"strings"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UpdateDisputeStatusInput struct {
	StudentEmail string `json:"student_email"`
	Description  string `json:"description"`
}

func (u *UpdateDisputeStatusInput) Validate() error {
	u.StudentEmail = strings.ToLower(u.StudentEmail)
	if err := v.ValidateStruct(u,
		v.Field(&u.StudentEmail, is.Email),
		v.Field(&u.Description),
	); err != nil {
		return err
	}
	return nil
}
