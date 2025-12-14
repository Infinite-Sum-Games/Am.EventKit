package models

import (
	"time"
)

type AddEventDetailsRequest struct {
	Name        string `json:"name"`
	Blurb       string `json:"blurb"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
	Price       int    `json:"price"`
	IsPerHead   bool   `json:"is_per_head"`
}

func (r AddEventDetailsRequest) Validate() error {
	return nil
}

type AddEventPosterRequest struct {
	PosterUrl string `json:"poster_url"`
}

func (r AddEventPosterRequest) Validate() error {
	return nil
}

type AddEventDimensionRequest struct {
	IsGroup     bool `json:"is_group"`
	MinTeamSize int  `json:"min_teamsize"`
	MaxTeamSize int  `json:"max_teamsize"`
	TotalSeats  int  `json:"total_seats"`
}

func (r AddEventDimensionRequest) Validate() error {
	return nil
}

type AddEventTogglesRequest struct {
	IsGroup            bool `json:"is_team_event"`
	IsPublished        bool `json:"is_published"`
	IsRegistrationOpen bool `json:"is_registration_open"`
}

func (r AddEventTogglesRequest) Validate() error {
	return nil
}

type ConnectEventAndOrganizerRequest struct {
	EventId     string `json:"id"`
	OrganizerId string `json:"organizer_id"`
}

func (r ConnectEventAndOrganizerRequest) Validate() error {
	return nil
}

type DisconnectEventAndOrganizerRequest struct {
	EventId     string `json:"id"`
	OrganizerId string `json:"organizer_id"`
}

func (r DisconnectEventAndOrganizerRequest) Validate() error {
	return nil
}

type ConnectEventAndTags struct {
	EventId string `json:"id"`
	TagId   string `json:"tag_id"`
}

func (r ConnectEventAndTags) Validate() error {
	return nil
}

type DisconnectEventAndTagsRequest struct {
	EventId string `json:"id"`
	TagID   string `json:"tag_id"`
}

func (r DisconnectEventAndTagsRequest) Validate() error {
	return nil
}

type AddNewEventScheduleRequest struct {
	Venue string `json:"venue"`
}

func (r AddNewEventScheduleRequest) Validate() error {
	return nil
}

type EditEventScheduleRequest struct {
	Round    string    `json:"round"`
	Datetime time.Time `json:"datetime"`
	Venue    string    `json:"venue"`
}

func (r EditEventScheduleRequest) Validate() error {
	return nil
}
