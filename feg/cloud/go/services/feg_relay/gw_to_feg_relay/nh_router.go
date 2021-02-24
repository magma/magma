/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// gw_to_feg_relay is h2c & GRPC server serving requests from AGWs to FeG
package gw_to_feg_relay

import (
	"sync"

	"google.golang.org/grpc"

	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

const (
	// Minimal length of PLMNID
	MinPlmnIdLen = 5
	MaxPlmnIdLen = 6
)

type connKey struct {
	service gateway_registry.GwServiceType
	addr    string
}

// Router is a service maintaining mapping and connections from Access Gateways to FeGs
type Router struct {
	sync.RWMutex
	connCache map[connKey]*grpc.ClientConn
}

// NewRouter returns a new instance of Gw to FeG router
func NewRouter() *Router {
	return &Router{connCache: map[connKey]*grpc.ClientConn{}}
}
