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

// WorkOrder is the model entity for the WorkOrder schema.
type WorkOrder struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty"`
	// Priority holds the value of the "priority" field.
	Priority string `json:"priority,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// OwnerName holds the value of the "owner_name" field.
	OwnerName string `json:"owner_name,omitempty"`
	// InstallDate holds the value of the "install_date" field.
	InstallDate time.Time `json:"install_date,omitempty"`
	// CreationDate holds the value of the "creation_date" field.
	CreationDate time.Time `json:"creation_date,omitempty"`
	// Assignee holds the value of the "assignee" field.
	Assignee string `json:"assignee,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
}

// FromRows scans the sql response data into WorkOrder.
func (wo *WorkOrder) FromRows(rows *sql.Rows) error {
	var scanwo struct {
		ID           int
		CreateTime   sql.NullTime
		UpdateTime   sql.NullTime
		Name         sql.NullString
		Status       sql.NullString
		Priority     sql.NullString
		Description  sql.NullString
		OwnerName    sql.NullString
		InstallDate  sql.NullTime
		CreationDate sql.NullTime
		Assignee     sql.NullString
		Index        sql.NullInt64
	}
	// the order here should be the same as in the `workorder.Columns`.
	if err := rows.Scan(
		&scanwo.ID,
		&scanwo.CreateTime,
		&scanwo.UpdateTime,
		&scanwo.Name,
		&scanwo.Status,
		&scanwo.Priority,
		&scanwo.Description,
		&scanwo.OwnerName,
		&scanwo.InstallDate,
		&scanwo.CreationDate,
		&scanwo.Assignee,
		&scanwo.Index,
	); err != nil {
		return err
	}
	wo.ID = strconv.Itoa(scanwo.ID)
	wo.CreateTime = scanwo.CreateTime.Time
	wo.UpdateTime = scanwo.UpdateTime.Time
	wo.Name = scanwo.Name.String
	wo.Status = scanwo.Status.String
	wo.Priority = scanwo.Priority.String
	wo.Description = scanwo.Description.String
	wo.OwnerName = scanwo.OwnerName.String
	wo.InstallDate = scanwo.InstallDate.Time
	wo.CreationDate = scanwo.CreationDate.Time
	wo.Assignee = scanwo.Assignee.String
	wo.Index = int(scanwo.Index.Int64)
	return nil
}

// QueryType queries the type edge of the WorkOrder.
func (wo *WorkOrder) QueryType() *WorkOrderTypeQuery {
	return (&WorkOrderClient{wo.config}).QueryType(wo)
}

// QueryEquipment queries the equipment edge of the WorkOrder.
func (wo *WorkOrder) QueryEquipment() *EquipmentQuery {
	return (&WorkOrderClient{wo.config}).QueryEquipment(wo)
}

// QueryLinks queries the links edge of the WorkOrder.
func (wo *WorkOrder) QueryLinks() *LinkQuery {
	return (&WorkOrderClient{wo.config}).QueryLinks(wo)
}

// QueryFiles queries the files edge of the WorkOrder.
func (wo *WorkOrder) QueryFiles() *FileQuery {
	return (&WorkOrderClient{wo.config}).QueryFiles(wo)
}

// QueryLocation queries the location edge of the WorkOrder.
func (wo *WorkOrder) QueryLocation() *LocationQuery {
	return (&WorkOrderClient{wo.config}).QueryLocation(wo)
}

// QueryComments queries the comments edge of the WorkOrder.
func (wo *WorkOrder) QueryComments() *CommentQuery {
	return (&WorkOrderClient{wo.config}).QueryComments(wo)
}

// QueryProperties queries the properties edge of the WorkOrder.
func (wo *WorkOrder) QueryProperties() *PropertyQuery {
	return (&WorkOrderClient{wo.config}).QueryProperties(wo)
}

// QueryCheckListItems queries the check_list_items edge of the WorkOrder.
func (wo *WorkOrder) QueryCheckListItems() *CheckListItemQuery {
	return (&WorkOrderClient{wo.config}).QueryCheckListItems(wo)
}

// QueryTechnician queries the technician edge of the WorkOrder.
func (wo *WorkOrder) QueryTechnician() *TechnicianQuery {
	return (&WorkOrderClient{wo.config}).QueryTechnician(wo)
}

// QueryProject queries the project edge of the WorkOrder.
func (wo *WorkOrder) QueryProject() *ProjectQuery {
	return (&WorkOrderClient{wo.config}).QueryProject(wo)
}

// Update returns a builder for updating this WorkOrder.
// Note that, you need to call WorkOrder.Unwrap() before calling this method, if this WorkOrder
// was returned from a transaction, and the transaction was committed or rolled back.
func (wo *WorkOrder) Update() *WorkOrderUpdateOne {
	return (&WorkOrderClient{wo.config}).UpdateOne(wo)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (wo *WorkOrder) Unwrap() *WorkOrder {
	tx, ok := wo.config.driver.(*txDriver)
	if !ok {
		panic("ent: WorkOrder is not a transactional entity")
	}
	wo.config.driver = tx.drv
	return wo
}

// String implements the fmt.Stringer.
func (wo *WorkOrder) String() string {
	var builder strings.Builder
	builder.WriteString("WorkOrder(")
	builder.WriteString(fmt.Sprintf("id=%v", wo.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(wo.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(wo.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(wo.Name)
	builder.WriteString(", status=")
	builder.WriteString(wo.Status)
	builder.WriteString(", priority=")
	builder.WriteString(wo.Priority)
	builder.WriteString(", description=")
	builder.WriteString(wo.Description)
	builder.WriteString(", owner_name=")
	builder.WriteString(wo.OwnerName)
	builder.WriteString(", install_date=")
	builder.WriteString(wo.InstallDate.Format(time.ANSIC))
	builder.WriteString(", creation_date=")
	builder.WriteString(wo.CreationDate.Format(time.ANSIC))
	builder.WriteString(", assignee=")
	builder.WriteString(wo.Assignee)
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", wo.Index))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (wo *WorkOrder) id() int {
	id, _ := strconv.Atoi(wo.ID)
	return id
}

// WorkOrders is a parsable slice of WorkOrder.
type WorkOrders []*WorkOrder

// FromRows scans the sql response data into WorkOrders.
func (wo *WorkOrders) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanwo := &WorkOrder{}
		if err := scanwo.FromRows(rows); err != nil {
			return err
		}
		*wo = append(*wo, scanwo)
	}
	return nil
}

func (wo WorkOrders) config(cfg config) {
	for _i := range wo {
		wo[_i].config = cfg
	}
}
