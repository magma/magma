// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/cloud/actions/action/mockaction"
	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/facebookincubator/symphony/cloud/actions/trigger/mocktrigger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testActionID1 core.ActionID = "action1"

	testTriggerID1 core.TriggerID = "trigger1"
	testTriggerID2 core.TriggerID = "trigger2"

	payload1 = map[string]interface{}{
		"networkID": "network1",
		"gatewayID": "gateway1",
	}
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogError(err error) {
	m.Called(err)
}

func TestExecutor(t *testing.T) {

	testRule := core.Rule{
		ID:        "rule1",
		TriggerID: testTriggerID1,
		ActionIDs: []core.ActionID{testActionID1},
	}

	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(testTriggerID1)
	trigger1.On("Evaluate", testRule).Return(true, nil)

	trigger2 := mocktrigger.New()
	trigger2.On("ID").Return(testTriggerID2)
	trigger2.On("Evaluate", testRule).Return(true, nil)

	action := mockaction.New()
	action.On("ID").Return(testActionID1)
	action.On("Execute", mock.Anything).Return(nil)

	registry := NewRegistry()
	registry.MustRegisterAction(action)
	registry.MustRegisterTrigger(trigger1)
	registry.MustRegisterTrigger(trigger2)

	exc := Executor{
		Context:  context.Background(),
		Registry: registry,
		DataLoader: BasicDataLoader{
			Rules: []core.Rule{testRule},
		},
		OnError: func(err error) {
			assert.Fail(t, "error in test when shouldnt be", err)
		},
	}

	triggers := map[core.TriggerID]map[string]interface{}{
		testTriggerID1: payload1,
		testTriggerID2: payload1,
	}
	exc.Execute(context.Background(), "id123", triggers)

	action.AssertNumberOfCalls(t, "Execute", 1)
}

func TestExecutorUnregisteredTrigger(t *testing.T) {

	var mockErrorHandler MockLogger
	mockErrorHandler.On("LogError", mock.Anything).Return()

	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(testTriggerID2)

	registry := NewRegistry()
	registry.MustRegisterTrigger(trigger1)

	exc := Executor{
		Context:  context.Background(),
		Registry: registry,
		DataLoader: BasicDataLoader{
			Rules: []core.Rule{},
		},
		OnError: mockErrorHandler.LogError,
	}

	triggers := map[core.TriggerID]map[string]interface{}{
		testTriggerID1: payload1,
		testTriggerID2: payload1,
	}
	exc.Execute(context.Background(), "id123", triggers)

	mockErrorHandler.AssertNumberOfCalls(t, "LogError", 1)
}
