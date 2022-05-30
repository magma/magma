// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// Msisdn Mobile station international subscriber directory number
// Example: 13109976224
//
// swagger:model msisdn
type Msisdn string

// Validate validates this msisdn
func (m Msisdn) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.MinLength("", "body", string(m), 1); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this msisdn based on context it is used
func (m Msisdn) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}