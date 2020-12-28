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

// Package sim implements EAP-SIM provider
package provider

import (
	"sync"

	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
)

// SIM Provider Implementation
type providerImpl struct {
	sync.RWMutex
	*servicers.EapSimSrv
}

func New() providers.Method {
	return &providerImpl{}
}

// String returns EAP SIM Provider name/info
func (*providerImpl) String() string {
	return "EAP-SIM"
}

// EAPType returns EAP SIM Type - 18
func (*providerImpl) EAPType() uint8 {
	return sim.TYPE
}
