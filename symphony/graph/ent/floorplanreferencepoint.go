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

// FloorPlanReferencePoint is the model entity for the FloorPlanReferencePoint schema.
type FloorPlanReferencePoint struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
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

// FromRows scans the sql response data into FloorPlanReferencePoint.
func (fprp *FloorPlanReferencePoint) FromRows(rows *sql.Rows) error {
	var scanfprp struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		X          sql.NullInt64
		Y          sql.NullInt64
		Latitude   sql.NullFloat64
		Longitude  sql.NullFloat64
	}
	// the order here should be the same as in the `floorplanreferencepoint.Columns`.
	if err := rows.Scan(
		&scanfprp.ID,
		&scanfprp.CreateTime,
		&scanfprp.UpdateTime,
		&scanfprp.X,
		&scanfprp.Y,
		&scanfprp.Latitude,
		&scanfprp.Longitude,
	); err != nil {
		return err
	}
	fprp.ID = strconv.Itoa(scanfprp.ID)
	fprp.CreateTime = scanfprp.CreateTime.Time
	fprp.UpdateTime = scanfprp.UpdateTime.Time
	fprp.X = int(scanfprp.X.Int64)
	fprp.Y = int(scanfprp.Y.Int64)
	fprp.Latitude = scanfprp.Latitude.Float64
	fprp.Longitude = scanfprp.Longitude.Float64
	return nil
}

// Update returns a builder for updating this FloorPlanReferencePoint.
// Note that, you need to call FloorPlanReferencePoint.Unwrap() before calling this method, if this FloorPlanReferencePoint
// was returned from a transaction, and the transaction was committed or rolled back.
func (fprp *FloorPlanReferencePoint) Update() *FloorPlanReferencePointUpdateOne {
	return (&FloorPlanReferencePointClient{fprp.config}).UpdateOne(fprp)
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

// id returns the int representation of the ID field.
func (fprp *FloorPlanReferencePoint) id() int {
	id, _ := strconv.Atoi(fprp.ID)
	return id
}

// FloorPlanReferencePoints is a parsable slice of FloorPlanReferencePoint.
type FloorPlanReferencePoints []*FloorPlanReferencePoint

// FromRows scans the sql response data into FloorPlanReferencePoints.
func (fprp *FloorPlanReferencePoints) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanfprp := &FloorPlanReferencePoint{}
		if err := scanfprp.FromRows(rows); err != nil {
			return err
		}
		*fprp = append(*fprp, scanfprp)
	}
	return nil
}

func (fprp FloorPlanReferencePoints) config(cfg config) {
	for _i := range fprp {
		fprp[_i].config = cfg
	}
}
