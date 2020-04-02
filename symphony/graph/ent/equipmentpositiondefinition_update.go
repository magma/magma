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
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPositionDefinitionUpdate is the builder for updating EquipmentPositionDefinition entities.
type EquipmentPositionDefinitionUpdate struct {
	config
	hooks      []Hook
	mutation   *EquipmentPositionDefinitionMutation
	predicates []predicate.EquipmentPositionDefinition
}

// Where adds a new predicate for the builder.
func (epdu *EquipmentPositionDefinitionUpdate) Where(ps ...predicate.EquipmentPositionDefinition) *EquipmentPositionDefinitionUpdate {
	epdu.predicates = append(epdu.predicates, ps...)
	return epdu
}

// SetName sets the name field.
func (epdu *EquipmentPositionDefinitionUpdate) SetName(s string) *EquipmentPositionDefinitionUpdate {
	epdu.mutation.SetName(s)
	return epdu
}

// SetIndex sets the index field.
func (epdu *EquipmentPositionDefinitionUpdate) SetIndex(i int) *EquipmentPositionDefinitionUpdate {
	epdu.mutation.ResetIndex()
	epdu.mutation.SetIndex(i)
	return epdu
}

// SetNillableIndex sets the index field if the given value is not nil.
func (epdu *EquipmentPositionDefinitionUpdate) SetNillableIndex(i *int) *EquipmentPositionDefinitionUpdate {
	if i != nil {
		epdu.SetIndex(*i)
	}
	return epdu
}

// AddIndex adds i to index.
func (epdu *EquipmentPositionDefinitionUpdate) AddIndex(i int) *EquipmentPositionDefinitionUpdate {
	epdu.mutation.AddIndex(i)
	return epdu
}

// ClearIndex clears the value of index.
func (epdu *EquipmentPositionDefinitionUpdate) ClearIndex() *EquipmentPositionDefinitionUpdate {
	epdu.mutation.ClearIndex()
	return epdu
}

// SetVisibilityLabel sets the visibility_label field.
func (epdu *EquipmentPositionDefinitionUpdate) SetVisibilityLabel(s string) *EquipmentPositionDefinitionUpdate {
	epdu.mutation.SetVisibilityLabel(s)
	return epdu
}

// SetNillableVisibilityLabel sets the visibility_label field if the given value is not nil.
func (epdu *EquipmentPositionDefinitionUpdate) SetNillableVisibilityLabel(s *string) *EquipmentPositionDefinitionUpdate {
	if s != nil {
		epdu.SetVisibilityLabel(*s)
	}
	return epdu
}

// ClearVisibilityLabel clears the value of visibility_label.
func (epdu *EquipmentPositionDefinitionUpdate) ClearVisibilityLabel() *EquipmentPositionDefinitionUpdate {
	epdu.mutation.ClearVisibilityLabel()
	return epdu
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (epdu *EquipmentPositionDefinitionUpdate) AddPositionIDs(ids ...int) *EquipmentPositionDefinitionUpdate {
	epdu.mutation.AddPositionIDs(ids...)
	return epdu
}

// AddPositions adds the positions edges to EquipmentPosition.
func (epdu *EquipmentPositionDefinitionUpdate) AddPositions(e ...*EquipmentPosition) *EquipmentPositionDefinitionUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdu.AddPositionIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epdu *EquipmentPositionDefinitionUpdate) SetEquipmentTypeID(id int) *EquipmentPositionDefinitionUpdate {
	epdu.mutation.SetEquipmentTypeID(id)
	return epdu
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epdu *EquipmentPositionDefinitionUpdate) SetNillableEquipmentTypeID(id *int) *EquipmentPositionDefinitionUpdate {
	if id != nil {
		epdu = epdu.SetEquipmentTypeID(*id)
	}
	return epdu
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (epdu *EquipmentPositionDefinitionUpdate) SetEquipmentType(e *EquipmentType) *EquipmentPositionDefinitionUpdate {
	return epdu.SetEquipmentTypeID(e.ID)
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (epdu *EquipmentPositionDefinitionUpdate) RemovePositionIDs(ids ...int) *EquipmentPositionDefinitionUpdate {
	epdu.mutation.RemovePositionIDs(ids...)
	return epdu
}

// RemovePositions removes positions edges to EquipmentPosition.
func (epdu *EquipmentPositionDefinitionUpdate) RemovePositions(e ...*EquipmentPosition) *EquipmentPositionDefinitionUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdu.RemovePositionIDs(ids...)
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (epdu *EquipmentPositionDefinitionUpdate) ClearEquipmentType() *EquipmentPositionDefinitionUpdate {
	epdu.mutation.ClearEquipmentType()
	return epdu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (epdu *EquipmentPositionDefinitionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := epdu.mutation.UpdateTime(); !ok {
		v := equipmentpositiondefinition.UpdateDefaultUpdateTime()
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
			mutation, ok := m.(*EquipmentPositionDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epdu.mutation = mutation
			affected, err = epdu.sqlSave(ctx)
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
func (epdu *EquipmentPositionDefinitionUpdate) SaveX(ctx context.Context) int {
	affected, err := epdu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (epdu *EquipmentPositionDefinitionUpdate) Exec(ctx context.Context) error {
	_, err := epdu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epdu *EquipmentPositionDefinitionUpdate) ExecX(ctx context.Context) {
	if err := epdu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epdu *EquipmentPositionDefinitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentpositiondefinition.Table,
			Columns: equipmentpositiondefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentpositiondefinition.FieldID,
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
			Column: equipmentpositiondefinition.FieldUpdateTime,
		})
	}
	if value, ok := epdu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentpositiondefinition.FieldName,
		})
	}
	if value, ok := epdu.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentpositiondefinition.FieldIndex,
		})
	}
	if value, ok := epdu.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentpositiondefinition.FieldIndex,
		})
	}
	if epdu.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: equipmentpositiondefinition.FieldIndex,
		})
	}
	if value, ok := epdu.mutation.VisibilityLabel(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentpositiondefinition.FieldVisibilityLabel,
		})
	}
	if epdu.mutation.VisibilityLabelCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipmentpositiondefinition.FieldVisibilityLabel,
		})
	}
	if nodes := epdu.mutation.RemovedPositionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentpositiondefinition.PositionsTable,
			Columns: []string{equipmentpositiondefinition.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epdu.mutation.PositionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentpositiondefinition.PositionsTable,
			Columns: []string{equipmentpositiondefinition.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
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
			Table:   equipmentpositiondefinition.EquipmentTypeTable,
			Columns: []string{equipmentpositiondefinition.EquipmentTypeColumn},
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
			Table:   equipmentpositiondefinition.EquipmentTypeTable,
			Columns: []string{equipmentpositiondefinition.EquipmentTypeColumn},
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
			err = &NotFoundError{equipmentpositiondefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentPositionDefinitionUpdateOne is the builder for updating a single EquipmentPositionDefinition entity.
type EquipmentPositionDefinitionUpdateOne struct {
	config
	hooks    []Hook
	mutation *EquipmentPositionDefinitionMutation
}

// SetName sets the name field.
func (epduo *EquipmentPositionDefinitionUpdateOne) SetName(s string) *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.SetName(s)
	return epduo
}

// SetIndex sets the index field.
func (epduo *EquipmentPositionDefinitionUpdateOne) SetIndex(i int) *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.ResetIndex()
	epduo.mutation.SetIndex(i)
	return epduo
}

// SetNillableIndex sets the index field if the given value is not nil.
func (epduo *EquipmentPositionDefinitionUpdateOne) SetNillableIndex(i *int) *EquipmentPositionDefinitionUpdateOne {
	if i != nil {
		epduo.SetIndex(*i)
	}
	return epduo
}

// AddIndex adds i to index.
func (epduo *EquipmentPositionDefinitionUpdateOne) AddIndex(i int) *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.AddIndex(i)
	return epduo
}

// ClearIndex clears the value of index.
func (epduo *EquipmentPositionDefinitionUpdateOne) ClearIndex() *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.ClearIndex()
	return epduo
}

// SetVisibilityLabel sets the visibility_label field.
func (epduo *EquipmentPositionDefinitionUpdateOne) SetVisibilityLabel(s string) *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.SetVisibilityLabel(s)
	return epduo
}

// SetNillableVisibilityLabel sets the visibility_label field if the given value is not nil.
func (epduo *EquipmentPositionDefinitionUpdateOne) SetNillableVisibilityLabel(s *string) *EquipmentPositionDefinitionUpdateOne {
	if s != nil {
		epduo.SetVisibilityLabel(*s)
	}
	return epduo
}

// ClearVisibilityLabel clears the value of visibility_label.
func (epduo *EquipmentPositionDefinitionUpdateOne) ClearVisibilityLabel() *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.ClearVisibilityLabel()
	return epduo
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (epduo *EquipmentPositionDefinitionUpdateOne) AddPositionIDs(ids ...int) *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.AddPositionIDs(ids...)
	return epduo
}

// AddPositions adds the positions edges to EquipmentPosition.
func (epduo *EquipmentPositionDefinitionUpdateOne) AddPositions(e ...*EquipmentPosition) *EquipmentPositionDefinitionUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epduo.AddPositionIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epduo *EquipmentPositionDefinitionUpdateOne) SetEquipmentTypeID(id int) *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.SetEquipmentTypeID(id)
	return epduo
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epduo *EquipmentPositionDefinitionUpdateOne) SetNillableEquipmentTypeID(id *int) *EquipmentPositionDefinitionUpdateOne {
	if id != nil {
		epduo = epduo.SetEquipmentTypeID(*id)
	}
	return epduo
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (epduo *EquipmentPositionDefinitionUpdateOne) SetEquipmentType(e *EquipmentType) *EquipmentPositionDefinitionUpdateOne {
	return epduo.SetEquipmentTypeID(e.ID)
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (epduo *EquipmentPositionDefinitionUpdateOne) RemovePositionIDs(ids ...int) *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.RemovePositionIDs(ids...)
	return epduo
}

// RemovePositions removes positions edges to EquipmentPosition.
func (epduo *EquipmentPositionDefinitionUpdateOne) RemovePositions(e ...*EquipmentPosition) *EquipmentPositionDefinitionUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epduo.RemovePositionIDs(ids...)
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (epduo *EquipmentPositionDefinitionUpdateOne) ClearEquipmentType() *EquipmentPositionDefinitionUpdateOne {
	epduo.mutation.ClearEquipmentType()
	return epduo
}

// Save executes the query and returns the updated entity.
func (epduo *EquipmentPositionDefinitionUpdateOne) Save(ctx context.Context) (*EquipmentPositionDefinition, error) {
	if _, ok := epduo.mutation.UpdateTime(); !ok {
		v := equipmentpositiondefinition.UpdateDefaultUpdateTime()
		epduo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *EquipmentPositionDefinition
	)
	if len(epduo.hooks) == 0 {
		node, err = epduo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPositionDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epduo.mutation = mutation
			node, err = epduo.sqlSave(ctx)
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
func (epduo *EquipmentPositionDefinitionUpdateOne) SaveX(ctx context.Context) *EquipmentPositionDefinition {
	epd, err := epduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return epd
}

// Exec executes the query on the entity.
func (epduo *EquipmentPositionDefinitionUpdateOne) Exec(ctx context.Context) error {
	_, err := epduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epduo *EquipmentPositionDefinitionUpdateOne) ExecX(ctx context.Context) {
	if err := epduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epduo *EquipmentPositionDefinitionUpdateOne) sqlSave(ctx context.Context) (epd *EquipmentPositionDefinition, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentpositiondefinition.Table,
			Columns: equipmentpositiondefinition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentpositiondefinition.FieldID,
			},
		},
	}
	id, ok := epduo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing EquipmentPositionDefinition.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := epduo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentpositiondefinition.FieldUpdateTime,
		})
	}
	if value, ok := epduo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentpositiondefinition.FieldName,
		})
	}
	if value, ok := epduo.mutation.Index(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentpositiondefinition.FieldIndex,
		})
	}
	if value, ok := epduo.mutation.AddedIndex(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentpositiondefinition.FieldIndex,
		})
	}
	if epduo.mutation.IndexCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: equipmentpositiondefinition.FieldIndex,
		})
	}
	if value, ok := epduo.mutation.VisibilityLabel(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentpositiondefinition.FieldVisibilityLabel,
		})
	}
	if epduo.mutation.VisibilityLabelCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: equipmentpositiondefinition.FieldVisibilityLabel,
		})
	}
	if nodes := epduo.mutation.RemovedPositionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentpositiondefinition.PositionsTable,
			Columns: []string{equipmentpositiondefinition.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epduo.mutation.PositionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentpositiondefinition.PositionsTable,
			Columns: []string{equipmentpositiondefinition.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentposition.FieldID,
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
			Table:   equipmentpositiondefinition.EquipmentTypeTable,
			Columns: []string{equipmentpositiondefinition.EquipmentTypeColumn},
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
			Table:   equipmentpositiondefinition.EquipmentTypeTable,
			Columns: []string{equipmentpositiondefinition.EquipmentTypeColumn},
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
	epd = &EquipmentPositionDefinition{config: epduo.config}
	_spec.Assign = epd.assignValues
	_spec.ScanValues = epd.scanValues()
	if err = sqlgraph.UpdateNode(ctx, epduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentpositiondefinition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return epd, nil
}
