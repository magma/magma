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

// package aka implements EAP-AKA provider
package provider

import (
	"errors"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
	managed_configs "magma/gateway/mconfig"
)

func NewService(srvsr *servicers.EapAkaSrv) providers.Method {
	return &providerImpl{EapAkaSrv: srvsr}
}

// Handle handles passed EAP-AKA payload & returns corresponding result
// this Handle implementation is using GRPC based AKA provider service
func (prov *providerImpl) Handle(msg *protos.Eap) (*protos.Eap, error) {
	if msg == nil {
		return nil, errors.New("Invalid EAP AKA Message")
	}
	prov.RLock()
	if prov.EapAkaSrv == nil {
		// servicer is not initialized, relock, recheck, create
		prov.RUnlock()
		prov.Lock()
		if prov.EapAkaSrv == nil {
			akaConfigs := &mconfig.EapAkaConfig{}
			err := managed_configs.GetServiceConfigs(aka.EapAkaServiceName, akaConfigs)
			if err != nil {
				glog.Errorf("Error getting EAP AKA service configs: %s", err)
				akaConfigs = nil
			}
			prov.EapAkaSrv, err = servicers.NewEapAkaService(akaConfigs)
			if err != nil || prov.EapAkaSrv == nil {
				glog.Fatalf("failed to create EAP AKA Service: %v", err) // should never happen
			}
		}
		prov.Unlock()
		prov.RLock()
	}
	defer prov.RUnlock()
	return prov.EapAkaSrv.HandleImpl(msg)
}
