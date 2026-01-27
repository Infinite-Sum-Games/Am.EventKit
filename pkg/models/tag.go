package models

import v "github.com/go-ozzo/ozzo-validation/v4"

type CreateTagRequest struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

func (r CreateTagRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Name, v.Required, v.RuneLength(3, 50)),
		v.Field(&r.Abbreviation, v.Required, v.RuneLength(3, 10)),
	)
}
