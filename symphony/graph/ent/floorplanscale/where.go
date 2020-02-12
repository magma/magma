// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package floorplanscale

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
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
func IDNotIn(ids ...string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
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
func IDGT(id string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// ReferencePoint1X applies equality check predicate on the "reference_point1_x" field. It's identical to ReferencePoint1XEQ.
func ReferencePoint1X(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldReferencePoint1X), v))
	})
}

// ReferencePoint1Y applies equality check predicate on the "reference_point1_y" field. It's identical to ReferencePoint1YEQ.
func ReferencePoint1Y(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldReferencePoint1Y), v))
	})
}

// ReferencePoint2X applies equality check predicate on the "reference_point2_x" field. It's identical to ReferencePoint2XEQ.
func ReferencePoint2X(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldReferencePoint2X), v))
	})
}

// ReferencePoint2Y applies equality check predicate on the "reference_point2_y" field. It's identical to ReferencePoint2YEQ.
func ReferencePoint2Y(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldReferencePoint2Y), v))
	})
}

// ScaleInMeters applies equality check predicate on the "scale_in_meters" field. It's identical to ScaleInMetersEQ.
func ScaleInMeters(v float64) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldScaleInMeters), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// ReferencePoint1XEQ applies the EQ predicate on the "reference_point1_x" field.
func ReferencePoint1XEQ(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldReferencePoint1X), v))
	})
}

// ReferencePoint1XNEQ applies the NEQ predicate on the "reference_point1_x" field.
func ReferencePoint1XNEQ(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldReferencePoint1X), v))
	})
}

// ReferencePoint1XIn applies the In predicate on the "reference_point1_x" field.
func ReferencePoint1XIn(vs ...int) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldReferencePoint1X), v...))
	})
}

// ReferencePoint1XNotIn applies the NotIn predicate on the "reference_point1_x" field.
func ReferencePoint1XNotIn(vs ...int) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldReferencePoint1X), v...))
	})
}

// ReferencePoint1XGT applies the GT predicate on the "reference_point1_x" field.
func ReferencePoint1XGT(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldReferencePoint1X), v))
	})
}

// ReferencePoint1XGTE applies the GTE predicate on the "reference_point1_x" field.
func ReferencePoint1XGTE(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldReferencePoint1X), v))
	})
}

// ReferencePoint1XLT applies the LT predicate on the "reference_point1_x" field.
func ReferencePoint1XLT(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldReferencePoint1X), v))
	})
}

// ReferencePoint1XLTE applies the LTE predicate on the "reference_point1_x" field.
func ReferencePoint1XLTE(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldReferencePoint1X), v))
	})
}

// ReferencePoint1YEQ applies the EQ predicate on the "reference_point1_y" field.
func ReferencePoint1YEQ(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldReferencePoint1Y), v))
	})
}

// ReferencePoint1YNEQ applies the NEQ predicate on the "reference_point1_y" field.
func ReferencePoint1YNEQ(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldReferencePoint1Y), v))
	})
}

// ReferencePoint1YIn applies the In predicate on the "reference_point1_y" field.
func ReferencePoint1YIn(vs ...int) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldReferencePoint1Y), v...))
	})
}

// ReferencePoint1YNotIn applies the NotIn predicate on the "reference_point1_y" field.
func ReferencePoint1YNotIn(vs ...int) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldReferencePoint1Y), v...))
	})
}

// ReferencePoint1YGT applies the GT predicate on the "reference_point1_y" field.
func ReferencePoint1YGT(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldReferencePoint1Y), v))
	})
}

// ReferencePoint1YGTE applies the GTE predicate on the "reference_point1_y" field.
func ReferencePoint1YGTE(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldReferencePoint1Y), v))
	})
}

// ReferencePoint1YLT applies the LT predicate on the "reference_point1_y" field.
func ReferencePoint1YLT(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldReferencePoint1Y), v))
	})
}

// ReferencePoint1YLTE applies the LTE predicate on the "reference_point1_y" field.
func ReferencePoint1YLTE(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldReferencePoint1Y), v))
	})
}

// ReferencePoint2XEQ applies the EQ predicate on the "reference_point2_x" field.
func ReferencePoint2XEQ(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldReferencePoint2X), v))
	})
}

// ReferencePoint2XNEQ applies the NEQ predicate on the "reference_point2_x" field.
func ReferencePoint2XNEQ(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldReferencePoint2X), v))
	})
}

// ReferencePoint2XIn applies the In predicate on the "reference_point2_x" field.
func ReferencePoint2XIn(vs ...int) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldReferencePoint2X), v...))
	})
}

// ReferencePoint2XNotIn applies the NotIn predicate on the "reference_point2_x" field.
func ReferencePoint2XNotIn(vs ...int) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldReferencePoint2X), v...))
	})
}

// ReferencePoint2XGT applies the GT predicate on the "reference_point2_x" field.
func ReferencePoint2XGT(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldReferencePoint2X), v))
	})
}

// ReferencePoint2XGTE applies the GTE predicate on the "reference_point2_x" field.
func ReferencePoint2XGTE(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldReferencePoint2X), v))
	})
}

// ReferencePoint2XLT applies the LT predicate on the "reference_point2_x" field.
func ReferencePoint2XLT(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldReferencePoint2X), v))
	})
}

// ReferencePoint2XLTE applies the LTE predicate on the "reference_point2_x" field.
func ReferencePoint2XLTE(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldReferencePoint2X), v))
	})
}

// ReferencePoint2YEQ applies the EQ predicate on the "reference_point2_y" field.
func ReferencePoint2YEQ(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldReferencePoint2Y), v))
	})
}

// ReferencePoint2YNEQ applies the NEQ predicate on the "reference_point2_y" field.
func ReferencePoint2YNEQ(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldReferencePoint2Y), v))
	})
}

// ReferencePoint2YIn applies the In predicate on the "reference_point2_y" field.
func ReferencePoint2YIn(vs ...int) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldReferencePoint2Y), v...))
	})
}

// ReferencePoint2YNotIn applies the NotIn predicate on the "reference_point2_y" field.
func ReferencePoint2YNotIn(vs ...int) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldReferencePoint2Y), v...))
	})
}

// ReferencePoint2YGT applies the GT predicate on the "reference_point2_y" field.
func ReferencePoint2YGT(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldReferencePoint2Y), v))
	})
}

// ReferencePoint2YGTE applies the GTE predicate on the "reference_point2_y" field.
func ReferencePoint2YGTE(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldReferencePoint2Y), v))
	})
}

// ReferencePoint2YLT applies the LT predicate on the "reference_point2_y" field.
func ReferencePoint2YLT(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldReferencePoint2Y), v))
	})
}

// ReferencePoint2YLTE applies the LTE predicate on the "reference_point2_y" field.
func ReferencePoint2YLTE(v int) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldReferencePoint2Y), v))
	})
}

// ScaleInMetersEQ applies the EQ predicate on the "scale_in_meters" field.
func ScaleInMetersEQ(v float64) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldScaleInMeters), v))
	})
}

// ScaleInMetersNEQ applies the NEQ predicate on the "scale_in_meters" field.
func ScaleInMetersNEQ(v float64) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldScaleInMeters), v))
	})
}

// ScaleInMetersIn applies the In predicate on the "scale_in_meters" field.
func ScaleInMetersIn(vs ...float64) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldScaleInMeters), v...))
	})
}

// ScaleInMetersNotIn applies the NotIn predicate on the "scale_in_meters" field.
func ScaleInMetersNotIn(vs ...float64) predicate.FloorPlanScale {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldScaleInMeters), v...))
	})
}

// ScaleInMetersGT applies the GT predicate on the "scale_in_meters" field.
func ScaleInMetersGT(v float64) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldScaleInMeters), v))
	})
}

// ScaleInMetersGTE applies the GTE predicate on the "scale_in_meters" field.
func ScaleInMetersGTE(v float64) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldScaleInMeters), v))
	})
}

// ScaleInMetersLT applies the LT predicate on the "scale_in_meters" field.
func ScaleInMetersLT(v float64) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldScaleInMeters), v))
	})
}

// ScaleInMetersLTE applies the LTE predicate on the "scale_in_meters" field.
func ScaleInMetersLTE(v float64) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldScaleInMeters), v))
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.FloorPlanScale) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.FloorPlanScale) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
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
func Not(p predicate.FloorPlanScale) predicate.FloorPlanScale {
	return predicate.FloorPlanScale(func(s *sql.Selector) {
		p(s.Not())
	})
}
