// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/service"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/servicetype"
)

// ServiceTypeUpdate is the builder for updating ServiceType entities.
type ServiceTypeUpdate struct {
	config
	hooks      []Hook
	mutation   *ServiceTypeMutation
	predicates []predicate.ServiceType
}

// Where adds a new predicate for the builder.
func (stu *ServiceTypeUpdate) Where(ps ...predicate.ServiceType) *ServiceTypeUpdate {
	stu.predicates = append(stu.predicates, ps...)
	return stu
}

// SetName sets the name field.
func (stu *ServiceTypeUpdate) SetName(s string) *ServiceTypeUpdate {
	stu.mutation.SetName(s)
	return stu
}

// SetHasCustomer sets the has_customer field.
func (stu *ServiceTypeUpdate) SetHasCustomer(b bool) *ServiceTypeUpdate {
	stu.mutation.SetHasCustomer(b)
	return stu
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stu *ServiceTypeUpdate) SetNillableHasCustomer(b *bool) *ServiceTypeUpdate {
	if b != nil {
		stu.SetHasCustomer(*b)
	}
	return stu
}

// SetIsDeleted sets the is_deleted field.
func (stu *ServiceTypeUpdate) SetIsDeleted(b bool) *ServiceTypeUpdate {
	stu.mutation.SetIsDeleted(b)
	return stu
}

// SetNillableIsDeleted sets the is_deleted field if the given value is not nil.
func (stu *ServiceTypeUpdate) SetNillableIsDeleted(b *bool) *ServiceTypeUpdate {
	if b != nil {
		stu.SetIsDeleted(*b)
	}
	return stu
}

// SetDiscoveryMethod sets the discovery_method field.
func (stu *ServiceTypeUpdate) SetDiscoveryMethod(sm servicetype.DiscoveryMethod) *ServiceTypeUpdate {
	stu.mutation.SetDiscoveryMethod(sm)
	return stu
}

// SetNillableDiscoveryMethod sets the discovery_method field if the given value is not nil.
func (stu *ServiceTypeUpdate) SetNillableDiscoveryMethod(sm *servicetype.DiscoveryMethod) *ServiceTypeUpdate {
	if sm != nil {
		stu.SetDiscoveryMethod(*sm)
	}
	return stu
}

// ClearDiscoveryMethod clears the value of discovery_method.
func (stu *ServiceTypeUpdate) ClearDiscoveryMethod() *ServiceTypeUpdate {
	stu.mutation.ClearDiscoveryMethod()
	return stu
}

// AddServiceIDs adds the services edge to Service by ids.
func (stu *ServiceTypeUpdate) AddServiceIDs(ids ...int) *ServiceTypeUpdate {
	stu.mutation.AddServiceIDs(ids...)
	return stu
}

// AddServices adds the services edges to Service.
func (stu *ServiceTypeUpdate) AddServices(s ...*Service) *ServiceTypeUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stu.AddServiceIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (stu *ServiceTypeUpdate) AddPropertyTypeIDs(ids ...int) *ServiceTypeUpdate {
	stu.mutation.AddPropertyTypeIDs(ids...)
	return stu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (stu *ServiceTypeUpdate) AddPropertyTypes(p ...*PropertyType) *ServiceTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stu.AddPropertyTypeIDs(ids...)
}

// AddEndpointDefinitionIDs adds the endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (stu *ServiceTypeUpdate) AddEndpointDefinitionIDs(ids ...int) *ServiceTypeUpdate {
	stu.mutation.AddEndpointDefinitionIDs(ids...)
	return stu
}

// AddEndpointDefinitions adds the endpoint_definitions edges to ServiceEndpointDefinition.
func (stu *ServiceTypeUpdate) AddEndpointDefinitions(s ...*ServiceEndpointDefinition) *ServiceTypeUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stu.AddEndpointDefinitionIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (stu *ServiceTypeUpdate) RemoveServiceIDs(ids ...int) *ServiceTypeUpdate {
	stu.mutation.RemoveServiceIDs(ids...)
	return stu
}

// RemoveServices removes services edges to Service.
func (stu *ServiceTypeUpdate) RemoveServices(s ...*Service) *ServiceTypeUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stu.RemoveServiceIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (stu *ServiceTypeUpdate) RemovePropertyTypeIDs(ids ...int) *ServiceTypeUpdate {
	stu.mutation.RemovePropertyTypeIDs(ids...)
	return stu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (stu *ServiceTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *ServiceTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stu.RemovePropertyTypeIDs(ids...)
}

// RemoveEndpointDefinitionIDs removes the endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (stu *ServiceTypeUpdate) RemoveEndpointDefinitionIDs(ids ...int) *ServiceTypeUpdate {
	stu.mutation.RemoveEndpointDefinitionIDs(ids...)
	return stu
}

// RemoveEndpointDefinitions removes endpoint_definitions edges to ServiceEndpointDefinition.
func (stu *ServiceTypeUpdate) RemoveEndpointDefinitions(s ...*ServiceEndpointDefinition) *ServiceTypeUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stu.RemoveEndpointDefinitionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (stu *ServiceTypeUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := stu.mutation.UpdateTime(); !ok {
		v := servicetype.UpdateDefaultUpdateTime()
		stu.mutation.SetUpdateTime(v)
	}
	if v, ok := stu.mutation.DiscoveryMethod(); ok {
		if err := servicetype.DiscoveryMethodValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"discovery_method\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(stu.hooks) == 0 {
		affected, err = stu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			stu.mutation = mutation
			affected, err = stu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(stu.hooks) - 1; i >= 0; i-- {
			mut = stu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, stu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
				Type:   field.TypeInt,
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
	if value, ok := stu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: servicetype.FieldUpdateTime,
		})
	}
	if value, ok := stu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: servicetype.FieldName,
		})
	}
	if value, ok := stu.mutation.HasCustomer(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: servicetype.FieldHasCustomer,
		})
	}
	if value, ok := stu.mutation.IsDeleted(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: servicetype.FieldIsDeleted,
		})
	}
	if value, ok := stu.mutation.DiscoveryMethod(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: servicetype.FieldDiscoveryMethod,
		})
	}
	if stu.mutation.DiscoveryMethodCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Column: servicetype.FieldDiscoveryMethod,
		})
	}
	if nodes := stu.mutation.RemovedServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stu.mutation.ServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stu.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stu.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := stu.mutation.RemovedEndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.EndpointDefinitionsTable,
			Columns: []string{servicetype.EndpointDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpointdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stu.mutation.EndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.EndpointDefinitionsTable,
			Columns: []string{servicetype.EndpointDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpointdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, stu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{servicetype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ServiceTypeUpdateOne is the builder for updating a single ServiceType entity.
type ServiceTypeUpdateOne struct {
	config
	hooks    []Hook
	mutation *ServiceTypeMutation
}

// SetName sets the name field.
func (stuo *ServiceTypeUpdateOne) SetName(s string) *ServiceTypeUpdateOne {
	stuo.mutation.SetName(s)
	return stuo
}

// SetHasCustomer sets the has_customer field.
func (stuo *ServiceTypeUpdateOne) SetHasCustomer(b bool) *ServiceTypeUpdateOne {
	stuo.mutation.SetHasCustomer(b)
	return stuo
}

// SetNillableHasCustomer sets the has_customer field if the given value is not nil.
func (stuo *ServiceTypeUpdateOne) SetNillableHasCustomer(b *bool) *ServiceTypeUpdateOne {
	if b != nil {
		stuo.SetHasCustomer(*b)
	}
	return stuo
}

// SetIsDeleted sets the is_deleted field.
func (stuo *ServiceTypeUpdateOne) SetIsDeleted(b bool) *ServiceTypeUpdateOne {
	stuo.mutation.SetIsDeleted(b)
	return stuo
}

// SetNillableIsDeleted sets the is_deleted field if the given value is not nil.
func (stuo *ServiceTypeUpdateOne) SetNillableIsDeleted(b *bool) *ServiceTypeUpdateOne {
	if b != nil {
		stuo.SetIsDeleted(*b)
	}
	return stuo
}

// SetDiscoveryMethod sets the discovery_method field.
func (stuo *ServiceTypeUpdateOne) SetDiscoveryMethod(sm servicetype.DiscoveryMethod) *ServiceTypeUpdateOne {
	stuo.mutation.SetDiscoveryMethod(sm)
	return stuo
}

// SetNillableDiscoveryMethod sets the discovery_method field if the given value is not nil.
func (stuo *ServiceTypeUpdateOne) SetNillableDiscoveryMethod(sm *servicetype.DiscoveryMethod) *ServiceTypeUpdateOne {
	if sm != nil {
		stuo.SetDiscoveryMethod(*sm)
	}
	return stuo
}

// ClearDiscoveryMethod clears the value of discovery_method.
func (stuo *ServiceTypeUpdateOne) ClearDiscoveryMethod() *ServiceTypeUpdateOne {
	stuo.mutation.ClearDiscoveryMethod()
	return stuo
}

// AddServiceIDs adds the services edge to Service by ids.
func (stuo *ServiceTypeUpdateOne) AddServiceIDs(ids ...int) *ServiceTypeUpdateOne {
	stuo.mutation.AddServiceIDs(ids...)
	return stuo
}

// AddServices adds the services edges to Service.
func (stuo *ServiceTypeUpdateOne) AddServices(s ...*Service) *ServiceTypeUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stuo.AddServiceIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (stuo *ServiceTypeUpdateOne) AddPropertyTypeIDs(ids ...int) *ServiceTypeUpdateOne {
	stuo.mutation.AddPropertyTypeIDs(ids...)
	return stuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (stuo *ServiceTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *ServiceTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stuo.AddPropertyTypeIDs(ids...)
}

// AddEndpointDefinitionIDs adds the endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (stuo *ServiceTypeUpdateOne) AddEndpointDefinitionIDs(ids ...int) *ServiceTypeUpdateOne {
	stuo.mutation.AddEndpointDefinitionIDs(ids...)
	return stuo
}

// AddEndpointDefinitions adds the endpoint_definitions edges to ServiceEndpointDefinition.
func (stuo *ServiceTypeUpdateOne) AddEndpointDefinitions(s ...*ServiceEndpointDefinition) *ServiceTypeUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stuo.AddEndpointDefinitionIDs(ids...)
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (stuo *ServiceTypeUpdateOne) RemoveServiceIDs(ids ...int) *ServiceTypeUpdateOne {
	stuo.mutation.RemoveServiceIDs(ids...)
	return stuo
}

// RemoveServices removes services edges to Service.
func (stuo *ServiceTypeUpdateOne) RemoveServices(s ...*Service) *ServiceTypeUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stuo.RemoveServiceIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (stuo *ServiceTypeUpdateOne) RemovePropertyTypeIDs(ids ...int) *ServiceTypeUpdateOne {
	stuo.mutation.RemovePropertyTypeIDs(ids...)
	return stuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (stuo *ServiceTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *ServiceTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return stuo.RemovePropertyTypeIDs(ids...)
}

// RemoveEndpointDefinitionIDs removes the endpoint_definitions edge to ServiceEndpointDefinition by ids.
func (stuo *ServiceTypeUpdateOne) RemoveEndpointDefinitionIDs(ids ...int) *ServiceTypeUpdateOne {
	stuo.mutation.RemoveEndpointDefinitionIDs(ids...)
	return stuo
}

// RemoveEndpointDefinitions removes endpoint_definitions edges to ServiceEndpointDefinition.
func (stuo *ServiceTypeUpdateOne) RemoveEndpointDefinitions(s ...*ServiceEndpointDefinition) *ServiceTypeUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return stuo.RemoveEndpointDefinitionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (stuo *ServiceTypeUpdateOne) Save(ctx context.Context) (*ServiceType, error) {
	if _, ok := stuo.mutation.UpdateTime(); !ok {
		v := servicetype.UpdateDefaultUpdateTime()
		stuo.mutation.SetUpdateTime(v)
	}
	if v, ok := stuo.mutation.DiscoveryMethod(); ok {
		if err := servicetype.DiscoveryMethodValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"discovery_method\": %v", err)
		}
	}

	var (
		err  error
		node *ServiceType
	)
	if len(stuo.hooks) == 0 {
		node, err = stuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			stuo.mutation = mutation
			node, err = stuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(stuo.hooks) - 1; i >= 0; i-- {
			mut = stuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, stuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: servicetype.FieldID,
			},
		},
	}
	id, ok := stuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing ServiceType.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := stuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: servicetype.FieldUpdateTime,
		})
	}
	if value, ok := stuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: servicetype.FieldName,
		})
	}
	if value, ok := stuo.mutation.HasCustomer(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: servicetype.FieldHasCustomer,
		})
	}
	if value, ok := stuo.mutation.IsDeleted(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: servicetype.FieldIsDeleted,
		})
	}
	if value, ok := stuo.mutation.DiscoveryMethod(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: servicetype.FieldDiscoveryMethod,
		})
	}
	if stuo.mutation.DiscoveryMethodCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Column: servicetype.FieldDiscoveryMethod,
		})
	}
	if nodes := stuo.mutation.RemovedServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stuo.mutation.ServicesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   servicetype.ServicesTable,
			Columns: []string{servicetype.ServicesColumn},
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
	if nodes := stuo.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stuo.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.PropertyTypesTable,
			Columns: []string{servicetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := stuo.mutation.RemovedEndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.EndpointDefinitionsTable,
			Columns: []string{servicetype.EndpointDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpointdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stuo.mutation.EndpointDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   servicetype.EndpointDefinitionsTable,
			Columns: []string{servicetype.EndpointDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: serviceendpointdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	st = &ServiceType{config: stuo.config}
	_spec.Assign = st.assignValues
	_spec.ScanValues = st.scanValues()
	if err = sqlgraph.UpdateNode(ctx, stuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{servicetype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return st, nil
}
