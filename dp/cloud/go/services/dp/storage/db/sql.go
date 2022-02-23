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
	"time"
)

func MakeInt(x int64) sql.NullInt64 {
	return sql.NullInt64{Int64: x, Valid: true}
}

func MakeFloat(x float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: x, Valid: true}
}

func MakeString(x string) sql.NullString {
	return sql.NullString{String: x, Valid: true}
}

func MakeBool(x bool) sql.NullBool {
	return sql.NullBool{Bool: x, Valid: true}
}

func MakeTime(x time.Time) sql.NullTime {
	return sql.NullTime{Time: x, Valid: true}
}
