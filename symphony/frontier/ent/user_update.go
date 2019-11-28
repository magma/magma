// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent/predicate"
	"github.com/facebookincubator/symphony/frontier/ent/user"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config

	updated_at *time.Time
	email      *string
	password   *string
	role       *int
	addrole    *int

	networks   *[]string
	tabs       *[]string
	cleartabs  bool
	predicates []predicate.User
}

// Where adds a new predicate for the builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.predicates = append(uu.predicates, ps...)
	return uu
}

// SetEmail sets the email field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.email = &s
	return uu
}

// SetPassword sets the password field.
func (uu *UserUpdate) SetPassword(s string) *UserUpdate {
	uu.password = &s
	return uu
}

// SetRole sets the role field.
func (uu *UserUpdate) SetRole(i int) *UserUpdate {
	uu.role = &i
	uu.addrole = nil
	return uu
}

// SetNillableRole sets the role field if the given value is not nil.
func (uu *UserUpdate) SetNillableRole(i *int) *UserUpdate {
	if i != nil {
		uu.SetRole(*i)
	}
	return uu
}

// AddRole adds i to role.
func (uu *UserUpdate) AddRole(i int) *UserUpdate {
	if uu.addrole == nil {
		uu.addrole = &i
	} else {
		*uu.addrole += i
	}
	return uu
}

// SetNetworks sets the networks field.
func (uu *UserUpdate) SetNetworks(s []string) *UserUpdate {
	uu.networks = &s
	return uu
}

// SetTabs sets the tabs field.
func (uu *UserUpdate) SetTabs(s []string) *UserUpdate {
	uu.tabs = &s
	return uu
}

// ClearTabs clears the value of tabs.
func (uu *UserUpdate) ClearTabs() *UserUpdate {
	uu.tabs = nil
	uu.cleartabs = true
	return uu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	if uu.updated_at == nil {
		v := user.UpdateDefaultUpdatedAt()
		uu.updated_at = &v
	}
	if uu.email != nil {
		if err := user.EmailValidator(*uu.email); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	if uu.password != nil {
		if err := user.PasswordValidator(*uu.password); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"password\": %v", err)
		}
	}
	if uu.role != nil {
		if err := user.RoleValidator(*uu.role); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
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
	var (
		builder  = sql.Dialect(uu.driver.Dialect())
		selector = builder.Select(user.FieldID).From(builder.Table(user.Table))
	)
	for _, p := range uu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = uu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := uu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(user.Table)
	)
	updater = updater.Where(sql.InInts(user.FieldID, ids...))
	if value := uu.updated_at; value != nil {
		updater.Set(user.FieldUpdatedAt, *value)
	}
	if value := uu.email; value != nil {
		updater.Set(user.FieldEmail, *value)
	}
	if value := uu.password; value != nil {
		updater.Set(user.FieldPassword, *value)
	}
	if value := uu.role; value != nil {
		updater.Set(user.FieldRole, *value)
	}
	if value := uu.addrole; value != nil {
		updater.Add(user.FieldRole, *value)
	}
	if value := uu.networks; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return 0, err
		}
		updater.Set(user.FieldNetworks, buf)
	}
	if value := uu.tabs; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return 0, err
		}
		updater.Set(user.FieldTabs, buf)
	}
	if uu.cleartabs {
		updater.SetNull(user.FieldTabs)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	id int

	updated_at *time.Time
	email      *string
	password   *string
	role       *int
	addrole    *int

	networks  *[]string
	tabs      *[]string
	cleartabs bool
}

// SetEmail sets the email field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.email = &s
	return uuo
}

// SetPassword sets the password field.
func (uuo *UserUpdateOne) SetPassword(s string) *UserUpdateOne {
	uuo.password = &s
	return uuo
}

// SetRole sets the role field.
func (uuo *UserUpdateOne) SetRole(i int) *UserUpdateOne {
	uuo.role = &i
	uuo.addrole = nil
	return uuo
}

// SetNillableRole sets the role field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableRole(i *int) *UserUpdateOne {
	if i != nil {
		uuo.SetRole(*i)
	}
	return uuo
}

// AddRole adds i to role.
func (uuo *UserUpdateOne) AddRole(i int) *UserUpdateOne {
	if uuo.addrole == nil {
		uuo.addrole = &i
	} else {
		*uuo.addrole += i
	}
	return uuo
}

// SetNetworks sets the networks field.
func (uuo *UserUpdateOne) SetNetworks(s []string) *UserUpdateOne {
	uuo.networks = &s
	return uuo
}

// SetTabs sets the tabs field.
func (uuo *UserUpdateOne) SetTabs(s []string) *UserUpdateOne {
	uuo.tabs = &s
	return uuo
}

// ClearTabs clears the value of tabs.
func (uuo *UserUpdateOne) ClearTabs() *UserUpdateOne {
	uuo.tabs = nil
	uuo.cleartabs = true
	return uuo
}

// Save executes the query and returns the updated entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	if uuo.updated_at == nil {
		v := user.UpdateDefaultUpdatedAt()
		uuo.updated_at = &v
	}
	if uuo.email != nil {
		if err := user.EmailValidator(*uuo.email); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"email\": %v", err)
		}
	}
	if uuo.password != nil {
		if err := user.PasswordValidator(*uuo.password); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"password\": %v", err)
		}
	}
	if uuo.role != nil {
		if err := user.RoleValidator(*uuo.role); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
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
	var (
		builder  = sql.Dialect(uuo.driver.Dialect())
		selector = builder.Select(user.Columns...).From(builder.Table(user.Table))
	)
	user.ID(uuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = uuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		u = &User{config: uuo.config}
		if err := u.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into User: %v", err)
		}
		id = u.ID
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("User with id: %v", uuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one User with the same id: %v", uuo.id)
	}

	tx, err := uuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(user.Table)
	)
	updater = updater.Where(sql.InInts(user.FieldID, ids...))
	if value := uuo.updated_at; value != nil {
		updater.Set(user.FieldUpdatedAt, *value)
		u.UpdatedAt = *value
	}
	if value := uuo.email; value != nil {
		updater.Set(user.FieldEmail, *value)
		u.Email = *value
	}
	if value := uuo.password; value != nil {
		updater.Set(user.FieldPassword, *value)
		u.Password = *value
	}
	if value := uuo.role; value != nil {
		updater.Set(user.FieldRole, *value)
		u.Role = *value
	}
	if value := uuo.addrole; value != nil {
		updater.Add(user.FieldRole, *value)
		u.Role += *value
	}
	if value := uuo.networks; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		updater.Set(user.FieldNetworks, buf)
		u.Networks = *value
	}
	if value := uuo.tabs; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		updater.Set(user.FieldTabs, buf)
		u.Tabs = *value
	}
	if uuo.cleartabs {
		var value []string
		u.Tabs = value
		updater.SetNull(user.FieldTabs)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return u, nil
}
