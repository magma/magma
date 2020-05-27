// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

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
	hooks      []Hook
	mutation   *ServiceMutation
	predicates []predicate.Service
}

// Where adds a new predicate for the builder.
func (su *ServiceUpdate) Where(ps ...predicate.Service) *ServiceUpdate {
	su.predicates = append(su.predicates, ps...)
	return su
}

// SetName sets the name field.
func (su *ServiceUpdate) SetName(s string) *ServiceUpdate {
	su.mutation.SetName(s)
	return su
}

// SetExternalID sets the external_id field.
func (su *ServiceUpdate) SetExternalID(s string) *ServiceUpdate {
	su.mutation.SetExternalID(s)
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
	su.mutation.ClearExternalID()
	return su
}

// SetStatus sets the status field.
func (su *ServiceUpdate) SetStatus(s string) *ServiceUpdate {
	su.mutation.SetStatus(s)
	return su
}

// SetTypeID sets the type edge to ServiceType by id.
func (su *ServiceUpdate) SetTypeID(id int) *ServiceUpdate {
	su.mutation.SetTypeID(id)
	return su
}

// SetType sets the type edge to ServiceType.
func (su *ServiceUpdate) SetType(s *ServiceType) *ServiceUpdate {
	return su.SetTypeID(s.ID)
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (su *ServiceUpdate) AddDownstreamIDs(ids ...int) *ServiceUpdate {
	su.mutation.AddDownstreamIDs(ids...)
	return su
}

// AddDownstream adds the downstream edges to Service.
func (su *ServiceUpdate) AddDownstream(s ...*Service) *ServiceUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddDownstreamIDs(ids...)
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (su *ServiceUpdate) AddUpstreamIDs(ids ...int) *ServiceUpdate {
	su.mutation.AddUpstreamIDs(ids...)
	return su
}

// AddUpstream adds the upstream edges to Service.
func (su *ServiceUpdate) AddUpstream(s ...*Service) *ServiceUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddUpstreamIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (su *ServiceUpdate) AddPropertyIDs(ids ...int) *ServiceUpdate {
	su.mutation.AddPropertyIDs(ids...)
	return su
}

// AddProperties adds the properties edges to Property.
func (su *ServiceUpdate) AddProperties(p ...*Property) *ServiceUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.AddPropertyIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (su *ServiceUpdate) AddLinkIDs(ids ...int) *ServiceUpdate {
	su.mutation.AddLinkIDs(ids...)
	return su
}

// AddLinks adds the links edges to Link.
func (su *ServiceUpdate) AddLinks(l ...*Link) *ServiceUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return su.AddLinkIDs(ids...)
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (su *ServiceUpdate) AddCustomerIDs(ids ...int) *ServiceUpdate {
	su.mutation.AddCustomerIDs(ids...)
	return su
}

// AddCustomer adds the customer edges to Customer.
func (su *ServiceUpdate) AddCustomer(c ...*Customer) *ServiceUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return su.AddCustomerIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (su *ServiceUpdate) AddEndpointIDs(ids ...int) *ServiceUpdate {
	su.mutation.AddEndpointIDs(ids...)
	return su
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (su *ServiceUpdate) AddEndpoints(s ...*ServiceEndpoint) *ServiceUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddEndpointIDs(ids...)
}

// ClearType clears the type edge to ServiceType.
func (su *ServiceUpdate) ClearType() *ServiceUpdate {
	su.mutation.ClearType()
	return su
}

// RemoveDownstreamIDs removes the downstream edge to Service by ids.
func (su *ServiceUpdate) RemoveDownstreamIDs(ids ...int) *ServiceUpdate {
	su.mutation.RemoveDownstreamIDs(ids...)
	return su
}

// RemoveDownstream removes downstream edges to Service.
func (su *ServiceUpdate) RemoveDownstream(s ...*Service) *ServiceUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveDownstreamIDs(ids...)
}

// RemoveUpstreamIDs removes the upstream edge to Service by ids.
func (su *ServiceUpdate) RemoveUpstreamIDs(ids ...int) *ServiceUpdate {
	su.mutation.RemoveUpstreamIDs(ids...)
	return su
}

// RemoveUpstream removes upstream edges to Service.
func (su *ServiceUpdate) RemoveUpstream(s ...*Service) *ServiceUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveUpstreamIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (su *ServiceUpdate) RemovePropertyIDs(ids ...int) *ServiceUpdate {
	su.mutation.RemovePropertyIDs(ids...)
	return su
}

// RemoveProperties removes properties edges to Property.
func (su *ServiceUpdate) RemoveProperties(p ...*Property) *ServiceUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.RemovePropertyIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (su *ServiceUpdate) RemoveLinkIDs(ids ...int) *ServiceUpdate {
	su.mutation.RemoveLinkIDs(ids...)
	return su
}

// RemoveLinks removes links edges to Link.
func (su *ServiceUpdate) RemoveLinks(l ...*Link) *ServiceUpdate {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return su.RemoveLinkIDs(ids...)
}

// RemoveCustomerIDs removes the customer edge to Customer by ids.
func (su *ServiceUpdate) RemoveCustomerIDs(ids ...int) *ServiceUpdate {
	su.mutation.RemoveCustomerIDs(ids...)
	return su
}

// RemoveCustomer removes customer edges to Customer.
func (su *ServiceUpdate) RemoveCustomer(c ...*Customer) *ServiceUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return su.RemoveCustomerIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (su *ServiceUpdate) RemoveEndpointIDs(ids ...int) *ServiceUpdate {
	su.mutation.RemoveEndpointIDs(ids...)
	return su
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (su *ServiceUpdate) RemoveEndpoints(s ...*ServiceEndpoint) *ServiceUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (su *ServiceUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := su.mutation.UpdateTime(); !ok {
		v := service.UpdateDefaultUpdateTime()
		su.mutation.SetUpdateTime(v)
	}
	if v, ok := su.mutation.Name(); ok {
		if err := service.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := su.mutation.ExternalID(); ok {
		if err := service.ExternalIDValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}

	if _, ok := su.mutation.TypeID(); su.mutation.TypeCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"type\"")
	}

	var (
		err      error
		affected int
	)
	if len(su.hooks) == 0 {
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			su.mutation = mutation
			affected, err = su.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(su.hooks) - 1; i >= 0; i-- {
			mut = su.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, su.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
				Type:   field.TypeInt,
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
	if value, ok := su.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: service.FieldUpdateTime,
		})
	}
	if value, ok := su.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldName,
		})
	}
	if value, ok := su.mutation.ExternalID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldExternalID,
		})
	}
	if su.mutation.ExternalIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: service.FieldExternalID,
		})
	}
	if value, ok := su.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldStatus,
		})
	}
	if su.mutation.TypeCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.TypeIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.mutation.RemovedDownstreamIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.DownstreamIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.mutation.RemovedUpstreamIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.UpstreamIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.PropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.mutation.RemovedLinksIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.LinksIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.mutation.RemovedCustomerIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.CustomerIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.mutation.RemovedEndpointsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.EndpointsIDs(); len(nodes) > 0 {
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
	hooks    []Hook
	mutation *ServiceMutation
}

// SetName sets the name field.
func (suo *ServiceUpdateOne) SetName(s string) *ServiceUpdateOne {
	suo.mutation.SetName(s)
	return suo
}

// SetExternalID sets the external_id field.
func (suo *ServiceUpdateOne) SetExternalID(s string) *ServiceUpdateOne {
	suo.mutation.SetExternalID(s)
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
	suo.mutation.ClearExternalID()
	return suo
}

// SetStatus sets the status field.
func (suo *ServiceUpdateOne) SetStatus(s string) *ServiceUpdateOne {
	suo.mutation.SetStatus(s)
	return suo
}

// SetTypeID sets the type edge to ServiceType by id.
func (suo *ServiceUpdateOne) SetTypeID(id int) *ServiceUpdateOne {
	suo.mutation.SetTypeID(id)
	return suo
}

// SetType sets the type edge to ServiceType.
func (suo *ServiceUpdateOne) SetType(s *ServiceType) *ServiceUpdateOne {
	return suo.SetTypeID(s.ID)
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (suo *ServiceUpdateOne) AddDownstreamIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.AddDownstreamIDs(ids...)
	return suo
}

// AddDownstream adds the downstream edges to Service.
func (suo *ServiceUpdateOne) AddDownstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddDownstreamIDs(ids...)
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (suo *ServiceUpdateOne) AddUpstreamIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.AddUpstreamIDs(ids...)
	return suo
}

// AddUpstream adds the upstream edges to Service.
func (suo *ServiceUpdateOne) AddUpstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddUpstreamIDs(ids...)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (suo *ServiceUpdateOne) AddPropertyIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.AddPropertyIDs(ids...)
	return suo
}

// AddProperties adds the properties edges to Property.
func (suo *ServiceUpdateOne) AddProperties(p ...*Property) *ServiceUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.AddPropertyIDs(ids...)
}

// AddLinkIDs adds the links edge to Link by ids.
func (suo *ServiceUpdateOne) AddLinkIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.AddLinkIDs(ids...)
	return suo
}

// AddLinks adds the links edges to Link.
func (suo *ServiceUpdateOne) AddLinks(l ...*Link) *ServiceUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return suo.AddLinkIDs(ids...)
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (suo *ServiceUpdateOne) AddCustomerIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.AddCustomerIDs(ids...)
	return suo
}

// AddCustomer adds the customer edges to Customer.
func (suo *ServiceUpdateOne) AddCustomer(c ...*Customer) *ServiceUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return suo.AddCustomerIDs(ids...)
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (suo *ServiceUpdateOne) AddEndpointIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.AddEndpointIDs(ids...)
	return suo
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (suo *ServiceUpdateOne) AddEndpoints(s ...*ServiceEndpoint) *ServiceUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddEndpointIDs(ids...)
}

// ClearType clears the type edge to ServiceType.
func (suo *ServiceUpdateOne) ClearType() *ServiceUpdateOne {
	suo.mutation.ClearType()
	return suo
}

// RemoveDownstreamIDs removes the downstream edge to Service by ids.
func (suo *ServiceUpdateOne) RemoveDownstreamIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.RemoveDownstreamIDs(ids...)
	return suo
}

// RemoveDownstream removes downstream edges to Service.
func (suo *ServiceUpdateOne) RemoveDownstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveDownstreamIDs(ids...)
}

// RemoveUpstreamIDs removes the upstream edge to Service by ids.
func (suo *ServiceUpdateOne) RemoveUpstreamIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.RemoveUpstreamIDs(ids...)
	return suo
}

// RemoveUpstream removes upstream edges to Service.
func (suo *ServiceUpdateOne) RemoveUpstream(s ...*Service) *ServiceUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveUpstreamIDs(ids...)
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (suo *ServiceUpdateOne) RemovePropertyIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.RemovePropertyIDs(ids...)
	return suo
}

// RemoveProperties removes properties edges to Property.
func (suo *ServiceUpdateOne) RemoveProperties(p ...*Property) *ServiceUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.RemovePropertyIDs(ids...)
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (suo *ServiceUpdateOne) RemoveLinkIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.RemoveLinkIDs(ids...)
	return suo
}

// RemoveLinks removes links edges to Link.
func (suo *ServiceUpdateOne) RemoveLinks(l ...*Link) *ServiceUpdateOne {
	ids := make([]int, len(l))
	for i := range l {
		ids[i] = l[i].ID
	}
	return suo.RemoveLinkIDs(ids...)
}

// RemoveCustomerIDs removes the customer edge to Customer by ids.
func (suo *ServiceUpdateOne) RemoveCustomerIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.RemoveCustomerIDs(ids...)
	return suo
}

// RemoveCustomer removes customer edges to Customer.
func (suo *ServiceUpdateOne) RemoveCustomer(c ...*Customer) *ServiceUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return suo.RemoveCustomerIDs(ids...)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (suo *ServiceUpdateOne) RemoveEndpointIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.RemoveEndpointIDs(ids...)
	return suo
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (suo *ServiceUpdateOne) RemoveEndpoints(s ...*ServiceEndpoint) *ServiceUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveEndpointIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (suo *ServiceUpdateOne) Save(ctx context.Context) (*Service, error) {
	if _, ok := suo.mutation.UpdateTime(); !ok {
		v := service.UpdateDefaultUpdateTime()
		suo.mutation.SetUpdateTime(v)
	}
	if v, ok := suo.mutation.Name(); ok {
		if err := service.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := suo.mutation.ExternalID(); ok {
		if err := service.ExternalIDValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}

	if _, ok := suo.mutation.TypeID(); suo.mutation.TypeCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"type\"")
	}

	var (
		err  error
		node *Service
	)
	if len(suo.hooks) == 0 {
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			suo.mutation = mutation
			node, err = suo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(suo.hooks) - 1; i >= 0; i-- {
			mut = suo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, suo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: service.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Service.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := suo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: service.FieldUpdateTime,
		})
	}
	if value, ok := suo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldName,
		})
	}
	if value, ok := suo.mutation.ExternalID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldExternalID,
		})
	}
	if suo.mutation.ExternalIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: service.FieldExternalID,
		})
	}
	if value, ok := suo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldStatus,
		})
	}
	if suo.mutation.TypeCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.TypeIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.mutation.RemovedDownstreamIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.DownstreamIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.mutation.RemovedUpstreamIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.UpstreamIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.mutation.RemovedPropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.PropertiesIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.mutation.RemovedLinksIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.LinksIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.mutation.RemovedCustomerIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.CustomerIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.mutation.RemovedEndpointsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.EndpointsIDs(); len(nodes) > 0 {
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
