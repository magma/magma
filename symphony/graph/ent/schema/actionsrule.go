// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/cloud/actions/core"
)

// ActionsRule defines the location type schema.
type ActionsRule struct {
	schema
}

// ActionsRuleAction are actions that will be executed for the rule
type ActionsRuleAction struct {
	ActionID core.ActionID `json:"actionID"`
	Data     string        `json:"data"`
}

// ActionsRuleFilter are filtered that are applied to a rule
type ActionsRuleFilter struct {
	FilterID   string `json:"filterID"`
	OperatorID string `json:"operatorID"`
	Data       string `json:"data"`
}

// Fields returns action rule fields.
func (ActionsRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("triggerID"),
		field.JSON("ruleFilters", []*ActionsRuleFilter{}),
		field.JSON("ruleActions", []*ActionsRuleAction{}),
	}
}
