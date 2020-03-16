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

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/user"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config

	update_time *time.Time

	first_name          *string
	clearfirst_name     bool
	last_name           *string
	clearlast_name      bool
	email               *string
	clearemail          bool
	status              *user.Status
	role                *user.Role
	profile_photo       map[int]struct{}
	clearedProfilePhoto bool
	predicates          []predicate.User
}

// Where adds a new predicate for the builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.predicates = append(uu.predicates, ps...)
	return uu
}

// SetFirstName sets the first_name field.
func (uu *UserUpdate) SetFirstName(s string) *UserUpdate {
	uu.first_name = &s
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
	uu.first_name = nil
	uu.clearfirst_name = true
	return uu
}

// SetLastName sets the last_name field.
func (uu *UserUpdate) SetLastName(s string) *UserUpdate {
	uu.last_name = &s
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
	uu.last_name = nil
	uu.clearlast_name = true
	return uu
}

// SetEmail sets the email field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.email = &s
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
	uu.email = nil
	uu.clearemail = true
	return uu
}

// SetStatus sets the status field.
func (uu *UserUpdate) SetStatus(u user.Status) *UserUpdate {
	uu.status = &u
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
	uu.role = &u
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
	if uu.profile_photo == nil {
		uu.profile_photo = make(map[int]struct{})
	}
	uu.profile_photo[id] = struct{}{}
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

// ClearProfilePhoto clears the profile_photo edge to File.
func (uu *UserUpdate) ClearProfilePhoto() *UserUpdate {
	uu.clearedProfilePhoto = true
	return uu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	if uu.update_time == nil {
		v := user.UpdateDefaultUpdateTime()
		uu.update_time = &v
	}
	if uu.first_name != nil {
		if err := user.FirstNameValidator(*uu.first_name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"first_name\": %v", err)
		}
	}
	if uu.last_name != nil {
		if err := user.LastNameValidator(*uu.last_name); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"last_name\": %v", err)
		}
	}
	if uu.email != nil {
		if err := user.EmailValidator(*uu.email); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	if uu.status != nil {
		if err := user.StatusValidator(*uu.status); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
		}
	}
	if uu.role != nil {
		if err := user.RoleValidator(*uu.role); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
	}
	if len(uu.profile_photo) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"profile_photo\"")
	}
	return uu.sqlSave(ctx)
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
	if value := uu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: user.FieldUpdateTime,
		})
	}
	if value := uu.first_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldFirstName,
		})
	}
	if uu.clearfirst_name {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldFirstName,
		})
	}
	if value := uu.last_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldLastName,
		})
	}
	if uu.clearlast_name {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldLastName,
		})
	}
	if value := uu.email; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldEmail,
		})
	}
	if uu.clearemail {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldEmail,
		})
	}
	if value := uu.status; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: user.FieldStatus,
		})
	}
	if value := uu.role; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: user.FieldRole,
		})
	}
	if uu.clearedProfilePhoto {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
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
	if nodes := uu.profile_photo; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
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
		for k, _ := range nodes {
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
	id int

	update_time *time.Time

	first_name          *string
	clearfirst_name     bool
	last_name           *string
	clearlast_name      bool
	email               *string
	clearemail          bool
	status              *user.Status
	role                *user.Role
	profile_photo       map[int]struct{}
	clearedProfilePhoto bool
}

// SetFirstName sets the first_name field.
func (uuo *UserUpdateOne) SetFirstName(s string) *UserUpdateOne {
	uuo.first_name = &s
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
	uuo.first_name = nil
	uuo.clearfirst_name = true
	return uuo
}

// SetLastName sets the last_name field.
func (uuo *UserUpdateOne) SetLastName(s string) *UserUpdateOne {
	uuo.last_name = &s
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
	uuo.last_name = nil
	uuo.clearlast_name = true
	return uuo
}

// SetEmail sets the email field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.email = &s
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
	uuo.email = nil
	uuo.clearemail = true
	return uuo
}

// SetStatus sets the status field.
func (uuo *UserUpdateOne) SetStatus(u user.Status) *UserUpdateOne {
	uuo.status = &u
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
	uuo.role = &u
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
	if uuo.profile_photo == nil {
		uuo.profile_photo = make(map[int]struct{})
	}
	uuo.profile_photo[id] = struct{}{}
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

// ClearProfilePhoto clears the profile_photo edge to File.
func (uuo *UserUpdateOne) ClearProfilePhoto() *UserUpdateOne {
	uuo.clearedProfilePhoto = true
	return uuo
}

// Save executes the query and returns the updated entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	if uuo.update_time == nil {
		v := user.UpdateDefaultUpdateTime()
		uuo.update_time = &v
	}
	if uuo.first_name != nil {
		if err := user.FirstNameValidator(*uuo.first_name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"first_name\": %v", err)
		}
	}
	if uuo.last_name != nil {
		if err := user.LastNameValidator(*uuo.last_name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"last_name\": %v", err)
		}
	}
	if uuo.email != nil {
		if err := user.EmailValidator(*uuo.email); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	if uuo.status != nil {
		if err := user.StatusValidator(*uuo.status); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
		}
	}
	if uuo.role != nil {
		if err := user.RoleValidator(*uuo.role); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
	}
	if len(uuo.profile_photo) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"profile_photo\"")
	}
	return uuo.sqlSave(ctx)
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
				Value:  uuo.id,
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		},
	}
	if value := uuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: user.FieldUpdateTime,
		})
	}
	if value := uuo.first_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldFirstName,
		})
	}
	if uuo.clearfirst_name {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldFirstName,
		})
	}
	if value := uuo.last_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldLastName,
		})
	}
	if uuo.clearlast_name {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldLastName,
		})
	}
	if value := uuo.email; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldEmail,
		})
	}
	if uuo.clearemail {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldEmail,
		})
	}
	if value := uuo.status; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: user.FieldStatus,
		})
	}
	if value := uuo.role; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: user.FieldRole,
		})
	}
	if uuo.clearedProfilePhoto {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
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
	if nodes := uuo.profile_photo; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
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
		for k, _ := range nodes {
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
