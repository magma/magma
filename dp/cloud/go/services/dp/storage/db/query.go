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
}

type arg struct {
	model    Model
	mask     FieldMask
	nullable bool
	filter   sq.Sqlizer
}

func (q *Query) WithBuilder(builder sq.StatementBuilderType) *Query {
	q.builder = builder
	return q
}

func (q *Query) Select(mask FieldMask) *Query {
	q.arg.mask = mask
	return q
}

func (q *Query) From(model Model) *Query {
	q.arg.model = model
	return q
}

func (q *Query) Where(filter sq.Sqlizer) *Query {
	q.arg.filter = filter
	return q
}

func (q *Query) Join(other *Query) *Query {
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
