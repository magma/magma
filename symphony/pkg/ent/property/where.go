// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package property

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// IntVal applies equality check predicate on the "int_val" field. It's identical to IntValEQ.
func IntVal(v int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIntVal), v))
	})
}

// BoolVal applies equality check predicate on the "bool_val" field. It's identical to BoolValEQ.
func BoolVal(v bool) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBoolVal), v))
	})
}

// FloatVal applies equality check predicate on the "float_val" field. It's identical to FloatValEQ.
func FloatVal(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFloatVal), v))
	})
}

// LatitudeVal applies equality check predicate on the "latitude_val" field. It's identical to LatitudeValEQ.
func LatitudeVal(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitudeVal), v))
	})
}

// LongitudeVal applies equality check predicate on the "longitude_val" field. It's identical to LongitudeValEQ.
func LongitudeVal(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitudeVal), v))
	})
}

// RangeFromVal applies equality check predicate on the "range_from_val" field. It's identical to RangeFromValEQ.
func RangeFromVal(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRangeFromVal), v))
	})
}

// RangeToVal applies equality check predicate on the "range_to_val" field. It's identical to RangeToValEQ.
func RangeToVal(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRangeToVal), v))
	})
}

// StringVal applies equality check predicate on the "string_val" field. It's identical to StringValEQ.
func StringVal(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStringVal), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// IntValEQ applies the EQ predicate on the "int_val" field.
func IntValEQ(v int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIntVal), v))
	})
}

// IntValNEQ applies the NEQ predicate on the "int_val" field.
func IntValNEQ(v int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIntVal), v))
	})
}

// IntValIn applies the In predicate on the "int_val" field.
func IntValIn(vs ...int) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldIntVal), v...))
	})
}

// IntValNotIn applies the NotIn predicate on the "int_val" field.
func IntValNotIn(vs ...int) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldIntVal), v...))
	})
}

// IntValGT applies the GT predicate on the "int_val" field.
func IntValGT(v int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIntVal), v))
	})
}

// IntValGTE applies the GTE predicate on the "int_val" field.
func IntValGTE(v int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIntVal), v))
	})
}

// IntValLT applies the LT predicate on the "int_val" field.
func IntValLT(v int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIntVal), v))
	})
}

// IntValLTE applies the LTE predicate on the "int_val" field.
func IntValLTE(v int) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIntVal), v))
	})
}

// IntValIsNil applies the IsNil predicate on the "int_val" field.
func IntValIsNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIntVal)))
	})
}

// IntValNotNil applies the NotNil predicate on the "int_val" field.
func IntValNotNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIntVal)))
	})
}

// BoolValEQ applies the EQ predicate on the "bool_val" field.
func BoolValEQ(v bool) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBoolVal), v))
	})
}

// BoolValNEQ applies the NEQ predicate on the "bool_val" field.
func BoolValNEQ(v bool) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBoolVal), v))
	})
}

// BoolValIsNil applies the IsNil predicate on the "bool_val" field.
func BoolValIsNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldBoolVal)))
	})
}

// BoolValNotNil applies the NotNil predicate on the "bool_val" field.
func BoolValNotNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldBoolVal)))
	})
}

// FloatValEQ applies the EQ predicate on the "float_val" field.
func FloatValEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFloatVal), v))
	})
}

// FloatValNEQ applies the NEQ predicate on the "float_val" field.
func FloatValNEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFloatVal), v))
	})
}

// FloatValIn applies the In predicate on the "float_val" field.
func FloatValIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldFloatVal), v...))
	})
}

// FloatValNotIn applies the NotIn predicate on the "float_val" field.
func FloatValNotIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldFloatVal), v...))
	})
}

// FloatValGT applies the GT predicate on the "float_val" field.
func FloatValGT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldFloatVal), v))
	})
}

// FloatValGTE applies the GTE predicate on the "float_val" field.
func FloatValGTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldFloatVal), v))
	})
}

// FloatValLT applies the LT predicate on the "float_val" field.
func FloatValLT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldFloatVal), v))
	})
}

// FloatValLTE applies the LTE predicate on the "float_val" field.
func FloatValLTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldFloatVal), v))
	})
}

// FloatValIsNil applies the IsNil predicate on the "float_val" field.
func FloatValIsNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldFloatVal)))
	})
}

// FloatValNotNil applies the NotNil predicate on the "float_val" field.
func FloatValNotNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldFloatVal)))
	})
}

// LatitudeValEQ applies the EQ predicate on the "latitude_val" field.
func LatitudeValEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValNEQ applies the NEQ predicate on the "latitude_val" field.
func LatitudeValNEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValIn applies the In predicate on the "latitude_val" field.
func LatitudeValIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLatitudeVal), v...))
	})
}

// LatitudeValNotIn applies the NotIn predicate on the "latitude_val" field.
func LatitudeValNotIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLatitudeVal), v...))
	})
}

// LatitudeValGT applies the GT predicate on the "latitude_val" field.
func LatitudeValGT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValGTE applies the GTE predicate on the "latitude_val" field.
func LatitudeValGTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValLT applies the LT predicate on the "latitude_val" field.
func LatitudeValLT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValLTE applies the LTE predicate on the "latitude_val" field.
func LatitudeValLTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValIsNil applies the IsNil predicate on the "latitude_val" field.
func LatitudeValIsNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLatitudeVal)))
	})
}

// LatitudeValNotNil applies the NotNil predicate on the "latitude_val" field.
func LatitudeValNotNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLatitudeVal)))
	})
}

// LongitudeValEQ applies the EQ predicate on the "longitude_val" field.
func LongitudeValEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValNEQ applies the NEQ predicate on the "longitude_val" field.
func LongitudeValNEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValIn applies the In predicate on the "longitude_val" field.
func LongitudeValIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLongitudeVal), v...))
	})
}

// LongitudeValNotIn applies the NotIn predicate on the "longitude_val" field.
func LongitudeValNotIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLongitudeVal), v...))
	})
}

// LongitudeValGT applies the GT predicate on the "longitude_val" field.
func LongitudeValGT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValGTE applies the GTE predicate on the "longitude_val" field.
func LongitudeValGTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValLT applies the LT predicate on the "longitude_val" field.
func LongitudeValLT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValLTE applies the LTE predicate on the "longitude_val" field.
func LongitudeValLTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValIsNil applies the IsNil predicate on the "longitude_val" field.
func LongitudeValIsNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLongitudeVal)))
	})
}

// LongitudeValNotNil applies the NotNil predicate on the "longitude_val" field.
func LongitudeValNotNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLongitudeVal)))
	})
}

// RangeFromValEQ applies the EQ predicate on the "range_from_val" field.
func RangeFromValEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValNEQ applies the NEQ predicate on the "range_from_val" field.
func RangeFromValNEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValIn applies the In predicate on the "range_from_val" field.
func RangeFromValIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldRangeFromVal), v...))
	})
}

// RangeFromValNotIn applies the NotIn predicate on the "range_from_val" field.
func RangeFromValNotIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldRangeFromVal), v...))
	})
}

// RangeFromValGT applies the GT predicate on the "range_from_val" field.
func RangeFromValGT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValGTE applies the GTE predicate on the "range_from_val" field.
func RangeFromValGTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValLT applies the LT predicate on the "range_from_val" field.
func RangeFromValLT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValLTE applies the LTE predicate on the "range_from_val" field.
func RangeFromValLTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValIsNil applies the IsNil predicate on the "range_from_val" field.
func RangeFromValIsNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldRangeFromVal)))
	})
}

// RangeFromValNotNil applies the NotNil predicate on the "range_from_val" field.
func RangeFromValNotNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldRangeFromVal)))
	})
}

// RangeToValEQ applies the EQ predicate on the "range_to_val" field.
func RangeToValEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRangeToVal), v))
	})
}

// RangeToValNEQ applies the NEQ predicate on the "range_to_val" field.
func RangeToValNEQ(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRangeToVal), v))
	})
}

// RangeToValIn applies the In predicate on the "range_to_val" field.
func RangeToValIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldRangeToVal), v...))
	})
}

// RangeToValNotIn applies the NotIn predicate on the "range_to_val" field.
func RangeToValNotIn(vs ...float64) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldRangeToVal), v...))
	})
}

// RangeToValGT applies the GT predicate on the "range_to_val" field.
func RangeToValGT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRangeToVal), v))
	})
}

// RangeToValGTE applies the GTE predicate on the "range_to_val" field.
func RangeToValGTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRangeToVal), v))
	})
}

// RangeToValLT applies the LT predicate on the "range_to_val" field.
func RangeToValLT(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRangeToVal), v))
	})
}

// RangeToValLTE applies the LTE predicate on the "range_to_val" field.
func RangeToValLTE(v float64) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRangeToVal), v))
	})
}

// RangeToValIsNil applies the IsNil predicate on the "range_to_val" field.
func RangeToValIsNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldRangeToVal)))
	})
}

// RangeToValNotNil applies the NotNil predicate on the "range_to_val" field.
func RangeToValNotNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldRangeToVal)))
	})
}

// StringValEQ applies the EQ predicate on the "string_val" field.
func StringValEQ(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStringVal), v))
	})
}

// StringValNEQ applies the NEQ predicate on the "string_val" field.
func StringValNEQ(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStringVal), v))
	})
}

// StringValIn applies the In predicate on the "string_val" field.
func StringValIn(vs ...string) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStringVal), v...))
	})
}

// StringValNotIn applies the NotIn predicate on the "string_val" field.
func StringValNotIn(vs ...string) predicate.Property {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Property(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStringVal), v...))
	})
}

// StringValGT applies the GT predicate on the "string_val" field.
func StringValGT(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStringVal), v))
	})
}

// StringValGTE applies the GTE predicate on the "string_val" field.
func StringValGTE(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStringVal), v))
	})
}

// StringValLT applies the LT predicate on the "string_val" field.
func StringValLT(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStringVal), v))
	})
}

// StringValLTE applies the LTE predicate on the "string_val" field.
func StringValLTE(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStringVal), v))
	})
}

// StringValContains applies the Contains predicate on the "string_val" field.
func StringValContains(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStringVal), v))
	})
}

// StringValHasPrefix applies the HasPrefix predicate on the "string_val" field.
func StringValHasPrefix(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStringVal), v))
	})
}

// StringValHasSuffix applies the HasSuffix predicate on the "string_val" field.
func StringValHasSuffix(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStringVal), v))
	})
}

// StringValIsNil applies the IsNil predicate on the "string_val" field.
func StringValIsNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldStringVal)))
	})
}

// StringValNotNil applies the NotNil predicate on the "string_val" field.
func StringValNotNil() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldStringVal)))
	})
}

// StringValEqualFold applies the EqualFold predicate on the "string_val" field.
func StringValEqualFold(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStringVal), v))
	})
}

// StringValContainsFold applies the ContainsFold predicate on the "string_val" field.
func StringValContainsFold(v string) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStringVal), v))
	})
}

// HasType applies the HasEdge predicate on the "type" edge.
func HasType() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TypeTable, TypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTypeWith applies the HasEdge predicate on the "type" edge with a given conditions (other predicates).
func HasTypeWith(preds ...predicate.PropertyType) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TypeTable, TypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLocation applies the HasEdge predicate on the "location" edge.
func HasLocation() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LocationTable, LocationColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLocationWith applies the HasEdge predicate on the "location" edge with a given conditions (other predicates).
func HasLocationWith(preds ...predicate.Location) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LocationTable, LocationColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEquipment applies the HasEdge predicate on the "equipment" edge.
func HasEquipment() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentTable, EquipmentColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentWith applies the HasEdge predicate on the "equipment" edge with a given conditions (other predicates).
func HasEquipmentWith(preds ...predicate.Equipment) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentTable, EquipmentColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasService applies the HasEdge predicate on the "service" edge.
func HasService() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServiceTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ServiceTable, ServiceColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceWith applies the HasEdge predicate on the "service" edge with a given conditions (other predicates).
func HasServiceWith(preds ...predicate.Service) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServiceInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ServiceTable, ServiceColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEquipmentPort applies the HasEdge predicate on the "equipment_port" edge.
func HasEquipmentPort() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentPortTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentPortTable, EquipmentPortColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentPortWith applies the HasEdge predicate on the "equipment_port" edge with a given conditions (other predicates).
func HasEquipmentPortWith(preds ...predicate.EquipmentPort) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentPortInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentPortTable, EquipmentPortColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLink applies the HasEdge predicate on the "link" edge.
func HasLink() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LinkTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LinkTable, LinkColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLinkWith applies the HasEdge predicate on the "link" edge with a given conditions (other predicates).
func HasLinkWith(preds ...predicate.Link) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LinkInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LinkTable, LinkColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWorkOrder applies the HasEdge predicate on the "work_order" edge.
func HasWorkOrder() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTable, WorkOrderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWorkOrderWith applies the HasEdge predicate on the "work_order" edge with a given conditions (other predicates).
func HasWorkOrderWith(preds ...predicate.WorkOrder) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTable, WorkOrderColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEquipmentValue applies the HasEdge predicate on the "equipment_value" edge.
func HasEquipmentValue() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentValueTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, EquipmentValueTable, EquipmentValueColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentValueWith applies the HasEdge predicate on the "equipment_value" edge with a given conditions (other predicates).
func HasEquipmentValueWith(preds ...predicate.Equipment) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentValueInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, EquipmentValueTable, EquipmentValueColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLocationValue applies the HasEdge predicate on the "location_value" edge.
func HasLocationValue() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationValueTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LocationValueTable, LocationValueColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLocationValueWith applies the HasEdge predicate on the "location_value" edge with a given conditions (other predicates).
func HasLocationValueWith(preds ...predicate.Location) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationValueInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LocationValueTable, LocationValueColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasServiceValue applies the HasEdge predicate on the "service_value" edge.
func HasServiceValue() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServiceValueTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ServiceValueTable, ServiceValueColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceValueWith applies the HasEdge predicate on the "service_value" edge with a given conditions (other predicates).
func HasServiceValueWith(preds ...predicate.Service) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServiceValueInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ServiceValueTable, ServiceValueColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWorkOrderValue applies the HasEdge predicate on the "work_order_value" edge.
func HasWorkOrderValue() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderValueTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, WorkOrderValueTable, WorkOrderValueColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWorkOrderValueWith applies the HasEdge predicate on the "work_order_value" edge with a given conditions (other predicates).
func HasWorkOrderValueWith(preds ...predicate.WorkOrder) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderValueInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, WorkOrderValueTable, WorkOrderValueColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasUserValue applies the HasEdge predicate on the "user_value" edge.
func HasUserValue() predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserValueTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, UserValueTable, UserValueColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserValueWith applies the HasEdge predicate on the "user_value" edge with a given conditions (other predicates).
func HasUserValueWith(preds ...predicate.User) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserValueInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, UserValueTable, UserValueColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Property) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Property) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
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
func Not(p predicate.Property) predicate.Property {
	return predicate.Property(func(s *sql.Selector) {
		p(s.Not())
	})
}
