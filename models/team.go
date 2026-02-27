package models

import (
	"strings"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type TeamBookingRequest struct {
	TeamName    string       `json:"team_name" binding:"required"`
	TeamMembers []TeamMember `json:"team_members" binding:"required"`
	ProblemStmt *string      `json:"ps,omitempty"`
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

func (e *TeamBookingRequest) Validate() error {
	for i := range e.TeamMembers {
		e.TeamMembers[i].StudentEmail = strings.ToLower(e.TeamMembers[i].StudentEmail)
	}
	return v.ValidateStruct(e,
		v.Field(&e.TeamName, v.Required, v.Length(2, 100)),
		v.Field(&e.TeamMembers,
			v.Length(0, 100),
			v.Each(v.By(func(value any) error {
				if tm, ok := value.(TeamMember); ok {
					return tm.Validate()
				}
				return nil
			})),
		),
		// asuming the only metadata we get is related to problem statement type
		v.Field(
			&e.ProblemStmt,
			v.When(
				e.ProblemStmt != nil,
				v.In("agentic_ai", "generative_ai", "aiot"),
			),
		),
	)
}
