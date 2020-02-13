// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveywifiscan

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func IDNotIn(ids ...string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func IDGT(id string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Ssid applies equality check predicate on the "ssid" field. It's identical to SsidEQ.
func Ssid(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSsid), v))
	})
}

// Bssid applies equality check predicate on the "bssid" field. It's identical to BssidEQ.
func Bssid(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBssid), v))
	})
}

// Timestamp applies equality check predicate on the "timestamp" field. It's identical to TimestampEQ.
func Timestamp(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTimestamp), v))
	})
}

// Frequency applies equality check predicate on the "frequency" field. It's identical to FrequencyEQ.
func Frequency(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFrequency), v))
	})
}

// Channel applies equality check predicate on the "channel" field. It's identical to ChannelEQ.
func Channel(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldChannel), v))
	})
}

// Band applies equality check predicate on the "band" field. It's identical to BandEQ.
func Band(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBand), v))
	})
}

// ChannelWidth applies equality check predicate on the "channel_width" field. It's identical to ChannelWidthEQ.
func ChannelWidth(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldChannelWidth), v))
	})
}

// Capabilities applies equality check predicate on the "capabilities" field. It's identical to CapabilitiesEQ.
func Capabilities(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCapabilities), v))
	})
}

// Strength applies equality check predicate on the "strength" field. It's identical to StrengthEQ.
func Strength(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStrength), v))
	})
}

// Latitude applies equality check predicate on the "latitude" field. It's identical to LatitudeEQ.
func Latitude(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitude), v))
	})
}

// Longitude applies equality check predicate on the "longitude" field. It's identical to LongitudeEQ.
func Longitude(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitude), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// SsidEQ applies the EQ predicate on the "ssid" field.
func SsidEQ(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSsid), v))
	})
}

// SsidNEQ applies the NEQ predicate on the "ssid" field.
func SsidNEQ(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSsid), v))
	})
}

// SsidIn applies the In predicate on the "ssid" field.
func SsidIn(vs ...string) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSsid), v...))
	})
}

// SsidNotIn applies the NotIn predicate on the "ssid" field.
func SsidNotIn(vs ...string) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSsid), v...))
	})
}

// SsidGT applies the GT predicate on the "ssid" field.
func SsidGT(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSsid), v))
	})
}

// SsidGTE applies the GTE predicate on the "ssid" field.
func SsidGTE(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSsid), v))
	})
}

// SsidLT applies the LT predicate on the "ssid" field.
func SsidLT(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSsid), v))
	})
}

// SsidLTE applies the LTE predicate on the "ssid" field.
func SsidLTE(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSsid), v))
	})
}

// SsidContains applies the Contains predicate on the "ssid" field.
func SsidContains(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSsid), v))
	})
}

// SsidHasPrefix applies the HasPrefix predicate on the "ssid" field.
func SsidHasPrefix(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSsid), v))
	})
}

// SsidHasSuffix applies the HasSuffix predicate on the "ssid" field.
func SsidHasSuffix(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSsid), v))
	})
}

// SsidIsNil applies the IsNil predicate on the "ssid" field.
func SsidIsNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldSsid)))
	})
}

// SsidNotNil applies the NotNil predicate on the "ssid" field.
func SsidNotNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldSsid)))
	})
}

// SsidEqualFold applies the EqualFold predicate on the "ssid" field.
func SsidEqualFold(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSsid), v))
	})
}

// SsidContainsFold applies the ContainsFold predicate on the "ssid" field.
func SsidContainsFold(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSsid), v))
	})
}

// BssidEQ applies the EQ predicate on the "bssid" field.
func BssidEQ(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBssid), v))
	})
}

// BssidNEQ applies the NEQ predicate on the "bssid" field.
func BssidNEQ(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBssid), v))
	})
}

// BssidIn applies the In predicate on the "bssid" field.
func BssidIn(vs ...string) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldBssid), v...))
	})
}

// BssidNotIn applies the NotIn predicate on the "bssid" field.
func BssidNotIn(vs ...string) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldBssid), v...))
	})
}

// BssidGT applies the GT predicate on the "bssid" field.
func BssidGT(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldBssid), v))
	})
}

// BssidGTE applies the GTE predicate on the "bssid" field.
func BssidGTE(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldBssid), v))
	})
}

// BssidLT applies the LT predicate on the "bssid" field.
func BssidLT(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldBssid), v))
	})
}

// BssidLTE applies the LTE predicate on the "bssid" field.
func BssidLTE(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldBssid), v))
	})
}

// BssidContains applies the Contains predicate on the "bssid" field.
func BssidContains(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldBssid), v))
	})
}

// BssidHasPrefix applies the HasPrefix predicate on the "bssid" field.
func BssidHasPrefix(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldBssid), v))
	})
}

// BssidHasSuffix applies the HasSuffix predicate on the "bssid" field.
func BssidHasSuffix(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldBssid), v))
	})
}

// BssidEqualFold applies the EqualFold predicate on the "bssid" field.
func BssidEqualFold(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldBssid), v))
	})
}

// BssidContainsFold applies the ContainsFold predicate on the "bssid" field.
func BssidContainsFold(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldBssid), v))
	})
}

// TimestampEQ applies the EQ predicate on the "timestamp" field.
func TimestampEQ(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTimestamp), v))
	})
}

// TimestampNEQ applies the NEQ predicate on the "timestamp" field.
func TimestampNEQ(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTimestamp), v))
	})
}

// TimestampIn applies the In predicate on the "timestamp" field.
func TimestampIn(vs ...time.Time) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func TimestampNotIn(vs ...time.Time) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func TimestampGT(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTimestamp), v))
	})
}

// TimestampGTE applies the GTE predicate on the "timestamp" field.
func TimestampGTE(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTimestamp), v))
	})
}

// TimestampLT applies the LT predicate on the "timestamp" field.
func TimestampLT(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTimestamp), v))
	})
}

// TimestampLTE applies the LTE predicate on the "timestamp" field.
func TimestampLTE(v time.Time) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTimestamp), v))
	})
}

// FrequencyEQ applies the EQ predicate on the "frequency" field.
func FrequencyEQ(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFrequency), v))
	})
}

// FrequencyNEQ applies the NEQ predicate on the "frequency" field.
func FrequencyNEQ(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFrequency), v))
	})
}

// FrequencyIn applies the In predicate on the "frequency" field.
func FrequencyIn(vs ...int) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldFrequency), v...))
	})
}

// FrequencyNotIn applies the NotIn predicate on the "frequency" field.
func FrequencyNotIn(vs ...int) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldFrequency), v...))
	})
}

// FrequencyGT applies the GT predicate on the "frequency" field.
func FrequencyGT(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldFrequency), v))
	})
}

// FrequencyGTE applies the GTE predicate on the "frequency" field.
func FrequencyGTE(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldFrequency), v))
	})
}

// FrequencyLT applies the LT predicate on the "frequency" field.
func FrequencyLT(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldFrequency), v))
	})
}

// FrequencyLTE applies the LTE predicate on the "frequency" field.
func FrequencyLTE(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldFrequency), v))
	})
}

// ChannelEQ applies the EQ predicate on the "channel" field.
func ChannelEQ(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldChannel), v))
	})
}

// ChannelNEQ applies the NEQ predicate on the "channel" field.
func ChannelNEQ(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldChannel), v))
	})
}

// ChannelIn applies the In predicate on the "channel" field.
func ChannelIn(vs ...int) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldChannel), v...))
	})
}

// ChannelNotIn applies the NotIn predicate on the "channel" field.
func ChannelNotIn(vs ...int) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldChannel), v...))
	})
}

// ChannelGT applies the GT predicate on the "channel" field.
func ChannelGT(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldChannel), v))
	})
}

// ChannelGTE applies the GTE predicate on the "channel" field.
func ChannelGTE(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldChannel), v))
	})
}

// ChannelLT applies the LT predicate on the "channel" field.
func ChannelLT(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldChannel), v))
	})
}

// ChannelLTE applies the LTE predicate on the "channel" field.
func ChannelLTE(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldChannel), v))
	})
}

// BandEQ applies the EQ predicate on the "band" field.
func BandEQ(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBand), v))
	})
}

// BandNEQ applies the NEQ predicate on the "band" field.
func BandNEQ(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBand), v))
	})
}

// BandIn applies the In predicate on the "band" field.
func BandIn(vs ...string) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldBand), v...))
	})
}

// BandNotIn applies the NotIn predicate on the "band" field.
func BandNotIn(vs ...string) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldBand), v...))
	})
}

// BandGT applies the GT predicate on the "band" field.
func BandGT(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldBand), v))
	})
}

// BandGTE applies the GTE predicate on the "band" field.
func BandGTE(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldBand), v))
	})
}

// BandLT applies the LT predicate on the "band" field.
func BandLT(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldBand), v))
	})
}

// BandLTE applies the LTE predicate on the "band" field.
func BandLTE(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldBand), v))
	})
}

// BandContains applies the Contains predicate on the "band" field.
func BandContains(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldBand), v))
	})
}

// BandHasPrefix applies the HasPrefix predicate on the "band" field.
func BandHasPrefix(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldBand), v))
	})
}

// BandHasSuffix applies the HasSuffix predicate on the "band" field.
func BandHasSuffix(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldBand), v))
	})
}

// BandIsNil applies the IsNil predicate on the "band" field.
func BandIsNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldBand)))
	})
}

// BandNotNil applies the NotNil predicate on the "band" field.
func BandNotNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldBand)))
	})
}

// BandEqualFold applies the EqualFold predicate on the "band" field.
func BandEqualFold(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldBand), v))
	})
}

// BandContainsFold applies the ContainsFold predicate on the "band" field.
func BandContainsFold(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldBand), v))
	})
}

// ChannelWidthEQ applies the EQ predicate on the "channel_width" field.
func ChannelWidthEQ(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldChannelWidth), v))
	})
}

// ChannelWidthNEQ applies the NEQ predicate on the "channel_width" field.
func ChannelWidthNEQ(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldChannelWidth), v))
	})
}

// ChannelWidthIn applies the In predicate on the "channel_width" field.
func ChannelWidthIn(vs ...int) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldChannelWidth), v...))
	})
}

// ChannelWidthNotIn applies the NotIn predicate on the "channel_width" field.
func ChannelWidthNotIn(vs ...int) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldChannelWidth), v...))
	})
}

// ChannelWidthGT applies the GT predicate on the "channel_width" field.
func ChannelWidthGT(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldChannelWidth), v))
	})
}

// ChannelWidthGTE applies the GTE predicate on the "channel_width" field.
func ChannelWidthGTE(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldChannelWidth), v))
	})
}

// ChannelWidthLT applies the LT predicate on the "channel_width" field.
func ChannelWidthLT(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldChannelWidth), v))
	})
}

// ChannelWidthLTE applies the LTE predicate on the "channel_width" field.
func ChannelWidthLTE(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldChannelWidth), v))
	})
}

// ChannelWidthIsNil applies the IsNil predicate on the "channel_width" field.
func ChannelWidthIsNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldChannelWidth)))
	})
}

// ChannelWidthNotNil applies the NotNil predicate on the "channel_width" field.
func ChannelWidthNotNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldChannelWidth)))
	})
}

// CapabilitiesEQ applies the EQ predicate on the "capabilities" field.
func CapabilitiesEQ(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesNEQ applies the NEQ predicate on the "capabilities" field.
func CapabilitiesNEQ(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesIn applies the In predicate on the "capabilities" field.
func CapabilitiesIn(vs ...string) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCapabilities), v...))
	})
}

// CapabilitiesNotIn applies the NotIn predicate on the "capabilities" field.
func CapabilitiesNotIn(vs ...string) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCapabilities), v...))
	})
}

// CapabilitiesGT applies the GT predicate on the "capabilities" field.
func CapabilitiesGT(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesGTE applies the GTE predicate on the "capabilities" field.
func CapabilitiesGTE(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesLT applies the LT predicate on the "capabilities" field.
func CapabilitiesLT(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesLTE applies the LTE predicate on the "capabilities" field.
func CapabilitiesLTE(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesContains applies the Contains predicate on the "capabilities" field.
func CapabilitiesContains(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesHasPrefix applies the HasPrefix predicate on the "capabilities" field.
func CapabilitiesHasPrefix(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesHasSuffix applies the HasSuffix predicate on the "capabilities" field.
func CapabilitiesHasSuffix(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesIsNil applies the IsNil predicate on the "capabilities" field.
func CapabilitiesIsNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldCapabilities)))
	})
}

// CapabilitiesNotNil applies the NotNil predicate on the "capabilities" field.
func CapabilitiesNotNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldCapabilities)))
	})
}

// CapabilitiesEqualFold applies the EqualFold predicate on the "capabilities" field.
func CapabilitiesEqualFold(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCapabilities), v))
	})
}

// CapabilitiesContainsFold applies the ContainsFold predicate on the "capabilities" field.
func CapabilitiesContainsFold(v string) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCapabilities), v))
	})
}

// StrengthEQ applies the EQ predicate on the "strength" field.
func StrengthEQ(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStrength), v))
	})
}

// StrengthNEQ applies the NEQ predicate on the "strength" field.
func StrengthNEQ(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStrength), v))
	})
}

// StrengthIn applies the In predicate on the "strength" field.
func StrengthIn(vs ...int) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStrength), v...))
	})
}

// StrengthNotIn applies the NotIn predicate on the "strength" field.
func StrengthNotIn(vs ...int) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStrength), v...))
	})
}

// StrengthGT applies the GT predicate on the "strength" field.
func StrengthGT(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStrength), v))
	})
}

// StrengthGTE applies the GTE predicate on the "strength" field.
func StrengthGTE(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStrength), v))
	})
}

// StrengthLT applies the LT predicate on the "strength" field.
func StrengthLT(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStrength), v))
	})
}

// StrengthLTE applies the LTE predicate on the "strength" field.
func StrengthLTE(v int) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStrength), v))
	})
}

// LatitudeEQ applies the EQ predicate on the "latitude" field.
func LatitudeEQ(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitude), v))
	})
}

// LatitudeNEQ applies the NEQ predicate on the "latitude" field.
func LatitudeNEQ(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLatitude), v))
	})
}

// LatitudeIn applies the In predicate on the "latitude" field.
func LatitudeIn(vs ...float64) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func LatitudeNotIn(vs ...float64) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func LatitudeGT(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLatitude), v))
	})
}

// LatitudeGTE applies the GTE predicate on the "latitude" field.
func LatitudeGTE(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLatitude), v))
	})
}

// LatitudeLT applies the LT predicate on the "latitude" field.
func LatitudeLT(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLatitude), v))
	})
}

// LatitudeLTE applies the LTE predicate on the "latitude" field.
func LatitudeLTE(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLatitude), v))
	})
}

// LatitudeIsNil applies the IsNil predicate on the "latitude" field.
func LatitudeIsNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLatitude)))
	})
}

// LatitudeNotNil applies the NotNil predicate on the "latitude" field.
func LatitudeNotNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLatitude)))
	})
}

// LongitudeEQ applies the EQ predicate on the "longitude" field.
func LongitudeEQ(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitude), v))
	})
}

// LongitudeNEQ applies the NEQ predicate on the "longitude" field.
func LongitudeNEQ(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLongitude), v))
	})
}

// LongitudeIn applies the In predicate on the "longitude" field.
func LongitudeIn(vs ...float64) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func LongitudeNotIn(vs ...float64) predicate.SurveyWiFiScan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func LongitudeGT(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLongitude), v))
	})
}

// LongitudeGTE applies the GTE predicate on the "longitude" field.
func LongitudeGTE(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLongitude), v))
	})
}

// LongitudeLT applies the LT predicate on the "longitude" field.
func LongitudeLT(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLongitude), v))
	})
}

// LongitudeLTE applies the LTE predicate on the "longitude" field.
func LongitudeLTE(v float64) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLongitude), v))
	})
}

// LongitudeIsNil applies the IsNil predicate on the "longitude" field.
func LongitudeIsNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLongitude)))
	})
}

// LongitudeNotNil applies the NotNil predicate on the "longitude" field.
func LongitudeNotNil() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLongitude)))
	})
}

// HasSurveyQuestion applies the HasEdge predicate on the "survey_question" edge.
func HasSurveyQuestion() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyQuestionTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SurveyQuestionTable, SurveyQuestionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSurveyQuestionWith applies the HasEdge predicate on the "survey_question" edge with a given conditions (other predicates).
func HasSurveyQuestionWith(preds ...predicate.SurveyQuestion) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func HasLocation() predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LocationTable, LocationColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLocationWith applies the HasEdge predicate on the "location" edge with a given conditions (other predicates).
func HasLocationWith(preds ...predicate.Location) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func And(predicates ...predicate.SurveyWiFiScan) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.SurveyWiFiScan) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
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
func Not(p predicate.SurveyWiFiScan) predicate.SurveyWiFiScan {
	return predicate.SurveyWiFiScan(func(s *sql.Selector) {
		p(s.Not())
	})
}
