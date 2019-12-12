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
	"github.com/facebookincubator/symphony/graph/ent/property"
)

// LinkCreate is the builder for creating a Link entity.
type LinkCreate struct {
	config
	create_time  *time.Time
	update_time  *time.Time
	future_state *string
	ports        map[string]struct{}
	work_order   map[string]struct{}
	properties   map[string]struct{}
	service      map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (lc *LinkCreate) SetCreateTime(t time.Time) *LinkCreate {
	lc.create_time = &t
	return lc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (lc *LinkCreate) SetNillableCreateTime(t *time.Time) *LinkCreate {
	if t != nil {
		lc.SetCreateTime(*t)
	}
	return lc
}

// SetUpdateTime sets the update_time field.
func (lc *LinkCreate) SetUpdateTime(t time.Time) *LinkCreate {
	lc.update_time = &t
	return lc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (lc *LinkCreate) SetNillableUpdateTime(t *time.Time) *LinkCreate {
	if t != nil {
		lc.SetUpdateTime(*t)
	}
	return lc
}

// SetFutureState sets the future_state field.
func (lc *LinkCreate) SetFutureState(s string) *LinkCreate {
	lc.future_state = &s
	return lc
}

// SetNillableFutureState sets the future_state field if the given value is not nil.
func (lc *LinkCreate) SetNillableFutureState(s *string) *LinkCreate {
	if s != nil {
		lc.SetFutureState(*s)
	}
	return lc
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (lc *LinkCreate) AddPortIDs(ids ...string) *LinkCreate {
	if lc.ports == nil {
		lc.ports = make(map[string]struct{})
	}
	for i := range ids {
		lc.ports[ids[i]] = struct{}{}
	}
	return lc
}

// AddPorts adds the ports edges to EquipmentPort.
func (lc *LinkCreate) AddPorts(e ...*EquipmentPort) *LinkCreate {
	ids := make([]string, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return lc.AddPortIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (lc *LinkCreate) SetWorkOrderID(id string) *LinkCreate {
	if lc.work_order == nil {
		lc.work_order = make(map[string]struct{})
	}
	lc.work_order[id] = struct{}{}
	return lc
}

// SetNillableWorkOrderID sets the work_order edge to WorkOrder by id if the given value is not nil.
func (lc *LinkCreate) SetNillableWorkOrderID(id *string) *LinkCreate {
	if id != nil {
		lc = lc.SetWorkOrderID(*id)
	}
	return lc
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (lc *LinkCreate) SetWorkOrder(w *WorkOrder) *LinkCreate {
	return lc.SetWorkOrderID(w.ID)
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (lc *LinkCreate) AddPropertyIDs(ids ...string) *LinkCreate {
	if lc.properties == nil {
		lc.properties = make(map[string]struct{})
	}
	for i := range ids {
		lc.properties[ids[i]] = struct{}{}
	}
	return lc
}

// AddProperties adds the properties edges to Property.
func (lc *LinkCreate) AddProperties(p ...*Property) *LinkCreate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lc.AddPropertyIDs(ids...)
}

// AddServiceIDs adds the service edge to Service by ids.
func (lc *LinkCreate) AddServiceIDs(ids ...string) *LinkCreate {
	if lc.service == nil {
		lc.service = make(map[string]struct{})
	}
	for i := range ids {
		lc.service[ids[i]] = struct{}{}
	}
	return lc
}

// AddService adds the service edges to Service.
func (lc *LinkCreate) AddService(s ...*Service) *LinkCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return lc.AddServiceIDs(ids...)
}

// Save creates the Link in the database.
func (lc *LinkCreate) Save(ctx context.Context) (*Link, error) {
	if lc.create_time == nil {
		v := link.DefaultCreateTime()
		lc.create_time = &v
	}
	if lc.update_time == nil {
		v := link.DefaultUpdateTime()
		lc.update_time = &v
	}
	if len(lc.work_order) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"work_order\"")
	}
	return lc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (lc *LinkCreate) SaveX(ctx context.Context) *Link {
	v, err := lc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (lc *LinkCreate) sqlSave(ctx context.Context) (*Link, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(lc.driver.Dialect())
		l       = &Link{config: lc.config}
	)
	tx, err := lc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(link.Table).Default()
	if value := lc.create_time; value != nil {
		insert.Set(link.FieldCreateTime, *value)
		l.CreateTime = *value
	}
	if value := lc.update_time; value != nil {
		insert.Set(link.FieldUpdateTime, *value)
		l.UpdateTime = *value
	}
	if value := lc.future_state; value != nil {
		insert.Set(link.FieldFutureState, *value)
		l.FutureState = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(link.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	l.ID = strconv.FormatInt(id, 10)
	if len(lc.ports) > 0 {
		p := sql.P()
		for eid := range lc.ports {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(lc.ports) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"ports\" %v already connected to a different \"Link\"", keys(lc.ports))})
		}
	}
	if len(lc.work_order) > 0 {
		for eid := range lc.work_order {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(link.WorkOrderTable).
				Set(link.WorkOrderColumn, eid).
				Where(sql.EQ(link.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(lc.properties) > 0 {
		p := sql.P()
		for eid := range lc.properties {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(lc.properties) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"properties\" %v already connected to a different \"Link\"", keys(lc.properties))})
		}
	}
	if len(lc.service) > 0 {
		for eid := range lc.service {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}

			query, args := builder.Insert(link.ServiceTable).
				Columns(link.ServicePrimaryKey[1], link.ServicePrimaryKey[0]).
				Values(id, eid).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return l, nil
}
