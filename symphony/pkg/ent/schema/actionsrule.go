// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

// ActionsRule defines the location type schema.
type ActionsRule struct {
	schema
}

// Fields returns action rule fields.
func (ActionsRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("triggerID"),
		field.JSON("ruleFilters", []*core.ActionsRuleFilter{}),
		field.JSON("ruleActions", []*core.ActionsRuleAction{}),
	}
}

// Policy returns ActionsRule policy.
func (ActionsRule) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			privacy.AlwaysAllowRule(),
		),
	)
}
