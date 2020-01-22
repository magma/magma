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
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
)

// EquipmentPositionDefinitionCreate is the builder for creating a EquipmentPositionDefinition entity.
type EquipmentPositionDefinitionCreate struct {
	config
	create_time      *time.Time
	update_time      *time.Time
	name             *string
	index            *int
	visibility_label *string
	positions        map[string]struct{}
	equipment_type   map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (epdc *EquipmentPositionDefinitionCreate) SetCreateTime(t time.Time) *EquipmentPositionDefinitionCreate {
	epdc.create_time = &t
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
	epdc.update_time = &t
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
	epdc.name = &s
	return epdc
}

// SetIndex sets the index field.
func (epdc *EquipmentPositionDefinitionCreate) SetIndex(i int) *EquipmentPositionDefinitionCreate {
	epdc.index = &i
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
	epdc.visibility_label = &s
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
func (epdc *EquipmentPositionDefinitionCreate) AddPositionIDs(ids ...string) *EquipmentPositionDefinitionCreate {
	if epdc.positions == nil {
		epdc.positions = make(map[string]struct{})
	}
	for i := range ids {
		epdc.positions[ids[i]] = struct{}{}
	}
	return epdc
}

// AddPositions adds the positions edges to EquipmentPosition.
func (epdc *EquipmentPositionDefinitionCreate) AddPositions(e ...*EquipmentPosition) *EquipmentPositionDefinitionCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdc.AddPositionIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epdc *EquipmentPositionDefinitionCreate) SetEquipmentTypeID(id string) *EquipmentPositionDefinitionCreate {
	if epdc.equipment_type == nil {
		epdc.equipment_type = make(map[string]struct{})
	}
	epdc.equipment_type[id] = struct{}{}
	return epdc
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epdc *EquipmentPositionDefinitionCreate) SetNillableEquipmentTypeID(id *string) *EquipmentPositionDefinitionCreate {
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
	if epdc.create_time == nil {
		v := equipmentpositiondefinition.DefaultCreateTime()
		epdc.create_time = &v
	}
	if epdc.update_time == nil {
		v := equipmentpositiondefinition.DefaultUpdateTime()
		epdc.update_time = &v
	}
	if epdc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if len(epdc.equipment_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_type\"")
	}
	return epdc.sqlSave(ctx)
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
				Type:   field.TypeString,
				Column: equipmentpositiondefinition.FieldID,
			},
		}
	)
	if value := epdc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentpositiondefinition.FieldCreateTime,
		})
		epd.CreateTime = *value
	}
	if value := epdc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentpositiondefinition.FieldUpdateTime,
		})
		epd.UpdateTime = *value
	}
	if value := epdc.name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentpositiondefinition.FieldName,
		})
		epd.Name = *value
	}
	if value := epdc.index; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: equipmentpositiondefinition.FieldIndex,
		})
		epd.Index = *value
	}
	if value := epdc.visibility_label; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentpositiondefinition.FieldVisibilityLabel,
		})
		epd.VisibilityLabel = *value
	}
	if nodes := epdc.positions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentpositiondefinition.PositionsTable,
			Columns: []string{equipmentpositiondefinition.PositionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentposition.FieldID,
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
	if nodes := epdc.equipment_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentpositiondefinition.EquipmentTypeTable,
			Columns: []string{equipmentpositiondefinition.EquipmentTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmenttype.FieldID,
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
	if err := sqlgraph.CreateNode(ctx, epdc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	epd.ID = strconv.FormatInt(id, 10)
	return epd, nil
}
