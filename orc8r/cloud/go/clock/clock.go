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

var (
	c clock   = &defaultClock{}
	s sleeper = &defaultSleep{}
)

// Now returns the current time or the time to which it's been set.
func Now() time.Time {
	return c.now()
}

// Sleep either the specified duration or a small, negligible duration.
func Sleep(d time.Duration) {
	s.sleep(d)
}

// SetAndFreezeClock will set the value to be returned by Now.
// This should only be called by test code.
func SetAndFreezeClock(t *testing.T, ti time.Time) {
	if t == nil {
		panic("for tests only")
	}
	c = &mockClock{mockTime: ti}
}

// UnfreezeClock will revert clock.Now's behavior to delegating to time.Now.
// This should only be called by test code.
func UnfreezeClock(t *testing.T) {
	r := recover()
	if t == nil {
		panic("for tests only")
	}
	c = &defaultClock{}
	if r != nil {
		panic(r)
	}
}

// SkipSleeps causes time.Sleep to sleep for only a small, negligible duration.
// This should only be used for test code.
func SkipSleeps(t *testing.T) {
	if t == nil {
		panic("for tests only")
	}
	s = &mockSleep{}
}

// ResumeSleeps causes time.Sleep to resume default behavior.
// This should only be used for test code.
func ResumeSleeps(t *testing.T) {
	if t == nil {
		panic("for tests only")
	}
	s = &defaultSleep{}
}
