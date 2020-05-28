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
	"github.com/facebookincubator/symphony/pkg/ent/equipmentposition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"
)

// EquipmentPositionDefinitionCreate is the builder for creating a EquipmentPositionDefinition entity.
type EquipmentPositionDefinitionCreate struct {
	config
	mutation *EquipmentPositionDefinitionMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (epdc *EquipmentPositionDefinitionCreate) SetCreateTime(t time.Time) *EquipmentPositionDefinitionCreate {
	epdc.mutation.SetCreateTime(t)
	return epdc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (epdc *EquipmentPositionDefinitionCreate) SetNillableCreateTime(t *time.Time) *EquipmentPositionDefinitionCreate {
	if t != nil {
		epdc.SetCreateTime(*t)
	}
	return epdc
}

// SetUpdateTime sets the update_time field.
func (epdc *EquipmentPositionDefinitionCreate) SetUpdateTime(t time.Time) *EquipmentPositionDefinitionCreate {
	epdc.mutation.SetUpdateTime(t)
	return epdc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (epdc *EquipmentPositionDefinitionCreate) SetNillableUpdateTime(t *time.Time) *EquipmentPositionDefinitionCreate {
	if t != nil {
		epdc.SetUpdateTime(*t)
	}
	return epdc
}

// SetName sets the name field.
func (epdc *EquipmentPositionDefinitionCreate) SetName(s string) *EquipmentPositionDefinitionCreate {
	epdc.mutation.SetName(s)
	return epdc
}

// SetIndex sets the index field.
func (epdc *EquipmentPositionDefinitionCreate) SetIndex(i int) *EquipmentPositionDefinitionCreate {
	epdc.mutation.SetIndex(i)
	return epdc
}

// SetNillableIndex sets the index field if the given value is not nil.
func (epdc *EquipmentPositionDefinitionCreate) SetNillableIndex(i *int) *EquipmentPositionDefinitionCreate {
	if i != nil {
		epdc.SetIndex(*i)
	}
	return epdc
}

// SetVisibilityLabel sets the visibility_label field.
func (epdc *EquipmentPositionDefinitionCreate) SetVisibilityLabel(s string) *EquipmentPositionDefinitionCreate {
	epdc.mutation.SetVisibilityLabel(s)
	return epdc
}

// SetNillableVisibilityLabel sets the visibility_label field if the given value is not nil.
func (epdc *EquipmentPositionDefinitionCreate) SetNillableVisibilityLabel(s *string) *EquipmentPositionDefinitionCreate {
	if s != nil {
		epdc.SetVisibilityLabel(*s)
	}
	return epdc
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (epdc *EquipmentPositionDefinitionCreate) AddPositionIDs(ids ...int) *EquipmentPositionDefinitionCreate {
	epdc.mutation.AddPositionIDs(ids...)
	return epdc
}

// AddPositions adds the positions edges to EquipmentPosition.
func (epdc *EquipmentPositionDefinitionCreate) AddPositions(e ...*EquipmentPosition) *EquipmentPositionDefinitionCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdc.AddPositionIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epdc *EquipmentPositionDefinitionCreate) SetEquipmentTypeID(id int) *EquipmentPositionDefinitionCreate {
	epdc.mutation.SetEquipmentTypeID(id)
	return epdc
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epdc *EquipmentPositionDefinitionCreate) SetNillableEquipmentTypeID(id *int) *EquipmentPositionDefinitionCreate {
	if id != nil {
		epdc = epdc.SetEquipmentTypeID(*id)
	}
	return epdc
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (epdc *EquipmentPositionDefinitionCreate) SetEquipmentType(e *EquipmentType) *EquipmentPositionDefinitionCreate {
	return epdc.SetEquipmentTypeID(e.ID)
}

// Save creates the EquipmentPositionDefinition in the database.
func (epdc *EquipmentPositionDefinitionCreate) Save(ctx context.Context) (*EquipmentPositionDefinition, error) {
	if _, ok := epdc.mutation.CreateTime(); !ok {
		v := equipmentpositiondefinition.DefaultCreateTime()
		epdc.mutation.SetCreateTime(v)
	}
	if _, ok := epdc.mutation.UpdateTime(); !ok {
		v := equipmentpositiondefinition.DefaultUpdateTime()
		epdc.mutation.SetUpdateTime(v)
	}
	if _, ok := epdc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	var (
		err  error
		node *EquipmentPositionDefinition
	)
	if len(epdc.hooks) == 0 {
		node, err = epdc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPositionDefinitionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			epdc.mutation = mutation
			node, err = epdc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(epdc.hooks) - 1; i >= 0; i-- {
			mut = epdc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, epdc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (epdc *EquipmentPositionDefinitionCreate) SaveX(ctx context.Context) *EquipmentPositionDefinition {
	v, err := epdc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epdc *EquipmentPositionDefinitionCreate) sqlSave(ctx context.Context) (*EquipmentPositionDefinition, error) {
	var (
		epd   = &EquipmentPositionDefinition{config: epdc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: equipmentpositiondefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentpositiondefinition.FieldID,
			},
		}
	)
	if value, ok := epdc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentpositiondefinition.FieldCreateTime,
		})
		epd.CreateTime = value
	}
	if value, ok := epdc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentpositiondefinition.FieldUpdateTime,
		})
		epd.UpdateTime = value
	}
	if value, ok := epdc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentpositiondefinition.FieldName,
		})
		epd.Name = value
	}
	if value, ok := epdc.mutation.Index(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentpositiondefinition.FieldIndex,
		})
		epd.Index = value
	}
	if value, ok := epdc.mutation.VisibilityLabel(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentpositiondefinition.FieldVisibilityLabel,
		})
		epd.VisibilityLabel = value
	}
	if nodes := epdc.mutation.PositionsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := epdc.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, epdc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	epd.ID = int(id)
	return epd, nil
}
