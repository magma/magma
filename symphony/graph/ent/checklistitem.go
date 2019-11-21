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

// CheckListItem is the model entity for the CheckListItem schema.
type CheckListItem struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// Checked holds the value of the "checked" field.
	Checked bool `json:"checked,omitempty"`
	// StringVal holds the value of the "string_val" field.
	StringVal string `json:"string_val,omitempty" gqlgen:"stringValue"`
	// EnumValues holds the value of the "enum_values" field.
	EnumValues string `json:"enum_values,omitempty" gqlgen:"enumValues"`
	// HelpText holds the value of the "help_text" field.
	HelpText *string `json:"help_text,omitempty" gqlgen:"helpText"`
}

// FromRows scans the sql response data into CheckListItem.
func (cli *CheckListItem) FromRows(rows *sql.Rows) error {
	var scancli struct {
		ID         int
		Title      sql.NullString
		Type       sql.NullString
		Index      sql.NullInt64
		Checked    sql.NullBool
		StringVal  sql.NullString
		EnumValues sql.NullString
		HelpText   sql.NullString
	}
	// the order here should be the same as in the `checklistitem.Columns`.
	if err := rows.Scan(
		&scancli.ID,
		&scancli.Title,
		&scancli.Type,
		&scancli.Index,
		&scancli.Checked,
		&scancli.StringVal,
		&scancli.EnumValues,
		&scancli.HelpText,
	); err != nil {
		return err
	}
	cli.ID = strconv.Itoa(scancli.ID)
	cli.Title = scancli.Title.String
	cli.Type = scancli.Type.String
	cli.Index = int(scancli.Index.Int64)
	cli.Checked = scancli.Checked.Bool
	cli.StringVal = scancli.StringVal.String
	cli.EnumValues = scancli.EnumValues.String
	if scancli.HelpText.Valid {
		cli.HelpText = new(string)
		*cli.HelpText = scancli.HelpText.String
	}
	return nil
}

// QueryWorkOrder queries the work_order edge of the CheckListItem.
func (cli *CheckListItem) QueryWorkOrder() *WorkOrderQuery {
	return (&CheckListItemClient{cli.config}).QueryWorkOrder(cli)
}

// Update returns a builder for updating this CheckListItem.
// Note that, you need to call CheckListItem.Unwrap() before calling this method, if this CheckListItem
// was returned from a transaction, and the transaction was committed or rolled back.
func (cli *CheckListItem) Update() *CheckListItemUpdateOne {
	return (&CheckListItemClient{cli.config}).UpdateOne(cli)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (cli *CheckListItem) Unwrap() *CheckListItem {
	tx, ok := cli.config.driver.(*txDriver)
	if !ok {
		panic("ent: CheckListItem is not a transactional entity")
	}
	cli.config.driver = tx.drv
	return cli
}

// String implements the fmt.Stringer.
func (cli *CheckListItem) String() string {
	var builder strings.Builder
	builder.WriteString("CheckListItem(")
	builder.WriteString(fmt.Sprintf("id=%v", cli.ID))
	builder.WriteString(", title=")
	builder.WriteString(cli.Title)
	builder.WriteString(", type=")
	builder.WriteString(cli.Type)
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", cli.Index))
	builder.WriteString(", checked=")
	builder.WriteString(fmt.Sprintf("%v", cli.Checked))
	builder.WriteString(", string_val=")
	builder.WriteString(cli.StringVal)
	builder.WriteString(", enum_values=")
	builder.WriteString(cli.EnumValues)
	if v := cli.HelpText; v != nil {
		builder.WriteString(", help_text=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (cli *CheckListItem) id() int {
	id, _ := strconv.Atoi(cli.ID)
	return id
}

// CheckListItems is a parsable slice of CheckListItem.
type CheckListItems []*CheckListItem

// FromRows scans the sql response data into CheckListItems.
func (cli *CheckListItems) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scancli := &CheckListItem{}
		if err := scancli.FromRows(rows); err != nil {
			return err
		}
		*cli = append(*cli, scancli)
	}
	return nil
}

func (cli CheckListItems) config(cfg config) {
	for _i := range cli {
		cli[_i].config = cfg
	}
}
