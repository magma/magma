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
	"github.com/facebookincubator/symphony/graph/ent/permissionspolicy"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
)

// UsersGroupCreate is the builder for creating a UsersGroup entity.
type UsersGroupCreate struct {
	config
	mutation *UsersGroupMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (ugc *UsersGroupCreate) SetCreateTime(t time.Time) *UsersGroupCreate {
	ugc.mutation.SetCreateTime(t)
	return ugc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (ugc *UsersGroupCreate) SetNillableCreateTime(t *time.Time) *UsersGroupCreate {
	if t != nil {
		ugc.SetCreateTime(*t)
	}
	return ugc
}

// SetUpdateTime sets the update_time field.
func (ugc *UsersGroupCreate) SetUpdateTime(t time.Time) *UsersGroupCreate {
	ugc.mutation.SetUpdateTime(t)
	return ugc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (ugc *UsersGroupCreate) SetNillableUpdateTime(t *time.Time) *UsersGroupCreate {
	if t != nil {
		ugc.SetUpdateTime(*t)
	}
	return ugc
}

// SetName sets the name field.
func (ugc *UsersGroupCreate) SetName(s string) *UsersGroupCreate {
	ugc.mutation.SetName(s)
	return ugc
}

// SetDescription sets the description field.
func (ugc *UsersGroupCreate) SetDescription(s string) *UsersGroupCreate {
	ugc.mutation.SetDescription(s)
	return ugc
}

// SetNillableDescription sets the description field if the given value is not nil.
func (ugc *UsersGroupCreate) SetNillableDescription(s *string) *UsersGroupCreate {
	if s != nil {
		ugc.SetDescription(*s)
	}
	return ugc
}

// SetStatus sets the status field.
func (ugc *UsersGroupCreate) SetStatus(u usersgroup.Status) *UsersGroupCreate {
	ugc.mutation.SetStatus(u)
	return ugc
}

// SetNillableStatus sets the status field if the given value is not nil.
func (ugc *UsersGroupCreate) SetNillableStatus(u *usersgroup.Status) *UsersGroupCreate {
	if u != nil {
		ugc.SetStatus(*u)
	}
	return ugc
}

// AddMemberIDs adds the members edge to User by ids.
func (ugc *UsersGroupCreate) AddMemberIDs(ids ...int) *UsersGroupCreate {
	ugc.mutation.AddMemberIDs(ids...)
	return ugc
}

// AddMembers adds the members edges to User.
func (ugc *UsersGroupCreate) AddMembers(u ...*User) *UsersGroupCreate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return ugc.AddMemberIDs(ids...)
}

// AddPolicyIDs adds the policies edge to PermissionsPolicy by ids.
func (ugc *UsersGroupCreate) AddPolicyIDs(ids ...int) *UsersGroupCreate {
	ugc.mutation.AddPolicyIDs(ids...)
	return ugc
}

// AddPolicies adds the policies edges to PermissionsPolicy.
func (ugc *UsersGroupCreate) AddPolicies(p ...*PermissionsPolicy) *UsersGroupCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ugc.AddPolicyIDs(ids...)
}

// Save creates the UsersGroup in the database.
func (ugc *UsersGroupCreate) Save(ctx context.Context) (*UsersGroup, error) {
	if _, ok := ugc.mutation.CreateTime(); !ok {
		v := usersgroup.DefaultCreateTime()
		ugc.mutation.SetCreateTime(v)
	}
	if _, ok := ugc.mutation.UpdateTime(); !ok {
		v := usersgroup.DefaultUpdateTime()
		ugc.mutation.SetUpdateTime(v)
	}
	if _, ok := ugc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := ugc.mutation.Name(); ok {
		if err := usersgroup.NameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"name\": %v", err)
		}
	}
	if _, ok := ugc.mutation.Status(); !ok {
		v := usersgroup.DefaultStatus
		ugc.mutation.SetStatus(v)
	}
	if v, ok := ugc.mutation.Status(); ok {
		if err := usersgroup.StatusValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
		}
	}
	var (
		err  error
		node *UsersGroup
	)
	if len(ugc.hooks) == 0 {
		node, err = ugc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UsersGroupMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ugc.mutation = mutation
			node, err = ugc.sqlSave(ctx)
			return node, err
		})
		for i := len(ugc.hooks) - 1; i >= 0; i-- {
			mut = ugc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ugc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ugc *UsersGroupCreate) SaveX(ctx context.Context) *UsersGroup {
	v, err := ugc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ugc *UsersGroupCreate) sqlSave(ctx context.Context) (*UsersGroup, error) {
	var (
		ug    = &UsersGroup{config: ugc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: usersgroup.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: usersgroup.FieldID,
			},
		}
	)
	if value, ok := ugc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: usersgroup.FieldCreateTime,
		})
		ug.CreateTime = value
	}
	if value, ok := ugc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: usersgroup.FieldUpdateTime,
		})
		ug.UpdateTime = value
	}
	if value, ok := ugc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usersgroup.FieldName,
		})
		ug.Name = value
	}
	if value, ok := ugc.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usersgroup.FieldDescription,
		})
		ug.Description = value
	}
	if value, ok := ugc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: usersgroup.FieldStatus,
		})
		ug.Status = value
	}
	if nodes := ugc.mutation.MembersIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ugc.mutation.PoliciesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   usersgroup.PoliciesTable,
			Columns: usersgroup.PoliciesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: permissionspolicy.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, ugc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	ug.ID = int(id)
	return ug, nil
}
