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
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderTypeUpdate is the builder for updating WorkOrderType entities.
type WorkOrderTypeUpdate struct {
	config

	update_time                 *time.Time
	name                        *string
	description                 *string
	cleardescription            bool
	work_orders                 map[string]struct{}
	property_types              map[string]struct{}
	definitions                 map[string]struct{}
	check_list_definitions      map[string]struct{}
	removedWorkOrders           map[string]struct{}
	removedPropertyTypes        map[string]struct{}
	removedDefinitions          map[string]struct{}
	removedCheckListDefinitions map[string]struct{}
	predicates                  []predicate.WorkOrderType
}

// Where adds a new predicate for the builder.
func (wotu *WorkOrderTypeUpdate) Where(ps ...predicate.WorkOrderType) *WorkOrderTypeUpdate {
	wotu.predicates = append(wotu.predicates, ps...)
	return wotu
}

// SetName sets the name field.
func (wotu *WorkOrderTypeUpdate) SetName(s string) *WorkOrderTypeUpdate {
	wotu.name = &s
	return wotu
}

// SetDescription sets the description field.
func (wotu *WorkOrderTypeUpdate) SetDescription(s string) *WorkOrderTypeUpdate {
	wotu.description = &s
	return wotu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wotu *WorkOrderTypeUpdate) SetNillableDescription(s *string) *WorkOrderTypeUpdate {
	if s != nil {
		wotu.SetDescription(*s)
	}
	return wotu
}

// ClearDescription clears the value of description.
func (wotu *WorkOrderTypeUpdate) ClearDescription() *WorkOrderTypeUpdate {
	wotu.description = nil
	wotu.cleardescription = true
	return wotu
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (wotu *WorkOrderTypeUpdate) AddWorkOrderIDs(ids ...string) *WorkOrderTypeUpdate {
	if wotu.work_orders == nil {
		wotu.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		wotu.work_orders[ids[i]] = struct{}{}
	}
	return wotu
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (wotu *WorkOrderTypeUpdate) AddWorkOrders(w ...*WorkOrder) *WorkOrderTypeUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotu.AddWorkOrderIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (wotu *WorkOrderTypeUpdate) AddPropertyTypeIDs(ids ...string) *WorkOrderTypeUpdate {
	if wotu.property_types == nil {
		wotu.property_types = make(map[string]struct{})
	}
	for i := range ids {
		wotu.property_types[ids[i]] = struct{}{}
	}
	return wotu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (wotu *WorkOrderTypeUpdate) AddPropertyTypes(p ...*PropertyType) *WorkOrderTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotu.AddPropertyTypeIDs(ids...)
}

// AddDefinitionIDs adds the definitions edge to WorkOrderDefinition by ids.
func (wotu *WorkOrderTypeUpdate) AddDefinitionIDs(ids ...string) *WorkOrderTypeUpdate {
	if wotu.definitions == nil {
		wotu.definitions = make(map[string]struct{})
	}
	for i := range ids {
		wotu.definitions[ids[i]] = struct{}{}
	}
	return wotu
}

// AddDefinitions adds the definitions edges to WorkOrderDefinition.
func (wotu *WorkOrderTypeUpdate) AddDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotu.AddDefinitionIDs(ids...)
}

// AddCheckListDefinitionIDs adds the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotu *WorkOrderTypeUpdate) AddCheckListDefinitionIDs(ids ...string) *WorkOrderTypeUpdate {
	if wotu.check_list_definitions == nil {
		wotu.check_list_definitions = make(map[string]struct{})
	}
	for i := range ids {
		wotu.check_list_definitions[ids[i]] = struct{}{}
	}
	return wotu
}

// AddCheckListDefinitions adds the check_list_definitions edges to CheckListItemDefinition.
func (wotu *WorkOrderTypeUpdate) AddCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotu.AddCheckListDefinitionIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (wotu *WorkOrderTypeUpdate) RemoveWorkOrderIDs(ids ...string) *WorkOrderTypeUpdate {
	if wotu.removedWorkOrders == nil {
		wotu.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		wotu.removedWorkOrders[ids[i]] = struct{}{}
	}
	return wotu
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (wotu *WorkOrderTypeUpdate) RemoveWorkOrders(w ...*WorkOrder) *WorkOrderTypeUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotu.RemoveWorkOrderIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (wotu *WorkOrderTypeUpdate) RemovePropertyTypeIDs(ids ...string) *WorkOrderTypeUpdate {
	if wotu.removedPropertyTypes == nil {
		wotu.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		wotu.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return wotu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (wotu *WorkOrderTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *WorkOrderTypeUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotu.RemovePropertyTypeIDs(ids...)
}

// RemoveDefinitionIDs removes the definitions edge to WorkOrderDefinition by ids.
func (wotu *WorkOrderTypeUpdate) RemoveDefinitionIDs(ids ...string) *WorkOrderTypeUpdate {
	if wotu.removedDefinitions == nil {
		wotu.removedDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		wotu.removedDefinitions[ids[i]] = struct{}{}
	}
	return wotu
}

// RemoveDefinitions removes definitions edges to WorkOrderDefinition.
func (wotu *WorkOrderTypeUpdate) RemoveDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotu.RemoveDefinitionIDs(ids...)
}

// RemoveCheckListDefinitionIDs removes the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotu *WorkOrderTypeUpdate) RemoveCheckListDefinitionIDs(ids ...string) *WorkOrderTypeUpdate {
	if wotu.removedCheckListDefinitions == nil {
		wotu.removedCheckListDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		wotu.removedCheckListDefinitions[ids[i]] = struct{}{}
	}
	return wotu
}

// RemoveCheckListDefinitions removes check_list_definitions edges to CheckListItemDefinition.
func (wotu *WorkOrderTypeUpdate) RemoveCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotu.RemoveCheckListDefinitionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (wotu *WorkOrderTypeUpdate) Save(ctx context.Context) (int, error) {
	if wotu.update_time == nil {
		v := workordertype.UpdateDefaultUpdateTime()
		wotu.update_time = &v
	}
	return wotu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (wotu *WorkOrderTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := wotu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wotu *WorkOrderTypeUpdate) Exec(ctx context.Context) error {
	_, err := wotu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wotu *WorkOrderTypeUpdate) ExecX(ctx context.Context) {
	if err := wotu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wotu *WorkOrderTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(wotu.driver.Dialect())
		selector = builder.Select(workordertype.FieldID).From(builder.Table(workordertype.Table))
	)
	for _, p := range wotu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = wotu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := wotu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(workordertype.Table)
	)
	updater = updater.Where(sql.InInts(workordertype.FieldID, ids...))
	if value := wotu.update_time; value != nil {
		updater.Set(workordertype.FieldUpdateTime, *value)
	}
	if value := wotu.name; value != nil {
		updater.Set(workordertype.FieldName, *value)
	}
	if value := wotu.description; value != nil {
		updater.Set(workordertype.FieldDescription, *value)
	}
	if wotu.cleardescription {
		updater.SetNull(workordertype.FieldDescription)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(wotu.removedWorkOrders) > 0 {
		eids := make([]int, len(wotu.removedWorkOrders))
		for eid := range wotu.removedWorkOrders {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(workordertype.WorkOrdersTable).
			SetNull(workordertype.WorkOrdersColumn).
			Where(sql.InInts(workordertype.WorkOrdersColumn, ids...)).
			Where(sql.InInts(workorder.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(wotu.work_orders) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range wotu.work_orders {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(workorder.FieldID, eid)
			}
			query, args := builder.Update(workordertype.WorkOrdersTable).
				Set(workordertype.WorkOrdersColumn, id).
				Where(sql.And(p, sql.IsNull(workordertype.WorkOrdersColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(wotu.work_orders) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"WorkOrderType\"", keys(wotu.work_orders))})
			}
		}
	}
	if len(wotu.removedPropertyTypes) > 0 {
		eids := make([]int, len(wotu.removedPropertyTypes))
		for eid := range wotu.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(workordertype.PropertyTypesTable).
			SetNull(workordertype.PropertyTypesColumn).
			Where(sql.InInts(workordertype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(wotu.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range wotu.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(workordertype.PropertyTypesTable).
				Set(workordertype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(workordertype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(wotu.property_types) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"WorkOrderType\"", keys(wotu.property_types))})
			}
		}
	}
	if len(wotu.removedDefinitions) > 0 {
		eids := make([]int, len(wotu.removedDefinitions))
		for eid := range wotu.removedDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(workordertype.DefinitionsTable).
			SetNull(workordertype.DefinitionsColumn).
			Where(sql.InInts(workordertype.DefinitionsColumn, ids...)).
			Where(sql.InInts(workorderdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(wotu.definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range wotu.definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(workorderdefinition.FieldID, eid)
			}
			query, args := builder.Update(workordertype.DefinitionsTable).
				Set(workordertype.DefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(workordertype.DefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(wotu.definitions) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"definitions\" %v already connected to a different \"WorkOrderType\"", keys(wotu.definitions))})
			}
		}
	}
	if len(wotu.removedCheckListDefinitions) > 0 {
		eids := make([]int, len(wotu.removedCheckListDefinitions))
		for eid := range wotu.removedCheckListDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(workordertype.CheckListDefinitionsTable).
			SetNull(workordertype.CheckListDefinitionsColumn).
			Where(sql.InInts(workordertype.CheckListDefinitionsColumn, ids...)).
			Where(sql.InInts(checklistitemdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(wotu.check_list_definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range wotu.check_list_definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(checklistitemdefinition.FieldID, eid)
			}
			query, args := builder.Update(workordertype.CheckListDefinitionsTable).
				Set(workordertype.CheckListDefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(workordertype.CheckListDefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(wotu.check_list_definitions) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"check_list_definitions\" %v already connected to a different \"WorkOrderType\"", keys(wotu.check_list_definitions))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// WorkOrderTypeUpdateOne is the builder for updating a single WorkOrderType entity.
type WorkOrderTypeUpdateOne struct {
	config
	id string

	update_time                 *time.Time
	name                        *string
	description                 *string
	cleardescription            bool
	work_orders                 map[string]struct{}
	property_types              map[string]struct{}
	definitions                 map[string]struct{}
	check_list_definitions      map[string]struct{}
	removedWorkOrders           map[string]struct{}
	removedPropertyTypes        map[string]struct{}
	removedDefinitions          map[string]struct{}
	removedCheckListDefinitions map[string]struct{}
}

// SetName sets the name field.
func (wotuo *WorkOrderTypeUpdateOne) SetName(s string) *WorkOrderTypeUpdateOne {
	wotuo.name = &s
	return wotuo
}

// SetDescription sets the description field.
func (wotuo *WorkOrderTypeUpdateOne) SetDescription(s string) *WorkOrderTypeUpdateOne {
	wotuo.description = &s
	return wotuo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wotuo *WorkOrderTypeUpdateOne) SetNillableDescription(s *string) *WorkOrderTypeUpdateOne {
	if s != nil {
		wotuo.SetDescription(*s)
	}
	return wotuo
}

// ClearDescription clears the value of description.
func (wotuo *WorkOrderTypeUpdateOne) ClearDescription() *WorkOrderTypeUpdateOne {
	wotuo.description = nil
	wotuo.cleardescription = true
	return wotuo
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddWorkOrderIDs(ids ...string) *WorkOrderTypeUpdateOne {
	if wotuo.work_orders == nil {
		wotuo.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		wotuo.work_orders[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (wotuo *WorkOrderTypeUpdateOne) AddWorkOrders(w ...*WorkOrder) *WorkOrderTypeUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotuo.AddWorkOrderIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddPropertyTypeIDs(ids ...string) *WorkOrderTypeUpdateOne {
	if wotuo.property_types == nil {
		wotuo.property_types = make(map[string]struct{})
	}
	for i := range ids {
		wotuo.property_types[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (wotuo *WorkOrderTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *WorkOrderTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotuo.AddPropertyTypeIDs(ids...)
}

// AddDefinitionIDs adds the definitions edge to WorkOrderDefinition by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddDefinitionIDs(ids ...string) *WorkOrderTypeUpdateOne {
	if wotuo.definitions == nil {
		wotuo.definitions = make(map[string]struct{})
	}
	for i := range ids {
		wotuo.definitions[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddDefinitions adds the definitions edges to WorkOrderDefinition.
func (wotuo *WorkOrderTypeUpdateOne) AddDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotuo.AddDefinitionIDs(ids...)
}

// AddCheckListDefinitionIDs adds the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotuo *WorkOrderTypeUpdateOne) AddCheckListDefinitionIDs(ids ...string) *WorkOrderTypeUpdateOne {
	if wotuo.check_list_definitions == nil {
		wotuo.check_list_definitions = make(map[string]struct{})
	}
	for i := range ids {
		wotuo.check_list_definitions[ids[i]] = struct{}{}
	}
	return wotuo
}

// AddCheckListDefinitions adds the check_list_definitions edges to CheckListItemDefinition.
func (wotuo *WorkOrderTypeUpdateOne) AddCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotuo.AddCheckListDefinitionIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemoveWorkOrderIDs(ids ...string) *WorkOrderTypeUpdateOne {
	if wotuo.removedWorkOrders == nil {
		wotuo.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		wotuo.removedWorkOrders[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (wotuo *WorkOrderTypeUpdateOne) RemoveWorkOrders(w ...*WorkOrder) *WorkOrderTypeUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotuo.RemoveWorkOrderIDs(ids...)
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemovePropertyTypeIDs(ids ...string) *WorkOrderTypeUpdateOne {
	if wotuo.removedPropertyTypes == nil {
		wotuo.removedPropertyTypes = make(map[string]struct{})
	}
	for i := range ids {
		wotuo.removedPropertyTypes[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (wotuo *WorkOrderTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *WorkOrderTypeUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotuo.RemovePropertyTypeIDs(ids...)
}

// RemoveDefinitionIDs removes the definitions edge to WorkOrderDefinition by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemoveDefinitionIDs(ids ...string) *WorkOrderTypeUpdateOne {
	if wotuo.removedDefinitions == nil {
		wotuo.removedDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		wotuo.removedDefinitions[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemoveDefinitions removes definitions edges to WorkOrderDefinition.
func (wotuo *WorkOrderTypeUpdateOne) RemoveDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotuo.RemoveDefinitionIDs(ids...)
}

// RemoveCheckListDefinitionIDs removes the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotuo *WorkOrderTypeUpdateOne) RemoveCheckListDefinitionIDs(ids ...string) *WorkOrderTypeUpdateOne {
	if wotuo.removedCheckListDefinitions == nil {
		wotuo.removedCheckListDefinitions = make(map[string]struct{})
	}
	for i := range ids {
		wotuo.removedCheckListDefinitions[ids[i]] = struct{}{}
	}
	return wotuo
}

// RemoveCheckListDefinitions removes check_list_definitions edges to CheckListItemDefinition.
func (wotuo *WorkOrderTypeUpdateOne) RemoveCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotuo.RemoveCheckListDefinitionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (wotuo *WorkOrderTypeUpdateOne) Save(ctx context.Context) (*WorkOrderType, error) {
	if wotuo.update_time == nil {
		v := workordertype.UpdateDefaultUpdateTime()
		wotuo.update_time = &v
	}
	return wotuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (wotuo *WorkOrderTypeUpdateOne) SaveX(ctx context.Context) *WorkOrderType {
	wot, err := wotuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return wot
}

// Exec executes the query on the entity.
func (wotuo *WorkOrderTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := wotuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wotuo *WorkOrderTypeUpdateOne) ExecX(ctx context.Context) {
	if err := wotuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (wotuo *WorkOrderTypeUpdateOne) sqlSave(ctx context.Context) (wot *WorkOrderType, err error) {
	var (
		builder  = sql.Dialect(wotuo.driver.Dialect())
		selector = builder.Select(workordertype.Columns...).From(builder.Table(workordertype.Table))
	)
	workordertype.ID(wotuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = wotuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		wot = &WorkOrderType{config: wotuo.config}
		if err := wot.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into WorkOrderType: %v", err)
		}
		id = wot.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("WorkOrderType with id: %v", wotuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one WorkOrderType with the same id: %v", wotuo.id)
	}

	tx, err := wotuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(workordertype.Table)
	)
	updater = updater.Where(sql.InInts(workordertype.FieldID, ids...))
	if value := wotuo.update_time; value != nil {
		updater.Set(workordertype.FieldUpdateTime, *value)
		wot.UpdateTime = *value
	}
	if value := wotuo.name; value != nil {
		updater.Set(workordertype.FieldName, *value)
		wot.Name = *value
	}
	if value := wotuo.description; value != nil {
		updater.Set(workordertype.FieldDescription, *value)
		wot.Description = *value
	}
	if wotuo.cleardescription {
		var value string
		wot.Description = value
		updater.SetNull(workordertype.FieldDescription)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(wotuo.removedWorkOrders) > 0 {
		eids := make([]int, len(wotuo.removedWorkOrders))
		for eid := range wotuo.removedWorkOrders {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(workordertype.WorkOrdersTable).
			SetNull(workordertype.WorkOrdersColumn).
			Where(sql.InInts(workordertype.WorkOrdersColumn, ids...)).
			Where(sql.InInts(workorder.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(wotuo.work_orders) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range wotuo.work_orders {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(workorder.FieldID, eid)
			}
			query, args := builder.Update(workordertype.WorkOrdersTable).
				Set(workordertype.WorkOrdersColumn, id).
				Where(sql.And(p, sql.IsNull(workordertype.WorkOrdersColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(wotuo.work_orders) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"WorkOrderType\"", keys(wotuo.work_orders))})
			}
		}
	}
	if len(wotuo.removedPropertyTypes) > 0 {
		eids := make([]int, len(wotuo.removedPropertyTypes))
		for eid := range wotuo.removedPropertyTypes {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(workordertype.PropertyTypesTable).
			SetNull(workordertype.PropertyTypesColumn).
			Where(sql.InInts(workordertype.PropertyTypesColumn, ids...)).
			Where(sql.InInts(propertytype.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(wotuo.property_types) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range wotuo.property_types {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(propertytype.FieldID, eid)
			}
			query, args := builder.Update(workordertype.PropertyTypesTable).
				Set(workordertype.PropertyTypesColumn, id).
				Where(sql.And(p, sql.IsNull(workordertype.PropertyTypesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(wotuo.property_types) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"WorkOrderType\"", keys(wotuo.property_types))})
			}
		}
	}
	if len(wotuo.removedDefinitions) > 0 {
		eids := make([]int, len(wotuo.removedDefinitions))
		for eid := range wotuo.removedDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(workordertype.DefinitionsTable).
			SetNull(workordertype.DefinitionsColumn).
			Where(sql.InInts(workordertype.DefinitionsColumn, ids...)).
			Where(sql.InInts(workorderdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(wotuo.definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range wotuo.definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(workorderdefinition.FieldID, eid)
			}
			query, args := builder.Update(workordertype.DefinitionsTable).
				Set(workordertype.DefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(workordertype.DefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(wotuo.definitions) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"definitions\" %v already connected to a different \"WorkOrderType\"", keys(wotuo.definitions))})
			}
		}
	}
	if len(wotuo.removedCheckListDefinitions) > 0 {
		eids := make([]int, len(wotuo.removedCheckListDefinitions))
		for eid := range wotuo.removedCheckListDefinitions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(workordertype.CheckListDefinitionsTable).
			SetNull(workordertype.CheckListDefinitionsColumn).
			Where(sql.InInts(workordertype.CheckListDefinitionsColumn, ids...)).
			Where(sql.InInts(checklistitemdefinition.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(wotuo.check_list_definitions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range wotuo.check_list_definitions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(checklistitemdefinition.FieldID, eid)
			}
			query, args := builder.Update(workordertype.CheckListDefinitionsTable).
				Set(workordertype.CheckListDefinitionsColumn, id).
				Where(sql.And(p, sql.IsNull(workordertype.CheckListDefinitionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(wotuo.check_list_definitions) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"check_list_definitions\" %v already connected to a different \"WorkOrderType\"", keys(wotuo.check_list_definitions))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return wot, nil
}
