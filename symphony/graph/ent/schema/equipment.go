// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"errors"
	"strings"

	"github.com/facebookincubator/symphony/graph/authz"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
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
		edge.To("property_types", PropertyType.Type).
			StructTag(`gqlgen:"propertyTypes"`),
		edge.To("link_property_types", PropertyType.Type).
			StructTag(`gqlgen:"linkPropertyTypes"`),
		edge.From("port_definitions", EquipmentPortDefinition.Type).
			Ref("equipment_port_type").
			StructTag(`gqlgen:"numberOfPortDefinitions"`),
	}
}

// Policy returns equipment port type policy.
func (EquipmentPortType) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.EquipmentPortTypeWritePolicyRule(),
		),
	)
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
			Unique().
			StructTag(`gqlgen:"portType"`),
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
			Unique().
			StructTag(`gqlgen:"definition"`),
		edge.From("parent", Equipment.Type).
			Ref("ports").
			Unique().
			StructTag(`gqlgen:"parentEquipment"`),
		edge.To("link", Link.Type).
			Unique().
			StructTag(`gqlgen:"link"`),
		edge.To("properties", Property.Type).
			StructTag(`gqlgen:"properties"`),
		edge.From("endpoints", ServiceEndpoint.Type).
			Ref("port").
			StructTag(`gqlgen:"serviceEndpoints"`),
	}
}

// Indexes returns equipment port indexes.
func (EquipmentPort) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("definition", "parent").
			Unique(),
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
			Unique().
			StructTag(`gqlgen:"definition"`),
		edge.From("parent", Equipment.Type).
			Ref("positions").
			Unique().
			StructTag(`gqlgen:"parentEquipment"`),
		edge.To("attachment", Equipment.Type).
			Unique().
			StructTag(`gqlgen:"attachedEquipment"`),
	}
}

// Indexes returns equipment position indexes.
func (EquipmentPosition) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("definition", "parent").
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
		edge.To("port_definitions", EquipmentPortDefinition.Type).
			StructTag(`gqlgen:"portDefinitions"`),
		edge.To("position_definitions", EquipmentPositionDefinition.Type).
			StructTag(`gqlgen:"positionDefinitions"`),
		edge.To("property_types", PropertyType.Type).
			StructTag(`gqlgen:"propertyTypes"`),
		edge.From("equipment", Equipment.Type).
			Ref("type").
			StructTag(`gqlgen:"equipments"`),
		edge.To("category", EquipmentCategory.Type).
			Unique().
			StructTag(`gqlgen:"category"`),
		edge.To("service_endpoint_definitions", ServiceEndpointDefinition.Type).
			StructTag(`gqlgen:"serviceEndpointDefinitions"`),
	}
}

// Policy returns equipment type policy.
func (EquipmentType) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.EquipmentTypeWritePolicyRule(),
		),
	)
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
			Optional().
			Validate(func(s string) error {
				if !strings.ContainsRune(s, '.') {
					return errors.New("invalid device id")
				}
				return nil
			}),
		field.String("external_id").
			Unique().
			Optional(),
	}
}

// Edges returns equipment edges.
func (Equipment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", EquipmentType.Type).
			Unique().
			Required().
			StructTag(`gqlgen:"equipmentType"`),
		edge.From("location", Location.Type).
			Ref("equipment").
			Unique().
			StructTag(`gqlgen:"parentLocation"`),
		edge.From("parent_position", EquipmentPosition.Type).
			Ref("attachment").
			Unique().
			StructTag(`gqlgen:"parentPosition"`),
		edge.To("positions", EquipmentPosition.Type).
			StructTag(`gqlgen:"positions"`),
		edge.To("ports", EquipmentPort.Type).
			StructTag(`gqlgen:"ports"`),
		edge.To("work_order", WorkOrder.Type).
			Unique().
			StructTag(`gqlgen:"workOrder"`),
		edge.To("properties", Property.Type).
			StructTag(`gqlgen:"properties"`),
		edge.To("files", File.Type).
			StructTag(`gqlgen:"files"`),
		edge.To("hyperlinks", Hyperlink.Type).
			StructTag(`gqlgen:"hyperlinks"`),
		edge.From("endpoints", ServiceEndpoint.Type).
			Ref("equipment"),
	}
}

// Policy returns equipment policy.
func (Equipment) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.EquipmentWritePolicyRule(),
		),
	)
}
