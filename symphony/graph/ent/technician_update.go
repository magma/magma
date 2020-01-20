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
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
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
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   technician.Table,
			Columns: technician.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: technician.FieldID,
			},
		},
	}
	if ps := tu.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := tu.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: technician.FieldUpdateTime,
		})
	}
	if value := tu.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: technician.FieldName,
		})
	}
	if value := tu.email; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: technician.FieldEmail,
		})
	}
	if nodes := tu.removedWorkOrders; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := tu.work_orders; len(nodes) > 0 {
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
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
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
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   technician.Table,
			Columns: technician.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  tuo.id,
				Type:   field.TypeString,
				Column: technician.FieldID,
			},
		},
	}
	if value := tuo.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: technician.FieldUpdateTime,
		})
	}
	if value := tuo.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: technician.FieldName,
		})
	}
	if value := tuo.email; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: technician.FieldEmail,
		})
	}
	if nodes := tuo.removedWorkOrders; len(nodes) > 0 {
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
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := tuo.work_orders; len(nodes) > 0 {
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
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	t = &Technician{config: tuo.config}
	spec.Assign = t.assignValues
	spec.ScanValues = t.scanValues()
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return t, nil
}
