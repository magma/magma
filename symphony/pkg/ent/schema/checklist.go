// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/authz"
)

// CheckListCategory defines the CheckListCategoryDefinition type schema.
type CheckListCategoryDefinition struct {
	schema
}

// Fields returns CheckListCategoryDefinition type fields.
func (CheckListCategoryDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			NotEmpty(),
		field.String("description").
			Optional(),
	}
}

// Edges returns CheckListCategoryDefinition type edges.
func (CheckListCategoryDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("check_list_item_definitions", CheckListItemDefinition.Type),
		edge.From("work_order_type", WorkOrderType.Type).
			Ref("check_list_category_definitions").
			Unique().
			Required(),
	}
}

// Policy returns equipment port definition policy.
func (CheckListCategoryDefinition) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.CheckListCategoryDefinitionWritePolicyRule(),
		),
	)
}

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
		edge.From("work_order", WorkOrder.Type).Ref("check_list_categories").
			Unique().
			Required(),
	}
}

// Policy returns checklist item policy.
func (CheckListCategory) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			authz.CheckListCategoryReadPolicyRule(),
		),
		authz.WithMutationRules(
			authz.CheckListCategoryWritePolicyRule(),
			authz.CheckListCategoryCreatePolicyRule(),
		),
	)
}

// CheckListItem defines the CheckListItemDefinition type schema.
type CheckListItemDefinition struct {
	schema
}

// Fields returns CheckListItemDefinition type fields.
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
		field.Enum("enum_selection_mode_value").
			Values("single", "multiple").
			Optional(),
		field.String("help_text").
			StructTag(`gqlgen:"helpText"`).
			Nillable().
			Optional(),
	}
}

// Edges returns CheckListItemDefinition type edges.
func (CheckListItemDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("check_list_category_definition", CheckListCategoryDefinition.Type).
			Ref("check_list_item_definitions").
			Unique().
			Required(),
	}
}

// Policy returns equipment port definition policy.
func (CheckListItemDefinition) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.CheckListItemDefinitionWritePolicyRule(),
		),
	)
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
		field.Enum("enum_selection_mode_value").
			Values("single", "multiple").
			Optional(),
		field.String("selected_enum_values").
			StructTag(`gqlgen:"selectedEnumValues"`).
			Optional(),
		field.Enum("yes_no_val").
			Values("YES", "NO").
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
		edge.From("wifi_scan", SurveyWiFiScan.Type).
			Ref("checklist_item"),
		edge.From("cell_scan", SurveyCellScan.Type).
			Ref("checklist_item"),
		edge.From("check_list_category", CheckListCategory.Type).
			Ref("check_list_items").
			Unique().
			Required(),
	}
}

// Policy returns equipment port definition policy.
func (CheckListItem) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			authz.CheckListItemReadPolicyRule(),
		),
		authz.WithMutationRules(
			authz.CheckListItemWritePolicyRule(),
			authz.CheckListItemCreatePolicyRule(),
		),
	)
}
