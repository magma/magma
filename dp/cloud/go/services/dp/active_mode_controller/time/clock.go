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

package time

import "time"

type Clock struct{}

func (*Clock) Now() time.Time {
	return time.Now()
}

func (*Clock) Tick(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}
