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
	"database/sql"

	"magma/orc8r/cloud/go/sqorc"
)

type Field interface {
	Value() interface{}
	Ptr() interface{}
	SqlType() sqorc.ColumnType
	IsNullable() bool
}

type IntField struct{ X *sql.NullInt64 }

func (x IntField) Value() interface{}        { return x.X.Int64 }
func (x IntField) Ptr() interface{}          { return x.X }
func (x IntField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeInt }
func (x IntField) IsNullable() bool          { return false }

type FloatField struct{ X *sql.NullFloat64 }

func (x FloatField) Value() interface{}        { return x.X.Float64 }
func (x FloatField) Ptr() interface{}          { return x.X }
func (x FloatField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeReal }
func (x FloatField) IsNullable() bool          { return false }

type StringField struct{ X *sql.NullString }

func (x StringField) Value() interface{}        { return x.X.String }
func (x StringField) Ptr() interface{}          { return x.X }
func (x StringField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeText }
func (x StringField) IsNullable() bool          { return false }

type BoolField struct{ X *sql.NullBool }

func (x BoolField) Value() interface{}        { return x.X.Bool }
func (x BoolField) Ptr() interface{}          { return x.X }
func (x BoolField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeBool }
func (x BoolField) IsNullable() bool          { return false }

type TimeField struct{ X *sql.NullTime }

func (x TimeField) Value() interface{}        { return x.X.Time }
func (x TimeField) Ptr() interface{}          { return x.X }
func (x TimeField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeDatetime }
func (x TimeField) IsNullable() bool          { return false }

type NullIntField struct{ X *sql.NullInt64 }

func (x NullIntField) Value() interface{}        { return *x.X }
func (x NullIntField) Ptr() interface{}          { return x.X }
func (x NullIntField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeInt }
func (x NullIntField) IsNullable() bool          { return true }

type NullFloatField struct{ X *sql.NullFloat64 }

func (x NullFloatField) Value() interface{}        { return *x.X }
func (x NullFloatField) Ptr() interface{}          { return x.X }
func (x NullFloatField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeReal }
func (x NullFloatField) IsNullable() bool          { return true }

type NullStringField struct{ X *sql.NullString }

func (x NullStringField) Value() interface{}        { return *x.X }
func (x NullStringField) Ptr() interface{}          { return x.X }
func (x NullStringField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeText }
func (x NullStringField) IsNullable() bool          { return true }

type NullBoolField struct{ X *sql.NullBool }

func (x NullBoolField) Value() interface{}        { return *x.X }
func (x NullBoolField) Ptr() interface{}          { return x.X }
func (x NullBoolField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeBool }
func (x NullBoolField) IsNullable() bool          { return true }

type NullTimeField struct{ X *sql.NullTime }

func (x NullTimeField) Value() interface{}        { return *x.X }
func (x NullTimeField) Ptr() interface{}          { return x.X }
func (x NullTimeField) SqlType() sqorc.ColumnType { return sqorc.ColumnTypeDatetime }
func (x NullTimeField) IsNullable() bool          { return true }
