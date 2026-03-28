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

package servicers

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"

	"magma/dp/cloud/go/services/dp/storage/db"
)

func dbFloat64OrNil(val *wrappers.DoubleValue) sql.NullFloat64 {
	if val == nil {
		return sql.NullFloat64{}
	}
	return db.MakeFloat(val.Value)
}

func dbBoolOrNil(val *wrappers.BoolValue) sql.NullBool {
	if val == nil {
		return sql.NullBool{}
	}
	return db.MakeBool(val.Value)
}

func dbStringOrNil(val *wrappers.StringValue) sql.NullString {
	if val == nil {
		return sql.NullString{}
	}
	return db.MakeString(val.Value)
}
