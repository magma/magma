/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package errors

import (
	"bytes"
	"fmt"
)

// MultiError defines an error wrapper for multiple errors
type MultiError struct {
	Errors []error
}

// Error returns a formatted string for MultiError list
func (me *MultiError) Error() string {
	if me == nil {
		return "<nil>"
	}
	switch len(me.Errors) {
	case 0:
		return ""
	case 1:
		return me.Errors[0].Error()
	default:
		var b bytes.Buffer
		fmtStr := "errors: [%d: %v"
		for i, e := range me.Errors {
			fmt.Fprintf(&b, fmtStr, i, e)
			fmtStr = "; %d: %v"
		}
		b.Write([]byte("]"))
		return b.String()
	}
}

// Set - sets multi set to errs
func (me *MultiError) Set(errs ...error) *MultiError {
	var errors []error
	for _, e := range errs {
		if e != nil {
			errors = append(errors, e)
		}
	}
	if me == nil {
		if len(errors) == 0 {
			return me
		}
		return &MultiError{Errors: errors}
	}
	me.Errors = errors
	return me
}

// Add appends errs to the existing MultiError set
func (me *MultiError) Add(errs ...error) *MultiError {
	if me == nil {
		return me.Set(errs...)
	}
	for _, e := range errs {
		if e != nil {
			me.Errors = append(me.Errors, e)
		}
	}
	return me
}
