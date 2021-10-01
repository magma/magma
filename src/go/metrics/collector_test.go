// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
)

func TestStartStop(t *testing.T) {
	mClock := clock.NewMock()
	ctrl := gomock.NewController(t)
	mSender := NewMockMetricSender(ctrl)
	defer ctrl.Finish()

	mSender.EXPECT().Send(gomock.Any())
	collector := NewCollector(mClock, time.Second, mSender)
	collector.Start()

	// Pass a second to check that collection is run
	mClock.Add(time.Second)
}
