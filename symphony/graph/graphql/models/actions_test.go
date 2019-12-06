// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"github.com/facebookincubator/symphony/cloud/actions"
	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestMagmaTriggersExist ensures that all triggers defined in graphql
// have an associated actions trigger
func TestMagmaTriggersExist(t *testing.T) {
	registry := actions.MainRegistry()

	for _, graphqlTriggerID := range AllTriggerID {
		actionsTriggerID := core.TriggerID(graphqlTriggerID)
		trigger, err := registry.TriggerForID(actionsTriggerID)
		assert.NoError(t, err)
		assert.NotNil(t, trigger)
	}
}

func TestMagmaActionsExist(t *testing.T) {
	registry := actions.MainRegistry()

	for _, graphqlActionID := range AllActionID {
		actionsActionID := core.ActionID(graphqlActionID)
		action, err := registry.ActionForID(actionsActionID)
		assert.NoError(t, err)
		assert.NotNil(t, action)
	}
}
