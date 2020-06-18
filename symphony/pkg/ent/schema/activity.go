// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

// Activity defines the location type schema.
type Activity struct {
	schema
}

// Fields returns Activity fields.
func (Activity) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("changed_field").
			Values("STATUS", "PRIORITY", "ASSIGNEE", "CREATION_DATE", "OWNER"),
		field.Bool("is_create").
			Default(false),
		field.String("old_value").
			Optional().
			Comment("raw value of the previous state (enum, entID ..)"),
		field.String("new_value").
			Optional().
			Comment("raw value of the next state (enum, entID ..)"),
	}
}

// Edges returns Activity edges.
func (Activity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("author", User.Type).
			Unique(),
		edge.From("work_order", WorkOrder.Type).
			Ref("activities").
			Unique(),
	}
}

// Policy returns Activity policy.
//TODO(T67933416): Replace the always read\write with the correct permissions
func (Activity) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			privacy.AlwaysAllowRule(),
			// authz.ActivityReadPolicyRule(),
		),
		authz.WithMutationRules(
			privacy.AlwaysAllowRule(),
			// authz.ActivityWritePolicyRule(),
			// authz.ActivityCreatePolicyRule(),
		),
	)
}
