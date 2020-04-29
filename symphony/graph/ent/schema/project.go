// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
	"github.com/facebookincubator/symphony/graph/authz"
)

// ProjectType defines the project type schema.
type ProjectType struct {
	schema
}

// Fields returns project fields.
func (ProjectType) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.Text("description").
			Optional().
			Nillable(),
	}
}

// Edges return project type edges.
func (ProjectType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
		edge.To("properties", PropertyType.Type),
		edge.To("work_orders", WorkOrderDefinition.Type),
	}
}

// Policy returns project type policy.
func (ProjectType) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.ProjectTypeWritePolicyRule(),
		),
	)
}

// Project defines the project schema.
type Project struct {
	schema
}

// Fields returns project fields.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.Text("description").
			Optional().
			Nillable(),
	}
}

// Edges returns project edges.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("type", ProjectType.Type).
			Ref("projects").
			Unique().
			Required(),
		edge.To("location", Location.Type).
			Unique(),
		edge.To("comments", Comment.Type),
		edge.To("work_orders", WorkOrder.Type),
		edge.To("properties", Property.Type),
		edge.To("creator", User.Type).
			Unique(),
	}
}

// Indexes return project indexes.
func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Edges("type").
			Unique(),
	}
}

// EquipmentPortDefinition defines the equipment port definition schema.
type WorkOrderDefinition struct {
	schema
}

// Fields returns equipment port definition fields.
func (WorkOrderDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.Int("index").
			Optional(),
	}
}

// Edges returns equipment port definition edges.
func (WorkOrderDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", WorkOrderType.Type).
			Unique(),
		edge.From("project_type", ProjectType.Type).
			Ref("work_orders").
			Unique(),
	}
}
