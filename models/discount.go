package models

import (
	"regexp"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateDiscountRequest struct {
	DiscountType        string `json:"discount_type"`
	DiscountedSoloSeats *int   `json:"discounted_solo_seats"` // nullable
	DiscountedTeamSeats *int   `json:"discounted_team_seats"` // nullable
	StartTime           string `json:"start_time"`            // ISO string
	EndTime             string `json:"end_time"`
}

type EditDiscountRequest struct {
	DiscountType        string `json:"discount_type"`
	DiscountedSoloSeats *int   `json:"discounted_solo_seats"`
	DiscountedTeamSeats *int   `json:"discounted_team_seats"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
}

func (r CreateDiscountRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.DiscountType, v.Required, v.Length(2, 50), v.Match(regexp.MustCompile(`^[A-Za-z_]+$`))),
		v.Field(&r.DiscountedSoloSeats, v.When(r.DiscountedSoloSeats != nil, v.Min(0))),
		v.Field(&r.DiscountedTeamSeats, v.When(r.DiscountedTeamSeats != nil, v.Min(0))),
		v.Field(&r.StartTime, v.Required),
		v.Field(&r.EndTime, v.Required),
	)
}

func (r EditDiscountRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.DiscountType, v.Required, v.Length(2, 50), v.Match(regexp.MustCompile(`^[A-Za-z_]+$`))),
		v.Field(&r.DiscountedSoloSeats, v.When(r.DiscountedSoloSeats != nil, v.Min(0))),
		v.Field(&r.DiscountedTeamSeats, v.When(r.DiscountedTeamSeats != nil, v.Min(0))),
		v.Field(&r.StartTime, v.Required),
		v.Field(&r.EndTime, v.Required),
	)
}
