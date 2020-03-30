// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/technician"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// TechnicianCreate is the builder for creating a Technician entity.
type TechnicianCreate struct {
	config
	mutation *TechnicianMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (tc *TechnicianCreate) SetCreateTime(t time.Time) *TechnicianCreate {
	tc.mutation.SetCreateTime(t)
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
	tc.mutation.SetUpdateTime(t)
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
	tc.mutation.SetName(s)
	return tc
}

// SetEmail sets the email field.
func (tc *TechnicianCreate) SetEmail(s string) *TechnicianCreate {
	tc.mutation.SetEmail(s)
	return tc
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (tc *TechnicianCreate) AddWorkOrderIDs(ids ...int) *TechnicianCreate {
	tc.mutation.AddWorkOrderIDs(ids...)
	return tc
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (tc *TechnicianCreate) AddWorkOrders(w ...*WorkOrder) *TechnicianCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tc.AddWorkOrderIDs(ids...)
}

// Save creates the Technician in the database.
func (tc *TechnicianCreate) Save(ctx context.Context) (*Technician, error) {
	if _, ok := tc.mutation.CreateTime(); !ok {
		v := technician.DefaultCreateTime()
		tc.mutation.SetCreateTime(v)
	}
	if _, ok := tc.mutation.UpdateTime(); !ok {
		v := technician.DefaultUpdateTime()
		tc.mutation.SetUpdateTime(v)
	}
	if _, ok := tc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := tc.mutation.Name(); ok {
		if err := technician.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if _, ok := tc.mutation.Email(); !ok {
		return nil, errors.New("ent: missing required field \"email\"")
	}
	if v, ok := tc.mutation.Email(); ok {
		if err := technician.EmailValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	var (
		err  error
		node *Technician
	)
	if len(tc.hooks) == 0 {
		node, err = tc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TechnicianMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tc.mutation = mutation
			node, err = tc.sqlSave(ctx)
			return node, err
		})
		for i := len(tc.hooks) - 1; i >= 0; i-- {
			mut = tc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
				Type:   field.TypeInt,
				Column: technician.FieldID,
			},
		}
	)
	if value, ok := tc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: technician.FieldCreateTime,
		})
		t.CreateTime = value
	}
	if value, ok := tc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: technician.FieldUpdateTime,
		})
		t.UpdateTime = value
	}
	if value, ok := tc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: technician.FieldName,
		})
		t.Name = value
	}
	if value, ok := tc.mutation.Email(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: technician.FieldEmail,
		})
		t.Email = value
	}
	if nodes := tc.mutation.WorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   technician.WorkOrdersTable,
			Columns: []string{technician.WorkOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		for _, k := range nodes {
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
	t.ID = int(id)
	return t, nil
}
