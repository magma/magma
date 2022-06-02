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

package servicers

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func protoDoubleOrNil(field sql.NullFloat64) *wrappers.DoubleValue {
	if field.Valid {
		return wrapperspb.Double(field.Float64)
	}
	return nil
}

func protoBoolOrNil(field sql.NullBool) *wrappers.BoolValue {
	if field.Valid {
		return wrapperspb.Bool(field.Bool)
	}
	return nil
}

func protoStringOrNil(field sql.NullString) *wrappers.StringValue {
	if field.Valid {
		return wrapperspb.String(field.String)
	}
	return nil
}
