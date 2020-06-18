// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphgrpc/schema"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/action/mockaction"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/actions/trigger/mocktrigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testNetworkID = "network1"
	testGatewayID = "gateway1"
)

func TestActionsAlertServiceClient_Receive(t *testing.T) {
	var testActionID1 core.ActionID = "action1"
	networkIDFilter := core.NewStringFieldFilter("networkID", "a string filter")
	gatewayIDFilter := core.NewStringFieldFilter("gatewayID", "a string filter")

	// Rule which triggers if a magma alert fires for `testNetworkID`
	testRule := core.Rule{
		ID:        "rule1",
		TriggerID: core.MagmaAlertTriggerID,
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
				Data:       testNetworkID,
			},
			{
				FilterID:   gatewayIDFilter.FilterID(),
				OperatorID: core.OperatorIsString.OperatorID(),
				Data:       testGatewayID,
			},
		},
	}

	// Mock magma alert trigger
	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(core.MagmaAlertTriggerID)
	trigger1.On("SupportedFilters").Return([]core.Filter{
		core.NewStringFieldFilter(
			"networkID",
			"the alerts networkID",
		),
		core.NewStringFieldFilter(
			"gatewayID",
			"the alerts gatewayID",
		),
	})

	// Mock action to be executed
	action := mockaction.New()
	action.On("ID").Return(testActionID1)
	action.On("Execute", mock.Anything).Return(nil)

	registry := executor.NewRegistry()
	registry.MustRegisterAction(action)
	registry.MustRegisterTrigger(trigger1)

	testExecutor := &executor.Executor{
		Registry: registry,
		DataLoader: executor.BasicDataLoader{
			Rules: []core.Rule{testRule},
		},
		OnError: func(ctx context.Context, err error) {
			assert.Fail(t, "error in test when there shouldn't be", err)
		},
	}

	as := NewActionsAlertService(func(ctx context.Context, tenantID string) (*actions.Client, error) {
		return actions.NewClient(testExecutor), nil
	})

	_, err := as.Trigger(context.Background(), &schema.AlertPayload{
		TenantID:  "tenant1",
		Alertname: "testAlert",
		NetworkID: testNetworkID,
		Labels:    map[string]string{"gatewayID": testGatewayID},
	})
	assert.NoError(t, err)
	action.AssertExpectations(t)

	_, err = as.Trigger(context.Background(), &schema.AlertPayload{
		TenantID:  "tenant1",
		Alertname: "testAlert",
		NetworkID: testNetworkID,
		Labels:    map[string]string{"gatewayID": "wrongGateway"},
	})
	assert.NoError(t, err)
	// Assert execute wasn't called again on trigger with wrong gateway label
	action.AssertNumberOfCalls(t, "Execute", 1)
}
