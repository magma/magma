// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/facebookincubator/ent/dialect/sql"
)

// CheckListItemDefinition is the model entity for the CheckListItemDefinition schema.
type CheckListItemDefinition struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// EnumValues holds the value of the "enum_values" field.
	EnumValues *string `json:"enum_values,omitempty" gqlgen:"enumValues"`
	// HelpText holds the value of the "help_text" field.
	HelpText *string `json:"help_text,omitempty" gqlgen:"helpText"`
}

// FromRows scans the sql response data into CheckListItemDefinition.
func (clid *CheckListItemDefinition) FromRows(rows *sql.Rows) error {
	var scanclid struct {
		ID         int
		Title      sql.NullString
		Type       sql.NullString
		Index      sql.NullInt64
		EnumValues sql.NullString
		HelpText   sql.NullString
	}
	// the order here should be the same as in the `checklistitemdefinition.Columns`.
	if err := rows.Scan(
		&scanclid.ID,
		&scanclid.Title,
		&scanclid.Type,
		&scanclid.Index,
		&scanclid.EnumValues,
		&scanclid.HelpText,
	); err != nil {
		return err
	}
	clid.ID = strconv.Itoa(scanclid.ID)
	clid.Title = scanclid.Title.String
	clid.Type = scanclid.Type.String
	clid.Index = int(scanclid.Index.Int64)
	if scanclid.EnumValues.Valid {
		clid.EnumValues = new(string)
		*clid.EnumValues = scanclid.EnumValues.String
	}
	if scanclid.HelpText.Valid {
		clid.HelpText = new(string)
		*clid.HelpText = scanclid.HelpText.String
	}
	return nil
}

// QueryWorkOrderType queries the work_order_type edge of the CheckListItemDefinition.
func (clid *CheckListItemDefinition) QueryWorkOrderType() *WorkOrderTypeQuery {
	return (&CheckListItemDefinitionClient{clid.config}).QueryWorkOrderType(clid)
}

// Update returns a builder for updating this CheckListItemDefinition.
// Note that, you need to call CheckListItemDefinition.Unwrap() before calling this method, if this CheckListItemDefinition
// was returned from a transaction, and the transaction was committed or rolled back.
func (clid *CheckListItemDefinition) Update() *CheckListItemDefinitionUpdateOne {
	return (&CheckListItemDefinitionClient{clid.config}).UpdateOne(clid)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (clid *CheckListItemDefinition) Unwrap() *CheckListItemDefinition {
	tx, ok := clid.config.driver.(*txDriver)
	if !ok {
		panic("ent: CheckListItemDefinition is not a transactional entity")
	}
	clid.config.driver = tx.drv
	return clid
}

// String implements the fmt.Stringer.
func (clid *CheckListItemDefinition) String() string {
	var builder strings.Builder
	builder.WriteString("CheckListItemDefinition(")
	builder.WriteString(fmt.Sprintf("id=%v", clid.ID))
	builder.WriteString(", title=")
	builder.WriteString(clid.Title)
	builder.WriteString(", type=")
	builder.WriteString(clid.Type)
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", clid.Index))
	if v := clid.EnumValues; v != nil {
		builder.WriteString(", enum_values=")
		builder.WriteString(*v)
	}
	if v := clid.HelpText; v != nil {
		builder.WriteString(", help_text=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (clid *CheckListItemDefinition) id() int {
	id, _ := strconv.Atoi(clid.ID)
	return id
}

// CheckListItemDefinitions is a parsable slice of CheckListItemDefinition.
type CheckListItemDefinitions []*CheckListItemDefinition

// FromRows scans the sql response data into CheckListItemDefinitions.
func (clid *CheckListItemDefinitions) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanclid := &CheckListItemDefinition{}
		if err := scanclid.FromRows(rows); err != nil {
			return err
		}
		*clid = append(*clid, scanclid)
	}
	return nil
}

func (clid CheckListItemDefinitions) config(cfg config) {
	for _i := range clid {
		clid[_i].config = cfg
	}
}
