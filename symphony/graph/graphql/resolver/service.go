// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

type serviceTypeResolver struct{}

func (serviceTypeResolver) PropertyTypes(ctx context.Context, obj *ent.ServiceType) ([]*ent.PropertyType, error) {
	return obj.QueryPropertyTypes().All(ctx)
}

func (serviceTypeResolver) Services(ctx context.Context, obj *ent.ServiceType) ([]*ent.Service, error) {
	return obj.QueryServices().All(ctx)
}

func (serviceTypeResolver) NumberOfServices(ctx context.Context, obj *ent.ServiceType) (int, error) {
	return obj.QueryServices().Count(ctx)
}

type serviceResolver struct{}

func (serviceResolver) Customer(ctx context.Context, obj *ent.Service) (*ent.Customer, error) {
	customer, err := obj.QueryCustomer().First(ctx)
	if err != nil {
		return nil, ent.MaskNotFound(err)
	}
	return customer, nil
}

func (serviceResolver) ServiceType(ctx context.Context, obj *ent.Service) (*ent.ServiceType, error) {
	return obj.QueryType().Only(ctx)
}

func (serviceResolver) Status(ctx context.Context, obj *ent.Service) (models.ServiceStatus, error) {
	return models.ServiceStatus(obj.Status), nil
}

func (serviceResolver) Upstream(ctx context.Context, obj *ent.Service) ([]*ent.Service, error) {
	return obj.QueryUpstream().All(ctx)
}

func (serviceResolver) Downstream(ctx context.Context, obj *ent.Service) ([]*ent.Service, error) {
	return obj.QueryDownstream().All(ctx)
}

func (serviceResolver) Properties(ctx context.Context, obj *ent.Service) ([]*ent.Property, error) {
	return obj.QueryProperties().All(ctx)
}

func (serviceResolver) Links(ctx context.Context, obj *ent.Service) ([]*ent.Link, error) {
	return obj.QueryLinks().All(ctx)
}

func (serviceResolver) Endpoints(ctx context.Context, obj *ent.Service) ([]*ent.ServiceEndpoint, error) {
	return obj.QueryEndpoints().All(ctx)
}

func (serviceResolver) rootNode(ctx context.Context, eq *ent.Equipment) *ent.Equipment {
	parent := eq
	for parent != nil {
		p, err := parent.QueryParentPosition().QueryParent().Only(ctx)
		if err != nil {
			break
		}

		parent = p
	}

	return parent
}

func (r serviceResolver) Topology(ctx context.Context, obj *ent.Service) (*models.NetworkTopology, error) {
	eqs, err := obj.QueryLinks().QueryPorts().QueryParent().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying links equipments")
	}

	var nodes []ent.Noder
	eqsMap := make(map[string]*ent.Equipment)
	for _, eq := range eqs {
		node := r.rootNode(ctx, eq)
		eqsMap[node.ID] = eq
		nodes = append(nodes, node)
	}

	eps, err := obj.QueryEndpoints().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying termination points")
	}

	for _, ep := range eps {
		equipment, err := ep.QueryPort().QueryParent().Only(ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return nil, errors.Wrap(err, "querying equipment of endpoint")
			}
		} else {
			node := r.rootNode(ctx, equipment)
			if _, ok := eqsMap[node.ID]; !ok {
				nodes = append(nodes, node)
			}
		}
	}

	lnks, err := obj.QueryLinks().All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querying service links")
	}

	var links []*models.TopologyLink

	for _, lnk := range lnks {
		leqs, err := lnk.
			QueryPorts().
			QueryParent().
			All(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "querying link equipments")
		}
		node0 := r.rootNode(ctx, leqs[0])
		node1 := r.rootNode(ctx, leqs[1])
		links = append(links, &models.TopologyLink{Type: models.TopologyLinkTypePhysical, Source: node0, Target: node1})
	}

	return &models.NetworkTopology{Nodes: nodes, Links: links}, nil
}

type serviceEndpointResolver struct{}

func (r serviceEndpointResolver) Port(ctx context.Context, obj *ent.ServiceEndpoint) (*ent.EquipmentPort, error) {
	return obj.QueryPort().Only(ctx)
}

func (r serviceEndpointResolver) Role(ctx context.Context, obj *ent.ServiceEndpoint) (models.ServiceEndpointRole, error) {
	return models.ServiceEndpointRole(obj.Role), nil
}

func (serviceEndpointResolver) Service(ctx context.Context, obj *ent.ServiceEndpoint) (*ent.Service, error) {
	return obj.QueryService().Only(ctx)
}

func (r mutationResolver) RemoveService(ctx context.Context, id string) (string, error) {
	client := r.ClientFrom(ctx)

	if _, err := client.ServiceEndpoint.Delete().
		Where(serviceendpoint.HasServiceWith(service.ID(id))).
		Exec(ctx); err != nil {
		return "", errors.Wrapf(err, "deleting service endpoints: id=%q", id)
	}

	if _, err := client.Property.Delete().
		Where(property.HasServiceWith(service.ID(id))).
		Exec(ctx); err != nil {
		return "", errors.Wrapf(err, "deleting service properties: id=%q", id)
	}
	if err := client.Service.DeleteOneID(id).Exec(ctx); err != nil {
		return "", errors.Wrapf(err, "deleting service: id=%q", id)
	}
	return id, nil
}

func (r mutationResolver) AddServiceEndpoint(ctx context.Context, input models.AddServiceEndpointInput) (*ent.Service, error) {
	client := r.ClientFrom(ctx)
	s, err := client.Service.Get(ctx, input.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "querying service: id=%q", input.ID)
	}

	if _, err := client.ServiceEndpoint.
		Create().
		SetRole(input.Role.String()).
		SetServiceID(input.ID).
		SetPortID(input.PortID).Save(ctx); err != nil {
		return nil, errors.Wrapf(err, "Creating service endpoint: service id=%q", input.ID)
	}

	return s, nil
}

func (r mutationResolver) RemoveServiceEndpoint(ctx context.Context, serviceEndpointID string) (*ent.Service, error) {
	client := r.ClientFrom(ctx)

	s, err := client.Service.Query().Where(service.HasEndpointsWith(serviceendpoint.ID(serviceEndpointID))).Only(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "query service")
	}

	if err := client.ServiceEndpoint.DeleteOneID(serviceEndpointID).Exec(ctx); err != nil {
		return nil, errors.Wrap(err, "query endpoint")
	}

	return s, nil
}
