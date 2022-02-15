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
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

func dfs(q *Query, v queryVisitor) {
	v.preVisit(q)
	for _, x := range q.join {
		dfs(x, v)
	}
	v.postVisit(q)
}

type queryVisitor interface {
	preVisit(*Query)
	postVisit(*Query)
}

type columnNamesCollector struct {
	order   []string
	columns map[string][]string
}

func collectColumns(q *Query) *columnNamesCollector {
	colsCollector := &columnNamesCollector{columns: map[string][]string{}}
	dfs(q, colsCollector)
	return colsCollector
}

func collectFields(q *Query, columns map[string][]string) *fieldPointersCollector {
	fieldsCollector := &fieldPointersCollector{columns: columns}
	dfs(q, fieldsCollector)
	return fieldsCollector
}

func (c *columnNamesCollector) getColumnNames() []string {
	var columns []string
	for _, table := range c.order {
		for _, col := range c.columns[table] {
			name := fmt.Sprintf("%s.%s", table, col)
			columns = append(columns, name)
		}
	}
	return columns
}

func (c *columnNamesCollector) preVisit(q *Query) {
	metadata := q.arg.model.GetMetadata()
	table := metadata.Table
	c.order = append(c.order, table)

	fields := metadata.CreateObject().Fields()
	cols := getColumns(applyMask(fields, q.arg.mask))
	c.columns[table] = cols
}

func (*columnNamesCollector) postVisit(_ *Query) {}

func getColumns(fields FieldMap) []string {
	cols := make([]string, 0, len(fields))
	for k := range fields {
		cols = append(cols, k)
	}
	return cols
}

type fieldPointersCollector struct {
	columns  map[string][]string
	models   []Model
	pointers []interface{}
}

func (f *fieldPointersCollector) preVisit(q *Query) {
	metadata := q.arg.model.GetMetadata()
	model := metadata.CreateObject()
	fields := model.Fields()
	f.models = append(f.models, model)
	for _, col := range f.columns[metadata.Table] {
		f.pointers = append(f.pointers, fields[col].BaseType.ptr())
	}
}

func (*fieldPointersCollector) postVisit(_ *Query) {}

type joinClause struct {
	query *Query
}

func (j *joinClause) ToSql() (string, []interface{}, error) {
	b := &joinBuilder{
		parent: map[*Query]*Query{},
		sql:    strings.Builder{},
	}
	dfs(j.query, b)
	return b.sql.String(), b.args, b.err
}

type joinBuilder struct {
	parent map[*Query]*Query
	sql    strings.Builder
	args   []interface{}
	err    error
}

func (j *joinBuilder) preVisit(q *Query) {
	if j.err != nil {
		return
	}
	if j.parent[q] != nil {
		j.sql.WriteString(getJoinType(q.arg.nullable))
		j.sql.WriteString(" JOIN ")
		if len(q.join) > 0 {
			j.sql.WriteString("(")
		}
		j.sql.WriteString(q.arg.model.GetMetadata().Table)
	}
	for _, x := range q.join {
		j.parent[x] = q
	}
}

func getJoinType(nullable bool) string {
	if nullable {
		return " LEFT"
	}
	return ""
}

func (j *joinBuilder) postVisit(q *Query) {
	if j.err != nil {
		return
	}
	parent, ok := j.parent[q]
	if !ok {
		return
	}
	if len(q.join) > 0 {
		j.sql.WriteString(")")
	}
	j.sql.WriteString(" ON ")
	on := buildOn(q.arg.model.GetMetadata(), parent.arg.model.GetMetadata())
	var expr sq.Sqlizer = sq.Expr(on)
	if q.arg.filter != nil {
		expr = sq.And{expr, q.arg.filter}
	}
	sql, args, err := expr.ToSql()
	j.sql.WriteString(sql)
	j.args = append(j.args, args...)
	j.err = err
}

func buildOn(first, second *ModelMetadata) string {
	if _, ok := second.Relations[first.Table]; !ok {
		first, second = second, first
	}
	return fmt.Sprintf(
		"%s.id = %s.%s",
		first.Table, second.Table, second.Relations[first.Table],
	)
}
