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
	"github.com/facebookincubator/symphony/graph/ent/technician"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// TechnicianCreate is the builder for creating a Technician entity.
type TechnicianCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	name        *string
	email       *string
	work_orders map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (tc *TechnicianCreate) SetCreateTime(t time.Time) *TechnicianCreate {
	tc.create_time = &t
	return tc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (tc *TechnicianCreate) SetNillableCreateTime(t *time.Time) *TechnicianCreate {
	if t != nil {
		tc.SetCreateTime(*t)
	}
	return tc
}

// SetUpdateTime sets the update_time field.
func (tc *TechnicianCreate) SetUpdateTime(t time.Time) *TechnicianCreate {
	tc.update_time = &t
	return tc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (tc *TechnicianCreate) SetNillableUpdateTime(t *time.Time) *TechnicianCreate {
	if t != nil {
		tc.SetUpdateTime(*t)
	}
	return tc
}

// SetName sets the name field.
func (tc *TechnicianCreate) SetName(s string) *TechnicianCreate {
	tc.name = &s
	return tc
}

// SetEmail sets the email field.
func (tc *TechnicianCreate) SetEmail(s string) *TechnicianCreate {
	tc.email = &s
	return tc
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (tc *TechnicianCreate) AddWorkOrderIDs(ids ...string) *TechnicianCreate {
	if tc.work_orders == nil {
		tc.work_orders = make(map[string]struct{})
	}
	for i := range ids {
		tc.work_orders[ids[i]] = struct{}{}
	}
	return tc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (tc *TechnicianCreate) AddWorkOrders(w ...*WorkOrder) *TechnicianCreate {
	ids := make([]string, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tc.AddWorkOrderIDs(ids...)
}

// Save creates the Technician in the database.
func (tc *TechnicianCreate) Save(ctx context.Context) (*Technician, error) {
	if tc.create_time == nil {
		v := technician.DefaultCreateTime()
		tc.create_time = &v
	}
	if tc.update_time == nil {
		v := technician.DefaultUpdateTime()
		tc.update_time = &v
	}
	if tc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if err := technician.NameValidator(*tc.name); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
	}
	if tc.email == nil {
		return nil, errors.New("ent: missing required field \"email\"")
	}
	if err := technician.EmailValidator(*tc.email); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
	}
	return tc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TechnicianCreate) SaveX(ctx context.Context) *Technician {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tc *TechnicianCreate) sqlSave(ctx context.Context) (*Technician, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(tc.driver.Dialect())
		t       = &Technician{config: tc.config}
	)
	tx, err := tc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(technician.Table).Default()
	if value := tc.create_time; value != nil {
		insert.Set(technician.FieldCreateTime, *value)
		t.CreateTime = *value
	}
	if value := tc.update_time; value != nil {
		insert.Set(technician.FieldUpdateTime, *value)
		t.UpdateTime = *value
	}
	if value := tc.name; value != nil {
		insert.Set(technician.FieldName, *value)
		t.Name = *value
	}
	if value := tc.email; value != nil {
		insert.Set(technician.FieldEmail, *value)
		t.Email = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(technician.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	t.ID = strconv.FormatInt(id, 10)
	if len(tc.work_orders) > 0 {
		p := sql.P()
		for eid := range tc.work_orders {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
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
		if int(affected) < len(tc.work_orders) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"work_orders\" %v already connected to a different \"Technician\"", keys(tc.work_orders))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}
