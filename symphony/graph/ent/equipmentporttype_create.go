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
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// EquipmentPortTypeCreate is the builder for creating a EquipmentPortType entity.
type EquipmentPortTypeCreate struct {
	config
	create_time         *time.Time
	update_time         *time.Time
	name                *string
	property_types      map[string]struct{}
	link_property_types map[string]struct{}
	port_definitions    map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (eptc *EquipmentPortTypeCreate) SetCreateTime(t time.Time) *EquipmentPortTypeCreate {
	eptc.create_time = &t
	return eptc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (eptc *EquipmentPortTypeCreate) SetNillableCreateTime(t *time.Time) *EquipmentPortTypeCreate {
	if t != nil {
		eptc.SetCreateTime(*t)
	}
	return eptc
}

// SetUpdateTime sets the update_time field.
func (eptc *EquipmentPortTypeCreate) SetUpdateTime(t time.Time) *EquipmentPortTypeCreate {
	eptc.update_time = &t
	return eptc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (eptc *EquipmentPortTypeCreate) SetNillableUpdateTime(t *time.Time) *EquipmentPortTypeCreate {
	if t != nil {
		eptc.SetUpdateTime(*t)
	}
	return eptc
}

// SetName sets the name field.
func (eptc *EquipmentPortTypeCreate) SetName(s string) *EquipmentPortTypeCreate {
	eptc.name = &s
	return eptc
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (eptc *EquipmentPortTypeCreate) AddPropertyTypeIDs(ids ...string) *EquipmentPortTypeCreate {
	if eptc.property_types == nil {
		eptc.property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptc.property_types[ids[i]] = struct{}{}
	}
	return eptc
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (eptc *EquipmentPortTypeCreate) AddPropertyTypes(p ...*PropertyType) *EquipmentPortTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptc.AddPropertyTypeIDs(ids...)
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (eptc *EquipmentPortTypeCreate) AddLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeCreate {
	if eptc.link_property_types == nil {
		eptc.link_property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptc.link_property_types[ids[i]] = struct{}{}
	}
	return eptc
}

// AddLinkPropertyTypes adds the link_property_types edges to PropertyType.
func (eptc *EquipmentPortTypeCreate) AddLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptc.AddLinkPropertyTypeIDs(ids...)
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (eptc *EquipmentPortTypeCreate) AddPortDefinitionIDs(ids ...string) *EquipmentPortTypeCreate {
	if eptc.port_definitions == nil {
		eptc.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		eptc.port_definitions[ids[i]] = struct{}{}
	}
	return eptc
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (eptc *EquipmentPortTypeCreate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptc.AddPortDefinitionIDs(ids...)
}

// Save creates the EquipmentPortType in the database.
func (eptc *EquipmentPortTypeCreate) Save(ctx context.Context) (*EquipmentPortType, error) {
	if eptc.create_time == nil {
		v := equipmentporttype.DefaultCreateTime()
		eptc.create_time = &v
	}
	if eptc.update_time == nil {
		v := equipmentporttype.DefaultUpdateTime()
		eptc.update_time = &v
	}
	if eptc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	return eptc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (eptc *EquipmentPortTypeCreate) SaveX(ctx context.Context) *EquipmentPortType {
	v, err := eptc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (eptc *EquipmentPortTypeCreate) sqlSave(ctx context.Context) (*EquipmentPortType, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(eptc.driver.Dialect())
		ept     = &EquipmentPortType{config: eptc.config}
	)
	tx, err := eptc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(equipmentporttype.Table).Default()
	if value := eptc.create_time; value != nil {
		insert.Set(equipmentporttype.FieldCreateTime, *value)
		ept.CreateTime = *value
	}
	if value := eptc.update_time; value != nil {
		insert.Set(equipmentporttype.FieldUpdateTime, *value)
		ept.UpdateTime = *value
	}
	if value := eptc.name; value != nil {
		insert.Set(equipmentporttype.FieldName, *value)
		ept.Name = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(equipmentporttype.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	ept.ID = strconv.FormatInt(id, 10)
	if len(eptc.property_types) > 0 {
		p := sql.P()
		for eid := range eptc.property_types {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(propertytype.FieldID, eid)
		}
		query, args := builder.Update(equipmentporttype.PropertyTypesTable).
			Set(equipmentporttype.PropertyTypesColumn, id).
			Where(sql.And(p, sql.IsNull(equipmentporttype.PropertyTypesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(eptc.property_types) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"EquipmentPortType\"", keys(eptc.property_types))})
		}
	}
	if len(eptc.link_property_types) > 0 {
		p := sql.P()
		for eid := range eptc.link_property_types {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(propertytype.FieldID, eid)
		}
		query, args := builder.Update(equipmentporttype.LinkPropertyTypesTable).
			Set(equipmentporttype.LinkPropertyTypesColumn, id).
			Where(sql.And(p, sql.IsNull(equipmentporttype.LinkPropertyTypesColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(eptc.link_property_types) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"link_property_types\" %v already connected to a different \"EquipmentPortType\"", keys(eptc.link_property_types))})
		}
	}
	if len(eptc.port_definitions) > 0 {
		p := sql.P()
		for eid := range eptc.port_definitions {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(equipmentportdefinition.FieldID, eid)
		}
		query, args := builder.Update(equipmentporttype.PortDefinitionsTable).
			Set(equipmentporttype.PortDefinitionsColumn, id).
			Where(sql.And(p, sql.IsNull(equipmentporttype.PortDefinitionsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(eptc.port_definitions) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"port_definitions\" %v already connected to a different \"EquipmentPortType\"", keys(eptc.port_definitions))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ept, nil
}
