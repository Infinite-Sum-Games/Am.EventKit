package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type TeamBookingRequest struct {
	TeamName    string       `json:"team_name" binding:"required"`
	TeamMembers []TeamMember `json:"team_members" binding:"required"`
}

type TeamMember struct {
	StudentEmail string `json:"student_email" binding:"required"`
	StudentRole  string `json:"student_role" binding:"required"`
}

// Validation for TeamMember
func (t TeamMember) Validate() error {
	return v.ValidateStruct(&t,
		v.Field(&t.StudentEmail, v.Required, v.Length(3, 200), is.Email),
		v.Field(&t.StudentRole, v.Required, v.Length(1, 100)),
	)
}

func (e TeamBookingRequest) Validate() error {
	return v.ValidateStruct(&e,
		v.Field(&e.TeamName, v.Required, v.Length(2, 100)),
		v.Field(&e.TeamMembers,
			v.Required,
			v.Length(1, 100),
			v.Each(v.By(func(value interface{}) error {
				if tm, ok := value.(TeamMember); ok {
					return tm.Validate()
				}
				return nil
			})),
		),
	)
}
