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
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// LinkUpdate is the builder for updating Link entities.
type LinkUpdate struct {
	config

	update_time       *time.Time
	future_state      *string
	clearfuture_state bool
	ports             map[string]struct{}
	work_order        map[string]struct{}
	properties        map[string]struct{}
	service           map[string]struct{}
	removedPorts      map[string]struct{}
	clearedWorkOrder  bool
	removedProperties map[string]struct{}
	removedService    map[string]struct{}
	predicates        []predicate.Link
}

// Where adds a new predicate for the builder.
func (lu *LinkUpdate) Where(ps ...predicate.Link) *LinkUpdate {
	lu.predicates = append(lu.predicates, ps...)
	return lu
}

// SetFutureState sets the future_state field.
func (lu *LinkUpdate) SetFutureState(s string) *LinkUpdate {
	lu.future_state = &s
	return lu
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (lu *LinkUpdate) SetNillableFutureState(s *string) *LinkUpdate {
	if s != nil {
		lu.SetFutureState(*s)
	}
	return lu
}

// ClearFutureState clears the value of future_state.
func (lu *LinkUpdate) ClearFutureState() *LinkUpdate {
	lu.future_state = nil
	lu.clearfuture_state = true
	return lu
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (lu *LinkUpdate) AddPortIDs(ids ...string) *LinkUpdate {
	if lu.ports == nil {
		lu.ports = make(map[string]struct{})
	}
	for i := range ids {
		lu.ports[ids[i]] = struct{}{}
	}
	return lu
}

// AddPorts adds the ports edges to EquipmentPort.
func (lu *LinkUpdate) AddPorts(e ...*EquipmentPort) *LinkUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lu.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (lu *LinkUpdate) SetWorkOrderID(id string) *LinkUpdate {
	if lu.work_order == nil {
		lu.work_order = make(map[string]struct{})
	}
	lu.work_order[id] = struct{}{}
	return lu
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (lu *LinkUpdate) SetNillableWorkOrderID(id *string) *LinkUpdate {
	if id != nil {
		lu = lu.SetWorkOrderID(*id)
	}
	return lu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (lu *LinkUpdate) SetWorkOrder(w *WorkOrder) *LinkUpdate {
	return lu.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lu *LinkUpdate) AddPropertyIDs(ids ...string) *LinkUpdate {
	if lu.properties == nil {
		lu.properties = make(map[string]struct{})
	}
	for i := range ids {
		lu.properties[ids[i]] = struct{}{}
	}
	return lu
}

// AddProperties adds the properties edges to Property.
func (lu *LinkUpdate) AddProperties(p ...*Property) *LinkUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.AddPropertyIDs(ids...)
}

// AddServiceIDs adds the service edge to Service by ids.
func (lu *LinkUpdate) AddServiceIDs(ids ...string) *LinkUpdate {
	if lu.service == nil {
		lu.service = make(map[string]struct{})
	}
	for i := range ids {
		lu.service[ids[i]] = struct{}{}
	}
	return lu
}

// AddService adds the service edges to Service.
func (lu *LinkUpdate) AddService(s ...*Service) *LinkUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.AddServiceIDs(ids...)
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (lu *LinkUpdate) RemovePortIDs(ids ...string) *LinkUpdate {
	if lu.removedPorts == nil {
		lu.removedPorts = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedPorts[ids[i]] = struct{}{}
	}
	return lu
}

// RemovePorts removes ports edges to EquipmentPort.
func (lu *LinkUpdate) RemovePorts(e ...*EquipmentPort) *LinkUpdate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lu.RemovePortIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (lu *LinkUpdate) ClearWorkOrder() *LinkUpdate {
	lu.clearedWorkOrder = true
	return lu
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (lu *LinkUpdate) RemovePropertyIDs(ids ...string) *LinkUpdate {
	if lu.removedProperties == nil {
		lu.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedProperties[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveProperties removes properties edges to Property.
func (lu *LinkUpdate) RemoveProperties(p ...*Property) *LinkUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lu.RemovePropertyIDs(ids...)
}

// RemoveServiceIDs removes the service edge to Service by ids.
func (lu *LinkUpdate) RemoveServiceIDs(ids ...string) *LinkUpdate {
	if lu.removedService == nil {
		lu.removedService = make(map[string]struct{})
	}
	for i := range ids {
		lu.removedService[ids[i]] = struct{}{}
	}
	return lu
}

// RemoveService removes service edges to Service.
func (lu *LinkUpdate) RemoveService(s ...*Service) *LinkUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lu.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (lu *LinkUpdate) Save(ctx context.Context) (int, error) {
	if lu.update_time == nil {
		v := link.UpdateDefaultUpdateTime()
		lu.update_time = &v
	}
	if len(lu.work_order) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return lu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (lu *LinkUpdate) SaveX(ctx context.Context) int {
	affected, err := lu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lu *LinkUpdate) Exec(ctx context.Context) error {
	_, err := lu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lu *LinkUpdate) ExecX(ctx context.Context) {
	if err := lu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (lu *LinkUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(lu.driver.Dialect())
		selector = builder.Select(link.FieldID).From(builder.Table(link.Table))
	)
	for _, p := range lu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = lu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := lu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(link.Table)
	)
	updater = updater.Where(sql.InInts(link.FieldID, ids...))
	if value := lu.update_time; value != nil {
		updater.Set(link.FieldUpdateTime, *value)
	}
	if value := lu.future_state; value != nil {
		updater.Set(link.FieldFutureState, *value)
	}
	if lu.clearfuture_state {
		updater.SetNull(link.FieldFutureState)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.removedPorts) > 0 {
		eids := make([]int, len(lu.removedPorts))
		for eid := range lu.removedPorts {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(link.PortsTable).
			SetNull(link.PortsColumn).
			Where(sql.InInts(link.PortsColumn, ids...)).
			Where(sql.InInts(equipmentport.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.ports) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.ports {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentport.FieldID, eid)
			}
			query, args := builder.Update(link.PortsTable).
				Set(link.PortsColumn, id).
				Where(sql.And(p, sql.IsNull(link.PortsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.ports) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"Link\"", keys(lu.ports))})
			}
		}
	}
	if lu.clearedWorkOrder {
		query, args := builder.Update(link.WorkOrderTable).
			SetNull(link.WorkOrderColumn).
			Where(sql.InInts(workorder.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.work_order) > 0 {
		for eid := range lu.work_order {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(link.WorkOrderTable).
				Set(link.WorkOrderColumn, eid).
				Where(sql.InInts(link.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if len(lu.removedProperties) > 0 {
		eids := make([]int, len(lu.removedProperties))
		for eid := range lu.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(link.PropertiesTable).
			SetNull(link.PropertiesColumn).
			Where(sql.InInts(link.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range lu.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(link.PropertiesTable).
				Set(link.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(link.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(lu.properties) {
				return 0, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Link\"", keys(lu.properties))})
			}
		}
	}
	if len(lu.removedService) > 0 {
		eids := make([]int, len(lu.removedService))
		for eid := range lu.removedService {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(link.ServiceTable).
			Where(sql.InInts(link.ServicePrimaryKey[1], ids...)).
			Where(sql.InInts(link.ServicePrimaryKey[0], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(lu.service) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range lu.service {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(link.ServiceTable).
			Columns(link.ServicePrimaryKey[1], link.ServicePrimaryKey[0])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// LinkUpdateOne is the builder for updating a single Link entity.
type LinkUpdateOne struct {
	config
	id string

	update_time       *time.Time
	future_state      *string
	clearfuture_state bool
	ports             map[string]struct{}
	work_order        map[string]struct{}
	properties        map[string]struct{}
	service           map[string]struct{}
	removedPorts      map[string]struct{}
	clearedWorkOrder  bool
	removedProperties map[string]struct{}
	removedService    map[string]struct{}
}

// SetFutureState sets the future_state field.
func (luo *LinkUpdateOne) SetFutureState(s string) *LinkUpdateOne {
	luo.future_state = &s
	return luo
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (luo *LinkUpdateOne) SetNillableFutureState(s *string) *LinkUpdateOne {
	if s != nil {
		luo.SetFutureState(*s)
	}
	return luo
}

// ClearFutureState clears the value of future_state.
func (luo *LinkUpdateOne) ClearFutureState() *LinkUpdateOne {
	luo.future_state = nil
	luo.clearfuture_state = true
	return luo
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (luo *LinkUpdateOne) AddPortIDs(ids ...string) *LinkUpdateOne {
	if luo.ports == nil {
		luo.ports = make(map[string]struct{})
	}
	for i := range ids {
		luo.ports[ids[i]] = struct{}{}
	}
	return luo
}

// AddPorts adds the ports edges to EquipmentPort.
func (luo *LinkUpdateOne) AddPorts(e ...*EquipmentPort) *LinkUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return luo.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (luo *LinkUpdateOne) SetWorkOrderID(id string) *LinkUpdateOne {
	if luo.work_order == nil {
		luo.work_order = make(map[string]struct{})
	}
	luo.work_order[id] = struct{}{}
	return luo
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (luo *LinkUpdateOne) SetNillableWorkOrderID(id *string) *LinkUpdateOne {
	if id != nil {
		luo = luo.SetWorkOrderID(*id)
	}
	return luo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (luo *LinkUpdateOne) SetWorkOrder(w *WorkOrder) *LinkUpdateOne {
	return luo.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (luo *LinkUpdateOne) AddPropertyIDs(ids ...string) *LinkUpdateOne {
	if luo.properties == nil {
		luo.properties = make(map[string]struct{})
	}
	for i := range ids {
		luo.properties[ids[i]] = struct{}{}
	}
	return luo
}

// AddProperties adds the properties edges to Property.
func (luo *LinkUpdateOne) AddProperties(p ...*Property) *LinkUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.AddPropertyIDs(ids...)
}

// AddServiceIDs adds the service edge to Service by ids.
func (luo *LinkUpdateOne) AddServiceIDs(ids ...string) *LinkUpdateOne {
	if luo.service == nil {
		luo.service = make(map[string]struct{})
	}
	for i := range ids {
		luo.service[ids[i]] = struct{}{}
	}
	return luo
}

// AddService adds the service edges to Service.
func (luo *LinkUpdateOne) AddService(s ...*Service) *LinkUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.AddServiceIDs(ids...)
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (luo *LinkUpdateOne) RemovePortIDs(ids ...string) *LinkUpdateOne {
	if luo.removedPorts == nil {
		luo.removedPorts = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedPorts[ids[i]] = struct{}{}
	}
	return luo
}

// RemovePorts removes ports edges to EquipmentPort.
func (luo *LinkUpdateOne) RemovePorts(e ...*EquipmentPort) *LinkUpdateOne {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return luo.RemovePortIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (luo *LinkUpdateOne) ClearWorkOrder() *LinkUpdateOne {
	luo.clearedWorkOrder = true
	return luo
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (luo *LinkUpdateOne) RemovePropertyIDs(ids ...string) *LinkUpdateOne {
	if luo.removedProperties == nil {
		luo.removedProperties = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedProperties[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveProperties removes properties edges to Property.
func (luo *LinkUpdateOne) RemoveProperties(p ...*Property) *LinkUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return luo.RemovePropertyIDs(ids...)
}

// RemoveServiceIDs removes the service edge to Service by ids.
func (luo *LinkUpdateOne) RemoveServiceIDs(ids ...string) *LinkUpdateOne {
	if luo.removedService == nil {
		luo.removedService = make(map[string]struct{})
	}
	for i := range ids {
		luo.removedService[ids[i]] = struct{}{}
	}
	return luo
}

// RemoveService removes service edges to Service.
func (luo *LinkUpdateOne) RemoveService(s ...*Service) *LinkUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return luo.RemoveServiceIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (luo *LinkUpdateOne) Save(ctx context.Context) (*Link, error) {
	if luo.update_time == nil {
		v := link.UpdateDefaultUpdateTime()
		luo.update_time = &v
	}
	if len(luo.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return luo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (luo *LinkUpdateOne) SaveX(ctx context.Context) *Link {
	l, err := luo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return l
}

// Exec executes the query on the entity.
func (luo *LinkUpdateOne) Exec(ctx context.Context) error {
	_, err := luo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (luo *LinkUpdateOne) ExecX(ctx context.Context) {
	if err := luo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (luo *LinkUpdateOne) sqlSave(ctx context.Context) (l *Link, err error) {
	var (
		builder  = sql.Dialect(luo.driver.Dialect())
		selector = builder.Select(link.Columns...).From(builder.Table(link.Table))
	)
	link.ID(luo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = luo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		l = &Link{config: luo.config}
		if err := l.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Link: %v", err)
		}
		id = l.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Link with id: %v", luo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Link with the same id: %v", luo.id)
	}

	tx, err := luo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(link.Table)
	)
	updater = updater.Where(sql.InInts(link.FieldID, ids...))
	if value := luo.update_time; value != nil {
		updater.Set(link.FieldUpdateTime, *value)
		l.UpdateTime = *value
	}
	if value := luo.future_state; value != nil {
		updater.Set(link.FieldFutureState, *value)
		l.FutureState = *value
	}
	if luo.clearfuture_state {
		var value string
		l.FutureState = value
		updater.SetNull(link.FieldFutureState)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.removedPorts) > 0 {
		eids := make([]int, len(luo.removedPorts))
		for eid := range luo.removedPorts {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(link.PortsTable).
			SetNull(link.PortsColumn).
			Where(sql.InInts(link.PortsColumn, ids...)).
			Where(sql.InInts(equipmentport.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.ports) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.ports {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(equipmentport.FieldID, eid)
			}
			query, args := builder.Update(link.PortsTable).
				Set(link.PortsColumn, id).
				Where(sql.And(p, sql.IsNull(link.PortsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(luo.ports) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"Link\"", keys(luo.ports))})
			}
		}
	}
	if luo.clearedWorkOrder {
		query, args := builder.Update(link.WorkOrderTable).
			SetNull(link.WorkOrderColumn).
			Where(sql.InInts(workorder.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.work_order) > 0 {
		for eid := range luo.work_order {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(link.WorkOrderTable).
				Set(link.WorkOrderColumn, eid).
				Where(sql.InInts(link.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(luo.removedProperties) > 0 {
		eids := make([]int, len(luo.removedProperties))
		for eid := range luo.removedProperties {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(link.PropertiesTable).
			SetNull(link.PropertiesColumn).
			Where(sql.InInts(link.PropertiesColumn, ids...)).
			Where(sql.InInts(property.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.properties) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range luo.properties {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(property.FieldID, eid)
			}
			query, args := builder.Update(link.PropertiesTable).
				Set(link.PropertiesColumn, id).
				Where(sql.And(p, sql.IsNull(link.PropertiesColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return nil, rollback(tx, err)
			}
			if int(affected) < len(luo.properties) {
				return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Link\"", keys(luo.properties))})
			}
		}
	}
	if len(luo.removedService) > 0 {
		eids := make([]int, len(luo.removedService))
		for eid := range luo.removedService {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Delete(link.ServiceTable).
			Where(sql.InInts(link.ServicePrimaryKey[1], ids...)).
			Where(sql.InInts(link.ServicePrimaryKey[0], eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(luo.service) > 0 {
		values := make([][]int, 0, len(ids))
		for _, id := range ids {
			for eid := range luo.service {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				values = append(values, []int{id, eid})
			}
		}
		builder := builder.Insert(link.ServiceTable).
			Columns(link.ServicePrimaryKey[1], link.ServicePrimaryKey[0])
		for _, v := range values {
			builder.Values(v[0], v[1])
		}
		query, args := builder.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return l, nil
}
