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
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// Hyperlink is the model entity for the Hyperlink schema.
type Hyperlink struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty" gqlgen:"displayName"`
	// Category holds the value of the "category" field.
	Category string `json:"category,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the HyperlinkQuery when eager-loading is set.
	Edges                 HyperlinkEdges `json:"edges"`
	equipment_hyperlinks  *int
	location_hyperlinks   *int
	work_order_hyperlinks *int
}

// HyperlinkEdges holds the relations/edges for other nodes in the graph.
type HyperlinkEdges struct {
	// Equipment holds the value of the equipment edge.
	Equipment *Equipment
	// Location holds the value of the location edge.
	Location *Location
	// WorkOrder holds the value of the work_order edge.
	WorkOrder *WorkOrder
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// EquipmentOrErr returns the Equipment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e HyperlinkEdges) EquipmentOrErr() (*Equipment, error) {
	if e.loadedTypes[0] {
		if e.Equipment == nil {
			// The edge equipment was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: equipment.Label}
		}
		return e.Equipment, nil
	}
	return nil, &NotLoadedError{edge: "equipment"}
}

// LocationOrErr returns the Location value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e HyperlinkEdges) LocationOrErr() (*Location, error) {
	if e.loadedTypes[1] {
		if e.Location == nil {
			// The edge location was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.Location, nil
	}
	return nil, &NotLoadedError{edge: "location"}
}

// WorkOrderOrErr returns the WorkOrder value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e HyperlinkEdges) WorkOrderOrErr() (*WorkOrder, error) {
	if e.loadedTypes[2] {
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
		&sql.NullInt64{}, // equipment_hyperlinks
		&sql.NullInt64{}, // location_hyperlinks
		&sql.NullInt64{}, // work_order_hyperlinks
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
	h.ID = int(value.Int64)
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
			return fmt.Errorf("unexpected type %T for edge-field equipment_hyperlinks", value)
		} else if value.Valid {
			h.equipment_hyperlinks = new(int)
			*h.equipment_hyperlinks = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_hyperlinks", value)
		} else if value.Valid {
			h.location_hyperlinks = new(int)
			*h.location_hyperlinks = int(value.Int64)
		}
		if value, ok := values[2].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_hyperlinks", value)
		} else if value.Valid {
			h.work_order_hyperlinks = new(int)
			*h.work_order_hyperlinks = int(value.Int64)
		}
	}
	return nil
}

// QueryEquipment queries the equipment edge of the Hyperlink.
func (h *Hyperlink) QueryEquipment() *EquipmentQuery {
	return (&HyperlinkClient{config: h.config}).QueryEquipment(h)
}

// QueryLocation queries the location edge of the Hyperlink.
func (h *Hyperlink) QueryLocation() *LocationQuery {
	return (&HyperlinkClient{config: h.config}).QueryLocation(h)
}

// QueryWorkOrder queries the work_order edge of the Hyperlink.
func (h *Hyperlink) QueryWorkOrder() *WorkOrderQuery {
	return (&HyperlinkClient{config: h.config}).QueryWorkOrder(h)
}

// Update returns a builder for updating this Hyperlink.
// Note that, you need to call Hyperlink.Unwrap() before calling this method, if this Hyperlink
// was returned from a transaction, and the transaction was committed or rolled back.
func (h *Hyperlink) Update() *HyperlinkUpdateOne {
	return (&HyperlinkClient{config: h.config}).UpdateOne(h)
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

// Hyperlinks is a parsable slice of Hyperlink.
type Hyperlinks []*Hyperlink

func (h Hyperlinks) config(cfg config) {
	for _i := range h {
		h[_i].config = cfg
	}
}
