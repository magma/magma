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
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
)

// Location is the model entity for the Location schema.
type Location struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
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
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LocationQuery when eager-loading is set.
	Edges             LocationEdges `json:"edges"`
	location_type     *int
	location_children *int
}

// LocationEdges holds the relations/edges for other nodes in the graph.
type LocationEdges struct {
	// Type holds the value of the type edge.
	Type *LocationType `gqlgen:"locationType"`
	// Parent holds the value of the parent edge.
	Parent *Location `gqlgen:"parentLocation"`
	// Children holds the value of the children edge.
	Children []*Location `gqlgen:"children"`
	// Files holds the value of the files edge.
	Files []*File `gqlgen:"files,images"`
	// Hyperlinks holds the value of the hyperlinks edge.
	Hyperlinks []*Hyperlink `gqlgen:"hyperlinks"`
	// Equipment holds the value of the equipment edge.
	Equipment []*Equipment `gqlgen:"equipments"`
	// Properties holds the value of the properties edge.
	Properties []*Property `gqlgen:"properties"`
	// Survey holds the value of the survey edge.
	Survey []*Survey `gqlgen:"surveys"`
	// WifiScan holds the value of the wifi_scan edge.
	WifiScan []*SurveyWiFiScan `gqlgen:"wifiData"`
	// CellScan holds the value of the cell_scan edge.
	CellScan []*SurveyCellScan `gqlgen:"cellData"`
	// WorkOrders holds the value of the work_orders edge.
	WorkOrders []*WorkOrder `gqlgen:"workOrders"`
	// FloorPlans holds the value of the floor_plans edge.
	FloorPlans []*FloorPlan `gqlgen:"floorPlans"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [12]bool
}

// TypeOrErr returns the Type value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LocationEdges) TypeOrErr() (*LocationType, error) {
	if e.loadedTypes[0] {
		if e.Type == nil {
			// The edge type was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: locationtype.Label}
		}
		return e.Type, nil
	}
	return nil, &NotLoadedError{edge: "type"}
}

// ParentOrErr returns the Parent value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LocationEdges) ParentOrErr() (*Location, error) {
	if e.loadedTypes[1] {
		if e.Parent == nil {
			// The edge parent was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.Parent, nil
	}
	return nil, &NotLoadedError{edge: "parent"}
}

// ChildrenOrErr returns the Children value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) ChildrenOrErr() ([]*Location, error) {
	if e.loadedTypes[2] {
		return e.Children, nil
	}
	return nil, &NotLoadedError{edge: "children"}
}

// FilesOrErr returns the Files value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) FilesOrErr() ([]*File, error) {
	if e.loadedTypes[3] {
		return e.Files, nil
	}
	return nil, &NotLoadedError{edge: "files"}
}

// HyperlinksOrErr returns the Hyperlinks value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) HyperlinksOrErr() ([]*Hyperlink, error) {
	if e.loadedTypes[4] {
		return e.Hyperlinks, nil
	}
	return nil, &NotLoadedError{edge: "hyperlinks"}
}

// EquipmentOrErr returns the Equipment value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) EquipmentOrErr() ([]*Equipment, error) {
	if e.loadedTypes[5] {
		return e.Equipment, nil
	}
	return nil, &NotLoadedError{edge: "equipment"}
}

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) PropertiesOrErr() ([]*Property, error) {
	if e.loadedTypes[6] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// SurveyOrErr returns the Survey value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) SurveyOrErr() ([]*Survey, error) {
	if e.loadedTypes[7] {
		return e.Survey, nil
	}
	return nil, &NotLoadedError{edge: "survey"}
}

// WifiScanOrErr returns the WifiScan value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) WifiScanOrErr() ([]*SurveyWiFiScan, error) {
	if e.loadedTypes[8] {
		return e.WifiScan, nil
	}
	return nil, &NotLoadedError{edge: "wifi_scan"}
}

// CellScanOrErr returns the CellScan value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) CellScanOrErr() ([]*SurveyCellScan, error) {
	if e.loadedTypes[9] {
		return e.CellScan, nil
	}
	return nil, &NotLoadedError{edge: "cell_scan"}
}

// WorkOrdersOrErr returns the WorkOrders value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) WorkOrdersOrErr() ([]*WorkOrder, error) {
	if e.loadedTypes[10] {
		return e.WorkOrders, nil
	}
	return nil, &NotLoadedError{edge: "work_orders"}
}

// FloorPlansOrErr returns the FloorPlans value or an error if the edge
// was not loaded in eager-loading.
func (e LocationEdges) FloorPlansOrErr() ([]*FloorPlan, error) {
	if e.loadedTypes[11] {
		return e.FloorPlans, nil
	}
	return nil, &NotLoadedError{edge: "floor_plans"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Location) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},   // id
		&sql.NullTime{},    // create_time
		&sql.NullTime{},    // update_time
		&sql.NullString{},  // name
		&sql.NullString{},  // external_id
		&sql.NullFloat64{}, // latitude
		&sql.NullFloat64{}, // longitude
		&sql.NullBool{},    // site_survey_needed
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Location) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // location_type
		&sql.NullInt64{}, // location_children
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Location fields.
func (l *Location) assignValues(values ...interface{}) error {
	if m, n := len(values), len(location.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	l.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		l.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		l.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		l.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field external_id", values[3])
	} else if value.Valid {
		l.ExternalID = value.String
	}
	if value, ok := values[4].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field latitude", values[4])
	} else if value.Valid {
		l.Latitude = value.Float64
	}
	if value, ok := values[5].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field longitude", values[5])
	} else if value.Valid {
		l.Longitude = value.Float64
	}
	if value, ok := values[6].(*sql.NullBool); !ok {
		return fmt.Errorf("unexpected type %T for field site_survey_needed", values[6])
	} else if value.Valid {
		l.SiteSurveyNeeded = value.Bool
	}
	values = values[7:]
	if len(values) == len(location.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_type", value)
		} else if value.Valid {
			l.location_type = new(int)
			*l.location_type = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_children", value)
		} else if value.Valid {
			l.location_children = new(int)
			*l.location_children = int(value.Int64)
		}
	}
	return nil
}

// QueryType queries the type edge of the Location.
func (l *Location) QueryType() *LocationTypeQuery {
	return (&LocationClient{config: l.config}).QueryType(l)
}

// QueryParent queries the parent edge of the Location.
func (l *Location) QueryParent() *LocationQuery {
	return (&LocationClient{config: l.config}).QueryParent(l)
}

// QueryChildren queries the children edge of the Location.
func (l *Location) QueryChildren() *LocationQuery {
	return (&LocationClient{config: l.config}).QueryChildren(l)
}

// QueryFiles queries the files edge of the Location.
func (l *Location) QueryFiles() *FileQuery {
	return (&LocationClient{config: l.config}).QueryFiles(l)
}

// QueryHyperlinks queries the hyperlinks edge of the Location.
func (l *Location) QueryHyperlinks() *HyperlinkQuery {
	return (&LocationClient{config: l.config}).QueryHyperlinks(l)
}

// QueryEquipment queries the equipment edge of the Location.
func (l *Location) QueryEquipment() *EquipmentQuery {
	return (&LocationClient{config: l.config}).QueryEquipment(l)
}

// QueryProperties queries the properties edge of the Location.
func (l *Location) QueryProperties() *PropertyQuery {
	return (&LocationClient{config: l.config}).QueryProperties(l)
}

// QuerySurvey queries the survey edge of the Location.
func (l *Location) QuerySurvey() *SurveyQuery {
	return (&LocationClient{config: l.config}).QuerySurvey(l)
}

// QueryWifiScan queries the wifi_scan edge of the Location.
func (l *Location) QueryWifiScan() *SurveyWiFiScanQuery {
	return (&LocationClient{config: l.config}).QueryWifiScan(l)
}

// QueryCellScan queries the cell_scan edge of the Location.
func (l *Location) QueryCellScan() *SurveyCellScanQuery {
	return (&LocationClient{config: l.config}).QueryCellScan(l)
}

// QueryWorkOrders queries the work_orders edge of the Location.
func (l *Location) QueryWorkOrders() *WorkOrderQuery {
	return (&LocationClient{config: l.config}).QueryWorkOrders(l)
}

// QueryFloorPlans queries the floor_plans edge of the Location.
func (l *Location) QueryFloorPlans() *FloorPlanQuery {
	return (&LocationClient{config: l.config}).QueryFloorPlans(l)
}

// Update returns a builder for updating this Location.
// Note that, you need to call Location.Unwrap() before calling this method, if this Location
// was returned from a transaction, and the transaction was committed or rolled back.
func (l *Location) Update() *LocationUpdateOne {
	return (&LocationClient{config: l.config}).UpdateOne(l)
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

// Locations is a parsable slice of Location.
type Locations []*Location

func (l Locations) config(cfg config) {
	for _i := range l {
		l[_i].config = cfg
	}
}
