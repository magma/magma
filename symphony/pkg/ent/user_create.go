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
	"github.com/facebookincubator/symphony/pkg/ent/file"
	"github.com/facebookincubator/symphony/pkg/ent/project"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/usersgroup"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

// UserCreate is the builder for creating a User entity.
type UserCreate struct {
	config
	mutation *UserMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (uc *UserCreate) SetCreateTime(t time.Time) *UserCreate {
	uc.mutation.SetCreateTime(t)
	return uc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (uc *UserCreate) SetNillableCreateTime(t *time.Time) *UserCreate {
	if t != nil {
		uc.SetCreateTime(*t)
	}
	return uc
}

// SetUpdateTime sets the update_time field.
func (uc *UserCreate) SetUpdateTime(t time.Time) *UserCreate {
	uc.mutation.SetUpdateTime(t)
	return uc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (uc *UserCreate) SetNillableUpdateTime(t *time.Time) *UserCreate {
	if t != nil {
		uc.SetUpdateTime(*t)
	}
	return uc
}

// SetAuthID sets the auth_id field.
func (uc *UserCreate) SetAuthID(s string) *UserCreate {
	uc.mutation.SetAuthID(s)
	return uc
}

// SetFirstName sets the first_name field.
func (uc *UserCreate) SetFirstName(s string) *UserCreate {
	uc.mutation.SetFirstName(s)
	return uc
}

// SetNillableFirstName sets the first_name field if the given value is not nil.
func (uc *UserCreate) SetNillableFirstName(s *string) *UserCreate {
	if s != nil {
		uc.SetFirstName(*s)
	}
	return uc
}

// SetLastName sets the last_name field.
func (uc *UserCreate) SetLastName(s string) *UserCreate {
	uc.mutation.SetLastName(s)
	return uc
}

// SetNillableLastName sets the last_name field if the given value is not nil.
func (uc *UserCreate) SetNillableLastName(s *string) *UserCreate {
	if s != nil {
		uc.SetLastName(*s)
	}
	return uc
}

// SetEmail sets the email field.
func (uc *UserCreate) SetEmail(s string) *UserCreate {
	uc.mutation.SetEmail(s)
	return uc
}

// SetNillableEmail sets the email field if the given value is not nil.
func (uc *UserCreate) SetNillableEmail(s *string) *UserCreate {
	if s != nil {
		uc.SetEmail(*s)
	}
	return uc
}

// SetStatus sets the status field.
func (uc *UserCreate) SetStatus(u user.Status) *UserCreate {
	uc.mutation.SetStatus(u)
	return uc
}

// SetNillableStatus sets the status field if the given value is not nil.
func (uc *UserCreate) SetNillableStatus(u *user.Status) *UserCreate {
	if u != nil {
		uc.SetStatus(*u)
	}
	return uc
}

// SetRole sets the role field.
func (uc *UserCreate) SetRole(u user.Role) *UserCreate {
	uc.mutation.SetRole(u)
	return uc
}

// SetNillableRole sets the role field if the given value is not nil.
func (uc *UserCreate) SetNillableRole(u *user.Role) *UserCreate {
	if u != nil {
		uc.SetRole(*u)
	}
	return uc
}

// SetProfilePhotoID sets the profile_photo edge to File by id.
func (uc *UserCreate) SetProfilePhotoID(id int) *UserCreate {
	uc.mutation.SetProfilePhotoID(id)
	return uc
}

// SetNillableProfilePhotoID sets the profile_photo edge to File by id if the given value is not nil.
func (uc *UserCreate) SetNillableProfilePhotoID(id *int) *UserCreate {
	if id != nil {
		uc = uc.SetProfilePhotoID(*id)
	}
	return uc
}

// SetProfilePhoto sets the profile_photo edge to File.
func (uc *UserCreate) SetProfilePhoto(f *File) *UserCreate {
	return uc.SetProfilePhotoID(f.ID)
}

// AddGroupIDs adds the groups edge to UsersGroup by ids.
func (uc *UserCreate) AddGroupIDs(ids ...int) *UserCreate {
	uc.mutation.AddGroupIDs(ids...)
	return uc
}

// AddGroups adds the groups edges to UsersGroup.
func (uc *UserCreate) AddGroups(u ...*UsersGroup) *UserCreate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uc.AddGroupIDs(ids...)
}

// AddOwnedWorkOrderIDs adds the owned_work_orders edge to WorkOrder by ids.
func (uc *UserCreate) AddOwnedWorkOrderIDs(ids ...int) *UserCreate {
	uc.mutation.AddOwnedWorkOrderIDs(ids...)
	return uc
}

// AddOwnedWorkOrders adds the owned_work_orders edges to WorkOrder.
func (uc *UserCreate) AddOwnedWorkOrders(w ...*WorkOrder) *UserCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uc.AddOwnedWorkOrderIDs(ids...)
}

// AddAssignedWorkOrderIDs adds the assigned_work_orders edge to WorkOrder by ids.
func (uc *UserCreate) AddAssignedWorkOrderIDs(ids ...int) *UserCreate {
	uc.mutation.AddAssignedWorkOrderIDs(ids...)
	return uc
}

// AddAssignedWorkOrders adds the assigned_work_orders edges to WorkOrder.
func (uc *UserCreate) AddAssignedWorkOrders(w ...*WorkOrder) *UserCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uc.AddAssignedWorkOrderIDs(ids...)
}

// AddCreatedProjectIDs adds the created_projects edge to Project by ids.
func (uc *UserCreate) AddCreatedProjectIDs(ids ...int) *UserCreate {
	uc.mutation.AddCreatedProjectIDs(ids...)
	return uc
}

// AddCreatedProjects adds the created_projects edges to Project.
func (uc *UserCreate) AddCreatedProjects(p ...*Project) *UserCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return uc.AddCreatedProjectIDs(ids...)
}

// Save creates the User in the database.
func (uc *UserCreate) Save(ctx context.Context) (*User, error) {
	if _, ok := uc.mutation.CreateTime(); !ok {
		v := user.DefaultCreateTime()
		uc.mutation.SetCreateTime(v)
	}
	if _, ok := uc.mutation.UpdateTime(); !ok {
		v := user.DefaultUpdateTime()
		uc.mutation.SetUpdateTime(v)
	}
	if _, ok := uc.mutation.AuthID(); !ok {
		return nil, errors.New("ent: missing required field \"auth_id\"")
	}
	if v, ok := uc.mutation.AuthID(); ok {
		if err := user.AuthIDValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"auth_id\": %v", err)
		}
	}
	if v, ok := uc.mutation.FirstName(); ok {
		if err := user.FirstNameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"first_name\": %v", err)
		}
	}
	if v, ok := uc.mutation.LastName(); ok {
		if err := user.LastNameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"last_name\": %v", err)
		}
	}
	if v, ok := uc.mutation.Email(); ok {
		if err := user.EmailValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	if _, ok := uc.mutation.Status(); !ok {
		v := user.DefaultStatus
		uc.mutation.SetStatus(v)
	}
	if v, ok := uc.mutation.Status(); ok {
		if err := user.StatusValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
		}
	}
	if _, ok := uc.mutation.Role(); !ok {
		v := user.DefaultRole
		uc.mutation.SetRole(v)
	}
	if v, ok := uc.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
	}
	var (
		err  error
		node *User
	)
	if len(uc.hooks) == 0 {
		node, err = uc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			uc.mutation = mutation
			node, err = uc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(uc.hooks) - 1; i >= 0; i-- {
			mut = uc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (uc *UserCreate) SaveX(ctx context.Context) *User {
	v, err := uc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (uc *UserCreate) sqlSave(ctx context.Context) (*User, error) {
	var (
		u     = &User{config: uc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: user.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		}
	)
	if value, ok := uc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldCreateTime,
		})
		u.CreateTime = value
	}
	if value, ok := uc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldUpdateTime,
		})
		u.UpdateTime = value
	}
	if value, ok := uc.mutation.AuthID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldAuthID,
		})
		u.AuthID = value
	}
	if value, ok := uc.mutation.FirstName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldFirstName,
		})
		u.FirstName = value
	}
	if value, ok := uc.mutation.LastName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldLastName,
		})
		u.LastName = value
	}
	if value, ok := uc.mutation.Email(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmail,
		})
		u.Email = value
	}
	if value, ok := uc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldStatus,
		})
		u.Status = value
	}
	if value, ok := uc.mutation.Role(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldRole,
		})
		u.Role = value
	}
	if nodes := uc.mutation.ProfilePhotoIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.ProfilePhotoTable,
			Columns: []string{user.ProfilePhotoColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.GroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GroupsTable,
			Columns: user.GroupsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersgroup.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.OwnedWorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   user.OwnedWorkOrdersTable,
			Columns: []string{user.OwnedWorkOrdersColumn},
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
	if nodes := uc.mutation.AssignedWorkOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   user.AssignedWorkOrdersTable,
			Columns: []string{user.AssignedWorkOrdersColumn},
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
	if nodes := uc.mutation.CreatedProjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   user.CreatedProjectsTable,
			Columns: []string{user.CreatedProjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: project.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, uc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	u.ID = int(id)
	return u, nil
}
