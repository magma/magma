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

	"github.com/magma/magma/src/go/capture"
	"github.com/magma/magma/src/go/log"
	pb "github.com/magma/magma/src/go/protos/magma/capture"
)

// CaptureServer handles CaptureServer RPCs.
type CaptureServer struct {
	log.Logger
	pb.CaptureServer
	*capture.Buffer
}

// NewCaptureServer returns a CaptureServer injected with the provided capture.Buffer and logger.
func NewCaptureServer(logger log.Logger, b *capture.Buffer) *CaptureServer {
	return &CaptureServer{
		Logger: logger,
		Buffer: b,
	}
}

// Flush calls buffer.Flush
func (c *CaptureServer) Flush(ctx context.Context, req *pb.FlushRequest) (*pb.FlushResponse, error) {
	c.Logger.
		Debug().Print("Flush")
	calls := c.Buffer.Flush()
	return &pb.FlushResponse{Recording: &pb.Recording{UnaryCalls: calls}}, nil
}
