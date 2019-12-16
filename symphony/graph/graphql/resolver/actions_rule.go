// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/cloud/actions"
	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

type actionsRuleResolver struct{}

func (actionsRuleResolver) Trigger(ctx context.Context, rule *ent.ActionsRule) (*models.ActionsTrigger, error) {
	ac := actions.FromContext(ctx)

	modelTriggerID := models.TriggerID(rule.TriggerID)
	actionsTriggerID := core.TriggerID(rule.TriggerID)

	if !modelTriggerID.IsValid() {
		return nil, errors.Errorf("triggerID %s not in models", rule.TriggerID)
	}

	trigger, err := ac.TriggerForID(actionsTriggerID)
	if err != nil {
		return nil, errors.Errorf("triggerID %s not a registered action", actionsTriggerID)
	}

	return &models.ActionsTrigger{
		TriggerID:   modelTriggerID,
		Description: trigger.Description(),
	}, nil
}

func (actionsRuleResolver) RuleActions(ctx context.Context, rule *ent.ActionsRule) ([]*core.ActionsRuleAction, error) {
	return rule.RuleActions, nil
}

func (actionsRuleResolver) RuleFilters(ctx context.Context, rule *ent.ActionsRule) ([]*core.ActionsRuleFilter, error) {
	return rule.RuleFilters, nil
}

func (actionsRuleResolver) TriggerID(ctx context.Context, rule *ent.ActionsRule) (models.TriggerID, error) {
	value := models.TriggerID(rule.TriggerID)
	if !value.IsValid() {
		return "", errors.Errorf("rule %s does not have a valid triggerID %s", rule.ID, rule.TriggerID)
	}
	return value, nil
}
