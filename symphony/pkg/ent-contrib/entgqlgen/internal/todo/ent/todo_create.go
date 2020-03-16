// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/todo"
)

// TodoCreate is the builder for creating a Todo entity.
type TodoCreate struct {
	config
	text     *string
	parent   map[int]struct{}
	children map[int]struct{}
}

// SetText sets the text field.
func (tc *TodoCreate) SetText(s string) *TodoCreate {
	tc.text = &s
	return tc
}

// SetParentID sets the parent edge to Todo by id.
func (tc *TodoCreate) SetParentID(id int) *TodoCreate {
	if tc.parent == nil {
		tc.parent = make(map[int]struct{})
	}
	tc.parent[id] = struct{}{}
	return tc
}

// SetNillableParentID sets the parent edge to Todo by id if the given value is not nil.
func (tc *TodoCreate) SetNillableParentID(id *int) *TodoCreate {
	if id != nil {
		tc = tc.SetParentID(*id)
	}
	return tc
}

// SetParent sets the parent edge to Todo.
func (tc *TodoCreate) SetParent(t *Todo) *TodoCreate {
	return tc.SetParentID(t.ID)
}

// AddChildIDs adds the children edge to Todo by ids.
func (tc *TodoCreate) AddChildIDs(ids ...int) *TodoCreate {
	if tc.children == nil {
		tc.children = make(map[int]struct{})
	}
	for i := range ids {
		tc.children[ids[i]] = struct{}{}
	}
	return tc
}

// AddChildren adds the children edges to Todo.
func (tc *TodoCreate) AddChildren(t ...*Todo) *TodoCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tc.AddChildIDs(ids...)
}

// Save creates the Todo in the database.
func (tc *TodoCreate) Save(ctx context.Context) (*Todo, error) {
	if tc.text == nil {
		return nil, errors.New("ent: missing required field \"text\"")
	}
	if err := todo.TextValidator(*tc.text); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"text\": %v", err)
	}
	if len(tc.parent) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"parent\"")
	}
	return tc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TodoCreate) SaveX(ctx context.Context) *Todo {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tc *TodoCreate) sqlSave(ctx context.Context) (*Todo, error) {
	var (
		t     = &Todo{config: tc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: todo.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: todo.FieldID,
			},
		}
	)
	if value := tc.text; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: todo.FieldText,
		})
		t.Text = *value
	}
	if nodes := tc.parent; len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.children; len(nodes) > 0 {
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
