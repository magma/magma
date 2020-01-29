// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
)

// EquipmentPositionCreate is the builder for creating a EquipmentPosition entity.
type EquipmentPositionCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	definition  map[string]struct{}
	parent      map[string]struct{}
	attachment  map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (epc *EquipmentPositionCreate) SetCreateTime(t time.Time) *EquipmentPositionCreate {
	epc.create_time = &t
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
	epc.update_time = &t
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
func (epc *EquipmentPositionCreate) SetDefinitionID(id string) *EquipmentPositionCreate {
	if epc.definition == nil {
		epc.definition = make(map[string]struct{})
	}
	epc.definition[id] = struct{}{}
	return epc
}

// SetDefinition sets the definition edge to EquipmentPositionDefinition.
func (epc *EquipmentPositionCreate) SetDefinition(e *EquipmentPositionDefinition) *EquipmentPositionCreate {
	return epc.SetDefinitionID(e.ID)
}

// SetParentID sets the parent edge to Equipment by id.
func (epc *EquipmentPositionCreate) SetParentID(id string) *EquipmentPositionCreate {
	if epc.parent == nil {
		epc.parent = make(map[string]struct{})
	}
	epc.parent[id] = struct{}{}
	return epc
}

// SetNillableParentID sets the parent edge to Equipment by id if the given value is not nil.
func (epc *EquipmentPositionCreate) SetNillableParentID(id *string) *EquipmentPositionCreate {
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
func (epc *EquipmentPositionCreate) SetAttachmentID(id string) *EquipmentPositionCreate {
	if epc.attachment == nil {
		epc.attachment = make(map[string]struct{})
	}
	epc.attachment[id] = struct{}{}
	return epc
}

// SetNillableAttachmentID sets the attachment edge to Equipment by id if the given value is not nil.
func (epc *EquipmentPositionCreate) SetNillableAttachmentID(id *string) *EquipmentPositionCreate {
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
	if epc.create_time == nil {
		v := equipmentposition.DefaultCreateTime()
		epc.create_time = &v
	}
	if epc.update_time == nil {
		v := equipmentposition.DefaultUpdateTime()
		epc.update_time = &v
	}
	if len(epc.definition) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"definition\"")
	}
	if epc.definition == nil {
		return nil, errors.New("ent: missing required edge \"definition\"")
	}
	if len(epc.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	if len(epc.attachment) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"attachment\"")
	}
	return epc.sqlSave(ctx)
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
				Type:   field.TypeString,
				Column: equipmentposition.FieldID,
			},
		}
	)
	if value := epc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentposition.FieldCreateTime,
		})
		ep.CreateTime = *value
	}
	if value := epc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentposition.FieldUpdateTime,
		})
		ep.UpdateTime = *value
	}
	if nodes := epc.definition; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentposition.DefinitionTable,
			Columns: []string{equipmentposition.DefinitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentpositiondefinition.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := epc.parent; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentposition.ParentTable,
			Columns: []string{equipmentposition.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := epc.attachment; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   equipmentposition.AttachmentTable,
			Columns: []string{equipmentposition.AttachmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipment.FieldID,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, epc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	ep.ID = strconv.FormatInt(id, 10)
	return ep, nil
}
