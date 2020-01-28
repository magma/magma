// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/facebookincubator/symphony/pkg/actions/action/mockaction"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/actions/trigger/mocktrigger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testActionID1 core.ActionID = "action1"
	testActionID2 core.ActionID = "action2"

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

func (m *MockLogger) LogError(ctx context.Context, err error) {
	m.Called(ctx, err)
}

func TestExecutor(t *testing.T) {
	networkIDFilter := core.NewStringFieldFilter("networkID", "a string filter")

	testRule := core.Rule{
		ID:        "rule1",
		TriggerID: testTriggerID1,
		RuleActions: []*core.ActionsRuleAction{
			{
				ActionID: testActionID1,
				Data:     "testdata",
			},
		},
		RuleFilters: []*core.ActionsRuleFilter{
			{
				FilterID:   networkIDFilter.FilterID(),
				OperatorID: core.OperatorIsString.OperatorID(),
				Data:       "network1",
			},
		},
	}

	testErrRule := core.Rule{
		ID:        "errRule",
		TriggerID: testTriggerID2,
		RuleActions: []*core.ActionsRuleAction{
			{
				ActionID: testActionID2,
				Data:     "testdata",
			},
		},
	}

	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(testTriggerID1)
	trigger1.On("SupportedFilters").Return([]core.Filter{networkIDFilter})

	trigger2 := mocktrigger.New()
	trigger2.On("ID").Return(testTriggerID2)
	trigger2.On("SupportedFilters").Return([]core.Filter{networkIDFilter})

	action := mockaction.New()
	action.On("ID").Return(testActionID1)
	action.On("Execute", mock.Anything).Return(nil)

	errAction := mockaction.New()
	errAction.On("ID").Return(testActionID2)
	errAction.On("Execute", mock.Anything).Return(errors.New("error"))

	registry := NewRegistry()
	registry.MustRegisterAction(action)
	registry.MustRegisterAction(errAction)
	registry.MustRegisterTrigger(trigger1)
	registry.MustRegisterTrigger(trigger2)

	exc := Executor{
		Registry: registry,
		DataLoader: BasicDataLoader{
			Rules: []core.Rule{testRule, testErrRule},
		},
		OnError: func(ctx context.Context, err error) {},
	}

	triggers := map[core.TriggerID]map[string]interface{}{
		testTriggerID1: payload1,
		testTriggerID2: payload1,
	}
	res := exc.Execute(context.Background(), "id123", triggers)
	assert.Equal(t, []core.ActionID{action.ID()}, res.Successes)
	assert.Equal(t, []ExecutionError{{
		ID:    errAction.ID(),
		Error: "executing action2: error",
	}}, res.Errors)

	action.AssertNumberOfCalls(t, "Execute", 1)
	errAction.AssertNumberOfCalls(t, "Execute", 1)
}

func TestExecutorRuleFilter(t *testing.T) {
	networkIDFilter := core.NewStringFieldFilter("networkID", "a string filter")

	testRule := core.Rule{
		ID:        "rule1",
		TriggerID: testTriggerID1,
		RuleActions: []*core.ActionsRuleAction{
			{
				ActionID: testActionID1,
				Data:     "testdata",
			},
		},
		RuleFilters: []*core.ActionsRuleFilter{
			{
				FilterID:   networkIDFilter.FilterID(),
				OperatorID: core.OperatorIsString.OperatorID(),
				Data:       "anotherNetwork",
			},
		},
	}

	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(testTriggerID1)
	trigger1.On("SupportedFilters").Return([]core.Filter{networkIDFilter})

	action := mockaction.New()
	action.On("ID").Return(testActionID1)
	action.On("Execute", mock.Anything).Return(nil)

	registry := NewRegistry()
	registry.MustRegisterAction(action)
	registry.MustRegisterTrigger(trigger1)

	exc := Executor{
		Registry: registry,
		DataLoader: BasicDataLoader{
			Rules: []core.Rule{testRule},
		},
		OnError: func(ctx context.Context, err error) {
			assert.Fail(t, "error in test when shouldnt be", err)
		},
	}

	triggers := map[core.TriggerID]map[string]interface{}{
		testTriggerID1: payload1,
	}
	res := exc.Execute(context.Background(), "id123", triggers)
	assert.Len(t, res.Successes, 0)
	assert.Len(t, res.Errors, 0)

	action.AssertNumberOfCalls(t, "Execute", 0)
}

func TestExecutorUnregisteredTrigger(t *testing.T) {
	var mockErrorHandler MockLogger
	mockErrorHandler.On("LogError", mock.Anything, mock.Anything).Return()

	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(testTriggerID2)

	registry := NewRegistry()
	registry.MustRegisterTrigger(trigger1)

	exc := Executor{
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
	res := exc.Execute(context.Background(), "id123", triggers)
	assert.Len(t, res.Successes, 0)
	assert.Len(t, res.Errors, 0)

	mockErrorHandler.AssertNumberOfCalls(t, "LogError", 1)
}
