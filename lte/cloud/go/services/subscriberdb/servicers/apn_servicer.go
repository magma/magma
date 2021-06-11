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

package servicers

import (
	"context"

	lte_protos "magma/lte/cloud/go/protos"
)

type apndbServicer struct{}

func NewApndbServicer() lte_protos.ApnDBCloudServer {
	return &apndbServicer{}
}

//ListApnConfigs returns a page of network wide APN configs
func (s *apndbServicer) ListApnConfigs(ctx context.Context, req *lte_protos.ListApnConfigRequest) (*lte_protos.ListApnConfigResponse, error) {
	//TODO - get the apn configs from configurator and build response
	resp := &lte_protos.ListApnConfigResponse{}
	return resp, nil
}

//ListGatewayApnConfigs returns a page of gateway specific APN configs
func (s *apndbServicer) ListGatewayApnConfigs(ctx context.Context, req *lte_protos.ListGatewayApnConfigRequest) (*lte_protos.ListGatewayApnConfigResponse, error) {
	//TODO - get the gateway apn configs from configurator and build response
	resp := &lte_protos.ListGatewayApnConfigResponse{}
	return resp, nil
}
