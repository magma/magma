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

// package servicers implements EAP-SIM GRPC service
package servicers

import (
	"sync"

	"github.com/golang/glog"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/sim"
)

// Handler - is an SIM Subtype
type Handler func(srvr *EapSimSrv, ctx *protos.Context, req eap.Packet) (eap.Packet, error)

var simHandlers struct {
	rwl sync.RWMutex
	hm  map[sim.Subtype]Handler
}

func AddHandler(st sim.Subtype, h Handler) {
	if h == nil {
		return
	}
	simHandlers.rwl.Lock()
	if simHandlers.hm == nil {
		simHandlers.hm = map[sim.Subtype]Handler{}
	}
	oldh, ok := simHandlers.hm[st]
	if ok && oldh != nil {
		glog.Warningf("EAP SIM Handler for subtype %d => %+v is already registered, will overwrite with %+v",
			st, oldh, h)
	}
	simHandlers.hm[st] = h
	simHandlers.rwl.Unlock()
}

func GetHandler(st sim.Subtype) Handler {
	simHandlers.rwl.RLock()
	defer simHandlers.rwl.RUnlock()
	res, ok := simHandlers.hm[st]
	if ok {
		return res
	}
	return nil
}
