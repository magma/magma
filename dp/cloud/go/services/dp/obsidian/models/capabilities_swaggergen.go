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

// Capabilities capabilities
//
// swagger:model capabilities
type Capabilities struct {

	// this is the maximum allowed difference in MHz between leftmost end of leftmost channel and rightmost end of rightmost channel used by a Base Station (eNB)
	// Required: true
	// Maximum: 150
	// Multiple Of: 5
	MaxIbwMhz int64 `json:"max_ibw_mhz"`

	// max tx power available on cbsd
	// Example: 30
	// Required: true
	MaxPower *float64 `json:"max_power"`

	// min tx power available on cbsd
	// Required: true
	MinPower *float64 `json:"min_power"`

	// number of antennas
	// Example: 2
	// Required: true
	// Minimum: 1
	NumberOfAntennas int64 `json:"number_of_antennas"`
}

// Validate validates this capabilities
func (m *Capabilities) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMaxIbwMhz(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMaxPower(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMinPower(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNumberOfAntennas(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Capabilities) validateMaxIbwMhz(formats strfmt.Registry) error {

	if err := validate.Required("max_ibw_mhz", "body", int64(m.MaxIbwMhz)); err != nil {
		return err
	}

	if err := validate.MaximumInt("max_ibw_mhz", "body", m.MaxIbwMhz, 150, false); err != nil {
		return err
	}

	if err := validate.MultipleOfInt("max_ibw_mhz", "body", m.MaxIbwMhz, 5); err != nil {
		return err
	}

	return nil
}

func (m *Capabilities) validateMaxPower(formats strfmt.Registry) error {

	if err := validate.Required("max_power", "body", m.MaxPower); err != nil {
		return err
	}

	return nil
}

func (m *Capabilities) validateMinPower(formats strfmt.Registry) error {

	if err := validate.Required("min_power", "body", m.MinPower); err != nil {
		return err
	}

	return nil
}

func (m *Capabilities) validateNumberOfAntennas(formats strfmt.Registry) error {

	if err := validate.Required("number_of_antennas", "body", int64(m.NumberOfAntennas)); err != nil {
		return err
	}

	if err := validate.MinimumInt("number_of_antennas", "body", m.NumberOfAntennas, 1, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this capabilities based on context it is used
func (m *Capabilities) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Capabilities) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Capabilities) UnmarshalBinary(b []byte) error {
	var res Capabilities
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}