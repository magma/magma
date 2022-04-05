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

package gx

import (
	"os"
	"testing"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"

	"magma/gateway/mconfig"
)

func TestGxClientConfig(t *testing.T) {
	// Create tmp mconfig test file & load configs from it
	fegConfigFmt := `{
		"configsByKey": {
			"session_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.SessionProxyConfig",
				"gx": {
					"servers": [{
						"protocol": "sctp",
						"address": "1.1.1.1:9999",
						"retransmits": 5,
						"watchdog_interval": 1,
						"retry_count": 5,
						"product_name": "magma_test",
						"realm": "local.openair4G.eur",
						"host": "local.magma-oai.openair4G.eur",
						"dest_host": "magma-oai.openair4G.eur",
						"dest_realm": "openair4G.eur",
						"request_timeout": 10
					}]
				}
			}
		}
	}`

	os.Setenv(GxSupportedVendorIDsEnv, "example-vendor-id")
	defer os.Unsetenv(GxSupportedVendorIDsEnv)

	err := mconfig.CreateLoadTempConfig(fegConfigFmt)
	assert.NoError(t, err)
	cliConfig := GetGxClientConfiguration()[0]

	assert.Equal(t, uint32(diam.GX_CHARGING_CONTROL_APP_ID), cliConfig.AppID)
	assert.Equal(t, "example-vendor-id", cliConfig.SupportedVendorIDs)
	assert.Equal(t, "local.magma-oai.openair4G.eur", cliConfig.Host)
	assert.Equal(t, "local.openair4G.eur", cliConfig.Realm)
	assert.Equal(t, "magma_test", cliConfig.ProductName)
	assert.Equal(t, uint(1), cliConfig.WatchdogInterval)
	assert.Equal(t, uint(5), cliConfig.RetryCount)
	assert.Equal(t, uint(5), cliConfig.Retransmits)
	assert.Equal(t, uint(10), cliConfig.RequestTimeout)
}
