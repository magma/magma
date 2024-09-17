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
	"encoding/json"
	"fmt"
)

type BaseType interface {
	value() any
	ptr() any
	isNull() bool
}

type IntType struct{ X *sql.NullInt64 }

func (x IntType) value() any   { return *x.X }
func (x IntType) ptr() any     { return x.X }
func (x IntType) isNull() bool { return !x.X.Valid }

type FloatType struct{ X *sql.NullFloat64 }

func (x FloatType) value() any   { return *x.X }
func (x FloatType) ptr() any     { return x.X }
func (x FloatType) isNull() bool { return !x.X.Valid }

type StringType struct{ X *sql.NullString }

func (x StringType) value() any   { return *x.X }
func (x StringType) ptr() any     { return x.X }
func (x StringType) isNull() bool { return !x.X.Valid }

type BoolType struct{ X *sql.NullBool }

func (x BoolType) value() any   { return *x.X }
func (x BoolType) ptr() any     { return x.X }
func (x BoolType) isNull() bool { return !x.X.Valid }

type TimeType struct{ X *sql.NullTime }

func (x TimeType) value() any   { return *x.X }
func (x TimeType) ptr() any     { return x.X }
func (x TimeType) isNull() bool { return !x.X.Valid }

type JsonType struct{ X any }

func (x JsonType) value() any {
	b, _ := json.Marshal(x.X)
	return b
}
func (x JsonType) ptr() any     { return jsonScanner{x: x.X} }
func (x JsonType) isNull() bool { return x.X == nil }

type jsonScanner struct{ x any }

func (j jsonScanner) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j.x)
	case string:
		return json.Unmarshal([]byte(v), j.x)
	default:
		return fmt.Errorf("unexpected type: %t", v)
	}
}
