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
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// EquipmentPortTypeUpdate is the builder for updating EquipmentPortType entities.
type EquipmentPortTypeUpdate struct {
	config
	hooks      []Hook
	mutation   *EquipmentPortTypeMutation
	predicates []predicate.EquipmentPortType
}

// Where adds a new predicate for the builder.
func (eptu *EquipmentPortTypeUpdate) Where(ps ...predicate.EquipmentPortType) *EquipmentPortTypeUpdate {
	eptu.predicates = append(eptu.predicates, ps...)
	return eptu
}

// SetName sets the name field.
func (eptu *EquipmentPortTypeUpdate) SetName(s string) *EquipmentPortTypeUpdate {
	eptu.mutation.SetName(s)
	return eptu
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) AddPropertyTypeIDs(ids ...int) *EquipmentPortTypeUpdate {
	eptu.mutation.AddPropertyTypeIDs(ids...)
	return eptu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) AddPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.AddPropertyTypeIDs(ids...)
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) AddLinkPropertyTypeIDs(ids ...int) *EquipmentPortTypeUpdate {
	eptu.mutation.AddLinkPropertyTypeIDs(ids...)
	return eptu
}

// AddLinkPropertyTypes adds the link_property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) AddLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.AddLinkPropertyTypeIDs(ids...)
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (eptu *EquipmentPortTypeUpdate) AddPortDefinitionIDs(ids ...int) *EquipmentPortTypeUpdate {
	eptu.mutation.AddPortDefinitionIDs(ids...)
	return eptu
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (eptu *EquipmentPortTypeUpdate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptu.AddPortDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) RemovePropertyTypeIDs(ids ...int) *EquipmentPortTypeUpdate {
	eptu.mutation.RemovePropertyTypeIDs(ids...)
	return eptu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.RemovePropertyTypeIDs(ids...)
}

// RemoveLinkPropertyTypeIDs removes the link_property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) RemoveLinkPropertyTypeIDs(ids ...int) *EquipmentPortTypeUpdate {
	eptu.mutation.RemoveLinkPropertyTypeIDs(ids...)
	return eptu
}

// RemoveLinkPropertyTypes removes link_property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) RemoveLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.RemoveLinkPropertyTypeIDs(ids...)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (eptu *EquipmentPortTypeUpdate) RemovePortDefinitionIDs(ids ...int) *EquipmentPortTypeUpdate {
	eptu.mutation.RemovePortDefinitionIDs(ids...)
	return eptu
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (eptu *EquipmentPortTypeUpdate) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptu.RemovePortDefinitionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (eptu *EquipmentPortTypeUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := eptu.mutation.UpdateTime(); !ok {
		v := equipmentporttype.UpdateDefaultUpdateTime()
		eptu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(eptu.hooks) == 0 {
		affected, err = eptu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			eptu.mutation = mutation
			affected, err = eptu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(eptu.hooks) - 1; i >= 0; i-- {
			mut = eptu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, eptu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (eptu *EquipmentPortTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := eptu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eptu *EquipmentPortTypeUpdate) Exec(ctx context.Context) error {
	_, err := eptu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eptu *EquipmentPortTypeUpdate) ExecX(ctx context.Context) {
	if err := eptu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eptu *EquipmentPortTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentporttype.Table,
			Columns: equipmentporttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentporttype.FieldID,
			},
		},
	}
	if ps := eptu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eptu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentporttype.FieldUpdateTime,
		})
	}
	if value, ok := eptu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentporttype.FieldName,
		})
	}
	if nodes := eptu.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptu.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptu.mutation.RemovedLinkPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptu.mutation.LinkPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptu.mutation.RemovedPortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eptu.mutation.PortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eptu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentporttype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentPortTypeUpdateOne is the builder for updating a single EquipmentPortType entity.
type EquipmentPortTypeUpdateOne struct {
	config
	hooks    []Hook
	mutation *EquipmentPortTypeMutation
}

// SetName sets the name field.
func (eptuo *EquipmentPortTypeUpdateOne) SetName(s string) *EquipmentPortTypeUpdateOne {
	eptuo.mutation.SetName(s)
	return eptuo
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddPropertyTypeIDs(ids ...int) *EquipmentPortTypeUpdateOne {
	eptuo.mutation.AddPropertyTypeIDs(ids...)
	return eptuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.AddPropertyTypeIDs(ids...)
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddLinkPropertyTypeIDs(ids ...int) *EquipmentPortTypeUpdateOne {
	eptuo.mutation.AddLinkPropertyTypeIDs(ids...)
	return eptuo
}

// AddLinkPropertyTypes adds the link_property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) AddLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.AddLinkPropertyTypeIDs(ids...)
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddPortDefinitionIDs(ids ...int) *EquipmentPortTypeUpdateOne {
	eptuo.mutation.AddPortDefinitionIDs(ids...)
	return eptuo
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (eptuo *EquipmentPortTypeUpdateOne) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptuo.AddPortDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePropertyTypeIDs(ids ...int) *EquipmentPortTypeUpdateOne {
	eptuo.mutation.RemovePropertyTypeIDs(ids...)
	return eptuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.RemovePropertyTypeIDs(ids...)
}

// RemoveLinkPropertyTypeIDs removes the link_property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemoveLinkPropertyTypeIDs(ids ...int) *EquipmentPortTypeUpdateOne {
	eptuo.mutation.RemoveLinkPropertyTypeIDs(ids...)
	return eptuo
}

// RemoveLinkPropertyTypes removes link_property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) RemoveLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.RemoveLinkPropertyTypeIDs(ids...)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePortDefinitionIDs(ids ...int) *EquipmentPortTypeUpdateOne {
	eptuo.mutation.RemovePortDefinitionIDs(ids...)
	return eptuo
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptuo.RemovePortDefinitionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (eptuo *EquipmentPortTypeUpdateOne) Save(ctx context.Context) (*EquipmentPortType, error) {
	if _, ok := eptuo.mutation.UpdateTime(); !ok {
		v := equipmentporttype.UpdateDefaultUpdateTime()
		eptuo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *EquipmentPortType
	)
	if len(eptuo.hooks) == 0 {
		node, err = eptuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			eptuo.mutation = mutation
			node, err = eptuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(eptuo.hooks) - 1; i >= 0; i-- {
			mut = eptuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, eptuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (eptuo *EquipmentPortTypeUpdateOne) SaveX(ctx context.Context) *EquipmentPortType {
	ept, err := eptuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ept
}

// Exec executes the query on the entity.
func (eptuo *EquipmentPortTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := eptuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eptuo *EquipmentPortTypeUpdateOne) ExecX(ctx context.Context) {
	if err := eptuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eptuo *EquipmentPortTypeUpdateOne) sqlSave(ctx context.Context) (ept *EquipmentPortType, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentporttype.Table,
			Columns: equipmentporttype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentporttype.FieldID,
			},
		},
	}
	id, ok := eptuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing EquipmentPortType.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := eptuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentporttype.FieldUpdateTime,
		})
	}
	if value, ok := eptuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentporttype.FieldName,
		})
	}
	if nodes := eptuo.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptuo.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.PropertyTypesTable,
			Columns: []string{equipmentporttype.PropertyTypesColumn},
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
	if nodes := eptuo.mutation.RemovedLinkPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptuo.mutation.LinkPropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   equipmentporttype.LinkPropertyTypesTable,
			Columns: []string{equipmentporttype.LinkPropertyTypesColumn},
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
	if nodes := eptuo.mutation.RemovedPortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eptuo.mutation.PortDefinitionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentporttype.PortDefinitionsTable,
			Columns: []string{equipmentporttype.PortDefinitionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentportdefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	ept = &EquipmentPortType{config: eptuo.config}
	_spec.Assign = ept.assignValues
	_spec.ScanValues = ept.scanValues()
	if err = sqlgraph.UpdateNode(ctx, eptuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentporttype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ept, nil
}
