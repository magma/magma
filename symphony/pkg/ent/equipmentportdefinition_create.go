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
	"github.com/facebookincubator/symphony/pkg/ent/equipmentport"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"
)

// EquipmentPortDefinitionCreate is the builder for creating a EquipmentPortDefinition entity.
type EquipmentPortDefinitionCreate struct {
	config
	mutation *EquipmentPortDefinitionMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (epdc *EquipmentPortDefinitionCreate) SetCreateTime(t time.Time) *EquipmentPortDefinitionCreate {
	epdc.mutation.SetCreateTime(t)
	return epdc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableCreateTime(t *time.Time) *EquipmentPortDefinitionCreate {
	if t != nil {
		epdc.SetCreateTime(*t)
	}
	return epdc
}

// SetUpdateTime sets the update_time field.
func (epdc *EquipmentPortDefinitionCreate) SetUpdateTime(t time.Time) *EquipmentPortDefinitionCreate {
	epdc.mutation.SetUpdateTime(t)
	return epdc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableUpdateTime(t *time.Time) *EquipmentPortDefinitionCreate {
	if t != nil {
		epdc.SetUpdateTime(*t)
	}
	return epdc
}

// SetName sets the name field.
func (epdc *EquipmentPortDefinitionCreate) SetName(s string) *EquipmentPortDefinitionCreate {
	epdc.mutation.SetName(s)
	return epdc
}

// SetIndex sets the index field.
func (epdc *EquipmentPortDefinitionCreate) SetIndex(i int) *EquipmentPortDefinitionCreate {
	epdc.mutation.SetIndex(i)
	return epdc
}

// SetNillableIndex sets the index field if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableIndex(i *int) *EquipmentPortDefinitionCreate {
	if i != nil {
		epdc.SetIndex(*i)
	}
	return epdc
}

// SetBandwidth sets the bandwidth field.
func (epdc *EquipmentPortDefinitionCreate) SetBandwidth(s string) *EquipmentPortDefinitionCreate {
	epdc.mutation.SetBandwidth(s)
	return epdc
}

// SetNillableBandwidth sets the bandwidth field if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableBandwidth(s *string) *EquipmentPortDefinitionCreate {
	if s != nil {
		epdc.SetBandwidth(*s)
	}
	return epdc
}

// SetVisibilityLabel sets the visibility_label field.
func (epdc *EquipmentPortDefinitionCreate) SetVisibilityLabel(s string) *EquipmentPortDefinitionCreate {
	epdc.mutation.SetVisibilityLabel(s)
	return epdc
}

// SetNillableVisibilityLabel sets the visibility_label field if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableVisibilityLabel(s *string) *EquipmentPortDefinitionCreate {
	if s != nil {
		epdc.SetVisibilityLabel(*s)
	}
	return epdc
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (epdc *EquipmentPortDefinitionCreate) SetEquipmentPortTypeID(id int) *EquipmentPortDefinitionCreate {
	epdc.mutation.SetEquipmentPortTypeID(id)
	return epdc
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableEquipmentPortTypeID(id *int) *EquipmentPortDefinitionCreate {
	if id != nil {
		epdc = epdc.SetEquipmentPortTypeID(*id)
	}
	return epdc
}

// SetEquipmentPortType sets the equipment_port_type edge to EquipmentPortType.
func (epdc *EquipmentPortDefinitionCreate) SetEquipmentPortType(e *EquipmentPortType) *EquipmentPortDefinitionCreate {
	return epdc.SetEquipmentPortTypeID(e.ID)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (epdc *EquipmentPortDefinitionCreate) AddPortIDs(ids ...int) *EquipmentPortDefinitionCreate {
	epdc.mutation.AddPortIDs(ids...)
	return epdc
}

// AddPorts adds the ports edges to EquipmentPort.
func (epdc *EquipmentPortDefinitionCreate) AddPorts(e ...*EquipmentPort) *EquipmentPortDefinitionCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdc.AddPortIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epdc *EquipmentPortDefinitionCreate) SetEquipmentTypeID(id int) *EquipmentPortDefinitionCreate {
	epdc.mutation.SetEquipmentTypeID(id)
	return epdc
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableEquipmentTypeID(id *int) *EquipmentPortDefinitionCreate {
	if id != nil {
		epdc = epdc.SetEquipmentTypeID(*id)
	}
	return epdc
}

// SetEquipmentType sets the equipment_type edge to EquipmentType.
func (epdc *EquipmentPortDefinitionCreate) SetEquipmentType(e *EquipmentType) *EquipmentPortDefinitionCreate {
	return epdc.SetEquipmentTypeID(e.ID)
}

// Save creates the EquipmentPortDefinition in the database.
func (epdc *EquipmentPortDefinitionCreate) Save(ctx context.Context) (*EquipmentPortDefinition, error) {
	if _, ok := epdc.mutation.CreateTime(); !ok {
		v := equipmentportdefinition.DefaultCreateTime()
		epdc.mutation.SetCreateTime(v)
	}
	if _, ok := epdc.mutation.UpdateTime(); !ok {
		v := equipmentportdefinition.DefaultUpdateTime()
		epdc.mutation.SetUpdateTime(v)
	}
	if _, ok := epdc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	var (
		err  error
		node *EquipmentPortDefinition
	)
	if len(epdc.hooks) == 0 {
		node, err = epdc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EquipmentPortDefinitionMutation)
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
func (epdc *EquipmentPortDefinitionCreate) SaveX(ctx context.Context) *EquipmentPortDefinition {
	v, err := epdc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (epdc *EquipmentPortDefinitionCreate) sqlSave(ctx context.Context) (*EquipmentPortDefinition, error) {
	var (
		epd   = &EquipmentPortDefinition{config: epdc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: equipmentportdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: equipmentportdefinition.FieldID,
			},
		}
	)
	if value, ok := epdc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentportdefinition.FieldCreateTime,
		})
		epd.CreateTime = value
	}
	if value, ok := epdc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: equipmentportdefinition.FieldUpdateTime,
		})
		epd.UpdateTime = value
	}
	if value, ok := epdc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldName,
		})
		epd.Name = value
	}
	if value, ok := epdc.mutation.Index(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: equipmentportdefinition.FieldIndex,
		})
		epd.Index = value
	}
	if value, ok := epdc.mutation.Bandwidth(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldBandwidth,
		})
		epd.Bandwidth = value
	}
	if value, ok := epdc.mutation.VisibilityLabel(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: equipmentportdefinition.FieldVisibilityLabel,
		})
		epd.VisibilityLabel = value
	}
	if nodes := epdc.mutation.EquipmentPortTypeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := epdc.mutation.PortsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := epdc.mutation.EquipmentTypeIDs(); len(nodes) > 0 {
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
