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

package testutil

import (
	"bytes"

	uber_zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"github.com/magma/magma/src/go/log"
	"github.com/magma/magma/src/go/log/zap"
)

// NewTestLogger returns a logger that logs to the returned buffer. Timestamps,
// caller, and stacktraces are disabled to make it easier to validate log
// output. Test coverage for timestamp, caller, and stacktrace output can be
// found in github.com/magma/magma/src/go/log/zap/zap_integ_test.go.
func NewTestLogger() (log.Logger, *bytes.Buffer) {
	ec := uber_zap.NewDevelopmentEncoderConfig()
	ec.TimeKey = ""
	enc := zapcore.NewConsoleEncoder(ec)
	buf := &zaptest.Buffer{}
	l := zap.New(enc, buf, log.DebugLevel)
	return l, &buf.Buffer
}
