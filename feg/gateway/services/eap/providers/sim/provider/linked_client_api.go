// +build link_local_service

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

// package sim implements EAP-SIM provider
package provider

import (
	"errors"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
	_ "magma/feg/gateway/services/eap/providers/sim/servicers/handlers"
	managed_configs "magma/gateway/mconfig"
)

func NewService(srvsr *servicers.EapSimSrv) providers.Method {
	return &providerImpl{EapSimSrv: srvsr}
}

// Handle handles passed EAP-SIM payload & returns corresponding result
// this Handle implementation is using GRPC based SIM provider service
func (prov *providerImpl) Handle(msg *protos.Eap) (*protos.Eap, error) {
	if msg == nil {
		return nil, errors.New("Invalid EAP SIM Message")
	}
	prov.RLock()
	if prov.EapSimSrv == nil {
		// servicer is not initialized, relock, recheck, create
		prov.RUnlock()
		prov.Lock()
		if prov.EapSimSrv == nil {
			simConfigs := &mconfig.EapSimConfig{}
			err := managed_configs.GetServiceConfigs(sim.EapSimServiceName, simConfigs)
			if err != nil {
				glog.Errorf("Error getting EAP SIM service configs: %s", err)
				simConfigs = nil
			}
			prov.EapSimSrv, err = servicers.NewEapSimService(simConfigs)
			if err != nil || prov.EapSimSrv == nil {
				glog.Fatalf("failed to create EAP SIM Service: %v", err) // should never happen
			}
		}
		prov.Unlock()
		prov.RLock()
	}
	defer prov.RUnlock()
	return prov.EapSimSrv.HandleImpl(msg)
}
