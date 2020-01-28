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

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
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
		t     = &Technician{config: tc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: technician.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: technician.FieldID,
			},
		}
	)
	if value := tc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: technician.FieldCreateTime,
		})
		t.CreateTime = *value
	}
	if value := tc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: technician.FieldUpdateTime,
		})
		t.UpdateTime = *value
	}
	if value := tc.name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: technician.FieldName,
		})
		t.Name = *value
	}
	if value := tc.email; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: technician.FieldEmail,
		})
		t.Email = *value
	}
	if nodes := tc.work_orders; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   technician.WorkOrdersTable,
			Columns: []string{technician.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: workorder.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	t.ID = strconv.FormatInt(id, 10)
	return t, nil
}
