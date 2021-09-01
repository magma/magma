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

package service

import (
	"github.com/magma/magma/protos/magma/sctpd"
)

// A Router is responsible for providing clients, where the connection to
// server implementations can be configured to be in-process or grpc.ClientConn
// targets.
type Router interface {
	// SctpdDownlinkClient provides a sctpd.SctpdDownlinkClient
	SctpdDownlinkClient() sctpd.SctpdDownlinkClient
	// SctpdUplinkClient provides a sctpd.SctpdUplinkClient
	SctpdUplinkClient() sctpd.SctpdUplinkClient
}

//go:generate go run github.com/golang/mock/mockgen -destination mock_service/mock_router.go . Router

type router struct {
	sctpdDownlinkClient sctpd.SctpdDownlinkClient
	sctpdUplinkClient   sctpd.SctpdUplinkClient
}

// NewRouter returns a Router with the provided clients.
func NewRouter(
	sctpdDownlinkClient sctpd.SctpdDownlinkClient,
	sctpdUplinkClient sctpd.SctpdUplinkClient,
) Router {
	return &router{
		sctpdDownlinkClient: sctpdDownlinkClient,
		sctpdUplinkClient:   sctpdUplinkClient,
	}
}

func (d *router) SctpdDownlinkClient() sctpd.SctpdDownlinkClient {
	return d.sctpdDownlinkClient
}

func (d *router) SctpdUplinkClient() sctpd.SctpdUplinkClient {
	return d.sctpdUplinkClient
}
