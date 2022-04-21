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
	fields := metadata.Properties
	tableBuilder = addColumns(tableBuilder, fields)
	_, err := tableBuilder.Exec()
	return err
}

func addColumns(builder sqorc.CreateTableBuilder, fields []*Field) sqorc.CreateTableBuilder {
	for _, field := range fields {
		colBuilder := builder.
			Column(field.Name).
			Type(field.SqlType)
		if !field.Nullable {
			colBuilder = colBuilder.NotNull()
		}
		if field.HasDefault {
			colBuilder = colBuilder.Default(field.DefaultValue)
		}
		builder = colBuilder.EndColumn()

		if field.Unique {
			builder = builder.Unique(field.Name)
		}
		if field.Relation != "" {
			builder = builder.ForeignKey(
				field.Relation,
				map[string]string{field.Name: "id"},
				sqorc.ColumnOnDeleteCascade,
			)
		}
	}
	return builder
}
