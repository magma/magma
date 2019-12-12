// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// Token is the model entity for the Token schema.
type Token struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Value holds the value of the "value" field.
	Value string `json:"-"`
}

// FromRows scans the sql response data into Token.
func (t *Token) FromRows(rows *sql.Rows) error {
	var scant struct {
		ID        int
		CreatedAt sql.NullTime
		UpdatedAt sql.NullTime
		Value     sql.NullString
	}
	// the order here should be the same as in the `token.Columns`.
	if err := rows.Scan(
		&scant.ID,
		&scant.CreatedAt,
		&scant.UpdatedAt,
		&scant.Value,
	); err != nil {
		return err
	}
	t.ID = scant.ID
	t.CreatedAt = scant.CreatedAt.Time
	t.UpdatedAt = scant.UpdatedAt.Time
	t.Value = scant.Value.String
	return nil
}

// QueryUser queries the user edge of the Token.
func (t *Token) QueryUser() *UserQuery {
	return (&TokenClient{t.config}).QueryUser(t)
}

// Update returns a builder for updating this Token.
// Note that, you need to call Token.Unwrap() before calling this method, if this Token
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Token) Update() *TokenUpdateOne {
	return (&TokenClient{t.config}).UpdateOne(t)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (t *Token) Unwrap() *Token {
	tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Token is not a transactional entity")
	}
	t.config.driver = tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Token) String() string {
	var builder strings.Builder
	builder.WriteString("Token(")
	builder.WriteString(fmt.Sprintf("id=%v", t.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(t.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(t.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", value=<sensitive>")
	builder.WriteByte(')')
	return builder.String()
}

// Tokens is a parsable slice of Token.
type Tokens []*Token

// FromRows scans the sql response data into Tokens.
func (t *Tokens) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scant := &Token{}
		if err := scant.FromRows(rows); err != nil {
			return err
		}
		*t = append(*t, scant)
	}
	return nil
}

func (t Tokens) config(cfg config) {
	for _i := range t {
		t[_i].config = cfg
	}
}
