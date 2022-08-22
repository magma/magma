/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sqorc

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	setValues []UpsertValue
	columns   []string

	expectedSql  string
	expectedArgs []interface{}
}

func TestPostgresInsertBuilder_OnConflict(t *testing.T) {
	ib := postgresStatementBuilder{squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}.Insert("table")

	cases := []testCase{
		{
			setValues:    nil,
			columns:      []string{"foo"},
			expectedSql:  "INSERT INTO table (foo,bar) VALUES ($1,$2) ON CONFLICT (foo) DO NOTHING",
			expectedArgs: []interface{}{},
		},
		{
			setValues:    []UpsertValue{},
			columns:      []string{"foo", "bar"},
			expectedSql:  "INSERT INTO table (foo,bar) VALUES ($1,$2) ON CONFLICT (foo, bar) DO NOTHING",
			expectedArgs: []interface{}{},
		},
		{
			setValues: []UpsertValue{
				{Column: "foo", Value: 1},
			},
			columns:      []string{"foo", "bar"},
			expectedSql:  "INSERT INTO table (foo,bar) VALUES ($1,$2) ON CONFLICT (foo, bar) DO UPDATE SET foo = $3",
			expectedArgs: []interface{}{1},
		},
		{
			setValues: []UpsertValue{
				{Column: "foo", Value: 1},
				{Column: "bar", Value: "baz"},
			},
			columns:      []string{"foo", "bar"},
			expectedSql:  "INSERT INTO table (foo,bar) VALUES ($1,$2) ON CONFLICT (foo, bar) DO UPDATE SET foo = $3, bar = $4",
			expectedArgs: []interface{}{1, "baz"},
		},
	}
	runCases(t, cases, ib)
}

func TestMysqlInsertBuilder_OnConflict(t *testing.T) {
	ib := mariaDBStatementBuilder{squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)}.Insert("table")

	cases := []testCase{
		{
			setValues:    nil,
			columns:      []string{"foo"},
			expectedSql:  "INSERT IGNORE INTO table (foo,bar) VALUES (?,?)",
			expectedArgs: []interface{}{},
		},
		{
			setValues:    []UpsertValue{},
			columns:      []string{"foo", "bar"},
			expectedSql:  "INSERT IGNORE INTO table (foo,bar) VALUES (?,?)",
			expectedArgs: []interface{}{},
		},
		{
			setValues: []UpsertValue{
				{Column: "foo", Value: 1},
			},
			columns:      []string{"foo", "bar"},
			expectedSql:  "INSERT INTO table (foo,bar) VALUES (?,?) ON DUPLICATE KEY UPDATE foo = ?",
			expectedArgs: []interface{}{1},
		},
		{
			setValues: []UpsertValue{
				{Column: "foo", Value: 1},
				{Column: "bar", Value: "baz"},
			},
			columns:      []string{"foo", "bar"},
			expectedSql:  "INSERT INTO table (foo,bar) VALUES (?,?) ON DUPLICATE KEY UPDATE foo = ?, bar = ?",
			expectedArgs: []interface{}{1, "baz"},
		},
	}
	runCases(t, cases, ib)
}

func runCases(t *testing.T, tcs []testCase, ib InsertBuilder) {
	for _, tc := range tcs {
		actualSql, actualArgs, err := ib.Columns("foo", "bar").Values("fooV", "barV").OnConflict(tc.setValues, tc.columns...).ToSql()
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedSql, actualSql)
		expectedArgs := []interface{}{"fooV", "barV"}
		expectedArgs = append(expectedArgs, tc.expectedArgs...)
		assert.Equal(t, expectedArgs, actualArgs)
	}
}
