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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// EquipmentPositionUpdate is the builder for updating EquipmentPosition entities.
type EquipmentPositionUpdate struct {
	config
	hooks      []Hook
	mutation   *EquipmentPositionMutation
	predicates []predicate.EquipmentPosition
}

// Where adds a new predicate for the builder.
func (epu *EquipmentPositionUpdate) Where(ps ...predicate.EquipmentPosition) *EquipmentPositionUpdate {
	epu.predicates = append(epu.predicates, ps...)
	return epu
}

// SetDefinitionID sets the definition edge to EquipmentPositionDefinition by id.
func (epu *EquipmentPositionUpdate) SetDefinitionID(id int) *EquipmentPositionUpdate {
	epu.mutation.SetDefinitionID(id)
	return epu
}

// SetDefinition sets the definition edge to EquipmentPositionDefinition.
func (epu *EquipmentPositionUpdate) SetDefinition(e *EquipmentPositionDefinition) *EquipmentPositionUpdate {
	return epu.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epu *EquipmentPositionUpdate) SetParentID(id int) *EquipmentPositionUpdate {
	epu.mutation.SetParentID(id)
	return epu
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epu *EquipmentPositionUpdate) SetNillableParentID(id *int) *EquipmentPositionUpdate {
	if id != nil {
		epu = epu.SetParentID(*id)
	}
	return epu
}

// SetParent sets the parent edge to Equipment.
func (epu *EquipmentPositionUpdate) SetParent(e *Equipment) *EquipmentPositionUpdate {
	return epu.SetParentID(e.ID)
}

// SetAttachmentID sets the attachment edge to Equipment by id.
func (epu *EquipmentPositionUpdate) SetAttachmentID(id int) *EquipmentPositionUpdate {
	epu.mutation.SetAttachmentID(id)
	return epu
}

// SetNillableAttachmentID sets the attachment edge to Equipment by id if the given value is not nil.
func (epu *EquipmentPositionUpdate) SetNillableAttachmentID(id *int) *EquipmentPositionUpdate {
	if id != nil {
		epu = epu.SetAttachmentID(*id)
	}
	return epu
}

// SetAttachment sets the attachment edge to Equipment.
func (epu *EquipmentPositionUpdate) SetAttachment(e *Equipment) *EquipmentPositionUpdate {
	return epu.SetAttachmentID(e.ID)
}

// ClearDefinition clears the definition edge to EquipmentPositionDefinition.
func (epu *EquipmentPositionUpdate) ClearDefinition() *EquipmentPositionUpdate {
	epu.mutation.ClearDefinition()
	return epu
}

// ClearParent clears the parent edge to Equipment.
func (epu *EquipmentPositionUpdate) ClearParent() *EquipmentPositionUpdate {
	epu.mutation.ClearParent()
	return epu
}

// ClearAttachment clears the attachment edge to Equipment.
func (epu *EquipmentPositionUpdate) ClearAttachment() *EquipmentPositionUpdate {
	epu.mutation.ClearAttachment()
	return epu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (epu *EquipmentPositionUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := epu.mutation.UpdateTime(); !ok {
		v := equipmentposition.UpdateDefaultUpdateTime()
		epu.mutation.SetUpdateTime(v)
	}

	if _, ok := epu.mutation.DefinitionID(); epu.mutation.DefinitionCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"definition\"")
	}

	var (
		err      error
		affected int
	)
	if len(epu.hooks) == 0 {
		affected, err = epu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPositionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epu.mutation = mutation
			affected, err = epu.sqlSave(ctx)
			return affected, err
		})
		for i := len(epu.hooks) - 1; i >= 0; i-- {
			mut = epu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (epu *EquipmentPositionUpdate) SaveX(ctx context.Context) int {
	affected, err := epu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (epu *EquipmentPositionUpdate) Exec(ctx context.Context) error {
	_, err := epu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epu *EquipmentPositionUpdate) ExecX(ctx context.Context) {
	if err := epu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epu *EquipmentPositionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentposition.Table,
			Columns: equipmentposition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentposition.FieldID,
			},
		},
	}
	if ps := epu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := epu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentposition.FieldUpdateTime,
		})
	}
	if epu.mutation.DefinitionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentposition.DefinitionTable,
			Columns: []string{equipmentposition.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentpositiondefinition.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.DefinitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentposition.DefinitionTable,
			Columns: []string{equipmentposition.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentpositiondefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epu.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentposition.ParentTable,
			Columns: []string{equipmentposition.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentposition.ParentTable,
			Columns: []string{equipmentposition.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epu.mutation.AttachmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   equipmentposition.AttachmentTable,
			Columns: []string{equipmentposition.AttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epu.mutation.AttachmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   equipmentposition.AttachmentTable,
			Columns: []string{equipmentposition.AttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, epu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentposition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// EquipmentPositionUpdateOne is the builder for updating a single EquipmentPosition entity.
type EquipmentPositionUpdateOne struct {
	config
	hooks    []Hook
	mutation *EquipmentPositionMutation
}

// SetDefinitionID sets the definition edge to EquipmentPositionDefinition by id.
func (epuo *EquipmentPositionUpdateOne) SetDefinitionID(id int) *EquipmentPositionUpdateOne {
	epuo.mutation.SetDefinitionID(id)
	return epuo
}

// SetDefinition sets the definition edge to EquipmentPositionDefinition.
func (epuo *EquipmentPositionUpdateOne) SetDefinition(e *EquipmentPositionDefinition) *EquipmentPositionUpdateOne {
	return epuo.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epuo *EquipmentPositionUpdateOne) SetParentID(id int) *EquipmentPositionUpdateOne {
	epuo.mutation.SetParentID(id)
	return epuo
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epuo *EquipmentPositionUpdateOne) SetNillableParentID(id *int) *EquipmentPositionUpdateOne {
	if id != nil {
		epuo = epuo.SetParentID(*id)
	}
	return epuo
}

// SetParent sets the parent edge to Equipment.
func (epuo *EquipmentPositionUpdateOne) SetParent(e *Equipment) *EquipmentPositionUpdateOne {
	return epuo.SetParentID(e.ID)
}

// SetAttachmentID sets the attachment edge to Equipment by id.
func (epuo *EquipmentPositionUpdateOne) SetAttachmentID(id int) *EquipmentPositionUpdateOne {
	epuo.mutation.SetAttachmentID(id)
	return epuo
}

// SetNillableAttachmentID sets the attachment edge to Equipment by id if the given value is not nil.
func (epuo *EquipmentPositionUpdateOne) SetNillableAttachmentID(id *int) *EquipmentPositionUpdateOne {
	if id != nil {
		epuo = epuo.SetAttachmentID(*id)
	}
	return epuo
}

// SetAttachment sets the attachment edge to Equipment.
func (epuo *EquipmentPositionUpdateOne) SetAttachment(e *Equipment) *EquipmentPositionUpdateOne {
	return epuo.SetAttachmentID(e.ID)
}

// ClearDefinition clears the definition edge to EquipmentPositionDefinition.
func (epuo *EquipmentPositionUpdateOne) ClearDefinition() *EquipmentPositionUpdateOne {
	epuo.mutation.ClearDefinition()
	return epuo
}

// ClearParent clears the parent edge to Equipment.
func (epuo *EquipmentPositionUpdateOne) ClearParent() *EquipmentPositionUpdateOne {
	epuo.mutation.ClearParent()
	return epuo
}

// ClearAttachment clears the attachment edge to Equipment.
func (epuo *EquipmentPositionUpdateOne) ClearAttachment() *EquipmentPositionUpdateOne {
	epuo.mutation.ClearAttachment()
	return epuo
}

// Save executes the query and returns the updated entity.
func (epuo *EquipmentPositionUpdateOne) Save(ctx context.Context) (*EquipmentPosition, error) {
	if _, ok := epuo.mutation.UpdateTime(); !ok {
		v := equipmentposition.UpdateDefaultUpdateTime()
		epuo.mutation.SetUpdateTime(v)
	}

	if _, ok := epuo.mutation.DefinitionID(); epuo.mutation.DefinitionCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"definition\"")
	}

	var (
		err  error
		node *EquipmentPosition
	)
	if len(epuo.hooks) == 0 {
		node, err = epuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPositionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epuo.mutation = mutation
			node, err = epuo.sqlSave(ctx)
			return node, err
		})
		for i := len(epuo.hooks) - 1; i >= 0; i-- {
			mut = epuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (epuo *EquipmentPositionUpdateOne) SaveX(ctx context.Context) *EquipmentPosition {
	ep, err := epuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ep
}

// Exec executes the query on the entity.
func (epuo *EquipmentPositionUpdateOne) Exec(ctx context.Context) error {
	_, err := epuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (epuo *EquipmentPositionUpdateOne) ExecX(ctx context.Context) {
	if err := epuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (epuo *EquipmentPositionUpdateOne) sqlSave(ctx context.Context) (ep *EquipmentPosition, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   equipmentposition.Table,
			Columns: equipmentposition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentposition.FieldID,
			},
		},
	}
	id, ok := epuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing EquipmentPosition.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := epuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentposition.FieldUpdateTime,
		})
	}
	if epuo.mutation.DefinitionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentposition.DefinitionTable,
			Columns: []string{equipmentposition.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentpositiondefinition.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.DefinitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentposition.DefinitionTable,
			Columns: []string{equipmentposition.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipmentpositiondefinition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epuo.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentposition.ParentTable,
			Columns: []string{equipmentposition.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentposition.ParentTable,
			Columns: []string{equipmentposition.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epuo.mutation.AttachmentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   equipmentposition.AttachmentTable,
			Columns: []string{equipmentposition.AttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := epuo.mutation.AttachmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   equipmentposition.AttachmentTable,
			Columns: []string{equipmentposition.AttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: equipment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	ep = &EquipmentPosition{config: epuo.config}
	_spec.Assign = ep.assignValues
	_spec.ScanValues = ep.scanValues()
	if err = sqlgraph.UpdateNode(ctx, epuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{equipmentposition.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ep, nil
}
