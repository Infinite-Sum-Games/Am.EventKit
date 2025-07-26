package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type StudentOnboardingRequest struct {
	Name             string `json:"name" binding:"required"`
	DepartmentName   string `json:"department_name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=8"`
	PhoneNumber      string `json:"phone_number" binding:"required"`
	IsAmritaStudent  bool   `json:"is_amrita_student"`
	AmritaRollNumber string `json:"amrita_roll_number,omitempty"`
	CollegeName      string `json:"college_name" binding:"required"`
	CollegeCity      string `json:"college_city" binding:"required"`
	AcademicYear     string `json:"academic_year" binding:"required"`
}

func (s StudentOnboardingRequest) Validate() error {
	return v.ValidateStruct(&s,
		v.Field(&s.Name, v.Required),
		v.Field(&s.DepartmentName, v.Required),
		v.Field(&s.Email, v.Required, is.Email),
		v.Field(&s.Password, v.Required, v.Length(8, 0)),
		v.Field(&s.PhoneNumber, v.Required, v.Length(10, 15)),
		v.Field(&s.CollegeName, v.Required),
		v.Field(&s.CollegeCity, v.Required),
		v.Field(&s.AcademicYear, v.Required))
}
