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

	"magma/orc8r/cloud/go/sqorc"
)

func CreateTable(tx sq.BaseRunner, builder sqorc.StatementBuilder, metadata *ModelMetadata) error {
	tableBuilder := builder.CreateTable(metadata.Table).
		IfNotExists().
		RunWith(tx).
		PrimaryKey("id")
	fields := metadata.CreateObject().Fields()
	tableBuilder = addColumns(tableBuilder, fields)
	tableBuilder = addRelations(tableBuilder, metadata)
	_, err := tableBuilder.Exec()
	return err
}

func addColumns(builder sqorc.CreateTableBuilder, fields FieldMap) sqorc.CreateTableBuilder {
	for column, field := range fields {
		colBuilder := builder.
			Column(column).
			Type(field.BaseType.sqlType())
		if !field.Nullable {
			colBuilder = colBuilder.NotNull()
		}
		if field.HasDefault {
			colBuilder = colBuilder.Default(field.DefaultValue)
		}
		builder = colBuilder.EndColumn()

		if field.Unique {
			builder = builder.Unique(column)
		}
	}
	return builder
}

func addRelations(builder sqorc.CreateTableBuilder, metadata *ModelMetadata) sqorc.CreateTableBuilder {
	for table, column := range metadata.Relations {
		builder = builder.ForeignKey(
			table,
			map[string]string{column: "id"},
			sqorc.ColumnOnDeleteCascade,
		)
	}
	return builder
}
