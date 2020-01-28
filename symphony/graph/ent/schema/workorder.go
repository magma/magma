// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// WorkOrderType defines the work order type schema.
type WorkOrderType struct {
	schema
}

// Fields returns work order type fields.
func (WorkOrderType) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.Text("description").
			Optional(),
	}
}

// Edges returns work order type edges.
func (WorkOrderType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("work_orders", WorkOrder.Type).
			Ref("type"),
		edge.To("property_types", PropertyType.Type),
		edge.From("definitions", WorkOrderDefinition.Type).
			Ref("type"),
		edge.To("check_list_definitions", CheckListItemDefinition.Type),
	}
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
		field.String("owner_name"),
		field.Time("install_date").
			Optional(),
		field.Time("creation_date"),
		field.String("assignee").
			Optional(),
		field.Int("index").
			Optional(),
	}
}

// Edges returns work order edges.
func (WorkOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", WorkOrderType.Type).
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
		edge.To("properties", Property.Type),
		edge.To("check_list_items", CheckListItem.Type),
		edge.To("technician", Technician.Type).
			Unique(),
		edge.From("project", Project.Type).
			Ref("work_orders").
			Unique(),
	}
}
