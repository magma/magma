// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveycellscan

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i], _ = strconv.Atoi(ids[i])
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i], _ = strconv.Atoi(ids[i])
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// NetworkType applies equality check predicate on the "network_type" field. It's identical to NetworkTypeEQ.
func NetworkType(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNetworkType), v))
	})
}

// SignalStrength applies equality check predicate on the "signal_strength" field. It's identical to SignalStrengthEQ.
func SignalStrength(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSignalStrength), v))
	})
}

// Timestamp applies equality check predicate on the "timestamp" field. It's identical to TimestampEQ.
func Timestamp(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTimestamp), v))
	})
}

// BaseStationID applies equality check predicate on the "base_station_id" field. It's identical to BaseStationIDEQ.
func BaseStationID(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBaseStationID), v))
	})
}

// NetworkID applies equality check predicate on the "network_id" field. It's identical to NetworkIDEQ.
func NetworkID(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNetworkID), v))
	})
}

// SystemID applies equality check predicate on the "system_id" field. It's identical to SystemIDEQ.
func SystemID(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSystemID), v))
	})
}

// CellID applies equality check predicate on the "cell_id" field. It's identical to CellIDEQ.
func CellID(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCellID), v))
	})
}

// LocationAreaCode applies equality check predicate on the "location_area_code" field. It's identical to LocationAreaCodeEQ.
func LocationAreaCode(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLocationAreaCode), v))
	})
}

// MobileCountryCode applies equality check predicate on the "mobile_country_code" field. It's identical to MobileCountryCodeEQ.
func MobileCountryCode(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMobileCountryCode), v))
	})
}

// MobileNetworkCode applies equality check predicate on the "mobile_network_code" field. It's identical to MobileNetworkCodeEQ.
func MobileNetworkCode(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMobileNetworkCode), v))
	})
}

// PrimaryScramblingCode applies equality check predicate on the "primary_scrambling_code" field. It's identical to PrimaryScramblingCodeEQ.
func PrimaryScramblingCode(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPrimaryScramblingCode), v))
	})
}

// Operator applies equality check predicate on the "operator" field. It's identical to OperatorEQ.
func Operator(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOperator), v))
	})
}

// Arfcn applies equality check predicate on the "arfcn" field. It's identical to ArfcnEQ.
func Arfcn(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldArfcn), v))
	})
}

// PhysicalCellID applies equality check predicate on the "physical_cell_id" field. It's identical to PhysicalCellIDEQ.
func PhysicalCellID(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPhysicalCellID), v))
	})
}

// TrackingAreaCode applies equality check predicate on the "tracking_area_code" field. It's identical to TrackingAreaCodeEQ.
func TrackingAreaCode(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTrackingAreaCode), v))
	})
}

// TimingAdvance applies equality check predicate on the "timing_advance" field. It's identical to TimingAdvanceEQ.
func TimingAdvance(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTimingAdvance), v))
	})
}

// Earfcn applies equality check predicate on the "earfcn" field. It's identical to EarfcnEQ.
func Earfcn(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEarfcn), v))
	})
}

// Uarfcn applies equality check predicate on the "uarfcn" field. It's identical to UarfcnEQ.
func Uarfcn(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUarfcn), v))
	})
}

// Latitude applies equality check predicate on the "latitude" field. It's identical to LatitudeEQ.
func Latitude(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitude), v))
	})
}

// Longitude applies equality check predicate on the "longitude" field. It's identical to LongitudeEQ.
func Longitude(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitude), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreateTime), v...))
	})
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreateTime), v...))
	})
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdateTime), v...))
	})
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdateTime), v...))
	})
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// NetworkTypeEQ applies the EQ predicate on the "network_type" field.
func NetworkTypeEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeNEQ applies the NEQ predicate on the "network_type" field.
func NetworkTypeNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeIn applies the In predicate on the "network_type" field.
func NetworkTypeIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldNetworkType), v...))
	})
}

// NetworkTypeNotIn applies the NotIn predicate on the "network_type" field.
func NetworkTypeNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldNetworkType), v...))
	})
}

// NetworkTypeGT applies the GT predicate on the "network_type" field.
func NetworkTypeGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeGTE applies the GTE predicate on the "network_type" field.
func NetworkTypeGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeLT applies the LT predicate on the "network_type" field.
func NetworkTypeLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeLTE applies the LTE predicate on the "network_type" field.
func NetworkTypeLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeContains applies the Contains predicate on the "network_type" field.
func NetworkTypeContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeHasPrefix applies the HasPrefix predicate on the "network_type" field.
func NetworkTypeHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeHasSuffix applies the HasSuffix predicate on the "network_type" field.
func NetworkTypeHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeEqualFold applies the EqualFold predicate on the "network_type" field.
func NetworkTypeEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldNetworkType), v))
	})
}

// NetworkTypeContainsFold applies the ContainsFold predicate on the "network_type" field.
func NetworkTypeContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldNetworkType), v))
	})
}

// SignalStrengthEQ applies the EQ predicate on the "signal_strength" field.
func SignalStrengthEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSignalStrength), v))
	})
}

// SignalStrengthNEQ applies the NEQ predicate on the "signal_strength" field.
func SignalStrengthNEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSignalStrength), v))
	})
}

// SignalStrengthIn applies the In predicate on the "signal_strength" field.
func SignalStrengthIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSignalStrength), v...))
	})
}

// SignalStrengthNotIn applies the NotIn predicate on the "signal_strength" field.
func SignalStrengthNotIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSignalStrength), v...))
	})
}

// SignalStrengthGT applies the GT predicate on the "signal_strength" field.
func SignalStrengthGT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSignalStrength), v))
	})
}

// SignalStrengthGTE applies the GTE predicate on the "signal_strength" field.
func SignalStrengthGTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSignalStrength), v))
	})
}

// SignalStrengthLT applies the LT predicate on the "signal_strength" field.
func SignalStrengthLT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSignalStrength), v))
	})
}

// SignalStrengthLTE applies the LTE predicate on the "signal_strength" field.
func SignalStrengthLTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSignalStrength), v))
	})
}

// TimestampEQ applies the EQ predicate on the "timestamp" field.
func TimestampEQ(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTimestamp), v))
	})
}

// TimestampNEQ applies the NEQ predicate on the "timestamp" field.
func TimestampNEQ(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTimestamp), v))
	})
}

// TimestampIn applies the In predicate on the "timestamp" field.
func TimestampIn(vs ...time.Time) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTimestamp), v...))
	})
}

// TimestampNotIn applies the NotIn predicate on the "timestamp" field.
func TimestampNotIn(vs ...time.Time) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTimestamp), v...))
	})
}

// TimestampGT applies the GT predicate on the "timestamp" field.
func TimestampGT(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTimestamp), v))
	})
}

// TimestampGTE applies the GTE predicate on the "timestamp" field.
func TimestampGTE(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTimestamp), v))
	})
}

// TimestampLT applies the LT predicate on the "timestamp" field.
func TimestampLT(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTimestamp), v))
	})
}

// TimestampLTE applies the LTE predicate on the "timestamp" field.
func TimestampLTE(v time.Time) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTimestamp), v))
	})
}

// TimestampIsNil applies the IsNil predicate on the "timestamp" field.
func TimestampIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldTimestamp)))
	})
}

// TimestampNotNil applies the NotNil predicate on the "timestamp" field.
func TimestampNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldTimestamp)))
	})
}

// BaseStationIDEQ applies the EQ predicate on the "base_station_id" field.
func BaseStationIDEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDNEQ applies the NEQ predicate on the "base_station_id" field.
func BaseStationIDNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDIn applies the In predicate on the "base_station_id" field.
func BaseStationIDIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldBaseStationID), v...))
	})
}

// BaseStationIDNotIn applies the NotIn predicate on the "base_station_id" field.
func BaseStationIDNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldBaseStationID), v...))
	})
}

// BaseStationIDGT applies the GT predicate on the "base_station_id" field.
func BaseStationIDGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDGTE applies the GTE predicate on the "base_station_id" field.
func BaseStationIDGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDLT applies the LT predicate on the "base_station_id" field.
func BaseStationIDLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDLTE applies the LTE predicate on the "base_station_id" field.
func BaseStationIDLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDContains applies the Contains predicate on the "base_station_id" field.
func BaseStationIDContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDHasPrefix applies the HasPrefix predicate on the "base_station_id" field.
func BaseStationIDHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDHasSuffix applies the HasSuffix predicate on the "base_station_id" field.
func BaseStationIDHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDIsNil applies the IsNil predicate on the "base_station_id" field.
func BaseStationIDIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldBaseStationID)))
	})
}

// BaseStationIDNotNil applies the NotNil predicate on the "base_station_id" field.
func BaseStationIDNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldBaseStationID)))
	})
}

// BaseStationIDEqualFold applies the EqualFold predicate on the "base_station_id" field.
func BaseStationIDEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldBaseStationID), v))
	})
}

// BaseStationIDContainsFold applies the ContainsFold predicate on the "base_station_id" field.
func BaseStationIDContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldBaseStationID), v))
	})
}

// NetworkIDEQ applies the EQ predicate on the "network_id" field.
func NetworkIDEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNetworkID), v))
	})
}

// NetworkIDNEQ applies the NEQ predicate on the "network_id" field.
func NetworkIDNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldNetworkID), v))
	})
}

// NetworkIDIn applies the In predicate on the "network_id" field.
func NetworkIDIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldNetworkID), v...))
	})
}

// NetworkIDNotIn applies the NotIn predicate on the "network_id" field.
func NetworkIDNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldNetworkID), v...))
	})
}

// NetworkIDGT applies the GT predicate on the "network_id" field.
func NetworkIDGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldNetworkID), v))
	})
}

// NetworkIDGTE applies the GTE predicate on the "network_id" field.
func NetworkIDGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldNetworkID), v))
	})
}

// NetworkIDLT applies the LT predicate on the "network_id" field.
func NetworkIDLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldNetworkID), v))
	})
}

// NetworkIDLTE applies the LTE predicate on the "network_id" field.
func NetworkIDLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldNetworkID), v))
	})
}

// NetworkIDContains applies the Contains predicate on the "network_id" field.
func NetworkIDContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldNetworkID), v))
	})
}

// NetworkIDHasPrefix applies the HasPrefix predicate on the "network_id" field.
func NetworkIDHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldNetworkID), v))
	})
}

// NetworkIDHasSuffix applies the HasSuffix predicate on the "network_id" field.
func NetworkIDHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldNetworkID), v))
	})
}

// NetworkIDIsNil applies the IsNil predicate on the "network_id" field.
func NetworkIDIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldNetworkID)))
	})
}

// NetworkIDNotNil applies the NotNil predicate on the "network_id" field.
func NetworkIDNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldNetworkID)))
	})
}

// NetworkIDEqualFold applies the EqualFold predicate on the "network_id" field.
func NetworkIDEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldNetworkID), v))
	})
}

// NetworkIDContainsFold applies the ContainsFold predicate on the "network_id" field.
func NetworkIDContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldNetworkID), v))
	})
}

// SystemIDEQ applies the EQ predicate on the "system_id" field.
func SystemIDEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSystemID), v))
	})
}

// SystemIDNEQ applies the NEQ predicate on the "system_id" field.
func SystemIDNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSystemID), v))
	})
}

// SystemIDIn applies the In predicate on the "system_id" field.
func SystemIDIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSystemID), v...))
	})
}

// SystemIDNotIn applies the NotIn predicate on the "system_id" field.
func SystemIDNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSystemID), v...))
	})
}

// SystemIDGT applies the GT predicate on the "system_id" field.
func SystemIDGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSystemID), v))
	})
}

// SystemIDGTE applies the GTE predicate on the "system_id" field.
func SystemIDGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSystemID), v))
	})
}

// SystemIDLT applies the LT predicate on the "system_id" field.
func SystemIDLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSystemID), v))
	})
}

// SystemIDLTE applies the LTE predicate on the "system_id" field.
func SystemIDLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSystemID), v))
	})
}

// SystemIDContains applies the Contains predicate on the "system_id" field.
func SystemIDContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSystemID), v))
	})
}

// SystemIDHasPrefix applies the HasPrefix predicate on the "system_id" field.
func SystemIDHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSystemID), v))
	})
}

// SystemIDHasSuffix applies the HasSuffix predicate on the "system_id" field.
func SystemIDHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSystemID), v))
	})
}

// SystemIDIsNil applies the IsNil predicate on the "system_id" field.
func SystemIDIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldSystemID)))
	})
}

// SystemIDNotNil applies the NotNil predicate on the "system_id" field.
func SystemIDNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldSystemID)))
	})
}

// SystemIDEqualFold applies the EqualFold predicate on the "system_id" field.
func SystemIDEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSystemID), v))
	})
}

// SystemIDContainsFold applies the ContainsFold predicate on the "system_id" field.
func SystemIDContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSystemID), v))
	})
}

// CellIDEQ applies the EQ predicate on the "cell_id" field.
func CellIDEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCellID), v))
	})
}

// CellIDNEQ applies the NEQ predicate on the "cell_id" field.
func CellIDNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCellID), v))
	})
}

// CellIDIn applies the In predicate on the "cell_id" field.
func CellIDIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCellID), v...))
	})
}

// CellIDNotIn applies the NotIn predicate on the "cell_id" field.
func CellIDNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCellID), v...))
	})
}

// CellIDGT applies the GT predicate on the "cell_id" field.
func CellIDGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCellID), v))
	})
}

// CellIDGTE applies the GTE predicate on the "cell_id" field.
func CellIDGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCellID), v))
	})
}

// CellIDLT applies the LT predicate on the "cell_id" field.
func CellIDLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCellID), v))
	})
}

// CellIDLTE applies the LTE predicate on the "cell_id" field.
func CellIDLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCellID), v))
	})
}

// CellIDContains applies the Contains predicate on the "cell_id" field.
func CellIDContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCellID), v))
	})
}

// CellIDHasPrefix applies the HasPrefix predicate on the "cell_id" field.
func CellIDHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCellID), v))
	})
}

// CellIDHasSuffix applies the HasSuffix predicate on the "cell_id" field.
func CellIDHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCellID), v))
	})
}

// CellIDIsNil applies the IsNil predicate on the "cell_id" field.
func CellIDIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldCellID)))
	})
}

// CellIDNotNil applies the NotNil predicate on the "cell_id" field.
func CellIDNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldCellID)))
	})
}

// CellIDEqualFold applies the EqualFold predicate on the "cell_id" field.
func CellIDEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCellID), v))
	})
}

// CellIDContainsFold applies the ContainsFold predicate on the "cell_id" field.
func CellIDContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCellID), v))
	})
}

// LocationAreaCodeEQ applies the EQ predicate on the "location_area_code" field.
func LocationAreaCodeEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeNEQ applies the NEQ predicate on the "location_area_code" field.
func LocationAreaCodeNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeIn applies the In predicate on the "location_area_code" field.
func LocationAreaCodeIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLocationAreaCode), v...))
	})
}

// LocationAreaCodeNotIn applies the NotIn predicate on the "location_area_code" field.
func LocationAreaCodeNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLocationAreaCode), v...))
	})
}

// LocationAreaCodeGT applies the GT predicate on the "location_area_code" field.
func LocationAreaCodeGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeGTE applies the GTE predicate on the "location_area_code" field.
func LocationAreaCodeGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeLT applies the LT predicate on the "location_area_code" field.
func LocationAreaCodeLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeLTE applies the LTE predicate on the "location_area_code" field.
func LocationAreaCodeLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeContains applies the Contains predicate on the "location_area_code" field.
func LocationAreaCodeContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeHasPrefix applies the HasPrefix predicate on the "location_area_code" field.
func LocationAreaCodeHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeHasSuffix applies the HasSuffix predicate on the "location_area_code" field.
func LocationAreaCodeHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeIsNil applies the IsNil predicate on the "location_area_code" field.
func LocationAreaCodeIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLocationAreaCode)))
	})
}

// LocationAreaCodeNotNil applies the NotNil predicate on the "location_area_code" field.
func LocationAreaCodeNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLocationAreaCode)))
	})
}

// LocationAreaCodeEqualFold applies the EqualFold predicate on the "location_area_code" field.
func LocationAreaCodeEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldLocationAreaCode), v))
	})
}

// LocationAreaCodeContainsFold applies the ContainsFold predicate on the "location_area_code" field.
func LocationAreaCodeContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldLocationAreaCode), v))
	})
}

// MobileCountryCodeEQ applies the EQ predicate on the "mobile_country_code" field.
func MobileCountryCodeEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeNEQ applies the NEQ predicate on the "mobile_country_code" field.
func MobileCountryCodeNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeIn applies the In predicate on the "mobile_country_code" field.
func MobileCountryCodeIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldMobileCountryCode), v...))
	})
}

// MobileCountryCodeNotIn applies the NotIn predicate on the "mobile_country_code" field.
func MobileCountryCodeNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldMobileCountryCode), v...))
	})
}

// MobileCountryCodeGT applies the GT predicate on the "mobile_country_code" field.
func MobileCountryCodeGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeGTE applies the GTE predicate on the "mobile_country_code" field.
func MobileCountryCodeGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeLT applies the LT predicate on the "mobile_country_code" field.
func MobileCountryCodeLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeLTE applies the LTE predicate on the "mobile_country_code" field.
func MobileCountryCodeLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeContains applies the Contains predicate on the "mobile_country_code" field.
func MobileCountryCodeContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeHasPrefix applies the HasPrefix predicate on the "mobile_country_code" field.
func MobileCountryCodeHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeHasSuffix applies the HasSuffix predicate on the "mobile_country_code" field.
func MobileCountryCodeHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeIsNil applies the IsNil predicate on the "mobile_country_code" field.
func MobileCountryCodeIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldMobileCountryCode)))
	})
}

// MobileCountryCodeNotNil applies the NotNil predicate on the "mobile_country_code" field.
func MobileCountryCodeNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldMobileCountryCode)))
	})
}

// MobileCountryCodeEqualFold applies the EqualFold predicate on the "mobile_country_code" field.
func MobileCountryCodeEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldMobileCountryCode), v))
	})
}

// MobileCountryCodeContainsFold applies the ContainsFold predicate on the "mobile_country_code" field.
func MobileCountryCodeContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldMobileCountryCode), v))
	})
}

// MobileNetworkCodeEQ applies the EQ predicate on the "mobile_network_code" field.
func MobileNetworkCodeEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeNEQ applies the NEQ predicate on the "mobile_network_code" field.
func MobileNetworkCodeNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeIn applies the In predicate on the "mobile_network_code" field.
func MobileNetworkCodeIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldMobileNetworkCode), v...))
	})
}

// MobileNetworkCodeNotIn applies the NotIn predicate on the "mobile_network_code" field.
func MobileNetworkCodeNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldMobileNetworkCode), v...))
	})
}

// MobileNetworkCodeGT applies the GT predicate on the "mobile_network_code" field.
func MobileNetworkCodeGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeGTE applies the GTE predicate on the "mobile_network_code" field.
func MobileNetworkCodeGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeLT applies the LT predicate on the "mobile_network_code" field.
func MobileNetworkCodeLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeLTE applies the LTE predicate on the "mobile_network_code" field.
func MobileNetworkCodeLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeContains applies the Contains predicate on the "mobile_network_code" field.
func MobileNetworkCodeContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeHasPrefix applies the HasPrefix predicate on the "mobile_network_code" field.
func MobileNetworkCodeHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeHasSuffix applies the HasSuffix predicate on the "mobile_network_code" field.
func MobileNetworkCodeHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeIsNil applies the IsNil predicate on the "mobile_network_code" field.
func MobileNetworkCodeIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldMobileNetworkCode)))
	})
}

// MobileNetworkCodeNotNil applies the NotNil predicate on the "mobile_network_code" field.
func MobileNetworkCodeNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldMobileNetworkCode)))
	})
}

// MobileNetworkCodeEqualFold applies the EqualFold predicate on the "mobile_network_code" field.
func MobileNetworkCodeEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldMobileNetworkCode), v))
	})
}

// MobileNetworkCodeContainsFold applies the ContainsFold predicate on the "mobile_network_code" field.
func MobileNetworkCodeContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldMobileNetworkCode), v))
	})
}

// PrimaryScramblingCodeEQ applies the EQ predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeNEQ applies the NEQ predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeIn applies the In predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldPrimaryScramblingCode), v...))
	})
}

// PrimaryScramblingCodeNotIn applies the NotIn predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldPrimaryScramblingCode), v...))
	})
}

// PrimaryScramblingCodeGT applies the GT predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeGTE applies the GTE predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeLT applies the LT predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeLTE applies the LTE predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeContains applies the Contains predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeHasPrefix applies the HasPrefix predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeHasSuffix applies the HasSuffix predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeIsNil applies the IsNil predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldPrimaryScramblingCode)))
	})
}

// PrimaryScramblingCodeNotNil applies the NotNil predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldPrimaryScramblingCode)))
	})
}

// PrimaryScramblingCodeEqualFold applies the EqualFold predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPrimaryScramblingCode), v))
	})
}

// PrimaryScramblingCodeContainsFold applies the ContainsFold predicate on the "primary_scrambling_code" field.
func PrimaryScramblingCodeContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPrimaryScramblingCode), v))
	})
}

// OperatorEQ applies the EQ predicate on the "operator" field.
func OperatorEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOperator), v))
	})
}

// OperatorNEQ applies the NEQ predicate on the "operator" field.
func OperatorNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldOperator), v))
	})
}

// OperatorIn applies the In predicate on the "operator" field.
func OperatorIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldOperator), v...))
	})
}

// OperatorNotIn applies the NotIn predicate on the "operator" field.
func OperatorNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldOperator), v...))
	})
}

// OperatorGT applies the GT predicate on the "operator" field.
func OperatorGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldOperator), v))
	})
}

// OperatorGTE applies the GTE predicate on the "operator" field.
func OperatorGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldOperator), v))
	})
}

// OperatorLT applies the LT predicate on the "operator" field.
func OperatorLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldOperator), v))
	})
}

// OperatorLTE applies the LTE predicate on the "operator" field.
func OperatorLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldOperator), v))
	})
}

// OperatorContains applies the Contains predicate on the "operator" field.
func OperatorContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldOperator), v))
	})
}

// OperatorHasPrefix applies the HasPrefix predicate on the "operator" field.
func OperatorHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldOperator), v))
	})
}

// OperatorHasSuffix applies the HasSuffix predicate on the "operator" field.
func OperatorHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldOperator), v))
	})
}

// OperatorIsNil applies the IsNil predicate on the "operator" field.
func OperatorIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldOperator)))
	})
}

// OperatorNotNil applies the NotNil predicate on the "operator" field.
func OperatorNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldOperator)))
	})
}

// OperatorEqualFold applies the EqualFold predicate on the "operator" field.
func OperatorEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldOperator), v))
	})
}

// OperatorContainsFold applies the ContainsFold predicate on the "operator" field.
func OperatorContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldOperator), v))
	})
}

// ArfcnEQ applies the EQ predicate on the "arfcn" field.
func ArfcnEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldArfcn), v))
	})
}

// ArfcnNEQ applies the NEQ predicate on the "arfcn" field.
func ArfcnNEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldArfcn), v))
	})
}

// ArfcnIn applies the In predicate on the "arfcn" field.
func ArfcnIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldArfcn), v...))
	})
}

// ArfcnNotIn applies the NotIn predicate on the "arfcn" field.
func ArfcnNotIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldArfcn), v...))
	})
}

// ArfcnGT applies the GT predicate on the "arfcn" field.
func ArfcnGT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldArfcn), v))
	})
}

// ArfcnGTE applies the GTE predicate on the "arfcn" field.
func ArfcnGTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldArfcn), v))
	})
}

// ArfcnLT applies the LT predicate on the "arfcn" field.
func ArfcnLT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldArfcn), v))
	})
}

// ArfcnLTE applies the LTE predicate on the "arfcn" field.
func ArfcnLTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldArfcn), v))
	})
}

// ArfcnIsNil applies the IsNil predicate on the "arfcn" field.
func ArfcnIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldArfcn)))
	})
}

// ArfcnNotNil applies the NotNil predicate on the "arfcn" field.
func ArfcnNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldArfcn)))
	})
}

// PhysicalCellIDEQ applies the EQ predicate on the "physical_cell_id" field.
func PhysicalCellIDEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDNEQ applies the NEQ predicate on the "physical_cell_id" field.
func PhysicalCellIDNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDIn applies the In predicate on the "physical_cell_id" field.
func PhysicalCellIDIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldPhysicalCellID), v...))
	})
}

// PhysicalCellIDNotIn applies the NotIn predicate on the "physical_cell_id" field.
func PhysicalCellIDNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldPhysicalCellID), v...))
	})
}

// PhysicalCellIDGT applies the GT predicate on the "physical_cell_id" field.
func PhysicalCellIDGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDGTE applies the GTE predicate on the "physical_cell_id" field.
func PhysicalCellIDGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDLT applies the LT predicate on the "physical_cell_id" field.
func PhysicalCellIDLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDLTE applies the LTE predicate on the "physical_cell_id" field.
func PhysicalCellIDLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDContains applies the Contains predicate on the "physical_cell_id" field.
func PhysicalCellIDContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDHasPrefix applies the HasPrefix predicate on the "physical_cell_id" field.
func PhysicalCellIDHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDHasSuffix applies the HasSuffix predicate on the "physical_cell_id" field.
func PhysicalCellIDHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDIsNil applies the IsNil predicate on the "physical_cell_id" field.
func PhysicalCellIDIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldPhysicalCellID)))
	})
}

// PhysicalCellIDNotNil applies the NotNil predicate on the "physical_cell_id" field.
func PhysicalCellIDNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldPhysicalCellID)))
	})
}

// PhysicalCellIDEqualFold applies the EqualFold predicate on the "physical_cell_id" field.
func PhysicalCellIDEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPhysicalCellID), v))
	})
}

// PhysicalCellIDContainsFold applies the ContainsFold predicate on the "physical_cell_id" field.
func PhysicalCellIDContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPhysicalCellID), v))
	})
}

// TrackingAreaCodeEQ applies the EQ predicate on the "tracking_area_code" field.
func TrackingAreaCodeEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeNEQ applies the NEQ predicate on the "tracking_area_code" field.
func TrackingAreaCodeNEQ(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeIn applies the In predicate on the "tracking_area_code" field.
func TrackingAreaCodeIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTrackingAreaCode), v...))
	})
}

// TrackingAreaCodeNotIn applies the NotIn predicate on the "tracking_area_code" field.
func TrackingAreaCodeNotIn(vs ...string) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTrackingAreaCode), v...))
	})
}

// TrackingAreaCodeGT applies the GT predicate on the "tracking_area_code" field.
func TrackingAreaCodeGT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeGTE applies the GTE predicate on the "tracking_area_code" field.
func TrackingAreaCodeGTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeLT applies the LT predicate on the "tracking_area_code" field.
func TrackingAreaCodeLT(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeLTE applies the LTE predicate on the "tracking_area_code" field.
func TrackingAreaCodeLTE(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeContains applies the Contains predicate on the "tracking_area_code" field.
func TrackingAreaCodeContains(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeHasPrefix applies the HasPrefix predicate on the "tracking_area_code" field.
func TrackingAreaCodeHasPrefix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeHasSuffix applies the HasSuffix predicate on the "tracking_area_code" field.
func TrackingAreaCodeHasSuffix(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeIsNil applies the IsNil predicate on the "tracking_area_code" field.
func TrackingAreaCodeIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldTrackingAreaCode)))
	})
}

// TrackingAreaCodeNotNil applies the NotNil predicate on the "tracking_area_code" field.
func TrackingAreaCodeNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldTrackingAreaCode)))
	})
}

// TrackingAreaCodeEqualFold applies the EqualFold predicate on the "tracking_area_code" field.
func TrackingAreaCodeEqualFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTrackingAreaCode), v))
	})
}

// TrackingAreaCodeContainsFold applies the ContainsFold predicate on the "tracking_area_code" field.
func TrackingAreaCodeContainsFold(v string) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTrackingAreaCode), v))
	})
}

// TimingAdvanceEQ applies the EQ predicate on the "timing_advance" field.
func TimingAdvanceEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTimingAdvance), v))
	})
}

// TimingAdvanceNEQ applies the NEQ predicate on the "timing_advance" field.
func TimingAdvanceNEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTimingAdvance), v))
	})
}

// TimingAdvanceIn applies the In predicate on the "timing_advance" field.
func TimingAdvanceIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTimingAdvance), v...))
	})
}

// TimingAdvanceNotIn applies the NotIn predicate on the "timing_advance" field.
func TimingAdvanceNotIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTimingAdvance), v...))
	})
}

// TimingAdvanceGT applies the GT predicate on the "timing_advance" field.
func TimingAdvanceGT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTimingAdvance), v))
	})
}

// TimingAdvanceGTE applies the GTE predicate on the "timing_advance" field.
func TimingAdvanceGTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTimingAdvance), v))
	})
}

// TimingAdvanceLT applies the LT predicate on the "timing_advance" field.
func TimingAdvanceLT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTimingAdvance), v))
	})
}

// TimingAdvanceLTE applies the LTE predicate on the "timing_advance" field.
func TimingAdvanceLTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTimingAdvance), v))
	})
}

// TimingAdvanceIsNil applies the IsNil predicate on the "timing_advance" field.
func TimingAdvanceIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldTimingAdvance)))
	})
}

// TimingAdvanceNotNil applies the NotNil predicate on the "timing_advance" field.
func TimingAdvanceNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldTimingAdvance)))
	})
}

// EarfcnEQ applies the EQ predicate on the "earfcn" field.
func EarfcnEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEarfcn), v))
	})
}

// EarfcnNEQ applies the NEQ predicate on the "earfcn" field.
func EarfcnNEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEarfcn), v))
	})
}

// EarfcnIn applies the In predicate on the "earfcn" field.
func EarfcnIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldEarfcn), v...))
	})
}

// EarfcnNotIn applies the NotIn predicate on the "earfcn" field.
func EarfcnNotIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldEarfcn), v...))
	})
}

// EarfcnGT applies the GT predicate on the "earfcn" field.
func EarfcnGT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldEarfcn), v))
	})
}

// EarfcnGTE applies the GTE predicate on the "earfcn" field.
func EarfcnGTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldEarfcn), v))
	})
}

// EarfcnLT applies the LT predicate on the "earfcn" field.
func EarfcnLT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldEarfcn), v))
	})
}

// EarfcnLTE applies the LTE predicate on the "earfcn" field.
func EarfcnLTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldEarfcn), v))
	})
}

// EarfcnIsNil applies the IsNil predicate on the "earfcn" field.
func EarfcnIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldEarfcn)))
	})
}

// EarfcnNotNil applies the NotNil predicate on the "earfcn" field.
func EarfcnNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldEarfcn)))
	})
}

// UarfcnEQ applies the EQ predicate on the "uarfcn" field.
func UarfcnEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUarfcn), v))
	})
}

// UarfcnNEQ applies the NEQ predicate on the "uarfcn" field.
func UarfcnNEQ(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUarfcn), v))
	})
}

// UarfcnIn applies the In predicate on the "uarfcn" field.
func UarfcnIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUarfcn), v...))
	})
}

// UarfcnNotIn applies the NotIn predicate on the "uarfcn" field.
func UarfcnNotIn(vs ...int) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUarfcn), v...))
	})
}

// UarfcnGT applies the GT predicate on the "uarfcn" field.
func UarfcnGT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUarfcn), v))
	})
}

// UarfcnGTE applies the GTE predicate on the "uarfcn" field.
func UarfcnGTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUarfcn), v))
	})
}

// UarfcnLT applies the LT predicate on the "uarfcn" field.
func UarfcnLT(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUarfcn), v))
	})
}

// UarfcnLTE applies the LTE predicate on the "uarfcn" field.
func UarfcnLTE(v int) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUarfcn), v))
	})
}

// UarfcnIsNil applies the IsNil predicate on the "uarfcn" field.
func UarfcnIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldUarfcn)))
	})
}

// UarfcnNotNil applies the NotNil predicate on the "uarfcn" field.
func UarfcnNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldUarfcn)))
	})
}

// LatitudeEQ applies the EQ predicate on the "latitude" field.
func LatitudeEQ(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitude), v))
	})
}

// LatitudeNEQ applies the NEQ predicate on the "latitude" field.
func LatitudeNEQ(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLatitude), v))
	})
}

// LatitudeIn applies the In predicate on the "latitude" field.
func LatitudeIn(vs ...float64) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLatitude), v...))
	})
}

// LatitudeNotIn applies the NotIn predicate on the "latitude" field.
func LatitudeNotIn(vs ...float64) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLatitude), v...))
	})
}

// LatitudeGT applies the GT predicate on the "latitude" field.
func LatitudeGT(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLatitude), v))
	})
}

// LatitudeGTE applies the GTE predicate on the "latitude" field.
func LatitudeGTE(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLatitude), v))
	})
}

// LatitudeLT applies the LT predicate on the "latitude" field.
func LatitudeLT(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLatitude), v))
	})
}

// LatitudeLTE applies the LTE predicate on the "latitude" field.
func LatitudeLTE(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLatitude), v))
	})
}

// LatitudeIsNil applies the IsNil predicate on the "latitude" field.
func LatitudeIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLatitude)))
	})
}

// LatitudeNotNil applies the NotNil predicate on the "latitude" field.
func LatitudeNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLatitude)))
	})
}

// LongitudeEQ applies the EQ predicate on the "longitude" field.
func LongitudeEQ(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitude), v))
	})
}

// LongitudeNEQ applies the NEQ predicate on the "longitude" field.
func LongitudeNEQ(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLongitude), v))
	})
}

// LongitudeIn applies the In predicate on the "longitude" field.
func LongitudeIn(vs ...float64) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLongitude), v...))
	})
}

// LongitudeNotIn applies the NotIn predicate on the "longitude" field.
func LongitudeNotIn(vs ...float64) predicate.SurveyCellScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLongitude), v...))
	})
}

// LongitudeGT applies the GT predicate on the "longitude" field.
func LongitudeGT(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLongitude), v))
	})
}

// LongitudeGTE applies the GTE predicate on the "longitude" field.
func LongitudeGTE(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLongitude), v))
	})
}

// LongitudeLT applies the LT predicate on the "longitude" field.
func LongitudeLT(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLongitude), v))
	})
}

// LongitudeLTE applies the LTE predicate on the "longitude" field.
func LongitudeLTE(v float64) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLongitude), v))
	})
}

// LongitudeIsNil applies the IsNil predicate on the "longitude" field.
func LongitudeIsNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLongitude)))
	})
}

// LongitudeNotNil applies the NotNil predicate on the "longitude" field.
func LongitudeNotNil() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLongitude)))
	})
}

// HasSurveyQuestion applies the HasEdge predicate on the "survey_question" edge.
func HasSurveyQuestion() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyQuestionTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SurveyQuestionTable, SurveyQuestionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSurveyQuestionWith applies the HasEdge predicate on the "survey_question" edge with a given conditions (other predicates).
func HasSurveyQuestionWith(preds ...predicate.SurveyQuestion) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyQuestionInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SurveyQuestionTable, SurveyQuestionColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLocation applies the HasEdge predicate on the "location" edge.
func HasLocation() predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LocationTable, LocationColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLocationWith applies the HasEdge predicate on the "location" edge with a given conditions (other predicates).
func HasLocationWith(preds ...predicate.Location) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LocationTable, LocationColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.SurveyCellScan) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.SurveyCellScan) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.SurveyCellScan) predicate.SurveyCellScan {
	return predicate.SurveyCellScan(func(s *sql.Selector) {
		p(s.Not())
	})
}
