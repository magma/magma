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

	"magma/orc8r/cloud/go/sqorc"
)

func (q *Query) Insert(mask FieldMask) (int64, error) {
	var id int64
	q.arg.inputMask = mask
	err := q.builder.
		Insert(buildFrom(q.arg)).
		SetMap(filterValues(q.arg)).
		Suffix("RETURNING id").
		QueryRow().
		Scan(&id)
	return id, err
}

func (q *Query) Update(mask FieldMask) ([]Model, error) {
	q.arg.inputMask = mask
	baseQuery := q.builder.
		Update(buildFrom(q.arg)).
		SetMap(filterValues(q.arg)).
		Where(q.arg.filter)

	c := collectColumns(q)
	cols := c.getColumnNames()
	if cols == nil {
		_, err := baseQuery.Exec()
		return nil, err
	}
	models, pointers := c.getPointers()
	suffix := getSuffix(cols)
	err := baseQuery.
		Suffix(suffix).
		QueryRow().
		Scan(pointers...)
	return models, err
}

func getSuffix(columns []string) string {
	return "RETURNING " + strings.Join(columns, ", ")
}

func (q *Query) Delete() error {
	_, err := q.builder.
		Delete(buildFrom(q.arg)).
		Where(q.arg.filter).
		Exec()
	return err
}

func (q *Query) Count() (int64, error) {
	var count int64
	err := q.builder.
		Select("COUNT(*)").
		From(buildFrom(q.arg)).
		Where(q.arg.filter).
		QueryRow().
		Scan(&count)
	return count, err
}

func (q *Query) Fetch() ([]Model, error) {
	c := collectColumns(q)
	models, pointers := c.getPointers()
	query := buildQuery(q.builder, q, q.arg.filter, c.getColumnNames())
	err := query.
		QueryRow().
		Scan(pointers...)
	return models, err
}

func (q *Query) List() ([][]Model, error) {
	c := collectColumns(q)
	query := buildQuery(q.builder, q, q.arg.filter, c.getColumnNames())
	query = addPagination(q.pagination, query)
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer sqorc.CloseRowsLogOnError(rows, "List")
	var result [][]Model
	for rows.Next() {
		models, pointers := c.getPointers()
		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}
		result = append(result, models)
	}
	return result, rows.Err()
}

func buildFrom(arg *arg) string {
	return fmt.Sprintf("%s %s", arg.metadata.Table, arg.alias)
}

func filterValues(arg *arg) map[string]interface{} {
	values := map[string]interface{}{}
	fields := arg.model.Fields()
	for i, p := range arg.metadata.Properties {
		if arg.inputMask.ShouldInclude(p.Name) {
			if p.HasDefault && fields[i].isNull() {
				values[p.Name] = p.DefaultValue
			} else {
				values[p.Name] = fields[i].value()
			}
		}
	}
	return values
}

func buildQuery(builder sq.StatementBuilderType, q *Query, filter sq.Sqlizer, columns []string) sq.SelectBuilder {
	return builder.
		Select(columns...).
		From(buildFrom(q.arg)).
		JoinClause(&joinClause{query: q}).
		Where(filter).
		Suffix(q.arg.lock)
}
