// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"testing"

	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/stretchr/testify/assert"
)

// TestMagmaTriggersExist ensures that all triggers defined in graphql
// have an associated actions trigger, at the same index (string set
// equivalence)
func TestMagmaTriggersExist(t *testing.T) {
	assert.Equal(t, len(core.AllTriggerIDs), len(AllTriggerID))

	for i, triggerID := range AllTriggerID {
		assert.Equal(t, string(triggerID), string(core.AllTriggerIDs[i]))
	}
}

func TestMagmaActionsExist(t *testing.T) {
	assert.Equal(t, len(core.AllActionIDs), len(AllActionID))

	for i, actionID := range AllActionID {
		assert.Equal(t, string(actionID), string(core.AllActionIDs[i]))
	}
}
