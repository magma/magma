// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// CheckListCategoryUpdate is the builder for updating CheckListCategory entities.
type CheckListCategoryUpdate struct {
	config

	update_time           *time.Time
	title                 *string
	description           *string
	cleardescription      bool
	check_list_items      map[int]struct{}
	removedCheckListItems map[int]struct{}
	predicates            []predicate.CheckListCategory
}

// Where adds a new predicate for the builder.
func (clcu *CheckListCategoryUpdate) Where(ps ...predicate.CheckListCategory) *CheckListCategoryUpdate {
	clcu.predicates = append(clcu.predicates, ps...)
	return clcu
}

// SetTitle sets the title field.
func (clcu *CheckListCategoryUpdate) SetTitle(s string) *CheckListCategoryUpdate {
	clcu.title = &s
	return clcu
}

// SetDescription sets the description field.
func (clcu *CheckListCategoryUpdate) SetDescription(s string) *CheckListCategoryUpdate {
	clcu.description = &s
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
	clcu.description = nil
	clcu.cleardescription = true
	return clcu
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (clcu *CheckListCategoryUpdate) AddCheckListItemIDs(ids ...int) *CheckListCategoryUpdate {
	if clcu.check_list_items == nil {
		clcu.check_list_items = make(map[int]struct{})
	}
	for i := range ids {
		clcu.check_list_items[ids[i]] = struct{}{}
	}
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

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (clcu *CheckListCategoryUpdate) RemoveCheckListItemIDs(ids ...int) *CheckListCategoryUpdate {
	if clcu.removedCheckListItems == nil {
		clcu.removedCheckListItems = make(map[int]struct{})
	}
	for i := range ids {
		clcu.removedCheckListItems[ids[i]] = struct{}{}
	}
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

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (clcu *CheckListCategoryUpdate) Save(ctx context.Context) (int, error) {
	if clcu.update_time == nil {
		v := checklistcategory.UpdateDefaultUpdateTime()
		clcu.update_time = &v
	}
	return clcu.sqlSave(ctx)
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
	if value := clcu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: checklistcategory.FieldUpdateTime,
		})
	}
	if value := clcu.title; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistcategory.FieldTitle,
		})
	}
	if value := clcu.description; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistcategory.FieldDescription,
		})
	}
	if clcu.cleardescription {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistcategory.FieldDescription,
		})
	}
	if nodes := clcu.removedCheckListItems; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clcu.check_list_items; len(nodes) > 0 {
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
		for k, _ := range nodes {
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
	id int

	update_time           *time.Time
	title                 *string
	description           *string
	cleardescription      bool
	check_list_items      map[int]struct{}
	removedCheckListItems map[int]struct{}
}

// SetTitle sets the title field.
func (clcuo *CheckListCategoryUpdateOne) SetTitle(s string) *CheckListCategoryUpdateOne {
	clcuo.title = &s
	return clcuo
}

// SetDescription sets the description field.
func (clcuo *CheckListCategoryUpdateOne) SetDescription(s string) *CheckListCategoryUpdateOne {
	clcuo.description = &s
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
	clcuo.description = nil
	clcuo.cleardescription = true
	return clcuo
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (clcuo *CheckListCategoryUpdateOne) AddCheckListItemIDs(ids ...int) *CheckListCategoryUpdateOne {
	if clcuo.check_list_items == nil {
		clcuo.check_list_items = make(map[int]struct{})
	}
	for i := range ids {
		clcuo.check_list_items[ids[i]] = struct{}{}
	}
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

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (clcuo *CheckListCategoryUpdateOne) RemoveCheckListItemIDs(ids ...int) *CheckListCategoryUpdateOne {
	if clcuo.removedCheckListItems == nil {
		clcuo.removedCheckListItems = make(map[int]struct{})
	}
	for i := range ids {
		clcuo.removedCheckListItems[ids[i]] = struct{}{}
	}
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

// Save executes the query and returns the updated entity.
func (clcuo *CheckListCategoryUpdateOne) Save(ctx context.Context) (*CheckListCategory, error) {
	if clcuo.update_time == nil {
		v := checklistcategory.UpdateDefaultUpdateTime()
		clcuo.update_time = &v
	}
	return clcuo.sqlSave(ctx)
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
				Value:  clcuo.id,
				Type:   field.TypeInt,
				Column: checklistcategory.FieldID,
			},
		},
	}
	if value := clcuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: checklistcategory.FieldUpdateTime,
		})
	}
	if value := clcuo.title; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistcategory.FieldTitle,
		})
	}
	if value := clcuo.description; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: checklistcategory.FieldDescription,
		})
	}
	if clcuo.cleardescription {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: checklistcategory.FieldDescription,
		})
	}
	if nodes := clcuo.removedCheckListItems; len(nodes) > 0 {
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
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := clcuo.check_list_items; len(nodes) > 0 {
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
		for k, _ := range nodes {
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
