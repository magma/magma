// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import "errors"

// Filter is an interface for implementing a filter
type Filter interface {
	// A unique identifier contextual to trigger for the filter
	FilterID() string

	// Ex. "the alerts network id"
	Description() string

	// The operators supported by this filter
	SupportedOperators() []Operator

	// Returns true if the filter passes
	Evaluate(ruleFilter *ActionsRuleFilter, inputParams map[string]interface{}) (bool, error)
}

// StringFieldFilter is a generic filter for filtering string fields
type StringFieldFilter struct {
	fieldName   string
	description string
}

// NewStringFieldFilter creates a new filter for strings
func NewStringFieldFilter(fieldName string, description string) Filter {
	return &StringFieldFilter{
		fieldName:   fieldName,
		description: description,
	}
}

// FilterID implements the Filter interface
func (f *StringFieldFilter) FilterID() string {
	return "stringfieldfilter_" + f.fieldName
}

// Description implements the Filter interface
func (f *StringFieldFilter) Description() string {
	return f.description
}

// SupportedOperators implements the Filter interface
func (f *StringFieldFilter) SupportedOperators() []Operator {
	return []Operator{
		OperatorIsString,
		OperatorIsNotString,
	}
}

// Evaluate implements the Filter interface
func (f *StringFieldFilter) Evaluate(filter *ActionsRuleFilter, inputParams map[string]interface{}) (bool, error) {
	if filter.OperatorID == OperatorIsString.OperatorID() {
		return filter.Data == inputParams[f.fieldName], nil
	}
	if filter.OperatorID == OperatorIsNotString.OperatorID() {
		return filter.Data != inputParams[f.fieldName], nil
	}
	return false, errors.New("invalid operatorID")
}
