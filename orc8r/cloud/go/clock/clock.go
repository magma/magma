/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package clock

import (
	"testing"
	"time"
)

var c Clock = &DefaultClock{}

// Now returns the current time or what it's been set to
func Now() time.Time {
	return c.Now()
}

// SetAndFreezeClock will set the value to be returned by Now()
// This should only be called by test code (hence the required but unused
// *testing.T parameter)
func SetAndFreezeClock(t *testing.T, ti time.Time) {
	if t == nil {
		panic("nice try")
	}
	c = &MockClock{mockTime: ti}
}

// UnfreezeClock will revert clock.Now()'s behavior to delegating to time.Now()
func UnfreezeClock(t *testing.T) {
	r := recover()
	if t == nil {
		panic("nice try")
	}
	c = &DefaultClock{}
	if r != nil {
		panic(r)
	}
}

// Clock is an interface for getting the current time
type Clock interface {
	// Now returns the current time (or what it's been set to)
	Now() time.Time
}

// DefaultClock is a Clock implementation which wraps time.Now()
type DefaultClock struct{}

func (d *DefaultClock) Now() time.Time {
	return time.Now()
}

// MockClock is a Clock implementation which always returns a fixed time
type MockClock struct {
	mockTime time.Time
}

func (m *MockClock) Now() time.Time {
	return m.mockTime
}
