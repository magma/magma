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
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/feg/cloud/go/protos/mconfig"

	"magma/feg/gateway/services/session_proxy/credit_control"

	managed_configs "magma/gateway/mconfig"
)

func TestGxConfig(t *testing.T) {
	// Create tmp mconfig test file & load configs from it
	fegConfigFmt := `{
		"configsByKey": {
			"session_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.SessionProxyConfig",
				"logLevel": "INFO",
				"gx": {
                                         "disableGx": false,
                                         "server": {
                                                  "protocol": "tcp",
                                                  "address": "1.1.1.1:9999",
                                                  "retransmits": 3,
                                                  "watchdogInterval": 1,
                                                  "retryCount": 5,
                                                  "productName": "magma",
                                                  "realm": "magma.com",
                                                  "host": "magma-fedgw.magma.com"
                                         }
                                 },
                                 "requestFailureThreshold": 0.5,
                                 "minimumRequestThreshold": 1
                         }
		}
	}`

	err := managed_configs.CreateLoadTempConfig(fegConfigFmt)
	assert.NoError(t, err)

	configsPtr := &mconfig.SessionProxyConfig{}
	managed_configs.GetServiceConfigs(credit_control.SessionProxyServiceName, configsPtr)

	gxConfig := configsPtr.GetGx()

	assert.Equal(t, false, gxConfig.DisableGx)
	assert.Equal(t, "tcp", gxConfig.Server.Protocol)
	assert.Equal(t, "1.1.1.1:9999", gxConfig.Server.Address)
	assert.Equal(t, uint32(3), gxConfig.Server.Retransmits)
	assert.Equal(t, uint32(1), gxConfig.Server.WatchdogInterval)
	assert.Equal(t, uint32(5), gxConfig.Server.RetryCount)
	assert.Equal(t, "magma", gxConfig.Server.ProductName)
	assert.Equal(t, "magma.com", gxConfig.Server.Realm)
	assert.Equal(t, "magma-fedgw.magma.com", gxConfig.Server.Host)
	assert.Equal(t, float32(0.5), configsPtr.GetRequestFailureThreshold())
	assert.Equal(t, uint32(1), configsPtr.GetMinimumRequestThreshold())
}
