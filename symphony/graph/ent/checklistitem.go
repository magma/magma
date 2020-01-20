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
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
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

// scanValues returns the types for scanning values from sql.Rows.
func (*CheckListItem) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},
		&sql.NullString{},
		&sql.NullString{},
		&sql.NullInt64{},
		&sql.NullBool{},
		&sql.NullString{},
		&sql.NullString{},
		&sql.NullString{},
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CheckListItem fields.
func (cli *CheckListItem) assignValues(values ...interface{}) error {
	if m, n := len(values), len(checklistitem.Columns); m != n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	cli.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field title", values[0])
	} else if value.Valid {
		cli.Title = value.String
	}
	if value, ok := values[1].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field type", values[1])
	} else if value.Valid {
		cli.Type = value.String
	}
	if value, ok := values[2].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field index", values[2])
	} else if value.Valid {
		cli.Index = int(value.Int64)
	}
	if value, ok := values[3].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field checked", values[3])
	} else if value.Valid {
		cli.Checked = value.Bool
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field string_val", values[4])
	} else if value.Valid {
		cli.StringVal = value.String
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field enum_values", values[5])
	} else if value.Valid {
		cli.EnumValues = value.String
	}
	if value, ok := values[6].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field help_text", values[6])
	} else if value.Valid {
		cli.HelpText = new(string)
		*cli.HelpText = value.String
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

func (cli CheckListItems) config(cfg config) {
	for _i := range cli {
		cli[_i].config = cfg
	}
}
