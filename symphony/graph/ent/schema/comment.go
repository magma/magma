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

// Comment defines the location type schema.
type Comment struct {
	schema
}

// Fields returns Comment fields.
func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.String("text"),
	}
}

// Edges returns Comment edges.
func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("author", User.Type).
			Required().
			Unique(),
		edge.From("work_order", WorkOrder.Type).
			Ref("comments").
			Unique(),
		edge.From("project", Project.Type).
			Ref("comments").
			Unique(),
	}
}

// Policy returns Comment policy.
func (Comment) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			authz.CommentReadPolicyRule(),
		),
		authz.WithMutationRules(
			authz.CommentWritePolicyRule(),
			authz.CommentCreatePolicyRule(),
		),
	)
}
