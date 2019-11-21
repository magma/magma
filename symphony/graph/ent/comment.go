// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// Comment is the model entity for the Comment schema.
type Comment struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// AuthorName holds the value of the "author_name" field.
	AuthorName string `json:"author_name,omitempty"`
	// Text holds the value of the "text" field.
	Text string `json:"text,omitempty"`
}

// FromRows scans the sql response data into Comment.
func (c *Comment) FromRows(rows *sql.Rows) error {
	var scanc struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		AuthorName sql.NullString
		Text       sql.NullString
	}
	// the order here should be the same as in the `comment.Columns`.
	if err := rows.Scan(
		&scanc.ID,
		&scanc.CreateTime,
		&scanc.UpdateTime,
		&scanc.AuthorName,
		&scanc.Text,
	); err != nil {
		return err
	}
	c.ID = strconv.Itoa(scanc.ID)
	c.CreateTime = scanc.CreateTime.Time
	c.UpdateTime = scanc.UpdateTime.Time
	c.AuthorName = scanc.AuthorName.String
	c.Text = scanc.Text.String
	return nil
}

// Update returns a builder for updating this Comment.
// Note that, you need to call Comment.Unwrap() before calling this method, if this Comment
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Comment) Update() *CommentUpdateOne {
	return (&CommentClient{c.config}).UpdateOne(c)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (c *Comment) Unwrap() *Comment {
	tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Comment is not a transactional entity")
	}
	c.config.driver = tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Comment) String() string {
	var builder strings.Builder
	builder.WriteString("Comment(")
	builder.WriteString(fmt.Sprintf("id=%v", c.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(c.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(c.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", author_name=")
	builder.WriteString(c.AuthorName)
	builder.WriteString(", text=")
	builder.WriteString(c.Text)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (c *Comment) id() int {
	id, _ := strconv.Atoi(c.ID)
	return id
}

// Comments is a parsable slice of Comment.
type Comments []*Comment

// FromRows scans the sql response data into Comments.
func (c *Comments) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanc := &Comment{}
		if err := scanc.FromRows(rows); err != nil {
			return err
		}
		*c = append(*c, scanc)
	}
	return nil
}

func (c Comments) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
