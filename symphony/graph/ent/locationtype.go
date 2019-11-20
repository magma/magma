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

// LocationType is the model entity for the LocationType schema.
type LocationType struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Site holds the value of the "site" field.
	Site bool `json:"site,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// MapType holds the value of the "map_type" field.
	MapType string `json:"map_type,omitempty"`
	// MapZoomLevel holds the value of the "map_zoom_level" field.
	MapZoomLevel int `json:"map_zoom_level,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
}

// FromRows scans the sql response data into LocationType.
func (lt *LocationType) FromRows(rows *sql.Rows) error {
	var scanlt struct {
		ID           int
		CreateTime   sql.NullTime
		UpdateTime   sql.NullTime
		Site         sql.NullBool
		Name         sql.NullString
		MapType      sql.NullString
		MapZoomLevel sql.NullInt64
		Index        sql.NullInt64
	}
	// the order here should be the same as in the `locationtype.Columns`.
	if err := rows.Scan(
		&scanlt.ID,
		&scanlt.CreateTime,
		&scanlt.UpdateTime,
		&scanlt.Site,
		&scanlt.Name,
		&scanlt.MapType,
		&scanlt.MapZoomLevel,
		&scanlt.Index,
	); err != nil {
		return err
	}
	lt.ID = strconv.Itoa(scanlt.ID)
	lt.CreateTime = scanlt.CreateTime.Time
	lt.UpdateTime = scanlt.UpdateTime.Time
	lt.Site = scanlt.Site.Bool
	lt.Name = scanlt.Name.String
	lt.MapType = scanlt.MapType.String
	lt.MapZoomLevel = int(scanlt.MapZoomLevel.Int64)
	lt.Index = int(scanlt.Index.Int64)
	return nil
}

// QueryLocations queries the locations edge of the LocationType.
func (lt *LocationType) QueryLocations() *LocationQuery {
	return (&LocationTypeClient{lt.config}).QueryLocations(lt)
}

// QueryPropertyTypes queries the property_types edge of the LocationType.
func (lt *LocationType) QueryPropertyTypes() *PropertyTypeQuery {
	return (&LocationTypeClient{lt.config}).QueryPropertyTypes(lt)
}

// QuerySurveyTemplateCategories queries the survey_template_categories edge of the LocationType.
func (lt *LocationType) QuerySurveyTemplateCategories() *SurveyTemplateCategoryQuery {
	return (&LocationTypeClient{lt.config}).QuerySurveyTemplateCategories(lt)
}

// Update returns a builder for updating this LocationType.
// Note that, you need to call LocationType.Unwrap() before calling this method, if this LocationType
// was returned from a transaction, and the transaction was committed or rolled back.
func (lt *LocationType) Update() *LocationTypeUpdateOne {
	return (&LocationTypeClient{lt.config}).UpdateOne(lt)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (lt *LocationType) Unwrap() *LocationType {
	tx, ok := lt.config.driver.(*txDriver)
	if !ok {
		panic("ent: LocationType is not a transactional entity")
	}
	lt.config.driver = tx.drv
	return lt
}

// String implements the fmt.Stringer.
func (lt *LocationType) String() string {
	var builder strings.Builder
	builder.WriteString("LocationType(")
	builder.WriteString(fmt.Sprintf("id=%v", lt.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(lt.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(lt.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", site=")
	builder.WriteString(fmt.Sprintf("%v", lt.Site))
	builder.WriteString(", name=")
	builder.WriteString(lt.Name)
	builder.WriteString(", map_type=")
	builder.WriteString(lt.MapType)
	builder.WriteString(", map_zoom_level=")
	builder.WriteString(fmt.Sprintf("%v", lt.MapZoomLevel))
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", lt.Index))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (lt *LocationType) id() int {
	id, _ := strconv.Atoi(lt.ID)
	return id
}

// LocationTypes is a parsable slice of LocationType.
type LocationTypes []*LocationType

// FromRows scans the sql response data into LocationTypes.
func (lt *LocationTypes) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanlt := &LocationType{}
		if err := scanlt.FromRows(rows); err != nil {
			return err
		}
		*lt = append(*lt, scanlt)
	}
	return nil
}

func (lt LocationTypes) config(cfg config) {
	for _i := range lt {
		lt[_i].config = cfg
	}
}
