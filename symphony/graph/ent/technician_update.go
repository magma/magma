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
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/technician"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// TechnicianUpdate is the builder for updating Technician entities.
type TechnicianUpdate struct {
	config

	update_time       *time.Time
	name              *string
	email             *string
	work_orders       map[string]struct{}
	removedWorkOrders map[string]struct{}
	predicates        []predicate.Technician
}

// Where adds a new predicate for the builder.
func (tu *TechnicianUpdate) Where(ps ...predicate.Technician) *TechnicianUpdate {
	tu.predicates = append(tu.predicates, ps...)
	return tu
}

// SetName sets the name field.
func (tu *TechnicianUpdate) SetName(s string) *TechnicianUpdate {
	tu.name = &s
	return tu
}

// SetEmail sets the email field.
func (tu *TechnicianUpdate) SetEmail(s string) *TechnicianUpdate {
	tu.email = &s
	return tu
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (tu *TechnicianUpdate) AddWorkOrderIDs(ids ...string) *TechnicianUpdate {
	if tu.work_orders == nil {
		tu.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		tu.work_orders[ids[i]] = struct{}{}
	}
	return tu
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (tu *TechnicianUpdate) AddWorkOrders(w ...*WorkOrder) *TechnicianUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tu.AddWorkOrderIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (tu *TechnicianUpdate) RemoveWorkOrderIDs(ids ...string) *TechnicianUpdate {
	if tu.removedWorkOrders == nil {
		tu.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		tu.removedWorkOrders[ids[i]] = struct{}{}
	}
	return tu
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (tu *TechnicianUpdate) RemoveWorkOrders(w ...*WorkOrder) *TechnicianUpdate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tu.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (tu *TechnicianUpdate) Save(ctx context.Context) (int, error) {
	if tu.update_time == nil {
		v := technician.UpdateDefaultUpdateTime()
		tu.update_time = &v
	}
	if tu.name != nil {
		if err := technician.NameValidator(*tu.name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if tu.email != nil {
		if err := technician.EmailValidator(*tu.email); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	return tu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TechnicianUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TechnicianUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TechnicianUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TechnicianUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(tu.driver.Dialect())
		selector = builder.Select(technician.FieldID).From(builder.Table(technician.Table))
	)
	for _, p := range tu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = tu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := tu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(technician.Table)
	)
	updater = updater.Where(sql.InInts(technician.FieldID, ids...))
	if value := tu.update_time; value != nil {
		updater.Set(technician.FieldUpdateTime, *value)
	}
	if value := tu.name; value != nil {
		updater.Set(technician.FieldName, *value)
	}
	if value := tu.email; value != nil {
		updater.Set(technician.FieldEmail, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(tu.removedWorkOrders) > 0 {
		eids := make([]int, len(tu.removedWorkOrders))
		for eid := range tu.removedWorkOrders {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(technician.WorkOrdersTable).
			SetNull(technician.WorkOrdersColumn).
			Where(sql.InInts(technician.WorkOrdersColumn, ids...)).
			Where(sql.InInts(workorder.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(tu.work_orders) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range tu.work_orders {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(workorder.FieldID, eid)
			}
			query, args := builder.Update(technician.WorkOrdersTable).
				Set(technician.WorkOrdersColumn, id).
				Where(sql.And(p, sql.IsNull(technician.WorkOrdersColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(tu.work_orders) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"Technician\"", keys(tu.work_orders))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// TechnicianUpdateOne is the builder for updating a single Technician entity.
type TechnicianUpdateOne struct {
	config
	id string

	update_time       *time.Time
	name              *string
	email             *string
	work_orders       map[string]struct{}
	removedWorkOrders map[string]struct{}
}

// SetName sets the name field.
func (tuo *TechnicianUpdateOne) SetName(s string) *TechnicianUpdateOne {
	tuo.name = &s
	return tuo
}

// SetEmail sets the email field.
func (tuo *TechnicianUpdateOne) SetEmail(s string) *TechnicianUpdateOne {
	tuo.email = &s
	return tuo
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (tuo *TechnicianUpdateOne) AddWorkOrderIDs(ids ...string) *TechnicianUpdateOne {
	if tuo.work_orders == nil {
		tuo.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		tuo.work_orders[ids[i]] = struct{}{}
	}
	return tuo
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (tuo *TechnicianUpdateOne) AddWorkOrders(w ...*WorkOrder) *TechnicianUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tuo.AddWorkOrderIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (tuo *TechnicianUpdateOne) RemoveWorkOrderIDs(ids ...string) *TechnicianUpdateOne {
	if tuo.removedWorkOrders == nil {
		tuo.removedWorkOrders = make(map[string]struct{})
	}
	for i := range ids {
		tuo.removedWorkOrders[ids[i]] = struct{}{}
	}
	return tuo
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (tuo *TechnicianUpdateOne) RemoveWorkOrders(w ...*WorkOrder) *TechnicianUpdateOne {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tuo.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (tuo *TechnicianUpdateOne) Save(ctx context.Context) (*Technician, error) {
	if tuo.update_time == nil {
		v := technician.UpdateDefaultUpdateTime()
		tuo.update_time = &v
	}
	if tuo.name != nil {
		if err := technician.NameValidator(*tuo.name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if tuo.email != nil {
		if err := technician.EmailValidator(*tuo.email); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	return tuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TechnicianUpdateOne) SaveX(ctx context.Context) *Technician {
	t, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return t
}

// Exec executes the query on the entity.
func (tuo *TechnicianUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TechnicianUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TechnicianUpdateOne) sqlSave(ctx context.Context) (t *Technician, err error) {
	var (
		builder  = sql.Dialect(tuo.driver.Dialect())
		selector = builder.Select(technician.Columns...).From(builder.Table(technician.Table))
	)
	technician.ID(tuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = tuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		t = &Technician{config: tuo.config}
		if err := t.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Technician: %v", err)
		}
		id = t.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Technician with id: %v", tuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Technician with the same id: %v", tuo.id)
	}

	tx, err := tuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(technician.Table)
	)
	updater = updater.Where(sql.InInts(technician.FieldID, ids...))
	if value := tuo.update_time; value != nil {
		updater.Set(technician.FieldUpdateTime, *value)
		t.UpdateTime = *value
	}
	if value := tuo.name; value != nil {
		updater.Set(technician.FieldName, *value)
		t.Name = *value
	}
	if value := tuo.email; value != nil {
		updater.Set(technician.FieldEmail, *value)
		t.Email = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(tuo.removedWorkOrders) > 0 {
		eids := make([]int, len(tuo.removedWorkOrders))
		for eid := range tuo.removedWorkOrders {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(technician.WorkOrdersTable).
			SetNull(technician.WorkOrdersColumn).
			Where(sql.InInts(technician.WorkOrdersColumn, ids...)).
			Where(sql.InInts(workorder.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(tuo.work_orders) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range tuo.work_orders {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(workorder.FieldID, eid)
			}
			query, args := builder.Update(technician.WorkOrdersTable).
				Set(technician.WorkOrdersColumn, id).
				Where(sql.And(p, sql.IsNull(technician.WorkOrdersColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(tuo.work_orders) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"Technician\"", keys(tuo.work_orders))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}
