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

// FloorPlan is the model entity for the FloorPlan schema.
type FloorPlan struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
}

// FromRows scans the sql response data into FloorPlan.
func (fp *FloorPlan) FromRows(rows *sql.Rows) error {
	var scanfp struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		Name       sql.NullString
	}
	// the order here should be the same as in the `floorplan.Columns`.
	if err := rows.Scan(
		&scanfp.ID,
		&scanfp.CreateTime,
		&scanfp.UpdateTime,
		&scanfp.Name,
	); err != nil {
		return err
	}
	fp.ID = strconv.Itoa(scanfp.ID)
	fp.CreateTime = scanfp.CreateTime.Time
	fp.UpdateTime = scanfp.UpdateTime.Time
	fp.Name = scanfp.Name.String
	return nil
}

// QueryLocation queries the location edge of the FloorPlan.
func (fp *FloorPlan) QueryLocation() *LocationQuery {
	return (&FloorPlanClient{fp.config}).QueryLocation(fp)
}

// QueryReferencePoint queries the reference_point edge of the FloorPlan.
func (fp *FloorPlan) QueryReferencePoint() *FloorPlanReferencePointQuery {
	return (&FloorPlanClient{fp.config}).QueryReferencePoint(fp)
}

// QueryScale queries the scale edge of the FloorPlan.
func (fp *FloorPlan) QueryScale() *FloorPlanScaleQuery {
	return (&FloorPlanClient{fp.config}).QueryScale(fp)
}

// QueryImage queries the image edge of the FloorPlan.
func (fp *FloorPlan) QueryImage() *FileQuery {
	return (&FloorPlanClient{fp.config}).QueryImage(fp)
}

// Update returns a builder for updating this FloorPlan.
// Note that, you need to call FloorPlan.Unwrap() before calling this method, if this FloorPlan
// was returned from a transaction, and the transaction was committed or rolled back.
func (fp *FloorPlan) Update() *FloorPlanUpdateOne {
	return (&FloorPlanClient{fp.config}).UpdateOne(fp)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (fp *FloorPlan) Unwrap() *FloorPlan {
	tx, ok := fp.config.driver.(*txDriver)
	if !ok {
		panic("ent: FloorPlan is not a transactional entity")
	}
	fp.config.driver = tx.drv
	return fp
}

// String implements the fmt.Stringer.
func (fp *FloorPlan) String() string {
	var builder strings.Builder
	builder.WriteString("FloorPlan(")
	builder.WriteString(fmt.Sprintf("id=%v", fp.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(fp.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(fp.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(fp.Name)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (fp *FloorPlan) id() int {
	id, _ := strconv.Atoi(fp.ID)
	return id
}

// FloorPlans is a parsable slice of FloorPlan.
type FloorPlans []*FloorPlan

// FromRows scans the sql response data into FloorPlans.
func (fp *FloorPlans) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanfp := &FloorPlan{}
		if err := scanfp.FromRows(rows); err != nil {
			return err
		}
		*fp = append(*fp, scanfp)
	}
	return nil
}

func (fp FloorPlans) config(cfg config) {
	for _i := range fp {
		fp[_i].config = cfg
	}
}
