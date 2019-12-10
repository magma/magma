// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package location

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.NEQ(s.C(FieldID), id))
		},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
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
		},
	)
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
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
		},
	)
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.GT(s.C(FieldID), id))
		},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.GTE(s.C(FieldID), id))
		},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.LT(s.C(FieldID), id))
		},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.LTE(s.C(FieldID), id))
		},
	)
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldCreateTime), v))
		},
	)
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldUpdateTime), v))
		},
	)
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldName), v))
		},
	)
}

// ExternalID applies equality check predicate on the "external_id" field. It's identical to ExternalIDEQ.
func ExternalID(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldExternalID), v))
		},
	)
}

// Latitude applies equality check predicate on the "latitude" field. It's identical to LatitudeEQ.
func Latitude(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldLatitude), v))
		},
	)
}

// Longitude applies equality check predicate on the "longitude" field. It's identical to LongitudeEQ.
func Longitude(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldLongitude), v))
		},
	)
}

// SiteSurveyNeeded applies equality check predicate on the "site_survey_needed" field. It's identical to SiteSurveyNeededEQ.
func SiteSurveyNeeded(v bool) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldSiteSurveyNeeded), v))
		},
	)
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldCreateTime), v...))
		},
	)
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldCreateTime), v...))
		},
	)
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldCreateTime), v))
		},
	)
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldUpdateTime), v...))
		},
	)
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldUpdateTime), v...))
		},
	)
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldUpdateTime), v))
		},
	)
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldName), v))
		},
	)
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldName), v))
		},
	)
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldName), v...))
		},
	)
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldName), v...))
		},
	)
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldName), v))
		},
	)
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldName), v))
		},
	)
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldName), v))
		},
	)
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldName), v))
		},
	)
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldName), v))
		},
	)
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldName), v))
		},
	)
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldName), v))
		},
	)
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldName), v))
		},
	)
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldName), v))
		},
	)
}

// ExternalIDEQ applies the EQ predicate on the "external_id" field.
func ExternalIDEQ(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDNEQ applies the NEQ predicate on the "external_id" field.
func ExternalIDNEQ(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDIn applies the In predicate on the "external_id" field.
func ExternalIDIn(vs ...string) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldExternalID), v...))
		},
	)
}

// ExternalIDNotIn applies the NotIn predicate on the "external_id" field.
func ExternalIDNotIn(vs ...string) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldExternalID), v...))
		},
	)
}

// ExternalIDGT applies the GT predicate on the "external_id" field.
func ExternalIDGT(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDGTE applies the GTE predicate on the "external_id" field.
func ExternalIDGTE(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDLT applies the LT predicate on the "external_id" field.
func ExternalIDLT(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDLTE applies the LTE predicate on the "external_id" field.
func ExternalIDLTE(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDContains applies the Contains predicate on the "external_id" field.
func ExternalIDContains(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDHasPrefix applies the HasPrefix predicate on the "external_id" field.
func ExternalIDHasPrefix(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDHasSuffix applies the HasSuffix predicate on the "external_id" field.
func ExternalIDHasSuffix(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDIsNil applies the IsNil predicate on the "external_id" field.
func ExternalIDIsNil() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.IsNull(s.C(FieldExternalID)))
		},
	)
}

// ExternalIDNotNil applies the NotNil predicate on the "external_id" field.
func ExternalIDNotNil() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NotNull(s.C(FieldExternalID)))
		},
	)
}

// ExternalIDEqualFold applies the EqualFold predicate on the "external_id" field.
func ExternalIDEqualFold(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDContainsFold applies the ContainsFold predicate on the "external_id" field.
func ExternalIDContainsFold(v string) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldExternalID), v))
		},
	)
}

// LatitudeEQ applies the EQ predicate on the "latitude" field.
func LatitudeEQ(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldLatitude), v))
		},
	)
}

// LatitudeNEQ applies the NEQ predicate on the "latitude" field.
func LatitudeNEQ(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldLatitude), v))
		},
	)
}

// LatitudeIn applies the In predicate on the "latitude" field.
func LatitudeIn(vs ...float64) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldLatitude), v...))
		},
	)
}

// LatitudeNotIn applies the NotIn predicate on the "latitude" field.
func LatitudeNotIn(vs ...float64) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldLatitude), v...))
		},
	)
}

// LatitudeGT applies the GT predicate on the "latitude" field.
func LatitudeGT(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldLatitude), v))
		},
	)
}

// LatitudeGTE applies the GTE predicate on the "latitude" field.
func LatitudeGTE(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldLatitude), v))
		},
	)
}

// LatitudeLT applies the LT predicate on the "latitude" field.
func LatitudeLT(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldLatitude), v))
		},
	)
}

// LatitudeLTE applies the LTE predicate on the "latitude" field.
func LatitudeLTE(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldLatitude), v))
		},
	)
}

// LongitudeEQ applies the EQ predicate on the "longitude" field.
func LongitudeEQ(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldLongitude), v))
		},
	)
}

// LongitudeNEQ applies the NEQ predicate on the "longitude" field.
func LongitudeNEQ(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldLongitude), v))
		},
	)
}

// LongitudeIn applies the In predicate on the "longitude" field.
func LongitudeIn(vs ...float64) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldLongitude), v...))
		},
	)
}

// LongitudeNotIn applies the NotIn predicate on the "longitude" field.
func LongitudeNotIn(vs ...float64) predicate.Location {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Location(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldLongitude), v...))
		},
	)
}

// LongitudeGT applies the GT predicate on the "longitude" field.
func LongitudeGT(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldLongitude), v))
		},
	)
}

// LongitudeGTE applies the GTE predicate on the "longitude" field.
func LongitudeGTE(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldLongitude), v))
		},
	)
}

// LongitudeLT applies the LT predicate on the "longitude" field.
func LongitudeLT(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldLongitude), v))
		},
	)
}

// LongitudeLTE applies the LTE predicate on the "longitude" field.
func LongitudeLTE(v float64) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldLongitude), v))
		},
	)
}

// SiteSurveyNeededEQ applies the EQ predicate on the "site_survey_needed" field.
func SiteSurveyNeededEQ(v bool) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldSiteSurveyNeeded), v))
		},
	)
}

// SiteSurveyNeededNEQ applies the NEQ predicate on the "site_survey_needed" field.
func SiteSurveyNeededNEQ(v bool) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldSiteSurveyNeeded), v))
		},
	)
}

// SiteSurveyNeededIsNil applies the IsNil predicate on the "site_survey_needed" field.
func SiteSurveyNeededIsNil() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.IsNull(s.C(FieldSiteSurveyNeeded)))
		},
	)
}

// SiteSurveyNeededNotNil applies the NotNil predicate on the "site_survey_needed" field.
func SiteSurveyNeededNotNil() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s.Where(sql.NotNull(s.C(FieldSiteSurveyNeeded)))
		},
	)
}

// HasType applies the HasEdge predicate on the "type" edge.
func HasType() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(TypeTable, FieldID),
				sql.Edge(sql.M2O, false, TypeTable, TypeColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasTypeWith applies the HasEdge predicate on the "type" edge with a given conditions (other predicates).
func HasTypeWith(preds ...predicate.LocationType) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(FieldID).From(builder.Table(TypeInverseTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(TypeColumn), t2))
		},
	)
}

// HasParent applies the HasEdge predicate on the "parent" edge.
func HasParent() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(ParentTable, FieldID),
				sql.Edge(sql.M2O, true, ParentTable, ParentColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasParentWith applies the HasEdge predicate on the "parent" edge with a given conditions (other predicates).
func HasParentWith(preds ...predicate.Location) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(FieldID).From(builder.Table(ParentTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(ParentColumn), t2))
		},
	)
}

// HasChildren applies the HasEdge predicate on the "children" edge.
func HasChildren() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(ChildrenTable, FieldID),
				sql.Edge(sql.O2M, false, ChildrenTable, ChildrenColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasChildrenWith applies the HasEdge predicate on the "children" edge with a given conditions (other predicates).
func HasChildrenWith(preds ...predicate.Location) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(ChildrenColumn).From(builder.Table(ChildrenTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasFiles applies the HasEdge predicate on the "files" edge.
func HasFiles() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(FilesTable, FieldID),
				sql.Edge(sql.O2M, false, FilesTable, FilesColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasFilesWith applies the HasEdge predicate on the "files" edge with a given conditions (other predicates).
func HasFilesWith(preds ...predicate.File) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(FilesColumn).From(builder.Table(FilesTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasEquipment applies the HasEdge predicate on the "equipment" edge.
func HasEquipment() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(EquipmentTable, FieldID),
				sql.Edge(sql.O2M, false, EquipmentTable, EquipmentColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasEquipmentWith applies the HasEdge predicate on the "equipment" edge with a given conditions (other predicates).
func HasEquipmentWith(preds ...predicate.Equipment) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(EquipmentColumn).From(builder.Table(EquipmentTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasProperties applies the HasEdge predicate on the "properties" edge.
func HasProperties() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(PropertiesTable, FieldID),
				sql.Edge(sql.O2M, false, PropertiesTable, PropertiesColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasPropertiesWith applies the HasEdge predicate on the "properties" edge with a given conditions (other predicates).
func HasPropertiesWith(preds ...predicate.Property) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(PropertiesColumn).From(builder.Table(PropertiesTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasSurvey applies the HasEdge predicate on the "survey" edge.
func HasSurvey() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(SurveyTable, FieldID),
				sql.Edge(sql.O2M, true, SurveyTable, SurveyColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasSurveyWith applies the HasEdge predicate on the "survey" edge with a given conditions (other predicates).
func HasSurveyWith(preds ...predicate.Survey) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(SurveyColumn).From(builder.Table(SurveyTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasWifiScan applies the HasEdge predicate on the "wifi_scan" edge.
func HasWifiScan() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(WifiScanTable, FieldID),
				sql.Edge(sql.O2M, true, WifiScanTable, WifiScanColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasWifiScanWith applies the HasEdge predicate on the "wifi_scan" edge with a given conditions (other predicates).
func HasWifiScanWith(preds ...predicate.SurveyWiFiScan) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(WifiScanColumn).From(builder.Table(WifiScanTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasCellScan applies the HasEdge predicate on the "cell_scan" edge.
func HasCellScan() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(CellScanTable, FieldID),
				sql.Edge(sql.O2M, true, CellScanTable, CellScanColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasCellScanWith applies the HasEdge predicate on the "cell_scan" edge with a given conditions (other predicates).
func HasCellScanWith(preds ...predicate.SurveyCellScan) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(CellScanColumn).From(builder.Table(CellScanTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasWorkOrders applies the HasEdge predicate on the "work_orders" edge.
func HasWorkOrders() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(WorkOrdersTable, FieldID),
				sql.Edge(sql.O2M, true, WorkOrdersTable, WorkOrdersColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasWorkOrdersWith applies the HasEdge predicate on the "work_orders" edge with a given conditions (other predicates).
func HasWorkOrdersWith(preds ...predicate.WorkOrder) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(WorkOrdersColumn).From(builder.Table(WorkOrdersTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasFloorPlans applies the HasEdge predicate on the "floor_plans" edge.
func HasFloorPlans() predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(FloorPlansTable, FieldID),
				sql.Edge(sql.O2M, true, FloorPlansTable, FloorPlansColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasFloorPlansWith applies the HasEdge predicate on the "floor_plans" edge with a given conditions (other predicates).
func HasFloorPlansWith(preds ...predicate.FloorPlan) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(FloorPlansColumn).From(builder.Table(FloorPlansTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Location) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s1 := s.Clone().SetP(nil)
			for _, p := range predicates {
				p(s1)
			}
			s.Where(s1.P())
		},
	)
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Location) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			s1 := s.Clone().SetP(nil)
			for i, p := range predicates {
				if i > 0 {
					s1.Or()
				}
				p(s1)
			}
			s.Where(s1.P())
		},
	)
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Location) predicate.Location {
	return predicate.Location(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
