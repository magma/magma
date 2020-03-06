// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphactions

import (
	"context"
	"strconv"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/pkg/errors"
)

// EntDataLoader is an implementation of loading rules using Ents
type EntDataLoader struct {
	Client *ent.Client
}

// QueryRules is an implementation of DataLoader.QueryRules
func (e EntDataLoader) QueryRules(ctx context.Context, triggerID core.TriggerID) ([]core.Rule, error) {
	entRules, err := e.Client.ActionsRule.Query().Where(
		actionsrule.TriggerIDEQ(string(triggerID)),
	).All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying rules")
	}

	rules := make([]core.Rule, 0, len(entRules))
	for _, rule := range entRules {
		rules = append(rules, entRuleToRule(rule))
	}
	return rules, nil
}

func entRuleToRule(rule *ent.ActionsRule) core.Rule {
	return core.Rule{
		ID:          strconv.Itoa(rule.ID),
		Name:        rule.Name,
		TriggerID:   core.TriggerID(rule.TriggerID),
		RuleActions: rule.RuleActions,
		RuleFilters: rule.RuleFilters,
	}
}
