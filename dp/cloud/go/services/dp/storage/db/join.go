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

func On(lhsTable, lhsField, rhsTable, rhsField string) sq.Sqlizer {
	e := fmt.Sprintf("%s.%s = %s.%s", lhsTable, lhsField, rhsTable, rhsField)
	return sq.Expr(e)
}

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
	args  []*arg
	masks []indexMask
	nCols int
}

func collectColumns(q *Query) *columnNamesCollector {
	colsCollector := &columnNamesCollector{}
	dfs(q, colsCollector)
	return colsCollector
}

func (c *columnNamesCollector) getColumnNames() []string {
	var columns []string
	for i, arg := range c.args {
		fields := c.masks[i].filterColumns(arg.metadata)
		table := getTableName(arg)
		for _, field := range fields {
			name := fmt.Sprintf("%s.%s", table, field)
			columns = append(columns, name)
		}
	}
	return columns
}

func getTableName(arg *arg) string {
	if arg.alias != "" {
		return arg.alias
	}
	return arg.metadata.Table
}

func (c *columnNamesCollector) getPointers() ([]Model, []any) {
	models := make([]Model, len(c.args))
	pointers := make([]any, 0, c.nCols)
	for i, arg := range c.args {
		models[i] = arg.metadata.CreateObject()
		fields := c.masks[i].filterPointers(models[i])
		pointers = append(pointers, fields...)
	}
	return models, pointers
}

func (c *columnNamesCollector) preVisit(q *Query) {
	c.args = append(c.args, q.arg)
	mask := makeIndexMask(q.arg.metadata, q.arg.outputMask)
	c.masks = append(c.masks, mask)
	c.nCols += len(mask)
}

func (*columnNamesCollector) postVisit(_ *Query) {}

type joinClause struct {
	query *Query
}

func (j *joinClause) ToSql() (string, []any, error) {
	b := &joinBuilder{
		sql: strings.Builder{},
	}
	dfs(j.query, b)
	return b.sql.String(), b.args, b.err
}

type joinBuilder struct {
	sql  strings.Builder
	args []any
	err  error
}

func (j *joinBuilder) preVisit(q *Query) {
	if j.err != nil {
		return
	}
	if q.parent != nil {
		j.sql.WriteString(getJoinType(q.arg.nullable))
		j.sql.WriteString(" JOIN ")
		if len(q.join) > 0 {
			j.sql.WriteString("(")
		}
		j.sql.WriteString(buildFrom(q.arg))
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
	if q.parent == nil {
		return
	}
	if len(q.join) > 0 {
		j.sql.WriteString(")")
	}
	j.sql.WriteString(" ON ")
	sql, args, err := q.arg.on.ToSql()
	j.sql.WriteString(sql)
	j.args = append(j.args, args...)
	j.err = err
}

func makeIndexMask(metadata *ModelMetadata, mask FieldMask) indexMask {
	var indices indexMask
	for i, field := range metadata.Properties {
		if mask.ShouldInclude(field.Name) {
			indices = append(indices, i)
		}
	}
	return indices
}

type indexMask []int

func (im indexMask) filterColumns(metadata *ModelMetadata) []string {
	columns := make([]string, len(im))
	for i, j := range im {
		columns[i] = metadata.Properties[j].Name
	}
	return columns
}

func (im indexMask) filterPointers(model Model) []any {
	fields := model.Fields()
	pointers := make([]any, len(im))
	for i, j := range im {
		pointers[i] = fields[j].ptr()
	}
	return pointers
}
