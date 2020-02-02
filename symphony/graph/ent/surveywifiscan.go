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
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyWiFiScan is the model entity for the SurveyWiFiScan schema.
type SurveyWiFiScan struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Ssid holds the value of the "ssid" field.
	Ssid string `json:"ssid,omitempty"`
	// Bssid holds the value of the "bssid" field.
	Bssid string `json:"bssid,omitempty"`
	// Timestamp holds the value of the "timestamp" field.
	Timestamp time.Time `json:"timestamp,omitempty"`
	// Frequency holds the value of the "frequency" field.
	Frequency int `json:"frequency,omitempty"`
	// Channel holds the value of the "channel" field.
	Channel int `json:"channel,omitempty"`
	// Band holds the value of the "band" field.
	Band string `json:"band,omitempty"`
	// ChannelWidth holds the value of the "channel_width" field.
	ChannelWidth int `json:"channel_width,omitempty"`
	// Capabilities holds the value of the "capabilities" field.
	Capabilities string `json:"capabilities,omitempty"`
	// Strength holds the value of the "strength" field.
	Strength int `json:"strength,omitempty"`
	// Latitude holds the value of the "latitude" field.
	Latitude float64 `json:"latitude,omitempty"`
	// Longitude holds the value of the "longitude" field.
	Longitude float64 `json:"longitude,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SurveyWiFiScanQuery when eager-loading is set.
	Edges              SurveyWiFiScanEdges `json:"edges"`
	survey_question_id *string
	location_id        *string
}

// SurveyWiFiScanEdges holds the relations/edges for other nodes in the graph.
type SurveyWiFiScanEdges struct {
	// SurveyQuestion holds the value of the survey_question edge.
	SurveyQuestion *SurveyQuestion
	// Location holds the value of the location edge.
	Location *Location
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// SurveyQuestionOrErr returns the SurveyQuestion value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SurveyWiFiScanEdges) SurveyQuestionOrErr() (*SurveyQuestion, error) {
	if e.loadedTypes[0] {
		if e.SurveyQuestion == nil {
			// The edge survey_question was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: surveyquestion.Label}
		}
		return e.SurveyQuestion, nil
	}
	return nil, &NotLoadedError{edge: "survey_question"}
}

// LocationOrErr returns the Location value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SurveyWiFiScanEdges) LocationOrErr() (*Location, error) {
	if e.loadedTypes[1] {
		if e.Location == nil {
			// The edge location was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: location.Label}
		}
		return e.Location, nil
	}
	return nil, &NotLoadedError{edge: "location"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SurveyWiFiScan) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},   // id
		&sql.NullTime{},    // create_time
		&sql.NullTime{},    // update_time
		&sql.NullString{},  // ssid
		&sql.NullString{},  // bssid
		&sql.NullTime{},    // timestamp
		&sql.NullInt64{},   // frequency
		&sql.NullInt64{},   // channel
		&sql.NullString{},  // band
		&sql.NullInt64{},   // channel_width
		&sql.NullString{},  // capabilities
		&sql.NullInt64{},   // strength
		&sql.NullFloat64{}, // latitude
		&sql.NullFloat64{}, // longitude
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*SurveyWiFiScan) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // survey_question_id
		&sql.NullInt64{}, // location_id
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SurveyWiFiScan fields.
func (swfs *SurveyWiFiScan) assignValues(values ...interface{}) error {
	if m, n := len(values), len(surveywifiscan.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	swfs.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		swfs.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		swfs.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field ssid", values[2])
	} else if value.Valid {
		swfs.Ssid = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field bssid", values[3])
	} else if value.Valid {
		swfs.Bssid = value.String
	}
	if value, ok := values[4].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field timestamp", values[4])
	} else if value.Valid {
		swfs.Timestamp = value.Time
	}
	if value, ok := values[5].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field frequency", values[5])
	} else if value.Valid {
		swfs.Frequency = int(value.Int64)
	}
	if value, ok := values[6].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field channel", values[6])
	} else if value.Valid {
		swfs.Channel = int(value.Int64)
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field band", values[7])
	} else if value.Valid {
		swfs.Band = value.String
	}
	if value, ok := values[8].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field channel_width", values[8])
	} else if value.Valid {
		swfs.ChannelWidth = int(value.Int64)
	}
	if value, ok := values[9].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field capabilities", values[9])
	} else if value.Valid {
		swfs.Capabilities = value.String
	}
	if value, ok := values[10].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field strength", values[10])
	} else if value.Valid {
		swfs.Strength = int(value.Int64)
	}
	if value, ok := values[11].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field latitude", values[11])
	} else if value.Valid {
		swfs.Latitude = value.Float64
	}
	if value, ok := values[12].(*sql.NullFloat64); !ok {
		return fmt.Errorf("unexpected type %T for field longitude", values[12])
	} else if value.Valid {
		swfs.Longitude = value.Float64
	}
	values = values[13:]
	if len(values) == len(surveywifiscan.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field survey_question_id", value)
		} else if value.Valid {
			swfs.survey_question_id = new(string)
			*swfs.survey_question_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_id", value)
		} else if value.Valid {
			swfs.location_id = new(string)
			*swfs.location_id = strconv.FormatInt(value.Int64, 10)
		}
	}
	return nil
}

// QuerySurveyQuestion queries the survey_question edge of the SurveyWiFiScan.
func (swfs *SurveyWiFiScan) QuerySurveyQuestion() *SurveyQuestionQuery {
	return (&SurveyWiFiScanClient{swfs.config}).QuerySurveyQuestion(swfs)
}

// QueryLocation queries the location edge of the SurveyWiFiScan.
func (swfs *SurveyWiFiScan) QueryLocation() *LocationQuery {
	return (&SurveyWiFiScanClient{swfs.config}).QueryLocation(swfs)
}

// Update returns a builder for updating this SurveyWiFiScan.
// Note that, you need to call SurveyWiFiScan.Unwrap() before calling this method, if this SurveyWiFiScan
// was returned from a transaction, and the transaction was committed or rolled back.
func (swfs *SurveyWiFiScan) Update() *SurveyWiFiScanUpdateOne {
	return (&SurveyWiFiScanClient{swfs.config}).UpdateOne(swfs)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (swfs *SurveyWiFiScan) Unwrap() *SurveyWiFiScan {
	tx, ok := swfs.config.driver.(*txDriver)
	if !ok {
		panic("ent: SurveyWiFiScan is not a transactional entity")
	}
	swfs.config.driver = tx.drv
	return swfs
}

// String implements the fmt.Stringer.
func (swfs *SurveyWiFiScan) String() string {
	var builder strings.Builder
	builder.WriteString("SurveyWiFiScan(")
	builder.WriteString(fmt.Sprintf("id=%v", swfs.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(swfs.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(swfs.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ssid=")
	builder.WriteString(swfs.Ssid)
	builder.WriteString(", bssid=")
	builder.WriteString(swfs.Bssid)
	builder.WriteString(", timestamp=")
	builder.WriteString(swfs.Timestamp.Format(time.ANSIC))
	builder.WriteString(", frequency=")
	builder.WriteString(fmt.Sprintf("%v", swfs.Frequency))
	builder.WriteString(", channel=")
	builder.WriteString(fmt.Sprintf("%v", swfs.Channel))
	builder.WriteString(", band=")
	builder.WriteString(swfs.Band)
	builder.WriteString(", channel_width=")
	builder.WriteString(fmt.Sprintf("%v", swfs.ChannelWidth))
	builder.WriteString(", capabilities=")
	builder.WriteString(swfs.Capabilities)
	builder.WriteString(", strength=")
	builder.WriteString(fmt.Sprintf("%v", swfs.Strength))
	builder.WriteString(", latitude=")
	builder.WriteString(fmt.Sprintf("%v", swfs.Latitude))
	builder.WriteString(", longitude=")
	builder.WriteString(fmt.Sprintf("%v", swfs.Longitude))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (swfs *SurveyWiFiScan) id() int {
	id, _ := strconv.Atoi(swfs.ID)
	return id
}

// SurveyWiFiScans is a parsable slice of SurveyWiFiScan.
type SurveyWiFiScans []*SurveyWiFiScan

func (swfs SurveyWiFiScans) config(cfg config) {
	for _i := range swfs {
		swfs[_i].config = cfg
	}
}
