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
}

// FromRows scans the sql response data into SurveyWiFiScan.
func (swfs *SurveyWiFiScan) FromRows(rows *sql.Rows) error {
	var scanswfs struct {
		ID           int
		CreateTime   sql.NullTime
		UpdateTime   sql.NullTime
		Ssid         sql.NullString
		Bssid        sql.NullString
		Timestamp    sql.NullTime
		Frequency    sql.NullInt64
		Channel      sql.NullInt64
		Band         sql.NullString
		ChannelWidth sql.NullInt64
		Capabilities sql.NullString
		Strength     sql.NullInt64
		Latitude     sql.NullFloat64
		Longitude    sql.NullFloat64
	}
	// the order here should be the same as in the `surveywifiscan.Columns`.
	if err := rows.Scan(
		&scanswfs.ID,
		&scanswfs.CreateTime,
		&scanswfs.UpdateTime,
		&scanswfs.Ssid,
		&scanswfs.Bssid,
		&scanswfs.Timestamp,
		&scanswfs.Frequency,
		&scanswfs.Channel,
		&scanswfs.Band,
		&scanswfs.ChannelWidth,
		&scanswfs.Capabilities,
		&scanswfs.Strength,
		&scanswfs.Latitude,
		&scanswfs.Longitude,
	); err != nil {
		return err
	}
	swfs.ID = strconv.Itoa(scanswfs.ID)
	swfs.CreateTime = scanswfs.CreateTime.Time
	swfs.UpdateTime = scanswfs.UpdateTime.Time
	swfs.Ssid = scanswfs.Ssid.String
	swfs.Bssid = scanswfs.Bssid.String
	swfs.Timestamp = scanswfs.Timestamp.Time
	swfs.Frequency = int(scanswfs.Frequency.Int64)
	swfs.Channel = int(scanswfs.Channel.Int64)
	swfs.Band = scanswfs.Band.String
	swfs.ChannelWidth = int(scanswfs.ChannelWidth.Int64)
	swfs.Capabilities = scanswfs.Capabilities.String
	swfs.Strength = int(scanswfs.Strength.Int64)
	swfs.Latitude = scanswfs.Latitude.Float64
	swfs.Longitude = scanswfs.Longitude.Float64
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

// FromRows scans the sql response data into SurveyWiFiScans.
func (swfs *SurveyWiFiScans) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanswfs := &SurveyWiFiScan{}
		if err := scanswfs.FromRows(rows); err != nil {
			return err
		}
		*swfs = append(*swfs, scanswfs)
	}
	return nil
}

func (swfs SurveyWiFiScans) config(cfg config) {
	for _i := range swfs {
		swfs[_i].config = cfg
	}
}
