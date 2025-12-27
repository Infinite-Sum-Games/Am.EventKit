package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type CreateDisputeInput struct {
	EventId       uuid.UUID `json:"event_id"`
	TransactionId uuid.UUID `json:"transaction_id"`
}

func (d CreateDisputeInput) Validate() error {
	if err := v.ValidateStruct(&d,
		v.Field(&d.EventId, v.Required, is.UUID),
		v.Field(&d.TransactionId, v.Required, is.UUID),
	); err != nil {
		return err
	}
	return nil
}

type UpdateDisputeStatusInput struct {
	StudentEmail      string         `json:"student_email"`
	Description       string         `json:"description"`
	TeamMemberDetails map[string]any `json:"team_member_details"`
}

func (u UpdateDisputeStatusInput) Validate() error {
	if err := v.ValidateStruct(&u,
		v.Field(&u.StudentEmail, is.Email),
		v.Field(&u.Description),
		v.Field(&u.TeamMemberDetails),
	); err != nil {
		return err
	}
	return nil
}
