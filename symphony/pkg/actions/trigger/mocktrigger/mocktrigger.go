// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mocktrigger

import (
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/stretchr/testify/mock"
)

// Trigger is a mock trigger
type Trigger struct {
	mock.Mock
}

// New returns a new trigger
func New() *Trigger {
	return &Trigger{}
}

// ID returns the ID of the trigger
func (m *Trigger) ID() core.TriggerID {
	args := m.Called()
	return args.Get(0).(core.TriggerID)
}

// Description returns the description of the trigger
func (m *Trigger) Description() string {
	args := m.Called()
	return args.String(0)
}

// SupportedActionIDs returns the supported actions for this trigger
func (m *Trigger) SupportedActionIDs() []core.ActionID {
	args := m.Called()
	return args.Get(0).([]core.ActionID)
}

// SupportedFilters returns the supported actions for this trigger
func (m *Trigger) SupportedFilters() []core.Filter {
	args := m.Called()
	return args.Get(0).([]core.Filter)
}

// Evaluate runs evaluations of this trigger
func (m *Trigger) Evaluate(r core.Rule, inputParams map[string]interface{}) (bool, error) {
	args := m.Called(r)
	return args.Bool(0), args.Error(1)
}
