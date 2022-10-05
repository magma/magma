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
	"log"
	"testing"

	"github.com/lann/builder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestColumnBuilder_ToSql(t *testing.T) {
	// psql

	actual, err := columnBuilder(postgresColumnTypeMap).
		Name("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		References("foo", "bar").
		Default("\"hello world\"").
		NotNull().
		OnDelete(ColumnOnDeleteCascade).
		ToSql()
	assert.NoError(t, err)
	expected := "pk TEXT PRIMARY KEY NOT NULL DEFAULT \"hello world\" REFERENCES foo (bar) ON DELETE CASCADE"
	assert.Equal(t, expected, actual)

	actual, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		ToSql()
	assert.NoError(t, err)
	expected = "foo BYTEA"
	assert.Equal(t, expected, actual)

	actual, err = columnBuilder(postgresColumnTypeMap).
		Name("version").
		Type(ColumnTypeInt).
		NotNull().
		Default(0).
		ToSql()
	assert.NoError(t, err)
	expected = "version INTEGER NOT NULL DEFAULT 0"
	assert.Equal(t, expected, actual)

	// maria
	actual, err = columnBuilder(mariaColumnTypeMap).
		Name("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		References("foo", "bar").
		Default("\"hello world\"").
		NotNull().
		OnDelete(ColumnOnDeleteCascade).
		ToSql()
	assert.NoError(t, err)
	expected = "pk VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT \"hello world\" REFERENCES foo (bar) ON DELETE CASCADE"
	assert.Equal(t, expected, actual)

	actual, err = columnBuilder(mariaColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		ToSql()
	assert.NoError(t, err)
	expected = "foo LONGBLOB"
	assert.Equal(t, expected, actual)

	actual, err = columnBuilder(mariaColumnTypeMap).
		Name("version").
		Type(ColumnTypeInt).
		NotNull().
		Default(0).
		ToSql()
	assert.NoError(t, err)
	expected = "version INT NOT NULL DEFAULT 0"
	assert.Equal(t, expected, actual)
}

func TestColumnBuilder_ToSql_Errors(t *testing.T) {
	_, err := columnBuilder(postgresColumnTypeMap).
		Type(ColumnTypeBytes).
		ToSql()
	assert.EqualError(t, err, "column name must be specified")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		ToSql()
	assert.EqualError(t, err, "column type must be specified")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes + 100).
		ToSql()
	assert.EqualError(t, err, "column type 103 not recognized")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		References("", "bar").
		ToSql()
	assert.EqualError(t, err, "reference table name and column of foreign key must not be empty")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		References("bar", "").
		ToSql()
	assert.EqualError(t, err, "reference table name and column of foreign key must not be empty")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		OnDelete(ColumnOnDeleteCascade).
		ToSql()
	assert.EqualError(t, err, "cannot specify an ON DELETE without a REFERENCES")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		References("bar", "baz").
		OnDelete(ColumnOnDeleteOption(255)).
		ToSql()
	assert.EqualError(t, err, "unrecognized on delete option 255")
}

func TestCreateTableBuilder_ToSql(t *testing.T) {
	// psql
	// we allow setting primary key constraint on a column and at the table
	// level even though it would be rejected by the engine
	actual, _, err := tableBuilder(postgresColumnTypeMap).
		Name("foobar").
		IfNotExists().
		Column("pk").Type(ColumnTypeText).PrimaryKey().EndColumn().
		Column("foo").Type(ColumnTypeBytes).NotNull().References("barbaz", "bites").OnDelete(ColumnOnDeleteCascade).EndColumn().
		Column("bar").Type(ColumnTypeInt).Default(42).EndColumn().
		PrimaryKey("pk", "foo").
		Unique("foo", "bar").
		ForeignKey("othert", map[string]string{"foo": "ofoo", "bar": "obar"}, ColumnOnDeleteDoNothing).
		ForeignKey("othert", map[string]string{"bar": "zbar"}, ColumnOnDeleteCascade).
		ToSql()
	assert.NoError(t, err)
	expected := "CREATE TABLE IF NOT EXISTS foobar (\n" +
		"pk TEXT PRIMARY KEY,\n" +
		"foo BYTEA NOT NULL REFERENCES barbaz (bites) ON DELETE CASCADE,\n" +
		"bar INTEGER DEFAULT 42,\n" +
		"PRIMARY KEY (pk, foo),\n" +
		"FOREIGN KEY (bar, foo) REFERENCES othert (obar, ofoo),\n" +
		"FOREIGN KEY (bar) REFERENCES othert (zbar) ON DELETE CASCADE,\n" +
		"UNIQUE (foo, bar)\n" +
		")"
	assert.Equal(t, expected, actual)

	actual, _, err = tableBuilder(postgresColumnTypeMap).
		Name("foobar").
		Column("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		EndColumn().
		ToSql()
	assert.NoError(t, err)
	expected = "CREATE TABLE foobar (\n" +
		"pk TEXT PRIMARY KEY\n" +
		")"
	assert.Equal(t, expected, actual)

	// maria
	actual, _, err = tableBuilder(mariaColumnTypeMap).
		Name("foobar").
		IfNotExists().
		Column("pk").Type(ColumnTypeText).PrimaryKey().EndColumn().
		Column("foo").Type(ColumnTypeBytes).NotNull().References("barbaz", "bites").OnDelete(ColumnOnDeleteCascade).EndColumn().
		Column("bar").Type(ColumnTypeInt).Default(42).EndColumn().
		PrimaryKey("pk", "foo").
		Unique("foo", "bar").
		ForeignKey("othert", map[string]string{"foo": "ofoo", "bar": "obar"}, ColumnOnDeleteDoNothing).
		ForeignKey("othert", map[string]string{"bar": "zbar"}, ColumnOnDeleteCascade).
		ToSql()
	assert.NoError(t, err)
	expected = "CREATE TABLE IF NOT EXISTS foobar (\n" +
		"pk VARCHAR(255) PRIMARY KEY,\n" +
		"foo LONGBLOB NOT NULL REFERENCES barbaz (bites) ON DELETE CASCADE,\n" +
		"bar INT DEFAULT 42,\n" +
		"PRIMARY KEY (pk, foo),\n" +
		"FOREIGN KEY (bar, foo) REFERENCES othert (obar, ofoo),\n" +
		"FOREIGN KEY (bar) REFERENCES othert (zbar) ON DELETE CASCADE,\n" +
		"UNIQUE (foo, bar)\n" +
		")"
	assert.Equal(t, expected, actual)

	actual, _, err = tableBuilder(mariaColumnTypeMap).
		Name("foobar").
		Column("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		EndColumn().
		ToSql()
	assert.NoError(t, err)
	expected = "CREATE TABLE foobar (\n" +
		"pk VARCHAR(255) PRIMARY KEY\n" +
		")"
	assert.Equal(t, expected, actual)
}

func TestCreateTableBuilder_Exec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error oppening stub DB conn: %s", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing stub DB: %s", err)
		}
	}()

	mock.ExpectBegin()
	mock.ExpectExec(
		"CREATE TABLE IF NOT EXISTS foobar \\(\n" +
			"pk TEXT PRIMARY KEY\n" +
			"\\)",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	_, err = NewPostgresStatementBuilder().CreateTable("foobar").
		IfNotExists().
		Column("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		EndColumn().
		RunWith(tx).
		Exec()
	assert.NoError(t, err)
}

func TestCreateIndexBuilder_ToSql(t *testing.T) {
	actual, _, err := CreateIndexBuilder(builder.EmptyBuilder).
		Name("foo").
		IfNotExists().
		On("foobar").
		Columns("bar", "baz").
		ToSql()
	assert.NoError(t, err)
	expected := "CREATE INDEX IF NOT EXISTS foo ON foobar (bar, baz)"
	assert.Equal(t, expected, actual)

	actual, _, err = CreateIndexBuilder(builder.EmptyBuilder).
		Name("foo").
		On("foobar").
		Columns("bar").
		ToSql()
	assert.NoError(t, err)
	expected = "CREATE INDEX foo ON foobar (bar)"
	assert.Equal(t, expected, actual)
}

func columnBuilder(mapping map[ColumnType]string) ColumnBuilder {
	return ColumnBuilder(builder.EmptyBuilder).columnTypeNames(mapping)
}

func tableBuilder(mapping map[ColumnType]string) CreateTableBuilder {
	return CreateTableBuilder(builder.EmptyBuilder).columnTypeNames(mapping)
}
