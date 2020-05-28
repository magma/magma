// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/ent/locationtype"
)

// LocationType is the model entity for the LocationType schema.
type LocationType struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Site holds the value of the "site" field.
	Site bool `json:"site,omitempty" gqlgen:"isSite"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// MapType holds the value of the "map_type" field.
	MapType string `json:"map_type,omitempty"`
	// MapZoomLevel holds the value of the "map_zoom_level" field.
	MapZoomLevel int `json:"map_zoom_level,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LocationTypeQuery when eager-loading is set.
	Edges LocationTypeEdges `json:"edges"`
}

// LocationTypeEdges holds the relations/edges for other nodes in the graph.
type LocationTypeEdges struct {
	// Locations holds the value of the locations edge.
	Locations []*Location `gqlgen:"locations"`
	// PropertyTypes holds the value of the property_types edge.
	PropertyTypes []*PropertyType `gqlgen:"propertyTypes"`
	// SurveyTemplateCategories holds the value of the survey_template_categories edge.
	SurveyTemplateCategories []*SurveyTemplateCategory `gqlgen:"surveyTemplateCategories"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// LocationsOrErr returns the Locations value or an error if the edge
// was not loaded in eager-loading.
func (e LocationTypeEdges) LocationsOrErr() ([]*Location, error) {
	if e.loadedTypes[0] {
		return e.Locations, nil
	}
	return nil, &NotLoadedError{edge: "locations"}
}

// PropertyTypesOrErr returns the PropertyTypes value or an error if the edge
// was not loaded in eager-loading.
func (e LocationTypeEdges) PropertyTypesOrErr() ([]*PropertyType, error) {
	if e.loadedTypes[1] {
		return e.PropertyTypes, nil
	}
	return nil, &NotLoadedError{edge: "property_types"}
}

// SurveyTemplateCategoriesOrErr returns the SurveyTemplateCategories value or an error if the edge
// was not loaded in eager-loading.
func (e LocationTypeEdges) SurveyTemplateCategoriesOrErr() ([]*SurveyTemplateCategory, error) {
	if e.loadedTypes[2] {
		return e.SurveyTemplateCategories, nil
	}
	return nil, &NotLoadedError{edge: "survey_template_categories"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LocationType) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullBool{},   // site
		&sql.NullString{}, // name
		&sql.NullString{}, // map_type
		&sql.NullInt64{},  // map_zoom_level
		&sql.NullInt64{},  // index
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LocationType fields.
func (lt *LocationType) assignValues(values ...interface{}) error {
	if m, n := len(values), len(locationtype.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	lt.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		lt.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		lt.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field site", values[2])
	} else if value.Valid {
		lt.Site = value.Bool
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[3])
	} else if value.Valid {
		lt.Name = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field map_type", values[4])
	} else if value.Valid {
		lt.MapType = value.String
	}
	if value, ok := values[5].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field map_zoom_level", values[5])
	} else if value.Valid {
		lt.MapZoomLevel = int(value.Int64)
	}
	if value, ok := values[6].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field index", values[6])
	} else if value.Valid {
		lt.Index = int(value.Int64)
	}
	return nil
}

// QueryLocations queries the locations edge of the LocationType.
func (lt *LocationType) QueryLocations() *LocationQuery {
	return (&LocationTypeClient{config: lt.config}).QueryLocations(lt)
}

// QueryPropertyTypes queries the property_types edge of the LocationType.
func (lt *LocationType) QueryPropertyTypes() *PropertyTypeQuery {
	return (&LocationTypeClient{config: lt.config}).QueryPropertyTypes(lt)
}

// QuerySurveyTemplateCategories queries the survey_template_categories edge of the LocationType.
func (lt *LocationType) QuerySurveyTemplateCategories() *SurveyTemplateCategoryQuery {
	return (&LocationTypeClient{config: lt.config}).QuerySurveyTemplateCategories(lt)
}

// Update returns a builder for updating this LocationType.
// Note that, you need to call LocationType.Unwrap() before calling this method, if this LocationType
// was returned from a transaction, and the transaction was committed or rolled back.
func (lt *LocationType) Update() *LocationTypeUpdateOne {
	return (&LocationTypeClient{config: lt.config}).UpdateOne(lt)
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

// LocationTypes is a parsable slice of LocationType.
type LocationTypes []*LocationType

func (lt LocationTypes) config(cfg config) {
	for _i := range lt {
		lt[_i].config = cfg
	}
}
