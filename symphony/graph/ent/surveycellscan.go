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

// SurveyCellScan is the model entity for the SurveyCellScan schema.
type SurveyCellScan struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// NetworkType holds the value of the "network_type" field.
	NetworkType string `json:"network_type,omitempty"`
	// SignalStrength holds the value of the "signal_strength" field.
	SignalStrength int `json:"signal_strength,omitempty"`
	// Timestamp holds the value of the "timestamp" field.
	Timestamp time.Time `json:"timestamp,omitempty"`
	// BaseStationID holds the value of the "base_station_id" field.
	BaseStationID string `json:"base_station_id,omitempty"`
	// NetworkID holds the value of the "network_id" field.
	NetworkID string `json:"network_id,omitempty"`
	// SystemID holds the value of the "system_id" field.
	SystemID string `json:"system_id,omitempty"`
	// CellID holds the value of the "cell_id" field.
	CellID string `json:"cell_id,omitempty"`
	// LocationAreaCode holds the value of the "location_area_code" field.
	LocationAreaCode string `json:"location_area_code,omitempty"`
	// MobileCountryCode holds the value of the "mobile_country_code" field.
	MobileCountryCode string `json:"mobile_country_code,omitempty"`
	// MobileNetworkCode holds the value of the "mobile_network_code" field.
	MobileNetworkCode string `json:"mobile_network_code,omitempty"`
	// PrimaryScramblingCode holds the value of the "primary_scrambling_code" field.
	PrimaryScramblingCode string `json:"primary_scrambling_code,omitempty"`
	// Operator holds the value of the "operator" field.
	Operator string `json:"operator,omitempty"`
	// Arfcn holds the value of the "arfcn" field.
	Arfcn int `json:"arfcn,omitempty"`
	// PhysicalCellID holds the value of the "physical_cell_id" field.
	PhysicalCellID string `json:"physical_cell_id,omitempty"`
	// TrackingAreaCode holds the value of the "tracking_area_code" field.
	TrackingAreaCode string `json:"tracking_area_code,omitempty"`
	// TimingAdvance holds the value of the "timing_advance" field.
	TimingAdvance int `json:"timing_advance,omitempty"`
	// Earfcn holds the value of the "earfcn" field.
	Earfcn int `json:"earfcn,omitempty"`
	// Uarfcn holds the value of the "uarfcn" field.
	Uarfcn int `json:"uarfcn,omitempty"`
	// Latitude holds the value of the "latitude" field.
	Latitude float64 `json:"latitude,omitempty"`
	// Longitude holds the value of the "longitude" field.
	Longitude float64 `json:"longitude,omitempty"`
}

// FromRows scans the sql response data into SurveyCellScan.
func (scs *SurveyCellScan) FromRows(rows *sql.Rows) error {
	var scanscs struct {
		ID                    int
		CreateTime            sql.NullTime
		UpdateTime            sql.NullTime
		NetworkType           sql.NullString
		SignalStrength        sql.NullInt64
		Timestamp             sql.NullTime
		BaseStationID         sql.NullString
		NetworkID             sql.NullString
		SystemID              sql.NullString
		CellID                sql.NullString
		LocationAreaCode      sql.NullString
		MobileCountryCode     sql.NullString
		MobileNetworkCode     sql.NullString
		PrimaryScramblingCode sql.NullString
		Operator              sql.NullString
		Arfcn                 sql.NullInt64
		PhysicalCellID        sql.NullString
		TrackingAreaCode      sql.NullString
		TimingAdvance         sql.NullInt64
		Earfcn                sql.NullInt64
		Uarfcn                sql.NullInt64
		Latitude              sql.NullFloat64
		Longitude             sql.NullFloat64
	}
	// the order here should be the same as in the `surveycellscan.Columns`.
	if err := rows.Scan(
		&scanscs.ID,
		&scanscs.CreateTime,
		&scanscs.UpdateTime,
		&scanscs.NetworkType,
		&scanscs.SignalStrength,
		&scanscs.Timestamp,
		&scanscs.BaseStationID,
		&scanscs.NetworkID,
		&scanscs.SystemID,
		&scanscs.CellID,
		&scanscs.LocationAreaCode,
		&scanscs.MobileCountryCode,
		&scanscs.MobileNetworkCode,
		&scanscs.PrimaryScramblingCode,
		&scanscs.Operator,
		&scanscs.Arfcn,
		&scanscs.PhysicalCellID,
		&scanscs.TrackingAreaCode,
		&scanscs.TimingAdvance,
		&scanscs.Earfcn,
		&scanscs.Uarfcn,
		&scanscs.Latitude,
		&scanscs.Longitude,
	); err != nil {
		return err
	}
	scs.ID = strconv.Itoa(scanscs.ID)
	scs.CreateTime = scanscs.CreateTime.Time
	scs.UpdateTime = scanscs.UpdateTime.Time
	scs.NetworkType = scanscs.NetworkType.String
	scs.SignalStrength = int(scanscs.SignalStrength.Int64)
	scs.Timestamp = scanscs.Timestamp.Time
	scs.BaseStationID = scanscs.BaseStationID.String
	scs.NetworkID = scanscs.NetworkID.String
	scs.SystemID = scanscs.SystemID.String
	scs.CellID = scanscs.CellID.String
	scs.LocationAreaCode = scanscs.LocationAreaCode.String
	scs.MobileCountryCode = scanscs.MobileCountryCode.String
	scs.MobileNetworkCode = scanscs.MobileNetworkCode.String
	scs.PrimaryScramblingCode = scanscs.PrimaryScramblingCode.String
	scs.Operator = scanscs.Operator.String
	scs.Arfcn = int(scanscs.Arfcn.Int64)
	scs.PhysicalCellID = scanscs.PhysicalCellID.String
	scs.TrackingAreaCode = scanscs.TrackingAreaCode.String
	scs.TimingAdvance = int(scanscs.TimingAdvance.Int64)
	scs.Earfcn = int(scanscs.Earfcn.Int64)
	scs.Uarfcn = int(scanscs.Uarfcn.Int64)
	scs.Latitude = scanscs.Latitude.Float64
	scs.Longitude = scanscs.Longitude.Float64
	return nil
}

// QuerySurveyQuestion queries the survey_question edge of the SurveyCellScan.
func (scs *SurveyCellScan) QuerySurveyQuestion() *SurveyQuestionQuery {
	return (&SurveyCellScanClient{scs.config}).QuerySurveyQuestion(scs)
}

// QueryLocation queries the location edge of the SurveyCellScan.
func (scs *SurveyCellScan) QueryLocation() *LocationQuery {
	return (&SurveyCellScanClient{scs.config}).QueryLocation(scs)
}

// Update returns a builder for updating this SurveyCellScan.
// Note that, you need to call SurveyCellScan.Unwrap() before calling this method, if this SurveyCellScan
// was returned from a transaction, and the transaction was committed or rolled back.
func (scs *SurveyCellScan) Update() *SurveyCellScanUpdateOne {
	return (&SurveyCellScanClient{scs.config}).UpdateOne(scs)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (scs *SurveyCellScan) Unwrap() *SurveyCellScan {
	tx, ok := scs.config.driver.(*txDriver)
	if !ok {
		panic("ent: SurveyCellScan is not a transactional entity")
	}
	scs.config.driver = tx.drv
	return scs
}

// String implements the fmt.Stringer.
func (scs *SurveyCellScan) String() string {
	var builder strings.Builder
	builder.WriteString("SurveyCellScan(")
	builder.WriteString(fmt.Sprintf("id=%v", scs.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(scs.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(scs.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", network_type=")
	builder.WriteString(scs.NetworkType)
	builder.WriteString(", signal_strength=")
	builder.WriteString(fmt.Sprintf("%v", scs.SignalStrength))
	builder.WriteString(", timestamp=")
	builder.WriteString(scs.Timestamp.Format(time.ANSIC))
	builder.WriteString(", base_station_id=")
	builder.WriteString(scs.BaseStationID)
	builder.WriteString(", network_id=")
	builder.WriteString(scs.NetworkID)
	builder.WriteString(", system_id=")
	builder.WriteString(scs.SystemID)
	builder.WriteString(", cell_id=")
	builder.WriteString(scs.CellID)
	builder.WriteString(", location_area_code=")
	builder.WriteString(scs.LocationAreaCode)
	builder.WriteString(", mobile_country_code=")
	builder.WriteString(scs.MobileCountryCode)
	builder.WriteString(", mobile_network_code=")
	builder.WriteString(scs.MobileNetworkCode)
	builder.WriteString(", primary_scrambling_code=")
	builder.WriteString(scs.PrimaryScramblingCode)
	builder.WriteString(", operator=")
	builder.WriteString(scs.Operator)
	builder.WriteString(", arfcn=")
	builder.WriteString(fmt.Sprintf("%v", scs.Arfcn))
	builder.WriteString(", physical_cell_id=")
	builder.WriteString(scs.PhysicalCellID)
	builder.WriteString(", tracking_area_code=")
	builder.WriteString(scs.TrackingAreaCode)
	builder.WriteString(", timing_advance=")
	builder.WriteString(fmt.Sprintf("%v", scs.TimingAdvance))
	builder.WriteString(", earfcn=")
	builder.WriteString(fmt.Sprintf("%v", scs.Earfcn))
	builder.WriteString(", uarfcn=")
	builder.WriteString(fmt.Sprintf("%v", scs.Uarfcn))
	builder.WriteString(", latitude=")
	builder.WriteString(fmt.Sprintf("%v", scs.Latitude))
	builder.WriteString(", longitude=")
	builder.WriteString(fmt.Sprintf("%v", scs.Longitude))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (scs *SurveyCellScan) id() int {
	id, _ := strconv.Atoi(scs.ID)
	return id
}

// SurveyCellScans is a parsable slice of SurveyCellScan.
type SurveyCellScans []*SurveyCellScan

// FromRows scans the sql response data into SurveyCellScans.
func (scs *SurveyCellScans) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanscs := &SurveyCellScan{}
		if err := scanscs.FromRows(rows); err != nil {
			return err
		}
		*scs = append(*scs, scanscs)
	}
	return nil
}

func (scs SurveyCellScans) config(cfg config) {
	for _i := range scs {
		scs[_i].config = cfg
	}
}
