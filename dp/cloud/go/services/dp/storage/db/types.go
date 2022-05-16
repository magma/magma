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
)

type BaseType interface {
	value() interface{}
	ptr() interface{}
	isNull() bool
}

type IntType struct{ X *sql.NullInt64 }

func (x IntType) value() interface{} { return *x.X }
func (x IntType) ptr() interface{}   { return x.X }
func (x IntType) isNull() bool       { return !x.X.Valid }

type FloatType struct{ X *sql.NullFloat64 }

func (x FloatType) value() interface{} { return *x.X }
func (x FloatType) ptr() interface{}   { return x.X }
func (x FloatType) isNull() bool       { return !x.X.Valid }

type StringType struct{ X *sql.NullString }

func (x StringType) value() interface{} { return *x.X }
func (x StringType) ptr() interface{}   { return x.X }
func (x StringType) isNull() bool       { return !x.X.Valid }

type BoolType struct{ X *sql.NullBool }

func (x BoolType) value() interface{} { return *x.X }
func (x BoolType) ptr() interface{}   { return x.X }
func (x BoolType) isNull() bool       { return !x.X.Valid }

type TimeType struct{ X *sql.NullTime }

func (x TimeType) value() interface{} { return *x.X }
func (x TimeType) ptr() interface{}   { return x.X }
func (x TimeType) isNull() bool       { return !x.X.Valid }
