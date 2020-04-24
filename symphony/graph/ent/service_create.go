// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceCreate is the builder for creating a Service entity.
type ServiceCreate struct {
	config
	mutation *ServiceMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (sc *ServiceCreate) SetCreateTime(t time.Time) *ServiceCreate {
	sc.mutation.SetCreateTime(t)
	return sc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (sc *ServiceCreate) SetNillableCreateTime(t *time.Time) *ServiceCreate {
	if t != nil {
		sc.SetCreateTime(*t)
	}
	return sc
}

// SetUpdateTime sets the update_time field.
func (sc *ServiceCreate) SetUpdateTime(t time.Time) *ServiceCreate {
	sc.mutation.SetUpdateTime(t)
	return sc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (sc *ServiceCreate) SetNillableUpdateTime(t *time.Time) *ServiceCreate {
	if t != nil {
		sc.SetUpdateTime(*t)
	}
	return sc
}

// SetName sets the name field.
func (sc *ServiceCreate) SetName(s string) *ServiceCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetExternalID sets the external_id field.
func (sc *ServiceCreate) SetExternalID(s string) *ServiceCreate {
	sc.mutation.SetExternalID(s)
	return sc
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (sc *ServiceCreate) SetNillableExternalID(s *string) *ServiceCreate {
	if s != nil {
		sc.SetExternalID(*s)
	}
	return sc
}

// SetStatus sets the status field.
func (sc *ServiceCreate) SetStatus(s string) *ServiceCreate {
	sc.mutation.SetStatus(s)
	return sc
}

// SetDiscoveryMethod sets the discovery_method field.
func (sc *ServiceCreate) SetDiscoveryMethod(sm service.DiscoveryMethod) *ServiceCreate {
	sc.mutation.SetDiscoveryMethod(sm)
	return sc
}

// SetNillableDiscoveryMethod sets the discovery_method field if the given value is not nil.
func (sc *ServiceCreate) SetNillableDiscoveryMethod(sm *service.DiscoveryMethod) *ServiceCreate {
	if sm != nil {
		sc.SetDiscoveryMethod(*sm)
	}
	return sc
}

// SetTypeID sets the type edge to ServiceType by id.
func (sc *ServiceCreate) SetTypeID(id int) *ServiceCreate {
	sc.mutation.SetTypeID(id)
	return sc
}

// SetType sets the type edge to ServiceType.
func (sc *ServiceCreate) SetType(s *ServiceType) *ServiceCreate {
	return sc.SetTypeID(s.ID)
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (sc *ServiceCreate) AddDownstreamIDs(ids ...int) *ServiceCreate {
	sc.mutation.AddDownstreamIDs(ids...)
	return sc
}

// AddDownstream adds the downstream edges to Service.
func (sc *ServiceCreate) AddDownstream(s ...*Service) *ServiceCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddDownstreamIDs(ids...)
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (sc *ServiceCreate) AddUpstreamIDs(ids ...int) *ServiceCreate {
	sc.mutation.AddUpstreamIDs(ids...)
	return sc
}

// AddUpstream adds the upstream edges to Service.
func (sc *ServiceCreate) AddUpstream(s ...*Service) *ServiceCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddUpstreamIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (sc *ServiceCreate) AddPropertyIDs(ids ...int) *ServiceCreate {
	sc.mutation.AddPropertyIDs(ids...)
	return sc
}

// AddProperties adds the properties edges to Property.
func (sc *ServiceCreate) AddProperties(p ...*Property) *ServiceCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return sc.AddPropertyIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (sc *ServiceCreate) AddLinkIDs(ids ...int) *ServiceCreate {
	sc.mutation.AddLinkIDs(ids...)
	return sc
}

// AddLinks adds the links edges to Link.
func (sc *ServiceCreate) AddLinks(l ...*Link) *ServiceCreate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return sc.AddLinkIDs(ids...)
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (sc *ServiceCreate) AddCustomerIDs(ids ...int) *ServiceCreate {
	sc.mutation.AddCustomerIDs(ids...)
	return sc
}

// AddCustomer adds the customer edges to Customer.
func (sc *ServiceCreate) AddCustomer(c ...*Customer) *ServiceCreate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return sc.AddCustomerIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (sc *ServiceCreate) AddEndpointIDs(ids ...int) *ServiceCreate {
	sc.mutation.AddEndpointIDs(ids...)
	return sc
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (sc *ServiceCreate) AddEndpoints(s ...*ServiceEndpoint) *ServiceCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddEndpointIDs(ids...)
}

// Save creates the Service in the database.
func (sc *ServiceCreate) Save(ctx context.Context) (*Service, error) {
	if _, ok := sc.mutation.CreateTime(); !ok {
		v := service.DefaultCreateTime()
		sc.mutation.SetCreateTime(v)
	}
	if _, ok := sc.mutation.UpdateTime(); !ok {
		v := service.DefaultUpdateTime()
		sc.mutation.SetUpdateTime(v)
	}
	if _, ok := sc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := sc.mutation.Name(); ok {
		if err := service.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := sc.mutation.ExternalID(); ok {
		if err := service.ExternalIDValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	if _, ok := sc.mutation.Status(); !ok {
		return nil, errors.New("ent: missing required field \"status\"")
	}
	if v, ok := sc.mutation.DiscoveryMethod(); ok {
		if err := service.DiscoveryMethodValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"discovery_method\": %v", err)
		}
	}
	if _, ok := sc.mutation.TypeID(); !ok {
		return nil, errors.New("ent: missing required edge \"type\"")
	}
	var (
		err  error
		node *Service
	)
	if len(sc.hooks) == 0 {
		node, err = sc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sc.mutation = mutation
			node, err = sc.sqlSave(ctx)
			return node, err
		})
		for i := len(sc.hooks) - 1; i >= 0; i-- {
			mut = sc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sc *ServiceCreate) SaveX(ctx context.Context) *Service {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sc *ServiceCreate) sqlSave(ctx context.Context) (*Service, error) {
	var (
		s     = &Service{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: service.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: service.FieldID,
			},
		}
	)
	if value, ok := sc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: service.FieldCreateTime,
		})
		s.CreateTime = value
	}
	if value, ok := sc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: service.FieldUpdateTime,
		})
		s.UpdateTime = value
	}
	if value, ok := sc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldName,
		})
		s.Name = value
	}
	if value, ok := sc.mutation.ExternalID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldExternalID,
		})
		s.ExternalID = &value
	}
	if value, ok := sc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldStatus,
		})
		s.Status = value
	}
	if value, ok := sc.mutation.DiscoveryMethod(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: service.FieldDiscoveryMethod,
		})
		s.DiscoveryMethod = value
	}
	if nodes := sc.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   service.TypeTable,
			Columns: []string{service.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: servicetype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.DownstreamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   service.DownstreamTable,
			Columns: service.DownstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.UpstreamIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.UpstreamTable,
			Columns: service.UpstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.PropertiesTable,
			Columns: []string{service.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.LinksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.LinksTable,
			Columns: service.LinksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: link.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.CustomerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.CustomerTable,
			Columns: service.CustomerPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: customer.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.EndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.EndpointsTable,
			Columns: []string{service.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	s.ID = int(id)
	return s, nil
}
