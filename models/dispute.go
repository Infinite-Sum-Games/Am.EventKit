package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UpdateDisputeStatusInput struct {
	StudentEmail      string         `json:"student_email"`
	Description       string         `json:"description"`
	TeamMemberDetails map[string]any `json:"team_member_details"`
}

func (u UpdateDisputeStatusInput) Validate() error {
	if err := v.ValidateStruct(&u,
		v.Field(&u.StudentEmail, is.Email),
		v.Field(&u.Description),
		v.Field(&u.TeamMemberDetails),
	); err != nil {
		return err
	}
	return nil
}
