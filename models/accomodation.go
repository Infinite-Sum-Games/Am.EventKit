package models

import (
	"regexp"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type AccomodationForm struct {
	IsAmritaCampus    bool   `json:"is_amrita_campus"`
	IsMale            bool   `json:"is_male"` // Gender
	IsHosteller       bool   `json:"is_hosteller"`
	RoomPreference    string `json:"room_preference"` // Single | 4 Sharing | Dormitory
	CollegeName       string `json:"college_name"`
	CollegeRollNumber string `json:"college_roll_number"`
	CheckInDate       string `json:"check_in_date"`
	CheckInTime       string `json:"check_in_time"`
	CheckOutDate      string `json:"check_out_date"`
	CheckOutTime      string `json:"check_out_time"`
}

func (r AccomodationForm) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.RoomPreference, v.Required, v.In("single", "4 sharing", "dormitory")),
		v.Field(&r.CollegeName, v.Required.When(!r.IsAmritaCampus)),
		v.Field(&r.CollegeRollNumber, v.Required.When(r.IsAmritaCampus == false)),
		v.Field(&r.CheckInDate, v.Required, v.Date("2006-01-02")),
		v.Field(&r.CheckInTime, v.Required, v.Match(regexp.MustCompile("^([01][0-9]|2[0-3]):[0-5][0-9]$"))),
		v.Field(&r.CheckOutDate, v.Required, v.Date("2006-01-02")),
		v.Field(&r.CheckOutTime, v.Required, v.Match(regexp.MustCompile("^([01][0-9]|2[0-3]):[0-5][0-9]$"))),
	)
}
