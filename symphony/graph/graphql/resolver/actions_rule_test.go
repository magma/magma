// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/action/mockaction"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/actions/trigger/mocktrigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func actionsContext(t *testing.T) (*TestResolver, context.Context) {

	trigger1 := mocktrigger.New()
	trigger1.On("ID").Return(core.TriggerID("trigger1"))
	trigger2 := mocktrigger.New()
	trigger2.On("ID").Return(core.TriggerID("trigger2"))

	action1 := mockaction.New()
	action1.On("ID").Return(core.ActionID("action1"))

	r, ctx := resolverctx(t)
	registry := executor.NewRegistry()
	registry.MustRegisterTrigger(trigger1)
	registry.MustRegisterTrigger(trigger2)
	registry.MustRegisterAction(action1)

	exc := &executor.Executor{
		Registry: registry,
		DataLoader: executor.BasicDataLoader{
			Rules: []core.Rule{},
		},
	}
	ctx = actions.NewContext(ctx, exc)
	return r.(*TestResolver), ctx
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

func TestQueryActionsRules(t *testing.T) {
	r, ctx := actionsContext(t)
	actions := []*core.ActionsRuleAction{
		{
			ActionID: "action1",
			Data:     "testdata",
		},
	}

	filters := []*core.ActionsRuleFilter{
		{
			FilterID:   "filter1",
			OperatorID: "eq",
			Data:       "testdata",
		},
	}

	rule, err := r.client.
		ActionsRule.
		Create().
		SetName("testInput").
		SetTriggerID("trigger1").
		SetRuleActions(actions).
		SetRuleFilters(filters).
		Save(ctx)
	require.NoError(t, err)

	rules, err := r.Query().ActionsRules(ctx)
	require.NoError(t, err)
	assert.Equal(t, len(rules.Results), 1)
	assert.Equal(t, rules.Results[0].ID, rule.ID)
	assert.Equal(t, rules.Results[0].Name, rule.Name)
	assert.Equal(t, rules.Results[0].TriggerID, rule.TriggerID)
	assert.Equal(t, rules.Results[0].RuleActions[0].ActionID, rule.RuleActions[0].ActionID)
	assert.Equal(t, rules.Results[0].RuleFilters[0].FilterID, rule.RuleFilters[0].FilterID)
}

func TestEditActionsRule(t *testing.T) {
	r, ctx := actionsContext(t)

	originalRule, err := r.client.
		ActionsRule.Create().
		SetName("testInput").
		SetTriggerID("trigger1").
		SetRuleActions([]*core.ActionsRuleAction{}).
		SetRuleFilters([]*core.ActionsRuleFilter{}).
		Save(ctx)
	assert.NoError(t, err)

	editedRule, err := r.Mutation().EditActionsRule(ctx, originalRule.ID, models.AddActionsRuleInput{
		Name:        "testInput",
		TriggerID:   "trigger2",
		RuleActions: []*models.ActionsRuleActionInput{},
		RuleFilters: []*models.ActionsRuleFilterInput{},
	})
	require.NoError(t, err)

	assert.Equal(t, editedRule.ID, originalRule.ID)
	assert.Equal(t, editedRule.Name, "testInput")
	assert.Equal(t, editedRule.TriggerID, "trigger2")
}

func TestRemoveActionsRule(t *testing.T) {
	r, ctx := actionsContext(t)

	originalRule, err := r.client.
		ActionsRule.Create().
		SetName("testInput").
		SetTriggerID("trigger1").
		SetRuleActions([]*core.ActionsRuleAction{}).
		SetRuleFilters([]*core.ActionsRuleFilter{}).
		Save(ctx)
	assert.NoError(t, err)

	success, err := r.Mutation().RemoveActionsRule(ctx, originalRule.ID)
	require.NoError(t, err)

	assert.True(t, success)
}
