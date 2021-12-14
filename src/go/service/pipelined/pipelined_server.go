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

package pipelined

import (
	"context"

	"github.com/magma/magma/src/go/log"
	pb "github.com/magma/magma/src/go/protos/magma/pipelined"
	"github.com/magma/magma/src/go/protos/magma/session_manager"
)

// PipelinedServer handles PipelinedServer RPCs by proxying RPCs to
// local OpenFlow calls
type PipelinedServer struct {
	// TODO: add OF client
	log.Logger
	*pb.UnimplementedPipelinedServer
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
