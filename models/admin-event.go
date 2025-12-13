package models

import "github.com/google/uuid"

type AddEventDetailsRequest struct {
	EventId uuid.UUID `json:"event_id"`
}

func (r *AddEventDetailsRequest) Validate() error {
	return nil
}

type AddEventPosterRequest struct {
	EventId        uuid.UUID `json:"event_id"`
	EventPosterUrl string    `json:"poster_url"`
}

func (r *AddEventPosterRequest) Validate() error {
	return nil
}

type DeletePosterRequest struct {
	EventId uuid.UUID `json:"event_id"`
}

func (r *DeletePosterRequest) Validate() error {
	return nil
}

type AddEventDimensionRequest struct {
}

func (r *AddEventDimensionRequest) Validate() error {
	return nil
}

type AddEventTogglesRequest struct {
}

func (r *AddEventTogglesRequest) Validate() error {
	return nil
}

type ConnectEventAndOrganizerRequest struct {
}

func (r *ConnectEventAndOrganizerRequest) Validate() error {
	return nil
}

type DisconnectEventAndOrganizerRequest struct {
}

func (r *DisconnectEventAndOrganizerRequest) Validate() error {
	return nil
}

type ConnectEventAndTags struct {
}

func (r *ConnectEventAndTags) Validate() error {
	return nil
}

type DisconnectEventAndTagsRequest struct {
}

func (r *DisconnectEventAndTagsRequest) Validate() error {
	return nil
}

type AttachNewEventScheduleRequest struct {
}

func (r *AttachNewEventScheduleRequest) Validate() error {
	return nil
}

type EditEventScheduleRequest struct {
}

func (r *EditEventScheduleRequest) Validate() error {
	return nil
}

type DeleteEventScheduleRequest struct {
}

func (r *DeleteEventScheduleRequest) Validate() error {
	return nil
}