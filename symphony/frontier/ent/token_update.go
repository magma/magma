// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
	"github.com/facebookincubator/symphony/frontier/ent/token"
	"github.com/facebookincubator/symphony/frontier/ent/user"
)

// TokenUpdate is the builder for updating Token entities.
type TokenUpdate struct {
	config

	updated_at *time.Time

	user        map[int]struct{}
	clearedUser bool
	predicates  []predicate.Token
}

// Where adds a new predicate for the builder.
func (tu *TokenUpdate) Where(ps ...predicate.Token) *TokenUpdate {
	tu.predicates = append(tu.predicates, ps...)
	return tu
}

// SetUserID sets the user edge to User by id.
func (tu *TokenUpdate) SetUserID(id int) *TokenUpdate {
	if tu.user == nil {
		tu.user = make(map[int]struct{})
	}
	tu.user[id] = struct{}{}
	return tu
}

// SetUser sets the user edge to User.
func (tu *TokenUpdate) SetUser(u *User) *TokenUpdate {
	return tu.SetUserID(u.ID)
}

// ClearUser clears the user edge to User.
func (tu *TokenUpdate) ClearUser() *TokenUpdate {
	tu.clearedUser = true
	return tu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (tu *TokenUpdate) Save(ctx context.Context) (int, error) {
	if tu.updated_at == nil {
		v := token.UpdateDefaultUpdatedAt()
		tu.updated_at = &v
	}
	if len(tu.user) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"user\"")
	}
	if tu.clearedUser && tu.user == nil {
		return 0, errors.New("ent: clearing a unique edge \"user\"")
	}
	return tu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TokenUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TokenUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TokenUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TokenUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   token.Table,
			Columns: token.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: token.FieldID,
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
	if value := tu.updated_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: token.FieldUpdatedAt,
		})
	}
	if tu.clearedUser {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   token.UserTable,
			Columns: []string{token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.user; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   token.UserTable,
			Columns: []string{token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
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

// TokenUpdateOne is the builder for updating a single Token entity.
type TokenUpdateOne struct {
	config
	id int

	updated_at *time.Time

	user        map[int]struct{}
	clearedUser bool
}

// SetUserID sets the user edge to User by id.
func (tuo *TokenUpdateOne) SetUserID(id int) *TokenUpdateOne {
	if tuo.user == nil {
		tuo.user = make(map[int]struct{})
	}
	tuo.user[id] = struct{}{}
	return tuo
}

// SetUser sets the user edge to User.
func (tuo *TokenUpdateOne) SetUser(u *User) *TokenUpdateOne {
	return tuo.SetUserID(u.ID)
}

// ClearUser clears the user edge to User.
func (tuo *TokenUpdateOne) ClearUser() *TokenUpdateOne {
	tuo.clearedUser = true
	return tuo
}

// Save executes the query and returns the updated entity.
func (tuo *TokenUpdateOne) Save(ctx context.Context) (*Token, error) {
	if tuo.updated_at == nil {
		v := token.UpdateDefaultUpdatedAt()
		tuo.updated_at = &v
	}
	if len(tuo.user) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"user\"")
	}
	if tuo.clearedUser && tuo.user == nil {
		return nil, errors.New("ent: clearing a unique edge \"user\"")
	}
	return tuo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TokenUpdateOne) SaveX(ctx context.Context) *Token {
	t, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return t
}

// Exec executes the query on the entity.
func (tuo *TokenUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TokenUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TokenUpdateOne) sqlSave(ctx context.Context) (t *Token, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   token.Table,
			Columns: token.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  tuo.id,
				Type:   field.TypeInt,
				Column: token.FieldID,
			},
		},
	}
	if value := tuo.updated_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: token.FieldUpdatedAt,
		})
	}
	if tuo.clearedUser {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   token.UserTable,
			Columns: []string{token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.user; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   token.UserTable,
			Columns: []string{token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	t = &Token{config: tuo.config}
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
