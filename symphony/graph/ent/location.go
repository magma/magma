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

// Location is the model entity for the Location schema.
type Location struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// ExternalID holds the value of the "external_id" field.
	ExternalID string `json:"external_id,omitempty"`
	// Latitude holds the value of the "latitude" field.
	Latitude float64 `json:"latitude,omitempty"`
	// Longitude holds the value of the "longitude" field.
	Longitude float64 `json:"longitude,omitempty"`
	// SiteSurveyNeeded holds the value of the "site_survey_needed" field.
	SiteSurveyNeeded bool `json:"site_survey_needed,omitempty" gqlgen:"siteSurveyNeeded"`
}

// FromRows scans the sql response data into Location.
func (l *Location) FromRows(rows *sql.Rows) error {
	var scanl struct {
		ID               int
		CreateTime       sql.NullTime
		UpdateTime       sql.NullTime
		Name             sql.NullString
		ExternalID       sql.NullString
		Latitude         sql.NullFloat64
		Longitude        sql.NullFloat64
		SiteSurveyNeeded sql.NullBool
	}
	// the order here should be the same as in the `location.Columns`.
	if err := rows.Scan(
		&scanl.ID,
		&scanl.CreateTime,
		&scanl.UpdateTime,
		&scanl.Name,
		&scanl.ExternalID,
		&scanl.Latitude,
		&scanl.Longitude,
		&scanl.SiteSurveyNeeded,
	); err != nil {
		return err
	}
	l.ID = strconv.Itoa(scanl.ID)
	l.CreateTime = scanl.CreateTime.Time
	l.UpdateTime = scanl.UpdateTime.Time
	l.Name = scanl.Name.String
	l.ExternalID = scanl.ExternalID.String
	l.Latitude = scanl.Latitude.Float64
	l.Longitude = scanl.Longitude.Float64
	l.SiteSurveyNeeded = scanl.SiteSurveyNeeded.Bool
	return nil
}

// QueryType queries the type edge of the Location.
func (l *Location) QueryType() *LocationTypeQuery {
	return (&LocationClient{l.config}).QueryType(l)
}

// QueryParent queries the parent edge of the Location.
func (l *Location) QueryParent() *LocationQuery {
	return (&LocationClient{l.config}).QueryParent(l)
}

// QueryChildren queries the children edge of the Location.
func (l *Location) QueryChildren() *LocationQuery {
	return (&LocationClient{l.config}).QueryChildren(l)
}

// QueryFiles queries the files edge of the Location.
func (l *Location) QueryFiles() *FileQuery {
	return (&LocationClient{l.config}).QueryFiles(l)
}

// QueryEquipment queries the equipment edge of the Location.
func (l *Location) QueryEquipment() *EquipmentQuery {
	return (&LocationClient{l.config}).QueryEquipment(l)
}

// QueryProperties queries the properties edge of the Location.
func (l *Location) QueryProperties() *PropertyQuery {
	return (&LocationClient{l.config}).QueryProperties(l)
}

// QuerySurvey queries the survey edge of the Location.
func (l *Location) QuerySurvey() *SurveyQuery {
	return (&LocationClient{l.config}).QuerySurvey(l)
}

// QueryWifiScan queries the wifi_scan edge of the Location.
func (l *Location) QueryWifiScan() *SurveyWiFiScanQuery {
	return (&LocationClient{l.config}).QueryWifiScan(l)
}

// QueryCellScan queries the cell_scan edge of the Location.
func (l *Location) QueryCellScan() *SurveyCellScanQuery {
	return (&LocationClient{l.config}).QueryCellScan(l)
}

// QueryWorkOrders queries the work_orders edge of the Location.
func (l *Location) QueryWorkOrders() *WorkOrderQuery {
	return (&LocationClient{l.config}).QueryWorkOrders(l)
}

// Update returns a builder for updating this Location.
// Note that, you need to call Location.Unwrap() before calling this method, if this Location
// was returned from a transaction, and the transaction was committed or rolled back.
func (l *Location) Update() *LocationUpdateOne {
	return (&LocationClient{l.config}).UpdateOne(l)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (l *Location) Unwrap() *Location {
	tx, ok := l.config.driver.(*txDriver)
	if !ok {
		panic("ent: Location is not a transactional entity")
	}
	l.config.driver = tx.drv
	return l
}

// String implements the fmt.Stringer.
func (l *Location) String() string {
	var builder strings.Builder
	builder.WriteString("Location(")
	builder.WriteString(fmt.Sprintf("id=%v", l.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(l.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(l.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(l.Name)
	builder.WriteString(", external_id=")
	builder.WriteString(l.ExternalID)
	builder.WriteString(", latitude=")
	builder.WriteString(fmt.Sprintf("%v", l.Latitude))
	builder.WriteString(", longitude=")
	builder.WriteString(fmt.Sprintf("%v", l.Longitude))
	builder.WriteString(", site_survey_needed=")
	builder.WriteString(fmt.Sprintf("%v", l.SiteSurveyNeeded))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (l *Location) id() int {
	id, _ := strconv.Atoi(l.ID)
	return id
}

// Locations is a parsable slice of Location.
type Locations []*Location

// FromRows scans the sql response data into Locations.
func (l *Locations) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanl := &Location{}
		if err := scanl.FromRows(rows); err != nil {
			return err
		}
		*l = append(*l, scanl)
	}
	return nil
}

func (l Locations) config(cfg config) {
	for _i := range l {
		l[_i].config = cfg
	}
}
