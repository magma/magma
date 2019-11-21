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
}

// FromRows scans the sql response data into EquipmentCategory.
func (ec *EquipmentCategory) FromRows(rows *sql.Rows) error {
	var scanec struct {
		ID         int
		CreateTime sql.NullTime
		UpdateTime sql.NullTime
		Name       sql.NullString
	}
	// the order here should be the same as in the `equipmentcategory.Columns`.
	if err := rows.Scan(
		&scanec.ID,
		&scanec.CreateTime,
		&scanec.UpdateTime,
		&scanec.Name,
	); err != nil {
		return err
	}
	ec.ID = strconv.Itoa(scanec.ID)
	ec.CreateTime = scanec.CreateTime.Time
	ec.UpdateTime = scanec.UpdateTime.Time
	ec.Name = scanec.Name.String
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

// FromRows scans the sql response data into EquipmentCategories.
func (ec *EquipmentCategories) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanec := &EquipmentCategory{}
		if err := scanec.FromRows(rows); err != nil {
			return err
		}
		*ec = append(*ec, scanec)
	}
	return nil
}

func (ec EquipmentCategories) config(cfg config) {
	for _i := range ec {
		ec[_i].config = cfg
	}
}
