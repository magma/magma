// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipment

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
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
func IDGT(id int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// FutureState applies equality check predicate on the "future_state" field. It's identical to FutureStateEQ.
func FutureState(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFutureState), v))
	})
}

// DeviceID applies equality check predicate on the "device_id" field. It's identical to DeviceIDEQ.
func DeviceID(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeviceID), v))
	})
}

// ExternalID applies equality check predicate on the "external_id" field. It's identical to ExternalIDEQ.
func ExternalID(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExternalID), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
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
func NameGT(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// FutureStateEQ applies the EQ predicate on the "future_state" field.
func FutureStateEQ(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFutureState), v))
	})
}

// FutureStateNEQ applies the NEQ predicate on the "future_state" field.
func FutureStateNEQ(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFutureState), v))
	})
}

// FutureStateIn applies the In predicate on the "future_state" field.
func FutureStateIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldFutureState), v...))
	})
}

// FutureStateNotIn applies the NotIn predicate on the "future_state" field.
func FutureStateNotIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldFutureState), v...))
	})
}

// FutureStateGT applies the GT predicate on the "future_state" field.
func FutureStateGT(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldFutureState), v))
	})
}

// FutureStateGTE applies the GTE predicate on the "future_state" field.
func FutureStateGTE(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldFutureState), v))
	})
}

// FutureStateLT applies the LT predicate on the "future_state" field.
func FutureStateLT(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldFutureState), v))
	})
}

// FutureStateLTE applies the LTE predicate on the "future_state" field.
func FutureStateLTE(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldFutureState), v))
	})
}

// FutureStateContains applies the Contains predicate on the "future_state" field.
func FutureStateContains(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldFutureState), v))
	})
}

// FutureStateHasPrefix applies the HasPrefix predicate on the "future_state" field.
func FutureStateHasPrefix(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldFutureState), v))
	})
}

// FutureStateHasSuffix applies the HasSuffix predicate on the "future_state" field.
func FutureStateHasSuffix(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldFutureState), v))
	})
}

// FutureStateIsNil applies the IsNil predicate on the "future_state" field.
func FutureStateIsNil() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldFutureState)))
	})
}

// FutureStateNotNil applies the NotNil predicate on the "future_state" field.
func FutureStateNotNil() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldFutureState)))
	})
}

// FutureStateEqualFold applies the EqualFold predicate on the "future_state" field.
func FutureStateEqualFold(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldFutureState), v))
	})
}

// FutureStateContainsFold applies the ContainsFold predicate on the "future_state" field.
func FutureStateContainsFold(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldFutureState), v))
	})
}

// DeviceIDEQ applies the EQ predicate on the "device_id" field.
func DeviceIDEQ(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeviceID), v))
	})
}

// DeviceIDNEQ applies the NEQ predicate on the "device_id" field.
func DeviceIDNEQ(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDeviceID), v))
	})
}

// DeviceIDIn applies the In predicate on the "device_id" field.
func DeviceIDIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldDeviceID), v...))
	})
}

// DeviceIDNotIn applies the NotIn predicate on the "device_id" field.
func DeviceIDNotIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldDeviceID), v...))
	})
}

// DeviceIDGT applies the GT predicate on the "device_id" field.
func DeviceIDGT(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDeviceID), v))
	})
}

// DeviceIDGTE applies the GTE predicate on the "device_id" field.
func DeviceIDGTE(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDeviceID), v))
	})
}

// DeviceIDLT applies the LT predicate on the "device_id" field.
func DeviceIDLT(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDeviceID), v))
	})
}

// DeviceIDLTE applies the LTE predicate on the "device_id" field.
func DeviceIDLTE(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDeviceID), v))
	})
}

// DeviceIDContains applies the Contains predicate on the "device_id" field.
func DeviceIDContains(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldDeviceID), v))
	})
}

// DeviceIDHasPrefix applies the HasPrefix predicate on the "device_id" field.
func DeviceIDHasPrefix(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldDeviceID), v))
	})
}

// DeviceIDHasSuffix applies the HasSuffix predicate on the "device_id" field.
func DeviceIDHasSuffix(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldDeviceID), v))
	})
}

// DeviceIDIsNil applies the IsNil predicate on the "device_id" field.
func DeviceIDIsNil() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDeviceID)))
	})
}

// DeviceIDNotNil applies the NotNil predicate on the "device_id" field.
func DeviceIDNotNil() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDeviceID)))
	})
}

// DeviceIDEqualFold applies the EqualFold predicate on the "device_id" field.
func DeviceIDEqualFold(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldDeviceID), v))
	})
}

// DeviceIDContainsFold applies the ContainsFold predicate on the "device_id" field.
func DeviceIDContainsFold(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldDeviceID), v))
	})
}

// ExternalIDEQ applies the EQ predicate on the "external_id" field.
func ExternalIDEQ(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExternalID), v))
	})
}

// ExternalIDNEQ applies the NEQ predicate on the "external_id" field.
func ExternalIDNEQ(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldExternalID), v))
	})
}

// ExternalIDIn applies the In predicate on the "external_id" field.
func ExternalIDIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldExternalID), v...))
	})
}

// ExternalIDNotIn applies the NotIn predicate on the "external_id" field.
func ExternalIDNotIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldExternalID), v...))
	})
}

// ExternalIDGT applies the GT predicate on the "external_id" field.
func ExternalIDGT(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldExternalID), v))
	})
}

// ExternalIDGTE applies the GTE predicate on the "external_id" field.
func ExternalIDGTE(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldExternalID), v))
	})
}

// ExternalIDLT applies the LT predicate on the "external_id" field.
func ExternalIDLT(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldExternalID), v))
	})
}

// ExternalIDLTE applies the LTE predicate on the "external_id" field.
func ExternalIDLTE(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldExternalID), v))
	})
}

// ExternalIDContains applies the Contains predicate on the "external_id" field.
func ExternalIDContains(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldExternalID), v))
	})
}

// ExternalIDHasPrefix applies the HasPrefix predicate on the "external_id" field.
func ExternalIDHasPrefix(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldExternalID), v))
	})
}

// ExternalIDHasSuffix applies the HasSuffix predicate on the "external_id" field.
func ExternalIDHasSuffix(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldExternalID), v))
	})
}

// ExternalIDIsNil applies the IsNil predicate on the "external_id" field.
func ExternalIDIsNil() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldExternalID)))
	})
}

// ExternalIDNotNil applies the NotNil predicate on the "external_id" field.
func ExternalIDNotNil() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldExternalID)))
	})
}

// ExternalIDEqualFold applies the EqualFold predicate on the "external_id" field.
func ExternalIDEqualFold(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldExternalID), v))
	})
}

// ExternalIDContainsFold applies the ContainsFold predicate on the "external_id" field.
func ExternalIDContainsFold(v string) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldExternalID), v))
	})
}

// HasType applies the HasEdge predicate on the "type" edge.
func HasType() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TypeTable, TypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTypeWith applies the HasEdge predicate on the "type" edge with a given conditions (other predicates).
func HasTypeWith(preds ...predicate.EquipmentType) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
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
func HasLocation() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LocationTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, LocationTable, LocationColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLocationWith applies the HasEdge predicate on the "location" edge with a given conditions (other predicates).
func HasLocationWith(preds ...predicate.Location) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
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

// HasParentPosition applies the HasEdge predicate on the "parent_position" edge.
func HasParentPosition() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ParentPositionTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, ParentPositionTable, ParentPositionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasParentPositionWith applies the HasEdge predicate on the "parent_position" edge with a given conditions (other predicates).
func HasParentPositionWith(preds ...predicate.EquipmentPosition) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ParentPositionInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, ParentPositionTable, ParentPositionColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPositions applies the HasEdge predicate on the "positions" edge.
func HasPositions() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PositionsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PositionsTable, PositionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPositionsWith applies the HasEdge predicate on the "positions" edge with a given conditions (other predicates).
func HasPositionsWith(preds ...predicate.EquipmentPosition) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PositionsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PositionsTable, PositionsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPorts applies the HasEdge predicate on the "ports" edge.
func HasPorts() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PortsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PortsTable, PortsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPortsWith applies the HasEdge predicate on the "ports" edge with a given conditions (other predicates).
func HasPortsWith(preds ...predicate.EquipmentPort) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PortsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PortsTable, PortsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWorkOrder applies the HasEdge predicate on the "work_order" edge.
func HasWorkOrder() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, WorkOrderTable, WorkOrderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWorkOrderWith applies the HasEdge predicate on the "work_order" edge with a given conditions (other predicates).
func HasWorkOrderWith(preds ...predicate.WorkOrder) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WorkOrderInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, WorkOrderTable, WorkOrderColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProperties applies the HasEdge predicate on the "properties" edge.
func HasProperties() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertiesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertiesTable, PropertiesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPropertiesWith applies the HasEdge predicate on the "properties" edge with a given conditions (other predicates).
func HasPropertiesWith(preds ...predicate.Property) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PropertiesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PropertiesTable, PropertiesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasFiles applies the HasEdge predicate on the "files" edge.
func HasFiles() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(FilesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, FilesTable, FilesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasFilesWith applies the HasEdge predicate on the "files" edge with a given conditions (other predicates).
func HasFilesWith(preds ...predicate.File) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(FilesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, FilesTable, FilesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasHyperlinks applies the HasEdge predicate on the "hyperlinks" edge.
func HasHyperlinks() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(HyperlinksTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, HyperlinksTable, HyperlinksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasHyperlinksWith applies the HasEdge predicate on the "hyperlinks" edge with a given conditions (other predicates).
func HasHyperlinksWith(preds ...predicate.Hyperlink) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(HyperlinksInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, HyperlinksTable, HyperlinksColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEndpoints applies the HasEdge predicate on the "endpoints" edge.
func HasEndpoints() predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EndpointsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, EndpointsTable, EndpointsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEndpointsWith applies the HasEdge predicate on the "endpoints" edge with a given conditions (other predicates).
func HasEndpointsWith(preds ...predicate.ServiceEndpoint) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EndpointsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, EndpointsTable, EndpointsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Equipment) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Equipment) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
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
func Not(p predicate.Equipment) predicate.Equipment {
	return predicate.Equipment(func(s *sql.Selector) {
		p(s.Not())
	})
}
