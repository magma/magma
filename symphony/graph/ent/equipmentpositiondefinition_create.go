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
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
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
		res     sql.Result
		builder = sql.Dialect(epdc.driver.Dialect())
		epd     = &EquipmentPositionDefinition{config: epdc.config}
	)
	tx, err := epdc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(equipmentpositiondefinition.Table).Default()
	if value := epdc.create_time; value != nil {
		insert.Set(equipmentpositiondefinition.FieldCreateTime, *value)
		epd.CreateTime = *value
	}
	if value := epdc.update_time; value != nil {
		insert.Set(equipmentpositiondefinition.FieldUpdateTime, *value)
		epd.UpdateTime = *value
	}
	if value := epdc.name; value != nil {
		insert.Set(equipmentpositiondefinition.FieldName, *value)
		epd.Name = *value
	}
	if value := epdc.index; value != nil {
		insert.Set(equipmentpositiondefinition.FieldIndex, *value)
		epd.Index = *value
	}
	if value := epdc.visibility_label; value != nil {
		insert.Set(equipmentpositiondefinition.FieldVisibilityLabel, *value)
		epd.VisibilityLabel = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(equipmentpositiondefinition.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	epd.ID = strconv.FormatInt(id, 10)
	if len(epdc.positions) > 0 {
		p := sql.P()
		for eid := range epdc.positions {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(equipmentposition.FieldID, eid)
		}
		query, args := builder.Update(equipmentpositiondefinition.PositionsTable).
			Set(equipmentpositiondefinition.PositionsColumn, id).
			Where(sql.And(p, sql.IsNull(equipmentpositiondefinition.PositionsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(epdc.positions) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"positions\" %v already connected to a different \"EquipmentPositionDefinition\"", keys(epdc.positions))})
		}
	}
	if len(epdc.equipment_type) > 0 {
		for eid := range epdc.equipment_type {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(equipmentpositiondefinition.EquipmentTypeTable).
				Set(equipmentpositiondefinition.EquipmentTypeColumn, eid).
				Where(sql.EQ(equipmentpositiondefinition.FieldID, id)).
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
