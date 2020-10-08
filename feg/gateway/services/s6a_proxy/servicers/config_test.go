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
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/gateway/mconfig"
)

func TestS6aConfig(t *testing.T) {
	// Create tmp mconfig test file & load configs from it
	fegConfigFmt := `{
		"configsByKey": {
			"s6a_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.S6aConfig",
				"logLevel": "INFO",
				"server": {
					"protocol": "sctp",
					"address": "1.1.1.1:9999",
					"retransmits": 3,
					"watchdog_interval": 1,
					"retry_count": 5,
					"product_name": "magma_test",
					"realm": "local.openair4G.eur",
					"host": "local.magma-oai.openair4G.eur",
					"dest_host":"magma-oai.openair4G.eur",
					"dest_realm":"openair4G.eur"
				}
			}
		}
	}`

	err := mconfig.CreateLoadTempConfig(fegConfigFmt)
	assert.NoError(t, err)
	config := GetS6aProxyConfigs()
	srvConfig := config.ServerCfg
	cliConfig := config.ClientCfg

	assert.Equal(t, "1.1.1.1:9999", srvConfig.Addr)
	assert.Equal(t, "sctp", srvConfig.Protocol)
	assert.Equal(t, "magma-oai.openair4G.eur", srvConfig.DestHost)
	assert.Equal(t, "openair4G.eur", srvConfig.DestRealm)
	assert.Equal(t, "local.magma-oai.openair4G.eur", cliConfig.Host)
	assert.Equal(t, "local.openair4G.eur", cliConfig.Realm)
	assert.Equal(t, uint(1), cliConfig.WatchdogInterval)
	assert.Equal(t, "magma_test", cliConfig.ProductName)
}
