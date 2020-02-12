// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package equipmentpositiondefinition

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func IDNotIn(ids ...string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func IDGT(id string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Index applies equality check predicate on the "index" field. It's identical to IndexEQ.
func Index(v int) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// VisibilityLabel applies equality check predicate on the "visibility_label" field. It's identical to VisibilityLabelEQ.
func VisibilityLabel(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVisibilityLabel), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func CreateTimeNotIn(vs ...time.Time) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func CreateTimeGT(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func UpdateTimeNotIn(vs ...time.Time) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func UpdateTimeGT(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func NameNotIn(vs ...string) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func NameGT(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// IndexEQ applies the EQ predicate on the "index" field.
func IndexEQ(v int) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIndex), v))
	})
}

// IndexNEQ applies the NEQ predicate on the "index" field.
func IndexNEQ(v int) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIndex), v))
	})
}

// IndexIn applies the In predicate on the "index" field.
func IndexIn(vs ...int) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func IndexNotIn(vs ...int) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func IndexGT(v int) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIndex), v))
	})
}

// IndexGTE applies the GTE predicate on the "index" field.
func IndexGTE(v int) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIndex), v))
	})
}

// IndexLT applies the LT predicate on the "index" field.
func IndexLT(v int) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIndex), v))
	})
}

// IndexLTE applies the LTE predicate on the "index" field.
func IndexLTE(v int) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIndex), v))
	})
}

// IndexIsNil applies the IsNil predicate on the "index" field.
func IndexIsNil() predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIndex)))
	})
}

// IndexNotNil applies the NotNil predicate on the "index" field.
func IndexNotNil() predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIndex)))
	})
}

// VisibilityLabelEQ applies the EQ predicate on the "visibility_label" field.
func VisibilityLabelEQ(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelNEQ applies the NEQ predicate on the "visibility_label" field.
func VisibilityLabelNEQ(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelIn applies the In predicate on the "visibility_label" field.
func VisibilityLabelIn(vs ...string) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func VisibilityLabelNotIn(vs ...string) predicate.EquipmentPositionDefinition {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func VisibilityLabelGT(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelGTE applies the GTE predicate on the "visibility_label" field.
func VisibilityLabelGTE(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelLT applies the LT predicate on the "visibility_label" field.
func VisibilityLabelLT(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelLTE applies the LTE predicate on the "visibility_label" field.
func VisibilityLabelLTE(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelContains applies the Contains predicate on the "visibility_label" field.
func VisibilityLabelContains(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelHasPrefix applies the HasPrefix predicate on the "visibility_label" field.
func VisibilityLabelHasPrefix(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelHasSuffix applies the HasSuffix predicate on the "visibility_label" field.
func VisibilityLabelHasSuffix(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelIsNil applies the IsNil predicate on the "visibility_label" field.
func VisibilityLabelIsNil() predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldVisibilityLabel)))
	})
}

// VisibilityLabelNotNil applies the NotNil predicate on the "visibility_label" field.
func VisibilityLabelNotNil() predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldVisibilityLabel)))
	})
}

// VisibilityLabelEqualFold applies the EqualFold predicate on the "visibility_label" field.
func VisibilityLabelEqualFold(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldVisibilityLabel), v))
	})
}

// VisibilityLabelContainsFold applies the ContainsFold predicate on the "visibility_label" field.
func VisibilityLabelContainsFold(v string) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldVisibilityLabel), v))
	})
}

// HasPositions applies the HasEdge predicate on the "positions" edge.
func HasPositions() predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PositionsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, PositionsTable, PositionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPositionsWith applies the HasEdge predicate on the "positions" edge with a given conditions (other predicates).
func HasPositionsWith(preds ...predicate.EquipmentPosition) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PositionsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, PositionsTable, PositionsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEquipmentType applies the HasEdge predicate on the "equipment_type" edge.
func HasEquipmentType() predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EquipmentTypeTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EquipmentTypeTable, EquipmentTypeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEquipmentTypeWith applies the HasEdge predicate on the "equipment_type" edge with a given conditions (other predicates).
func HasEquipmentTypeWith(preds ...predicate.EquipmentType) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func And(predicates ...predicate.EquipmentPositionDefinition) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.EquipmentPositionDefinition) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
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
func Not(p predicate.EquipmentPositionDefinition) predicate.EquipmentPositionDefinition {
	return predicate.EquipmentPositionDefinition(func(s *sql.Selector) {
		p(s.Not())
	})
}
