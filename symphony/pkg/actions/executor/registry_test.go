// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"testing"

	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/actions/trigger/mocktrigger"
	"github.com/stretchr/testify/assert"
)

func TestRegistryTriggers(t *testing.T) {
	var triggerID core.TriggerID = "id123"

	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(triggerID)

	registry := NewRegistry()
	registry.MustRegisterTrigger(trigger1)

	assert.Equal(t, 1, len(registry.Triggers()))
	assert.Equal(t, triggerID, registry.Triggers()[0].ID())

	result, err := registry.TriggerForID(triggerID)
	assert.NoError(t, err)
	assert.Equal(t, trigger1, result)
}
