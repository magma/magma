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
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
)

// EquipmentCategory is the model entity for the EquipmentCategory schema.
type EquipmentCategory struct {
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
	// The values are being populated by the EquipmentCategoryQuery when eager-loading is set.
	Edges EquipmentCategoryEdges `json:"edges"`
}

// EquipmentCategoryEdges holds the relations/edges for other nodes in the graph.
type EquipmentCategoryEdges struct {
	// Types holds the value of the types edge.
	Types []*EquipmentType
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// TypesOrErr returns the Types value or an error if the edge
// was not loaded in eager-loading.
func (e EquipmentCategoryEdges) TypesOrErr() ([]*EquipmentType, error) {
	if e.loadedTypes[0] {
		return e.Types, nil
	}
	return nil, &NotLoadedError{edge: "types"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*EquipmentCategory) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the EquipmentCategory fields.
func (ec *EquipmentCategory) assignValues(values ...interface{}) error {
	if m, n := len(values), len(equipmentcategory.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	ec.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		ec.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		ec.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		ec.Name = value.String
	}
	return nil
}

// QueryTypes queries the types edge of the EquipmentCategory.
func (ec *EquipmentCategory) QueryTypes() *EquipmentTypeQuery {
	return (&EquipmentCategoryClient{ec.config}).QueryTypes(ec)
}

// Update returns a builder for updating this EquipmentCategory.
// Note that, you need to call EquipmentCategory.Unwrap() before calling this method, if this EquipmentCategory
// was returned from a transaction, and the transaction was committed or rolled back.
func (ec *EquipmentCategory) Update() *EquipmentCategoryUpdateOne {
	return (&EquipmentCategoryClient{ec.config}).UpdateOne(ec)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (ec *EquipmentCategory) Unwrap() *EquipmentCategory {
	tx, ok := ec.config.driver.(*txDriver)
	if !ok {
		panic("ent: EquipmentCategory is not a transactional entity")
	}
	ec.config.driver = tx.drv
	return ec
}

// String implements the fmt.Stringer.
func (ec *EquipmentCategory) String() string {
	var builder strings.Builder
	builder.WriteString("EquipmentCategory(")
	builder.WriteString(fmt.Sprintf("id=%v", ec.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(ec.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(ec.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(ec.Name)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (ec *EquipmentCategory) id() int {
	id, _ := strconv.Atoi(ec.ID)
	return id
}

// EquipmentCategories is a parsable slice of EquipmentCategory.
type EquipmentCategories []*EquipmentCategory

func (ec EquipmentCategories) config(cfg config) {
	for _i := range ec {
		ec[_i].config = cfg
	}
}
