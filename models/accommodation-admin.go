package models

import (
	"regexp"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type AddHostelRequest struct {
	RoomCount       int32  `json:"room_count"`
	IsMale          bool   `json:"is_male"`
	WardenEmail     string `json:"warden_email"`
	Latitude        string `json:"latitude"`
	Longtitude      string `json:"longtitude"`
	MapUrl          string `json:"map_url"`
	HostelName      string `json:"hostel_name"`
	DayScholarPrice int32  `json:"day_scholar_price"`
	OutsiderPrice   int32  `json:"outsider_price"`
}

func (r AddHostelRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.RoomCount, v.Required, v.Min(1)),
		v.Field(&r.IsMale),
		v.Field(&r.WardenEmail, v.Required, is.Email),
		v.Field(&r.Latitude),
		v.Field(&r.Longtitude),
		v.Field(&r.MapUrl, is.URL),
		v.Field(
			&r.HostelName,
			v.Required,
			v.Match(
				regexp.MustCompile(
					`^[A-Z0-9 ]+ BHAVANAM - (SINGLE|DORM|4 SHARING)$`,
				),
			).Error(
				"hostel_name must be uppercase and end with 'BHAVANAM - SINGLE', 'BHAVANAM - DORM', or 'BHAVANAM - 4 SHARING'",
			),
		),
		v.Field(&r.DayScholarPrice, v.Required, v.Min(0)),
		v.Field(&r.OutsiderPrice, v.Required, v.Min(0)),
	)
}

type UpdateHostelRequest struct {
	HostelID        string `json:"hostel_id"`
	RoomCount       int32  `json:"room_count"`
	IsMale          bool   `json:"is_male"`
	WardenEmail     string `json:"warden_email"`
	Latitude        string `json:"latitude"`
	Longtitude      string `json:"longtitude"`
	MapUrl          string `json:"map_url"`
	DayScholarPrice int32  `json:"day_scholar_price"`
	OutsiderPrice   int32  `json:"outsider_price"`
}

func (r UpdateHostelRequest) Validate() error {
	decimalRegex := regexp.MustCompile(`^-?\d+(\.\d+)?$`)

	return v.ValidateStruct(&r,
		v.Field(&r.HostelID, v.Required, is.UUID),
		v.Field(&r.RoomCount, v.Min(0)),
		v.Field(&r.IsMale),
		v.Field(&r.WardenEmail, is.Email),
		v.Field(
			&r.Latitude,
			v.When(r.Latitude != "",
				v.Match(decimalRegex),
			),
		),
		v.Field(
			&r.Longtitude,
			v.When(r.Longtitude != "",
				v.Match(decimalRegex),
			),
		),
		v.Field(&r.MapUrl, is.URL),
		v.Field(&r.DayScholarPrice, v.Min(0)),
		v.Field(&r.OutsiderPrice, v.Min(0)),
	)
}

type AllotHostelRequest struct {
	HostelID string `json:"hostel_id"`
	DayCount int32  `json:"day_count"`
}

func (r AllotHostelRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.HostelID, v.Required, is.UUID),
		v.Field(&r.DayCount, v.Required, v.Min(1), v.Max(4)),
	)
}

type UpdateAccommodationByIdRequest struct {
	IsMale            bool   `json:"is_male"`
	IsHosteller       bool   `json:"is_hosteller"`
	CollegeRollNumber string `json:"college_roll_number"`
	CollegeName       string `json:"college_name"`
	RoomPreference    string `json:"room_preference"`
	IsAmritaCampus    bool   `json:"is_amrita_campus"`
	CheckInDate       string `json:"check_in_date"`
	CheckInTime       string `json:"check_in_time"`
	CheckOutDate      string `json:"check_out_date"`
	CheckOutTime      string `json:"check_out_time"`
}

func (r UpdateAccommodationByIdRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.IsMale),
		v.Field(&r.IsHosteller),
		v.Field(&r.CollegeRollNumber),
		v.Field(&r.CollegeName),
		v.Field(&r.RoomPreference),
		v.Field(&r.IsAmritaCampus),
		v.Field(&r.CheckInDate),
		v.Field(&r.CheckInTime),
		v.Field(&r.CheckOutDate),
		v.Field(&r.CheckOutTime),
	)
}

type MapQrStudentIdRequest struct {
	StudentID     string `json:"student_id"`
	HospitalityId string `json:"hospitality_id"`
}

func (r MapQrStudentIdRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.StudentID, v.Required, is.UUID),
		v.Field(
			&r.HospitalityId,
			v.Required,
			// v.Match(regexp.MustCompile(`^A\d{4}CBE$`)).
			// 	Error("hospitality_id must be in the format A1234CBE"),
			//P001
			v.Match(regexp.MustCompile(`^P\d{3}$`)).
				Error("hospitality_id must be in the format P001"),
		),
	)
}
