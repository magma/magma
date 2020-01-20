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
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
)

// EquipmentPosition is the model entity for the EquipmentPosition schema.
type EquipmentPosition struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*EquipmentPosition) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},
		&sql.NullTime{},
		&sql.NullTime{},
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the EquipmentPosition fields.
func (ep *EquipmentPosition) assignValues(values ...interface{}) error {
	if m, n := len(values), len(equipmentposition.Columns); m != n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	ep.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		ep.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		ep.UpdateTime = value.Time
	}
	return nil
}

// QueryDefinition queries the definition edge of the EquipmentPosition.
func (ep *EquipmentPosition) QueryDefinition() *EquipmentPositionDefinitionQuery {
	return (&EquipmentPositionClient{ep.config}).QueryDefinition(ep)
}

// QueryParent queries the parent edge of the EquipmentPosition.
func (ep *EquipmentPosition) QueryParent() *EquipmentQuery {
	return (&EquipmentPositionClient{ep.config}).QueryParent(ep)
}

// QueryAttachment queries the attachment edge of the EquipmentPosition.
func (ep *EquipmentPosition) QueryAttachment() *EquipmentQuery {
	return (&EquipmentPositionClient{ep.config}).QueryAttachment(ep)
}

// Update returns a builder for updating this EquipmentPosition.
// Note that, you need to call EquipmentPosition.Unwrap() before calling this method, if this EquipmentPosition
// was returned from a transaction, and the transaction was committed or rolled back.
func (ep *EquipmentPosition) Update() *EquipmentPositionUpdateOne {
	return (&EquipmentPositionClient{ep.config}).UpdateOne(ep)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (ep *EquipmentPosition) Unwrap() *EquipmentPosition {
	tx, ok := ep.config.driver.(*txDriver)
	if !ok {
		panic("ent: EquipmentPosition is not a transactional entity")
	}
	ep.config.driver = tx.drv
	return ep
}

// String implements the fmt.Stringer.
func (ep *EquipmentPosition) String() string {
	var builder strings.Builder
	builder.WriteString("EquipmentPosition(")
	builder.WriteString(fmt.Sprintf("id=%v", ep.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(ep.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(ep.UpdateTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (ep *EquipmentPosition) id() int {
	id, _ := strconv.Atoi(ep.ID)
	return id
}

// EquipmentPositions is a parsable slice of EquipmentPosition.
type EquipmentPositions []*EquipmentPosition

func (ep EquipmentPositions) config(cfg config) {
	for _i := range ep {
		ep[_i].config = cfg
	}
}
