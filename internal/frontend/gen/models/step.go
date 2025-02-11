// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Step Individual task within a DAG that performs a specific operation
//
// swagger:model Step
type Step struct {

	// List of arguments to pass to the command
	// Required: true
	Args []string `json:"Args"`

	// Complete command string including arguments to execute
	// Required: true
	CmdWithArgs *string `json:"CmdWithArgs"`

	// Base command to execute without arguments
	// Required: true
	Command *string `json:"Command"`

	// List of step names that must complete before this step can start
	// Required: true
	Depends []string `json:"Depends"`

	// Human-readable description of what the step does
	// Required: true
	Description *string `json:"Description"`

	// Working directory for executing the step's command
	// Required: true
	Dir *string `json:"Dir"`

	// Whether to send email notifications on step failure
	// Required: true
	MailOnError *bool `json:"MailOnError"`

	// Unique identifier for the step within the DAG
	// Required: true
	Name *string `json:"Name"`

	// Variable name to store the step's output
	// Required: true
	Output *string `json:"Output"`

	// Parameters to pass to the sub DAG
	Params string `json:"Params,omitempty"`

	// Conditions that must be met before the step can start
	// Required: true
	Preconditions []*Precondition `json:"Preconditions"`

	// repeat policy
	// Required: true
	RepeatPolicy *RepeatPolicy `json:"RepeatPolicy"`

	// Sub DAG to run
	Run string `json:"Run,omitempty"`

	// Script content if the step executes a script file
	// Required: true
	Script *string `json:"Script"`

	// File path for capturing standard error
	// Required: true
	Stderr *string `json:"Stderr"`

	// File path for capturing standard output
	// Required: true
	Stdout *string `json:"Stdout"`
}

// Validate validates this step
func (m *Step) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateArgs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCmdWithArgs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCommand(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDepends(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDir(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMailOnError(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOutput(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePreconditions(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRepeatPolicy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateScript(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStderr(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStdout(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Step) validateArgs(formats strfmt.Registry) error {

	if err := validate.Required("Args", "body", m.Args); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateCmdWithArgs(formats strfmt.Registry) error {

	if err := validate.Required("CmdWithArgs", "body", m.CmdWithArgs); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateCommand(formats strfmt.Registry) error {

	if err := validate.Required("Command", "body", m.Command); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateDepends(formats strfmt.Registry) error {

	if err := validate.Required("Depends", "body", m.Depends); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateDescription(formats strfmt.Registry) error {

	if err := validate.Required("Description", "body", m.Description); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateDir(formats strfmt.Registry) error {

	if err := validate.Required("Dir", "body", m.Dir); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateMailOnError(formats strfmt.Registry) error {

	if err := validate.Required("MailOnError", "body", m.MailOnError); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateName(formats strfmt.Registry) error {

	if err := validate.Required("Name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateOutput(formats strfmt.Registry) error {

	if err := validate.Required("Output", "body", m.Output); err != nil {
		return err
	}

	return nil
}

func (m *Step) validatePreconditions(formats strfmt.Registry) error {

	if err := validate.Required("Preconditions", "body", m.Preconditions); err != nil {
		return err
	}

	for i := 0; i < len(m.Preconditions); i++ {
		if swag.IsZero(m.Preconditions[i]) { // not required
			continue
		}

		if m.Preconditions[i] != nil {
			if err := m.Preconditions[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Preconditions" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("Preconditions" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Step) validateRepeatPolicy(formats strfmt.Registry) error {

	if err := validate.Required("RepeatPolicy", "body", m.RepeatPolicy); err != nil {
		return err
	}

	if m.RepeatPolicy != nil {
		if err := m.RepeatPolicy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("RepeatPolicy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("RepeatPolicy")
			}
			return err
		}
	}

	return nil
}

func (m *Step) validateScript(formats strfmt.Registry) error {

	if err := validate.Required("Script", "body", m.Script); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateStderr(formats strfmt.Registry) error {

	if err := validate.Required("Stderr", "body", m.Stderr); err != nil {
		return err
	}

	return nil
}

func (m *Step) validateStdout(formats strfmt.Registry) error {

	if err := validate.Required("Stdout", "body", m.Stdout); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this step based on the context it is used
func (m *Step) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePreconditions(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRepeatPolicy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Step) contextValidatePreconditions(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Preconditions); i++ {

		if m.Preconditions[i] != nil {

			if swag.IsZero(m.Preconditions[i]) { // not required
				return nil
			}

			if err := m.Preconditions[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Preconditions" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("Preconditions" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Step) contextValidateRepeatPolicy(ctx context.Context, formats strfmt.Registry) error {

	if m.RepeatPolicy != nil {

		if err := m.RepeatPolicy.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("RepeatPolicy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("RepeatPolicy")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Step) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Step) UnmarshalBinary(b []byte) error {
	var res Step
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
