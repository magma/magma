// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// EquipmentPortTypeUpdate is the builder for updating EquipmentPortType entities.
type EquipmentPortTypeUpdate struct {
	config

	update_time              *time.Time
	name                     *string
	property_types           map[string]struct{}
	link_property_types      map[string]struct{}
	port_definitions         map[string]struct{}
	removedPropertyTypes     map[string]struct{}
	removedLinkPropertyTypes map[string]struct{}
	removedPortDefinitions   map[string]struct{}
	predicates               []predicate.EquipmentPortType
}

// Where adds a new predicate for the builder.
func (eptu *EquipmentPortTypeUpdate) Where(ps ...predicate.EquipmentPortType) *EquipmentPortTypeUpdate {
	eptu.predicates = append(eptu.predicates, ps...)
	return eptu
}

// SetName sets the name field.
func (eptu *EquipmentPortTypeUpdate) SetName(s string) *EquipmentPortTypeUpdate {
	eptu.name = &s
	return eptu
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) AddPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.property_types == nil {
		eptu.property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptu.property_types[ids[i]] = struct{}{}
	}
	return eptu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) AddPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.AddPropertyTypeIDs(ids...)
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) AddLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.link_property_types == nil {
		eptu.link_property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptu.link_property_types[ids[i]] = struct{}{}
	}
	return eptu
}

// AddLinkPropertyTypes adds the link_property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) AddLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.AddLinkPropertyTypeIDs(ids...)
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (eptu *EquipmentPortTypeUpdate) AddPortDefinitionIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.port_definitions == nil {
		eptu.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		eptu.port_definitions[ids[i]] = struct{}{}
	}
	return eptu
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (eptu *EquipmentPortTypeUpdate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptu.AddPortDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) RemovePropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.removedPropertyTypes == nil {
		eptu.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		eptu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return eptu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.RemovePropertyTypeIDs(ids...)
}

// RemoveLinkPropertyTypeIDs removes the link_property_types edge to PropertyType by ids.
func (eptu *EquipmentPortTypeUpdate) RemoveLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.removedLinkPropertyTypes == nil {
		eptu.removedLinkPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		eptu.removedLinkPropertyTypes[ids[i]] = struct{}{}
	}
	return eptu
}

// RemoveLinkPropertyTypes removes link_property_types edges to PropertyType.
func (eptu *EquipmentPortTypeUpdate) RemoveLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptu.RemoveLinkPropertyTypeIDs(ids...)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (eptu *EquipmentPortTypeUpdate) RemovePortDefinitionIDs(ids ...string) *EquipmentPortTypeUpdate {
	if eptu.removedPortDefinitions == nil {
		eptu.removedPortDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		eptu.removedPortDefinitions[ids[i]] = struct{}{}
	}
	return eptu
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (eptu *EquipmentPortTypeUpdate) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptu.RemovePortDefinitionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (eptu *EquipmentPortTypeUpdate) Save(ctx context.Context) (int, error) {
	if eptu.update_time == nil {
		v := equipmentporttype.UpdateDefaultUpdateTime()
		eptu.update_time = &v
	}
	return eptu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (eptu *EquipmentPortTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := eptu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eptu *EquipmentPortTypeUpdate) Exec(ctx context.Context) error {
	_, err := eptu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eptu *EquipmentPortTypeUpdate) ExecX(ctx context.Context) {
	if err := eptu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eptu *EquipmentPortTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(eptu.driver.Dialect())
		selector = builder.Select(equipmentporttype.FieldID).From(builder.Table(equipmentporttype.Table))
	)
	for _, p := range eptu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = eptu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := eptu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmentporttype.Table)
	)
	updater = updater.Where(sql.InInts(equipmentporttype.FieldID, ids...))
	if value := eptu.update_time; value != nil {
		updater.Set(equipmentporttype.FieldUpdateTime, *value)
	}
	if value := eptu.name; value != nil {
		updater.Set(equipmentporttype.FieldName, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eptu.removedPropertyTypes) > 0 {
		eids := make([]int, len(eptu.removedPropertyTypes))
		for eid := range eptu.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentporttype.PropertyTypesTable).
			SetNull(equipmentporttype.PropertyTypesColumn).
			Where(sql.InInts(equipmentporttype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eptu.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eptu.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(equipmentporttype.PropertyTypesTable).
				Set(equipmentporttype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentporttype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(eptu.property_types) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"EquipmentPortType\"", keys(eptu.property_types))})
			}
		}
	}
	if len(eptu.removedLinkPropertyTypes) > 0 {
		eids := make([]int, len(eptu.removedLinkPropertyTypes))
		for eid := range eptu.removedLinkPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentporttype.LinkPropertyTypesTable).
			SetNull(equipmentporttype.LinkPropertyTypesColumn).
			Where(sql.InInts(equipmentporttype.LinkPropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eptu.link_property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eptu.link_property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(equipmentporttype.LinkPropertyTypesTable).
				Set(equipmentporttype.LinkPropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentporttype.LinkPropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(eptu.link_property_types) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"link_property_types\" %v already connected to a different \"EquipmentPortType\"", keys(eptu.link_property_types))})
			}
		}
	}
	if len(eptu.removedPortDefinitions) > 0 {
		eids := make([]int, len(eptu.removedPortDefinitions))
		for eid := range eptu.removedPortDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentporttype.PortDefinitionsTable).
			SetNull(equipmentporttype.PortDefinitionsColumn).
			Where(sql.InInts(equipmentporttype.PortDefinitionsColumn, ids...)).
			Where(sql.InInts(equipmentportdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(eptu.port_definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eptu.port_definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentportdefinition.FieldID, eid)
			}
			query, args := builder.Update(equipmentporttype.PortDefinitionsTable).
				Set(equipmentporttype.PortDefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmentporttype.PortDefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(eptu.port_definitions) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"port_definitions\" %v already connected to a different \"EquipmentPortType\"", keys(eptu.port_definitions))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// EquipmentPortTypeUpdateOne is the builder for updating a single EquipmentPortType entity.
type EquipmentPortTypeUpdateOne struct {
	config
	id string

	update_time              *time.Time
	name                     *string
	property_types           map[string]struct{}
	link_property_types      map[string]struct{}
	port_definitions         map[string]struct{}
	removedPropertyTypes     map[string]struct{}
	removedLinkPropertyTypes map[string]struct{}
	removedPortDefinitions   map[string]struct{}
}

// SetName sets the name field.
func (eptuo *EquipmentPortTypeUpdateOne) SetName(s string) *EquipmentPortTypeUpdateOne {
	eptuo.name = &s
	return eptuo
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.property_types == nil {
		eptuo.property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.property_types[ids[i]] = struct{}{}
	}
	return eptuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.AddPropertyTypeIDs(ids...)
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.link_property_types == nil {
		eptuo.link_property_types = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.link_property_types[ids[i]] = struct{}{}
	}
	return eptuo
}

// AddLinkPropertyTypes adds the link_property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) AddLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.AddLinkPropertyTypeIDs(ids...)
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (eptuo *EquipmentPortTypeUpdateOne) AddPortDefinitionIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.port_definitions == nil {
		eptuo.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.port_definitions[ids[i]] = struct{}{}
	}
	return eptuo
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (eptuo *EquipmentPortTypeUpdateOne) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptuo.AddPortDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.removedPropertyTypes == nil {
		eptuo.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return eptuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.RemovePropertyTypeIDs(ids...)
}

// RemoveLinkPropertyTypeIDs removes the link_property_types edge to PropertyType by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemoveLinkPropertyTypeIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.removedLinkPropertyTypes == nil {
		eptuo.removedLinkPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.removedLinkPropertyTypes[ids[i]] = struct{}{}
	}
	return eptuo
}

// RemoveLinkPropertyTypes removes link_property_types edges to PropertyType.
func (eptuo *EquipmentPortTypeUpdateOne) RemoveLinkPropertyTypes(p ...*PropertyType) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return eptuo.RemoveLinkPropertyTypeIDs(ids...)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePortDefinitionIDs(ids ...string) *EquipmentPortTypeUpdateOne {
	if eptuo.removedPortDefinitions == nil {
		eptuo.removedPortDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		eptuo.removedPortDefinitions[ids[i]] = struct{}{}
	}
	return eptuo
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (eptuo *EquipmentPortTypeUpdateOne) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentPortTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return eptuo.RemovePortDefinitionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (eptuo *EquipmentPortTypeUpdateOne) Save(ctx context.Context) (*EquipmentPortType, error) {
	if eptuo.update_time == nil {
		v := equipmentporttype.UpdateDefaultUpdateTime()
		eptuo.update_time = &v
	}
	return eptuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (eptuo *EquipmentPortTypeUpdateOne) SaveX(ctx context.Context) *EquipmentPortType {
	ept, err := eptuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ept
}

// Exec executes the query on the entity.
func (eptuo *EquipmentPortTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := eptuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eptuo *EquipmentPortTypeUpdateOne) ExecX(ctx context.Context) {
	if err := eptuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eptuo *EquipmentPortTypeUpdateOne) sqlSave(ctx context.Context) (ept *EquipmentPortType, err error) {
	var (
		builder  = sql.Dialect(eptuo.driver.Dialect())
		selector = builder.Select(equipmentporttype.Columns...).From(builder.Table(equipmentporttype.Table))
	)
	equipmentporttype.ID(eptuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = eptuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		ept = &EquipmentPortType{config: eptuo.config}
		if err := ept.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into EquipmentPortType: %v", err)
		}
		id = ept.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("EquipmentPortType with id: %v", eptuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one EquipmentPortType with the same id: %v", eptuo.id)
	}

	tx, err := eptuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmentporttype.Table)
	)
	updater = updater.Where(sql.InInts(equipmentporttype.FieldID, ids...))
	if value := eptuo.update_time; value != nil {
		updater.Set(equipmentporttype.FieldUpdateTime, *value)
		ept.UpdateTime = *value
	}
	if value := eptuo.name; value != nil {
		updater.Set(equipmentporttype.FieldName, *value)
		ept.Name = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(eptuo.removedPropertyTypes) > 0 {
		eids := make([]int, len(eptuo.removedPropertyTypes))
		for eid := range eptuo.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentporttype.PropertyTypesTable).
			SetNull(equipmentporttype.PropertyTypesColumn).
			Where(sql.InInts(equipmentporttype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(eptuo.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eptuo.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(eptuo.property_types) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"EquipmentPortType\"", keys(eptuo.property_types))})
			}
		}
	}
	if len(eptuo.removedLinkPropertyTypes) > 0 {
		eids := make([]int, len(eptuo.removedLinkPropertyTypes))
		for eid := range eptuo.removedLinkPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentporttype.LinkPropertyTypesTable).
			SetNull(equipmentporttype.LinkPropertyTypesColumn).
			Where(sql.InInts(equipmentporttype.LinkPropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(eptuo.link_property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eptuo.link_property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(eptuo.link_property_types) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"link_property_types\" %v already connected to a different \"EquipmentPortType\"", keys(eptuo.link_property_types))})
			}
		}
	}
	if len(eptuo.removedPortDefinitions) > 0 {
		eids := make([]int, len(eptuo.removedPortDefinitions))
		for eid := range eptuo.removedPortDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmentporttype.PortDefinitionsTable).
			SetNull(equipmentporttype.PortDefinitionsColumn).
			Where(sql.InInts(equipmentporttype.PortDefinitionsColumn, ids...)).
			Where(sql.InInts(equipmentportdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(eptuo.port_definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range eptuo.port_definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(eptuo.port_definitions) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"port_definitions\" %v already connected to a different \"EquipmentPortType\"", keys(eptuo.port_definitions))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ept, nil
}
