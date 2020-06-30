// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
	"github.com/facebookincubator/ent/schema/mixin"
	"github.com/facebookincubator/symphony/pkg/authz"
)

// WorkOrderTemplateMixin defines the work order template mixin schema.
type WorkOrderTemplateMixin struct {
	mixin.Schema
}

// Fields returns work order template mixin fields.
func (WorkOrderTemplateMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Text("description").
			Optional(),
	}
}

// Edges returns work order template mixin edges.
func (WorkOrderTemplateMixin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("property_types", PropertyType.Type),
		edge.To("check_list_category_definitions", CheckListCategoryDefinition.Type),
	}
}

// WorkOrderType defines the work order type schema.
type WorkOrderType struct {
	schema
}

// Mixin returns work order type mixins.
func (WorkOrderType) Mixin() []ent.Mixin {
	return []ent.Mixin{
		WorkOrderTemplateMixin{},
	}
}

// Indexes returns work order type indexes.
func (WorkOrderType) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Edges returns work order type edges.
func (WorkOrderType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("work_orders", WorkOrder.Type).
			Ref("type"),
		edge.From("definitions", WorkOrderDefinition.Type).
			Ref("type"),
	}
}

// Policy returns work order type policy.
func (WorkOrderType) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.WorkOrderTypeWritePolicyRule(),
		),
	)
}

// WorkOrderTemplate defines the work order template schema.
type WorkOrderTemplate struct {
	schema
}

// Mixin returns work order template mixins.
func (WorkOrderTemplate) Mixin() []ent.Mixin {
	return []ent.Mixin{
		WorkOrderTemplateMixin{},
	}
}

// Edges returns work order template edges.
func (WorkOrderTemplate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", WorkOrderType.Type).
			Unique(),
	}
}

// Policy returns work order template policy.
func (WorkOrderTemplate) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.WorkOrderTemplateWritePolicyRule(),
		),
	)
}

// WorkOrder defines the work order schema.
type WorkOrder struct {
	schema
}

// Fields returns work order fields.
func (WorkOrder) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("status").
			Default("PLANNED"),
		field.String("priority").
			Default("NONE"),
		field.Text("description").
			Optional(),
		field.Time("install_date").
			Optional(),
		field.Time("creation_date"),
		field.Int("index").
			Optional(),
		field.Time("close_date").
			Optional(),
	}
}

// Edges returns work order edges.
func (WorkOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", WorkOrderType.Type).
			Unique(),
		edge.To("template", WorkOrderTemplate.Type).
			Unique(),
		edge.From("equipment", Equipment.Type).
			Ref("work_order"),
		edge.From("links", Link.Type).
			Ref("work_order"),
		edge.To("files", File.Type),
		edge.To("hyperlinks", Hyperlink.Type),
		edge.To("location", Location.Type).
			Unique(),
		edge.To("comments", Comment.Type),
		edge.To("activities", Activity.Type),
		edge.To("properties", Property.Type),
		edge.To("check_list_categories", CheckListCategory.Type),
		edge.From("project", Project.Type).
			Ref("work_orders").
			Unique(),
		edge.To("owner", User.Type).
			Required().
			Unique(),
		edge.To("assignee", User.Type).
			Unique(),
	}
}

// Indexes returns work order indexes.
func (WorkOrder) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("creation_date"),
		index.Fields("close_date"),
	}
}

// Policy returns work order policy.
func (WorkOrder) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			authz.WorkOrderReadPolicyRule(),
		),
		authz.WithMutationRules(
			authz.WorkOrderWritePolicyRule(),
			authz.AllowWorkOrderOwnerOrAssigneeWrite(),
		),
	)
}
