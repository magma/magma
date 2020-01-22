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
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
)

// Hyperlink is the model entity for the Hyperlink schema.
type Hyperlink struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty" gqlgen:"displayName"`
	// Category holds the value of the "category" field.
	Category                string `json:"category,omitempty"`
	equipment_hyperlink_id  *string
	location_hyperlink_id   *string
	work_order_hyperlink_id *string
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Hyperlink) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // url
		&sql.NullString{}, // name
		&sql.NullString{}, // category
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Hyperlink) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // equipment_hyperlink_id
		&sql.NullInt64{}, // location_hyperlink_id
		&sql.NullInt64{}, // work_order_hyperlink_id
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Hyperlink fields.
func (h *Hyperlink) assignValues(values ...interface{}) error {
	if m, n := len(values), len(hyperlink.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	h.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		h.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		h.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field url", values[2])
	} else if value.Valid {
		h.URL = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[3])
	} else if value.Valid {
		h.Name = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field category", values[4])
	} else if value.Valid {
		h.Category = value.String
	}
	values = values[5:]
	if len(values) == len(hyperlink.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field equipment_hyperlink_id", value)
		} else if value.Valid {
			h.equipment_hyperlink_id = new(string)
			*h.equipment_hyperlink_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_hyperlink_id", value)
		} else if value.Valid {
			h.location_hyperlink_id = new(string)
			*h.location_hyperlink_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_hyperlink_id", value)
		} else if value.Valid {
			h.work_order_hyperlink_id = new(string)
			*h.work_order_hyperlink_id = strconv.FormatInt(value.Int64, 10)
		}
	}
	return nil
}

// Update returns a builder for updating this Hyperlink.
// Note that, you need to call Hyperlink.Unwrap() before calling this method, if this Hyperlink
// was returned from a transaction, and the transaction was committed or rolled back.
func (h *Hyperlink) Update() *HyperlinkUpdateOne {
	return (&HyperlinkClient{h.config}).UpdateOne(h)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (h *Hyperlink) Unwrap() *Hyperlink {
	tx, ok := h.config.driver.(*txDriver)
	if !ok {
		panic("ent: Hyperlink is not a transactional entity")
	}
	h.config.driver = tx.drv
	return h
}

// String implements the fmt.Stringer.
func (h *Hyperlink) String() string {
	var builder strings.Builder
	builder.WriteString("Hyperlink(")
	builder.WriteString(fmt.Sprintf("id=%v", h.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(h.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(h.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", url=")
	builder.WriteString(h.URL)
	builder.WriteString(", name=")
	builder.WriteString(h.Name)
	builder.WriteString(", category=")
	builder.WriteString(h.Category)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (h *Hyperlink) id() int {
	id, _ := strconv.Atoi(h.ID)
	return id
}

// Hyperlinks is a parsable slice of Hyperlink.
type Hyperlinks []*Hyperlink

func (h Hyperlinks) config(cfg config) {
	for _i := range h {
		h[_i].config = cfg
	}
}
