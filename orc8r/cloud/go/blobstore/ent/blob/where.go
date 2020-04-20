/*
 * Copyright 2020 The Magma Authors
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated (@generated) by entc, DO NOT EDIT.

package blob

import (
	"magma/orc8r/cloud/go/blobstore/ent/predicate"

	"github.com/facebookincubator/ent/dialect/sql"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldID), id))
		},
	)
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldID), id))
		},
	)
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
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
		},
	)
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
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
		},
	)
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldID), id))
		},
	)
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldID), id))
		},
	)
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldID), id))
		},
	)
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldID), id))
		},
	)
}

// NetworkID applies equality check predicate on the "network_id" field. It's identical to NetworkIDEQ.
func NetworkID(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldNetworkID), v))
		},
	)
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldType), v))
		},
	)
}

// Key applies equality check predicate on the "key" field. It's identical to KeyEQ.
func Key(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldKey), v))
		},
	)
}

// Value applies equality check predicate on the "value" field. It's identical to ValueEQ.
func Value(v []byte) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldValue), v))
		},
	)
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v uint64) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldVersion), v))
		},
	)
}

// NetworkIDEQ applies the EQ predicate on the "network_id" field.
func NetworkIDEQ(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDNEQ applies the NEQ predicate on the "network_id" field.
func NetworkIDNEQ(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDIn applies the In predicate on the "network_id" field.
func NetworkIDIn(vs ...string) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldNetworkID), v...))
		},
	)
}

// NetworkIDNotIn applies the NotIn predicate on the "network_id" field.
func NetworkIDNotIn(vs ...string) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldNetworkID), v...))
		},
	)
}

// NetworkIDGT applies the GT predicate on the "network_id" field.
func NetworkIDGT(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDGTE applies the GTE predicate on the "network_id" field.
func NetworkIDGTE(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDLT applies the LT predicate on the "network_id" field.
func NetworkIDLT(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDLTE applies the LTE predicate on the "network_id" field.
func NetworkIDLTE(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDContains applies the Contains predicate on the "network_id" field.
func NetworkIDContains(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDHasPrefix applies the HasPrefix predicate on the "network_id" field.
func NetworkIDHasPrefix(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDHasSuffix applies the HasSuffix predicate on the "network_id" field.
func NetworkIDHasSuffix(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDEqualFold applies the EqualFold predicate on the "network_id" field.
func NetworkIDEqualFold(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldNetworkID), v))
		},
	)
}

// NetworkIDContainsFold applies the ContainsFold predicate on the "network_id" field.
func NetworkIDContainsFold(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldNetworkID), v))
		},
	)
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldType), v))
		},
	)
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldType), v))
		},
	)
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldType), v...))
		},
	)
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldType), v...))
		},
	)
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldType), v))
		},
	)
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldType), v))
		},
	)
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldType), v))
		},
	)
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldType), v))
		},
	)
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldType), v))
		},
	)
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldType), v))
		},
	)
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldType), v))
		},
	)
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldType), v))
		},
	)
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldType), v))
		},
	)
}

// KeyEQ applies the EQ predicate on the "key" field.
func KeyEQ(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldKey), v))
		},
	)
}

// KeyNEQ applies the NEQ predicate on the "key" field.
func KeyNEQ(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldKey), v))
		},
	)
}

// KeyIn applies the In predicate on the "key" field.
func KeyIn(vs ...string) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldKey), v...))
		},
	)
}

// KeyNotIn applies the NotIn predicate on the "key" field.
func KeyNotIn(vs ...string) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldKey), v...))
		},
	)
}

// KeyGT applies the GT predicate on the "key" field.
func KeyGT(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldKey), v))
		},
	)
}

// KeyGTE applies the GTE predicate on the "key" field.
func KeyGTE(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldKey), v))
		},
	)
}

// KeyLT applies the LT predicate on the "key" field.
func KeyLT(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldKey), v))
		},
	)
}

// KeyLTE applies the LTE predicate on the "key" field.
func KeyLTE(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldKey), v))
		},
	)
}

// KeyContains applies the Contains predicate on the "key" field.
func KeyContains(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.Contains(s.C(FieldKey), v))
		},
	)
}

// KeyHasPrefix applies the HasPrefix predicate on the "key" field.
func KeyHasPrefix(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.HasPrefix(s.C(FieldKey), v))
		},
	)
}

// KeyHasSuffix applies the HasSuffix predicate on the "key" field.
func KeyHasSuffix(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.HasSuffix(s.C(FieldKey), v))
		},
	)
}

// KeyEqualFold applies the EqualFold predicate on the "key" field.
func KeyEqualFold(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EqualFold(s.C(FieldKey), v))
		},
	)
}

// KeyContainsFold applies the ContainsFold predicate on the "key" field.
func KeyContainsFold(v string) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.ContainsFold(s.C(FieldKey), v))
		},
	)
}

// ValueEQ applies the EQ predicate on the "value" field.
func ValueEQ(v []byte) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldValue), v))
		},
	)
}

// ValueNEQ applies the NEQ predicate on the "value" field.
func ValueNEQ(v []byte) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldValue), v))
		},
	)
}

// ValueIn applies the In predicate on the "value" field.
func ValueIn(vs ...[]byte) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldValue), v...))
		},
	)
}

// ValueNotIn applies the NotIn predicate on the "value" field.
func ValueNotIn(vs ...[]byte) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldValue), v...))
		},
	)
}

// ValueGT applies the GT predicate on the "value" field.
func ValueGT(v []byte) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldValue), v))
		},
	)
}

// ValueGTE applies the GTE predicate on the "value" field.
func ValueGTE(v []byte) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldValue), v))
		},
	)
}

// ValueLT applies the LT predicate on the "value" field.
func ValueLT(v []byte) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldValue), v))
		},
	)
}

// ValueLTE applies the LTE predicate on the "value" field.
func ValueLTE(v []byte) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldValue), v))
		},
	)
}

// ValueIsNil applies the IsNil predicate on the "value" field.
func ValueIsNil() predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.IsNull(s.C(FieldValue)))
		},
	)
}

// ValueNotNil applies the NotNil predicate on the "value" field.
func ValueNotNil() predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.NotNull(s.C(FieldValue)))
		},
	)
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v uint64) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.EQ(s.C(FieldVersion), v))
		},
	)
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v uint64) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.NEQ(s.C(FieldVersion), v))
		},
	)
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...uint64) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.In(s.C(FieldVersion), v...))
		},
	)
}

// VersionNotIn applies the NotIn predicate on the "version" field.
func VersionNotIn(vs ...uint64) predicate.Blob {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Blob(
		func(s *sql.Selector) {
			// if not arguments were provided, append the FALSE constants,
			// since we can't apply "IN ()". This will make this predicate falsy.
			if len(vs) == 0 {
				s.Where(sql.False())
				return
			}
			s.Where(sql.NotIn(s.C(FieldVersion), v...))
		},
	)
}

// VersionGT applies the GT predicate on the "version" field.
func VersionGT(v uint64) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GT(s.C(FieldVersion), v))
		},
	)
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v uint64) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.GTE(s.C(FieldVersion), v))
		},
	)
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v uint64) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LT(s.C(FieldVersion), v))
		},
	)
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v uint64) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			s.Where(sql.LTE(s.C(FieldVersion), v))
		},
	)
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Blob) predicate.Blob {
	return predicate.Blob(
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
func Or(predicates ...predicate.Blob) predicate.Blob {
	return predicate.Blob(
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
func Not(p predicate.Blob) predicate.Blob {
	return predicate.Blob(
		func(s *sql.Selector) {
			p(s.Not())
		},
	)
}
