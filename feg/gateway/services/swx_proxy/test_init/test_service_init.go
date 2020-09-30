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

package test_init

import (
	"fmt"
	"os"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/swx_proxy/cache"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/swx_proxy/servicers/test"
	"magma/gateway/mconfig"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/service"
)

func StartTestService(t *testing.T) (*service.Service, error) {
	return StartTestServiceWithCache(t, cache.New())
}

func StartTestServiceWithCache(t *testing.T, cache *cache.Impl) (*service.Service, error) {
	os.Setenv("USE_REMOTE_SWX_PROXY", "false")
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)

	// Note we get only get index 0
	config := servicers.GetSwxProxyConfig()[0]
	serverAddr, err := test.StartTestSwxServer(config.ServerCfg.Protocol, config.ServerCfg.Addr)
	if err != nil {
		return nil, err
	}

	// Update server config with chosen port of swx test server
	config.ServerCfg.Addr = serverAddr
	service, err := servicers.NewSwxProxyWithCache(config, cache)
	if err != nil {
		return nil, err
	}
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)
	return srv, nil
}

func InitTestMconfig(t *testing.T, addr string, verify bool) error {
	// Create tmp mconfig test file & load configs from it
	fegConfigFmt := `{
		"configsByKey": {
			"swx_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.SwxConfig",
				"logLevel": "INFO",
				"servers": [
					{
						"protocol": "sctp",
						"address": "%s",
						"retransmits": 3,
						"watchdogInterval": 1,
						"retryCount": 5,
						"productName": "magma_test",
						"realm": "openair4G.eur",
						"host": "magma-oai.openair4G.eur"
						}
				],
				"verifyAuthorization": %t,
				"hlr_plmn_ids": [ "00102", "00103" ]
			}
		}
	}`
	stringConfig := fmt.Sprintf(fegConfigFmt, addr, verify)
	res := mconfig.CreateLoadTempConfig(stringConfig)

	// Note we get only get index 0
	cfg := servicers.GetSwxProxyConfig()[0]
	if !plmn_filter.CheckImsiOnPlmnIdListIfAny("001020000000055", cfg.HlrPlmnIds) {
		t.Fatalf("IMSI 001020000000055 should be HLR IMSI, HLR PLMN ID Map: %+v", cfg.HlrPlmnIds)
	}
	if !plmn_filter.CheckImsiOnPlmnIdListIfAny("001030000000055", cfg.HlrPlmnIds) {
		t.Fatalf("IMSI 001030000000055 should be HLR IMSI, HLR PLMN ID Map: %+v", cfg.HlrPlmnIds)
	}
	if plmn_filter.CheckImsiOnPlmnIdListIfAny("001010000000055", cfg.HlrPlmnIds) {
		t.Fatalf("IMSI 001010000000055 should NOT be HLR IMSI, HLR PLMN ID Map: %+v", cfg.HlrPlmnIds)
	}

	return res
}
