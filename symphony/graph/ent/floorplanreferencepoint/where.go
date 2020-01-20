// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package floorplanreferencepoint

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func IDNotIn(ids ...string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func IDGT(id string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	},
	)
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	},
	)
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	},
	)
}

// X applies equality check predicate on the "x" field. It's identical to XEQ.
func X(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldX), v))
	},
	)
}

// Y applies equality check predicate on the "y" field. It's identical to YEQ.
func Y(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldY), v))
	},
	)
}

// Latitude applies equality check predicate on the "latitude" field. It's identical to LatitudeEQ.
func Latitude(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitude), v))
	},
	)
}

// Longitude applies equality check predicate on the "longitude" field. It's identical to LongitudeEQ.
func Longitude(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitude), v))
	},
	)
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	},
	)
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	},
	)
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	},
	)
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	},
	)
}

// XEQ applies the EQ predicate on the "x" field.
func XEQ(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldX), v))
	},
	)
}

// XNEQ applies the NEQ predicate on the "x" field.
func XNEQ(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldX), v))
	},
	)
}

// XIn applies the In predicate on the "x" field.
func XIn(vs ...int) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldX), v...))
	},
	)
}

// XNotIn applies the NotIn predicate on the "x" field.
func XNotIn(vs ...int) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldX), v...))
	},
	)
}

// XGT applies the GT predicate on the "x" field.
func XGT(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldX), v))
	},
	)
}

// XGTE applies the GTE predicate on the "x" field.
func XGTE(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldX), v))
	},
	)
}

// XLT applies the LT predicate on the "x" field.
func XLT(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldX), v))
	},
	)
}

// XLTE applies the LTE predicate on the "x" field.
func XLTE(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldX), v))
	},
	)
}

// YEQ applies the EQ predicate on the "y" field.
func YEQ(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldY), v))
	},
	)
}

// YNEQ applies the NEQ predicate on the "y" field.
func YNEQ(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldY), v))
	},
	)
}

// YIn applies the In predicate on the "y" field.
func YIn(vs ...int) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldY), v...))
	},
	)
}

// YNotIn applies the NotIn predicate on the "y" field.
func YNotIn(vs ...int) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldY), v...))
	},
	)
}

// YGT applies the GT predicate on the "y" field.
func YGT(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldY), v))
	},
	)
}

// YGTE applies the GTE predicate on the "y" field.
func YGTE(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldY), v))
	},
	)
}

// YLT applies the LT predicate on the "y" field.
func YLT(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldY), v))
	},
	)
}

// YLTE applies the LTE predicate on the "y" field.
func YLTE(v int) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldY), v))
	},
	)
}

// LatitudeEQ applies the EQ predicate on the "latitude" field.
func LatitudeEQ(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitude), v))
	},
	)
}

// LatitudeNEQ applies the NEQ predicate on the "latitude" field.
func LatitudeNEQ(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLatitude), v))
	},
	)
}

// LatitudeIn applies the In predicate on the "latitude" field.
func LatitudeIn(vs ...float64) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func LatitudeNotIn(vs ...float64) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func LatitudeGT(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLatitude), v))
	},
	)
}

// LatitudeGTE applies the GTE predicate on the "latitude" field.
func LatitudeGTE(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLatitude), v))
	},
	)
}

// LatitudeLT applies the LT predicate on the "latitude" field.
func LatitudeLT(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLatitude), v))
	},
	)
}

// LatitudeLTE applies the LTE predicate on the "latitude" field.
func LatitudeLTE(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLatitude), v))
	},
	)
}

// LongitudeEQ applies the EQ predicate on the "longitude" field.
func LongitudeEQ(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitude), v))
	},
	)
}

// LongitudeNEQ applies the NEQ predicate on the "longitude" field.
func LongitudeNEQ(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLongitude), v))
	},
	)
}

// LongitudeIn applies the In predicate on the "longitude" field.
func LongitudeIn(vs ...float64) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func LongitudeNotIn(vs ...float64) predicate.FloorPlanReferencePoint {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
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
func LongitudeGT(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLongitude), v))
	},
	)
}

// LongitudeGTE applies the GTE predicate on the "longitude" field.
func LongitudeGTE(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLongitude), v))
	},
	)
}

// LongitudeLT applies the LT predicate on the "longitude" field.
func LongitudeLT(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLongitude), v))
	},
	)
}

// LongitudeLTE applies the LTE predicate on the "longitude" field.
func LongitudeLTE(v float64) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLongitude), v))
	},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.FloorPlanReferencePoint) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(
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
func Or(predicates ...predicate.FloorPlanReferencePoint) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(
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
func Not(p predicate.FloorPlanReferencePoint) predicate.FloorPlanReferencePoint {
	return predicate.FloorPlanReferencePoint(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
