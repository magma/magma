// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipment

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.NEQ(s.C(FieldID), id))
		},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Equipment {
	return predicate.Equipment(
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
func IDNotIn(ids ...string) predicate.Equipment {
	return predicate.Equipment(
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
func IDGT(id string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.GT(s.C(FieldID), id))
		},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.GTE(s.C(FieldID), id))
		},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.LT(s.C(FieldID), id))
		},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.LTE(s.C(FieldID), id))
		},
	)
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldCreateTime), v))
		},
	)
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldUpdateTime), v))
		},
	)
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldName), v))
		},
	)
}

// FutureState applies equality check predicate on the "future_state" field. It's identical to FutureStateEQ.
func FutureState(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldFutureState), v))
		},
	)
}

// DeviceID applies equality check predicate on the "device_id" field. It's identical to DeviceIDEQ.
func DeviceID(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldDeviceID), v))
		},
	)
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
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
func CreateTimeNotIn(vs ...time.Time) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
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
func CreateTimeGT(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldCreateTime), v))
		},
	)
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
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
func UpdateTimeNotIn(vs ...time.Time) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
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
func UpdateTimeGT(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldUpdateTime), v))
		},
	)
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldName), v))
		},
	)
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldName), v))
		},
	)
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
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
func NameNotIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
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
func NameGT(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldName), v))
		},
	)
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldName), v))
		},
	)
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldName), v))
		},
	)
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldName), v))
		},
	)
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldName), v))
		},
	)
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldName), v))
		},
	)
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldName), v))
		},
	)
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldName), v))
		},
	)
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldName), v))
		},
	)
}

// FutureStateEQ applies the EQ predicate on the "future_state" field.
func FutureStateEQ(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateNEQ applies the NEQ predicate on the "future_state" field.
func FutureStateNEQ(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateIn applies the In predicate on the "future_state" field.
func FutureStateIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldFutureState), v...))
		},
	)
}

// FutureStateNotIn applies the NotIn predicate on the "future_state" field.
func FutureStateNotIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldFutureState), v...))
		},
	)
}

// FutureStateGT applies the GT predicate on the "future_state" field.
func FutureStateGT(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateGTE applies the GTE predicate on the "future_state" field.
func FutureStateGTE(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateLT applies the LT predicate on the "future_state" field.
func FutureStateLT(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateLTE applies the LTE predicate on the "future_state" field.
func FutureStateLTE(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateContains applies the Contains predicate on the "future_state" field.
func FutureStateContains(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateHasPrefix applies the HasPrefix predicate on the "future_state" field.
func FutureStateHasPrefix(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateHasSuffix applies the HasSuffix predicate on the "future_state" field.
func FutureStateHasSuffix(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateIsNil applies the IsNil predicate on the "future_state" field.
func FutureStateIsNil() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.IsNull(s.C(FieldFutureState)))
		},
	)
}

// FutureStateNotNil applies the NotNil predicate on the "future_state" field.
func FutureStateNotNil() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.NotNull(s.C(FieldFutureState)))
		},
	)
}

// FutureStateEqualFold applies the EqualFold predicate on the "future_state" field.
func FutureStateEqualFold(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldFutureState), v))
		},
	)
}

// FutureStateContainsFold applies the ContainsFold predicate on the "future_state" field.
func FutureStateContainsFold(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldFutureState), v))
		},
	)
}

// DeviceIDEQ applies the EQ predicate on the "device_id" field.
func DeviceIDEQ(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDNEQ applies the NEQ predicate on the "device_id" field.
func DeviceIDNEQ(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDIn applies the In predicate on the "device_id" field.
func DeviceIDIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldDeviceID), v...))
		},
	)
}

// DeviceIDNotIn applies the NotIn predicate on the "device_id" field.
func DeviceIDNotIn(vs ...string) predicate.Equipment {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Equipment(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldDeviceID), v...))
		},
	)
}

// DeviceIDGT applies the GT predicate on the "device_id" field.
func DeviceIDGT(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDGTE applies the GTE predicate on the "device_id" field.
func DeviceIDGTE(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDLT applies the LT predicate on the "device_id" field.
func DeviceIDLT(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDLTE applies the LTE predicate on the "device_id" field.
func DeviceIDLTE(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDContains applies the Contains predicate on the "device_id" field.
func DeviceIDContains(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDHasPrefix applies the HasPrefix predicate on the "device_id" field.
func DeviceIDHasPrefix(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDHasSuffix applies the HasSuffix predicate on the "device_id" field.
func DeviceIDHasSuffix(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDIsNil applies the IsNil predicate on the "device_id" field.
func DeviceIDIsNil() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.IsNull(s.C(FieldDeviceID)))
		},
	)
}

// DeviceIDNotNil applies the NotNil predicate on the "device_id" field.
func DeviceIDNotNil() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.NotNull(s.C(FieldDeviceID)))
		},
	)
}

// DeviceIDEqualFold applies the EqualFold predicate on the "device_id" field.
func DeviceIDEqualFold(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldDeviceID), v))
		},
	)
}

// DeviceIDContainsFold applies the ContainsFold predicate on the "device_id" field.
func DeviceIDContainsFold(v string) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldDeviceID), v))
		},
	)
}

// HasType applies the HasEdge predicate on the "type" edge.
func HasType() predicate.Equipment {
	return predicate.Equipment(
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
func HasTypeWith(preds ...predicate.EquipmentType) predicate.Equipment {
	return predicate.Equipment(
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

// HasLocation applies the HasEdge predicate on the "location" edge.
func HasLocation() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(LocationTable, FieldID),
				sql.Edge(sql.M2O, true, LocationTable, LocationColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasLocationWith applies the HasEdge predicate on the "location" edge with a given conditions (other predicates).
func HasLocationWith(preds ...predicate.Location) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(FieldID).From(builder.Table(LocationInverseTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(LocationColumn), t2))
		},
	)
}

// HasParentPosition applies the HasEdge predicate on the "parent_position" edge.
func HasParentPosition() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(ParentPositionTable, FieldID),
				sql.Edge(sql.O2O, true, ParentPositionTable, ParentPositionColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasParentPositionWith applies the HasEdge predicate on the "parent_position" edge with a given conditions (other predicates).
func HasParentPositionWith(preds ...predicate.EquipmentPosition) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(FieldID).From(builder.Table(ParentPositionInverseTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(ParentPositionColumn), t2))
		},
	)
}

// HasPositions applies the HasEdge predicate on the "positions" edge.
func HasPositions() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(PositionsTable, FieldID),
				sql.Edge(sql.O2M, false, PositionsTable, PositionsColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasPositionsWith applies the HasEdge predicate on the "positions" edge with a given conditions (other predicates).
func HasPositionsWith(preds ...predicate.EquipmentPosition) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(PositionsColumn).From(builder.Table(PositionsTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasPorts applies the HasEdge predicate on the "ports" edge.
func HasPorts() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(PortsTable, FieldID),
				sql.Edge(sql.O2M, false, PortsTable, PortsColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasPortsWith applies the HasEdge predicate on the "ports" edge with a given conditions (other predicates).
func HasPortsWith(preds ...predicate.EquipmentPort) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(PortsColumn).From(builder.Table(PortsTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(FieldID), t2))
		},
	)
}

// HasWorkOrder applies the HasEdge predicate on the "work_order" edge.
func HasWorkOrder() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(WorkOrderTable, FieldID),
				sql.Edge(sql.M2O, false, WorkOrderTable, WorkOrderColumn),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasWorkOrderWith applies the HasEdge predicate on the "work_order" edge with a given conditions (other predicates).
func HasWorkOrderWith(preds ...predicate.WorkOrder) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Select(FieldID).From(builder.Table(WorkOrderInverseTable))
			for _, p := range preds {
				p(t2)
			}
			s.Where(sql.In(t1.C(WorkOrderColumn), t2))
		},
	)
}

// HasProperties applies the HasEdge predicate on the "properties" edge.
func HasProperties() predicate.Equipment {
	return predicate.Equipment(
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
func HasPropertiesWith(preds ...predicate.Property) predicate.Equipment {
	return predicate.Equipment(
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

// HasService applies the HasEdge predicate on the "service" edge.
func HasService() predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			step := sql.NewStep(
				sql.From(Table, FieldID),
				sql.To(ServiceTable, FieldID),
				sql.Edge(sql.M2M, true, ServiceTable, ServicePrimaryKey...),
			)
			sql.HasNeighbors(s, step)
		},
	)
}

// HasServiceWith applies the HasEdge predicate on the "service" edge with a given conditions (other predicates).
func HasServiceWith(preds ...predicate.Service) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Table(ServiceInverseTable)
			t3 := builder.Table(ServiceTable)
			t4 := builder.Select(t3.C(ServicePrimaryKey[1])).
				From(t3).
				Join(t2).
				On(t3.C(ServicePrimaryKey[0]), t2.C(FieldID))
			t5 := builder.Select().From(t2)
			for _, p := range preds {
				p(t5)
			}
			t4.FromSelect(t5)
			s.Where(sql.In(t1.C(FieldID), t4))
		},
	)
}

// HasFiles applies the HasEdge predicate on the "files" edge.
func HasFiles() predicate.Equipment {
	return predicate.Equipment(
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
func HasFilesWith(preds ...predicate.File) predicate.Equipment {
	return predicate.Equipment(
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

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Equipment) predicate.Equipment {
	return predicate.Equipment(
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
func Or(predicates ...predicate.Equipment) predicate.Equipment {
	return predicate.Equipment(
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
func Not(p predicate.Equipment) predicate.Equipment {
	return predicate.Equipment(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
