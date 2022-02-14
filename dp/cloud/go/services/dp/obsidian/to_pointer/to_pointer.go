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
)

func Float(x float64) *float64 {
	return &x
}

func Int(x int64) *int64 {
	return &x
}

func Str(x string) *string {
	return &x
}

func Bool(x bool) *bool {
	return &x
}

func TimeToDateTime(t int64) *strfmt.DateTime {
	tm := time.Unix(t, 0)
	dt := strfmt.DateTime(tm)
	return &dt
}

func TimeMilliToDate(t int64) *strfmt.DateTime {
	const milli = 1000
	sec, msec := t/milli, t%milli
	tm := time.Unix(sec, msec/int64(time.Millisecond))
	dt := strfmt.DateTime(tm)
	return &dt
}
