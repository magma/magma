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

// FromRows scans the sql response data into FloorPlanScale.
func (fps *FloorPlanScale) FromRows(rows *sql.Rows) error {
	var scanfps struct {
		ID               int
		CreateTime       sql.NullTime
		UpdateTime       sql.NullTime
		ReferencePoint1X sql.NullInt64
		ReferencePoint1Y sql.NullInt64
		ReferencePoint2X sql.NullInt64
		ReferencePoint2Y sql.NullInt64
		ScaleInMeters    sql.NullFloat64
	}
	// the order here should be the same as in the `floorplanscale.Columns`.
	if err := rows.Scan(
		&scanfps.ID,
		&scanfps.CreateTime,
		&scanfps.UpdateTime,
		&scanfps.ReferencePoint1X,
		&scanfps.ReferencePoint1Y,
		&scanfps.ReferencePoint2X,
		&scanfps.ReferencePoint2Y,
		&scanfps.ScaleInMeters,
	); err != nil {
		return err
	}
	fps.ID = strconv.Itoa(scanfps.ID)
	fps.CreateTime = scanfps.CreateTime.Time
	fps.UpdateTime = scanfps.UpdateTime.Time
	fps.ReferencePoint1X = int(scanfps.ReferencePoint1X.Int64)
	fps.ReferencePoint1Y = int(scanfps.ReferencePoint1Y.Int64)
	fps.ReferencePoint2X = int(scanfps.ReferencePoint2X.Int64)
	fps.ReferencePoint2Y = int(scanfps.ReferencePoint2Y.Int64)
	fps.ScaleInMeters = scanfps.ScaleInMeters.Float64
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

// FromRows scans the sql response data into FloorPlanScales.
func (fps *FloorPlanScales) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanfps := &FloorPlanScale{}
		if err := scanfps.FromRows(rows); err != nil {
			return err
		}
		*fps = append(*fps, scanfps)
	}
	return nil
}

func (fps FloorPlanScales) config(cfg config) {
	for _i := range fps {
		fps[_i].config = cfg
	}
}
