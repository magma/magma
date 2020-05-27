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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks      []Hook
	mutation   *UserMutation
	predicates []predicate.User
}

// Where adds a new predicate for the builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.predicates = append(uu.predicates, ps...)
	return uu
}

// SetFirstName sets the first_name field.
func (uu *UserUpdate) SetFirstName(s string) *UserUpdate {
	uu.mutation.SetFirstName(s)
	return uu
}

// SetNillableFirstName sets the first_name field if the given value is not nil.
func (uu *UserUpdate) SetNillableFirstName(s *string) *UserUpdate {
	if s != nil {
		uu.SetFirstName(*s)
	}
	return uu
}

// ClearFirstName clears the value of first_name.
func (uu *UserUpdate) ClearFirstName() *UserUpdate {
	uu.mutation.ClearFirstName()
	return uu
}

// SetLastName sets the last_name field.
func (uu *UserUpdate) SetLastName(s string) *UserUpdate {
	uu.mutation.SetLastName(s)
	return uu
}

// SetNillableLastName sets the last_name field if the given value is not nil.
func (uu *UserUpdate) SetNillableLastName(s *string) *UserUpdate {
	if s != nil {
		uu.SetLastName(*s)
	}
	return uu
}

// ClearLastName clears the value of last_name.
func (uu *UserUpdate) ClearLastName() *UserUpdate {
	uu.mutation.ClearLastName()
	return uu
}

// SetEmail sets the email field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.mutation.SetEmail(s)
	return uu
}

// SetNillableEmail sets the email field if the given value is not nil.
func (uu *UserUpdate) SetNillableEmail(s *string) *UserUpdate {
	if s != nil {
		uu.SetEmail(*s)
	}
	return uu
}

// ClearEmail clears the value of email.
func (uu *UserUpdate) ClearEmail() *UserUpdate {
	uu.mutation.ClearEmail()
	return uu
}

// SetStatus sets the status field.
func (uu *UserUpdate) SetStatus(u user.Status) *UserUpdate {
	uu.mutation.SetStatus(u)
	return uu
}

// SetNillableStatus sets the status field if the given value is not nil.
func (uu *UserUpdate) SetNillableStatus(u *user.Status) *UserUpdate {
	if u != nil {
		uu.SetStatus(*u)
	}
	return uu
}

// SetRole sets the role field.
func (uu *UserUpdate) SetRole(u user.Role) *UserUpdate {
	uu.mutation.SetRole(u)
	return uu
}

// SetNillableRole sets the role field if the given value is not nil.
func (uu *UserUpdate) SetNillableRole(u *user.Role) *UserUpdate {
	if u != nil {
		uu.SetRole(*u)
	}
	return uu
}

// SetProfilePhotoID sets the profile_photo edge to File by id.
func (uu *UserUpdate) SetProfilePhotoID(id int) *UserUpdate {
	uu.mutation.SetProfilePhotoID(id)
	return uu
}

// SetNillableProfilePhotoID sets the profile_photo edge to File by id if the given value is not nil.
func (uu *UserUpdate) SetNillableProfilePhotoID(id *int) *UserUpdate {
	if id != nil {
		uu = uu.SetProfilePhotoID(*id)
	}
	return uu
}

// SetProfilePhoto sets the profile_photo edge to File.
func (uu *UserUpdate) SetProfilePhoto(f *File) *UserUpdate {
	return uu.SetProfilePhotoID(f.ID)
}

// AddGroupIDs adds the groups edge to UsersGroup by ids.
func (uu *UserUpdate) AddGroupIDs(ids ...int) *UserUpdate {
	uu.mutation.AddGroupIDs(ids...)
	return uu
}

// AddGroups adds the groups edges to UsersGroup.
func (uu *UserUpdate) AddGroups(u ...*UsersGroup) *UserUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uu.AddGroupIDs(ids...)
}

// AddOwnedWorkOrderIDs adds the owned_work_orders edge to WorkOrder by ids.
func (uu *UserUpdate) AddOwnedWorkOrderIDs(ids ...int) *UserUpdate {
	uu.mutation.AddOwnedWorkOrderIDs(ids...)
	return uu
}

// AddOwnedWorkOrders adds the owned_work_orders edges to WorkOrder.
func (uu *UserUpdate) AddOwnedWorkOrders(w ...*WorkOrder) *UserUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uu.AddOwnedWorkOrderIDs(ids...)
}

// AddAssignedWorkOrderIDs adds the assigned_work_orders edge to WorkOrder by ids.
func (uu *UserUpdate) AddAssignedWorkOrderIDs(ids ...int) *UserUpdate {
	uu.mutation.AddAssignedWorkOrderIDs(ids...)
	return uu
}

// AddAssignedWorkOrders adds the assigned_work_orders edges to WorkOrder.
func (uu *UserUpdate) AddAssignedWorkOrders(w ...*WorkOrder) *UserUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uu.AddAssignedWorkOrderIDs(ids...)
}

// AddCreatedProjectIDs adds the created_projects edge to Project by ids.
func (uu *UserUpdate) AddCreatedProjectIDs(ids ...int) *UserUpdate {
	uu.mutation.AddCreatedProjectIDs(ids...)
	return uu
}

// AddCreatedProjects adds the created_projects edges to Project.
func (uu *UserUpdate) AddCreatedProjects(p ...*Project) *UserUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return uu.AddCreatedProjectIDs(ids...)
}

// ClearProfilePhoto clears the profile_photo edge to File.
func (uu *UserUpdate) ClearProfilePhoto() *UserUpdate {
	uu.mutation.ClearProfilePhoto()
	return uu
}

// RemoveGroupIDs removes the groups edge to UsersGroup by ids.
func (uu *UserUpdate) RemoveGroupIDs(ids ...int) *UserUpdate {
	uu.mutation.RemoveGroupIDs(ids...)
	return uu
}

// RemoveGroups removes groups edges to UsersGroup.
func (uu *UserUpdate) RemoveGroups(u ...*UsersGroup) *UserUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uu.RemoveGroupIDs(ids...)
}

// RemoveOwnedWorkOrderIDs removes the owned_work_orders edge to WorkOrder by ids.
func (uu *UserUpdate) RemoveOwnedWorkOrderIDs(ids ...int) *UserUpdate {
	uu.mutation.RemoveOwnedWorkOrderIDs(ids...)
	return uu
}

// RemoveOwnedWorkOrders removes owned_work_orders edges to WorkOrder.
func (uu *UserUpdate) RemoveOwnedWorkOrders(w ...*WorkOrder) *UserUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uu.RemoveOwnedWorkOrderIDs(ids...)
}

// RemoveAssignedWorkOrderIDs removes the assigned_work_orders edge to WorkOrder by ids.
func (uu *UserUpdate) RemoveAssignedWorkOrderIDs(ids ...int) *UserUpdate {
	uu.mutation.RemoveAssignedWorkOrderIDs(ids...)
	return uu
}

// RemoveAssignedWorkOrders removes assigned_work_orders edges to WorkOrder.
func (uu *UserUpdate) RemoveAssignedWorkOrders(w ...*WorkOrder) *UserUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uu.RemoveAssignedWorkOrderIDs(ids...)
}

// RemoveCreatedProjectIDs removes the created_projects edge to Project by ids.
func (uu *UserUpdate) RemoveCreatedProjectIDs(ids ...int) *UserUpdate {
	uu.mutation.RemoveCreatedProjectIDs(ids...)
	return uu
}

// RemoveCreatedProjects removes created_projects edges to Project.
func (uu *UserUpdate) RemoveCreatedProjects(p ...*Project) *UserUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return uu.RemoveCreatedProjectIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := uu.mutation.UpdateTime(); !ok {
		v := user.UpdateDefaultUpdateTime()
		uu.mutation.SetUpdateTime(v)
	}
	if v, ok := uu.mutation.FirstName(); ok {
		if err := user.FirstNameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"first_name\": %v", err)
		}
	}
	if v, ok := uu.mutation.LastName(); ok {
		if err := user.LastNameValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"last_name\": %v", err)
		}
	}
	if v, ok := uu.mutation.Email(); ok {
		if err := user.EmailValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	if v, ok := uu.mutation.Status(); ok {
		if err := user.StatusValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
		}
	}
	if v, ok := uu.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(uu.hooks) == 0 {
		affected, err = uu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			uu.mutation = mutation
			affected, err = uu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(uu.hooks) - 1; i >= 0; i-- {
			mut = uu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		},
	}
	if ps := uu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldUpdateTime,
		})
	}
	if value, ok := uu.mutation.FirstName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldFirstName,
		})
	}
	if uu.mutation.FirstNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldFirstName,
		})
	}
	if value, ok := uu.mutation.LastName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldLastName,
		})
	}
	if uu.mutation.LastNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldLastName,
		})
	}
	if value, ok := uu.mutation.Email(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmail,
		})
	}
	if uu.mutation.EmailCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldEmail,
		})
	}
	if value, ok := uu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldStatus,
		})
	}
	if value, ok := uu.mutation.Role(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldRole,
		})
	}
	if uu.mutation.ProfilePhotoCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.ProfilePhotoIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uu.mutation.RemovedGroupsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.GroupsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uu.mutation.RemovedOwnedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.OwnedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uu.mutation.RemovedAssignedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.AssignedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uu.mutation.RemovedCreatedProjectsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.CreatedProjectsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// SetFirstName sets the first_name field.
func (uuo *UserUpdateOne) SetFirstName(s string) *UserUpdateOne {
	uuo.mutation.SetFirstName(s)
	return uuo
}

// SetNillableFirstName sets the first_name field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableFirstName(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetFirstName(*s)
	}
	return uuo
}

// ClearFirstName clears the value of first_name.
func (uuo *UserUpdateOne) ClearFirstName() *UserUpdateOne {
	uuo.mutation.ClearFirstName()
	return uuo
}

// SetLastName sets the last_name field.
func (uuo *UserUpdateOne) SetLastName(s string) *UserUpdateOne {
	uuo.mutation.SetLastName(s)
	return uuo
}

// SetNillableLastName sets the last_name field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableLastName(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetLastName(*s)
	}
	return uuo
}

// ClearLastName clears the value of last_name.
func (uuo *UserUpdateOne) ClearLastName() *UserUpdateOne {
	uuo.mutation.ClearLastName()
	return uuo
}

// SetEmail sets the email field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.mutation.SetEmail(s)
	return uuo
}

// SetNillableEmail sets the email field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableEmail(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetEmail(*s)
	}
	return uuo
}

// ClearEmail clears the value of email.
func (uuo *UserUpdateOne) ClearEmail() *UserUpdateOne {
	uuo.mutation.ClearEmail()
	return uuo
}

// SetStatus sets the status field.
func (uuo *UserUpdateOne) SetStatus(u user.Status) *UserUpdateOne {
	uuo.mutation.SetStatus(u)
	return uuo
}

// SetNillableStatus sets the status field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableStatus(u *user.Status) *UserUpdateOne {
	if u != nil {
		uuo.SetStatus(*u)
	}
	return uuo
}

// SetRole sets the role field.
func (uuo *UserUpdateOne) SetRole(u user.Role) *UserUpdateOne {
	uuo.mutation.SetRole(u)
	return uuo
}

// SetNillableRole sets the role field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableRole(u *user.Role) *UserUpdateOne {
	if u != nil {
		uuo.SetRole(*u)
	}
	return uuo
}

// SetProfilePhotoID sets the profile_photo edge to File by id.
func (uuo *UserUpdateOne) SetProfilePhotoID(id int) *UserUpdateOne {
	uuo.mutation.SetProfilePhotoID(id)
	return uuo
}

// SetNillableProfilePhotoID sets the profile_photo edge to File by id if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableProfilePhotoID(id *int) *UserUpdateOne {
	if id != nil {
		uuo = uuo.SetProfilePhotoID(*id)
	}
	return uuo
}

// SetProfilePhoto sets the profile_photo edge to File.
func (uuo *UserUpdateOne) SetProfilePhoto(f *File) *UserUpdateOne {
	return uuo.SetProfilePhotoID(f.ID)
}

// AddGroupIDs adds the groups edge to UsersGroup by ids.
func (uuo *UserUpdateOne) AddGroupIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.AddGroupIDs(ids...)
	return uuo
}

// AddGroups adds the groups edges to UsersGroup.
func (uuo *UserUpdateOne) AddGroups(u ...*UsersGroup) *UserUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uuo.AddGroupIDs(ids...)
}

// AddOwnedWorkOrderIDs adds the owned_work_orders edge to WorkOrder by ids.
func (uuo *UserUpdateOne) AddOwnedWorkOrderIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.AddOwnedWorkOrderIDs(ids...)
	return uuo
}

// AddOwnedWorkOrders adds the owned_work_orders edges to WorkOrder.
func (uuo *UserUpdateOne) AddOwnedWorkOrders(w ...*WorkOrder) *UserUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uuo.AddOwnedWorkOrderIDs(ids...)
}

// AddAssignedWorkOrderIDs adds the assigned_work_orders edge to WorkOrder by ids.
func (uuo *UserUpdateOne) AddAssignedWorkOrderIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.AddAssignedWorkOrderIDs(ids...)
	return uuo
}

// AddAssignedWorkOrders adds the assigned_work_orders edges to WorkOrder.
func (uuo *UserUpdateOne) AddAssignedWorkOrders(w ...*WorkOrder) *UserUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uuo.AddAssignedWorkOrderIDs(ids...)
}

// AddCreatedProjectIDs adds the created_projects edge to Project by ids.
func (uuo *UserUpdateOne) AddCreatedProjectIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.AddCreatedProjectIDs(ids...)
	return uuo
}

// AddCreatedProjects adds the created_projects edges to Project.
func (uuo *UserUpdateOne) AddCreatedProjects(p ...*Project) *UserUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return uuo.AddCreatedProjectIDs(ids...)
}

// ClearProfilePhoto clears the profile_photo edge to File.
func (uuo *UserUpdateOne) ClearProfilePhoto() *UserUpdateOne {
	uuo.mutation.ClearProfilePhoto()
	return uuo
}

// RemoveGroupIDs removes the groups edge to UsersGroup by ids.
func (uuo *UserUpdateOne) RemoveGroupIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.RemoveGroupIDs(ids...)
	return uuo
}

// RemoveGroups removes groups edges to UsersGroup.
func (uuo *UserUpdateOne) RemoveGroups(u ...*UsersGroup) *UserUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uuo.RemoveGroupIDs(ids...)
}

// RemoveOwnedWorkOrderIDs removes the owned_work_orders edge to WorkOrder by ids.
func (uuo *UserUpdateOne) RemoveOwnedWorkOrderIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.RemoveOwnedWorkOrderIDs(ids...)
	return uuo
}

// RemoveOwnedWorkOrders removes owned_work_orders edges to WorkOrder.
func (uuo *UserUpdateOne) RemoveOwnedWorkOrders(w ...*WorkOrder) *UserUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uuo.RemoveOwnedWorkOrderIDs(ids...)
}

// RemoveAssignedWorkOrderIDs removes the assigned_work_orders edge to WorkOrder by ids.
func (uuo *UserUpdateOne) RemoveAssignedWorkOrderIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.RemoveAssignedWorkOrderIDs(ids...)
	return uuo
}

// RemoveAssignedWorkOrders removes assigned_work_orders edges to WorkOrder.
func (uuo *UserUpdateOne) RemoveAssignedWorkOrders(w ...*WorkOrder) *UserUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return uuo.RemoveAssignedWorkOrderIDs(ids...)
}

// RemoveCreatedProjectIDs removes the created_projects edge to Project by ids.
func (uuo *UserUpdateOne) RemoveCreatedProjectIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.RemoveCreatedProjectIDs(ids...)
	return uuo
}

// RemoveCreatedProjects removes created_projects edges to Project.
func (uuo *UserUpdateOne) RemoveCreatedProjects(p ...*Project) *UserUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return uuo.RemoveCreatedProjectIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	if _, ok := uuo.mutation.UpdateTime(); !ok {
		v := user.UpdateDefaultUpdateTime()
		uuo.mutation.SetUpdateTime(v)
	}
	if v, ok := uuo.mutation.FirstName(); ok {
		if err := user.FirstNameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"first_name\": %v", err)
		}
	}
	if v, ok := uuo.mutation.LastName(); ok {
		if err := user.LastNameValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"last_name\": %v", err)
		}
	}
	if v, ok := uuo.mutation.Email(); ok {
		if err := user.EmailValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	if v, ok := uuo.mutation.Status(); ok {
		if err := user.StatusValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
		}
	}
	if v, ok := uuo.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
	}

	var (
		err  error
		node *User
	)
	if len(uuo.hooks) == 0 {
		node, err = uuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			uuo.mutation = mutation
			node, err = uuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(uuo.hooks) - 1; i >= 0; i-- {
			mut = uuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	u, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return u
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (u *User, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		},
	}
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing User.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := uuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldUpdateTime,
		})
	}
	if value, ok := uuo.mutation.FirstName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldFirstName,
		})
	}
	if uuo.mutation.FirstNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldFirstName,
		})
	}
	if value, ok := uuo.mutation.LastName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldLastName,
		})
	}
	if uuo.mutation.LastNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldLastName,
		})
	}
	if value, ok := uuo.mutation.Email(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmail,
		})
	}
	if uuo.mutation.EmailCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldEmail,
		})
	}
	if value, ok := uuo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldStatus,
		})
	}
	if value, ok := uuo.mutation.Role(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldRole,
		})
	}
	if uuo.mutation.ProfilePhotoCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.ProfilePhotoIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uuo.mutation.RemovedGroupsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.GroupsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uuo.mutation.RemovedOwnedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.OwnedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uuo.mutation.RemovedAssignedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.AssignedWorkOrdersIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uuo.mutation.RemovedCreatedProjectsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.CreatedProjectsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	u = &User{config: uuo.config}
	_spec.Assign = u.assignValues
	_spec.ScanValues = u.scanValues()
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return u, nil
}
