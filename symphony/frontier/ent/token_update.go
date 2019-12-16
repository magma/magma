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
	var (
		builder  = sql.Dialect(tu.driver.Dialect())
		selector = builder.Select(token.FieldID).From(builder.Table(token.Table))
	)
	for _, p := range tu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = tu.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := tu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(token.Table)
	)
	updater = updater.Where(sql.InInts(token.FieldID, ids...))
	if value := tu.updated_at; value != nil {
		updater.Set(token.FieldUpdatedAt, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if tu.clearedUser {
		query, args := builder.Update(token.UserTable).
			SetNull(token.UserColumn).
			Where(sql.InInts(user.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(tu.user) > 0 {
		for eid := range tu.user {
			query, args := builder.Update(token.UserTable).
				Set(token.UserColumn, eid).
				Where(sql.InInts(token.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
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
	var (
		builder  = sql.Dialect(tuo.driver.Dialect())
		selector = builder.Select(token.Columns...).From(builder.Table(token.Table))
	)
	token.ID(tuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = tuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		t = &Token{config: tuo.config}
		if err := t.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Token: %v", err)
		}
		id = t.ID
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Token with id: %v", tuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Token with the same id: %v", tuo.id)
	}

	tx, err := tuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(token.Table)
	)
	updater = updater.Where(sql.InInts(token.FieldID, ids...))
	if value := tuo.updated_at; value != nil {
		updater.Set(token.FieldUpdatedAt, *value)
		t.UpdatedAt = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if tuo.clearedUser {
		query, args := builder.Update(token.UserTable).
			SetNull(token.UserColumn).
			Where(sql.InInts(user.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(tuo.user) > 0 {
		for eid := range tuo.user {
			query, args := builder.Update(token.UserTable).
				Set(token.UserColumn, eid).
				Where(sql.InInts(token.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}
