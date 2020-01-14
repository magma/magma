// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/ent-integrations/relay/internal/todo/ent/todo"
)

// TodoCreate is the builder for creating a Todo entity.
type TodoCreate struct {
	config
	text *string
}

// SetText sets the text field.
func (tc *TodoCreate) SetText(s string) *TodoCreate {
	tc.text = &s
	return tc
}

// Save creates the Todo in the database.
func (tc *TodoCreate) Save(ctx context.Context) (*Todo, error) {
	if tc.text == nil {
		return nil, errors.New("ent: missing required field \"text\"")
	}
	if err := todo.TextValidator(*tc.text); err != nil {
		return nil, fmt.Errorf("ent: validator failed for field \"text\": %v", err)
	}
	return tc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TodoCreate) SaveX(ctx context.Context) *Todo {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (tc *TodoCreate) sqlSave(ctx context.Context) (*Todo, error) {
	var (
		t    = &Todo{config: tc.config}
		spec = &sqlgraph.CreateSpec{
			Table: todo.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: todo.FieldID,
			},
		}
	)
	if value := tc.text; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: todo.FieldText,
		})
		t.Text = *value
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
