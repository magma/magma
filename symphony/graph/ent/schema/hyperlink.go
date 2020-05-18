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

// Hyperlink defines the hyperlink schema.
type Hyperlink struct {
	schema
}

// Fields returns hyperlink fields.
func (Hyperlink) Fields() []ent.Field {
	return []ent.Field{
		field.String("url"),
		field.String("name").
			StructTag(`gqlgen:"displayName"`).
			Optional(),
		field.String("category").
			Optional(),
	}
}

// Edges returns hyperlink edges.
func (Hyperlink) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("work_order", WorkOrder.Type).
			Ref("hyperlinks").
			Unique(),
	}
}

// Policy returns hyperlink policy.
func (Hyperlink) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			authz.HyperlinkReadPolicyRule(),
		),
	)
}
