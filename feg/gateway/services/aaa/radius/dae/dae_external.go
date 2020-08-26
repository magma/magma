// Copyright 2020 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !with_builtin_radius

// package dae implements Radius Dynamic Authorization Extensions API (https://tools.ietf.org/html/rfc5176)
package dae

import (
	"context"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
)

type extDAEServer struct{}

// NewDAEServicer returns default servicer using external, registry based service
func NewDAEServicer(*mconfig.RadiusConfig) DAE {
	return extDAEServer{}
}

// Disconnect is DAE's Disconnect Messages equivalent
func (extDAEServer) Disconnect(aaaCtx *protos.Context) error {
	conn, err := registry.GetConnection(registry.RADIUS)
	if err == nil {
		_, err = protos.NewAuthorizationClient(conn).Disconnect(
			context.Background(), &protos.DisconnectRequest{Ctx: aaaCtx})
	}
	return err
}
