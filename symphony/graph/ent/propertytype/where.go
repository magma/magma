// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package propertytype

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
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
func IDGT(id int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldType), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// Category applies equality check predicate on the "category" field. It's identical to CategoryEQ.
func Category(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCategory), v))
	})
}

// IntVal applies equality check predicate on the "int_val" field. It's identical to IntValEQ.
func IntVal(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIntVal), v))
	})
}

// BoolVal applies equality check predicate on the "bool_val" field. It's identical to BoolValEQ.
func BoolVal(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBoolVal), v))
	})
}

// FloatVal applies equality check predicate on the "float_val" field. It's identical to FloatValEQ.
func FloatVal(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFloatVal), v))
	})
}

// LatitudeVal applies equality check predicate on the "latitude_val" field. It's identical to LatitudeValEQ.
func LatitudeVal(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitudeVal), v))
	})
}

// LongitudeVal applies equality check predicate on the "longitude_val" field. It's identical to LongitudeValEQ.
func LongitudeVal(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitudeVal), v))
	})
}

// StringVal applies equality check predicate on the "string_val" field. It's identical to StringValEQ.
func StringVal(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStringVal), v))
	})
}

// RangeFromVal applies equality check predicate on the "range_from_val" field. It's identical to RangeFromValEQ.
func RangeFromVal(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRangeFromVal), v))
	})
}

// RangeToVal applies equality check predicate on the "range_to_val" field. It's identical to RangeToValEQ.
func RangeToVal(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRangeToVal), v))
	})
}

// IsInstanceProperty applies equality check predicate on the "is_instance_property" field. It's identical to IsInstancePropertyEQ.
func IsInstanceProperty(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIsInstanceProperty), v))
	})
}

// Editable applies equality check predicate on the "editable" field. It's identical to EditableEQ.
func Editable(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEditable), v))
	})
}

// Mandatory applies equality check predicate on the "mandatory" field. It's identical to MandatoryEQ.
func Mandatory(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMandatory), v))
	})
}

// Deleted applies equality check predicate on the "deleted" field. It's identical to DeletedEQ.
func Deleted(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeleted), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldType), v))
	})
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldType), v))
	})
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldType), v...))
	})
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldType), v...))
	})
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldType), v))
	})
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldType), v))
	})
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldType), v))
	})
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldType), v))
	})
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldType), v))
	})
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldType), v))
	})
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldType), v))
	})
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldType), v))
	})
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldType), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	})
}

// IndexIn applies the In predicate on the "index" field.
func IndexIn(vs ...int) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldIndex), v...))
	})
}

// IndexNotIn applies the NotIn predicate on the "index" field.
func IndexNotIn(vs ...int) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldIndex), v...))
	})
}

// IndexGT applies the GT predicate on the "index" field.
func IndexGT(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	})
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	})
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	})
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	})
}

// IndexIsNil applies the IsNil predicate on the "index" field.
func IndexIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIndex)))
	})
}

// IndexNotNil applies the NotNil predicate on the "index" field.
func IndexNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIndex)))
	})
}

// CategoryEQ applies the EQ predicate on the "category" field.
func CategoryEQ(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCategory), v))
	})
}

// CategoryNEQ applies the NEQ predicate on the "category" field.
func CategoryNEQ(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCategory), v))
	})
}

// CategoryIn applies the In predicate on the "category" field.
func CategoryIn(vs ...string) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCategory), v...))
	})
}

// CategoryNotIn applies the NotIn predicate on the "category" field.
func CategoryNotIn(vs ...string) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCategory), v...))
	})
}

// CategoryGT applies the GT predicate on the "category" field.
func CategoryGT(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCategory), v))
	})
}

// CategoryGTE applies the GTE predicate on the "category" field.
func CategoryGTE(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCategory), v))
	})
}

// CategoryLT applies the LT predicate on the "category" field.
func CategoryLT(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCategory), v))
	})
}

// CategoryLTE applies the LTE predicate on the "category" field.
func CategoryLTE(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCategory), v))
	})
}

// CategoryContains applies the Contains predicate on the "category" field.
func CategoryContains(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCategory), v))
	})
}

// CategoryHasPrefix applies the HasPrefix predicate on the "category" field.
func CategoryHasPrefix(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCategory), v))
	})
}

// CategoryHasSuffix applies the HasSuffix predicate on the "category" field.
func CategoryHasSuffix(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCategory), v))
	})
}

// CategoryIsNil applies the IsNil predicate on the "category" field.
func CategoryIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldCategory)))
	})
}

// CategoryNotNil applies the NotNil predicate on the "category" field.
func CategoryNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldCategory)))
	})
}

// CategoryEqualFold applies the EqualFold predicate on the "category" field.
func CategoryEqualFold(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCategory), v))
	})
}

// CategoryContainsFold applies the ContainsFold predicate on the "category" field.
func CategoryContainsFold(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCategory), v))
	})
}

// IntValEQ applies the EQ predicate on the "int_val" field.
func IntValEQ(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIntVal), v))
	})
}

// IntValNEQ applies the NEQ predicate on the "int_val" field.
func IntValNEQ(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIntVal), v))
	})
}

// IntValIn applies the In predicate on the "int_val" field.
func IntValIn(vs ...int) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func IntValNotIn(vs ...int) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func IntValGT(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIntVal), v))
	})
}

// IntValGTE applies the GTE predicate on the "int_val" field.
func IntValGTE(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIntVal), v))
	})
}

// IntValLT applies the LT predicate on the "int_val" field.
func IntValLT(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIntVal), v))
	})
}

// IntValLTE applies the LTE predicate on the "int_val" field.
func IntValLTE(v int) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIntVal), v))
	})
}

// IntValIsNil applies the IsNil predicate on the "int_val" field.
func IntValIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIntVal)))
	})
}

// IntValNotNil applies the NotNil predicate on the "int_val" field.
func IntValNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIntVal)))
	})
}

// BoolValEQ applies the EQ predicate on the "bool_val" field.
func BoolValEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBoolVal), v))
	})
}

// BoolValNEQ applies the NEQ predicate on the "bool_val" field.
func BoolValNEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBoolVal), v))
	})
}

// BoolValIsNil applies the IsNil predicate on the "bool_val" field.
func BoolValIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldBoolVal)))
	})
}

// BoolValNotNil applies the NotNil predicate on the "bool_val" field.
func BoolValNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldBoolVal)))
	})
}

// FloatValEQ applies the EQ predicate on the "float_val" field.
func FloatValEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFloatVal), v))
	})
}

// FloatValNEQ applies the NEQ predicate on the "float_val" field.
func FloatValNEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFloatVal), v))
	})
}

// FloatValIn applies the In predicate on the "float_val" field.
func FloatValIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func FloatValNotIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func FloatValGT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldFloatVal), v))
	})
}

// FloatValGTE applies the GTE predicate on the "float_val" field.
func FloatValGTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldFloatVal), v))
	})
}

// FloatValLT applies the LT predicate on the "float_val" field.
func FloatValLT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldFloatVal), v))
	})
}

// FloatValLTE applies the LTE predicate on the "float_val" field.
func FloatValLTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldFloatVal), v))
	})
}

// FloatValIsNil applies the IsNil predicate on the "float_val" field.
func FloatValIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldFloatVal)))
	})
}

// FloatValNotNil applies the NotNil predicate on the "float_val" field.
func FloatValNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldFloatVal)))
	})
}

// LatitudeValEQ applies the EQ predicate on the "latitude_val" field.
func LatitudeValEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValNEQ applies the NEQ predicate on the "latitude_val" field.
func LatitudeValNEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValIn applies the In predicate on the "latitude_val" field.
func LatitudeValIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func LatitudeValNotIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func LatitudeValGT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValGTE applies the GTE predicate on the "latitude_val" field.
func LatitudeValGTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValLT applies the LT predicate on the "latitude_val" field.
func LatitudeValLT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValLTE applies the LTE predicate on the "latitude_val" field.
func LatitudeValLTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLatitudeVal), v))
	})
}

// LatitudeValIsNil applies the IsNil predicate on the "latitude_val" field.
func LatitudeValIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLatitudeVal)))
	})
}

// LatitudeValNotNil applies the NotNil predicate on the "latitude_val" field.
func LatitudeValNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLatitudeVal)))
	})
}

// LongitudeValEQ applies the EQ predicate on the "longitude_val" field.
func LongitudeValEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValNEQ applies the NEQ predicate on the "longitude_val" field.
func LongitudeValNEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValIn applies the In predicate on the "longitude_val" field.
func LongitudeValIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func LongitudeValNotIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func LongitudeValGT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValGTE applies the GTE predicate on the "longitude_val" field.
func LongitudeValGTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValLT applies the LT predicate on the "longitude_val" field.
func LongitudeValLT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValLTE applies the LTE predicate on the "longitude_val" field.
func LongitudeValLTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLongitudeVal), v))
	})
}

// LongitudeValIsNil applies the IsNil predicate on the "longitude_val" field.
func LongitudeValIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLongitudeVal)))
	})
}

// LongitudeValNotNil applies the NotNil predicate on the "longitude_val" field.
func LongitudeValNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLongitudeVal)))
	})
}

// StringValEQ applies the EQ predicate on the "string_val" field.
func StringValEQ(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStringVal), v))
	})
}

// StringValNEQ applies the NEQ predicate on the "string_val" field.
func StringValNEQ(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStringVal), v))
	})
}

// StringValIn applies the In predicate on the "string_val" field.
func StringValIn(vs ...string) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func StringValNotIn(vs ...string) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func StringValGT(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStringVal), v))
	})
}

// StringValGTE applies the GTE predicate on the "string_val" field.
func StringValGTE(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStringVal), v))
	})
}

// StringValLT applies the LT predicate on the "string_val" field.
func StringValLT(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStringVal), v))
	})
}

// StringValLTE applies the LTE predicate on the "string_val" field.
func StringValLTE(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStringVal), v))
	})
}

// StringValContains applies the Contains predicate on the "string_val" field.
func StringValContains(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStringVal), v))
	})
}

// StringValHasPrefix applies the HasPrefix predicate on the "string_val" field.
func StringValHasPrefix(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStringVal), v))
	})
}

// StringValHasSuffix applies the HasSuffix predicate on the "string_val" field.
func StringValHasSuffix(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStringVal), v))
	})
}

// StringValIsNil applies the IsNil predicate on the "string_val" field.
func StringValIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldStringVal)))
	})
}

// StringValNotNil applies the NotNil predicate on the "string_val" field.
func StringValNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldStringVal)))
	})
}

// StringValEqualFold applies the EqualFold predicate on the "string_val" field.
func StringValEqualFold(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStringVal), v))
	})
}

// StringValContainsFold applies the ContainsFold predicate on the "string_val" field.
func StringValContainsFold(v string) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStringVal), v))
	})
}

// RangeFromValEQ applies the EQ predicate on the "range_from_val" field.
func RangeFromValEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValNEQ applies the NEQ predicate on the "range_from_val" field.
func RangeFromValNEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValIn applies the In predicate on the "range_from_val" field.
func RangeFromValIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func RangeFromValNotIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func RangeFromValGT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValGTE applies the GTE predicate on the "range_from_val" field.
func RangeFromValGTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValLT applies the LT predicate on the "range_from_val" field.
func RangeFromValLT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValLTE applies the LTE predicate on the "range_from_val" field.
func RangeFromValLTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRangeFromVal), v))
	})
}

// RangeFromValIsNil applies the IsNil predicate on the "range_from_val" field.
func RangeFromValIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldRangeFromVal)))
	})
}

// RangeFromValNotNil applies the NotNil predicate on the "range_from_val" field.
func RangeFromValNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldRangeFromVal)))
	})
}

// RangeToValEQ applies the EQ predicate on the "range_to_val" field.
func RangeToValEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRangeToVal), v))
	})
}

// RangeToValNEQ applies the NEQ predicate on the "range_to_val" field.
func RangeToValNEQ(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRangeToVal), v))
	})
}

// RangeToValIn applies the In predicate on the "range_to_val" field.
func RangeToValIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func RangeToValNotIn(vs ...float64) predicate.PropertyType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.PropertyType(func(s *sql.Selector) {
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
func RangeToValGT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRangeToVal), v))
	})
}

// RangeToValGTE applies the GTE predicate on the "range_to_val" field.
func RangeToValGTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRangeToVal), v))
	})
}

// RangeToValLT applies the LT predicate on the "range_to_val" field.
func RangeToValLT(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRangeToVal), v))
	})
}

// RangeToValLTE applies the LTE predicate on the "range_to_val" field.
func RangeToValLTE(v float64) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRangeToVal), v))
	})
}

// RangeToValIsNil applies the IsNil predicate on the "range_to_val" field.
func RangeToValIsNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldRangeToVal)))
	})
}

// RangeToValNotNil applies the NotNil predicate on the "range_to_val" field.
func RangeToValNotNil() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldRangeToVal)))
	})
}

// IsInstancePropertyEQ applies the EQ predicate on the "is_instance_property" field.
func IsInstancePropertyEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIsInstanceProperty), v))
	})
}

// IsInstancePropertyNEQ applies the NEQ predicate on the "is_instance_property" field.
func IsInstancePropertyNEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIsInstanceProperty), v))
	})
}

// EditableEQ applies the EQ predicate on the "editable" field.
func EditableEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEditable), v))
	})
}

// EditableNEQ applies the NEQ predicate on the "editable" field.
func EditableNEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEditable), v))
	})
}

// MandatoryEQ applies the EQ predicate on the "mandatory" field.
func MandatoryEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMandatory), v))
	})
}

// MandatoryNEQ applies the NEQ predicate on the "mandatory" field.
func MandatoryNEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldMandatory), v))
	})
}

// DeletedEQ applies the EQ predicate on the "deleted" field.
func DeletedEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeleted), v))
	})
}

// DeletedNEQ applies the NEQ predicate on the "deleted" field.
func DeletedNEQ(v bool) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDeleted), v))
	})
}

// HasProperties applies the HasEdge predicate on the "properties" edge.
func HasProperties() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertiesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, PropertiesTable, PropertiesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPropertiesWith applies the HasEdge predicate on the "properties" edge with a given conditions (other predicates).
func HasPropertiesWith(preds ...predicate.Property) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertiesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, PropertiesTable, PropertiesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLocationType applies the HasEdge predicate on the "location_type" edge.
func HasLocationType() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LocationTypeTable, LocationTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLocationTypeWith applies the HasEdge predicate on the "location_type" edge with a given conditions (other predicates).
func HasLocationTypeWith(preds ...predicate.LocationType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LocationTypeTable, LocationTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEquipmentPortType applies the HasEdge predicate on the "equipment_port_type" edge.
func HasEquipmentPortType() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentPortTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentPortTypeTable, EquipmentPortTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentPortTypeWith applies the HasEdge predicate on the "equipment_port_type" edge with a given conditions (other predicates).
func HasEquipmentPortTypeWith(preds ...predicate.EquipmentPortType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentPortTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentPortTypeTable, EquipmentPortTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLinkEquipmentPortType applies the HasEdge predicate on the "link_equipment_port_type" edge.
func HasLinkEquipmentPortType() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LinkEquipmentPortTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LinkEquipmentPortTypeTable, LinkEquipmentPortTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLinkEquipmentPortTypeWith applies the HasEdge predicate on the "link_equipment_port_type" edge with a given conditions (other predicates).
func HasLinkEquipmentPortTypeWith(preds ...predicate.EquipmentPortType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LinkEquipmentPortTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LinkEquipmentPortTypeTable, LinkEquipmentPortTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEquipmentType applies the HasEdge predicate on the "equipment_type" edge.
func HasEquipmentType() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentTypeTable, EquipmentTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentTypeWith applies the HasEdge predicate on the "equipment_type" edge with a given conditions (other predicates).
func HasEquipmentTypeWith(preds ...predicate.EquipmentType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentTypeTable, EquipmentTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasServiceType applies the HasEdge predicate on the "service_type" edge.
func HasServiceType() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServiceTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ServiceTypeTable, ServiceTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceTypeWith applies the HasEdge predicate on the "service_type" edge with a given conditions (other predicates).
func HasServiceTypeWith(preds ...predicate.ServiceType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ServiceTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ServiceTypeTable, ServiceTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWorkOrderType applies the HasEdge predicate on the "work_order_type" edge.
func HasWorkOrderType() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTypeTable, WorkOrderTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWorkOrderTypeWith applies the HasEdge predicate on the "work_order_type" edge with a given conditions (other predicates).
func HasWorkOrderTypeWith(preds ...predicate.WorkOrderType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkOrderTypeTable, WorkOrderTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProjectType applies the HasEdge predicate on the "project_type" edge.
func HasProjectType() predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTypeTable, ProjectTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectTypeWith applies the HasEdge predicate on the "project_type" edge with a given conditions (other predicates).
func HasProjectTypeWith(preds ...predicate.ProjectType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTypeTable, ProjectTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.PropertyType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.PropertyType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
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
func Not(p predicate.PropertyType) predicate.PropertyType {
	return predicate.PropertyType(func(s *sql.Selector) {
		p(s.Not())
	})
}
