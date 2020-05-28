// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package locationtype

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
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
func IDGT(id int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Site applies equality check predicate on the "site" field. It's identical to SiteEQ.
func Site(v bool) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSite), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// MapType applies equality check predicate on the "map_type" field. It's identical to MapTypeEQ.
func MapType(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMapType), v))
	})
}

// MapZoomLevel applies equality check predicate on the "map_zoom_level" field. It's identical to MapZoomLevelEQ.
func MapZoomLevel(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMapZoomLevel), v))
	})
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// SiteEQ applies the EQ predicate on the "site" field.
func SiteEQ(v bool) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSite), v))
	})
}

// SiteNEQ applies the NEQ predicate on the "site" field.
func SiteNEQ(v bool) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSite), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
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
func NameGT(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// MapTypeEQ applies the EQ predicate on the "map_type" field.
func MapTypeEQ(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMapType), v))
	})
}

// MapTypeNEQ applies the NEQ predicate on the "map_type" field.
func MapTypeNEQ(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldMapType), v))
	})
}

// MapTypeIn applies the In predicate on the "map_type" field.
func MapTypeIn(vs ...string) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldMapType), v...))
	})
}

// MapTypeNotIn applies the NotIn predicate on the "map_type" field.
func MapTypeNotIn(vs ...string) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldMapType), v...))
	})
}

// MapTypeGT applies the GT predicate on the "map_type" field.
func MapTypeGT(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldMapType), v))
	})
}

// MapTypeGTE applies the GTE predicate on the "map_type" field.
func MapTypeGTE(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldMapType), v))
	})
}

// MapTypeLT applies the LT predicate on the "map_type" field.
func MapTypeLT(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldMapType), v))
	})
}

// MapTypeLTE applies the LTE predicate on the "map_type" field.
func MapTypeLTE(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldMapType), v))
	})
}

// MapTypeContains applies the Contains predicate on the "map_type" field.
func MapTypeContains(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldMapType), v))
	})
}

// MapTypeHasPrefix applies the HasPrefix predicate on the "map_type" field.
func MapTypeHasPrefix(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldMapType), v))
	})
}

// MapTypeHasSuffix applies the HasSuffix predicate on the "map_type" field.
func MapTypeHasSuffix(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldMapType), v))
	})
}

// MapTypeIsNil applies the IsNil predicate on the "map_type" field.
func MapTypeIsNil() predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldMapType)))
	})
}

// MapTypeNotNil applies the NotNil predicate on the "map_type" field.
func MapTypeNotNil() predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldMapType)))
	})
}

// MapTypeEqualFold applies the EqualFold predicate on the "map_type" field.
func MapTypeEqualFold(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldMapType), v))
	})
}

// MapTypeContainsFold applies the ContainsFold predicate on the "map_type" field.
func MapTypeContainsFold(v string) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldMapType), v))
	})
}

// MapZoomLevelEQ applies the EQ predicate on the "map_zoom_level" field.
func MapZoomLevelEQ(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMapZoomLevel), v))
	})
}

// MapZoomLevelNEQ applies the NEQ predicate on the "map_zoom_level" field.
func MapZoomLevelNEQ(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldMapZoomLevel), v))
	})
}

// MapZoomLevelIn applies the In predicate on the "map_zoom_level" field.
func MapZoomLevelIn(vs ...int) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldMapZoomLevel), v...))
	})
}

// MapZoomLevelNotIn applies the NotIn predicate on the "map_zoom_level" field.
func MapZoomLevelNotIn(vs ...int) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldMapZoomLevel), v...))
	})
}

// MapZoomLevelGT applies the GT predicate on the "map_zoom_level" field.
func MapZoomLevelGT(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldMapZoomLevel), v))
	})
}

// MapZoomLevelGTE applies the GTE predicate on the "map_zoom_level" field.
func MapZoomLevelGTE(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldMapZoomLevel), v))
	})
}

// MapZoomLevelLT applies the LT predicate on the "map_zoom_level" field.
func MapZoomLevelLT(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldMapZoomLevel), v))
	})
}

// MapZoomLevelLTE applies the LTE predicate on the "map_zoom_level" field.
func MapZoomLevelLTE(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldMapZoomLevel), v))
	})
}

// MapZoomLevelIsNil applies the IsNil predicate on the "map_zoom_level" field.
func MapZoomLevelIsNil() predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldMapZoomLevel)))
	})
}

// MapZoomLevelNotNil applies the NotNil predicate on the "map_zoom_level" field.
func MapZoomLevelNotNil() predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldMapZoomLevel)))
	})
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	})
}

// IndexIn applies the In predicate on the "index" field.
func IndexIn(vs ...int) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
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
func IndexNotIn(vs ...int) predicate.LocationType {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LocationType(func(s *sql.Selector) {
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
func IndexGT(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	})
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	})
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	})
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	})
}

// HasLocations applies the HasEdge predicate on the "locations" edge.
func HasLocations() predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, LocationsTable, LocationsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLocationsWith applies the HasEdge predicate on the "locations" edge with a given conditions (other predicates).
func HasLocationsWith(preds ...predicate.Location) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, LocationsTable, LocationsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPropertyTypes applies the HasEdge predicate on the "property_types" edge.
func HasPropertyTypes() predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertyTypesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertyTypesTable, PropertyTypesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPropertyTypesWith applies the HasEdge predicate on the "property_types" edge with a given conditions (other predicates).
func HasPropertyTypesWith(preds ...predicate.PropertyType) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertyTypesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertyTypesTable, PropertyTypesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasSurveyTemplateCategories applies the HasEdge predicate on the "survey_template_categories" edge.
func HasSurveyTemplateCategories() predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyTemplateCategoriesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SurveyTemplateCategoriesTable, SurveyTemplateCategoriesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSurveyTemplateCategoriesWith applies the HasEdge predicate on the "survey_template_categories" edge with a given conditions (other predicates).
func HasSurveyTemplateCategoriesWith(preds ...predicate.SurveyTemplateCategory) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyTemplateCategoriesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SurveyTemplateCategoriesTable, SurveyTemplateCategoriesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.LocationType) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.LocationType) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
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
func Not(p predicate.LocationType) predicate.LocationType {
	return predicate.LocationType(func(s *sql.Selector) {
		p(s.Not())
	})
}
