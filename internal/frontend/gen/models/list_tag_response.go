// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ListTagResponse list tag response
//
// swagger:model ListTagResponse
type ListTagResponse struct {

	// errors
	// Required: true
	Errors []string `json:"Errors"`

	// tags
	// Required: true
	Tags []string `json:"Tags"`
}

// Validate validates this list tag response
func (m *ListTagResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateErrors(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ListTagResponse) validateErrors(formats strfmt.Registry) error {

	if err := validate.Required("Errors", "body", m.Errors); err != nil {
		return err
	}

	return nil
}

func (m *ListTagResponse) validateTags(formats strfmt.Registry) error {

	if err := validate.Required("Tags", "body", m.Tags); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this list tag response based on context it is used
func (m *ListTagResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ListTagResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ListTagResponse) UnmarshalBinary(b []byte) error {
	var res ListTagResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
