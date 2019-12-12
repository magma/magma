// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"context"

	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/pkg/errors"
)

// DataLoader is an interface for querying data for the executor
type DataLoader interface {
	QueryRules(context.Context, core.TriggerID) ([]core.Rule, error)
}

// BasicDataLoader is a simple implementation for querying rules.
// In the real world, this will query a database
type BasicDataLoader struct {
	Rules []core.Rule
}

// QueryRules returns all rules matches the specified TriggerID
func (b BasicDataLoader) QueryRules(ctx context.Context, triggerID core.TriggerID) (ret []core.Rule, _err error) {
	for _, rule := range b.Rules {
		if rule.TriggerID == triggerID {
			ret = append(ret, rule)
		}
	}
	return
}

// EntDataLoader is an implementation of loading rules using Ents
type EntDataLoader struct {
	client *ent.Client
}

// QueryRules is an implementation of DataLoader.QueryRules
func (e EntDataLoader) QueryRules(ctx context.Context, triggerID core.TriggerID) ([]core.Rule, error) {
	entRules, err := e.client.ActionsRule.Query().Where(
		actionsrule.TriggerIDEQ(string(triggerID)),
	).All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying rules")
	}

	rules := []core.Rule{}
	for _, rule := range entRules {
		rules = append(rules, entRuleToRule(rule))
	}
	return rules, nil
}

func entRuleToRule(rule *ent.ActionsRule) core.Rule {
	return core.Rule{
		ID:          rule.ID,
		Name:        rule.Name,
		TriggerID:   core.TriggerID(rule.TriggerID),
		RuleActions: rule.RuleActions,
		RuleFilters: rule.RuleFilters,
	}
}
