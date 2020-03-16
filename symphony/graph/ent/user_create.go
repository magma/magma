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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/user"
)

// UserCreate is the builder for creating a User entity.
type UserCreate struct {
	config
	create_time   *time.Time
	update_time   *time.Time
	auth_id       *string
	first_name    *string
	last_name     *string
	email         *string
	status        *user.Status
	role          *user.Role
	profile_photo map[int]struct{}
}

// SetCreateTime sets the create_time field.
func (uc *UserCreate) SetCreateTime(t time.Time) *UserCreate {
	uc.create_time = &t
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
	uc.update_time = &t
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
	uc.auth_id = &s
	return uc
}

// SetFirstName sets the first_name field.
func (uc *UserCreate) SetFirstName(s string) *UserCreate {
	uc.first_name = &s
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
	uc.last_name = &s
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
	uc.email = &s
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
	uc.status = &u
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
	uc.role = &u
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
	if uc.profile_photo == nil {
		uc.profile_photo = make(map[int]struct{})
	}
	uc.profile_photo[id] = struct{}{}
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

// Save creates the User in the database.
func (uc *UserCreate) Save(ctx context.Context) (*User, error) {
	if uc.create_time == nil {
		v := user.DefaultCreateTime()
		uc.create_time = &v
	}
	if uc.update_time == nil {
		v := user.DefaultUpdateTime()
		uc.update_time = &v
	}
	if uc.auth_id == nil {
		return nil, errors.New("ent: missing required field \"auth_id\"")
	}
	if err := user.AuthIDValidator(*uc.auth_id); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"auth_id\": %v", err)
	}
	if uc.first_name != nil {
		if err := user.FirstNameValidator(*uc.first_name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"first_name\": %v", err)
		}
	}
	if uc.last_name != nil {
		if err := user.LastNameValidator(*uc.last_name); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"last_name\": %v", err)
		}
	}
	if uc.email != nil {
		if err := user.EmailValidator(*uc.email); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	if uc.status == nil {
		v := user.DefaultStatus
		uc.status = &v
	}
	if err := user.StatusValidator(*uc.status); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"status\": %v", err)
	}
	if uc.role == nil {
		v := user.DefaultRole
		uc.role = &v
	}
	if err := user.RoleValidator(*uc.role); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
	}
	if len(uc.profile_photo) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"profile_photo\"")
	}
	return uc.sqlSave(ctx)
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
	if value := uc.create_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: user.FieldCreateTime,
		})
		u.CreateTime = *value
	}
	if value := uc.update_time; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: user.FieldUpdateTime,
		})
		u.UpdateTime = *value
	}
	if value := uc.auth_id; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldAuthID,
		})
		u.AuthID = *value
	}
	if value := uc.first_name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldFirstName,
		})
		u.FirstName = *value
	}
	if value := uc.last_name; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldLastName,
		})
		u.LastName = *value
	}
	if value := uc.email; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: user.FieldEmail,
		})
		u.Email = *value
	}
	if value := uc.status; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: user.FieldStatus,
		})
		u.Status = *value
	}
	if value := uc.role; value != nil {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  *value,
			Column: user.FieldRole,
		})
		u.Role = *value
	}
	if nodes := uc.profile_photo; len(nodes) > 0 {
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
