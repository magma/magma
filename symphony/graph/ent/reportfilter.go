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
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
)

// ReportFilter is the model entity for the ReportFilter schema.
type ReportFilter struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Entity holds the value of the "entity" field.
	Entity reportfilter.Entity `json:"entity,omitempty"`
	// Filters holds the value of the "filters" field.
	Filters string `json:"filters,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ReportFilter) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // entity
		&sql.NullString{}, // filters
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ReportFilter fields.
func (rf *ReportFilter) assignValues(values ...interface{}) error {
	if m, n := len(values), len(reportfilter.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	rf.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		rf.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		rf.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		rf.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field entity", values[3])
	} else if value.Valid {
		rf.Entity = reportfilter.Entity(value.String)
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field filters", values[4])
	} else if value.Valid {
		rf.Filters = value.String
	}
	return nil
}

// Update returns a builder for updating this ReportFilter.
// Note that, you need to call ReportFilter.Unwrap() before calling this method, if this ReportFilter
// was returned from a transaction, and the transaction was committed or rolled back.
func (rf *ReportFilter) Update() *ReportFilterUpdateOne {
	return (&ReportFilterClient{config: rf.config}).UpdateOne(rf)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (rf *ReportFilter) Unwrap() *ReportFilter {
	tx, ok := rf.config.driver.(*txDriver)
	if !ok {
		panic("ent: ReportFilter is not a transactional entity")
	}
	rf.config.driver = tx.drv
	return rf
}

// String implements the fmt.Stringer.
func (rf *ReportFilter) String() string {
	var builder strings.Builder
	builder.WriteString("ReportFilter(")
	builder.WriteString(fmt.Sprintf("id=%v", rf.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(rf.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(rf.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(rf.Name)
	builder.WriteString(", entity=")
	builder.WriteString(fmt.Sprintf("%v", rf.Entity))
	builder.WriteString(", filters=")
	builder.WriteString(rf.Filters)
	builder.WriteByte(')')
	return builder.String()
}

// ReportFilters is a parsable slice of ReportFilter.
type ReportFilters []*ReportFilter

func (rf ReportFilters) config(cfg config) {
	for _i := range rf {
		rf[_i].config = cfg
	}
}
