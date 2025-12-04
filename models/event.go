package models

import (
	"time"

	"github.com/google/uuid"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type Event struct {
	ID   uuid.UUID `json:"event_id"`
	Name string    `json:"event_name"`
}

type EventScheduleInput struct {
	EventDate string `json:"event_date"` // YYYY-MM-DD
	StartTime string `json:"start_time"` // RFC3339 or HH:MM
	EndTime   string `json:"end_time"`   // RFC3339 or HH:MM
	Venue     string `json:"venue"`
}

func (s EventScheduleInput) Validate() error {
	if err := v.ValidateStruct(&s,
		v.Field(&s.EventDate, v.Required, v.Length(10, 10)),
		v.Field(&s.StartTime, v.Required, v.Length(4, 50)),
		v.Field(&s.EndTime, v.Required, v.Length(4, 50)),
		v.Field(&s.Venue, v.Required, v.RuneLength(2, 200)),
	); err != nil {
		return err
	}

	if _, err := time.Parse("2006-01-02", s.EventDate); err != nil {
		return err
	}

	if _, err := time.Parse(time.RFC3339, s.StartTime); err != nil {
		if _, err2 := time.Parse("15:04", s.StartTime); err2 != nil {
			return err
		}
	}

	if _, err := time.Parse(time.RFC3339, s.EndTime); err != nil {
		if _, err2 := time.Parse("15:04", s.EndTime); err2 != nil {
			return err
		}
	}

	return nil
}

type CreateEventRequest struct {
	Name           string               `json:"name"`
	Blurb          string               `json:"blurb"`
	Description    string               `json:"description"`
	CoverImageURL  string               `json:"cover_image_url"`
	Price          float64              `json:"price"`
	IsPerHead      bool                 `json:"is_per_head"`
	Rules          string               `json:"rules"`
	EventType      string               `json:"event_type"` // EVENT | WORKSHOP
	IsGroup        bool                 `json:"is_group"`
	MaxTeamSize    int                  `json:"max_teamsize"`
	MinTeamSize    int                  `json:"min_teamsize"`
	TotalSeats     int32                `json:"total_seats"`
	SeatsFilled    int32                `json:"seats_filled"`
	EventStatus    string               `json:"event_status"`    // CLOSED | ACTIVE | COMPLETED
	EventMode      string               `json:"event_mode"`      // ONLINE | OFFLINE
	AttendanceMode string               `json:"attendance_mode"` // SOLO | DUO
	OrganizerIDs   []string             `json:"organizer_ids"`   // UUIDs
	TagIDs         []string             `json:"tag_ids"`         // UUIDs
	PeopleIDs      []string             `json:"people_ids"`      // UUIDs
	Schedules      []EventScheduleInput `json:"schedules"`
}

func (r CreateEventRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Name, v.Required, v.RuneLength(3, 200)),
		v.Field(&r.Blurb, v.Required, v.RuneLength(3, 500)),
		v.Field(&r.Description, v.Required, v.RuneLength(3, 4000)),
		v.Field(&r.Price, v.Required),
		v.Field(&r.Rules, v.Required, v.RuneLength(1, 4000)),
		v.Field(&r.EventType, v.Required, v.In("EVENT", "WORKSHOP")),
		v.Field(&r.TotalSeats, v.Required, v.Min(1)),
		v.Field(&r.EventStatus, v.Required, v.In("CLOSED", "ACTIVE", "COMPLETED")),
		v.Field(&r.EventMode, v.Required, v.In("ONLINE", "OFFLINE")),
		v.Field(&r.AttendanceMode, v.Required, v.In("SOLO", "DUO")),
		v.Field(&r.OrganizerIDs, v.Required),
		v.Field(&r.TagIDs, v.Required),
		v.Field(&r.Schedules, v.Required),
	)
}

type UpdateEventRequest struct {
	Name           string               `json:"name"`
	Blurb          string               `json:"blurb"`
	Description    string               `json:"description"`
	CoverImageURL  string               `json:"cover_image_url"`
	Price          float64              `json:"price"`
	IsPerHead      bool                 `json:"is_per_head"`
	Rules          string               `json:"rules"`
	EventType      string               `json:"event_type"`
	IsGroup        bool                 `json:"is_group"`
	MaxTeamSize    int                  `json:"max_teamsize"`
	MinTeamSize    int                  `json:"min_teamsize"`
	TotalSeats     int32                `json:"total_seats"`
	SeatsFilled    int32                `json:"seats_filled"`
	EventStatus    string               `json:"event_status"`
	EventMode      string               `json:"event_mode"`
	AttendanceMode string               `json:"attendance_mode"`
	OrganizerIDs   []string             `json:"organizer_ids"`
	TagIDs         []string             `json:"tag_ids"`
	PeopleIDs      []string             `json:"people_ids"`
	Schedules      []EventScheduleInput `json:"schedules"`
}

func (r UpdateEventRequest) Validate() error {
	// reuse same rules as create
	return CreateEventRequest(r).Validate()
}
