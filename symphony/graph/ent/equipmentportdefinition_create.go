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
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
)

// EquipmentPortDefinitionCreate is the builder for creating a EquipmentPortDefinition entity.
type EquipmentPortDefinitionCreate struct {
	config
	create_time         *time.Time
	update_time         *time.Time
	name                *string
	index               *int
	bandwidth           *string
	visibility_label    *string
	equipment_port_type map[string]struct{}
	ports               map[string]struct{}
	equipment_type      map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (epdc *EquipmentPortDefinitionCreate) SetCreateTime(t time.Time) *EquipmentPortDefinitionCreate {
	epdc.create_time = &t
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
	epdc.update_time = &t
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
	epdc.name = &s
	return epdc
}

// SetIndex sets the index field.
func (epdc *EquipmentPortDefinitionCreate) SetIndex(i int) *EquipmentPortDefinitionCreate {
	epdc.index = &i
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
	epdc.bandwidth = &s
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
	epdc.visibility_label = &s
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
func (epdc *EquipmentPortDefinitionCreate) SetEquipmentPortTypeID(id string) *EquipmentPortDefinitionCreate {
	if epdc.equipment_port_type == nil {
		epdc.equipment_port_type = make(map[string]struct{})
	}
	epdc.equipment_port_type[id] = struct{}{}
	return epdc
}

// SetNillableEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableEquipmentPortTypeID(id *string) *EquipmentPortDefinitionCreate {
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
func (epdc *EquipmentPortDefinitionCreate) AddPortIDs(ids ...string) *EquipmentPortDefinitionCreate {
	if epdc.ports == nil {
		epdc.ports = make(map[string]struct{})
	}
	for i := range ids {
		epdc.ports[ids[i]] = struct{}{}
	}
	return epdc
}

// AddPorts adds the ports edges to EquipmentPort.
func (epdc *EquipmentPortDefinitionCreate) AddPorts(e ...*EquipmentPort) *EquipmentPortDefinitionCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return epdc.AddPortIDs(ids...)
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (epdc *EquipmentPortDefinitionCreate) SetEquipmentTypeID(id string) *EquipmentPortDefinitionCreate {
	if epdc.equipment_type == nil {
		epdc.equipment_type = make(map[string]struct{})
	}
	epdc.equipment_type[id] = struct{}{}
	return epdc
}

// SetNillableEquipmentTypeID sets the equipment_type edge to EquipmentType by id if the given value is not nil.
func (epdc *EquipmentPortDefinitionCreate) SetNillableEquipmentTypeID(id *string) *EquipmentPortDefinitionCreate {
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
	if epdc.create_time == nil {
		v := equipmentportdefinition.DefaultCreateTime()
		epdc.create_time = &v
	}
	if epdc.update_time == nil {
		v := equipmentportdefinition.DefaultUpdateTime()
		epdc.update_time = &v
	}
	if epdc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if len(epdc.equipment_port_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_port_type\"")
	}
	if len(epdc.equipment_type) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"equipment_type\"")
	}
	return epdc.sqlSave(ctx)
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
		epd  = &EquipmentPortDefinition{config: epdc.config}
		spec = &sqlgraph.CreateSpec{
			Table: equipmentportdefinition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: equipmentportdefinition.FieldID,
			},
		}
	)
	if value := epdc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentportdefinition.FieldCreateTime,
		})
		epd.CreateTime = *value
	}
	if value := epdc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: equipmentportdefinition.FieldUpdateTime,
		})
		epd.UpdateTime = *value
	}
	if value := epdc.name; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentportdefinition.FieldName,
		})
		epd.Name = *value
	}
	if value := epdc.index; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: equipmentportdefinition.FieldIndex,
		})
		epd.Index = *value
	}
	if value := epdc.bandwidth; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentportdefinition.FieldBandwidth,
		})
		epd.Bandwidth = *value
	}
	if value := epdc.visibility_label; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: equipmentportdefinition.FieldVisibilityLabel,
		})
		epd.VisibilityLabel = *value
	}
	if nodes := epdc.equipment_port_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   equipmentportdefinition.EquipmentPortTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentPortTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentporttype.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := epdc.ports; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   equipmentportdefinition.PortsTable,
			Columns: []string{equipmentportdefinition.PortsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: equipmentport.FieldID,
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := epdc.equipment_type; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   equipmentportdefinition.EquipmentTypeTable,
			Columns: []string{equipmentportdefinition.EquipmentTypeColumn},
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
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, epdc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	epd.ID = strconv.FormatInt(id, 10)
	return epd, nil
}
