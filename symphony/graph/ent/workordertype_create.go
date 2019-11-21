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
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
)

// WorkOrderTypeCreate is the builder for creating a WorkOrderType entity.
type WorkOrderTypeCreate struct {
	config
	create_time            *time.Time
	update_time            *time.Time
	name                   *string
	description            *string
	work_orders            map[string]struct{}
	property_types         map[string]struct{}
	definitions            map[string]struct{}
	check_list_definitions map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (wotc *WorkOrderTypeCreate) SetCreateTime(t time.Time) *WorkOrderTypeCreate {
	wotc.create_time = &t
	return wotc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (wotc *WorkOrderTypeCreate) SetNillableCreateTime(t *time.Time) *WorkOrderTypeCreate {
	if t != nil {
		wotc.SetCreateTime(*t)
	}
	return wotc
}

// SetUpdateTime sets the update_time field.
func (wotc *WorkOrderTypeCreate) SetUpdateTime(t time.Time) *WorkOrderTypeCreate {
	wotc.update_time = &t
	return wotc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (wotc *WorkOrderTypeCreate) SetNillableUpdateTime(t *time.Time) *WorkOrderTypeCreate {
	if t != nil {
		wotc.SetUpdateTime(*t)
	}
	return wotc
}

// SetName sets the name field.
func (wotc *WorkOrderTypeCreate) SetName(s string) *WorkOrderTypeCreate {
	wotc.name = &s
	return wotc
}

// SetDescription sets the description field.
func (wotc *WorkOrderTypeCreate) SetDescription(s string) *WorkOrderTypeCreate {
	wotc.description = &s
	return wotc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (wotc *WorkOrderTypeCreate) SetNillableDescription(s *string) *WorkOrderTypeCreate {
	if s != nil {
		wotc.SetDescription(*s)
	}
	return wotc
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (wotc *WorkOrderTypeCreate) AddWorkOrderIDs(ids ...string) *WorkOrderTypeCreate {
	if wotc.work_orders == nil {
		wotc.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		wotc.work_orders[ids[i]] = struct{}{}
	}
	return wotc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (wotc *WorkOrderTypeCreate) AddWorkOrders(w ...*WorkOrder) *WorkOrderTypeCreate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotc.AddWorkOrderIDs(ids...)
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (wotc *WorkOrderTypeCreate) AddPropertyTypeIDs(ids ...string) *WorkOrderTypeCreate {
	if wotc.property_types == nil {
		wotc.property_types = make(map[string]struct{})
	}
	for i := range ids {
		wotc.property_types[ids[i]] = struct{}{}
	}
	return wotc
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (wotc *WorkOrderTypeCreate) AddPropertyTypes(p ...*PropertyType) *WorkOrderTypeCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return wotc.AddPropertyTypeIDs(ids...)
}

// AddDefinitionIDs adds the definitions edge to WorkOrderDefinition by ids.
func (wotc *WorkOrderTypeCreate) AddDefinitionIDs(ids ...string) *WorkOrderTypeCreate {
	if wotc.definitions == nil {
		wotc.definitions = make(map[string]struct{})
	}
	for i := range ids {
		wotc.definitions[ids[i]] = struct{}{}
	}
	return wotc
}

// AddDefinitions adds the definitions edges to WorkOrderDefinition.
func (wotc *WorkOrderTypeCreate) AddDefinitions(w ...*WorkOrderDefinition) *WorkOrderTypeCreate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return wotc.AddDefinitionIDs(ids...)
}

// AddCheckListDefinitionIDs adds the check_list_definitions edge to CheckListItemDefinition by ids.
func (wotc *WorkOrderTypeCreate) AddCheckListDefinitionIDs(ids ...string) *WorkOrderTypeCreate {
	if wotc.check_list_definitions == nil {
		wotc.check_list_definitions = make(map[string]struct{})
	}
	for i := range ids {
		wotc.check_list_definitions[ids[i]] = struct{}{}
	}
	return wotc
}

// AddCheckListDefinitions adds the check_list_definitions edges to CheckListItemDefinition.
func (wotc *WorkOrderTypeCreate) AddCheckListDefinitions(c ...*CheckListItemDefinition) *WorkOrderTypeCreate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wotc.AddCheckListDefinitionIDs(ids...)
}

// Save creates the WorkOrderType in the database.
func (wotc *WorkOrderTypeCreate) Save(ctx context.Context) (*WorkOrderType, error) {
	if wotc.create_time == nil {
		v := workordertype.DefaultCreateTime()
		wotc.create_time = &v
	}
	if wotc.update_time == nil {
		v := workordertype.DefaultUpdateTime()
		wotc.update_time = &v
	}
	if wotc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	return wotc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (wotc *WorkOrderTypeCreate) SaveX(ctx context.Context) *WorkOrderType {
	v, err := wotc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (wotc *WorkOrderTypeCreate) sqlSave(ctx context.Context) (*WorkOrderType, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(wotc.driver.Dialect())
		wot     = &WorkOrderType{config: wotc.config}
	)
	tx, err := wotc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(workordertype.Table).Default()
	if value := wotc.create_time; value != nil {
		insert.Set(workordertype.FieldCreateTime, *value)
		wot.CreateTime = *value
	}
	if value := wotc.update_time; value != nil {
		insert.Set(workordertype.FieldUpdateTime, *value)
		wot.UpdateTime = *value
	}
	if value := wotc.name; value != nil {
		insert.Set(workordertype.FieldName, *value)
		wot.Name = *value
	}
	if value := wotc.description; value != nil {
		insert.Set(workordertype.FieldDescription, *value)
		wot.Description = *value
	}
	id, err := insertLastID(ctx, tx, insert.Returning(workordertype.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	wot.ID = strconv.FormatInt(id, 10)
	if len(wotc.work_orders) > 0 {
		p := sql.P()
		for eid := range wotc.work_orders {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(wotc.work_orders) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"WorkOrderType\"", keys(wotc.work_orders))})
		}
	}
	if len(wotc.property_types) > 0 {
		p := sql.P()
		for eid := range wotc.property_types {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(wotc.property_types) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"property_types\" %v already connected to a different \"WorkOrderType\"", keys(wotc.property_types))})
		}
	}
	if len(wotc.definitions) > 0 {
		p := sql.P()
		for eid := range wotc.definitions {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(wotc.definitions) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"definitions\" %v already connected to a different \"WorkOrderType\"", keys(wotc.definitions))})
		}
	}
	if len(wotc.check_list_definitions) > 0 {
		p := sql.P()
		for eid := range wotc.check_list_definitions {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(wotc.check_list_definitions) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"check_list_definitions\" %v already connected to a different \"WorkOrderType\"", keys(wotc.check_list_definitions))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return wot, nil
}
