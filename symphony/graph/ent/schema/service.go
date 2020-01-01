// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Customer holds the schema definition for the ServiceType entity.
type Customer struct {
	schema
}

// Fields of the Customer.
func (Customer) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("external_id").
			Optional().
			Nillable().
			NotEmpty().
			Unique(),
	}
}

// Edges of the Customer.
func (Customer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("services", Service.Type).
			Ref("customer"),
	}
}

// ServiceType holds the schema definition for the ServiceType entity.
type ServiceType struct {
	schema
}

// Fields of the ServiceType.
func (ServiceType) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.Bool("has_customer").Default(false),
	}
}

// Edges of the ServiceType.
func (ServiceType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("services", Service.Type).
			Ref("type"),
		edge.To("property_types", PropertyType.Type),
	}
}

// ServiceEndpoint holds the schema definition for the ServiceEndpoint entity.
type ServiceEndpoint struct {
	schema
}

// Fields of the ServiceType.
func (ServiceEndpoint) Fields() []ent.Field {
	return []ent.Field{
		field.String("role"),
	}
}

// Edges of the ServiceType.
func (ServiceEndpoint) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("port", EquipmentPort.Type).Unique(),
		edge.From("service", Service.Type).Ref("endpoints").Unique(),
	}
}

// Service holds the schema definition for the Service entity.
type Service struct {
	schema
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("external_id").
			Optional().
			Nillable().
			NotEmpty().
			Unique(),
		field.String("status"),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", ServiceType.Type).
			Unique().
			Required(),
		edge.To("upstream", Service.Type).
			From("downstream"),
		edge.To("properties", Property.Type),
		edge.To("links", Link.Type),
		edge.To("customer", Customer.Type),
		edge.To("endpoints", ServiceEndpoint.Type),
	}
}
