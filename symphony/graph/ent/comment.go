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
	"github.com/facebookincubator/symphony/graph/ent/comment"
)

// Comment is the model entity for the Comment schema.
type Comment struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// AuthorName holds the value of the "author_name" field.
	AuthorName string `json:"author_name,omitempty"`
	// Text holds the value of the "text" field.
	Text                string `json:"text,omitempty"`
	project_comments    *int
	work_order_comments *int
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Comment) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // author_name
		&sql.NullString{}, // text
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Comment) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // project_comments
		&sql.NullInt64{}, // work_order_comments
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Comment fields.
func (c *Comment) assignValues(values ...interface{}) error {
	if m, n := len(values), len(comment.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	c.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		c.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		c.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field author_name", values[2])
	} else if value.Valid {
		c.AuthorName = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field text", values[3])
	} else if value.Valid {
		c.Text = value.String
	}
	values = values[4:]
	if len(values) == len(comment.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field project_comments", value)
		} else if value.Valid {
			c.project_comments = new(int)
			*c.project_comments = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_comments", value)
		} else if value.Valid {
			c.work_order_comments = new(int)
			*c.work_order_comments = int(value.Int64)
		}
	}
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

// Comments is a parsable slice of Comment.
type Comments []*Comment

func (c Comments) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
