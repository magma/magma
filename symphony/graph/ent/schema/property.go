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

// PropertyType defines the property type schema.
type PropertyType struct {
	schema
}

// Fields returns property type fields.
func (PropertyType) Fields() []ent.Field {
	return []ent.Field{
		field.String("type"),
		field.String("name"),
		field.Int("index").
			Optional(),
		field.String("category").
			Optional(),
		field.Int("int_val").
			StructTag(`json:"intValue" gqlgen:"intValue"`).
			Optional(),
		field.Bool("bool_val").
			StructTag(`json:"booleanValue" gqlgen:"booleanValue"`).
			Optional(),
		field.Float("float_val").
			StructTag(`json:"floatValue" gqlgen:"floatValue"`).
			Optional(),
		field.Float("latitude_val").
			StructTag(`json:"latitudeValue" gqlgen:"latitudeValue"`).
			Optional(),
		field.Float("longitude_val").
			StructTag(`json:"longitudeValue" gqlgen:"longitudeValue"`).
			Optional(),
		field.Text("string_val").
			StructTag(`json:"stringValue" gqlgen:"stringValue"`).
			Optional(),
		field.Float("range_from_val").
			StructTag(`json:"rangeFromValue" gqlgen:"rangeFromValue"`).
			Optional(),
		field.Float("range_to_val").
			StructTag(`json:"rangeToValue" gqlgen:"rangeToValue"`).
			Optional(),
		field.Bool("is_instance_property").
			StructTag(`gqlgen:"isInstanceProperty"`).
			Default(true),
		field.Bool("editable").
			StructTag(`gqlgen:"isEditable"`).
			Default(true),
		field.Bool("mandatory").
			StructTag(`gqlgen:"isMandatory"`).
			Default(false),
		field.Bool("deleted").
			StructTag(`gqlgen:"isDeleted"`).
			Default(false),
	}
}

// Edges returns property type edges.
func (PropertyType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("properties", Property.Type).
			Ref("type"),
		edge.From("location_type", LocationType.Type).
			Ref("property_types").
			Unique(),
		edge.From("equipment_port_type", EquipmentPortType.Type).
			Ref("property_types").
			Unique(),
		edge.From("link_equipment_port_type", EquipmentPortType.Type).
			Ref("link_property_types").
			Unique(),
		edge.From("equipment_type", EquipmentType.Type).
			Ref("property_types").
			Unique(),
		edge.From("service_type", ServiceType.Type).
			Ref("property_types").
			Unique(),
		edge.From("work_order_type", WorkOrderType.Type).
			Ref("property_types").
			Unique(),
		edge.From("project_type", ProjectType.Type).
			Ref("properties").
			Unique(),
	}
}

// Indexes returns property type indexes.
func (PropertyType) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Edges("location_type").
			Unique(),
		index.Fields("name").
			Edges("equipment_port_type").
			Unique(),
		index.Fields("name").
			Edges("equipment_type").
			Unique(),
		index.Fields("name").
			Edges("link_equipment_port_type").
			Unique(),
		index.Fields("name").
			Edges("work_order_type").
			Unique(),
	}
}

// Property defines the property schema.
type Property struct {
	schema
}

// Fields returns property fields.
func (Property) Fields() []ent.Field {
	return []ent.Field{
		field.Int("int_val").
			StructTag(`json:"intValue" gqlgen:"intValue"`).
			Optional(),
		field.Bool("bool_val").
			StructTag(`json:"booleanValue" gqlgen:"booleanValue"`).
			Optional(),
		field.Float("float_val").
			StructTag(`json:"floatValue" gqlgen:"floatValue"`).
			Optional(),
		field.Float("latitude_val").
			StructTag(`json:"latitudeValue" gqlgen:"latitudeValue"`).
			Optional(),
		field.Float("longitude_val").
			StructTag(`json:"longitudeValue" gqlgen:"longitudeValue"`).
			Optional(),
		field.Float("range_from_val").
			StructTag(`json:"rangeFromValue" gqlgen:"rangeFromValue"`).
			Optional(),
		field.Float("range_to_val").
			StructTag(`json:"rangeToValue" gqlgen:"rangeToValue"`).
			Optional(),
		field.String("string_val").
			StructTag(`json:"stringValue" gqlgen:"stringValue"`).
			Optional(),
	}
}

// Edges returns property edges.
func (Property) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", PropertyType.Type).
			Unique().
			Required().
			StructTag(`gqlgen:"propertyType"`),
		edge.From("location", Location.Type).
			Unique().
			Ref("properties").
			StructTag(`gqlgen:"locationValue"`),
		edge.From("equipment", Equipment.Type).
			Unique().
			Ref("properties").
			StructTag(`gqlgen:"equipmentValue"`),
		edge.From("service", Service.Type).
			Unique().
			Ref("properties").
			StructTag(`gqlgen:"serviceValue"`),
		edge.From("equipment_port", EquipmentPort.Type).
			Unique().
			Ref("properties"),
		edge.From("link", Link.Type).
			Unique().
			Ref("properties"),
		edge.From("work_order", WorkOrder.Type).
			Unique().
			Ref("properties"),
		edge.From("project", Project.Type).
			Ref("properties").
			Unique(),
		edge.To("equipment_value", Equipment.Type).
			Unique(),
		edge.To("location_value", Location.Type).
			Unique(),
		edge.To("service_value", Service.Type).
			Unique(),
	}
}
