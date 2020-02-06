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
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
)

// FloorPlanScale is the model entity for the FloorPlanScale schema.
type FloorPlanScale struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// ReferencePoint1X holds the value of the "reference_point1_x" field.
	ReferencePoint1X int `json:"reference_point1_x,omitempty"`
	// ReferencePoint1Y holds the value of the "reference_point1_y" field.
	ReferencePoint1Y int `json:"reference_point1_y,omitempty"`
	// ReferencePoint2X holds the value of the "reference_point2_x" field.
	ReferencePoint2X int `json:"reference_point2_x,omitempty"`
	// ReferencePoint2Y holds the value of the "reference_point2_y" field.
	ReferencePoint2Y int `json:"reference_point2_y,omitempty"`
	// ScaleInMeters holds the value of the "scale_in_meters" field.
	ScaleInMeters float64 `json:"scale_in_meters,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*FloorPlanScale) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},   // id
		&sql.NullTime{},    // create_time
		&sql.NullTime{},    // update_time
		&sql.NullInt64{},   // reference_point1_x
		&sql.NullInt64{},   // reference_point1_y
		&sql.NullInt64{},   // reference_point2_x
		&sql.NullInt64{},   // reference_point2_y
		&sql.NullFloat64{}, // scale_in_meters
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the FloorPlanScale fields.
func (fps *FloorPlanScale) assignValues(values ...interface{}) error {
	if m, n := len(values), len(floorplanscale.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	fps.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		fps.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		fps.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field reference_point1_x", values[2])
	} else if value.Valid {
		fps.ReferencePoint1X = int(value.Int64)
	}
	if value, ok := values[3].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field reference_point1_y", values[3])
	} else if value.Valid {
		fps.ReferencePoint1Y = int(value.Int64)
	}
	if value, ok := values[4].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field reference_point2_x", values[4])
	} else if value.Valid {
		fps.ReferencePoint2X = int(value.Int64)
	}
	if value, ok := values[5].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field reference_point2_y", values[5])
	} else if value.Valid {
		fps.ReferencePoint2Y = int(value.Int64)
	}
	if value, ok := values[6].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field scale_in_meters", values[6])
	} else if value.Valid {
		fps.ScaleInMeters = value.Float64
	}
	return nil
}

// Update returns a builder for updating this FloorPlanScale.
// Note that, you need to call FloorPlanScale.Unwrap() before calling this method, if this FloorPlanScale
// was returned from a transaction, and the transaction was committed or rolled back.
func (fps *FloorPlanScale) Update() *FloorPlanScaleUpdateOne {
	return (&FloorPlanScaleClient{fps.config}).UpdateOne(fps)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (fps *FloorPlanScale) Unwrap() *FloorPlanScale {
	tx, ok := fps.config.driver.(*txDriver)
	if !ok {
		panic("ent: FloorPlanScale is not a transactional entity")
	}
	fps.config.driver = tx.drv
	return fps
}

// String implements the fmt.Stringer.
func (fps *FloorPlanScale) String() string {
	var builder strings.Builder
	builder.WriteString("FloorPlanScale(")
	builder.WriteString(fmt.Sprintf("id=%v", fps.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(fps.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(fps.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", reference_point1_x=")
	builder.WriteString(fmt.Sprintf("%v", fps.ReferencePoint1X))
	builder.WriteString(", reference_point1_y=")
	builder.WriteString(fmt.Sprintf("%v", fps.ReferencePoint1Y))
	builder.WriteString(", reference_point2_x=")
	builder.WriteString(fmt.Sprintf("%v", fps.ReferencePoint2X))
	builder.WriteString(", reference_point2_y=")
	builder.WriteString(fmt.Sprintf("%v", fps.ReferencePoint2Y))
	builder.WriteString(", scale_in_meters=")
	builder.WriteString(fmt.Sprintf("%v", fps.ScaleInMeters))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (fps *FloorPlanScale) id() int {
	id, _ := strconv.Atoi(fps.ID)
	return id
}

// FloorPlanScales is a parsable slice of FloorPlanScale.
type FloorPlanScales []*FloorPlanScale

func (fps FloorPlanScales) config(cfg config) {
	for _i := range fps {
		fps[_i].config = cfg
	}
}
