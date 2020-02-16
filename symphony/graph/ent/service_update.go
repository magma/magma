// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceUpdate is the builder for updating Service entities.
type ServiceUpdate struct {
	config

	update_time       *time.Time
	name              *string
	external_id       *string
	clearexternal_id  bool
	status            *string
	_type             map[string]struct{}
	downstream        map[string]struct{}
	upstream          map[string]struct{}
	properties        map[string]struct{}
	links             map[string]struct{}
	customer          map[string]struct{}
	endpoints         map[string]struct{}
	clearedType       bool
	removedDownstream map[string]struct{}
	removedUpstream   map[string]struct{}
	removedProperties map[string]struct{}
	removedLinks      map[string]struct{}
	removedCustomer   map[string]struct{}
	removedEndpoints  map[string]struct{}
	predicates        []predicate.Service
}

// Where adds a new predicate for the builder.
func (su *ServiceUpdate) Where(ps ...predicate.Service) *ServiceUpdate {
	su.predicates = append(su.predicates, ps...)
	return su
}

// SetName sets the name field.
func (su *ServiceUpdate) SetName(s string) *ServiceUpdate {
	su.name = &s
	return su
}

// SetExternalID sets the external_id field.
func (su *ServiceUpdate) SetExternalID(s string) *ServiceUpdate {
	su.external_id = &s
	return su
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (su *ServiceUpdate) SetNillableExternalID(s *string) *ServiceUpdate {
	if s != nil {
		su.SetExternalID(*s)
	}
	return su
}

// ClearExternalID clears the value of external_id.
func (su *ServiceUpdate) ClearExternalID() *ServiceUpdate {
	su.external_id = nil
	su.clearexternal_id = true
	return su
}

// SetStatus sets the status field.
func (su *ServiceUpdate) SetStatus(s string) *ServiceUpdate {
	su.status = &s
	return su
}

// SetTypeID sets the type edge to ServiceType by id.
func (su *ServiceUpdate) SetTypeID(id string) *ServiceUpdate {
	if su._type == nil {
		su._type = make(map[string]struct{})
	}
	su._type[id] = struct{}{}
	return su
}

// SetType sets the type edge to ServiceType.
func (su *ServiceUpdate) SetType(s *ServiceType) *ServiceUpdate {
	return su.SetTypeID(s.ID)
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (su *ServiceUpdate) AddDownstreamIDs(ids ...string) *ServiceUpdate {
	if su.downstream == nil {
		su.downstream = make(map[string]struct{})
	}
	for i := range ids {
		su.downstream[ids[i]] = struct{}{}
	}
	return su
}

// AddDownstream adds the downstream edges to Service.
func (su *ServiceUpdate) AddDownstream(s ...*Service) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddDownstreamIDs(ids...)
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (su *ServiceUpdate) AddUpstreamIDs(ids ...string) *ServiceUpdate {
	if su.upstream == nil {
		su.upstream = make(map[string]struct{})
	}
	for i := range ids {
		su.upstream[ids[i]] = struct{}{}
	}
	return su
}

// AddUpstream adds the upstream edges to Service.
func (su *ServiceUpdate) AddUpstream(s ...*Service) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddUpstreamIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (su *ServiceUpdate) AddPropertyIDs(ids ...string) *ServiceUpdate {
	if su.properties == nil {
		su.properties = make(map[string]struct{})
	}
	for i := range ids {
		su.properties[ids[i]] = struct{}{}
	}
	return su
}

// AddProperties adds the properties edges to Property.
func (su *ServiceUpdate) AddProperties(p ...*Property) *ServiceUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.AddPropertyIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (su *ServiceUpdate) AddLinkIDs(ids ...string) *ServiceUpdate {
	if su.links == nil {
		su.links = make(map[string]struct{})
	}
	for i := range ids {
		su.links[ids[i]] = struct{}{}
	}
	return su
}

// AddLinks adds the links edges to Link.
func (su *ServiceUpdate) AddLinks(l ...*Link) *ServiceUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return su.AddLinkIDs(ids...)
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (su *ServiceUpdate) AddCustomerIDs(ids ...string) *ServiceUpdate {
	if su.customer == nil {
		su.customer = make(map[string]struct{})
	}
	for i := range ids {
		su.customer[ids[i]] = struct{}{}
	}
	return su
}

// AddCustomer adds the customer edges to Customer.
func (su *ServiceUpdate) AddCustomer(c ...*Customer) *ServiceUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return su.AddCustomerIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (su *ServiceUpdate) AddEndpointIDs(ids ...string) *ServiceUpdate {
	if su.endpoints == nil {
		su.endpoints = make(map[string]struct{})
	}
	for i := range ids {
		su.endpoints[ids[i]] = struct{}{}
	}
	return su
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (su *ServiceUpdate) AddEndpoints(s ...*ServiceEndpoint) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddEndpointIDs(ids...)
}

// ClearType clears the type edge to ServiceType.
func (su *ServiceUpdate) ClearType() *ServiceUpdate {
	su.clearedType = true
	return su
}

// RemoveDownstreamIDs removes the downstream edge to Service by ids.
func (su *ServiceUpdate) RemoveDownstreamIDs(ids ...string) *ServiceUpdate {
	if su.removedDownstream == nil {
		su.removedDownstream = make(map[string]struct{})
	}
	for i := range ids {
		su.removedDownstream[ids[i]] = struct{}{}
	}
	return su
}

// RemoveDownstream removes downstream edges to Service.
func (su *ServiceUpdate) RemoveDownstream(s ...*Service) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveDownstreamIDs(ids...)
}

// RemoveUpstreamIDs removes the upstream edge to Service by ids.
func (su *ServiceUpdate) RemoveUpstreamIDs(ids ...string) *ServiceUpdate {
	if su.removedUpstream == nil {
		su.removedUpstream = make(map[string]struct{})
	}
	for i := range ids {
		su.removedUpstream[ids[i]] = struct{}{}
	}
	return su
}

// RemoveUpstream removes upstream edges to Service.
func (su *ServiceUpdate) RemoveUpstream(s ...*Service) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveUpstreamIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (su *ServiceUpdate) RemovePropertyIDs(ids ...string) *ServiceUpdate {
	if su.removedProperties == nil {
		su.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		su.removedProperties[ids[i]] = struct{}{}
	}
	return su
}

// RemoveProperties removes properties edges to Property.
func (su *ServiceUpdate) RemoveProperties(p ...*Property) *ServiceUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.RemovePropertyIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (su *ServiceUpdate) RemoveLinkIDs(ids ...string) *ServiceUpdate {
	if su.removedLinks == nil {
		su.removedLinks = make(map[string]struct{})
	}
	for i := range ids {
		su.removedLinks[ids[i]] = struct{}{}
	}
	return su
}

// RemoveLinks removes links edges to Link.
func (su *ServiceUpdate) RemoveLinks(l ...*Link) *ServiceUpdate {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return su.RemoveLinkIDs(ids...)
}

// RemoveCustomerIDs removes the customer edge to Customer by ids.
func (su *ServiceUpdate) RemoveCustomerIDs(ids ...string) *ServiceUpdate {
	if su.removedCustomer == nil {
		su.removedCustomer = make(map[string]struct{})
	}
	for i := range ids {
		su.removedCustomer[ids[i]] = struct{}{}
	}
	return su
}

// RemoveCustomer removes customer edges to Customer.
func (su *ServiceUpdate) RemoveCustomer(c ...*Customer) *ServiceUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return su.RemoveCustomerIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (su *ServiceUpdate) RemoveEndpointIDs(ids ...string) *ServiceUpdate {
	if su.removedEndpoints == nil {
		su.removedEndpoints = make(map[string]struct{})
	}
	for i := range ids {
		su.removedEndpoints[ids[i]] = struct{}{}
	}
	return su
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (su *ServiceUpdate) RemoveEndpoints(s ...*ServiceEndpoint) *ServiceUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (su *ServiceUpdate) Save(ctx context.Context) (int, error) {
	if su.update_time == nil {
		v := service.UpdateDefaultUpdateTime()
		su.update_time = &v
	}
	if su.name != nil {
		if err := service.NameValidator(*su.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if su.external_id != nil {
		if err := service.ExternalIDValidator(*su.external_id); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	if len(su._type) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if su.clearedType && su._type == nil {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}
	return su.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (su *ServiceUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *ServiceUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *ServiceUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

func (su *ServiceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   service.Table,
			Columns: service.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: service.FieldID,
			},
		},
	}
	if ps := su.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := su.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: service.FieldUpdateTime,
		})
	}
	if value := su.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: service.FieldName,
		})
	}
	if value := su.external_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: service.FieldExternalID,
		})
	}
	if su.clearexternal_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: service.FieldExternalID,
		})
	}
	if value := su.status; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: service.FieldStatus,
		})
	}
	if su.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   service.TypeTable,
			Columns: []string{service.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: servicetype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   service.TypeTable,
			Columns: []string{service.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: servicetype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.removedDownstream; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   service.DownstreamTable,
			Columns: service.DownstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.downstream; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   service.DownstreamTable,
			Columns: service.DownstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.removedUpstream; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.UpstreamTable,
			Columns: service.UpstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.upstream; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.UpstreamTable,
			Columns: service.UpstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.PropertiesTable,
			Columns: []string{service.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.PropertiesTable,
			Columns: []string{service.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.removedLinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.LinksTable,
			Columns: service.LinksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.links; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.LinksTable,
			Columns: service.LinksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.removedCustomer; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.CustomerTable,
			Columns: service.CustomerPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: customer.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.customer; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.CustomerTable,
			Columns: service.CustomerPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: customer.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.removedEndpoints; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.EndpointsTable,
			Columns: []string{service.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.endpoints; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.EndpointsTable,
			Columns: []string{service.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{service.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ServiceUpdateOne is the builder for updating a single Service entity.
type ServiceUpdateOne struct {
	config
	id string

	update_time       *time.Time
	name              *string
	external_id       *string
	clearexternal_id  bool
	status            *string
	_type             map[string]struct{}
	downstream        map[string]struct{}
	upstream          map[string]struct{}
	properties        map[string]struct{}
	links             map[string]struct{}
	customer          map[string]struct{}
	endpoints         map[string]struct{}
	clearedType       bool
	removedDownstream map[string]struct{}
	removedUpstream   map[string]struct{}
	removedProperties map[string]struct{}
	removedLinks      map[string]struct{}
	removedCustomer   map[string]struct{}
	removedEndpoints  map[string]struct{}
}

// SetName sets the name field.
func (suo *ServiceUpdateOne) SetName(s string) *ServiceUpdateOne {
	suo.name = &s
	return suo
}

// SetExternalID sets the external_id field.
func (suo *ServiceUpdateOne) SetExternalID(s string) *ServiceUpdateOne {
	suo.external_id = &s
	return suo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableExternalID(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetExternalID(*s)
	}
	return suo
}

// ClearExternalID clears the value of external_id.
func (suo *ServiceUpdateOne) ClearExternalID() *ServiceUpdateOne {
	suo.external_id = nil
	suo.clearexternal_id = true
	return suo
}

// SetStatus sets the status field.
func (suo *ServiceUpdateOne) SetStatus(s string) *ServiceUpdateOne {
	suo.status = &s
	return suo
}

// SetTypeID sets the type edge to ServiceType by id.
func (suo *ServiceUpdateOne) SetTypeID(id string) *ServiceUpdateOne {
	if suo._type == nil {
		suo._type = make(map[string]struct{})
	}
	suo._type[id] = struct{}{}
	return suo
}

// SetType sets the type edge to ServiceType.
func (suo *ServiceUpdateOne) SetType(s *ServiceType) *ServiceUpdateOne {
	return suo.SetTypeID(s.ID)
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (suo *ServiceUpdateOne) AddDownstreamIDs(ids ...string) *ServiceUpdateOne {
	if suo.downstream == nil {
		suo.downstream = make(map[string]struct{})
	}
	for i := range ids {
		suo.downstream[ids[i]] = struct{}{}
	}
	return suo
}

// AddDownstream adds the downstream edges to Service.
func (suo *ServiceUpdateOne) AddDownstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddDownstreamIDs(ids...)
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (suo *ServiceUpdateOne) AddUpstreamIDs(ids ...string) *ServiceUpdateOne {
	if suo.upstream == nil {
		suo.upstream = make(map[string]struct{})
	}
	for i := range ids {
		suo.upstream[ids[i]] = struct{}{}
	}
	return suo
}

// AddUpstream adds the upstream edges to Service.
func (suo *ServiceUpdateOne) AddUpstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddUpstreamIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (suo *ServiceUpdateOne) AddPropertyIDs(ids ...string) *ServiceUpdateOne {
	if suo.properties == nil {
		suo.properties = make(map[string]struct{})
	}
	for i := range ids {
		suo.properties[ids[i]] = struct{}{}
	}
	return suo
}

// AddProperties adds the properties edges to Property.
func (suo *ServiceUpdateOne) AddProperties(p ...*Property) *ServiceUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.AddPropertyIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (suo *ServiceUpdateOne) AddLinkIDs(ids ...string) *ServiceUpdateOne {
	if suo.links == nil {
		suo.links = make(map[string]struct{})
	}
	for i := range ids {
		suo.links[ids[i]] = struct{}{}
	}
	return suo
}

// AddLinks adds the links edges to Link.
func (suo *ServiceUpdateOne) AddLinks(l ...*Link) *ServiceUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return suo.AddLinkIDs(ids...)
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (suo *ServiceUpdateOne) AddCustomerIDs(ids ...string) *ServiceUpdateOne {
	if suo.customer == nil {
		suo.customer = make(map[string]struct{})
	}
	for i := range ids {
		suo.customer[ids[i]] = struct{}{}
	}
	return suo
}

// AddCustomer adds the customer edges to Customer.
func (suo *ServiceUpdateOne) AddCustomer(c ...*Customer) *ServiceUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return suo.AddCustomerIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (suo *ServiceUpdateOne) AddEndpointIDs(ids ...string) *ServiceUpdateOne {
	if suo.endpoints == nil {
		suo.endpoints = make(map[string]struct{})
	}
	for i := range ids {
		suo.endpoints[ids[i]] = struct{}{}
	}
	return suo
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (suo *ServiceUpdateOne) AddEndpoints(s ...*ServiceEndpoint) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddEndpointIDs(ids...)
}

// ClearType clears the type edge to ServiceType.
func (suo *ServiceUpdateOne) ClearType() *ServiceUpdateOne {
	suo.clearedType = true
	return suo
}

// RemoveDownstreamIDs removes the downstream edge to Service by ids.
func (suo *ServiceUpdateOne) RemoveDownstreamIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedDownstream == nil {
		suo.removedDownstream = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedDownstream[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveDownstream removes downstream edges to Service.
func (suo *ServiceUpdateOne) RemoveDownstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveDownstreamIDs(ids...)
}

// RemoveUpstreamIDs removes the upstream edge to Service by ids.
func (suo *ServiceUpdateOne) RemoveUpstreamIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedUpstream == nil {
		suo.removedUpstream = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedUpstream[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveUpstream removes upstream edges to Service.
func (suo *ServiceUpdateOne) RemoveUpstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveUpstreamIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (suo *ServiceUpdateOne) RemovePropertyIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedProperties == nil {
		suo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedProperties[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveProperties removes properties edges to Property.
func (suo *ServiceUpdateOne) RemoveProperties(p ...*Property) *ServiceUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.RemovePropertyIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (suo *ServiceUpdateOne) RemoveLinkIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedLinks == nil {
		suo.removedLinks = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedLinks[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveLinks removes links edges to Link.
func (suo *ServiceUpdateOne) RemoveLinks(l ...*Link) *ServiceUpdateOne {
	ids := make([]string, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return suo.RemoveLinkIDs(ids...)
}

// RemoveCustomerIDs removes the customer edge to Customer by ids.
func (suo *ServiceUpdateOne) RemoveCustomerIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedCustomer == nil {
		suo.removedCustomer = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedCustomer[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveCustomer removes customer edges to Customer.
func (suo *ServiceUpdateOne) RemoveCustomer(c ...*Customer) *ServiceUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return suo.RemoveCustomerIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (suo *ServiceUpdateOne) RemoveEndpointIDs(ids ...string) *ServiceUpdateOne {
	if suo.removedEndpoints == nil {
		suo.removedEndpoints = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedEndpoints[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (suo *ServiceUpdateOne) RemoveEndpoints(s ...*ServiceEndpoint) *ServiceUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (suo *ServiceUpdateOne) Save(ctx context.Context) (*Service, error) {
	if suo.update_time == nil {
		v := service.UpdateDefaultUpdateTime()
		suo.update_time = &v
	}
	if suo.name != nil {
		if err := service.NameValidator(*suo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if suo.external_id != nil {
		if err := service.ExternalIDValidator(*suo.external_id); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	if len(suo._type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"type\"")
	}
	if suo.clearedType && suo._type == nil {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}
	return suo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *ServiceUpdateOne) SaveX(ctx context.Context) *Service {
	s, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return s
}

// Exec executes the query on the entity.
func (suo *ServiceUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *ServiceUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (suo *ServiceUpdateOne) sqlSave(ctx context.Context) (s *Service, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   service.Table,
			Columns: service.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  suo.id,
				Type:   field.TypeString,
				Column: service.FieldID,
			},
		},
	}
	if value := suo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: service.FieldUpdateTime,
		})
	}
	if value := suo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: service.FieldName,
		})
	}
	if value := suo.external_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: service.FieldExternalID,
		})
	}
	if suo.clearexternal_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: service.FieldExternalID,
		})
	}
	if value := suo.status; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: service.FieldStatus,
		})
	}
	if suo.clearedType {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   service.TypeTable,
			Columns: []string{service.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: servicetype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo._type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   service.TypeTable,
			Columns: []string{service.TypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: servicetype.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.removedDownstream; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   service.DownstreamTable,
			Columns: service.DownstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.downstream; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   service.DownstreamTable,
			Columns: service.DownstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.removedUpstream; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.UpstreamTable,
			Columns: service.UpstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.upstream; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.UpstreamTable,
			Columns: service.UpstreamPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.removedProperties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.PropertiesTable,
			Columns: []string{service.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.properties; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.PropertiesTable,
			Columns: []string{service.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: property.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.removedLinks; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.LinksTable,
			Columns: service.LinksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.links; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.LinksTable,
			Columns: service.LinksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: link.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.removedCustomer; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.CustomerTable,
			Columns: service.CustomerPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: customer.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.customer; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   service.CustomerTable,
			Columns: service.CustomerPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: customer.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.removedEndpoints; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.EndpointsTable,
			Columns: []string{service.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.endpoints; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.EndpointsTable,
			Columns: []string{service.EndpointsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: serviceendpoint.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	s = &Service{config: suo.config}
	_spec.Assign = s.assignValues
	_spec.ScanValues = s.scanValues()
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{service.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return s, nil
}
