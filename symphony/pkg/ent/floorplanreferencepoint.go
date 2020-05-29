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
	"github.com/facebookincubator/symphony/pkg/ent/floorplanreferencepoint"
)

// FloorPlanReferencePoint is the model entity for the FloorPlanReferencePoint schema.
type FloorPlanReferencePoint struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// X holds the value of the "x" field.
	X int `json:"x,omitempty"`
	// Y holds the value of the "y" field.
	Y int `json:"y,omitempty"`
	// Latitude holds the value of the "latitude" field.
	Latitude float64 `json:"latitude,omitempty"`
	// Longitude holds the value of the "longitude" field.
	Longitude float64 `json:"longitude,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*FloorPlanReferencePoint) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},   // id
		&sql.NullTime{},    // create_time
		&sql.NullTime{},    // update_time
		&sql.NullInt64{},   // x
		&sql.NullInt64{},   // y
		&sql.NullFloat64{}, // latitude
		&sql.NullFloat64{}, // longitude
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the FloorPlanReferencePoint fields.
func (fprp *FloorPlanReferencePoint) assignValues(values ...interface{}) error {
	if m, n := len(values), len(floorplanreferencepoint.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	fprp.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		fprp.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		fprp.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field x", values[2])
	} else if value.Valid {
		fprp.X = int(value.Int64)
	}
	if value, ok := values[3].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field y", values[3])
	} else if value.Valid {
		fprp.Y = int(value.Int64)
	}
	if value, ok := values[4].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field latitude", values[4])
	} else if value.Valid {
		fprp.Latitude = value.Float64
	}
	if value, ok := values[5].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field longitude", values[5])
	} else if value.Valid {
		fprp.Longitude = value.Float64
	}
	return nil
}

// Update returns a builder for updating this FloorPlanReferencePoint.
// Note that, you need to call FloorPlanReferencePoint.Unwrap() before calling this method, if this FloorPlanReferencePoint
// was returned from a transaction, and the transaction was committed or rolled back.
func (fprp *FloorPlanReferencePoint) Update() *FloorPlanReferencePointUpdateOne {
	return (&FloorPlanReferencePointClient{config: fprp.config}).UpdateOne(fprp)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (fprp *FloorPlanReferencePoint) Unwrap() *FloorPlanReferencePoint {
	tx, ok := fprp.config.driver.(*txDriver)
	if !ok {
		panic("ent: FloorPlanReferencePoint is not a transactional entity")
	}
	fprp.config.driver = tx.drv
	return fprp
}

// String implements the fmt.Stringer.
func (fprp *FloorPlanReferencePoint) String() string {
	var builder strings.Builder
	builder.WriteString("FloorPlanReferencePoint(")
	builder.WriteString(fmt.Sprintf("id=%v", fprp.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(fprp.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(fprp.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", x=")
	builder.WriteString(fmt.Sprintf("%v", fprp.X))
	builder.WriteString(", y=")
	builder.WriteString(fmt.Sprintf("%v", fprp.Y))
	builder.WriteString(", latitude=")
	builder.WriteString(fmt.Sprintf("%v", fprp.Latitude))
	builder.WriteString(", longitude=")
	builder.WriteString(fmt.Sprintf("%v", fprp.Longitude))
	builder.WriteByte(')')
	return builder.String()
}

// FloorPlanReferencePoints is a parsable slice of FloorPlanReferencePoint.
type FloorPlanReferencePoints []*FloorPlanReferencePoint

func (fprp FloorPlanReferencePoints) config(cfg config) {
	for _i := range fprp {
		fprp[_i].config = cfg
	}
}
