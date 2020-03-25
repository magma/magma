// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentportdefinition

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func IDGT(id int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// Bandwidth applies equality check predicate on the "bandwidth" field. It's identical to BandwidthEQ.
func Bandwidth(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBandwidth), v))
	})
}

// VisibilityLabel applies equality check predicate on the "visibility_label" field. It's identical to VisibilityLabelEQ.
func VisibilityLabel(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVisibilityLabel), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func NameGT(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	})
}

// IndexIn applies the In predicate on the "index" field.
func IndexIn(vs ...int) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func IndexNotIn(vs ...int) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func IndexGT(v int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	})
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	})
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	})
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	})
}

// IndexIsNil applies the IsNil predicate on the "index" field.
func IndexIsNil() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIndex)))
	})
}

// IndexNotNil applies the NotNil predicate on the "index" field.
func IndexNotNil() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIndex)))
	})
}

// BandwidthEQ applies the EQ predicate on the "bandwidth" field.
func BandwidthEQ(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBandwidth), v))
	})
}

// BandwidthNEQ applies the NEQ predicate on the "bandwidth" field.
func BandwidthNEQ(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBandwidth), v))
	})
}

// BandwidthIn applies the In predicate on the "bandwidth" field.
func BandwidthIn(vs ...string) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldBandwidth), v...))
	})
}

// BandwidthNotIn applies the NotIn predicate on the "bandwidth" field.
func BandwidthNotIn(vs ...string) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldBandwidth), v...))
	})
}

// BandwidthGT applies the GT predicate on the "bandwidth" field.
func BandwidthGT(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldBandwidth), v))
	})
}

// BandwidthGTE applies the GTE predicate on the "bandwidth" field.
func BandwidthGTE(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldBandwidth), v))
	})
}

// BandwidthLT applies the LT predicate on the "bandwidth" field.
func BandwidthLT(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldBandwidth), v))
	})
}

// BandwidthLTE applies the LTE predicate on the "bandwidth" field.
func BandwidthLTE(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldBandwidth), v))
	})
}

// BandwidthContains applies the Contains predicate on the "bandwidth" field.
func BandwidthContains(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldBandwidth), v))
	})
}

// BandwidthHasPrefix applies the HasPrefix predicate on the "bandwidth" field.
func BandwidthHasPrefix(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldBandwidth), v))
	})
}

// BandwidthHasSuffix applies the HasSuffix predicate on the "bandwidth" field.
func BandwidthHasSuffix(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldBandwidth), v))
	})
}

// BandwidthIsNil applies the IsNil predicate on the "bandwidth" field.
func BandwidthIsNil() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldBandwidth)))
	})
}

// BandwidthNotNil applies the NotNil predicate on the "bandwidth" field.
func BandwidthNotNil() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldBandwidth)))
	})
}

// BandwidthEqualFold applies the EqualFold predicate on the "bandwidth" field.
func BandwidthEqualFold(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldBandwidth), v))
	})
}

// BandwidthContainsFold applies the ContainsFold predicate on the "bandwidth" field.
func BandwidthContainsFold(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldBandwidth), v))
	})
}

// VisibilityLabelEQ applies the EQ predicate on the "visibility_label" field.
func VisibilityLabelEQ(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelNEQ applies the NEQ predicate on the "visibility_label" field.
func VisibilityLabelNEQ(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelIn applies the In predicate on the "visibility_label" field.
func VisibilityLabelIn(vs ...string) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldVisibilityLabel), v...))
	})
}

// VisibilityLabelNotIn applies the NotIn predicate on the "visibility_label" field.
func VisibilityLabelNotIn(vs ...string) predicate.EquipmentPortDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldVisibilityLabel), v...))
	})
}

// VisibilityLabelGT applies the GT predicate on the "visibility_label" field.
func VisibilityLabelGT(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelGTE applies the GTE predicate on the "visibility_label" field.
func VisibilityLabelGTE(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelLT applies the LT predicate on the "visibility_label" field.
func VisibilityLabelLT(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelLTE applies the LTE predicate on the "visibility_label" field.
func VisibilityLabelLTE(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelContains applies the Contains predicate on the "visibility_label" field.
func VisibilityLabelContains(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelHasPrefix applies the HasPrefix predicate on the "visibility_label" field.
func VisibilityLabelHasPrefix(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelHasSuffix applies the HasSuffix predicate on the "visibility_label" field.
func VisibilityLabelHasSuffix(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelIsNil applies the IsNil predicate on the "visibility_label" field.
func VisibilityLabelIsNil() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldVisibilityLabel)))
	})
}

// VisibilityLabelNotNil applies the NotNil predicate on the "visibility_label" field.
func VisibilityLabelNotNil() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldVisibilityLabel)))
	})
}

// VisibilityLabelEqualFold applies the EqualFold predicate on the "visibility_label" field.
func VisibilityLabelEqualFold(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelContainsFold applies the ContainsFold predicate on the "visibility_label" field.
func VisibilityLabelContainsFold(v string) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldVisibilityLabel), v))
	})
}

// HasEquipmentPortType applies the HasEdge predicate on the "equipment_port_type" edge.
func HasEquipmentPortType() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentPortTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, EquipmentPortTypeTable, EquipmentPortTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentPortTypeWith applies the HasEdge predicate on the "equipment_port_type" edge with a given conditions (other predicates).
func HasEquipmentPortTypeWith(preds ...predicate.EquipmentPortType) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentPortTypeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, EquipmentPortTypeTable, EquipmentPortTypeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPorts applies the HasEdge predicate on the "ports" edge.
func HasPorts() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PortsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, PortsTable, PortsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPortsWith applies the HasEdge predicate on the "ports" edge with a given conditions (other predicates).
func HasPortsWith(preds ...predicate.EquipmentPort) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PortsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, PortsTable, PortsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEquipmentType applies the HasEdge predicate on the "equipment_type" edge.
func HasEquipmentType() predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentTypeTable, EquipmentTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentTypeWith applies the HasEdge predicate on the "equipment_type" edge with a given conditions (other predicates).
func HasEquipmentTypeWith(preds ...predicate.EquipmentType) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.EquipmentPortDefinition) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.EquipmentPortDefinition) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
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
func Not(p predicate.EquipmentPortDefinition) predicate.EquipmentPortDefinition {
	return predicate.EquipmentPortDefinition(func(s *sql.Selector) {
		p(s.Not())
	})
}
