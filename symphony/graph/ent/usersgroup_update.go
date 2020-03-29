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
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
)

// UsersGroupUpdate is the builder for updating UsersGroup entities.
type UsersGroupUpdate struct {
	config
	hooks      []Hook
	mutation   *UsersGroupMutation
	predicates []predicate.UsersGroup
}

// Where adds a new predicate for the builder.
func (ugu *UsersGroupUpdate) Where(ps ...predicate.UsersGroup) *UsersGroupUpdate {
	ugu.predicates = append(ugu.predicates, ps...)
	return ugu
}

// SetName sets the name field.
func (ugu *UsersGroupUpdate) SetName(s string) *UsersGroupUpdate {
	ugu.mutation.SetName(s)
	return ugu
}

// SetDescription sets the description field.
func (ugu *UsersGroupUpdate) SetDescription(s string) *UsersGroupUpdate {
	ugu.mutation.SetDescription(s)
	return ugu
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ugu *UsersGroupUpdate) SetNillableDescription(s *string) *UsersGroupUpdate {
	if s != nil {
		ugu.SetDescription(*s)
	}
	return ugu
}

// ClearDescription clears the value of description.
func (ugu *UsersGroupUpdate) ClearDescription() *UsersGroupUpdate {
	ugu.mutation.ClearDescription()
	return ugu
}

// SetStatus sets the status field.
func (ugu *UsersGroupUpdate) SetStatus(u usersgroup.Status) *UsersGroupUpdate {
	ugu.mutation.SetStatus(u)
	return ugu
}

// SetNillableStatus sets the status field if the given value is not nil.
func (ugu *UsersGroupUpdate) SetNillableStatus(u *usersgroup.Status) *UsersGroupUpdate {
	if u != nil {
		ugu.SetStatus(*u)
	}
	return ugu
}

// AddMemberIDs adds the members edge to User by ids.
func (ugu *UsersGroupUpdate) AddMemberIDs(ids ...int) *UsersGroupUpdate {
	ugu.mutation.AddMemberIDs(ids...)
	return ugu
}

// AddMembers adds the members edges to User.
func (ugu *UsersGroupUpdate) AddMembers(u ...*User) *UsersGroupUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ugu.AddMemberIDs(ids...)
}

// RemoveMemberIDs removes the members edge to User by ids.
func (ugu *UsersGroupUpdate) RemoveMemberIDs(ids ...int) *UsersGroupUpdate {
	ugu.mutation.RemoveMemberIDs(ids...)
	return ugu
}

// RemoveMembers removes members edges to User.
func (ugu *UsersGroupUpdate) RemoveMembers(u ...*User) *UsersGroupUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ugu.RemoveMemberIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ugu *UsersGroupUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := ugu.mutation.UpdateTime(); !ok {
		v := usersgroup.UpdateDefaultUpdateTime()
		ugu.mutation.SetUpdateTime(v)
	}
	if v, ok := ugu.mutation.Name(); ok {
		if err := usersgroup.NameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := ugu.mutation.Status(); ok {
		if err := usersgroup.StatusValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(ugu.hooks) == 0 {
		affected, err = ugu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UsersGroupMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ugu.mutation = mutation
			affected, err = ugu.sqlSave(ctx)
			return affected, err
		})
		for i := len(ugu.hooks); i > 0; i-- {
			mut = ugu.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, ugu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ugu *UsersGroupUpdate) SaveX(ctx context.Context) int {
	affected, err := ugu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ugu *UsersGroupUpdate) Exec(ctx context.Context) error {
	_, err := ugu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ugu *UsersGroupUpdate) ExecX(ctx context.Context) {
	if err := ugu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ugu *UsersGroupUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   usersgroup.Table,
			Columns: usersgroup.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: usersgroup.FieldID,
			},
		},
	}
	if ps := ugu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ugu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: usersgroup.FieldUpdateTime,
		})
	}
	if value, ok := ugu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usersgroup.FieldName,
		})
	}
	if value, ok := ugu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usersgroup.FieldDescription,
		})
	}
	if ugu.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: usersgroup.FieldDescription,
		})
	}
	if value, ok := ugu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: usersgroup.FieldStatus,
		})
	}
	if nodes := ugu.mutation.RemovedMembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   usersgroup.MembersTable,
			Columns: usersgroup.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ugu.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   usersgroup.MembersTable,
			Columns: usersgroup.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ugu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{usersgroup.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// UsersGroupUpdateOne is the builder for updating a single UsersGroup entity.
type UsersGroupUpdateOne struct {
	config
	hooks    []Hook
	mutation *UsersGroupMutation
}

// SetName sets the name field.
func (uguo *UsersGroupUpdateOne) SetName(s string) *UsersGroupUpdateOne {
	uguo.mutation.SetName(s)
	return uguo
}

// SetDescription sets the description field.
func (uguo *UsersGroupUpdateOne) SetDescription(s string) *UsersGroupUpdateOne {
	uguo.mutation.SetDescription(s)
	return uguo
}

// SetNillableDescription sets the description field if the given value is not nil.
func (uguo *UsersGroupUpdateOne) SetNillableDescription(s *string) *UsersGroupUpdateOne {
	if s != nil {
		uguo.SetDescription(*s)
	}
	return uguo
}

// ClearDescription clears the value of description.
func (uguo *UsersGroupUpdateOne) ClearDescription() *UsersGroupUpdateOne {
	uguo.mutation.ClearDescription()
	return uguo
}

// SetStatus sets the status field.
func (uguo *UsersGroupUpdateOne) SetStatus(u usersgroup.Status) *UsersGroupUpdateOne {
	uguo.mutation.SetStatus(u)
	return uguo
}

// SetNillableStatus sets the status field if the given value is not nil.
func (uguo *UsersGroupUpdateOne) SetNillableStatus(u *usersgroup.Status) *UsersGroupUpdateOne {
	if u != nil {
		uguo.SetStatus(*u)
	}
	return uguo
}

// AddMemberIDs adds the members edge to User by ids.
func (uguo *UsersGroupUpdateOne) AddMemberIDs(ids ...int) *UsersGroupUpdateOne {
	uguo.mutation.AddMemberIDs(ids...)
	return uguo
}

// AddMembers adds the members edges to User.
func (uguo *UsersGroupUpdateOne) AddMembers(u ...*User) *UsersGroupUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uguo.AddMemberIDs(ids...)
}

// RemoveMemberIDs removes the members edge to User by ids.
func (uguo *UsersGroupUpdateOne) RemoveMemberIDs(ids ...int) *UsersGroupUpdateOne {
	uguo.mutation.RemoveMemberIDs(ids...)
	return uguo
}

// RemoveMembers removes members edges to User.
func (uguo *UsersGroupUpdateOne) RemoveMembers(u ...*User) *UsersGroupUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uguo.RemoveMemberIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (uguo *UsersGroupUpdateOne) Save(ctx context.Context) (*UsersGroup, error) {
	if _, ok := uguo.mutation.UpdateTime(); !ok {
		v := usersgroup.UpdateDefaultUpdateTime()
		uguo.mutation.SetUpdateTime(v)
	}
	if v, ok := uguo.mutation.Name(); ok {
		if err := usersgroup.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if v, ok := uguo.mutation.Status(); ok {
		if err := usersgroup.StatusValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
		}
	}

	var (
		err  error
		node *UsersGroup
	)
	if len(uguo.hooks) == 0 {
		node, err = uguo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UsersGroupMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			uguo.mutation = mutation
			node, err = uguo.sqlSave(ctx)
			return node, err
		})
		for i := len(uguo.hooks); i > 0; i-- {
			mut = uguo.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, uguo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (uguo *UsersGroupUpdateOne) SaveX(ctx context.Context) *UsersGroup {
	ug, err := uguo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ug
}

// Exec executes the query on the entity.
func (uguo *UsersGroupUpdateOne) Exec(ctx context.Context) error {
	_, err := uguo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uguo *UsersGroupUpdateOne) ExecX(ctx context.Context) {
	if err := uguo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (uguo *UsersGroupUpdateOne) sqlSave(ctx context.Context) (ug *UsersGroup, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   usersgroup.Table,
			Columns: usersgroup.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: usersgroup.FieldID,
			},
		},
	}
	id, ok := uguo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing UsersGroup.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := uguo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: usersgroup.FieldUpdateTime,
		})
	}
	if value, ok := uguo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usersgroup.FieldName,
		})
	}
	if value, ok := uguo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usersgroup.FieldDescription,
		})
	}
	if uguo.mutation.DescriptionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: usersgroup.FieldDescription,
		})
	}
	if value, ok := uguo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: usersgroup.FieldStatus,
		})
	}
	if nodes := uguo.mutation.RemovedMembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   usersgroup.MembersTable,
			Columns: usersgroup.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uguo.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   usersgroup.MembersTable,
			Columns: usersgroup.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	ug = &UsersGroup{config: uguo.config}
	_spec.Assign = ug.assignValues
	_spec.ScanValues = ug.scanValues()
	if err = sqlgraph.UpdateNode(ctx, uguo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{usersgroup.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ug, nil
}
