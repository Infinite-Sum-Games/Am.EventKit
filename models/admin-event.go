package models

import (
	"time"

	"github.com/google/uuid"
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
	MinTeamSize int `json:"min_team_size"`
	MaxTeamSize int `json:"max_team_size"`
}

func (r AddEventDimensionRequest) Validate() error {
	return nil
}

type AddEventTogglesRequest struct {
	EventId            uuid.UUID `json:"event_id"`
	IsTeamEvent        *bool     `json:"is_team_event"`
	IsPublished        *bool     `json:"is_published"`
	IsRegistrationOpen *bool     `json:"is_registration_open"`
}

func (r AddEventTogglesRequest) Validate() error {
	return nil
}

type ConnectEventAndOrganizerRequest struct {
	EventId     uuid.UUID `json:"event_id"`
	OrganizerId uuid.UUID `json:"organizer_id"`
}

func (r ConnectEventAndOrganizerRequest) Validate() error {
	return nil
}

type DisconnectEventAndOrganizerRequest struct {
	EventId     uuid.UUID `json:"event_id"`
	OrganizerId uuid.UUID `json:"organizer_id"`
}

func (r DisconnectEventAndOrganizerRequest) Validate() error {
	return nil
}

type ConnectEventAndTags struct {
	EventId uuid.UUID   `json:"event_id"`
	TagIDs  []uuid.UUID `json:"tag_ids"`
}

func (r ConnectEventAndTags) Validate() error {
	return nil
}

type DisconnectEventAndTagsRequest struct {
	EventId uuid.UUID `json:"event_id"`
	TagID   uuid.UUID `json:"tag_id"`
}

func (r *DisconnectEventAndTagsRequest) Validate() error {
	return nil
}

type AttachNewEventScheduleRequest struct {
	EventId  uuid.UUID `json:"event_id"`
	Round    string    `json:"round"`
	Datetime time.Time `json:"datetime"`
	Venue    string    `json:"venue"`
}

func (r *AttachNewEventScheduleRequest) Validate() error {
	return nil
}

type EditEventScheduleRequest struct {
	ScheduleId uuid.UUID `json:"schedule_id"`
	Round      string    `json:"round"`
	Datetime   time.Time `json:"datetime"`
	Venue      string    `json:"venue"`
}

func (r *EditEventScheduleRequest) Validate() error {
	return nil
}

type DeleteEventScheduleRequest struct {
	ScheduleId uuid.UUID `json:"schedule_id"`
}

func (r *DeleteEventScheduleRequest) Validate() error {
	return nil
}

