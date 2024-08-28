// Code generated by go-swagger; DO NOT EDIT.

package dags

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/dagu-org/dagu/internal/frontend/gen/models"
)

// CreateDagOKCode is the HTTP code returned for type CreateDagOK
const CreateDagOKCode int = 200

/*
CreateDagOK A successful response.

swagger:response createDagOK
*/
type CreateDagOK struct {

	/*
	  In: Body
	*/
	Payload *models.CreateDagResponse `json:"body,omitempty"`
}

// NewCreateDagOK creates CreateDagOK with default headers values
func NewCreateDagOK() *CreateDagOK {

	return &CreateDagOK{}
}

// WithPayload adds the payload to the create dag o k response
func (o *CreateDagOK) WithPayload(payload *models.CreateDagResponse) *CreateDagOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create dag o k response
func (o *CreateDagOK) SetPayload(payload *models.CreateDagResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateDagOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
CreateDagDefault Generic error response.

swagger:response createDagDefault
*/
type CreateDagDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.APIError `json:"body,omitempty"`
}

// NewCreateDagDefault creates CreateDagDefault with default headers values
func NewCreateDagDefault(code int) *CreateDagDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateDagDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create dag default response
func (o *CreateDagDefault) WithStatusCode(code int) *CreateDagDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create dag default response
func (o *CreateDagDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the create dag default response
func (o *CreateDagDefault) WithPayload(payload *models.APIError) *CreateDagDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create dag default response
func (o *CreateDagDefault) SetPayload(payload *models.APIError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateDagDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
