// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
	"github.com/facebookincubator/symphony/pkg/authz"
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
		edge.To("projects", Project.Type).
			StructTag(`gqlgen:"projects"`),
		edge.To("properties", PropertyType.Type).
			StructTag(`gqlgen:"properties"`),
		edge.To("work_orders", WorkOrderDefinition.Type).
			StructTag(`gqlgen:"workOrders"`),
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
			Required().
			StructTag(`gqlgen:"type"`),
		edge.To("location", Location.Type).
			Unique().
			StructTag(`gqlgen:"location"`),
		edge.To("comments", Comment.Type).
			StructTag(`gqlgen:"comments"`),
		edge.To("work_orders", WorkOrder.Type).
			StructTag(`gqlgen:"workOrders"`),
		edge.To("properties", Property.Type).
			StructTag(`gqlgen:"properties"`),
		edge.To("creator", User.Type).
			Unique().
			StructTag(`gqlgen:"createdBy"`),
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

// Policy returns project policy.
func (Project) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			authz.ProjectReadPolicyRule(),
		),
		authz.WithMutationRules(
			authz.ProjectWritePolicyRule(),
			authz.AllowProjectCreatorWrite(),
		),
	)
}

// WorkOrderDefinition defines the work order definition schema.
type WorkOrderDefinition struct {
	schema
}

// Fields returns work order definition fields.
func (WorkOrderDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.Int("index").
			Optional(),
	}
}

// Edges returns work order definition edges.
func (WorkOrderDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", WorkOrderType.Type).
			Unique().
			StructTag(`gqlgen:"type"`),
		edge.From("project_type", ProjectType.Type).
			Ref("work_orders").
			Unique(),
	}
}

// Policy returns work order definition policy.
func (WorkOrderDefinition) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.WorkOrderDefinitionWritePolicyRule(),
		),
	)
}
