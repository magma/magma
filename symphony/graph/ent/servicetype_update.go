// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
)

// ServiceTypeUpdate is the builder for updating ServiceType entities.
type ServiceTypeUpdate struct {
	config

	update_time          *time.Time
	name                 *string
	has_customer         *bool
	services             map[string]struct{}
	property_types       map[string]struct{}
	removedServices      map[string]struct{}
	removedPropertyTypes map[string]struct{}
	predicates           []predicate.ServiceType
}

// Where adds a new predicate for the builder.
func (stu *ServiceTypeUpdate) Where(ps ...predicate.ServiceType) *ServiceTypeUpdate {
	stu.predicates = append(stu.predicates, ps...)
	return stu
}

// SetName sets the name field.
func (stu *ServiceTypeUpdate) SetName(s string) *ServiceTypeUpdate {
	stu.name = &s
	return stu
}

// SetHasCustomer sets the has_customer field.
func (stu *ServiceTypeUpdate) SetHasCustomer(b bool) *ServiceTypeUpdate {
	stu.has_customer = &b
	return stu
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stu *ServiceTypeUpdate) SetNillableHasCustomer(b *bool) *ServiceTypeUpdate {
	if b != nil {
		stu.SetHasCustomer(*b)
	}
	return stu
}

// AddServiceIDs adds the services edge to Service by ids.
func (stu *ServiceTypeUpdate) AddServiceIDs(ids ...string) *ServiceTypeUpdate {
	if stu.services == nil {
		stu.services = make(map[string]struct{})
	}
	for i := range ids {
		stu.services[ids[i]] = struct{}{}
	}
	return stu
}

// AddServices adds the services edges to Service.
func (stu *ServiceTypeUpdate) AddServices(s ...*Service) *ServiceTypeUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stu.AddServiceIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (stu *ServiceTypeUpdate) AddPropertyTypeIDs(ids ...string) *ServiceTypeUpdate {
	if stu.property_types == nil {
		stu.property_types = make(map[string]struct{})
	}
	for i := range ids {
		stu.property_types[ids[i]] = struct{}{}
	}
	return stu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (stu *ServiceTypeUpdate) AddPropertyTypes(p ...*PropertyType) *ServiceTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stu.AddPropertyTypeIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (stu *ServiceTypeUpdate) RemoveServiceIDs(ids ...string) *ServiceTypeUpdate {
	if stu.removedServices == nil {
		stu.removedServices = make(map[string]struct{})
	}
	for i := range ids {
		stu.removedServices[ids[i]] = struct{}{}
	}
	return stu
}

// RemoveServices removes services edges to Service.
func (stu *ServiceTypeUpdate) RemoveServices(s ...*Service) *ServiceTypeUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stu.RemoveServiceIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (stu *ServiceTypeUpdate) RemovePropertyTypeIDs(ids ...string) *ServiceTypeUpdate {
	if stu.removedPropertyTypes == nil {
		stu.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		stu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return stu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (stu *ServiceTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *ServiceTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stu.RemovePropertyTypeIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (stu *ServiceTypeUpdate) Save(ctx context.Context) (int, error) {
	if stu.update_time == nil {
		v := servicetype.UpdateDefaultUpdateTime()
		stu.update_time = &v
	}
	return stu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stu *ServiceTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := stu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (stu *ServiceTypeUpdate) Exec(ctx context.Context) error {
	_, err := stu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stu *ServiceTypeUpdate) ExecX(ctx context.Context) {
	if err := stu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stu *ServiceTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   servicetype.Table,
			Columns: servicetype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: servicetype.FieldID,
			},
		},
	}
	if ps := stu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := stu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: servicetype.FieldUpdateTime,
		})
	}
	if value := stu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: servicetype.FieldName,
		})
	}
	if value := stu.has_customer; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: servicetype.FieldHasCustomer,
		})
	}
	if nodes := stu.removedServices; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stu.services; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stu.removedPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
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
	if nodes := stu.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
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
	if n, err = sqlgraph.UpdateNodes(ctx, stu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ServiceTypeUpdateOne is the builder for updating a single ServiceType entity.
type ServiceTypeUpdateOne struct {
	config
	id string

	update_time          *time.Time
	name                 *string
	has_customer         *bool
	services             map[string]struct{}
	property_types       map[string]struct{}
	removedServices      map[string]struct{}
	removedPropertyTypes map[string]struct{}
}

// SetName sets the name field.
func (stuo *ServiceTypeUpdateOne) SetName(s string) *ServiceTypeUpdateOne {
	stuo.name = &s
	return stuo
}

// SetHasCustomer sets the has_customer field.
func (stuo *ServiceTypeUpdateOne) SetHasCustomer(b bool) *ServiceTypeUpdateOne {
	stuo.has_customer = &b
	return stuo
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stuo *ServiceTypeUpdateOne) SetNillableHasCustomer(b *bool) *ServiceTypeUpdateOne {
	if b != nil {
		stuo.SetHasCustomer(*b)
	}
	return stuo
}

// AddServiceIDs adds the services edge to Service by ids.
func (stuo *ServiceTypeUpdateOne) AddServiceIDs(ids ...string) *ServiceTypeUpdateOne {
	if stuo.services == nil {
		stuo.services = make(map[string]struct{})
	}
	for i := range ids {
		stuo.services[ids[i]] = struct{}{}
	}
	return stuo
}

// AddServices adds the services edges to Service.
func (stuo *ServiceTypeUpdateOne) AddServices(s ...*Service) *ServiceTypeUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stuo.AddServiceIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (stuo *ServiceTypeUpdateOne) AddPropertyTypeIDs(ids ...string) *ServiceTypeUpdateOne {
	if stuo.property_types == nil {
		stuo.property_types = make(map[string]struct{})
	}
	for i := range ids {
		stuo.property_types[ids[i]] = struct{}{}
	}
	return stuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (stuo *ServiceTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *ServiceTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stuo.AddPropertyTypeIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (stuo *ServiceTypeUpdateOne) RemoveServiceIDs(ids ...string) *ServiceTypeUpdateOne {
	if stuo.removedServices == nil {
		stuo.removedServices = make(map[string]struct{})
	}
	for i := range ids {
		stuo.removedServices[ids[i]] = struct{}{}
	}
	return stuo
}

// RemoveServices removes services edges to Service.
func (stuo *ServiceTypeUpdateOne) RemoveServices(s ...*Service) *ServiceTypeUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stuo.RemoveServiceIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (stuo *ServiceTypeUpdateOne) RemovePropertyTypeIDs(ids ...string) *ServiceTypeUpdateOne {
	if stuo.removedPropertyTypes == nil {
		stuo.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		stuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return stuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (stuo *ServiceTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *ServiceTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stuo.RemovePropertyTypeIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (stuo *ServiceTypeUpdateOne) Save(ctx context.Context) (*ServiceType, error) {
	if stuo.update_time == nil {
		v := servicetype.UpdateDefaultUpdateTime()
		stuo.update_time = &v
	}
	return stuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stuo *ServiceTypeUpdateOne) SaveX(ctx context.Context) *ServiceType {
	st, err := stuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return st
}

// Exec executes the query on the entity.
func (stuo *ServiceTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := stuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stuo *ServiceTypeUpdateOne) ExecX(ctx context.Context) {
	if err := stuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stuo *ServiceTypeUpdateOne) sqlSave(ctx context.Context) (st *ServiceType, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   servicetype.Table,
			Columns: servicetype.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  stuo.id,
				Type:   field.TypeString,
				Column: servicetype.FieldID,
			},
		},
	}
	if value := stuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: servicetype.FieldUpdateTime,
		})
	}
	if value := stuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: servicetype.FieldName,
		})
	}
	if value := stuo.has_customer; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  *value,
			Column: servicetype.FieldHasCustomer,
		})
	}
	if nodes := stuo.removedServices; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stuo.services; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stuo.removedPropertyTypes; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
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
	if nodes := stuo.property_types; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: propertytype.FieldID,
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
	st = &ServiceType{config: stuo.config}
	_spec.Assign = st.assignValues
	_spec.ScanValues = st.scanValues()
	if err = sqlgraph.UpdateNode(ctx, stuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return st, nil
}
