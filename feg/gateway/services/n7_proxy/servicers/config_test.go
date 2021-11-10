/*
Copyright 2021 The Magma Authors.

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

	"github.com/stretchr/testify/assert"

	"magma/feg/gateway/services/n7_proxy/servicers"
	"magma/gateway/mconfig"
)

var (
	config = `{
		"configsByKey": {
			"n7_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.N7Config",
				"logLevel": "INFO",
				"disableN7": false,
				"servers": [
					{
						"apiRoot": "https://mockpcf/npcf-smpolicycontrol/v1",
						"tokenUrl": "https://mockpcf/oauth2/token",
						"clientId": "feg_magma_client",
						"clientSecret": "feg_mamga_secret"
					}
				]
			}
		}
	}`
	err_config = `{
		"configsByKey": {
			"n7_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.N7Config",
				"logLevel": "INFO",
				"disableN7": false,
				"servers": [
					{
						"apiRoot": "mockpcf/npcf-smpolicycontrol/v1",
						"tokenUrl": "https://mockpcf/oauth2/token",
						"clientId": "feg_magma_client",
						"clientSecret": "feg_mamga_secret"
					}
				]
			}
		}
	}`
)

func TestGetN7Config(t *testing.T) {
	conf := generateN7Mconfig(t, config)
	assert.NotNil(t, conf)
	assert.Equal(t, 1, len(conf.Servers))
	assert.Equal(t, false, conf.DisableN7)
	assert.Equal(t, "https://mockpcf/npcf-smpolicycontrol/v1", conf.Servers[0].ApiRoot)
	assert.Equal(t, "https://mockpcf/oauth2/token", conf.Servers[0].TokenUrl)
	assert.Equal(t, "feg_magma_client", conf.Servers[0].ClientId)
	assert.Equal(t, "feg_mamga_secret", conf.Servers[0].ClientSecret)
}

func TestInvalidConfig(t *testing.T) {
	conf := generateN7Mconfig(t, err_config)
	assert.Nil(t, conf)
}

func generateN7Mconfig(t *testing.T, configString string) *servicers.N7ProxyConfig {
	err := mconfig.CreateLoadTempConfig(configString)
	assert.NoError(t, err)
	return servicers.GetN7ProxyConfig()
}
