// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/cloud/actions"
	"github.com/facebookincubator/symphony/cloud/actions/action/mockaction"
	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/facebookincubator/symphony/cloud/actions/executor"
	"github.com/facebookincubator/symphony/cloud/actions/trigger/mocktrigger"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func actionsContext(t *testing.T) (generated.ResolverRoot, context.Context) {

	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(core.TriggerID("trigger1"))

	action1 := mockaction.New()
	action1.On("ID").Return(core.ActionID("action1"))

	r, ctx := resolverctx(t)
	registry := actions.MainRegistry()
	registry.MustRegisterTrigger(trigger1)
	registry.MustRegisterAction(action1)

	exc := &executor.Executor{
		Context:  ctx,
		Registry: registry,
		DataLoader: executor.BasicDataLoader{
			Rules: []core.Rule{},
		},
	}
	ctx = actions.NewContext(ctx, exc)
	return r, ctx
}

func TestAddActionsRuleInvalidTrigger(t *testing.T) {
	r, ctx := actionsContext(t)

	input := models.AddActionsRuleInput{
		Name:        "testInput",
		TriggerID:   "triggerInvalid",
		RuleActions: []*models.ActionsRuleActionInput{},
		RuleFilters: []*models.ActionsRuleFilterInput{},
	}
	_, err := r.Mutation().AddActionsRule(ctx, input)
	require.Error(t, err)
}

func TestAddActionsRule(t *testing.T) {
	r, ctx := actionsContext(t)

	input := models.AddActionsRuleInput{
		Name:      "testInput",
		TriggerID: "trigger1",
		RuleActions: []*models.ActionsRuleActionInput{
			{
				ActionID: "action1",
				Data:     "testdata",
			},
		},
		RuleFilters: []*models.ActionsRuleFilterInput{
			{
				FilterID:   "filter1",
				OperatorID: "eq",
				Data:       "testdata",
			},
		},
	}
	rule, err := r.Mutation().AddActionsRule(ctx, input)
	require.NoError(t, err)

	assert.Equal(t, rule.Name, "testInput")
	assert.Equal(t, rule.TriggerID, "trigger1")

	assert.Equal(t, len(rule.RuleActions), 1)
	assert.Equal(t, rule.RuleActions[0].ActionID, core.ActionID("action1"))
	assert.Equal(t, rule.RuleActions[0].Data, "testdata")

	assert.Equal(t, len(rule.RuleFilters), 1)
	assert.Equal(t, rule.RuleFilters[0].FilterID, "filter1")
	assert.Equal(t, rule.RuleFilters[0].OperatorID, "eq")
	assert.Equal(t, rule.RuleFilters[0].Data, "testdata")
}
