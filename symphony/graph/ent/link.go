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
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// Link is the model entity for the Link schema.
type Link struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// FutureState holds the value of the "future_state" field.
	FutureState string `json:"future_state,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LinkQuery when eager-loading is set.
	Edges           LinkEdges `json:"edges"`
	link_work_order *int
}

// LinkEdges holds the relations/edges for other nodes in the graph.
type LinkEdges struct {
	// Ports holds the value of the ports edge.
	Ports []*EquipmentPort
	// WorkOrder holds the value of the work_order edge.
	WorkOrder *WorkOrder
	// Properties holds the value of the properties edge.
	Properties []*Property
	// Service holds the value of the service edge.
	Service []*Service
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// PortsOrErr returns the Ports value or an error if the edge
// was not loaded in eager-loading.
func (e LinkEdges) PortsOrErr() ([]*EquipmentPort, error) {
	if e.loadedTypes[0] {
		return e.Ports, nil
	}
	return nil, &NotLoadedError{edge: "ports"}
}

// WorkOrderOrErr returns the WorkOrder value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LinkEdges) WorkOrderOrErr() (*WorkOrder, error) {
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

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e LinkEdges) PropertiesOrErr() ([]*Property, error) {
	if e.loadedTypes[2] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// ServiceOrErr returns the Service value or an error if the edge
// was not loaded in eager-loading.
func (e LinkEdges) ServiceOrErr() ([]*Service, error) {
	if e.loadedTypes[3] {
		return e.Service, nil
	}
	return nil, &NotLoadedError{edge: "service"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Link) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // future_state
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Link) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // link_work_order
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Link fields.
func (l *Link) assignValues(values ...interface{}) error {
	if m, n := len(values), len(link.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	l.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		l.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		l.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field future_state", values[2])
	} else if value.Valid {
		l.FutureState = value.String
	}
	values = values[3:]
	if len(values) == len(link.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field link_work_order", value)
		} else if value.Valid {
			l.link_work_order = new(int)
			*l.link_work_order = int(value.Int64)
		}
	}
	return nil
}

// QueryPorts queries the ports edge of the Link.
func (l *Link) QueryPorts() *EquipmentPortQuery {
	return (&LinkClient{l.config}).QueryPorts(l)
}

// QueryWorkOrder queries the work_order edge of the Link.
func (l *Link) QueryWorkOrder() *WorkOrderQuery {
	return (&LinkClient{l.config}).QueryWorkOrder(l)
}

// QueryProperties queries the properties edge of the Link.
func (l *Link) QueryProperties() *PropertyQuery {
	return (&LinkClient{l.config}).QueryProperties(l)
}

// QueryService queries the service edge of the Link.
func (l *Link) QueryService() *ServiceQuery {
	return (&LinkClient{l.config}).QueryService(l)
}

// Update returns a builder for updating this Link.
// Note that, you need to call Link.Unwrap() before calling this method, if this Link
// was returned from a transaction, and the transaction was committed or rolled back.
func (l *Link) Update() *LinkUpdateOne {
	return (&LinkClient{l.config}).UpdateOne(l)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (l *Link) Unwrap() *Link {
	tx, ok := l.config.driver.(*txDriver)
	if !ok {
		panic("ent: Link is not a transactional entity")
	}
	l.config.driver = tx.drv
	return l
}

// String implements the fmt.Stringer.
func (l *Link) String() string {
	var builder strings.Builder
	builder.WriteString("Link(")
	builder.WriteString(fmt.Sprintf("id=%v", l.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(l.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(l.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", future_state=")
	builder.WriteString(l.FutureState)
	builder.WriteByte(')')
	return builder.String()
}

// Links is a parsable slice of Link.
type Links []*Link

func (l Links) config(cfg config) {
	for _i := range l {
		l[_i].config = cfg
	}
}
