// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mockaction

import (
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/stretchr/testify/mock"
)

// Action is a mock action used for testing
type Action struct {
	mock.Mock
}

// New returns a new action
func New() *Action {
	return &Action{}
}

// ID returns the ID of the action
func (m *Action) ID() core.ActionID {
	args := m.Called()
	return args.Get(0).(core.ActionID)
}

// Description returns the description of the action
func (m *Action) Description() string {
	args := m.Called()
	return args.String(0)
}

// DataType implements core.Action.DataType()
func (m *Action) DataType() core.DataType {
	args := m.Called()
	return args.Get(0).(core.DataType)
}

// Execute executes the action
func (m *Action) Execute(ac core.ActionContext) error {
	args := m.Called(ac)
	return args.Error(0)
}
