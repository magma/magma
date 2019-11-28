// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/user"
)

// UserCreate is the builder for creating a User entity.
type UserCreate struct {
	config
	created_at *time.Time
	updated_at *time.Time
	email      *string
	password   *string
	role       *int
	tenant     *string
	networks   *[]string
	tabs       *[]string
}

// SetCreatedAt sets the created_at field.
func (uc *UserCreate) SetCreatedAt(t time.Time) *UserCreate {
	uc.created_at = &t
	return uc
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (uc *UserCreate) SetNillableCreatedAt(t *time.Time) *UserCreate {
	if t != nil {
		uc.SetCreatedAt(*t)
	}
	return uc
}

// SetUpdatedAt sets the updated_at field.
func (uc *UserCreate) SetUpdatedAt(t time.Time) *UserCreate {
	uc.updated_at = &t
	return uc
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (uc *UserCreate) SetNillableUpdatedAt(t *time.Time) *UserCreate {
	if t != nil {
		uc.SetUpdatedAt(*t)
	}
	return uc
}

// SetEmail sets the email field.
func (uc *UserCreate) SetEmail(s string) *UserCreate {
	uc.email = &s
	return uc
}

// SetPassword sets the password field.
func (uc *UserCreate) SetPassword(s string) *UserCreate {
	uc.password = &s
	return uc
}

// SetRole sets the role field.
func (uc *UserCreate) SetRole(i int) *UserCreate {
	uc.role = &i
	return uc
}

// SetNillableRole sets the role field if the given value is not nil.
func (uc *UserCreate) SetNillableRole(i *int) *UserCreate {
	if i != nil {
		uc.SetRole(*i)
	}
	return uc
}

// SetTenant sets the tenant field.
func (uc *UserCreate) SetTenant(s string) *UserCreate {
	uc.tenant = &s
	return uc
}

// SetNillableTenant sets the tenant field if the given value is not nil.
func (uc *UserCreate) SetNillableTenant(s *string) *UserCreate {
	if s != nil {
		uc.SetTenant(*s)
	}
	return uc
}

// SetNetworks sets the networks field.
func (uc *UserCreate) SetNetworks(s []string) *UserCreate {
	uc.networks = &s
	return uc
}

// SetTabs sets the tabs field.
func (uc *UserCreate) SetTabs(s []string) *UserCreate {
	uc.tabs = &s
	return uc
}

// Save creates the User in the database.
func (uc *UserCreate) Save(ctx context.Context) (*User, error) {
	if uc.created_at == nil {
		v := user.DefaultCreatedAt()
		uc.created_at = &v
	}
	if uc.updated_at == nil {
		v := user.DefaultUpdatedAt()
		uc.updated_at = &v
	}
	if uc.email == nil {
		return nil, errors.New("ent: missing required field \"email\"")
	}
	if err := user.EmailValidator(*uc.email); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
	}
	if uc.password == nil {
		return nil, errors.New("ent: missing required field \"password\"")
	}
	if err := user.PasswordValidator(*uc.password); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"password\": %v", err)
	}
	if uc.role == nil {
		v := user.DefaultRole
		uc.role = &v
	}
	if err := user.RoleValidator(*uc.role); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
	}
	if uc.tenant == nil {
		v := user.DefaultTenant
		uc.tenant = &v
	}
	if uc.networks == nil {
		return nil, errors.New("ent: missing required field \"networks\"")
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
		builder = sql.Dialect(uc.driver.Dialect())
		u       = &User{config: uc.config}
	)
	tx, err := uc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(user.Table).Default()
	if value := uc.created_at; value != nil {
		insert.Set(user.FieldCreatedAt, *value)
		u.CreatedAt = *value
	}
	if value := uc.updated_at; value != nil {
		insert.Set(user.FieldUpdatedAt, *value)
		u.UpdatedAt = *value
	}
	if value := uc.email; value != nil {
		insert.Set(user.FieldEmail, *value)
		u.Email = *value
	}
	if value := uc.password; value != nil {
		insert.Set(user.FieldPassword, *value)
		u.Password = *value
	}
	if value := uc.role; value != nil {
		insert.Set(user.FieldRole, *value)
		u.Role = *value
	}
	if value := uc.tenant; value != nil {
		insert.Set(user.FieldTenant, *value)
		u.Tenant = *value
	}
	if value := uc.networks; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		insert.Set(user.FieldNetworks, buf)
		u.Networks = *value
	}
	if value := uc.tabs; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		insert.Set(user.FieldTabs, buf)
		u.Tabs = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(user.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	u.ID = int(id)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}
