// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
)

// CheckListCategory defines the CheckListCategory type schema.
type CheckListCategory struct {
	schema
}

// Fields returns CheckListCategory type fields.
func (CheckListCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("description").
			Optional(),
	}
}

// Edges returns CheckListCategory type edges.
func (CheckListCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("check_list_items", CheckListItem.Type),
	}
}

// CheckListItem defines the CheckListItem type schema.
type CheckListItemDefinition struct {
	schema
}

// Fields returns CheckListItem type fields.
func (CheckListItemDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("type"),
		field.Int("index").
			Optional(),
		field.String("enum_values").
			StructTag(`gqlgen:"enumValues"`).
			Nillable().
			Optional(),
		field.String("help_text").
			StructTag(`gqlgen:"helpText"`).
			Nillable().
			Optional(),
	}
}

// Edges returns CheckListItem type edges.
func (CheckListItemDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("work_order_type", WorkOrderType.Type).
			Ref("check_list_definitions").
			Unique(),
	}
}

// Indexes returns CheckListItem type indexes.
func (CheckListItemDefinition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title").
			Edges("work_order_type").
			Unique(),
	}
}

// CheckListItem defines the CheckListItem schema.
type CheckListItem struct {
	ent.Schema
}

// Fields returns CheckListItem fields.
func (CheckListItem) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("type"),
		field.Int("index").
			Optional(),
		field.Bool("checked").
			Optional(),
		field.String("string_val").
			StructTag(`gqlgen:"stringValue"`).
			Optional(),
		field.String("enum_values").
			StructTag(`gqlgen:"enumValues"`).
			Optional(),
		field.String("enum_selection_mode").
			StructTag(`gqlgen:"enumSelectionMode"`).
			Optional(),
		field.String("selected_enum_values").
			StructTag(`gqlgen:"selectedEnumValues"`).
			Optional(),
		field.String("help_text").
			StructTag(`gqlgen:"helpText"`).
			Nillable().
			Optional(),
	}
}

// Edges returns CheckListItem edges.
func (CheckListItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
		edge.From("work_order", WorkOrder.Type).
			Unique().
			Ref("check_list_items"),
	}
}

// Indexes returns CheckListItem type indexes.
func (CheckListItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title").
			Edges("work_order").
			Unique(),
	}
}
