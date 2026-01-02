package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type AddHostelRequest struct {
	RoomCount   int32  `json:"room_count"`
	IsMale      bool   `json:"is_male"`
	WardenEmail string `json:"warden_email"`
	Latitude    string `json:"latitude"`
	Longtitude  string `json:"longtitude"`
	MapUrl      string `json:"map_url"`
	HostelName  string `json:"hostel_name"`
}

func (r AddHostelRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.RoomCount, v.Required, v.Min(1)),
		v.Field(&r.IsMale),
		v.Field(&r.WardenEmail, is.Email),
		v.Field(&r.Latitude),
		v.Field(&r.Longtitude),
		v.Field(&r.MapUrl, is.URL),
		v.Field(&r.HostelName, v.Required),
	)
}
