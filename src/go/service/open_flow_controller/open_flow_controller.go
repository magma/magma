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

package open_flow_controller

import (
	"context"

	"github.com/magma/magma/src/go/log"
	pb "github.com/magma/magma/src/go/protos/magma/pipelined"
	"github.com/magma/magma/src/go/protos/magma/session_manager"
)

// OpenFlowControllerServer runs an OpenFlowController and provides a gRPC
// interface for interacting with it
type OpenFlowControllerServer struct {
	log.Logger
	*pb.U
}

// NewPipelinedServer returns a PipelinedServer injected with the provided logger
func NewPipelinedServer(logger log.Logger) *PipelinedServer {
	return &PipelinedServer{Logger: logger}
}

// GetStats returns a RuleRecordTable filtering records based on cookie and cookie mask request parameters
func (p *PipelinedServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*session_manager.RuleRecordTable, error) {
	p.Logger.
		With("cookie", req.GetCookie()).
		With("cookie_mask", req.GetCookieMask()).
		Debug().Print("GetStats")
	return &session_manager.RuleRecordTable{}, nil
}
