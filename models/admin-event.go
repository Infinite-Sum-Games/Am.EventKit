package models

import (
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type AddEventDetailsRequest struct {
	Name        string `json:"name"`
	Blurb       string `json:"blurb"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
	Price       int32  `json:"price"`
	IsPerHead   bool   `json:"is_per_head"`
}

func (r AddEventDetailsRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Name, v.Required, v.RuneLength(3, 200)),
		// v.Field(&r.Blurb),
		// v.Field(&r.Description),
		// v.Field(&r.Rules),
		v.Field(&r.IsPerHead, v.In(true, false)),
	)
}

type AddEventPosterRequest struct {
	PosterUrl string `json:"poster_url"`
}

func (r AddEventPosterRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.PosterUrl, v.Required, is.URL),
	)
}

type AddEventDimensionRequest struct {
	IsGroup     bool `json:"is_group"`
	MinTeamSize int  `json:"min_teamsize"`
	MaxTeamSize int  `json:"max_teamsize"`
	TotalSeats  int  `json:"total_seats"`
	IsPerHead   bool `json:"is_per_head"`
}

func (r AddEventDimensionRequest) Validate() error {
	if err := v.ValidateStruct(&r,
		v.Field(&r.IsGroup, v.In(true, false)),
		v.Field(&r.MinTeamSize, v.Min(0)),
		v.Field(&r.MaxTeamSize, v.Min(0)),
		v.Field(&r.TotalSeats, v.Required, v.Min(1)),
	); err != nil {
		return err
	}

	// Logical checks
	if r.IsGroup && r.MinTeamSize > r.MaxTeamSize {
		return v.Errors{
			"min_teamsize": v.NewError("validation", "must be <= max_teamsize"),
		}
	}
	return nil
}

type AddEventTogglesRequest struct {
	EventType      string `json:"event_type"`
	AttendanceMode string `json:"attendance_mode"`
	IsOffline      bool   `json:"is_offline"`
	IsTechnical    bool   `json:"is_technical"`
}

func (r AddEventTogglesRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.EventType, v.Required, v.In("EVENT", "WORKSHOP")),
		v.Field(&r.AttendanceMode, v.Required, v.In("SOLO", "DUO")),
		v.Field(&r.IsOffline, v.In(true, false)),
		v.Field(&r.IsTechnical, v.In(true, false)),
	)
}

type ConnectEventAndOrganizerRequest struct {
	EventId     string `json:"id"`
	OrganizerId string `json:"organizer_id"`
}

func (r ConnectEventAndOrganizerRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.EventId, v.Required, is.UUID),
		v.Field(&r.OrganizerId, v.Required, is.UUID),
	)
}

type DisconnectEventAndOrganizerRequest struct {
	EventId     string `json:"id"`
	OrganizerId string `json:"organizer_id"`
}

func (r DisconnectEventAndOrganizerRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.EventId, v.Required, is.UUID),
		v.Field(&r.OrganizerId, v.Required, is.UUID),
	)
}

type ConnectEventAndTagsRequest struct {
	EventId string `json:"id"`
	TagId   string `json:"tag_id"`
}

func (r ConnectEventAndTagsRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.EventId, v.Required, is.UUID),
		v.Field(&r.TagId, v.Required, is.UUID),
	)
}

type DisconnectEventAndTagsRequest struct {
	EventId string `json:"id"`
	TagId   string `json:"tag_id"`
}

func (r DisconnectEventAndTagsRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.EventId, v.Required, is.UUID),
		v.Field(&r.TagId, v.Required, is.UUID),
	)
}

type ConnectEventAndPeopleRequest struct {
	EventId  string `json:"id"`
	PersonId string `json:"person_id"`
}

func (r ConnectEventAndPeopleRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.EventId, v.Required, is.UUID),
		v.Field(&r.PersonId, v.Required, is.UUID),
	)
}

type DisconnectEventAndPeopleRequest struct {
	EventId  string `json:"id"`
	PersonId string `json:"person_id"`
}

func (r DisconnectEventAndPeopleRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.EventId, v.Required, is.UUID),
		v.Field(&r.PersonId, v.Required, is.UUID),
	)
}

type AddEventScheduleRequest struct {
	EventDate time.Time `json:"event_date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Venue     string    `json:"venue"`
}

func (r AddEventScheduleRequest) Validate() error {
	return validateSchedule(
		r.EventDate,
		r.StartTime,
		r.EndTime,
		r.Venue,
	)
}

type EditEventScheduleRequest struct {
	EventDate time.Time `json:"event_date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Venue     string    `json:"venue"`
}

func (r EditEventScheduleRequest) Validate() error {
	return validateSchedule(
		r.EventDate,
		r.StartTime,
		r.EndTime,
		r.Venue,
	)
}

func validateSchedule(eventDate, startTime, endTime time.Time, venue string) error {
	if err := v.ValidateStruct(&struct {
		EventDate time.Time
		StartTime time.Time
		EndTime   time.Time
		Venue     string
	}{
		eventDate,
		startTime,
		endTime,
		venue,
	},
		v.Field(&eventDate, v.Required),
		v.Field(&startTime, v.Required),
		v.Field(&endTime, v.Required),
		v.Field(&venue, v.Required, v.RuneLength(2, 200)),
	); err != nil {
		return err
	}

	if !endTime.After(startTime) {
		return v.Errors{
			"end_time": v.NewError("validation", "must be after start_time"),
		}
	}

	return nil
}
