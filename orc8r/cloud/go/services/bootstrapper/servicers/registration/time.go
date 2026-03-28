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

package registration

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// GetTime turns Timestamp into Time
func GetTime(timestamp *timestamp.Timestamp) time.Time {
	return time.Unix(timestamp.Seconds, int64(timestamp.Nanos))
}

// GetTimestamp turns Time into Timestamp
func GetTimestamp(time time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.Nanosecond()),
	}
}
