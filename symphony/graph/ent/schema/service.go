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
	"github.com/facebookincubator/symphony/graph/ent/privacy"
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

// Policy returns Customer policy.
func (Customer) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			privacy.AlwaysAllowRule(),
		),
	)
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
		field.Bool("is_deleted").Default(false),
		field.Enum("discovery_method").
			Comment("how will service of this type be discovered? (null means manual adding and not discovery)").
			Values("INVENTORY").
			Optional(),
	}
}

// Edges of the ServiceType.
func (ServiceType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("services", Service.Type).
			Ref("type"),
		edge.To("property_types", PropertyType.Type),
		edge.To("endpoint_definitions", ServiceEndpointDefinition.Type),
	}
}

// Policy returns service type policy.
func (ServiceType) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.ServiceTypeWritePolicyRule(),
		),
	)
}

// ServiceEndpoint holds the schema definition for the ServiceEndpoint entity.
type ServiceEndpoint struct {
	schema
}

// Edges of the ServiceEndpoint.
func (ServiceEndpoint) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("port", EquipmentPort.Type).Unique(),
		edge.To("equipment", Equipment.Type).Unique().Required(),
		edge.From("service", Service.Type).Ref("endpoints").Unique().Required(),
		edge.From("definition", ServiceEndpointDefinition.Type).Ref("endpoints").Unique(),
	}
}

// Policy returns ServiceEndPoint policy.
func (ServiceEndpoint) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.ServiceEndpointWritePolicyRule(),
		),
	)
}

// ServiceEndpointDefinition holds the schema definition for the ServiceEndpointDefinition entity.
type ServiceEndpointDefinition struct {
	schema
}

// Fields of the ServiceEndpointDefinition.
func (ServiceEndpointDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("role").
			Optional(),
		field.String("name").NotEmpty(),
		field.Int("index"),
	}
}

// Edges of the ServiceEndpointDefinition.
func (ServiceEndpointDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("endpoints", ServiceEndpoint.Type),
		edge.From("service_type", ServiceType.Type).
			Unique().
			Ref("endpoint_definitions"),
		edge.From("equipment_type", EquipmentType.Type).
			Unique().
			Ref("service_endpoint_definitions"),
	}
}

// Indexes returns ServiceEndpointDefinition indexes.
func (ServiceEndpointDefinition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("index").
			Edges("service_type").
			Unique(),
		index.Fields("name").
			Edges("service_type").
			Unique(),
	}
}

// Policy returns ServiceEndpointDefinition policy.
func (ServiceEndpointDefinition) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.ServiceEndpointDefinitionWritePolicyRule(),
		),
	)
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

// Policy returns service policy.
func (Service) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithMutationRules(
			authz.ServiceWritePolicyRule(),
		),
	)
}
