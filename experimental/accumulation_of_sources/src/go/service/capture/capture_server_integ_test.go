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

package capture

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/magma/magma/src/go/capture"
	"github.com/magma/magma/src/go/internal/testutil"
	capturepb "github.com/magma/magma/src/go/protos/magma/capture"
	"github.com/stretchr/testify/assert"
)

func TestNewCaptureServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := testutil.NewTestLogger()
	buffer := capture.NewBuffer()

	cs := NewCaptureServer(logger, buffer)
	assert.Equal(t, logger, cs.Logger)
	assert.Equal(t, buffer, cs.Buffer)
}

func TestCaptureServer_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger, _ := testutil.NewTestLogger()
	buf := capture.NewBuffer()
	call := &capturepb.UnaryCall{
		Method: "test",
	}
	buf.Write(call)
	cs := NewCaptureServer(logger, buf)
	resp, err := cs.Flush(context.Background(), &capturepb.FlushRequest{})
	assert.NoError(t, err)
	assert.Equal(t, call, resp.GetRecording().GetUnaryCalls()[0])
}
