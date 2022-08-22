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

package n7_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/gateway/mconfig"
)

const (
	URL1            = "https://mockpcf/npcf-smpolicycontrol/v1"
	TOKEN_URL       = "https://mockpcf/oauth2/token"
	CLIENT_ID       = "feg_magma_client"
	CLIENT_SECRET   = "feg_magma_secret"
	LOCAL_ADDR      = "127.0.0.1:10100"
	NOTIFY_API_ROOT = "https://magma-feg.magam.com/npcf-smpolicycontrol/v1"
)

var (
	config = `{
		"configsByKey": {
			"n7_n40_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.N7N40ProxyConfig",
				"logLevel": "INFO",
				"n7_config": {
					"disableN7": false,
					"server": {
							"apiRoot": "https://mockpcf/npcf-smpolicycontrol/v1",
							"tokenUrl": "https://mockpcf/oauth2/token",
							"clientId": "feg_magma_client",
							"clientSecret": "feg_magma_secret"
					},
					"client": {
						"local_addr": "127.0.0.1:10100",
						"notify_api_root": "https://magma-feg.magam.com/npcf-smpolicycontrol/v1"
					}
				}
			}
		}
	}`
	err_config = `{
		"configsByKey": {
			"n7_n40_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.N7N40ProxyConfig",
				"logLevel": "INFO",
				"n7_config": {
					"disableN7": false,
					"server": {
						"apiRoot": "mockpcf/npcf-smpolicycontrol/v1",
						"tokenUrl": "https://mockpcf/oauth2/token",
						"clientId": "feg_magma_client",
						"clientSecret": "feg_magma_secret"
					},
					"client": {
						"local_addr": "127.0.0.1:10100",
						"notify_api_root": "https://magma-feg.magam.com/npcf-smpolicycontrol/v1"
					}
				}
			}
		}
	}`
	empty_config = `{
		"configsByKey": {
			"n7_n40_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.N7N40ProxyConfig"
			}
		}
	}`
)

func TestGetN7Config(t *testing.T) {
	conf, err := generateN7Mconfig(t, config)
	require.NoError(t, err)
	assert.Equal(t, false, conf.DisableN7)
	url1, _ := url.ParseRequestURI(URL1)
	assert.Equal(t, *url1, conf.ServerConfig.ApiRoot)
	assert.Equal(t, TOKEN_URL, conf.ServerConfig.TokenUrl)
	assert.Equal(t, CLIENT_ID, conf.ServerConfig.ClientId)
	assert.Equal(t, CLIENT_SECRET, conf.ServerConfig.ClientSecret)
	assert.Equal(t, LOCAL_ADDR, conf.ClientConfig.LocalAddr)
	assert.Equal(t, NOTIFY_API_ROOT, conf.ClientConfig.NotifyApiRoot)
}

func TestInvalidConfig(t *testing.T) {
	_, err := generateN7Mconfig(t, err_config)
	assert.Error(t, err)
}

func TestGetFromEnv(t *testing.T) {
	conf, err := generateN7Mconfig(t, empty_config)
	require.NoError(t, err)
	assert.Equal(t, false, conf.DisableN7)
	url1, _ := url.ParseRequestURI(n7.DefaultPcfApiRoot)
	assert.Equal(t, *url1, conf.ServerConfig.ApiRoot)
	assert.Equal(t, n7.DefaultPcfTokenUrl, conf.ServerConfig.TokenUrl)
	assert.Equal(t, n7.DefaultClientId, conf.ServerConfig.ClientId)
	assert.Equal(t, n7.DefaultClientSecret, conf.ServerConfig.ClientSecret)
	assert.Equal(t, n7.DefaultN7ClientAddr, conf.ClientConfig.LocalAddr)
	assert.Equal(t, n7.DefaultN7ClientApiRoot, conf.ClientConfig.NotifyApiRoot)
}

func generateN7Mconfig(t *testing.T, configString string) (*n7.N7Config, error) {
	err := mconfig.CreateLoadTempConfig(configString)
	assert.NoError(t, err)
	return n7.GetN7Config()
}
