// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package service

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.NEQ(s.C(FieldID), id))
		},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Service {
	return predicate.Service(
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
func IDNotIn(ids ...string) predicate.Service {
	return predicate.Service(
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
func IDGT(id string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.GT(s.C(FieldID), id))
		},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.GTE(s.C(FieldID), id))
		},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.LT(s.C(FieldID), id))
		},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			id, _ := strconv.Atoi(id)
			s.Where(sql.LTE(s.C(FieldID), id))
		},
	)
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldCreateTime), v))
		},
	)
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldUpdateTime), v))
		},
	)
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldName), v))
		},
	)
}

// ExternalID applies equality check predicate on the "external_id" field. It's identical to ExternalIDEQ.
func ExternalID(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldExternalID), v))
		},
	)
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(
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
func CreateTimeNotIn(vs ...time.Time) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(
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
func CreateTimeGT(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldCreateTime), v))
		},
	)
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldCreateTime), v))
		},
	)
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(
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
func UpdateTimeNotIn(vs ...time.Time) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(
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
func UpdateTimeGT(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldUpdateTime), v))
		},
	)
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldUpdateTime), v))
		},
	)
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldName), v))
		},
	)
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldName), v))
		},
	)
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(
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
func NameNotIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(
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
func NameGT(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldName), v))
		},
	)
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldName), v))
		},
	)
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldName), v))
		},
	)
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldName), v))
		},
	)
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldName), v))
		},
	)
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldName), v))
		},
	)
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldName), v))
		},
	)
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldName), v))
		},
	)
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldName), v))
		},
	)
}

// ExternalIDEQ applies the EQ predicate on the "external_id" field.
func ExternalIDEQ(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDNEQ applies the NEQ predicate on the "external_id" field.
func ExternalIDNEQ(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDIn applies the In predicate on the "external_id" field.
func ExternalIDIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(
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
func ExternalIDNotIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(
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
func ExternalIDGT(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDGTE applies the GTE predicate on the "external_id" field.
func ExternalIDGTE(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDLT applies the LT predicate on the "external_id" field.
func ExternalIDLT(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDLTE applies the LTE predicate on the "external_id" field.
func ExternalIDLTE(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDContains applies the Contains predicate on the "external_id" field.
func ExternalIDContains(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDHasPrefix applies the HasPrefix predicate on the "external_id" field.
func ExternalIDHasPrefix(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDHasSuffix applies the HasSuffix predicate on the "external_id" field.
func ExternalIDHasSuffix(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDIsNil applies the IsNil predicate on the "external_id" field.
func ExternalIDIsNil() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.IsNull(s.C(FieldExternalID)))
		},
	)
}

// ExternalIDNotNil applies the NotNil predicate on the "external_id" field.
func ExternalIDNotNil() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.NotNull(s.C(FieldExternalID)))
		},
	)
}

// ExternalIDEqualFold applies the EqualFold predicate on the "external_id" field.
func ExternalIDEqualFold(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldExternalID), v))
		},
	)
}

// ExternalIDContainsFold applies the ContainsFold predicate on the "external_id" field.
func ExternalIDContainsFold(v string) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldExternalID), v))
		},
	)
}

// HasType applies the HasEdge predicate on the "type" edge.
func HasType() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			t1 := s.Table()
			s.Where(sql.NotNull(t1.C(TypeColumn)))
		},
	)
}

// HasTypeWith applies the HasEdge predicate on the "type" edge with a given conditions (other predicates).
func HasTypeWith(preds ...predicate.ServiceType) predicate.Service {
	return predicate.Service(
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

// HasDownstream applies the HasEdge predicate on the "downstream" edge.
func HasDownstream() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			t1 := s.Table()
			builder := sql.Dialect(s.Dialect())
			s.Where(
				sql.In(
					t1.C(FieldID),
					builder.Select(DownstreamPrimaryKey[1]).
						From(builder.Table(DownstreamTable)),
				),
			)
		},
	)
}

// HasDownstreamWith applies the HasEdge predicate on the "downstream" edge with a given conditions (other predicates).
func HasDownstreamWith(preds ...predicate.Service) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Table(Table)
			t3 := builder.Table(DownstreamTable)
			t4 := builder.Select(t3.C(DownstreamPrimaryKey[1])).
				From(t3).
				Join(t2).
				On(t3.C(DownstreamPrimaryKey[0]), t2.C(FieldID))
			t5 := builder.Select().From(t2)
			for _, p := range preds {
				p(t5)
			}
			t4.FromSelect(t5)
			s.Where(sql.In(t1.C(FieldID), t4))
		},
	)
}

// HasUpstream applies the HasEdge predicate on the "upstream" edge.
func HasUpstream() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			t1 := s.Table()
			builder := sql.Dialect(s.Dialect())
			s.Where(
				sql.In(
					t1.C(FieldID),
					builder.Select(UpstreamPrimaryKey[0]).
						From(builder.Table(UpstreamTable)),
				),
			)
		},
	)
}

// HasUpstreamWith applies the HasEdge predicate on the "upstream" edge with a given conditions (other predicates).
func HasUpstreamWith(preds ...predicate.Service) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Table(Table)
			t3 := builder.Table(UpstreamTable)
			t4 := builder.Select(t3.C(UpstreamPrimaryKey[0])).
				From(t3).
				Join(t2).
				On(t3.C(UpstreamPrimaryKey[1]), t2.C(FieldID))
			t5 := builder.Select().From(t2)
			for _, p := range preds {
				p(t5)
			}
			t4.FromSelect(t5)
			s.Where(sql.In(t1.C(FieldID), t4))
		},
	)
}

// HasProperties applies the HasEdge predicate on the "properties" edge.
func HasProperties() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			t1 := s.Table()
			builder := sql.Dialect(s.Dialect())
			s.Where(
				sql.In(
					t1.C(FieldID),
					builder.Select(PropertiesColumn).
						From(builder.Table(PropertiesTable)).
						Where(sql.NotNull(PropertiesColumn)),
				),
			)
		},
	)
}

// HasPropertiesWith applies the HasEdge predicate on the "properties" edge with a given conditions (other predicates).
func HasPropertiesWith(preds ...predicate.Property) predicate.Service {
	return predicate.Service(
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

// HasTerminationPoints applies the HasEdge predicate on the "termination_points" edge.
func HasTerminationPoints() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			t1 := s.Table()
			builder := sql.Dialect(s.Dialect())
			s.Where(
				sql.In(
					t1.C(FieldID),
					builder.Select(TerminationPointsPrimaryKey[0]).
						From(builder.Table(TerminationPointsTable)),
				),
			)
		},
	)
}

// HasTerminationPointsWith applies the HasEdge predicate on the "termination_points" edge with a given conditions (other predicates).
func HasTerminationPointsWith(preds ...predicate.Equipment) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Table(TerminationPointsInverseTable)
			t3 := builder.Table(TerminationPointsTable)
			t4 := builder.Select(t3.C(TerminationPointsPrimaryKey[0])).
				From(t3).
				Join(t2).
				On(t3.C(TerminationPointsPrimaryKey[1]), t2.C(FieldID))
			t5 := builder.Select().From(t2)
			for _, p := range preds {
				p(t5)
			}
			t4.FromSelect(t5)
			s.Where(sql.In(t1.C(FieldID), t4))
		},
	)
}

// HasLinks applies the HasEdge predicate on the "links" edge.
func HasLinks() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			t1 := s.Table()
			builder := sql.Dialect(s.Dialect())
			s.Where(
				sql.In(
					t1.C(FieldID),
					builder.Select(LinksPrimaryKey[0]).
						From(builder.Table(LinksTable)),
				),
			)
		},
	)
}

// HasLinksWith applies the HasEdge predicate on the "links" edge with a given conditions (other predicates).
func HasLinksWith(preds ...predicate.Link) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Table(LinksInverseTable)
			t3 := builder.Table(LinksTable)
			t4 := builder.Select(t3.C(LinksPrimaryKey[0])).
				From(t3).
				Join(t2).
				On(t3.C(LinksPrimaryKey[1]), t2.C(FieldID))
			t5 := builder.Select().From(t2)
			for _, p := range preds {
				p(t5)
			}
			t4.FromSelect(t5)
			s.Where(sql.In(t1.C(FieldID), t4))
		},
	)
}

// HasCustomer applies the HasEdge predicate on the "customer" edge.
func HasCustomer() predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			t1 := s.Table()
			builder := sql.Dialect(s.Dialect())
			s.Where(
				sql.In(
					t1.C(FieldID),
					builder.Select(CustomerPrimaryKey[0]).
						From(builder.Table(CustomerTable)),
				),
			)
		},
	)
}

// HasCustomerWith applies the HasEdge predicate on the "customer" edge with a given conditions (other predicates).
func HasCustomerWith(preds ...predicate.Customer) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			builder := sql.Dialect(s.Dialect())
			t1 := s.Table()
			t2 := builder.Table(CustomerInverseTable)
			t3 := builder.Table(CustomerTable)
			t4 := builder.Select(t3.C(CustomerPrimaryKey[0])).
				From(t3).
				Join(t2).
				On(t3.C(CustomerPrimaryKey[1]), t2.C(FieldID))
			t5 := builder.Select().From(t2)
			for _, p := range preds {
				p(t5)
			}
			t4.FromSelect(t5)
			s.Where(sql.In(t1.C(FieldID), t4))
		},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Service) predicate.Service {
	return predicate.Service(
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
func Or(predicates ...predicate.Service) predicate.Service {
	return predicate.Service(
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
func Not(p predicate.Service) predicate.Service {
	return predicate.Service(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
