// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/facebookincubator/ent/dialect/sql"
)

// Todo is the model entity for the Todo schema.
type Todo struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Text holds the value of the "text" field.
	Text string `json:"text,omitempty"`
}

// FromRows scans the sql response data into Todo.
func (t *Todo) FromRows(rows *sql.Rows) error {
	var scant struct {
		ID   int
		Text sql.NullString
	}
	// the order here should be the same as in the `todo.Columns`.
	if err := rows.Scan(
		&scant.ID,
		&scant.Text,
	); err != nil {
		return err
	}
	t.ID = strconv.Itoa(scant.ID)
	t.Text = scant.Text.String
	return nil
}

// Update returns a builder for updating this Todo.
// Note that, you need to call Todo.Unwrap() before calling this method, if this Todo
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Todo) Update() *TodoUpdateOne {
	return (&TodoClient{t.config}).UpdateOne(t)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (t *Todo) Unwrap() *Todo {
	tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Todo is not a transactional entity")
	}
	t.config.driver = tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Todo) String() string {
	var builder strings.Builder
	builder.WriteString("Todo(")
	builder.WriteString(fmt.Sprintf("id=%v", t.ID))
	builder.WriteString(", text=")
	builder.WriteString(t.Text)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (t *Todo) id() int {
	id, _ := strconv.Atoi(t.ID)
	return id
}

// Todos is a parsable slice of Todo.
type Todos []*Todo

// FromRows scans the sql response data into Todos.
func (t *Todos) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scant := &Todo{}
		if err := scant.FromRows(rows); err != nil {
			return err
		}
		*t = append(*t, scant)
	}
	return nil
}

func (t Todos) config(cfg config) {
	for _i := range t {
		t[_i].config = cfg
	}
}
