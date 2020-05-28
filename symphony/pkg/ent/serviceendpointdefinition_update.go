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
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/servicetype"
)

// ServiceEndpointDefinitionUpdate is the builder for updating ServiceEndpointDefinition entities.
type ServiceEndpointDefinitionUpdate struct {
	config
	hooks      []Hook
	mutation   *ServiceEndpointDefinitionMutation
	predicates []predicate.ServiceEndpointDefinition
}

// Where adds a new predicate for the builder.
func (sedu *ServiceEndpointDefinitionUpdate) Where(ps ...predicate.ServiceEndpointDefinition) *ServiceEndpointDefinitionUpdate {
	sedu.predicates = append(sedu.predicates, ps...)
	return sedu
}

// SetRole sets the role field.
func (sedu *ServiceEndpointDefinitionUpdate) SetRole(s string) *ServiceEndpointDefinitionUpdate {
	sedu.mutation.SetRole(s)
	return sedu
}

// SetNillableRole sets the role field if the given value is not nil.
func (sedu *ServiceEndpointDefinitionUpdate) SetNillableRole(s *string) *ServiceEndpointDefinitionUpdate {
	if s != nil {
		sedu.SetRole(*s)
	}
	return sedu
}

// ClearRole clears the value of role.
func (sedu *ServiceEndpointDefinitionUpdate) ClearRole() *ServiceEndpointDefinitionUpdate {
	sedu.mutation.ClearRole()
	return sedu
}

// SetName sets the name field.
func (sedu *ServiceEndpointDefinitionUpdate) SetName(s string) *ServiceEndpointDefinitionUpdate {
	sedu.mutation.SetName(s)
	return sedu
}

// SetIndex sets the index field.
func (sedu *ServiceEndpointDefinitionUpdate) SetIndex(i int) *ServiceEndpointDefinitionUpdate {
	sedu.mutation.ResetIndex()
	sedu.mutation.SetIndex(i)
	return sedu
}

// AddIndex adds i to index.
func (sedu *ServiceEndpointDefinitionUpdate) AddIndex(i int) *ServiceEndpointDefinitionUpdate {
	sedu.mutation.AddIndex(i)
	return sedu
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (sedu *ServiceEndpointDefinitionUpdate) AddEndpointIDs(ids ...int) *ServiceEndpointDefinitionUpdate {
	sedu.mutation.AddEndpointIDs(ids...)
	return sedu
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (sedu *ServiceEndpointDefinitionUpdate) AddEndpoints(s ...*ServiceEndpoint) *ServiceEndpointDefinitionUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sedu.AddEndpointIDs(ids...)
}

// SetServiceTypeID sets the service_type edge to ServiceType by id.
func (sedu *ServiceEndpointDefinitionUpdate) SetServiceTypeID(id int) *ServiceEndpointDefinitionUpdate {
	sedu.mutation.SetServiceTypeID(id)
	return sedu
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (sedu *ServiceEndpointDefinitionUpdate) SetNillableServiceTypeID(id *int) *ServiceEndpointDefinitionUpdate {
	if id != nil {
		sedu = sedu.SetServiceTypeID(*id)
	}
	return sedu
}

// SetServiceType sets the service_type edge to ServiceType.
func (sedu *ServiceEndpointDefinitionUpdate) SetServiceType(s *ServiceType) *ServiceEndpointDefinitionUpdate {
	return sedu.SetServiceTypeID(s.ID)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (sedu *ServiceEndpointDefinitionUpdate) SetEquipmentTypeID(id int) *ServiceEndpointDefinitionUpdate {
	sedu.mutation.SetEquipmentTypeID(id)
	return sedu
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (sedu *ServiceEndpointDefinitionUpdate) SetNillableEquipmentTypeID(id *int) *ServiceEndpointDefinitionUpdate {
	if id != nil {
		sedu = sedu.SetEquipmentTypeID(*id)
	}
	return sedu
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (sedu *ServiceEndpointDefinitionUpdate) SetEquipmentType(e *EquipmentType) *ServiceEndpointDefinitionUpdate {
	return sedu.SetEquipmentTypeID(e.ID)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (sedu *ServiceEndpointDefinitionUpdate) RemoveEndpointIDs(ids ...int) *ServiceEndpointDefinitionUpdate {
	sedu.mutation.RemoveEndpointIDs(ids...)
	return sedu
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (sedu *ServiceEndpointDefinitionUpdate) RemoveEndpoints(s ...*ServiceEndpoint) *ServiceEndpointDefinitionUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sedu.RemoveEndpointIDs(ids...)
}

// ClearServiceType clears the service_type edge to ServiceType.
func (sedu *ServiceEndpointDefinitionUpdate) ClearServiceType() *ServiceEndpointDefinitionUpdate {
	sedu.mutation.ClearServiceType()
	return sedu
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (sedu *ServiceEndpointDefinitionUpdate) ClearEquipmentType() *ServiceEndpointDefinitionUpdate {
	sedu.mutation.ClearEquipmentType()
	return sedu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (sedu *ServiceEndpointDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := sedu.mutation.UpdateTime(); !ok {
		v := serviceendpointdefinition.UpdateDefaultUpdateTime()
		sedu.mutation.SetUpdateTime(v)
	}
	if v, ok := sedu.mutation.Name(); ok {
		if err := serviceendpointdefinition.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(sedu.hooks) == 0 {
		affected, err = sedu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceEndpointDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sedu.mutation = mutation
			affected, err = sedu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(sedu.hooks) - 1; i >= 0; i-- {
			mut = sedu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sedu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (sedu *ServiceEndpointDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := sedu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (sedu *ServiceEndpointDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := sedu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sedu *ServiceEndpointDefinitionUpdate) ExecX(ctx context.Context) {
	if err := sedu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (sedu *ServiceEndpointDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   serviceendpointdefinition.Table,
			Columns: serviceendpointdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: serviceendpointdefinition.FieldID,
			},
		},
	}
	if ps := sedu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := sedu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: serviceendpointdefinition.FieldUpdateTime,
		})
	}
	if value, ok := sedu.mutation.Role(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: serviceendpointdefinition.FieldRole,
		})
	}
	if sedu.mutation.RoleCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: serviceendpointdefinition.FieldRole,
		})
	}
	if value, ok := sedu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: serviceendpointdefinition.FieldName,
		})
	}
	if value, ok := sedu.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: serviceendpointdefinition.FieldIndex,
		})
	}
	if value, ok := sedu.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: serviceendpointdefinition.FieldIndex,
		})
	}
	if nodes := sedu.mutation.RemovedEndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   serviceendpointdefinition.EndpointsTable,
			Columns: []string{serviceendpointdefinition.EndpointsColumn},
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
	if nodes := sedu.mutation.EndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   serviceendpointdefinition.EndpointsTable,
			Columns: []string{serviceendpointdefinition.EndpointsColumn},
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
	if sedu.mutation.ServiceTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.ServiceTypeTable,
			Columns: []string{serviceendpointdefinition.ServiceTypeColumn},
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
	if nodes := sedu.mutation.ServiceTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.ServiceTypeTable,
			Columns: []string{serviceendpointdefinition.ServiceTypeColumn},
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
	if sedu.mutation.EquipmentTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.EquipmentTypeTable,
			Columns: []string{serviceendpointdefinition.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := sedu.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.EquipmentTypeTable,
			Columns: []string{serviceendpointdefinition.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, sedu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{serviceendpointdefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ServiceEndpointDefinitionUpdateOne is the builder for updating a single ServiceEndpointDefinition entity.
type ServiceEndpointDefinitionUpdateOne struct {
	config
	hooks    []Hook
	mutation *ServiceEndpointDefinitionMutation
}

// SetRole sets the role field.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetRole(s string) *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.SetRole(s)
	return seduo
}

// SetNillableRole sets the role field if the given value is not nil.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetNillableRole(s *string) *ServiceEndpointDefinitionUpdateOne {
	if s != nil {
		seduo.SetRole(*s)
	}
	return seduo
}

// ClearRole clears the value of role.
func (seduo *ServiceEndpointDefinitionUpdateOne) ClearRole() *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.ClearRole()
	return seduo
}

// SetName sets the name field.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetName(s string) *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.SetName(s)
	return seduo
}

// SetIndex sets the index field.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetIndex(i int) *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.ResetIndex()
	seduo.mutation.SetIndex(i)
	return seduo
}

// AddIndex adds i to index.
func (seduo *ServiceEndpointDefinitionUpdateOne) AddIndex(i int) *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.AddIndex(i)
	return seduo
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (seduo *ServiceEndpointDefinitionUpdateOne) AddEndpointIDs(ids ...int) *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.AddEndpointIDs(ids...)
	return seduo
}

// AddEndpoints adds the endpoints edges to ServiceEndpoint.
func (seduo *ServiceEndpointDefinitionUpdateOne) AddEndpoints(s ...*ServiceEndpoint) *ServiceEndpointDefinitionUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return seduo.AddEndpointIDs(ids...)
}

// SetServiceTypeID sets the service_type edge to ServiceType by id.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetServiceTypeID(id int) *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.SetServiceTypeID(id)
	return seduo
}

// SetNillableServiceTypeID sets the service_type edge to ServiceType by id if the given value is not nil.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetNillableServiceTypeID(id *int) *ServiceEndpointDefinitionUpdateOne {
	if id != nil {
		seduo = seduo.SetServiceTypeID(*id)
	}
	return seduo
}

// SetServiceType sets the service_type edge to ServiceType.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetServiceType(s *ServiceType) *ServiceEndpointDefinitionUpdateOne {
	return seduo.SetServiceTypeID(s.ID)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetEquipmentTypeID(id int) *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.SetEquipmentTypeID(id)
	return seduo
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetNillableEquipmentTypeID(id *int) *ServiceEndpointDefinitionUpdateOne {
	if id != nil {
		seduo = seduo.SetEquipmentTypeID(*id)
	}
	return seduo
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (seduo *ServiceEndpointDefinitionUpdateOne) SetEquipmentType(e *EquipmentType) *ServiceEndpointDefinitionUpdateOne {
	return seduo.SetEquipmentTypeID(e.ID)
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (seduo *ServiceEndpointDefinitionUpdateOne) RemoveEndpointIDs(ids ...int) *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.RemoveEndpointIDs(ids...)
	return seduo
}

// RemoveEndpoints removes endpoints edges to ServiceEndpoint.
func (seduo *ServiceEndpointDefinitionUpdateOne) RemoveEndpoints(s ...*ServiceEndpoint) *ServiceEndpointDefinitionUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return seduo.RemoveEndpointIDs(ids...)
}

// ClearServiceType clears the service_type edge to ServiceType.
func (seduo *ServiceEndpointDefinitionUpdateOne) ClearServiceType() *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.ClearServiceType()
	return seduo
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (seduo *ServiceEndpointDefinitionUpdateOne) ClearEquipmentType() *ServiceEndpointDefinitionUpdateOne {
	seduo.mutation.ClearEquipmentType()
	return seduo
}

// Save executes the query and returns the updated entity.
func (seduo *ServiceEndpointDefinitionUpdateOne) Save(ctx context.Context) (*ServiceEndpointDefinition, error) {
	if _, ok := seduo.mutation.UpdateTime(); !ok {
		v := serviceendpointdefinition.UpdateDefaultUpdateTime()
		seduo.mutation.SetUpdateTime(v)
	}
	if v, ok := seduo.mutation.Name(); ok {
		if err := serviceendpointdefinition.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}

	var (
		err  error
		node *ServiceEndpointDefinition
	)
	if len(seduo.hooks) == 0 {
		node, err = seduo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceEndpointDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			seduo.mutation = mutation
			node, err = seduo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(seduo.hooks) - 1; i >= 0; i-- {
			mut = seduo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, seduo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (seduo *ServiceEndpointDefinitionUpdateOne) SaveX(ctx context.Context) *ServiceEndpointDefinition {
	sed, err := seduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return sed
}

// Exec executes the query on the entity.
func (seduo *ServiceEndpointDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := seduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (seduo *ServiceEndpointDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := seduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (seduo *ServiceEndpointDefinitionUpdateOne) sqlSave(ctx context.Context) (sed *ServiceEndpointDefinition, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   serviceendpointdefinition.Table,
			Columns: serviceendpointdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: serviceendpointdefinition.FieldID,
			},
		},
	}
	id, ok := seduo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing ServiceEndpointDefinition.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := seduo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: serviceendpointdefinition.FieldUpdateTime,
		})
	}
	if value, ok := seduo.mutation.Role(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: serviceendpointdefinition.FieldRole,
		})
	}
	if seduo.mutation.RoleCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: serviceendpointdefinition.FieldRole,
		})
	}
	if value, ok := seduo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: serviceendpointdefinition.FieldName,
		})
	}
	if value, ok := seduo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: serviceendpointdefinition.FieldIndex,
		})
	}
	if value, ok := seduo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: serviceendpointdefinition.FieldIndex,
		})
	}
	if nodes := seduo.mutation.RemovedEndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   serviceendpointdefinition.EndpointsTable,
			Columns: []string{serviceendpointdefinition.EndpointsColumn},
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
	if nodes := seduo.mutation.EndpointsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   serviceendpointdefinition.EndpointsTable,
			Columns: []string{serviceendpointdefinition.EndpointsColumn},
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
	if seduo.mutation.ServiceTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.ServiceTypeTable,
			Columns: []string{serviceendpointdefinition.ServiceTypeColumn},
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
	if nodes := seduo.mutation.ServiceTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.ServiceTypeTable,
			Columns: []string{serviceendpointdefinition.ServiceTypeColumn},
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
	if seduo.mutation.EquipmentTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.EquipmentTypeTable,
			Columns: []string{serviceendpointdefinition.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := seduo.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceendpointdefinition.EquipmentTypeTable,
			Columns: []string{serviceendpointdefinition.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmenttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	sed = &ServiceEndpointDefinition{config: seduo.config}
	_spec.Assign = sed.assignValues
	_spec.ScanValues = sed.scanValues()
	if err = sqlgraph.UpdateNode(ctx, seduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{serviceendpointdefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return sed, nil
}
