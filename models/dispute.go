package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateDisputeStatusInput struct {
	Description string `json:"description"`
}

func (u UpdateDisputeStatusInput) Validate() error {
	if err := v.ValidateStruct(&u,
		v.Field(&u.Description),
	); err != nil {
		return err
	}
	return nil
}
