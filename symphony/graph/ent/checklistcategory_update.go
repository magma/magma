// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// CheckListCategoryUpdate is the builder for updating CheckListCategory entities.
type CheckListCategoryUpdate struct {
	config
	hooks      []Hook
	mutation   *CheckListCategoryMutation
	predicates []predicate.CheckListCategory
}

// Where adds a new predicate for the builder.
func (clcu *CheckListCategoryUpdate) Where(ps ...predicate.CheckListCategory) *CheckListCategoryUpdate {
	clcu.predicates = append(clcu.predicates, ps...)
	return clcu
}

// SetTitle sets the title field.
func (clcu *CheckListCategoryUpdate) SetTitle(s string) *CheckListCategoryUpdate {
	clcu.mutation.SetTitle(s)
	return clcu
}

// SetDescription sets the description field.
func (clcu *CheckListCategoryUpdate) SetDescription(s string) *CheckListCategoryUpdate {
	clcu.mutation.SetDescription(s)
	return clcu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (clcu *CheckListCategoryUpdate) SetNillableDescription(s *string) *CheckListCategoryUpdate {
	if s != nil {
		clcu.SetDescription(*s)
	}
	return clcu
}

// ClearDescription clears the value of description.
func (clcu *CheckListCategoryUpdate) ClearDescription() *CheckListCategoryUpdate {
	clcu.mutation.ClearDescription()
	return clcu
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (clcu *CheckListCategoryUpdate) AddCheckListItemIDs(ids ...int) *CheckListCategoryUpdate {
	clcu.mutation.AddCheckListItemIDs(ids...)
	return clcu
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (clcu *CheckListCategoryUpdate) AddCheckListItems(c ...*CheckListItem) *CheckListCategoryUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcu.AddCheckListItemIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (clcu *CheckListCategoryUpdate) SetWorkOrderID(id int) *CheckListCategoryUpdate {
	clcu.mutation.SetWorkOrderID(id)
	return clcu
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (clcu *CheckListCategoryUpdate) SetWorkOrder(w *WorkOrder) *CheckListCategoryUpdate {
	return clcu.SetWorkOrderID(w.ID)
}

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (clcu *CheckListCategoryUpdate) RemoveCheckListItemIDs(ids ...int) *CheckListCategoryUpdate {
	clcu.mutation.RemoveCheckListItemIDs(ids...)
	return clcu
}

// RemoveCheckListItems removes check_list_items edges to CheckListItem.
func (clcu *CheckListCategoryUpdate) RemoveCheckListItems(c ...*CheckListItem) *CheckListCategoryUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcu.RemoveCheckListItemIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (clcu *CheckListCategoryUpdate) ClearWorkOrder() *CheckListCategoryUpdate {
	clcu.mutation.ClearWorkOrder()
	return clcu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (clcu *CheckListCategoryUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := clcu.mutation.UpdateTime(); !ok {
		v := checklistcategory.UpdateDefaultUpdateTime()
		clcu.mutation.SetUpdateTime(v)
	}

	if _, ok := clcu.mutation.WorkOrderID(); clcu.mutation.WorkOrderCleared() && !ok {
		return 0, errors.New("ent: clearing a unique edge \"work_order\"")
	}
	var (
		err      error
		affected int
	)
	if len(clcu.hooks) == 0 {
		affected, err = clcu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clcu.mutation = mutation
			affected, err = clcu.sqlSave(ctx)
			return affected, err
		})
		for i := len(clcu.hooks) - 1; i >= 0; i-- {
			mut = clcu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clcu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (clcu *CheckListCategoryUpdate) SaveX(ctx context.Context) int {
	affected, err := clcu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (clcu *CheckListCategoryUpdate) Exec(ctx context.Context) error {
	_, err := clcu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (clcu *CheckListCategoryUpdate) ExecX(ctx context.Context) {
	if err := clcu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (clcu *CheckListCategoryUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistcategory.Table,
			Columns: checklistcategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategory.FieldID,
			},
		},
	}
	if ps := clcu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := clcu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistcategory.FieldUpdateTime,
		})
	}
	if value, ok := clcu.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategory.FieldTitle,
		})
	}
	if value, ok := clcu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategory.FieldDescription,
		})
	}
	if clcu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistcategory.FieldDescription,
		})
	}
	if nodes := clcu.mutation.RemovedCheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategory.CheckListItemsTable,
			Columns: []string{checklistcategory.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitem.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clcu.mutation.CheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategory.CheckListItemsTable,
			Columns: []string{checklistcategory.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitem.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if clcu.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategory.WorkOrderTable,
			Columns: []string{checklistcategory.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clcu.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategory.WorkOrderTable,
			Columns: []string{checklistcategory.WorkOrderColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, clcu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{checklistcategory.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// CheckListCategoryUpdateOne is the builder for updating a single CheckListCategory entity.
type CheckListCategoryUpdateOne struct {
	config
	hooks    []Hook
	mutation *CheckListCategoryMutation
}

// SetTitle sets the title field.
func (clcuo *CheckListCategoryUpdateOne) SetTitle(s string) *CheckListCategoryUpdateOne {
	clcuo.mutation.SetTitle(s)
	return clcuo
}

// SetDescription sets the description field.
func (clcuo *CheckListCategoryUpdateOne) SetDescription(s string) *CheckListCategoryUpdateOne {
	clcuo.mutation.SetDescription(s)
	return clcuo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (clcuo *CheckListCategoryUpdateOne) SetNillableDescription(s *string) *CheckListCategoryUpdateOne {
	if s != nil {
		clcuo.SetDescription(*s)
	}
	return clcuo
}

// ClearDescription clears the value of description.
func (clcuo *CheckListCategoryUpdateOne) ClearDescription() *CheckListCategoryUpdateOne {
	clcuo.mutation.ClearDescription()
	return clcuo
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (clcuo *CheckListCategoryUpdateOne) AddCheckListItemIDs(ids ...int) *CheckListCategoryUpdateOne {
	clcuo.mutation.AddCheckListItemIDs(ids...)
	return clcuo
}

// AddCheckListItems adds the check_list_items edges to CheckListItem.
func (clcuo *CheckListCategoryUpdateOne) AddCheckListItems(c ...*CheckListItem) *CheckListCategoryUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcuo.AddCheckListItemIDs(ids...)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (clcuo *CheckListCategoryUpdateOne) SetWorkOrderID(id int) *CheckListCategoryUpdateOne {
	clcuo.mutation.SetWorkOrderID(id)
	return clcuo
}

// SetWorkOrder sets the work_order edge to WorkOrder.
func (clcuo *CheckListCategoryUpdateOne) SetWorkOrder(w *WorkOrder) *CheckListCategoryUpdateOne {
	return clcuo.SetWorkOrderID(w.ID)
}

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (clcuo *CheckListCategoryUpdateOne) RemoveCheckListItemIDs(ids ...int) *CheckListCategoryUpdateOne {
	clcuo.mutation.RemoveCheckListItemIDs(ids...)
	return clcuo
}

// RemoveCheckListItems removes check_list_items edges to CheckListItem.
func (clcuo *CheckListCategoryUpdateOne) RemoveCheckListItems(c ...*CheckListItem) *CheckListCategoryUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return clcuo.RemoveCheckListItemIDs(ids...)
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (clcuo *CheckListCategoryUpdateOne) ClearWorkOrder() *CheckListCategoryUpdateOne {
	clcuo.mutation.ClearWorkOrder()
	return clcuo
}

// Save executes the query and returns the updated entity.
func (clcuo *CheckListCategoryUpdateOne) Save(ctx context.Context) (*CheckListCategory, error) {
	if _, ok := clcuo.mutation.UpdateTime(); !ok {
		v := checklistcategory.UpdateDefaultUpdateTime()
		clcuo.mutation.SetUpdateTime(v)
	}

	if _, ok := clcuo.mutation.WorkOrderID(); clcuo.mutation.WorkOrderCleared() && !ok {
		return nil, errors.New("ent: clearing a unique edge \"work_order\"")
	}
	var (
		err  error
		node *CheckListCategory
	)
	if len(clcuo.hooks) == 0 {
		node, err = clcuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CheckListCategoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			clcuo.mutation = mutation
			node, err = clcuo.sqlSave(ctx)
			return node, err
		})
		for i := len(clcuo.hooks) - 1; i >= 0; i-- {
			mut = clcuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, clcuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (clcuo *CheckListCategoryUpdateOne) SaveX(ctx context.Context) *CheckListCategory {
	clc, err := clcuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return clc
}

// Exec executes the query on the entity.
func (clcuo *CheckListCategoryUpdateOne) Exec(ctx context.Context) error {
	_, err := clcuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (clcuo *CheckListCategoryUpdateOne) ExecX(ctx context.Context) {
	if err := clcuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (clcuo *CheckListCategoryUpdateOne) sqlSave(ctx context.Context) (clc *CheckListCategory, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   checklistcategory.Table,
			Columns: checklistcategory.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: checklistcategory.FieldID,
			},
		},
	}
	id, ok := clcuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing CheckListCategory.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := clcuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: checklistcategory.FieldUpdateTime,
		})
	}
	if value, ok := clcuo.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategory.FieldTitle,
		})
	}
	if value, ok := clcuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: checklistcategory.FieldDescription,
		})
	}
	if clcuo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistcategory.FieldDescription,
		})
	}
	if nodes := clcuo.mutation.RemovedCheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategory.CheckListItemsTable,
			Columns: []string{checklistcategory.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitem.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clcuo.mutation.CheckListItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   checklistcategory.CheckListItemsTable,
			Columns: []string{checklistcategory.CheckListItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: checklistitem.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if clcuo.mutation.WorkOrderCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategory.WorkOrderTable,
			Columns: []string{checklistcategory.WorkOrderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: workorder.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clcuo.mutation.WorkOrderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   checklistcategory.WorkOrderTable,
			Columns: []string{checklistcategory.WorkOrderColumn},
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
	clc = &CheckListCategory{config: clcuo.config}
	_spec.Assign = clc.assignValues
	_spec.ScanValues = clc.scanValues()
	if err = sqlgraph.UpdateNode(ctx, clcuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{checklistcategory.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return clc, nil
}
