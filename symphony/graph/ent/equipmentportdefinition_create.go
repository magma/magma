// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
)

// EquipmentPortDefinitionCreate is the builder for creating a EquipmentPortDefinition entity.
type EquipmentPortDefinitionCreate struct {
	config
	create_time         *time.Time
	update_time         *time.Time
	name                *string
	_type               *string
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

// SetType sets the type field.
func (epdc *EquipmentPortDefinitionCreate) SetType(s string) *EquipmentPortDefinitionCreate {
	epdc._type = &s
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
	if epdc._type == nil {
		return nil, errors.New("ent: missing required field \"type\"")
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
		res     sql.Result
		builder = sql.Dialect(epdc.driver.Dialect())
		epd     = &EquipmentPortDefinition{config: epdc.config}
	)
	tx, err := epdc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(equipmentportdefinition.Table).Default()
	if value := epdc.create_time; value != nil {
		insert.Set(equipmentportdefinition.FieldCreateTime, *value)
		epd.CreateTime = *value
	}
	if value := epdc.update_time; value != nil {
		insert.Set(equipmentportdefinition.FieldUpdateTime, *value)
		epd.UpdateTime = *value
	}
	if value := epdc.name; value != nil {
		insert.Set(equipmentportdefinition.FieldName, *value)
		epd.Name = *value
	}
	if value := epdc._type; value != nil {
		insert.Set(equipmentportdefinition.FieldType, *value)
		epd.Type = *value
	}
	if value := epdc.index; value != nil {
		insert.Set(equipmentportdefinition.FieldIndex, *value)
		epd.Index = *value
	}
	if value := epdc.bandwidth; value != nil {
		insert.Set(equipmentportdefinition.FieldBandwidth, *value)
		epd.Bandwidth = *value
	}
	if value := epdc.visibility_label; value != nil {
		insert.Set(equipmentportdefinition.FieldVisibilityLabel, *value)
		epd.VisibilityLabel = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(equipmentportdefinition.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	epd.ID = strconv.FormatInt(id, 10)
	if len(epdc.equipment_port_type) > 0 {
		for eid := range epdc.equipment_port_type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipmentportdefinition.EquipmentPortTypeTable).
				Set(equipmentportdefinition.EquipmentPortTypeColumn, eid).
				Where(sql.EQ(equipmentportdefinition.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(epdc.ports) > 0 {
		p := sql.P()
		for eid := range epdc.ports {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(equipmentport.FieldID, eid)
		}
		query, args := builder.Update(equipmentportdefinition.PortsTable).
			Set(equipmentportdefinition.PortsColumn, id).
			Where(sql.And(p, sql.IsNull(equipmentportdefinition.PortsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(epdc.ports) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"EquipmentPortDefinition\"", keys(epdc.ports))})
		}
	}
	if len(epdc.equipment_type) > 0 {
		for eid := range epdc.equipment_type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipmentportdefinition.EquipmentTypeTable).
				Set(equipmentportdefinition.EquipmentTypeColumn, eid).
				Where(sql.EQ(equipmentportdefinition.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return epd, nil
}
