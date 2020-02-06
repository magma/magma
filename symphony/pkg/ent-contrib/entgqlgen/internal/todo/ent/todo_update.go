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
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/todo"
)

// TodoUpdate is the builder for updating Todo entities.
type TodoUpdate struct {
	config
	text            *string
	parent          map[int]struct{}
	children        map[int]struct{}
	clearedParent   bool
	removedChildren map[int]struct{}
	predicates      []predicate.Todo
}

// Where adds a new predicate for the builder.
func (tu *TodoUpdate) Where(ps ...predicate.Todo) *TodoUpdate {
	tu.predicates = append(tu.predicates, ps...)
	return tu
}

// SetText sets the text field.
func (tu *TodoUpdate) SetText(s string) *TodoUpdate {
	tu.text = &s
	return tu
}

// SetParentID sets the parent edge to Todo by id.
func (tu *TodoUpdate) SetParentID(id int) *TodoUpdate {
	if tu.parent == nil {
		tu.parent = make(map[int]struct{})
	}
	tu.parent[id] = struct{}{}
	return tu
}

// SetNillableParentID sets the parent edge to Todo by id if the given value is not nil.
func (tu *TodoUpdate) SetNillableParentID(id *int) *TodoUpdate {
	if id != nil {
		tu = tu.SetParentID(*id)
	}
	return tu
}

// SetParent sets the parent edge to Todo.
func (tu *TodoUpdate) SetParent(t *Todo) *TodoUpdate {
	return tu.SetParentID(t.ID)
}

// AddChildIDs adds the children edge to Todo by ids.
func (tu *TodoUpdate) AddChildIDs(ids ...int) *TodoUpdate {
	if tu.children == nil {
		tu.children = make(map[int]struct{})
	}
	for i := range ids {
		tu.children[ids[i]] = struct{}{}
	}
	return tu
}

// AddChildren adds the children edges to Todo.
func (tu *TodoUpdate) AddChildren(t ...*Todo) *TodoUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tu.AddChildIDs(ids...)
}

// ClearParent clears the parent edge to Todo.
func (tu *TodoUpdate) ClearParent() *TodoUpdate {
	tu.clearedParent = true
	return tu
}

// RemoveChildIDs removes the children edge to Todo by ids.
func (tu *TodoUpdate) RemoveChildIDs(ids ...int) *TodoUpdate {
	if tu.removedChildren == nil {
		tu.removedChildren = make(map[int]struct{})
	}
	for i := range ids {
		tu.removedChildren[ids[i]] = struct{}{}
	}
	return tu
}

// RemoveChildren removes children edges to Todo.
func (tu *TodoUpdate) RemoveChildren(t ...*Todo) *TodoUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tu.RemoveChildIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (tu *TodoUpdate) Save(ctx context.Context) (int, error) {
	if tu.text != nil {
		if err := todo.TextValidator(*tu.text); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"text\": %v", err)
		}
	}
	if len(tu.parent) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	return tu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TodoUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TodoUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TodoUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TodoUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   todo.Table,
			Columns: todo.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: todo.FieldID,
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
	if value := tu.text; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: todo.FieldText,
		})
	}
	if tu.clearedParent {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   todo.ParentTable,
			Columns: []string{todo.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: todo.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.parent; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   todo.ParentTable,
			Columns: []string{todo.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: todo.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := tu.removedChildren; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   todo.ChildrenTable,
			Columns: []string{todo.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: todo.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.children; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   todo.ChildrenTable,
			Columns: []string{todo.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: todo.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// TodoUpdateOne is the builder for updating a single Todo entity.
type TodoUpdateOne struct {
	config
	id              int
	text            *string
	parent          map[int]struct{}
	children        map[int]struct{}
	clearedParent   bool
	removedChildren map[int]struct{}
}

// SetText sets the text field.
func (tuo *TodoUpdateOne) SetText(s string) *TodoUpdateOne {
	tuo.text = &s
	return tuo
}

// SetParentID sets the parent edge to Todo by id.
func (tuo *TodoUpdateOne) SetParentID(id int) *TodoUpdateOne {
	if tuo.parent == nil {
		tuo.parent = make(map[int]struct{})
	}
	tuo.parent[id] = struct{}{}
	return tuo
}

// SetNillableParentID sets the parent edge to Todo by id if the given value is not nil.
func (tuo *TodoUpdateOne) SetNillableParentID(id *int) *TodoUpdateOne {
	if id != nil {
		tuo = tuo.SetParentID(*id)
	}
	return tuo
}

// SetParent sets the parent edge to Todo.
func (tuo *TodoUpdateOne) SetParent(t *Todo) *TodoUpdateOne {
	return tuo.SetParentID(t.ID)
}

// AddChildIDs adds the children edge to Todo by ids.
func (tuo *TodoUpdateOne) AddChildIDs(ids ...int) *TodoUpdateOne {
	if tuo.children == nil {
		tuo.children = make(map[int]struct{})
	}
	for i := range ids {
		tuo.children[ids[i]] = struct{}{}
	}
	return tuo
}

// AddChildren adds the children edges to Todo.
func (tuo *TodoUpdateOne) AddChildren(t ...*Todo) *TodoUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tuo.AddChildIDs(ids...)
}

// ClearParent clears the parent edge to Todo.
func (tuo *TodoUpdateOne) ClearParent() *TodoUpdateOne {
	tuo.clearedParent = true
	return tuo
}

// RemoveChildIDs removes the children edge to Todo by ids.
func (tuo *TodoUpdateOne) RemoveChildIDs(ids ...int) *TodoUpdateOne {
	if tuo.removedChildren == nil {
		tuo.removedChildren = make(map[int]struct{})
	}
	for i := range ids {
		tuo.removedChildren[ids[i]] = struct{}{}
	}
	return tuo
}

// RemoveChildren removes children edges to Todo.
func (tuo *TodoUpdateOne) RemoveChildren(t ...*Todo) *TodoUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tuo.RemoveChildIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (tuo *TodoUpdateOne) Save(ctx context.Context) (*Todo, error) {
	if tuo.text != nil {
		if err := todo.TextValidator(*tuo.text); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"text\": %v", err)
		}
	}
	if len(tuo.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	return tuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TodoUpdateOne) SaveX(ctx context.Context) *Todo {
	t, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return t
}

// Exec executes the query on the entity.
func (tuo *TodoUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TodoUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TodoUpdateOne) sqlSave(ctx context.Context) (t *Todo, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   todo.Table,
			Columns: todo.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  tuo.id,
				Type:   field.TypeInt,
				Column: todo.FieldID,
			},
		},
	}
	if value := tuo.text; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: todo.FieldText,
		})
	}
	if tuo.clearedParent {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   todo.ParentTable,
			Columns: []string{todo.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: todo.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.parent; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   todo.ParentTable,
			Columns: []string{todo.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: todo.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := tuo.removedChildren; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   todo.ChildrenTable,
			Columns: []string{todo.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: todo.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.children; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   todo.ChildrenTable,
			Columns: []string{todo.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: todo.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	t = &Todo{config: tuo.config}
	_spec.Assign = t.assignValues
	_spec.ScanValues = t.scanValues()
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return t, nil
}
