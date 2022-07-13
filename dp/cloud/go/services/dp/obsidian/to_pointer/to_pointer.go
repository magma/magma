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

package to_pointer

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Float(x float64) *float64 {
	return &x
}

func Int64(x int64) *int64 {
	return &x
}

func Int(x int) *int {
	return &x
}

func String(x string) *string {
	return &x
}

func Bool(x bool) *bool {
	return &x
}

func TimeToDateTime(t int64) strfmt.DateTime {
	tm := time.Unix(t, 0)
	return strfmt.DateTime(tm)
}

func DoubleValueToFloat(v *wrappers.DoubleValue) *float64 {
	if v == nil {
		return nil
	}
	return &v.Value
}

func BoolValueToBool(v *wrappers.BoolValue) *bool {
	if v == nil {
		return nil
	}
	return &v.Value
}

func StringValueToString(v *wrappers.StringValue) *string {
	if v == nil {
		return nil
	}
	return &v.Value
}

func FloatToDoubleValue(v *float64) *wrappers.DoubleValue {
	if v == nil {
		return nil
	}
	return wrapperspb.Double(*v)
}

func BoolToBoolValue(v *bool) *wrapperspb.BoolValue {
	if v == nil {
		return nil
	}
	return wrapperspb.Bool(*v)
}

func StringToStringValue(v *string) *wrapperspb.StringValue {
	if v == nil {
		return nil
	}
	return wrapperspb.String(*v)
}
