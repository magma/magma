// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// EquipmentPortType defines the equipment port definition schema.
type EquipmentPortType struct {
	schema
}

// Fields returns equipment type fields.
func (EquipmentPortType) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
	}
}

// Edges returns equipment type edges.
func (EquipmentPortType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("property_types", PropertyType.Type),
		edge.To("link_property_types", PropertyType.Type),
		edge.From("port_definitions", EquipmentPortDefinition.Type).
			Ref("equipment_port_type"),
	}
}

// EquipmentPortDefinition defines the equipment port definition schema.
type EquipmentPortDefinition struct {
	schema
}

// Fields returns equipment port definition fields.
func (EquipmentPortDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("index").
			Optional(),
		field.String("bandwidth").
			Optional(),
		field.String("visibility_label").
			StructTag(`gqlgen:"visibleLabel"`).
			Optional(),
	}
}

// Edges returns equipment port definition edges.
func (EquipmentPortDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("equipment_port_type", EquipmentPortType.Type).
			Unique(),
		edge.From("ports", EquipmentPort.Type).
			Ref("definition"),
		edge.From("equipment_type", EquipmentType.Type).
			Ref("port_definitions").
			Unique(),
	}
}

// EquipmentPort defines the equipment port schema.
type EquipmentPort struct {
	schema
}

// Edges returns equipment port Edges.
func (EquipmentPort) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("definition", EquipmentPortDefinition.Type).
			Required().
			Unique(),
		edge.From("parent", Equipment.Type).
			Ref("ports").
			Unique(),
		edge.To("link", Link.Type).
			Unique(),
		edge.To("properties", Property.Type),
		edge.From("endpoints", ServiceEndpoint.Type).
			Ref("port"),
	}
}

// EquipmentPositionDefinition defines the equipment position definition schema.
type EquipmentPositionDefinition struct {
	schema
}

// Fields returns equipment position definition fields.
func (EquipmentPositionDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("index").
			Optional(),
		field.String("visibility_label").
			StructTag(`gqlgen:"visibleLabel"`).
			Optional(),
	}
}

// Edges returns equipment position definition edges.
func (EquipmentPositionDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("positions", EquipmentPosition.Type).
			Ref("definition"),
		edge.From("equipment_type", EquipmentType.Type).
			Ref("position_definitions").
			Unique(),
	}
}

// EquipmentPosition defines the equipment position schema.
type EquipmentPosition struct {
	schema
}

// Edges returns equipment position Edges.
func (EquipmentPosition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("definition", EquipmentPositionDefinition.Type).
			Required().
			Unique(),
		edge.From("parent", Equipment.Type).
			Ref("positions").
			Unique(),
		edge.To("attachment", Equipment.Type).
			Unique(),
	}
}

// EquipmentCategory defines the equipment category schema.
type EquipmentCategory struct {
	schema
}

// Fields returns equipment category fields.
func (EquipmentCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
	}
}

// Edges returns equipment category edges.
func (EquipmentCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("types", EquipmentType.Type).
			Ref("category"),
	}
}

// EquipmentType defines the equipment type schema.
type EquipmentType struct {
	schema
}

// Fields returns equipment type fields.
func (EquipmentType) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
	}
}

// Edges returns equipment type edges.
func (EquipmentType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("port_definitions", EquipmentPortDefinition.Type),
		edge.To("position_definitions", EquipmentPositionDefinition.Type),
		edge.To("property_types", PropertyType.Type),
		edge.From("equipment", Equipment.Type).
			Ref("type"),
		edge.To("category", EquipmentCategory.Type).
			Unique(),
	}
}

// Equipment defines the equipment schema.
type Equipment struct {
	schema
}

// Fields returns equipment fields.
func (Equipment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("future_state").
			Optional(),
		field.String("device_id").
			Optional(),
		field.String("external_id").
			Optional(),
	}
}

// Edges returns equipment edges.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", EquipmentType.Type).
			Unique().
			Required(),
		edge.From("location", Location.Type).
			Ref("equipment").
			Unique(),
		edge.From("parent_position", EquipmentPosition.Type).
			Ref("attachment").
			Unique(),
		edge.To("positions", EquipmentPosition.Type),
		edge.To("ports", EquipmentPort.Type),
		edge.To("work_order", WorkOrder.Type).
			Unique(),
		edge.To("properties", Property.Type),
		edge.From("service", Service.Type).
			Ref("termination_points"),
		edge.To("files", File.Type),
	}
}
