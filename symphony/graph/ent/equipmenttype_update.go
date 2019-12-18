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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// EquipmentTypeUpdate is the builder for updating EquipmentType entities.
type EquipmentTypeUpdate struct {
	config

	update_time                *time.Time
	name                       *string
	port_definitions           map[string]struct{}
	position_definitions       map[string]struct{}
	property_types             map[string]struct{}
	equipment                  map[string]struct{}
	category                   map[string]struct{}
	removedPortDefinitions     map[string]struct{}
	removedPositionDefinitions map[string]struct{}
	removedPropertyTypes       map[string]struct{}
	removedEquipment           map[string]struct{}
	clearedCategory            bool
	predicates                 []predicate.EquipmentType
}

// Where adds a new predicate for the builder.
func (etu *EquipmentTypeUpdate) Where(ps ...predicate.EquipmentType) *EquipmentTypeUpdate {
	etu.predicates = append(etu.predicates, ps...)
	return etu
}

// SetName sets the name field.
func (etu *EquipmentTypeUpdate) SetName(s string) *EquipmentTypeUpdate {
	etu.name = &s
	return etu
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (etu *EquipmentTypeUpdate) AddPortDefinitionIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.port_definitions == nil {
		etu.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etu.port_definitions[ids[i]] = struct{}{}
	}
	return etu
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (etu *EquipmentTypeUpdate) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddPortDefinitionIDs(ids...)
}

// AddPositionDefinitionIDs adds the position_definitions edge to EquipmentPositionDefinition by ids.
func (etu *EquipmentTypeUpdate) AddPositionDefinitionIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.position_definitions == nil {
		etu.position_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etu.position_definitions[ids[i]] = struct{}{}
	}
	return etu
}

// AddPositionDefinitions adds the position_definitions edges to EquipmentPositionDefinition.
func (etu *EquipmentTypeUpdate) AddPositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddPositionDefinitionIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (etu *EquipmentTypeUpdate) AddPropertyTypeIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.property_types == nil {
		etu.property_types = make(map[string]struct{})
	}
	for i := range ids {
		etu.property_types[ids[i]] = struct{}{}
	}
	return etu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (etu *EquipmentTypeUpdate) AddPropertyTypes(p ...*PropertyType) *EquipmentTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etu.AddPropertyTypeIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (etu *EquipmentTypeUpdate) AddEquipmentIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.equipment == nil {
		etu.equipment = make(map[string]struct{})
	}
	for i := range ids {
		etu.equipment[ids[i]] = struct{}{}
	}
	return etu
}

// AddEquipment adds the equipment edges to Equipment.
func (etu *EquipmentTypeUpdate) AddEquipment(e ...*Equipment) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.AddEquipmentIDs(ids...)
}

// SetCategoryID sets the category edge to EquipmentCategory by id.
func (etu *EquipmentTypeUpdate) SetCategoryID(id string) *EquipmentTypeUpdate {
	if etu.category == nil {
		etu.category = make(map[string]struct{})
	}
	etu.category[id] = struct{}{}
	return etu
}

// SetNillableCategoryID sets the category edge to EquipmentCategory by id if the given value is not nil.
func (etu *EquipmentTypeUpdate) SetNillableCategoryID(id *string) *EquipmentTypeUpdate {
	if id != nil {
		etu = etu.SetCategoryID(*id)
	}
	return etu
}

// SetCategory sets the category edge to EquipmentCategory.
func (etu *EquipmentTypeUpdate) SetCategory(e *EquipmentCategory) *EquipmentTypeUpdate {
	return etu.SetCategoryID(e.ID)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (etu *EquipmentTypeUpdate) RemovePortDefinitionIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.removedPortDefinitions == nil {
		etu.removedPortDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		etu.removedPortDefinitions[ids[i]] = struct{}{}
	}
	return etu
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (etu *EquipmentTypeUpdate) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemovePortDefinitionIDs(ids...)
}

// RemovePositionDefinitionIDs removes the position_definitions edge to EquipmentPositionDefinition by ids.
func (etu *EquipmentTypeUpdate) RemovePositionDefinitionIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.removedPositionDefinitions == nil {
		etu.removedPositionDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		etu.removedPositionDefinitions[ids[i]] = struct{}{}
	}
	return etu
}

// RemovePositionDefinitions removes position_definitions edges to EquipmentPositionDefinition.
func (etu *EquipmentTypeUpdate) RemovePositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemovePositionDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (etu *EquipmentTypeUpdate) RemovePropertyTypeIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.removedPropertyTypes == nil {
		etu.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		etu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return etu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (etu *EquipmentTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *EquipmentTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etu.RemovePropertyTypeIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (etu *EquipmentTypeUpdate) RemoveEquipmentIDs(ids ...string) *EquipmentTypeUpdate {
	if etu.removedEquipment == nil {
		etu.removedEquipment = make(map[string]struct{})
	}
	for i := range ids {
		etu.removedEquipment[ids[i]] = struct{}{}
	}
	return etu
}

// RemoveEquipment removes equipment edges to Equipment.
func (etu *EquipmentTypeUpdate) RemoveEquipment(e ...*Equipment) *EquipmentTypeUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etu.RemoveEquipmentIDs(ids...)
}

// ClearCategory clears the category edge to EquipmentCategory.
func (etu *EquipmentTypeUpdate) ClearCategory() *EquipmentTypeUpdate {
	etu.clearedCategory = true
	return etu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (etu *EquipmentTypeUpdate) Save(ctx context.Context) (int, error) {
	if etu.update_time == nil {
		v := equipmenttype.UpdateDefaultUpdateTime()
		etu.update_time = &v
	}
	if len(etu.category) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return etu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (etu *EquipmentTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := etu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (etu *EquipmentTypeUpdate) Exec(ctx context.Context) error {
	_, err := etu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (etu *EquipmentTypeUpdate) ExecX(ctx context.Context) {
	if err := etu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (etu *EquipmentTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(etu.driver.Dialect())
		selector = builder.Select(equipmenttype.FieldID).From(builder.Table(equipmenttype.Table))
	)
	for _, p := range etu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = etu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := etu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmenttype.Table)
	)
	updater = updater.Where(sql.InInts(equipmenttype.FieldID, ids...))
	if value := etu.update_time; value != nil {
		updater.Set(equipmenttype.FieldUpdateTime, *value)
	}
	if value := etu.name; value != nil {
		updater.Set(equipmenttype.FieldName, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(etu.removedPortDefinitions) > 0 {
		eids := make([]int, len(etu.removedPortDefinitions))
		for eid := range etu.removedPortDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmenttype.PortDefinitionsTable).
			SetNull(equipmenttype.PortDefinitionsColumn).
			Where(sql.InInts(equipmenttype.PortDefinitionsColumn, ids...)).
			Where(sql.InInts(equipmentportdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(etu.port_definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range etu.port_definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentportdefinition.FieldID, eid)
			}
			query, args := builder.Update(equipmenttype.PortDefinitionsTable).
				Set(equipmenttype.PortDefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmenttype.PortDefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(etu.port_definitions) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"port_definitions\" %v already connected to a different \"EquipmentType\"", keys(etu.port_definitions))})
			}
		}
	}
	if len(etu.removedPositionDefinitions) > 0 {
		eids := make([]int, len(etu.removedPositionDefinitions))
		for eid := range etu.removedPositionDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmenttype.PositionDefinitionsTable).
			SetNull(equipmenttype.PositionDefinitionsColumn).
			Where(sql.InInts(equipmenttype.PositionDefinitionsColumn, ids...)).
			Where(sql.InInts(equipmentpositiondefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(etu.position_definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range etu.position_definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentpositiondefinition.FieldID, eid)
			}
			query, args := builder.Update(equipmenttype.PositionDefinitionsTable).
				Set(equipmenttype.PositionDefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmenttype.PositionDefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(etu.position_definitions) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"position_definitions\" %v already connected to a different \"EquipmentType\"", keys(etu.position_definitions))})
			}
		}
	}
	if len(etu.removedPropertyTypes) > 0 {
		eids := make([]int, len(etu.removedPropertyTypes))
		for eid := range etu.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmenttype.PropertyTypesTable).
			SetNull(equipmenttype.PropertyTypesColumn).
			Where(sql.InInts(equipmenttype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(etu.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range etu.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(equipmenttype.PropertyTypesTable).
				Set(equipmenttype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(equipmenttype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(etu.property_types) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"EquipmentType\"", keys(etu.property_types))})
			}
		}
	}
	if len(etu.removedEquipment) > 0 {
		eids := make([]int, len(etu.removedEquipment))
		for eid := range etu.removedEquipment {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmenttype.EquipmentTable).
			SetNull(equipmenttype.EquipmentColumn).
			Where(sql.InInts(equipmenttype.EquipmentColumn, ids...)).
			Where(sql.InInts(equipment.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(etu.equipment) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range etu.equipment {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipment.FieldID, eid)
			}
			query, args := builder.Update(equipmenttype.EquipmentTable).
				Set(equipmenttype.EquipmentColumn, id).
				Where(sql.And(p, sql.IsNull(equipmenttype.EquipmentColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(etu.equipment) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"equipment\" %v already connected to a different \"EquipmentType\"", keys(etu.equipment))})
			}
		}
	}
	if etu.clearedCategory {
		query, args := builder.Update(equipmenttype.CategoryTable).
			SetNull(equipmenttype.CategoryColumn).
			Where(sql.InInts(equipmentcategory.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(etu.category) > 0 {
		for eid := range etu.category {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmenttype.CategoryTable).
				Set(equipmenttype.CategoryColumn, eid).
				Where(sql.InInts(equipmenttype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// EquipmentTypeUpdateOne is the builder for updating a single EquipmentType entity.
type EquipmentTypeUpdateOne struct {
	config
	id string

	update_time                *time.Time
	name                       *string
	port_definitions           map[string]struct{}
	position_definitions       map[string]struct{}
	property_types             map[string]struct{}
	equipment                  map[string]struct{}
	category                   map[string]struct{}
	removedPortDefinitions     map[string]struct{}
	removedPositionDefinitions map[string]struct{}
	removedPropertyTypes       map[string]struct{}
	removedEquipment           map[string]struct{}
	clearedCategory            bool
}

// SetName sets the name field.
func (etuo *EquipmentTypeUpdateOne) SetName(s string) *EquipmentTypeUpdateOne {
	etuo.name = &s
	return etuo
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) AddPortDefinitionIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.port_definitions == nil {
		etuo.port_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etuo.port_definitions[ids[i]] = struct{}{}
	}
	return etuo
}

// AddPortDefinitions adds the port_definitions edges to EquipmentPortDefinition.
func (etuo *EquipmentTypeUpdateOne) AddPortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddPortDefinitionIDs(ids...)
}

// AddPositionDefinitionIDs adds the position_definitions edge to EquipmentPositionDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) AddPositionDefinitionIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.position_definitions == nil {
		etuo.position_definitions = make(map[string]struct{})
	}
	for i := range ids {
		etuo.position_definitions[ids[i]] = struct{}{}
	}
	return etuo
}

// AddPositionDefinitions adds the position_definitions edges to EquipmentPositionDefinition.
func (etuo *EquipmentTypeUpdateOne) AddPositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddPositionDefinitionIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (etuo *EquipmentTypeUpdateOne) AddPropertyTypeIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.property_types == nil {
		etuo.property_types = make(map[string]struct{})
	}
	for i := range ids {
		etuo.property_types[ids[i]] = struct{}{}
	}
	return etuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (etuo *EquipmentTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *EquipmentTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etuo.AddPropertyTypeIDs(ids...)
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (etuo *EquipmentTypeUpdateOne) AddEquipmentIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.equipment == nil {
		etuo.equipment = make(map[string]struct{})
	}
	for i := range ids {
		etuo.equipment[ids[i]] = struct{}{}
	}
	return etuo
}

// AddEquipment adds the equipment edges to Equipment.
func (etuo *EquipmentTypeUpdateOne) AddEquipment(e ...*Equipment) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.AddEquipmentIDs(ids...)
}

// SetCategoryID sets the category edge to EquipmentCategory by id.
func (etuo *EquipmentTypeUpdateOne) SetCategoryID(id string) *EquipmentTypeUpdateOne {
	if etuo.category == nil {
		etuo.category = make(map[string]struct{})
	}
	etuo.category[id] = struct{}{}
	return etuo
}

// SetNillableCategoryID sets the category edge to EquipmentCategory by id if the given value is not nil.
func (etuo *EquipmentTypeUpdateOne) SetNillableCategoryID(id *string) *EquipmentTypeUpdateOne {
	if id != nil {
		etuo = etuo.SetCategoryID(*id)
	}
	return etuo
}

// SetCategory sets the category edge to EquipmentCategory.
func (etuo *EquipmentTypeUpdateOne) SetCategory(e *EquipmentCategory) *EquipmentTypeUpdateOne {
	return etuo.SetCategoryID(e.ID)
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePortDefinitionIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.removedPortDefinitions == nil {
		etuo.removedPortDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		etuo.removedPortDefinitions[ids[i]] = struct{}{}
	}
	return etuo
}

// RemovePortDefinitions removes port_definitions edges to EquipmentPortDefinition.
func (etuo *EquipmentTypeUpdateOne) RemovePortDefinitions(e ...*EquipmentPortDefinition) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemovePortDefinitionIDs(ids...)
}

// RemovePositionDefinitionIDs removes the position_definitions edge to EquipmentPositionDefinition by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePositionDefinitionIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.removedPositionDefinitions == nil {
		etuo.removedPositionDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		etuo.removedPositionDefinitions[ids[i]] = struct{}{}
	}
	return etuo
}

// RemovePositionDefinitions removes position_definitions edges to EquipmentPositionDefinition.
func (etuo *EquipmentTypeUpdateOne) RemovePositionDefinitions(e ...*EquipmentPositionDefinition) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemovePositionDefinitionIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (etuo *EquipmentTypeUpdateOne) RemovePropertyTypeIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.removedPropertyTypes == nil {
		etuo.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		etuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return etuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (etuo *EquipmentTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *EquipmentTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return etuo.RemovePropertyTypeIDs(ids...)
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (etuo *EquipmentTypeUpdateOne) RemoveEquipmentIDs(ids ...string) *EquipmentTypeUpdateOne {
	if etuo.removedEquipment == nil {
		etuo.removedEquipment = make(map[string]struct{})
	}
	for i := range ids {
		etuo.removedEquipment[ids[i]] = struct{}{}
	}
	return etuo
}

// RemoveEquipment removes equipment edges to Equipment.
func (etuo *EquipmentTypeUpdateOne) RemoveEquipment(e ...*Equipment) *EquipmentTypeUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return etuo.RemoveEquipmentIDs(ids...)
}

// ClearCategory clears the category edge to EquipmentCategory.
func (etuo *EquipmentTypeUpdateOne) ClearCategory() *EquipmentTypeUpdateOne {
	etuo.clearedCategory = true
	return etuo
}

// Save executes the query and returns the updated entity.
func (etuo *EquipmentTypeUpdateOne) Save(ctx context.Context) (*EquipmentType, error) {
	if etuo.update_time == nil {
		v := equipmenttype.UpdateDefaultUpdateTime()
		etuo.update_time = &v
	}
	if len(etuo.category) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return etuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (etuo *EquipmentTypeUpdateOne) SaveX(ctx context.Context) *EquipmentType {
	et, err := etuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return et
}

// Exec executes the query on the entity.
func (etuo *EquipmentTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := etuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (etuo *EquipmentTypeUpdateOne) ExecX(ctx context.Context) {
	if err := etuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (etuo *EquipmentTypeUpdateOne) sqlSave(ctx context.Context) (et *EquipmentType, err error) {
	var (
		builder  = sql.Dialect(etuo.driver.Dialect())
		selector = builder.Select(equipmenttype.Columns...).From(builder.Table(equipmenttype.Table))
	)
	equipmenttype.ID(etuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = etuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		et = &EquipmentType{config: etuo.config}
		if err := et.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into EquipmentType: %v", err)
		}
		id = et.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("EquipmentType with id: %v", etuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one EquipmentType with the same id: %v", etuo.id)
	}

	tx, err := etuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(equipmenttype.Table)
	)
	updater = updater.Where(sql.InInts(equipmenttype.FieldID, ids...))
	if value := etuo.update_time; value != nil {
		updater.Set(equipmenttype.FieldUpdateTime, *value)
		et.UpdateTime = *value
	}
	if value := etuo.name; value != nil {
		updater.Set(equipmenttype.FieldName, *value)
		et.Name = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(etuo.removedPortDefinitions) > 0 {
		eids := make([]int, len(etuo.removedPortDefinitions))
		for eid := range etuo.removedPortDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmenttype.PortDefinitionsTable).
			SetNull(equipmenttype.PortDefinitionsColumn).
			Where(sql.InInts(equipmenttype.PortDefinitionsColumn, ids...)).
			Where(sql.InInts(equipmentportdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(etuo.port_definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range etuo.port_definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentportdefinition.FieldID, eid)
			}
			query, args := builder.Update(equipmenttype.PortDefinitionsTable).
				Set(equipmenttype.PortDefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmenttype.PortDefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(etuo.port_definitions) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"port_definitions\" %v already connected to a different \"EquipmentType\"", keys(etuo.port_definitions))})
			}
		}
	}
	if len(etuo.removedPositionDefinitions) > 0 {
		eids := make([]int, len(etuo.removedPositionDefinitions))
		for eid := range etuo.removedPositionDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmenttype.PositionDefinitionsTable).
			SetNull(equipmenttype.PositionDefinitionsColumn).
			Where(sql.InInts(equipmenttype.PositionDefinitionsColumn, ids...)).
			Where(sql.InInts(equipmentpositiondefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(etuo.position_definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range etuo.position_definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentpositiondefinition.FieldID, eid)
			}
			query, args := builder.Update(equipmenttype.PositionDefinitionsTable).
				Set(equipmenttype.PositionDefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(equipmenttype.PositionDefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(etuo.position_definitions) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"position_definitions\" %v already connected to a different \"EquipmentType\"", keys(etuo.position_definitions))})
			}
		}
	}
	if len(etuo.removedPropertyTypes) > 0 {
		eids := make([]int, len(etuo.removedPropertyTypes))
		for eid := range etuo.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmenttype.PropertyTypesTable).
			SetNull(equipmenttype.PropertyTypesColumn).
			Where(sql.InInts(equipmenttype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(etuo.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range etuo.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(equipmenttype.PropertyTypesTable).
				Set(equipmenttype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(equipmenttype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(etuo.property_types) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"EquipmentType\"", keys(etuo.property_types))})
			}
		}
	}
	if len(etuo.removedEquipment) > 0 {
		eids := make([]int, len(etuo.removedEquipment))
		for eid := range etuo.removedEquipment {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(equipmenttype.EquipmentTable).
			SetNull(equipmenttype.EquipmentColumn).
			Where(sql.InInts(equipmenttype.EquipmentColumn, ids...)).
			Where(sql.InInts(equipment.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(etuo.equipment) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range etuo.equipment {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipment.FieldID, eid)
			}
			query, args := builder.Update(equipmenttype.EquipmentTable).
				Set(equipmenttype.EquipmentColumn, id).
				Where(sql.And(p, sql.IsNull(equipmenttype.EquipmentColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(etuo.equipment) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"equipment\" %v already connected to a different \"EquipmentType\"", keys(etuo.equipment))})
			}
		}
	}
	if etuo.clearedCategory {
		query, args := builder.Update(equipmenttype.CategoryTable).
			SetNull(equipmenttype.CategoryColumn).
			Where(sql.InInts(equipmentcategory.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(etuo.category) > 0 {
		for eid := range etuo.category {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(equipmenttype.CategoryTable).
				Set(equipmenttype.CategoryColumn, eid).
				Where(sql.InInts(equipmenttype.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return et, nil
}
