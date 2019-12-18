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
	"github.com/facebookincubator/symphony/frontier/ent/token"
	"github.com/facebookincubator/symphony/frontier/ent/user"
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
		t    = &Token{config: tc.config}
		spec = &sqlgraph.CreateSpec{
			Table: token.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: token.FieldID,
			},
		}
	)
	if value := tc.created_at; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: token.FieldCreatedAt,
		})
		t.CreatedAt = *value
	}
	if value := tc.updated_at; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: token.FieldUpdatedAt,
		})
		t.UpdatedAt = *value
	}
	if value := tc.value; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: token.FieldValue,
		})
		t.Value = *value
	}
	if nodes := tc.user; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   token.UserTable,
			Columns: []string{token.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, tc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	t.ID = int(id)
	return t, nil
}
