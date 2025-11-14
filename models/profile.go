package models

import (
	"regexp"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type EditProfileRequest struct {
	Name           string `json:"name"`
	DepartmentName string `json:"department_name"`
	PhoneNumber    string `json:"phone_number"`
	CollegeName    string `json:"college_name"`
	CollegeCity    string `json:"college_city"`
	AcademicYear   string `json:"academic_year"`
}

func (e EditProfileRequest) Validate() error {
	return v.ValidateStruct(&e,
		v.Field(&e.Name, v.Required, v.Length(2, 100)),
		v.Field(&e.DepartmentName, v.Required, v.Length(2, 100)),
		v.Field(&e.PhoneNumber, v.Required, v.Match(regexp.MustCompile(`^[0-9]{10}$`))),
		v.Field(&e.CollegeName, v.Required, v.Length(2, 200)),
		v.Field(&e.CollegeCity, v.Required, v.Length(2, 100)),
		v.Field(&e.AcademicYear, v.Required, v.In("1", "2", "3", "4", "5")),
	)
}
