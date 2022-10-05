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

package pipelined_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/magma/magma/src/go/internal/testutil"
	pb "github.com/magma/magma/src/go/protos/magma/pipelined"
	"github.com/magma/magma/src/go/protos/magma/session_manager"
	"github.com/magma/magma/src/go/service/pipelined"
	"github.com/stretchr/testify/assert"
)

func TestPipelinedServer_GetStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &pb.GetStatsRequest{
		Cookie:     1,
		CookieMask: 0,
	}
	expect := &session_manager.RuleRecordTable{}

	logger, logBuffer := testutil.NewTestLogger()

	ps := pipelined.NewPipelinedServer(logger)
	got, err := ps.GetStats(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, got)
	assert.Equal(
		t,
		"DEBUG\tGetStats\t{\"cookie\": 1, \"cookie_mask\": 0}\n",
		logBuffer.String())
}
