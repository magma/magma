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

package gy

import (
	"os"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"

	"magma/gateway/mconfig"
)

// Create tmp mconfig test file & load configs from it
var fegConfigFmt = `{
	"configsByKey": {
		"session_proxy": {
			"@type": "type.googleapis.com/magma.mconfig.SessionProxyConfig",
			"gy": {
				"DisableGy": false,
				"OverwriteApn": "apn.magma.com",
				"servers": [{
					"protocol": "sctp",
					"address": "1.1.1.1:9999",
					"retransmits": 5,
					"watchdog_interval": 1,
					"retry_count": 5,
					"local_address": "2.2.2.2:9999",
					"product_name": "magma_test",
					"realm": "local.openair4G.eur",
					"host": "local.magma-oai.openair4G.eur",
					"dest_host": "magma-oai.openair4G.eur",
					"dest_realm": "openair4G.eur",
					"disable_dest_host": true,
					"overwrite_dest_host": false,
					"request_timeout": 10
				}],
				"virtual_apn_rules": [{
					"apn_filter": ".*",
					"charging_characteristics_filter": "1*",
					"apn_overwrite": "vApnGy.magma-fedgw.magma.com"
				}]
			}
		}
	}
}`

func TestGyClientConfig(t *testing.T) {

	os.Setenv(GySupportedVendorIDsEnv, "example-supported-vendor-id")
	defer os.Unsetenv(GySupportedVendorIDsEnv)

	os.Setenv(GyServiceContextIdEnv, "context-id")
	defer os.Unsetenv(GyServiceContextIdEnv)

	err := mconfig.CreateLoadTempConfig(fegConfigFmt)
	assert.NoError(t, err)
	cliConfig := GetGyClientConfiguration()[0]

	assert.Equal(t, uint32(diam.CHARGING_CONTROL_APP_ID), cliConfig.AppID)
	assert.Equal(t, "example-supported-vendor-id", cliConfig.SupportedVendorIDs)
	assert.Equal(t, "context-id", cliConfig.ServiceContextId)
	assert.Equal(t, "local.magma-oai.openair4G.eur", cliConfig.Host)
	assert.Equal(t, "local.openair4G.eur", cliConfig.Realm)
	assert.Equal(t, "magma_test", cliConfig.ProductName)
	assert.Equal(t, uint(1), cliConfig.WatchdogInterval)
	assert.Equal(t, uint(5), cliConfig.RetryCount)
	assert.Equal(t, uint(5), cliConfig.Retransmits)
	assert.Equal(t, uint(10), cliConfig.RequestTimeout)
}

func TestGetGyOCSConfig(t *testing.T) {

	err := mconfig.CreateLoadTempConfig(fegConfigFmt)
	assert.NoError(t, err)
	ocsConfig := GetOCSConfiguration()[0]

	assert.Equal(t, "1.1.1.1:9999", ocsConfig.Addr)
	assert.Equal(t, "sctp", ocsConfig.Protocol)
	assert.Equal(t, "2.2.2.2:9999", ocsConfig.LocalAddr)
	assert.Equal(t, "magma-oai.openair4G.eur", ocsConfig.DestHost)
	assert.Equal(t, "openair4G.eur", ocsConfig.DestRealm)
	assert.Equal(t, bool(true), ocsConfig.DisableDestHost)
	assert.Equal(t, bool(false), ocsConfig.OverwriteDestHost)
}

func TestGyGlobalConfig(t *testing.T) {

	os.Setenv(OCSServiceIdentifierEnv, "example-service-id")
	defer os.Unsetenv(OCSServiceIdentifierEnv)

	os.Setenv(DisableRequestedGrantedUnitsAVPEnv, "false")
	defer os.Unsetenv(DisableRequestedGrantedUnitsAVPEnv)

	err := mconfig.CreateLoadTempConfig(fegConfigFmt)
	assert.NoError(t, err)
	globalConfig := GetGyGlobalConfig()
	vApnRules := globalConfig.VirtualApnRules[0]

	assert.Equal(t, "apn.magma.com", globalConfig.OCSOverwriteApn)
	assert.Equal(t, "example-service-id", globalConfig.OCSServiceIdentifier)
	assert.Equal(t, bool(false), globalConfig.DisableGy)
	assert.Regexp(t, ".*", vApnRules.ApnFilter)
	assert.Regexp(t, "1*", vApnRules.ChargingCharacteristicsFilter)
	assert.Equal(t, "vApnGy.magma-fedgw.magma.com", vApnRules.ApnOverwrite)
	assert.Equal(t, bool(false), globalConfig.DisableServiceGrantedUnitsAVP)
}
