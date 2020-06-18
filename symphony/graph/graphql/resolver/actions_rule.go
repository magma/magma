// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/pkg/errors"
)

type actionsRuleResolver struct{}

func (actionsRuleResolver) Trigger(ctx context.Context, rule *ent.ActionsRule) (*models.ActionsTrigger, error) {
	ac := actions.FromContext(ctx)

	actionsTriggerID := core.TriggerID(rule.TriggerID)

	trigger, err := ac.TriggerForID(actionsTriggerID)
	if err != nil {
		return nil, errors.Errorf("triggerID %s not a registered action", actionsTriggerID)
	}

	return &models.ActionsTrigger{
		TriggerID:   core.TriggerID(rule.TriggerID),
		Description: trigger.Description(),
	}, nil
}

func (actionsRuleResolver) RuleActions(ctx context.Context, rule *ent.ActionsRule) ([]*core.ActionsRuleAction, error) {
	return rule.RuleActions, nil
}

func (actionsRuleResolver) RuleFilters(ctx context.Context, rule *ent.ActionsRule) ([]*core.ActionsRuleFilter, error) {
	return rule.RuleFilters, nil
}

func (actionsRuleResolver) TriggerID(ctx context.Context, rule *ent.ActionsRule) (core.TriggerID, error) {
	return core.TriggerID(rule.TriggerID), nil
}
