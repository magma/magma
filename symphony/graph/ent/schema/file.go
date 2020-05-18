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

// File defines the file schema.
type File struct {
	schema
}

// Fields returns file fields.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("type"),
		field.String("name").
			StructTag(`gqlgen:"fileName"`),
		field.Int("size").
			StructTag(`gqlgen:"sizeInBytes"`).
			NonNegative().
			Optional(),
		field.Time("modified_at").
			StructTag(`gqlgen:"modified"`).
			Optional(),
		field.Time("uploaded_at").
			StructTag(`gqlgen:"uploaded"`).
			Optional(),
		field.String("content_type").
			StructTag(`gqlgen:"mimeType"`),
		field.String("store_key"),
		field.String("category").
			Optional(),
	}
}

// Edges of the File
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("location", Location.Type).
			Ref("files").
			Unique(),
		edge.From("equipment", Equipment.Type).
			Ref("files").
			Unique(),
		edge.From("user", User.Type).
			Ref("profile_photo").
			Unique(),
		edge.From("work_order", WorkOrder.Type).
			Ref("files").
			Unique(),
		edge.From("checklist_item", CheckListItem.Type).
			Ref("files").
			Unique(),
	}
}

// Policy returns file policy.
func (File) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			authz.FileReadPolicyRule(),
		),
		authz.WithMutationRules(
			authz.FileWritePolicyRule(),
			authz.FileCreatePolicyRule(),
		))
}
