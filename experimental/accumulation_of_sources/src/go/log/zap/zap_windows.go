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

package zap

import (
	"net/url"
	"os"

	"go.uber.org/zap"
)

// newWinFileSink is a workaround from uber-go/zap #621
// TL;DR: raw file paths do not work since zap calls url.Parse; file:///<path>
// is the canonical way to specify a file URL, but this leaves us with a
// `/<path>` path after URL parse, which is invalid on Windows.
func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// Remove leading slash left by url.Parse()
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

func init() {
	zap.RegisterSink("winfile", newWinFileSink)
}
