// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentposition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentpositiondefinition"
)

// EquipmentPositionCreate is the builder for creating a EquipmentPosition entity.
type EquipmentPositionCreate struct {
	config
	mutation *EquipmentPositionMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (epc *EquipmentPositionCreate) SetCreateTime(t time.Time) *EquipmentPositionCreate {
	epc.mutation.SetCreateTime(t)
	return epc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (epc *EquipmentPositionCreate) SetNillableCreateTime(t *time.Time) *EquipmentPositionCreate {
	if t != nil {
		epc.SetCreateTime(*t)
	}
	return epc
}

// SetUpdateTime sets the update_time field.
func (epc *EquipmentPositionCreate) SetUpdateTime(t time.Time) *EquipmentPositionCreate {
	epc.mutation.SetUpdateTime(t)
	return epc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (epc *EquipmentPositionCreate) SetNillableUpdateTime(t *time.Time) *EquipmentPositionCreate {
	if t != nil {
		epc.SetUpdateTime(*t)
	}
	return epc
}

// SetDefinitionID sets the definition edge to EquipmentPositionDefinition by id.
func (epc *EquipmentPositionCreate) SetDefinitionID(id int) *EquipmentPositionCreate {
	epc.mutation.SetDefinitionID(id)
	return epc
}

// SetDefinition sets the definition edge to EquipmentPositionDefinition.
func (epc *EquipmentPositionCreate) SetDefinition(e *EquipmentPositionDefinition) *EquipmentPositionCreate {
	return epc.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epc *EquipmentPositionCreate) SetParentID(id int) *EquipmentPositionCreate {
	epc.mutation.SetParentID(id)
	return epc
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epc *EquipmentPositionCreate) SetNillableParentID(id *int) *EquipmentPositionCreate {
	if id != nil {
		epc = epc.SetParentID(*id)
	}
	return epc
}

// SetParent sets the parent edge to Equipment.
func (epc *EquipmentPositionCreate) SetParent(e *Equipment) *EquipmentPositionCreate {
	return epc.SetParentID(e.ID)
}

// SetAttachmentID sets the attachment edge to Equipment by id.
func (epc *EquipmentPositionCreate) SetAttachmentID(id int) *EquipmentPositionCreate {
	epc.mutation.SetAttachmentID(id)
	return epc
}

// SetNillableAttachmentID sets the attachment edge to Equipment by id if the given value is not nil.
func (epc *EquipmentPositionCreate) SetNillableAttachmentID(id *int) *EquipmentPositionCreate {
	if id != nil {
		epc = epc.SetAttachmentID(*id)
	}
	return epc
}

// SetAttachment sets the attachment edge to Equipment.
func (epc *EquipmentPositionCreate) SetAttachment(e *Equipment) *EquipmentPositionCreate {
	return epc.SetAttachmentID(e.ID)
}

// Save creates the EquipmentPosition in the database.
func (epc *EquipmentPositionCreate) Save(ctx context.Context) (*EquipmentPosition, error) {
	if _, ok := epc.mutation.CreateTime(); !ok {
		v := equipmentposition.DefaultCreateTime()
		epc.mutation.SetCreateTime(v)
	}
	if _, ok := epc.mutation.UpdateTime(); !ok {
		v := equipmentposition.DefaultUpdateTime()
		epc.mutation.SetUpdateTime(v)
	}
	if _, ok := epc.mutation.DefinitionID(); !ok {
		return nil, errors.New("ent: missing required edge \"definition\"")
	}
	var (
		err  error
		node *EquipmentPosition
	)
	if len(epc.hooks) == 0 {
		node, err = epc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPositionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epc.mutation = mutation
			node, err = epc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(epc.hooks) - 1; i >= 0; i-- {
			mut = epc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (epc *EquipmentPositionCreate) SaveX(ctx context.Context) *EquipmentPosition {
	v, err := epc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epc *EquipmentPositionCreate) sqlSave(ctx context.Context) (*EquipmentPosition, error) {
	var (
		ep    = &EquipmentPosition{config: epc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: equipmentposition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentposition.FieldID,
			},
		}
	)
	if value, ok := epc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentposition.FieldCreateTime,
		})
		ep.CreateTime = value
	}
	if value, ok := epc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentposition.FieldUpdateTime,
		})
		ep.UpdateTime = value
	}
	if nodes := epc.mutation.DefinitionIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := epc.mutation.ParentIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := epc.mutation.AttachmentIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, epc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	ep.ID = int(id)
	return ep, nil
}
