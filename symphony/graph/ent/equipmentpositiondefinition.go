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

// EquipmentPositionDefinition is the model entity for the EquipmentPositionDefinition schema.
type EquipmentPositionDefinition struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// VisibilityLabel holds the value of the "visibility_label" field.
	VisibilityLabel string `json:"visibility_label,omitempty" gqlgen:"visibleLabel"`
}

// FromRows scans the sql response data into EquipmentPositionDefinition.
func (epd *EquipmentPositionDefinition) FromRows(rows *sql.Rows) error {
	var scanepd struct {
		ID              int
		CreateTime      sql.NullTime
		UpdateTime      sql.NullTime
		Name            sql.NullString
		Index           sql.NullInt64
		VisibilityLabel sql.NullString
	}
	// the order here should be the same as in the `equipmentpositiondefinition.Columns`.
	if err := rows.Scan(
		&scanepd.ID,
		&scanepd.CreateTime,
		&scanepd.UpdateTime,
		&scanepd.Name,
		&scanepd.Index,
		&scanepd.VisibilityLabel,
	); err != nil {
		return err
	}
	epd.ID = strconv.Itoa(scanepd.ID)
	epd.CreateTime = scanepd.CreateTime.Time
	epd.UpdateTime = scanepd.UpdateTime.Time
	epd.Name = scanepd.Name.String
	epd.Index = int(scanepd.Index.Int64)
	epd.VisibilityLabel = scanepd.VisibilityLabel.String
	return nil
}

// QueryPositions queries the positions edge of the EquipmentPositionDefinition.
func (epd *EquipmentPositionDefinition) QueryPositions() *EquipmentPositionQuery {
	return (&EquipmentPositionDefinitionClient{epd.config}).QueryPositions(epd)
}

// QueryEquipmentType queries the equipment_type edge of the EquipmentPositionDefinition.
func (epd *EquipmentPositionDefinition) QueryEquipmentType() *EquipmentTypeQuery {
	return (&EquipmentPositionDefinitionClient{epd.config}).QueryEquipmentType(epd)
}

// Update returns a builder for updating this EquipmentPositionDefinition.
// Note that, you need to call EquipmentPositionDefinition.Unwrap() before calling this method, if this EquipmentPositionDefinition
// was returned from a transaction, and the transaction was committed or rolled back.
func (epd *EquipmentPositionDefinition) Update() *EquipmentPositionDefinitionUpdateOne {
	return (&EquipmentPositionDefinitionClient{epd.config}).UpdateOne(epd)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (epd *EquipmentPositionDefinition) Unwrap() *EquipmentPositionDefinition {
	tx, ok := epd.config.driver.(*txDriver)
	if !ok {
		panic("ent: EquipmentPositionDefinition is not a transactional entity")
	}
	epd.config.driver = tx.drv
	return epd
}

// String implements the fmt.Stringer.
func (epd *EquipmentPositionDefinition) String() string {
	var builder strings.Builder
	builder.WriteString("EquipmentPositionDefinition(")
	builder.WriteString(fmt.Sprintf("id=%v", epd.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(epd.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(epd.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(epd.Name)
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", epd.Index))
	builder.WriteString(", visibility_label=")
	builder.WriteString(epd.VisibilityLabel)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (epd *EquipmentPositionDefinition) id() int {
	id, _ := strconv.Atoi(epd.ID)
	return id
}

// EquipmentPositionDefinitions is a parsable slice of EquipmentPositionDefinition.
type EquipmentPositionDefinitions []*EquipmentPositionDefinition

// FromRows scans the sql response data into EquipmentPositionDefinitions.
func (epd *EquipmentPositionDefinitions) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanepd := &EquipmentPositionDefinition{}
		if err := scanepd.FromRows(rows); err != nil {
			return err
		}
		*epd = append(*epd, scanepd)
	}
	return nil
}

func (epd EquipmentPositionDefinitions) config(cfg config) {
	for _i := range epd {
		epd[_i].config = cfg
	}
}
