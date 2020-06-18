// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
)

// PermissionsPolicy defines the policy schema.
type PermissionsPolicy struct {
	schema
}

// Fields returns policy fields.
func (PermissionsPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("description").
			Optional(),
		field.Bool("is_global").
			Optional().
			Default(false),
		field.JSON("inventory_policy", &models.InventoryPolicyInput{}).
			Optional(),
		field.JSON("workforce_policy", &models.WorkforcePolicyInput{}).
			Optional(),
	}
}

// Edges returns policy edges.
func (PermissionsPolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("groups", UsersGroup.Type).
			Ref("policies"),
	}
}

// Policy returns PermissionsPolicy policies.
func (PermissionsPolicy) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.PermissionsPolicyWritePolicyRule(),
		),
	)
}
