// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/service"
)

// CustomerUpdate is the builder for updating Customer entities.
type CustomerUpdate struct {
	config

	update_time      *time.Time
	name             *string
	external_id      *string
	clearexternal_id bool
	services         map[int]struct{}
	removedServices  map[int]struct{}
	predicates       []predicate.Customer
}

// Where adds a new predicate for the builder.
func (cu *CustomerUpdate) Where(ps ...predicate.Customer) *CustomerUpdate {
	cu.predicates = append(cu.predicates, ps...)
	return cu
}

// SetName sets the name field.
func (cu *CustomerUpdate) SetName(s string) *CustomerUpdate {
	cu.name = &s
	return cu
}

// SetExternalID sets the external_id field.
func (cu *CustomerUpdate) SetExternalID(s string) *CustomerUpdate {
	cu.external_id = &s
	return cu
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (cu *CustomerUpdate) SetNillableExternalID(s *string) *CustomerUpdate {
	if s != nil {
		cu.SetExternalID(*s)
	}
	return cu
}

// ClearExternalID clears the value of external_id.
func (cu *CustomerUpdate) ClearExternalID() *CustomerUpdate {
	cu.external_id = nil
	cu.clearexternal_id = true
	return cu
}

// AddServiceIDs adds the services edge to Service by ids.
func (cu *CustomerUpdate) AddServiceIDs(ids ...int) *CustomerUpdate {
	if cu.services == nil {
		cu.services = make(map[int]struct{})
	}
	for i := range ids {
		cu.services[ids[i]] = struct{}{}
	}
	return cu
}

// AddServices adds the services edges to Service.
func (cu *CustomerUpdate) AddServices(s ...*Service) *CustomerUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cu.AddServiceIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (cu *CustomerUpdate) RemoveServiceIDs(ids ...int) *CustomerUpdate {
	if cu.removedServices == nil {
		cu.removedServices = make(map[int]struct{})
	}
	for i := range ids {
		cu.removedServices[ids[i]] = struct{}{}
	}
	return cu
}

// RemoveServices removes services edges to Service.
func (cu *CustomerUpdate) RemoveServices(s ...*Service) *CustomerUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cu.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (cu *CustomerUpdate) Save(ctx context.Context) (int, error) {
	if cu.update_time == nil {
		v := customer.UpdateDefaultUpdateTime()
		cu.update_time = &v
	}
	if cu.name != nil {
		if err := customer.NameValidator(*cu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if cu.external_id != nil {
		if err := customer.ExternalIDValidator(*cu.external_id); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	return cu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CustomerUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CustomerUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CustomerUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cu *CustomerUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   customer.Table,
			Columns: customer.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: customer.FieldID,
			},
		},
	}
	if ps := cu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := cu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: customer.FieldUpdateTime,
		})
	}
	if value := cu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: customer.FieldName,
		})
	}
	if value := cu.external_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: customer.FieldExternalID,
		})
	}
	if cu.clearexternal_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: customer.FieldExternalID,
		})
	}
	if nodes := cu.removedServices; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   customer.ServicesTable,
			Columns: customer.ServicesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.services; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   customer.ServicesTable,
			Columns: customer.ServicesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{customer.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// CustomerUpdateOne is the builder for updating a single Customer entity.
type CustomerUpdateOne struct {
	config
	id int

	update_time      *time.Time
	name             *string
	external_id      *string
	clearexternal_id bool
	services         map[int]struct{}
	removedServices  map[int]struct{}
}

// SetName sets the name field.
func (cuo *CustomerUpdateOne) SetName(s string) *CustomerUpdateOne {
	cuo.name = &s
	return cuo
}

// SetExternalID sets the external_id field.
func (cuo *CustomerUpdateOne) SetExternalID(s string) *CustomerUpdateOne {
	cuo.external_id = &s
	return cuo
}

// SetNillableExternalID sets the external_id field if the given value is not nil.
func (cuo *CustomerUpdateOne) SetNillableExternalID(s *string) *CustomerUpdateOne {
	if s != nil {
		cuo.SetExternalID(*s)
	}
	return cuo
}

// ClearExternalID clears the value of external_id.
func (cuo *CustomerUpdateOne) ClearExternalID() *CustomerUpdateOne {
	cuo.external_id = nil
	cuo.clearexternal_id = true
	return cuo
}

// AddServiceIDs adds the services edge to Service by ids.
func (cuo *CustomerUpdateOne) AddServiceIDs(ids ...int) *CustomerUpdateOne {
	if cuo.services == nil {
		cuo.services = make(map[int]struct{})
	}
	for i := range ids {
		cuo.services[ids[i]] = struct{}{}
	}
	return cuo
}

// AddServices adds the services edges to Service.
func (cuo *CustomerUpdateOne) AddServices(s ...*Service) *CustomerUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cuo.AddServiceIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (cuo *CustomerUpdateOne) RemoveServiceIDs(ids ...int) *CustomerUpdateOne {
	if cuo.removedServices == nil {
		cuo.removedServices = make(map[int]struct{})
	}
	for i := range ids {
		cuo.removedServices[ids[i]] = struct{}{}
	}
	return cuo
}

// RemoveServices removes services edges to Service.
func (cuo *CustomerUpdateOne) RemoveServices(s ...*Service) *CustomerUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return cuo.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (cuo *CustomerUpdateOne) Save(ctx context.Context) (*Customer, error) {
	if cuo.update_time == nil {
		v := customer.UpdateDefaultUpdateTime()
		cuo.update_time = &v
	}
	if cuo.name != nil {
		if err := customer.NameValidator(*cuo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if cuo.external_id != nil {
		if err := customer.ExternalIDValidator(*cuo.external_id); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"external_id\": %v", err)
		}
	}
	return cuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CustomerUpdateOne) SaveX(ctx context.Context) *Customer {
	c, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return c
}

// Exec executes the query on the entity.
func (cuo *CustomerUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CustomerUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cuo *CustomerUpdateOne) sqlSave(ctx context.Context) (c *Customer, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   customer.Table,
			Columns: customer.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  cuo.id,
				Type:   field.TypeInt,
				Column: customer.FieldID,
			},
		},
	}
	if value := cuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: customer.FieldUpdateTime,
		})
	}
	if value := cuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: customer.FieldName,
		})
	}
	if value := cuo.external_id; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: customer.FieldExternalID,
		})
	}
	if cuo.clearexternal_id {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: customer.FieldExternalID,
		})
	}
	if nodes := cuo.removedServices; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   customer.ServicesTable,
			Columns: customer.ServicesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.services; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   customer.ServicesTable,
			Columns: customer.ServicesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	c = &Customer{config: cuo.config}
	_spec.Assign = c.assignValues
	_spec.ScanValues = c.scanValues()
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{customer.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return c, nil
}
