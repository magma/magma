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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/location"
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
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the FloorPlanQuery when eager-loading is set.
	Edges                      FloorPlanEdges `json:"edges"`
	floor_plan_location        *string
	floor_plan_reference_point *string
	floor_plan_scale           *string
	floor_plan_image           *string
}

// FloorPlanEdges holds the relations/edges for other nodes in the graph.
type FloorPlanEdges struct {
	// Location holds the value of the location edge.
	Location *Location
	// ReferencePoint holds the value of the reference_point edge.
	ReferencePoint *FloorPlanReferencePoint
	// Scale holds the value of the scale edge.
	Scale *FloorPlanScale
	// Image holds the value of the image edge.
	Image *File
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// LocationOrErr returns the Location value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FloorPlanEdges) LocationOrErr() (*Location, error) {
	if e.loadedTypes[0] {
		if e.Location == nil {
			// The edge location was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.Location, nil
	}
	return nil, &NotLoadedError{edge: "location"}
}

// ReferencePointOrErr returns the ReferencePoint value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FloorPlanEdges) ReferencePointOrErr() (*FloorPlanReferencePoint, error) {
	if e.loadedTypes[1] {
		if e.ReferencePoint == nil {
			// The edge reference_point was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: floorplanreferencepoint.Label}
		}
		return e.ReferencePoint, nil
	}
	return nil, &NotLoadedError{edge: "reference_point"}
}

// ScaleOrErr returns the Scale value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FloorPlanEdges) ScaleOrErr() (*FloorPlanScale, error) {
	if e.loadedTypes[2] {
		if e.Scale == nil {
			// The edge scale was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: floorplanscale.Label}
		}
		return e.Scale, nil
	}
	return nil, &NotLoadedError{edge: "scale"}
}

// ImageOrErr returns the Image value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e FloorPlanEdges) ImageOrErr() (*File, error) {
	if e.loadedTypes[3] {
		if e.Image == nil {
			// The edge image was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: file.Label}
		}
		return e.Image, nil
	}
	return nil, &NotLoadedError{edge: "image"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*FloorPlan) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*FloorPlan) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // floor_plan_location
		&sql.NullInt64{}, // floor_plan_reference_point
		&sql.NullInt64{}, // floor_plan_scale
		&sql.NullInt64{}, // floor_plan_image
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the FloorPlan fields.
func (fp *FloorPlan) assignValues(values ...interface{}) error {
	if m, n := len(values), len(floorplan.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	fp.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		fp.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		fp.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		fp.Name = value.String
	}
	values = values[3:]
	if len(values) == len(floorplan.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field floor_plan_location", value)
		} else if value.Valid {
			fp.floor_plan_location = new(string)
			*fp.floor_plan_location = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field floor_plan_reference_point", value)
		} else if value.Valid {
			fp.floor_plan_reference_point = new(string)
			*fp.floor_plan_reference_point = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field floor_plan_scale", value)
		} else if value.Valid {
			fp.floor_plan_scale = new(string)
			*fp.floor_plan_scale = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[3].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field floor_plan_image", value)
		} else if value.Valid {
			fp.floor_plan_image = new(string)
			*fp.floor_plan_image = strconv.FormatInt(value.Int64, 10)
		}
	}
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

func (fp FloorPlans) config(cfg config) {
	for _i := range fp {
		fp[_i].config = cfg
	}
}
