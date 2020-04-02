// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/technician"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// TechnicianUpdate is the builder for updating Technician entities.
type TechnicianUpdate struct {
	config
	hooks      []Hook
	mutation   *TechnicianMutation
	predicates []predicate.Technician
}

// Where adds a new predicate for the builder.
func (tu *TechnicianUpdate) Where(ps ...predicate.Technician) *TechnicianUpdate {
	tu.predicates = append(tu.predicates, ps...)
	return tu
}

// SetName sets the name field.
func (tu *TechnicianUpdate) SetName(s string) *TechnicianUpdate {
	tu.mutation.SetName(s)
	return tu
}

// SetEmail sets the email field.
func (tu *TechnicianUpdate) SetEmail(s string) *TechnicianUpdate {
	tu.mutation.SetEmail(s)
	return tu
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (tu *TechnicianUpdate) AddWorkOrderIDs(ids ...int) *TechnicianUpdate {
	tu.mutation.AddWorkOrderIDs(ids...)
	return tu
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (tu *TechnicianUpdate) AddWorkOrders(w ...*WorkOrder) *TechnicianUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tu.AddWorkOrderIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (tu *TechnicianUpdate) RemoveWorkOrderIDs(ids ...int) *TechnicianUpdate {
	tu.mutation.RemoveWorkOrderIDs(ids...)
	return tu
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (tu *TechnicianUpdate) RemoveWorkOrders(w ...*WorkOrder) *TechnicianUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tu.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (tu *TechnicianUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := tu.mutation.UpdateTime(); !ok {
		v := technician.UpdateDefaultUpdateTime()
		tu.mutation.SetUpdateTime(v)
	}
	if v, ok := tu.mutation.Name(); ok {
		if err := technician.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := tu.mutation.Email(); ok {
		if err := technician.EmailValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(tu.hooks) == 0 {
		affected, err = tu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TechnicianMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tu.mutation = mutation
			affected, err = tu.sqlSave(ctx)
			return affected, err
		})
		for i := len(tu.hooks) - 1; i >= 0; i-- {
			mut = tu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   technician.Table,
			Columns: technician.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: technician.FieldID,
			},
		},
	}
	if ps := tu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: technician.FieldUpdateTime,
		})
	}
	if value, ok := tu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: technician.FieldName,
		})
	}
	if value, ok := tu.mutation.Email(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: technician.FieldEmail,
		})
	}
	if nodes := tu.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{technician.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// TechnicianUpdateOne is the builder for updating a single Technician entity.
type TechnicianUpdateOne struct {
	config
	hooks    []Hook
	mutation *TechnicianMutation
}

// SetName sets the name field.
func (tuo *TechnicianUpdateOne) SetName(s string) *TechnicianUpdateOne {
	tuo.mutation.SetName(s)
	return tuo
}

// SetEmail sets the email field.
func (tuo *TechnicianUpdateOne) SetEmail(s string) *TechnicianUpdateOne {
	tuo.mutation.SetEmail(s)
	return tuo
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (tuo *TechnicianUpdateOne) AddWorkOrderIDs(ids ...int) *TechnicianUpdateOne {
	tuo.mutation.AddWorkOrderIDs(ids...)
	return tuo
}

// AddWorkOrders adds the work_orders edges to WorkOrder.
func (tuo *TechnicianUpdateOne) AddWorkOrders(w ...*WorkOrder) *TechnicianUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tuo.AddWorkOrderIDs(ids...)
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (tuo *TechnicianUpdateOne) RemoveWorkOrderIDs(ids ...int) *TechnicianUpdateOne {
	tuo.mutation.RemoveWorkOrderIDs(ids...)
	return tuo
}

// RemoveWorkOrders removes work_orders edges to WorkOrder.
func (tuo *TechnicianUpdateOne) RemoveWorkOrders(w ...*WorkOrder) *TechnicianUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return tuo.RemoveWorkOrderIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (tuo *TechnicianUpdateOne) Save(ctx context.Context) (*Technician, error) {
	if _, ok := tuo.mutation.UpdateTime(); !ok {
		v := technician.UpdateDefaultUpdateTime()
		tuo.mutation.SetUpdateTime(v)
	}
	if v, ok := tuo.mutation.Name(); ok {
		if err := technician.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := tuo.mutation.Email(); ok {
		if err := technician.EmailValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}

	var (
		err  error
		node *Technician
	)
	if len(tuo.hooks) == 0 {
		node, err = tuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TechnicianMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tuo.mutation = mutation
			node, err = tuo.sqlSave(ctx)
			return node, err
		})
		for i := len(tuo.hooks) - 1; i >= 0; i-- {
			mut = tuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   technician.Table,
			Columns: technician.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: technician.FieldID,
			},
		},
	}
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Technician.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := tuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: technician.FieldUpdateTime,
		})
	}
	if value, ok := tuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: technician.FieldName,
		})
	}
	if value, ok := tuo.mutation.Email(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: technician.FieldEmail,
		})
	}
	if nodes := tuo.mutation.RemovedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.WorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	t = &Technician{config: tuo.config}
	_spec.Assign = t.assignValues
	_spec.ScanValues = t.scanValues()
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{technician.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return t, nil
}
