// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// User defines the user schema.
type User struct {
	schema
}

// Fields returns user fields.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("auth_id").
			NotEmpty().
			Immutable().
			Unique(),
		field.String("first_name").
			NotEmpty().
			Optional(),
		field.String("last_name").
			NotEmpty().
			Optional(),
		field.String("email").
			NotEmpty().
			Optional(),
		field.Enum("status").
			Values("ACTIVE", "DEACTIVATED").
			Default("ACTIVE"),
		field.Enum("role").
			Values("USER", "ADMIN", "OWNER").
			Default("USER"),
	}
}

// Edges returns user edges.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("profile_photo", File.Type).
			Unique(),
		edge.From("groups", UsersGroup.Type).
			Ref("members"),
	}
}

// Policy returns user policy.
func (User) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			privacy.DenyMutationOperationRule(
				ent.OpDelete|ent.OpDeleteOne,
			),
			authz.UserWritePolicyRule(),
		),
	)
}

// Hooks of the User.
func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		viewer.UpdateCurrentUser(),
	}
}
