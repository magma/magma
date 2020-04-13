// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
)

// CheckListItem is the model entity for the CheckListItem schema.
type CheckListItem struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
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
	// EnumSelectionMode holds the value of the "enum_selection_mode" field.
	EnumSelectionMode string `json:"enum_selection_mode,omitempty" gqlgen:"enumSelectionMode"`
	// SelectedEnumValues holds the value of the "selected_enum_values" field.
	SelectedEnumValues string `json:"selected_enum_values,omitempty" gqlgen:"selectedEnumValues"`
	// YesNoVal holds the value of the "yes_no_val" field.
	YesNoVal checklistitem.YesNoVal `json:"yes_no_val,omitempty"`
	// HelpText holds the value of the "help_text" field.
	HelpText *string `json:"help_text,omitempty" gqlgen:"helpText"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CheckListItemQuery when eager-loading is set.
	Edges                                CheckListItemEdges `json:"edges"`
	check_list_category_check_list_items *int
	work_order_check_list_items          *int
}

// CheckListItemEdges holds the relations/edges for other nodes in the graph.
type CheckListItemEdges struct {
	// Files holds the value of the files edge.
	Files []*File
	// WifiScan holds the value of the wifi_scan edge.
	WifiScan []*SurveyWiFiScan
	// CellScan holds the value of the cell_scan edge.
	CellScan []*SurveyCellScan
	// WorkOrder holds the value of the work_order edge.
	WorkOrder *WorkOrder
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// FilesOrErr returns the Files value or an error if the edge
// was not loaded in eager-loading.
func (e CheckListItemEdges) FilesOrErr() ([]*File, error) {
	if e.loadedTypes[0] {
		return e.Files, nil
	}
	return nil, &NotLoadedError{edge: "files"}
}

// WifiScanOrErr returns the WifiScan value or an error if the edge
// was not loaded in eager-loading.
func (e CheckListItemEdges) WifiScanOrErr() ([]*SurveyWiFiScan, error) {
	if e.loadedTypes[1] {
		return e.WifiScan, nil
	}
	return nil, &NotLoadedError{edge: "wifi_scan"}
}

// CellScanOrErr returns the CellScan value or an error if the edge
// was not loaded in eager-loading.
func (e CheckListItemEdges) CellScanOrErr() ([]*SurveyCellScan, error) {
	if e.loadedTypes[2] {
		return e.CellScan, nil
	}
	return nil, &NotLoadedError{edge: "cell_scan"}
}

// WorkOrderOrErr returns the WorkOrder value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CheckListItemEdges) WorkOrderOrErr() (*WorkOrder, error) {
	if e.loadedTypes[3] {
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
func (*CheckListItem) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullString{}, // title
		&sql.NullString{}, // type
		&sql.NullInt64{},  // index
		&sql.NullBool{},   // checked
		&sql.NullString{}, // string_val
		&sql.NullString{}, // enum_values
		&sql.NullString{}, // enum_selection_mode
		&sql.NullString{}, // selected_enum_values
		&sql.NullString{}, // yes_no_val
		&sql.NullString{}, // help_text
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*CheckListItem) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // check_list_category_check_list_items
		&sql.NullInt64{}, // work_order_check_list_items
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CheckListItem fields.
func (cli *CheckListItem) assignValues(values ...interface{}) error {
	if m, n := len(values), len(checklistitem.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	cli.ID = int(value.Int64)
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
		return fmt.Errorf("unexpected type %T for field enum_selection_mode", values[6])
	} else if value.Valid {
		cli.EnumSelectionMode = value.String
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field selected_enum_values", values[7])
	} else if value.Valid {
		cli.SelectedEnumValues = value.String
	}
	if value, ok := values[8].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field yes_no_val", values[8])
	} else if value.Valid {
		cli.YesNoVal = checklistitem.YesNoVal(value.String)
	}
	if value, ok := values[9].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field help_text", values[9])
	} else if value.Valid {
		cli.HelpText = new(string)
		*cli.HelpText = value.String
	}
	values = values[10:]
	if len(values) == len(checklistitem.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field check_list_category_check_list_items", value)
		} else if value.Valid {
			cli.check_list_category_check_list_items = new(int)
			*cli.check_list_category_check_list_items = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field work_order_check_list_items", value)
		} else if value.Valid {
			cli.work_order_check_list_items = new(int)
			*cli.work_order_check_list_items = int(value.Int64)
		}
	}
	return nil
}

// QueryFiles queries the files edge of the CheckListItem.
func (cli *CheckListItem) QueryFiles() *FileQuery {
	return (&CheckListItemClient{config: cli.config}).QueryFiles(cli)
}

// QueryWifiScan queries the wifi_scan edge of the CheckListItem.
func (cli *CheckListItem) QueryWifiScan() *SurveyWiFiScanQuery {
	return (&CheckListItemClient{config: cli.config}).QueryWifiScan(cli)
}

// QueryCellScan queries the cell_scan edge of the CheckListItem.
func (cli *CheckListItem) QueryCellScan() *SurveyCellScanQuery {
	return (&CheckListItemClient{config: cli.config}).QueryCellScan(cli)
}

// QueryWorkOrder queries the work_order edge of the CheckListItem.
func (cli *CheckListItem) QueryWorkOrder() *WorkOrderQuery {
	return (&CheckListItemClient{config: cli.config}).QueryWorkOrder(cli)
}

// Update returns a builder for updating this CheckListItem.
// Note that, you need to call CheckListItem.Unwrap() before calling this method, if this CheckListItem
// was returned from a transaction, and the transaction was committed or rolled back.
func (cli *CheckListItem) Update() *CheckListItemUpdateOne {
	return (&CheckListItemClient{config: cli.config}).UpdateOne(cli)
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
	builder.WriteString(", enum_selection_mode=")
	builder.WriteString(cli.EnumSelectionMode)
	builder.WriteString(", selected_enum_values=")
	builder.WriteString(cli.SelectedEnumValues)
	builder.WriteString(", yes_no_val=")
	builder.WriteString(fmt.Sprintf("%v", cli.YesNoVal))
	if v := cli.HelpText; v != nil {
		builder.WriteString(", help_text=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// CheckListItems is a parsable slice of CheckListItem.
type CheckListItems []*CheckListItem

func (cli CheckListItems) config(cfg config) {
	for _i := range cli {
		cli[_i].config = cfg
	}
}
