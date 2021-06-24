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
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/lann/builder"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func init() {
	builder.Register(ColumnBuilder{}, createColumnData{})
	builder.Register(CreateTableBuilder{}, createTableData{})
	builder.Register(CreateIndexBuilder{}, createIndexData{})
}

/*
Because we are only supporting psql and mysql right now and the only difference
between those dialects for table creation is the names of column types, we can
use concrete types for the CREATE TABLE builder and column builder.

The parameterized difference between the dialects is stored as a mapping of
column type to name inside the data structure for each builder.
*/

var postgresColumnTypeMap = map[ColumnType]string{
	ColumnTypeText:   "TEXT",
	ColumnTypeInt:    "INTEGER",
	ColumnTypeBigInt: "BIGINT",
	// BYTEA is effectively limited to 1GB
	ColumnTypeBytes: "BYTEA",
	ColumnTypeBool:  "BOOLEAN",
}

var mariaColumnTypeMap = map[ColumnType]string{
	// Mysql won't index TEXT columns, so choose VARCHAR(255) for text type
	ColumnTypeText:   "VARCHAR(255)",
	ColumnTypeInt:    "INT",
	ColumnTypeBigInt: "BIGINT",
	// LONGBLOB stores up to 4GB and the cost is a flat extra 2 bytes of
	// storage over BLOB, which is limited to 64KB
	ColumnTypeBytes: "LONGBLOB",
	ColumnTypeBool:  "BOOLEAN",
}

// ColumnOnDeleteOption is an enum type to specify ON DELETE behavior for
// foreign keys
type ColumnOnDeleteOption uint8

const (
	ColumnOnDeleteDoNothing ColumnOnDeleteOption = iota
	ColumnOnDeleteCascade
	// Fill in other behaviors as needed
)

// ColumnType is an enum type to specify table column types
type ColumnType uint8

const (
	ColumnTypeText ColumnType = iota
	ColumnTypeInt
	ColumnTypeBigInt
	ColumnTypeBytes
	ColumnTypeBool
	// Fill in other types as needed
)

//=============================================================================
// Tables
//=============================================================================

// CreateTableBuilder is a builder for DDL table creation statements.
// This builder is immutable and all operations will return a new instance
// with the requested fields set.
type CreateTableBuilder builder.Builder

// Name sets the name of the table to be created
func (b CreateTableBuilder) Name(name string) CreateTableBuilder {
	return builder.Set(b, "Name", name).(CreateTableBuilder)
}

// IfNotExists sets the table creation to run only if the table does not
// already exist
func (b CreateTableBuilder) IfNotExists() CreateTableBuilder {
	return builder.Set(b, "IfNotExists", true).(CreateTableBuilder)
}

// PrimaryKey specifies columns to create a primary key on
// Note that the column builder also has a PrimaryKey. We do not cross-validate
// primary key constraints, so avoid specifying PKs in multiple places.
func (b CreateTableBuilder) PrimaryKey(columns ...string) CreateTableBuilder {
	return builder.Set(b, "PrimaryKey", columns).(CreateTableBuilder)
}

// ForeignKey adds a foreign key constraint on a table. Note that the column
// builder also supports foreign key constraints for individual columns for
// engines which support it.
// TODO: pull this into a builder
func (b CreateTableBuilder) ForeignKey(on string, columnMap map[string]string, onDelete ColumnOnDeleteOption) CreateTableBuilder {
	return builder.Append(b, "ForeignKeys", foreignKey{Table: on, ColumnMapping: columnMap, OnDelete: onDelete}).(CreateTableBuilder)
}

// Unique adds a unique constraint on a set of columns.
func (b CreateTableBuilder) Unique(columns ...string) CreateTableBuilder {
	return builder.Append(b, "Unique", columns).(CreateTableBuilder)
}

// RunWith sets the runner for the statement.
func (b CreateTableBuilder) RunWith(runner squirrel.BaseRunner) CreateTableBuilder {
	return builder.Set(b, "RunWith", runner).(CreateTableBuilder)
}

// Column returns a ColumnBuilder to build a column for the table. The returned
// ColumnBuilder will have a reference back to this CreateTableBuilder which
// will be returned when EndColumn() is called so you can chain column creation
// into table creation.
func (b CreateTableBuilder) Column(name string) ColumnBuilder {
	val, _ := builder.Get(b, "ColumnTypeNames")
	return ColumnBuilder(builder.EmptyBuilder).
		parentBuilder(&b).
		columnTypeNames(val.(map[ColumnType]string)).
		Name(name)
}

// Exec runs the statement using the set runner.
func (b CreateTableBuilder) Exec() (sql.Result, error) {
	d := builder.GetStruct(b).(createTableData)
	return d.Exec()
}

// ToSql returns the SQL string and arguments for the statement.
func (b CreateTableBuilder) ToSql() (string, []interface{}, error) {
	d := builder.GetStruct(b).(createTableData)
	return d.ToSql()
}

// unexported internal function for multi-dialect support
func (b CreateTableBuilder) columnTypeNames(m map[ColumnType]string) CreateTableBuilder {
	return builder.Set(b, "ColumnTypeNames", m).(CreateTableBuilder)
}

//=============================================================================
// Columns
//=============================================================================

// ColumnBuilder is a builder for columns within a table creation statement
// This builder is immutable and all methods will return a new instance of the
// builder with the requested fields set.
type ColumnBuilder builder.Builder

// Name sets the name of the column.
func (b ColumnBuilder) Name(name string) ColumnBuilder {
	return builder.Set(b, "Name", name).(ColumnBuilder)
}

// Type sets the type of the column.
func (b ColumnBuilder) Type(columnType ColumnType) ColumnBuilder {
	return builder.Set(b, "Type", &columnType).(ColumnBuilder)
}

// NotNull marks the column as not nullable.
func (b ColumnBuilder) NotNull() ColumnBuilder {
	return builder.Set(b, "NotNull", true).(ColumnBuilder)
}

// Default sets the default value for the column. This value is not escaped so
// SQL expressions are valid.
func (b ColumnBuilder) Default(value interface{}) ColumnBuilder {
	return builder.Set(b, "Default", value).(ColumnBuilder)
}

// PrimaryKey marks the column as a PK for the table.
func (b ColumnBuilder) PrimaryKey() ColumnBuilder {
	return builder.Set(b, "PrimaryKey", true).(ColumnBuilder)
}

// References marks the column as a foreign key to the specified table and
// foreign column.
func (b ColumnBuilder) References(table string, column string) ColumnBuilder {
	return builder.Set(b, "References", &[2]string{table, column}).(ColumnBuilder)
}

// OnDelete sets the deletion behavior for a foreign key column.
func (b ColumnBuilder) OnDelete(onDelete ColumnOnDeleteOption) ColumnBuilder {
	return builder.Set(b, "OnDelete", &onDelete).(ColumnBuilder)
}

// EndColumn returns the parent CreateTableBuilder to continue building the
// table creation statement.
func (b ColumnBuilder) EndColumn() CreateTableBuilder {
	pb, _ := builder.Get(b, "ParentBuilder")
	return builder.Append(*(pb.(*CreateTableBuilder)), "Columns", &b).(CreateTableBuilder)
}

// ToSql returns the column creation as a SQL string.
func (b ColumnBuilder) ToSql() (string, error) {
	d := builder.GetStruct(b).(createColumnData)
	return d.ToSql()
}

// unexported internal methods for multi-dialect support

func (b ColumnBuilder) parentBuilder(pb *CreateTableBuilder) ColumnBuilder {
	return builder.Set(b, "ParentBuilder", pb).(ColumnBuilder)
}

func (b ColumnBuilder) columnTypeNames(m map[ColumnType]string) ColumnBuilder {
	return builder.Set(b, "ColumnTypeNames", m).(ColumnBuilder)
}

//=============================================================================
// Indexes
//=============================================================================

// CreateIndexBuilder is a builder for CREATE INDEX statements
type CreateIndexBuilder builder.Builder

// Name sets the name of the index
func (b CreateIndexBuilder) Name(name string) CreateIndexBuilder {
	return builder.Set(b, "Name", name).(CreateIndexBuilder)
}

// IfNotExists sets the index creation to run only if it doesn't already exist
func (b CreateIndexBuilder) IfNotExists() CreateIndexBuilder {
	return builder.Set(b, "IfNotExists", true).(CreateIndexBuilder)
}

// On sets the table for the index
func (b CreateIndexBuilder) On(table string) CreateIndexBuilder {
	return builder.Set(b, "Table", table).(CreateIndexBuilder)
}

// Columns sets the columns the index is on
func (b CreateIndexBuilder) Columns(columns ...string) CreateIndexBuilder {
	return builder.Set(b, "Columns", columns).(CreateIndexBuilder)
}

// RunWith sets the runner to Exec the statement with
func (b CreateIndexBuilder) RunWith(runner squirrel.BaseRunner) CreateIndexBuilder {
	return builder.Set(b, "RunWith", runner).(CreateIndexBuilder)
}

// Exec runs the statement using the set runner
func (b CreateIndexBuilder) Exec() (sql.Result, error) {
	d := builder.GetStruct(b).(createIndexData)
	return d.Exec()
}

// ToSql returns the sql string and args for to the statement
func (b CreateIndexBuilder) ToSql() (string, []interface{}, error) {
	d := builder.GetStruct(b).(createIndexData)
	return d.ToSql()
}

//=============================================================================
// Builder data types
//=============================================================================

type foreignKey struct {
	Table         string
	ColumnMapping map[string]string
	OnDelete      ColumnOnDeleteOption
}

type createTableData struct {
	RunWith squirrel.BaseRunner
	// this should be set by the statement builder entry point
	ColumnTypeNames map[ColumnType]string

	Name        string
	IfNotExists bool
	Columns     []*ColumnBuilder

	// Constraints
	PrimaryKey  []string
	ForeignKeys []foreignKey
	Unique      [][]string
}

func (d createTableData) Exec() (sql.Result, error) {
	if d.RunWith == nil {
		return nil, squirrel.RunnerNotSet
	}
	return squirrel.ExecWith(d.RunWith, d)
}

func (d createTableData) ToSql() (string, []interface{}, error) {
	if funk.IsEmpty(d.Name) {
		return "", nil, errors.New("table name must be specified")
	}
	if funk.IsEmpty(d.Columns) {
		return "", nil, errors.New("table must specify at least one column")
	}
	if funk.IsEmpty(d.ColumnTypeNames) {
		return "", nil, errors.New("column type-name mapping is empty")
	}

	sb := strings.Builder{}
	sb.WriteString("CREATE TABLE ")
	if d.IfNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(d.Name)
	sb.WriteString(" (\n")

	colSqls := make([]string, 0, len(d.Columns))
	for i, col := range d.Columns {
		colSql, err := col.ToSql()
		if err != nil {
			return "", nil, errors.Wrapf(err, "could not sqlize column at position %d", i)
		}
		colSqls = append(colSqls, colSql)
	}
	sb.WriteString(strings.Join(colSqls, ",\n"))

	if !funk.IsEmpty(d.PrimaryKey) {
		sb.WriteString(",\n")
		sb.WriteString(fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(d.PrimaryKey, ", ")))
	}

	for _, fk := range d.ForeignKeys {
		fkCols := funk.Keys(fk.ColumnMapping).([]string)
		sort.Strings(fkCols)
		refdCols := funk.Map(fkCols, func(s string) string { return fk.ColumnMapping[s] }).([]string)

		sb.WriteString(",\n")
		sb.WriteString("FOREIGN KEY (")
		sb.WriteString(strings.Join(fkCols, ", "))
		sb.WriteString(") REFERENCES ")
		sb.WriteString(fk.Table)
		sb.WriteString(" (")
		sb.WriteString(strings.Join(refdCols, ", "))
		sb.WriteString(")")

		switch fk.OnDelete {
		case ColumnOnDeleteDoNothing:
			break
		case ColumnOnDeleteCascade:
			sb.WriteString(" ON DELETE CASCADE")
		default:
			return "", nil, errors.Errorf("unrecognized on delete behavior %v", fk.OnDelete)
		}
	}

	for _, uniqConstr := range d.Unique {
		sb.WriteString(",\n")
		sb.WriteString(fmt.Sprintf("UNIQUE (%s)", strings.Join(uniqConstr, ", ")))
	}
	sb.WriteString("\n)")
	return sb.String(), []interface{}{}, nil
}

type createColumnData struct {
	// These should be set by the table builder when it spawns a column builder
	ParentBuilder   *CreateTableBuilder
	ColumnTypeNames map[ColumnType]string

	Name    string
	Type    *ColumnType
	NotNull bool
	Default interface{}

	// Foreign key options
	References *[2]string
	OnDelete   *ColumnOnDeleteOption

	PrimaryKey bool
}

func (d createColumnData) ToSql() (string, error) {
	if funk.IsEmpty(d.Name) {
		return "", errors.New("column name must be specified")
	}
	if d.Type == nil {
		return "", errors.New("column type must be specified")
	}
	if _, typeOk := d.ColumnTypeNames[*d.Type]; !typeOk {
		return "", errors.Errorf("column type %v not recognized", *d.Type)
	}
	if d.References != nil && (len(d.References[0]) == 0 || len(d.References[1]) == 0) {
		return "", errors.Errorf("reference table name and column of foreign key must not be empty")
	}
	if d.OnDelete != nil && d.References == nil {
		return "", errors.Errorf("cannot specify an ON DELETE without a REFERENCES")
	}

	sb := strings.Builder{}
	sb.WriteString(d.Name)
	sb.WriteString(" ")
	sb.WriteString(d.ColumnTypeNames[*d.Type])

	if d.PrimaryKey {
		sb.WriteString(" PRIMARY KEY")
	}

	if d.NotNull {
		sb.WriteString(" NOT NULL")
	}
	if d.Default != nil {
		sb.WriteString(fmt.Sprintf(" DEFAULT %v", d.Default))
	}
	if d.References != nil {
		sb.WriteString(fmt.Sprintf(" REFERENCES %s (%s)", d.References[0], d.References[1]))
	}
	if d.OnDelete != nil {
		switch *d.OnDelete {
		case ColumnOnDeleteCascade:
			sb.WriteString(" ON DELETE CASCADE")
		default:
			return "", errors.Errorf("unrecognized on delete option %v", *d.OnDelete)
		}
	}
	return sb.String(), nil
}

type createIndexData struct {
	RunWith squirrel.BaseRunner

	Name        string
	IfNotExists bool
	Table       string
	Columns     []string
}

func (d createIndexData) Exec() (sql.Result, error) {
	if d.RunWith == nil {
		return nil, squirrel.RunnerNotSet
	}
	return squirrel.ExecWith(d.RunWith, d)
}

func (d createIndexData) ToSql() (string, []interface{}, error) {
	if funk.IsEmpty(d.Name) {
		return "", nil, errors.New("index name must be specified")
	}
	if funk.IsEmpty(d.Table) {
		return "", nil, errors.New("index table must be specified")
	}
	if funk.IsEmpty(d.Columns) {
		return "", nil, errors.New("index columns must be specified")
	}

	sb := strings.Builder{}
	sb.WriteString("CREATE INDEX ")
	if d.IfNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(d.Name)
	sb.WriteString(" ON ")
	sb.WriteString(d.Table)
	sb.WriteString(" (")
	sb.WriteString(strings.Join(d.Columns, ", "))
	sb.WriteString(")")
	return sb.String(), []interface{}{}, nil
}
