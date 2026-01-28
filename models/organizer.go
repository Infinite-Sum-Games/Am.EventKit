package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateOrganizerRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	OrgType       string `json:"org_type"` // DEPARTMENT | CLUB
	StudentHead   string `json:"student_head"`
	StudentCoHead string `json:"student_co_head"` // optional
	FacultyHead   string `json:"faculty_head"`
}

func (r CreateOrganizerRequest) Validate() error {
	return v.ValidateStruct(&r,

		// Required fields
		v.Field(&r.Name, v.Required, v.RuneLength(3, 100)),
		v.Field(&r.Email, v.Required, is.Email),
		v.Field(&r.Password, v.Required, v.RuneLength(6, 100)),
		v.Field(&r.OrgType, v.Required, v.In("DEPARTMENT", "CLUB")),
		v.Field(&r.StudentHead, v.Required, v.RuneLength(3, 100)),
		v.Field(&r.FacultyHead, v.Required, v.RuneLength(3, 100)),

		// Optional fields
		v.Field(&r.StudentCoHead, v.When(r.StudentCoHead != "", v.RuneLength(3, 100))),
	)
}

type EditOrganizerRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	OrgType       string `json:"org_type"` // DEPARTMENT | CLUB
	StudentHead   string `json:"student_head"`
	StudentCoHead string `json:"student_co_head"` // optional
	FacultyHead   string `json:"faculty_head"`
}

func (r EditOrganizerRequest) Validate() error {
	return v.ValidateStruct(&r,

		v.Field(&r.Name, v.Required, v.RuneLength(3, 100)),
		v.Field(&r.Email, v.Required, is.Email),
		v.Field(&r.OrgType, v.Required, v.In("DEPARTMENT", "CLUB")),
		v.Field(&r.StudentHead, v.Required, v.RuneLength(3, 100)),
		v.Field(&r.FacultyHead, v.Required, v.RuneLength(3, 100)),
		v.Field(&r.StudentCoHead, v.When(r.StudentCoHead != "", v.RuneLength(3, 100))),
	)
}

type ChangeOrganizerPasswordRequest struct {
	Password string `json:"password"`
}

func (r ChangeOrganizerPasswordRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Password, v.Required, v.RuneLength(3, 100)),
	)
}
