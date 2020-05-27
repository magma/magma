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
	"github.com/facebookincubator/symphony/graph/ent/activity"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// Activity is the model entity for the Activity schema.
type Activity struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// ChangedField holds the value of the "changed_field" field.
	ChangedField activity.ChangedField `json:"changed_field,omitempty"`
	// IsCreate holds the value of the "is_create" field.
	IsCreate bool `json:"is_create,omitempty"`
	// OldValue holds the value of the "old_value" field.
	OldValue string `json:"old_value,omitempty"`
	// NewValue holds the value of the "new_value" field.
	NewValue string `json:"new_value,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ActivityQuery when eager-loading is set.
	Edges                 ActivityEdges `json:"edges"`
	activity_author       *int
	work_order_activities *int
}

// ActivityEdges holds the relations/edges for other nodes in the graph.
type ActivityEdges struct {
	// Author holds the value of the author edge.
	Author *User
	// WorkOrder holds the value of the work_order edge.
	WorkOrder *WorkOrder
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// AuthorOrErr returns the Author value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ActivityEdges) AuthorOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.Author == nil {
			// The edge author was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.Author, nil
	}
	return nil, &NotLoadedError{edge: "author"}
}

// WorkOrderOrErr returns the WorkOrder value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ActivityEdges) WorkOrderOrErr() (*WorkOrder, error) {
	if e.loadedTypes[1] {
		if e.WorkOrder == nil {
			// The edge work_order was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: workorder.Label}
		}
		return e.WorkOrder, nil
	}
	return nil, &NotLoadedError{edge: "work_order"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Activity) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // changed_field
		&sql.NullBool{},   // is_create
		&sql.NullString{}, // old_value
		&sql.NullString{}, // new_value
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Activity) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // activity_author
		&sql.NullInt64{}, // work_order_activities
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Activity fields.
func (a *Activity) assignValues(values ...interface{}) error {
	if m, n := len(values), len(activity.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	a.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		a.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		a.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field changed_field", values[2])
	} else if value.Valid {
		a.ChangedField = activity.ChangedField(value.String)
	}
	if value, ok := values[3].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field is_create", values[3])
	} else if value.Valid {
		a.IsCreate = value.Bool
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field old_value", values[4])
	} else if value.Valid {
		a.OldValue = value.String
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field new_value", values[5])
	} else if value.Valid {
		a.NewValue = value.String
	}
	values = values[6:]
	if len(values) == len(activity.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field activity_author", value)
		} else if value.Valid {
			a.activity_author = new(int)
			*a.activity_author = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_activities", value)
		} else if value.Valid {
			a.work_order_activities = new(int)
			*a.work_order_activities = int(value.Int64)
		}
	}
	return nil
}

// QueryAuthor queries the author edge of the Activity.
func (a *Activity) QueryAuthor() *UserQuery {
	return (&ActivityClient{config: a.config}).QueryAuthor(a)
}

// QueryWorkOrder queries the work_order edge of the Activity.
func (a *Activity) QueryWorkOrder() *WorkOrderQuery {
	return (&ActivityClient{config: a.config}).QueryWorkOrder(a)
}

// Update returns a builder for updating this Activity.
// Note that, you need to call Activity.Unwrap() before calling this method, if this Activity
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Activity) Update() *ActivityUpdateOne {
	return (&ActivityClient{config: a.config}).UpdateOne(a)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (a *Activity) Unwrap() *Activity {
	tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Activity is not a transactional entity")
	}
	a.config.driver = tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Activity) String() string {
	var builder strings.Builder
	builder.WriteString("Activity(")
	builder.WriteString(fmt.Sprintf("id=%v", a.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(a.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(a.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", changed_field=")
	builder.WriteString(fmt.Sprintf("%v", a.ChangedField))
	builder.WriteString(", is_create=")
	builder.WriteString(fmt.Sprintf("%v", a.IsCreate))
	builder.WriteString(", old_value=")
	builder.WriteString(a.OldValue)
	builder.WriteString(", new_value=")
	builder.WriteString(a.NewValue)
	builder.WriteByte(')')
	return builder.String()
}

// Activities is a parsable slice of Activity.
type Activities []*Activity

func (a Activities) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}
