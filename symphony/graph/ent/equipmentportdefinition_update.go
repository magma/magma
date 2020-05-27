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
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPortDefinitionUpdate is the builder for updating EquipmentPortDefinition entities.
type EquipmentPortDefinitionUpdate struct {
	config
	hooks      []Hook
	mutation   *EquipmentPortDefinitionMutation
	predicates []predicate.EquipmentPortDefinition
}

// Where adds a new predicate for the builder.
func (epdu *EquipmentPortDefinitionUpdate) Where(ps ...predicate.EquipmentPortDefinition) *EquipmentPortDefinitionUpdate {
	epdu.predicates = append(epdu.predicates, ps...)
	return epdu
}

// SetName sets the name field.
func (epdu *EquipmentPortDefinitionUpdate) SetName(s string) *EquipmentPortDefinitionUpdate {
	epdu.mutation.SetName(s)
	return epdu
}

// SetIndex sets the index field.
func (epdu *EquipmentPortDefinitionUpdate) SetIndex(i int) *EquipmentPortDefinitionUpdate {
	epdu.mutation.ResetIndex()
	epdu.mutation.SetIndex(i)
	return epdu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableIndex(i *int) *EquipmentPortDefinitionUpdate {
	if i != nil {
		epdu.SetIndex(*i)
	}
	return epdu
}

// AddIndex adds i to index.
func (epdu *EquipmentPortDefinitionUpdate) AddIndex(i int) *EquipmentPortDefinitionUpdate {
	epdu.mutation.AddIndex(i)
	return epdu
}

// ClearIndex clears the value of index.
func (epdu *EquipmentPortDefinitionUpdate) ClearIndex() *EquipmentPortDefinitionUpdate {
	epdu.mutation.ClearIndex()
	return epdu
}

// SetBandwidth sets the bandwidth field.
func (epdu *EquipmentPortDefinitionUpdate) SetBandwidth(s string) *EquipmentPortDefinitionUpdate {
	epdu.mutation.SetBandwidth(s)
	return epdu
}

// SetNillableBandwidth sets the bandwidth field if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableBandwidth(s *string) *EquipmentPortDefinitionUpdate {
	if s != nil {
		epdu.SetBandwidth(*s)
	}
	return epdu
}

// ClearBandwidth clears the value of bandwidth.
func (epdu *EquipmentPortDefinitionUpdate) ClearBandwidth() *EquipmentPortDefinitionUpdate {
	epdu.mutation.ClearBandwidth()
	return epdu
}

// SetVisibilityLabel sets the visibility_label field.
func (epdu *EquipmentPortDefinitionUpdate) SetVisibilityLabel(s string) *EquipmentPortDefinitionUpdate {
	epdu.mutation.SetVisibilityLabel(s)
	return epdu
}

// SetNillableVisibilityLabel sets the visibility_label field if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableVisibilityLabel(s *string) *EquipmentPortDefinitionUpdate {
	if s != nil {
		epdu.SetVisibilityLabel(*s)
	}
	return epdu
}

// ClearVisibilityLabel clears the value of visibility_label.
func (epdu *EquipmentPortDefinitionUpdate) ClearVisibilityLabel() *EquipmentPortDefinitionUpdate {
	epdu.mutation.ClearVisibilityLabel()
	return epdu
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (epdu *EquipmentPortDefinitionUpdate) SetEquipmentPortTypeID(id int) *EquipmentPortDefinitionUpdate {
	epdu.mutation.SetEquipmentPortTypeID(id)
	return epdu
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableEquipmentPortTypeID(id *int) *EquipmentPortDefinitionUpdate {
	if id != nil {
		epdu = epdu.SetEquipmentPortTypeID(*id)
	}
	return epdu
}

// SetEquipmentPortType sets the equipment_port_type edge to EquipmentPortType.
func (epdu *EquipmentPortDefinitionUpdate) SetEquipmentPortType(e *EquipmentPortType) *EquipmentPortDefinitionUpdate {
	return epdu.SetEquipmentPortTypeID(e.ID)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (epdu *EquipmentPortDefinitionUpdate) AddPortIDs(ids ...int) *EquipmentPortDefinitionUpdate {
	epdu.mutation.AddPortIDs(ids...)
	return epdu
}

// AddPorts adds the ports edges to EquipmentPort.
func (epdu *EquipmentPortDefinitionUpdate) AddPorts(e ...*EquipmentPort) *EquipmentPortDefinitionUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdu.AddPortIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epdu *EquipmentPortDefinitionUpdate) SetEquipmentTypeID(id int) *EquipmentPortDefinitionUpdate {
	epdu.mutation.SetEquipmentTypeID(id)
	return epdu
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epdu *EquipmentPortDefinitionUpdate) SetNillableEquipmentTypeID(id *int) *EquipmentPortDefinitionUpdate {
	if id != nil {
		epdu = epdu.SetEquipmentTypeID(*id)
	}
	return epdu
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (epdu *EquipmentPortDefinitionUpdate) SetEquipmentType(e *EquipmentType) *EquipmentPortDefinitionUpdate {
	return epdu.SetEquipmentTypeID(e.ID)
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (epdu *EquipmentPortDefinitionUpdate) ClearEquipmentPortType() *EquipmentPortDefinitionUpdate {
	epdu.mutation.ClearEquipmentPortType()
	return epdu
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (epdu *EquipmentPortDefinitionUpdate) RemovePortIDs(ids ...int) *EquipmentPortDefinitionUpdate {
	epdu.mutation.RemovePortIDs(ids...)
	return epdu
}

// RemovePorts removes ports edges to EquipmentPort.
func (epdu *EquipmentPortDefinitionUpdate) RemovePorts(e ...*EquipmentPort) *EquipmentPortDefinitionUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdu.RemovePortIDs(ids...)
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (epdu *EquipmentPortDefinitionUpdate) ClearEquipmentType() *EquipmentPortDefinitionUpdate {
	epdu.mutation.ClearEquipmentType()
	return epdu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (epdu *EquipmentPortDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := epdu.mutation.UpdateTime(); !ok {
		v := equipmentportdefinition.UpdateDefaultUpdateTime()
		epdu.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(epdu.hooks) == 0 {
		affected, err = epdu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epdu.mutation = mutation
			affected, err = epdu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(epdu.hooks) - 1; i >= 0; i-- {
			mut = epdu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epdu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (epdu *EquipmentPortDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := epdu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (epdu *EquipmentPortDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := epdu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epdu *EquipmentPortDefinitionUpdate) ExecX(ctx context.Context) {
	if err := epdu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epdu *EquipmentPortDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentportdefinition.Table,
			Columns: equipmentportdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentportdefinition.FieldID,
			},
		},
	}
	if ps := epdu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := epdu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentportdefinition.FieldUpdateTime,
		})
	}
	if value, ok := epdu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldName,
		})
	}
	if value, ok := epdu.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentportdefinition.FieldIndex,
		})
	}
	if value, ok := epdu.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentportdefinition.FieldIndex,
		})
	}
	if epdu.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: equipmentportdefinition.FieldIndex,
		})
	}
	if value, ok := epdu.mutation.Bandwidth(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldBandwidth,
		})
	}
	if epdu.mutation.BandwidthCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipmentportdefinition.FieldBandwidth,
		})
	}
	if value, ok := epdu.mutation.VisibilityLabel(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldVisibilityLabel,
		})
	}
	if epdu.mutation.VisibilityLabelCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipmentportdefinition.FieldVisibilityLabel,
		})
	}
	if epdu.mutation.EquipmentPortTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentportdefinition.EquipmentPortTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epdu.mutation.EquipmentPortTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentportdefinition.EquipmentPortTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := epdu.mutation.RemovedPortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentportdefinition.PortsTable,
			Columns: []string{equipmentportdefinition.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epdu.mutation.PortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentportdefinition.PortsTable,
			Columns: []string{equipmentportdefinition.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epdu.mutation.EquipmentTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentportdefinition.EquipmentTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentTypeColumn},
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
	if nodes := epdu.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentportdefinition.EquipmentTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentTypeColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, epdu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentportdefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentPortDefinitionUpdateOne is the builder for updating a single EquipmentPortDefinition entity.
type EquipmentPortDefinitionUpdateOne struct {
	config
	hooks    []Hook
	mutation *EquipmentPortDefinitionMutation
}

// SetName sets the name field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetName(s string) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.SetName(s)
	return epduo
}

// SetIndex sets the index field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetIndex(i int) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.ResetIndex()
	epduo.mutation.SetIndex(i)
	return epduo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableIndex(i *int) *EquipmentPortDefinitionUpdateOne {
	if i != nil {
		epduo.SetIndex(*i)
	}
	return epduo
}

// AddIndex adds i to index.
func (epduo *EquipmentPortDefinitionUpdateOne) AddIndex(i int) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.AddIndex(i)
	return epduo
}

// ClearIndex clears the value of index.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearIndex() *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.ClearIndex()
	return epduo
}

// SetBandwidth sets the bandwidth field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetBandwidth(s string) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.SetBandwidth(s)
	return epduo
}

// SetNillableBandwidth sets the bandwidth field if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableBandwidth(s *string) *EquipmentPortDefinitionUpdateOne {
	if s != nil {
		epduo.SetBandwidth(*s)
	}
	return epduo
}

// ClearBandwidth clears the value of bandwidth.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearBandwidth() *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.ClearBandwidth()
	return epduo
}

// SetVisibilityLabel sets the visibility_label field.
func (epduo *EquipmentPortDefinitionUpdateOne) SetVisibilityLabel(s string) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.SetVisibilityLabel(s)
	return epduo
}

// SetNillableVisibilityLabel sets the visibility_label field if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableVisibilityLabel(s *string) *EquipmentPortDefinitionUpdateOne {
	if s != nil {
		epduo.SetVisibilityLabel(*s)
	}
	return epduo
}

// ClearVisibilityLabel clears the value of visibility_label.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearVisibilityLabel() *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.ClearVisibilityLabel()
	return epduo
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (epduo *EquipmentPortDefinitionUpdateOne) SetEquipmentPortTypeID(id int) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.SetEquipmentPortTypeID(id)
	return epduo
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableEquipmentPortTypeID(id *int) *EquipmentPortDefinitionUpdateOne {
	if id != nil {
		epduo = epduo.SetEquipmentPortTypeID(*id)
	}
	return epduo
}

// SetEquipmentPortType sets the equipment_port_type edge to EquipmentPortType.
func (epduo *EquipmentPortDefinitionUpdateOne) SetEquipmentPortType(e *EquipmentPortType) *EquipmentPortDefinitionUpdateOne {
	return epduo.SetEquipmentPortTypeID(e.ID)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (epduo *EquipmentPortDefinitionUpdateOne) AddPortIDs(ids ...int) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.AddPortIDs(ids...)
	return epduo
}

// AddPorts adds the ports edges to EquipmentPort.
func (epduo *EquipmentPortDefinitionUpdateOne) AddPorts(e ...*EquipmentPort) *EquipmentPortDefinitionUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epduo.AddPortIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epduo *EquipmentPortDefinitionUpdateOne) SetEquipmentTypeID(id int) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.SetEquipmentTypeID(id)
	return epduo
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epduo *EquipmentPortDefinitionUpdateOne) SetNillableEquipmentTypeID(id *int) *EquipmentPortDefinitionUpdateOne {
	if id != nil {
		epduo = epduo.SetEquipmentTypeID(*id)
	}
	return epduo
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (epduo *EquipmentPortDefinitionUpdateOne) SetEquipmentType(e *EquipmentType) *EquipmentPortDefinitionUpdateOne {
	return epduo.SetEquipmentTypeID(e.ID)
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearEquipmentPortType() *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.ClearEquipmentPortType()
	return epduo
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (epduo *EquipmentPortDefinitionUpdateOne) RemovePortIDs(ids ...int) *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.RemovePortIDs(ids...)
	return epduo
}

// RemovePorts removes ports edges to EquipmentPort.
func (epduo *EquipmentPortDefinitionUpdateOne) RemovePorts(e ...*EquipmentPort) *EquipmentPortDefinitionUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epduo.RemovePortIDs(ids...)
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (epduo *EquipmentPortDefinitionUpdateOne) ClearEquipmentType() *EquipmentPortDefinitionUpdateOne {
	epduo.mutation.ClearEquipmentType()
	return epduo
}

// Save executes the query and returns the updated entity.
func (epduo *EquipmentPortDefinitionUpdateOne) Save(ctx context.Context) (*EquipmentPortDefinition, error) {
	if _, ok := epduo.mutation.UpdateTime(); !ok {
		v := equipmentportdefinition.UpdateDefaultUpdateTime()
		epduo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *EquipmentPortDefinition
	)
	if len(epduo.hooks) == 0 {
		node, err = epduo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epduo.mutation = mutation
			node, err = epduo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(epduo.hooks) - 1; i >= 0; i-- {
			mut = epduo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epduo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (epduo *EquipmentPortDefinitionUpdateOne) SaveX(ctx context.Context) *EquipmentPortDefinition {
	epd, err := epduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return epd
}

// Exec executes the query on the entity.
func (epduo *EquipmentPortDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := epduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epduo *EquipmentPortDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := epduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epduo *EquipmentPortDefinitionUpdateOne) sqlSave(ctx context.Context) (epd *EquipmentPortDefinition, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentportdefinition.Table,
			Columns: equipmentportdefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentportdefinition.FieldID,
			},
		},
	}
	id, ok := epduo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing EquipmentPortDefinition.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := epduo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentportdefinition.FieldUpdateTime,
		})
	}
	if value, ok := epduo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldName,
		})
	}
	if value, ok := epduo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentportdefinition.FieldIndex,
		})
	}
	if value, ok := epduo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentportdefinition.FieldIndex,
		})
	}
	if epduo.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: equipmentportdefinition.FieldIndex,
		})
	}
	if value, ok := epduo.mutation.Bandwidth(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldBandwidth,
		})
	}
	if epduo.mutation.BandwidthCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipmentportdefinition.FieldBandwidth,
		})
	}
	if value, ok := epduo.mutation.VisibilityLabel(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldVisibilityLabel,
		})
	}
	if epduo.mutation.VisibilityLabelCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipmentportdefinition.FieldVisibilityLabel,
		})
	}
	if epduo.mutation.EquipmentPortTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentportdefinition.EquipmentPortTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epduo.mutation.EquipmentPortTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentportdefinition.EquipmentPortTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentporttype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := epduo.mutation.RemovedPortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentportdefinition.PortsTable,
			Columns: []string{equipmentportdefinition.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epduo.mutation.PortsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentportdefinition.PortsTable,
			Columns: []string{equipmentportdefinition.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentport.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epduo.mutation.EquipmentTypeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentportdefinition.EquipmentTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentTypeColumn},
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
	if nodes := epduo.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentportdefinition.EquipmentTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentTypeColumn},
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
	epd = &EquipmentPortDefinition{config: epduo.config}
	_spec.Assign = epd.assignValues
	_spec.ScanValues = epd.scanValues()
	if err = sqlgraph.UpdateNode(ctx, epduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentportdefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return epd, nil
}
