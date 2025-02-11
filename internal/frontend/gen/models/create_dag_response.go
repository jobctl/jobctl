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

// CreateDagResponse create dag response
//
// swagger:model CreateDagResponse
type CreateDagResponse struct {

	// dag ID
	// Required: true
	DagID *string `json:"DagID"`
}

// Validate validates this create dag response
func (m *CreateDagResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDagID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreateDagResponse) validateDagID(formats strfmt.Registry) error {

	if err := validate.Required("DagID", "body", m.DagID); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this create dag response based on context it is used
func (m *CreateDagResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreateDagResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateDagResponse) UnmarshalBinary(b []byte) error {
	var res CreateDagResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
