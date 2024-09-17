/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package db

import (
	sq "github.com/Masterminds/squirrel"
)

func NewQuery() *Query {
	return &Query{
		arg:        &arg{},
		pagination: &pagination{},
	}
}

type Query struct {
	builder    sq.StatementBuilderType
	arg        *arg
	join       []*Query
	pagination *pagination
	parent     *Query
}

type arg struct {
	alias      string
	model      Model
	metadata   *ModelMetadata
	inputMask  FieldMask
	outputMask FieldMask
	nullable   bool
	filter     sq.Sqlizer
	on         sq.Sqlizer
	lock       string
}

func (q *Query) WithBuilder(builder sq.StatementBuilderType) *Query {
	q.builder = builder
	return q
}

func (q *Query) Select(mask FieldMask) *Query {
	q.arg.outputMask = mask
	return q
}

func (q *Query) From(model Model) *Query {
	q.arg.model = model
	q.arg.metadata = model.GetMetadata()
	return q
}

func (q *Query) As(alias string) *Query {
	q.arg.alias = alias
	return q
}

func (q *Query) Where(filter sq.Sqlizer) *Query {
	q.arg.filter = filter
	return q
}

func (q *Query) Lock(lock string) *Query {
	q.arg.lock = lock
	return q
}

func (q *Query) On(cond sq.Sqlizer) *Query {
	q.arg.on = cond
	return q
}

func (q *Query) Join(other *Query) *Query {
	other.parent = q
	q.join = append(q.join, other)
	return q
}

func (q *Query) Nullable() *Query {
	q.arg.nullable = true
	return q
}

func (q *Query) Limit(limit int64) *Query {
	q.pagination.limit = MakeInt(limit)
	return q
}

func (q *Query) Offset(offset int64) *Query {
	q.pagination.offset = MakeInt(offset)
	return q
}

func (q *Query) OrderBy(column string, order Order) *Query {
	q.pagination.orderBy = column
	q.pagination.order = order
	return q
}
