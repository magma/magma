// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// Filter is an interface for implementing a filter
type Filter interface {
	// A unique identifier contextual to trigger for the filter
	FilterID() string

	// Ex. "the alerts network id"
	Description() string

	// The operators supported by this filter
	SupportedOperators() []Operator
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
func (sff *StringFieldFilter) FilterID() string {
	return "stringfieldfilter_" + sff.fieldName
}

// Description implements the Filter interface
func (sff *StringFieldFilter) Description() string {
	return sff.description
}

// SupportedOperators implements the Filter interface
func (sff *StringFieldFilter) SupportedOperators() []Operator {
	return []Operator{
		OperatorIsString,
		OperatorIsNotString,
	}
}
