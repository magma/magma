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

package servicers_test

import (
	"testing"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/gateway/mconfig"

	"github.com/stretchr/testify/assert"
)

func TestSwxProxyMultipleConfigurationMconfig(t *testing.T) {
	confs := generateSwxProxyConfigFromString(t, multipleServersMconfig)
	assertHLRClients(t, confs)
	// server on "server" tag should not appear
	assert.Equal(t, 2, len(confs))
	assert.Equal(t, "10.0.0.1:1", confs[0].ServerCfg.DiameterServerConnConfig.Addr)
	assert.Equal(t, "10.0.0.2:2", confs[1].ServerCfg.DiameterServerConnConfig.Addr)
	assert.Equal(t, "magma_test1", confs[0].ClientCfg.ProductName)
	assert.Equal(t, "magma_test2", confs[1].ClientCfg.ProductName)
}

// TODO: remove  once backwards compatibility is not needed for the field server
func TestSwxProxyLegacyConfigurationMconfig(t *testing.T) {
	confs := generateSwxProxyConfigFromString(t, legacyServerMconfigGen)
	assertHLRClients(t, confs)
	// server on "server" tag should not appear
	assert.Equal(t, 1, len(confs))
	assert.Equal(t, "10.0.0.0:0", confs[0].ServerCfg.DiameterServerConnConfig.Addr)
	assert.Equal(t, "magma_test0", confs[0].ClientCfg.ProductName)

}

func TestSwxProxyService_ValidateConfig(t *testing.T) {
	err := servicers.ValidateSwxProxyConfig(nil)
	assert.EqualError(t, err, "Nil SwxProxyConfig provided")

	validClientConfig := &diameter.DiameterClientConfig{
		Host:  "magma-oai.openair4G.eur", // diameter host
		Realm: "openair4G.eur",           // diameter realm,
	}
	validServerConfig := &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
		Addr:     "",      // to be filled in once server addr is started
		Protocol: "sctp"}, // tcp/sctp
	}
	nilClientConfig := &servicers.SwxProxyConfig{
		ClientCfg:           nil,
		ServerCfg:           validServerConfig,
		VerifyAuthorization: false,
	}
	nilServerConfig := &servicers.SwxProxyConfig{
		ClientCfg:           validClientConfig,
		ServerCfg:           nil,
		VerifyAuthorization: false,
	}
	err = servicers.ValidateSwxProxyConfig(nilClientConfig)
	assert.EqualError(t, err, "Nil client config provided")

	err = servicers.ValidateSwxProxyConfig(nilServerConfig)
	assert.EqualError(t, err, "Nil server config provided")

	invalidClientConfig := &servicers.SwxProxyConfig{
		ClientCfg: &diameter.DiameterClientConfig{
			Host:  "",              // diameter host
			Realm: "openair4G.eur", // diameter realm,
		},
		ServerCfg:           validServerConfig,
		VerifyAuthorization: false,
	}
	err = servicers.ValidateSwxProxyConfig(invalidClientConfig)
	assert.EqualError(t, err, "Invalid Diameter Host")

	invalidServerConfig := &servicers.SwxProxyConfig{
		ClientCfg: validClientConfig,
		ServerCfg: &diameter.DiameterServerConfig{DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:     "",     // to be filled in once server addr is started
			Protocol: "sss"}, // tcp/sctp
		},
		VerifyAuthorization: false,
	}
	err = servicers.ValidateSwxProxyConfig(invalidServerConfig)
	assert.EqualError(t, err, "Invalid Diameter Address (sss://): unknown network sss")

	validConfig := &servicers.SwxProxyConfig{
		ClientCfg:           validClientConfig,
		ServerCfg:           validServerConfig,
		VerifyAuthorization: false,
	}
	err = servicers.ValidateSwxProxyConfig(validConfig)
	assert.NoError(t, err)
}

func generateSwxProxyConfigFromString(t *testing.T, confString string) []*servicers.SwxProxyConfig {
	err := mconfig.CreateLoadTempConfig(confString)
	assert.NoError(t, err)
	// Note we get only get index 0
	return servicers.GetSwxProxyConfig()
}

func assertHLRClients(t *testing.T, confs []*servicers.SwxProxyConfig) {
	for _, cfg := range confs {
		assert.Truef(t, plmn_filter.CheckImsiOnPlmnIdListIfAny("001020000000055", cfg.HlrPlmnIds),
			"IMSI 001020000000055 should be HLR IMSI, HLR PLMN ID Map: %+v", cfg.HlrPlmnIds)
		assert.Truef(t, plmn_filter.CheckImsiOnPlmnIdListIfAny("001030000000055", cfg.HlrPlmnIds),
			"IMSI 001030000000055 should be HLR IMSI, HLR PLMN ID Map: %+v", cfg.HlrPlmnIds)
		assert.Falsef(t, plmn_filter.CheckImsiOnPlmnIdListIfAny("001010000000055", cfg.HlrPlmnIds),
			"IMSI 001010000000055 should NOT be HLR IMSI, HLR PLMN ID Map: %+v", cfg.HlrPlmnIds)
	}
}

// ---- CONFIGURATIONS ----
var (
	// TODO: remove "server" tag once backwards compatibility is not needed for the field server
	multipleServersMconfig = `{
		"configsByKey": {
			"swx_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.SwxConfig",
				"logLevel": "INFO",
				"server": {
						"protocol": "sctp",
						"address": "10.0.0.0:0",
						"retransmits": 3,
						"watchdogInterval": 1,
						"retryCount": 5,
						"productName": "magma_test0",
						"realm": "openair4G.eur",
						"host": "magma-oai.openair4G.eur"
					},
				"servers": [
					{
						"protocol": "sctp",
						"address": "10.0.0.1:1",
						"retransmits": 3,
						"watchdogInterval": 1,
						"retryCount": 5,
						"productName": "magma_test1",
						"realm": "openair4G.eur",
						"host": "magma-oai.openair4G.eur"
					},
					{
						"protocol": "sctp",
						"address": "10.0.0.2:2",
						"retransmits": 3,
						"watchdogInterval": 1,
						"retryCount": 5,
						"productName": "magma_test2",
						"realm": "openair4G.eur",
						"host": "magma-oai.openair4G.eur"
					}
				],
				"verifyAuthorization": true,
				"hlr_plmn_ids": [ "00102", "00103" ]
			}
		}
	}`
	// TODO: remove  once backwards compatibility is not needed for the field server
	legacyServerMconfigGen = `{
		"configsByKey": {
			"swx_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.SwxConfig",
				"logLevel": "INFO",
				"server": {
						"protocol": "sctp",
						"address": "10.0.0.0:0",
						"retransmits": 3,
						"watchdogInterval": 1,
						"retryCount": 5,
						"productName": "magma_test0",
						"realm": "openair4G.eur",
						"host": "magma-oai.openair4G.eur"
					},
				"verifyAuthorization": true,
				"hlr_plmn_ids": [ "00102", "00103" ]
			}
		}
	}`
)
