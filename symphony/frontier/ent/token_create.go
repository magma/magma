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
	"github.com/facebookincubator/symphony/frontier/ent/token"
)

// TokenCreate is the builder for creating a Token entity.
type TokenCreate struct {
	config
	created_at *time.Time
	updated_at *time.Time
	value      *string
	user       map[int]struct{}
}

// SetCreatedAt sets the created_at field.
func (tc *TokenCreate) SetCreatedAt(t time.Time) *TokenCreate {
	tc.created_at = &t
	return tc
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (tc *TokenCreate) SetNillableCreatedAt(t *time.Time) *TokenCreate {
	if t != nil {
		tc.SetCreatedAt(*t)
	}
	return tc
}

// SetUpdatedAt sets the updated_at field.
func (tc *TokenCreate) SetUpdatedAt(t time.Time) *TokenCreate {
	tc.updated_at = &t
	return tc
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (tc *TokenCreate) SetNillableUpdatedAt(t *time.Time) *TokenCreate {
	if t != nil {
		tc.SetUpdatedAt(*t)
	}
	return tc
}

// SetValue sets the value field.
func (tc *TokenCreate) SetValue(s string) *TokenCreate {
	tc.value = &s
	return tc
}

// SetUserID sets the user edge to User by id.
func (tc *TokenCreate) SetUserID(id int) *TokenCreate {
	if tc.user == nil {
		tc.user = make(map[int]struct{})
	}
	tc.user[id] = struct{}{}
	return tc
}

// SetUser sets the user edge to User.
func (tc *TokenCreate) SetUser(u *User) *TokenCreate {
	return tc.SetUserID(u.ID)
}

// Save creates the Token in the database.
func (tc *TokenCreate) Save(ctx context.Context) (*Token, error) {
	if tc.created_at == nil {
		v := token.DefaultCreatedAt()
		tc.created_at = &v
	}
	if tc.updated_at == nil {
		v := token.DefaultUpdatedAt()
		tc.updated_at = &v
	}
	if tc.value == nil {
		return nil, errors.New("ent: missing required field \"value\"")
	}
	if err := token.ValueValidator(*tc.value); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"value\": %v", err)
	}
	if len(tc.user) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"user\"")
	}
	if tc.user == nil {
		return nil, errors.New("ent: missing required edge \"user\"")
	}
	return tc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TokenCreate) SaveX(ctx context.Context) *Token {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tc *TokenCreate) sqlSave(ctx context.Context) (*Token, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(tc.driver.Dialect())
		t       = &Token{config: tc.config}
	)
	tx, err := tc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(token.Table).Default()
	if value := tc.created_at; value != nil {
		insert.Set(token.FieldCreatedAt, *value)
		t.CreatedAt = *value
	}
	if value := tc.updated_at; value != nil {
		insert.Set(token.FieldUpdatedAt, *value)
		t.UpdatedAt = *value
	}
	if value := tc.value; value != nil {
		insert.Set(token.FieldValue, *value)
		t.Value = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(token.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	t.ID = int(id)
	if len(tc.user) > 0 {
		for eid := range tc.user {
			query, args := builder.Update(token.UserTable).
				Set(token.UserColumn, eid).
				Where(sql.EQ(token.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}
