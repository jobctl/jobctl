// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostDagActionRequest Request body for posting an action to a DAG.
//
// swagger:model PostDagActionRequest
type PostDagActionRequest struct {

	// Action to be performed on the DAG.
	// Required: true
	// Enum: [start suspend stop retry mark-success mark-failed save rename]
	Action *string `json:"action"`

	// Additional parameters for the action.
	Params string `json:"params,omitempty"`

	// Unique request ID for the action.
	RequestID string `json:"requestId,omitempty"`

	// Step name if the action targets a specific step.
	Step string `json:"step,omitempty"`

	// Optional extra value for the action.
	Value string `json:"value,omitempty"`
}

// Validate validates this post dag action request
func (m *PostDagActionRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAction(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var postDagActionRequestTypeActionPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["start","suspend","stop","retry","mark-success","mark-failed","save","rename"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		postDagActionRequestTypeActionPropEnum = append(postDagActionRequestTypeActionPropEnum, v)
	}
}

const (

	// PostDagActionRequestActionStart captures enum value "start"
	PostDagActionRequestActionStart string = "start"

	// PostDagActionRequestActionSuspend captures enum value "suspend"
	PostDagActionRequestActionSuspend string = "suspend"

	// PostDagActionRequestActionStop captures enum value "stop"
	PostDagActionRequestActionStop string = "stop"

	// PostDagActionRequestActionRetry captures enum value "retry"
	PostDagActionRequestActionRetry string = "retry"

	// PostDagActionRequestActionMarkDashSuccess captures enum value "mark-success"
	PostDagActionRequestActionMarkDashSuccess string = "mark-success"

	// PostDagActionRequestActionMarkDashFailed captures enum value "mark-failed"
	PostDagActionRequestActionMarkDashFailed string = "mark-failed"

	// PostDagActionRequestActionSave captures enum value "save"
	PostDagActionRequestActionSave string = "save"

	// PostDagActionRequestActionRename captures enum value "rename"
	PostDagActionRequestActionRename string = "rename"
)

// prop value enum
func (m *PostDagActionRequest) validateActionEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, postDagActionRequestTypeActionPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *PostDagActionRequest) validateAction(formats strfmt.Registry) error {

	if err := validate.Required("action", "body", m.Action); err != nil {
		return err
	}

	// value enum
	if err := m.validateActionEnum("action", "body", *m.Action); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this post dag action request based on context it is used
func (m *PostDagActionRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PostDagActionRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostDagActionRequest) UnmarshalBinary(b []byte) error {
	var res PostDagActionRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
