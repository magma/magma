/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/lib/go/errors"
)

func TestMultiError(t *testing.T) {
	// Correct Multi error implementation tests
	if returnNil() != nil {
		t.Error("'if err == nil' check should succeed for returned error")
	}
	if returnNilMulti().AsError() != nil {
		t.Error("'if err == nil' check should succeed for returned MultiError")
	}
	me := errors.NewMulti()
	assert.Nil(t, me)
	var e error = nil
	me = me.Add(e)
	assert.Nil(t, me)
	if me == nil {
		t.Error("'if err == nil' is not expected to succeed for multi error")
	}
	if me.AsError() != nil {
		t.Error("'if err == nil' check should succeed for converted error")
	}
	me = me.Add(errors.ErrNotFound)
	me = me.Add(errors.ErrNotFound)
	assert.NotNil(t, me)
	assert.NotNil(t, me.AsError())
	assert.GreaterOrEqual(t, len(me.AsError().Error()), 20)
	assert.Equal(t, me.AsError().Error(), me.Error())
	assert.Equal(t, 2, len(me.Get()))

	multi := returnMulti(errors.ErrNotFound, errors.ErrAlreadyExists, nil)
	assert.NotNil(t, multi)
	assert.Equal(t, 2, len(multi.Get()))
	assert.NotEmpty(t, multi.Error())

	multiErr := returnMultiError(errors.ErrNotFound, nil, errors.ErrAlreadyExists)
	assert.NotNil(t, multiErr)
	assert.NotEmpty(t, multiErr.Error())

	// test AddFmt
	me = errors.NewMulti().AddFmt(fmt.Errorf("foo bar"), "Multi error has %s (%d)", "one", 1)
	assert.Len(t, me.Get(), 1)
	assert.Equal(t, "Multi error has one (1) foo bar", me.Get()[0].Error())
	assert.Equal(t, "Multi error has one (1) foo bar", me.Error())
	me = me.AddFmt(fmt.Errorf("foo bars"), "Multi error has %s (%d)", "two", 2)
	assert.Len(t, me.Get(), 2)
	assert.Equal(t, "Multi error has two (2) foo bars", me.Get()[1].Error())
	assert.Equal(t, "errors: [0: Multi error has one (1) foo bar; 1: Multi error has two (2) foo bars]", me.Error())

	// Wrong Multi error Implementation tests
	// check that Cast() 'fixing' wrong implementation
	if returnNilWrong() == nil {
		t.Error("'if err == nil' check should not succeed for returned error with wrong implementation")
	}
	if returnNilMultiWrong().AsError() == nil {
		t.Error("'if err == nil' check should not succeed for returned MultiError with wrong implementation")
	}
	// but - Cast should fix the 'if nil' issues
	if errors.Cast(returnNilWrong()) == nil {
		t.Error("'if err == nil' check should succeed for returned error with wrong implementation")
	}
	if errors.Cast(returnNilMultiWrong().AsError()) == nil {
		t.Error("'if err == nil' check should succeed for returned MultiError with wrong implementation")
	}

	multi = returnMultiWrong(errors.ErrNotFound, errors.ErrAlreadyExists, nil)
	assert.NotNil(t, multi)
	assert.Equal(t, 3, len(multi.Get()))
	assert.NotEmpty(t, multi.Error())

	multiErr = returnMultiErrorWrong(errors.ErrNotFound, nil, errors.ErrAlreadyExists)
	assert.NotNil(t, multiErr)
	assert.NotEmpty(t, multiErr.Error())

}

func returnNil() error {
	return errors.NewMulti().AsError()
}

func returnNilMulti() errors.Multi {
	return errors.NewMulti()
}

func returnMulti(e1, e2, e3 error) errors.Multi {
	return errors.NewMulti(e1, e2, e3)
}

func returnMultiError(e1, e2, e3 error) error {
	return errors.NewMulti(e1, e2, e3)
}

// IncorrectMultiError Multi error interface implementation
type IncorrectMultiError struct {
	Errors []error
}

// Error returns a formatted string for MultiError list
func (me *IncorrectMultiError) Error() string {
	if me == nil {
		return "<nil>"
	}
	return fmt.Sprintf("IncorrectMultiError: %v", me.Errors)
}

// Set - sets multi set to errs
func (me *IncorrectMultiError) Set(errs ...error) errors.Multi {
	if me != nil {
		errs = append(me.Errors, errs...)
	}
	return &IncorrectMultiError{Errors: errs}
}

// Get - returns a list of chained errors
func (me *IncorrectMultiError) Get() []error {
	if me != nil {
		return me.Errors
	}
	return nil
}

// Add appends errs to the existing MultiError set
func (me *IncorrectMultiError) Add(errs ...error) errors.Multi {
	if me == nil {
		return me.Set(errs...)
	}
	me.Errors = append(me.Errors, errs...)
	return me
}

// AddFmt adds a new formatted error if err is not nil, it's a noop if err == nil
func (me *IncorrectMultiError) AddFmt(err error, _ string, _ ...interface{}) errors.Multi {
	return me.Add(err)
}

func (me *IncorrectMultiError) AsError() error {
	return me
}

func returnNilWrong() error {
	var res *IncorrectMultiError
	return res.AsError()
}

func returnNilMultiWrong() errors.Multi {
	var res *IncorrectMultiError
	return res
}

func returnMultiWrong(e1, e2, e3 error) errors.Multi {
	return &IncorrectMultiError{Errors: []error{e1, e2, e3}}
}

func returnMultiErrorWrong(e1, e2, e3 error) error {
	return &IncorrectMultiError{Errors: []error{e1, e2, e3}}
}
