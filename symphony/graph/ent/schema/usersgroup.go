// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/authz"
)

// UsersGroup defines the users group schema.
type UsersGroup struct {
	schema
}

// Fields returns UsersGroup fields.
func (UsersGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("description").
			Optional(),
		field.Enum("status").
			Values("ACTIVE", "DEACTIVATED").
			Default("ACTIVE"),
	}
}

// Edges returns UsersGroup edges.
func (UsersGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("members", User.Type),
		edge.To("policies", PermissionsPolicy.Type),
	}
}

// Policy returns UserGroup policies.
func (UsersGroup) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.AllowAdminRule(),
		),
	)
}
