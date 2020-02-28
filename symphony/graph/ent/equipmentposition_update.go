// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"time"

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

	update_time       *time.Time
	definition        map[int]struct{}
	parent            map[int]struct{}
	attachment        map[int]struct{}
	clearedDefinition bool
	clearedParent     bool
	clearedAttachment bool
	predicates        []predicate.EquipmentPosition
}

// Where adds a new predicate for the builder.
func (epu *EquipmentPositionUpdate) Where(ps ...predicate.EquipmentPosition) *EquipmentPositionUpdate {
	epu.predicates = append(epu.predicates, ps...)
	return epu
}

// SetDefinitionID sets the definition edge to EquipmentPositionDefinition by id.
func (epu *EquipmentPositionUpdate) SetDefinitionID(id int) *EquipmentPositionUpdate {
	if epu.definition == nil {
		epu.definition = make(map[int]struct{})
	}
	epu.definition[id] = struct{}{}
	return epu
}

// SetDefinition sets the definition edge to EquipmentPositionDefinition.
func (epu *EquipmentPositionUpdate) SetDefinition(e *EquipmentPositionDefinition) *EquipmentPositionUpdate {
	return epu.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epu *EquipmentPositionUpdate) SetParentID(id int) *EquipmentPositionUpdate {
	if epu.parent == nil {
		epu.parent = make(map[int]struct{})
	}
	epu.parent[id] = struct{}{}
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
	if epu.attachment == nil {
		epu.attachment = make(map[int]struct{})
	}
	epu.attachment[id] = struct{}{}
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
	epu.clearedDefinition = true
	return epu
}

// ClearParent clears the parent edge to Equipment.
func (epu *EquipmentPositionUpdate) ClearParent() *EquipmentPositionUpdate {
	epu.clearedParent = true
	return epu
}

// ClearAttachment clears the attachment edge to Equipment.
func (epu *EquipmentPositionUpdate) ClearAttachment() *EquipmentPositionUpdate {
	epu.clearedAttachment = true
	return epu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (epu *EquipmentPositionUpdate) Save(ctx context.Context) (int, error) {
	if epu.update_time == nil {
		v := equipmentposition.UpdateDefaultUpdateTime()
		epu.update_time = &v
	}
	if len(epu.definition) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"definition\"")
	}
	if epu.clearedDefinition && epu.definition == nil {
		return 0, errors.New("ent: clearing a unique edge \"definition\"")
	}
	if len(epu.parent) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	if len(epu.attachment) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"attachment\"")
	}
	return epu.sqlSave(ctx)
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
	if value := epu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentposition.FieldUpdateTime,
		})
	}
	if epu.clearedDefinition {
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
	if nodes := epu.definition; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epu.clearedParent {
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
	if nodes := epu.parent; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epu.clearedAttachment {
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
	if nodes := epu.attachment; len(nodes) > 0 {
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
		for k, _ := range nodes {
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
	id int

	update_time       *time.Time
	definition        map[int]struct{}
	parent            map[int]struct{}
	attachment        map[int]struct{}
	clearedDefinition bool
	clearedParent     bool
	clearedAttachment bool
}

// SetDefinitionID sets the definition edge to EquipmentPositionDefinition by id.
func (epuo *EquipmentPositionUpdateOne) SetDefinitionID(id int) *EquipmentPositionUpdateOne {
	if epuo.definition == nil {
		epuo.definition = make(map[int]struct{})
	}
	epuo.definition[id] = struct{}{}
	return epuo
}

// SetDefinition sets the definition edge to EquipmentPositionDefinition.
func (epuo *EquipmentPositionUpdateOne) SetDefinition(e *EquipmentPositionDefinition) *EquipmentPositionUpdateOne {
	return epuo.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epuo *EquipmentPositionUpdateOne) SetParentID(id int) *EquipmentPositionUpdateOne {
	if epuo.parent == nil {
		epuo.parent = make(map[int]struct{})
	}
	epuo.parent[id] = struct{}{}
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
	if epuo.attachment == nil {
		epuo.attachment = make(map[int]struct{})
	}
	epuo.attachment[id] = struct{}{}
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
	epuo.clearedDefinition = true
	return epuo
}

// ClearParent clears the parent edge to Equipment.
func (epuo *EquipmentPositionUpdateOne) ClearParent() *EquipmentPositionUpdateOne {
	epuo.clearedParent = true
	return epuo
}

// ClearAttachment clears the attachment edge to Equipment.
func (epuo *EquipmentPositionUpdateOne) ClearAttachment() *EquipmentPositionUpdateOne {
	epuo.clearedAttachment = true
	return epuo
}

// Save executes the query and returns the updated entity.
func (epuo *EquipmentPositionUpdateOne) Save(ctx context.Context) (*EquipmentPosition, error) {
	if epuo.update_time == nil {
		v := equipmentposition.UpdateDefaultUpdateTime()
		epuo.update_time = &v
	}
	if len(epuo.definition) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"definition\"")
	}
	if epuo.clearedDefinition && epuo.definition == nil {
		return nil, errors.New("ent: clearing a unique edge \"definition\"")
	}
	if len(epuo.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	if len(epuo.attachment) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"attachment\"")
	}
	return epuo.sqlSave(ctx)
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
				Value:  epuo.id,
				Type:   field.TypeInt,
				Column: equipmentposition.FieldID,
			},
		},
	}
	if value := epuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentposition.FieldUpdateTime,
		})
	}
	if epuo.clearedDefinition {
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
	if nodes := epuo.definition; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epuo.clearedParent {
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
	if nodes := epuo.parent; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if epuo.clearedAttachment {
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
	if nodes := epuo.attachment; len(nodes) > 0 {
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
		for k, _ := range nodes {
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
