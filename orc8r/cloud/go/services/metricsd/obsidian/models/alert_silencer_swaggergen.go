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

// AlertSilencer alert silencer
//
// swagger:model alert_silencer
type AlertSilencer struct {

	// comment
	// Required: true
	Comment *string `json:"comment"`

	// created by
	// Required: true
	CreatedBy *string `json:"createdBy"`

	// ends at
	// Example: 2019-10-17T22:19:41.990Z
	// Required: true
	EndsAt *string `json:"endsAt"`

	// matchers
	// Required: true
	Matchers []*Matcher `json:"matchers"`

	// starts at
	// Example: 2019-10-17T22:19:41.990Z
	// Required: true
	StartsAt *string `json:"startsAt"`
}

// Validate validates this alert silencer
func (m *AlertSilencer) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateComment(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedBy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEndsAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMatchers(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartsAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertSilencer) validateComment(formats strfmt.Registry) error {

	if err := validate.Required("comment", "body", m.Comment); err != nil {
		return err
	}

	return nil
}

func (m *AlertSilencer) validateCreatedBy(formats strfmt.Registry) error {

	if err := validate.Required("createdBy", "body", m.CreatedBy); err != nil {
		return err
	}

	return nil
}

func (m *AlertSilencer) validateEndsAt(formats strfmt.Registry) error {

	if err := validate.Required("endsAt", "body", m.EndsAt); err != nil {
		return err
	}

	return nil
}

func (m *AlertSilencer) validateMatchers(formats strfmt.Registry) error {

	if err := validate.Required("matchers", "body", m.Matchers); err != nil {
		return err
	}

	for i := 0; i < len(m.Matchers); i++ {
		if swag.IsZero(m.Matchers[i]) { // not required
			continue
		}

		if m.Matchers[i] != nil {
			if err := m.Matchers[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("matchers" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("matchers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *AlertSilencer) validateStartsAt(formats strfmt.Registry) error {

	if err := validate.Required("startsAt", "body", m.StartsAt); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this alert silencer based on the context it is used
func (m *AlertSilencer) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMatchers(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertSilencer) contextValidateMatchers(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Matchers); i++ {

		if m.Matchers[i] != nil {
			if err := m.Matchers[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("matchers" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("matchers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *AlertSilencer) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AlertSilencer) UnmarshalBinary(b []byte) error {
	var res AlertSilencer
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}