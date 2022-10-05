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

// package service_test
package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/gateway/config"
	platform_cfg "magma/orc8r/lib/go/service/config"
)

const testCPConfigYaml = `
#
# Copyright (c) ...
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# nghttpx config will be generated here and used
nghttpx_config_location: /var/tmp/nghttpx.conf

# Location for certs
rootca_cert: /var/opt/magma/certs/rootCA.pem
gateway_cert: /var/opt/magma/secrets/certs/gateway.crt
gateway_key: /var/opt/magma/secrets/certs/gateway.key

# Listening port of the proxy for local services. The port would be closed
# for the rest of the world.
local_port: 8888

# Cloud address for reaching out to the cloud.
cloud_address: controller.magma.foobar.com
cloud_port: 9999

bootstrap_address: bootstrapper-controller.magma.foobar.com
bootstrap_port: 4444

# Option to use nghttpx for proxying. If disabled, the individual
# services would establish the TLS connections themselves.
proxy_cloud_connections: False

# Allows http_proxy usage if the environment variable is present
allow_http_proxy: True
`

func TestControlProxyConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "magma_cfg_test")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	cfgFilePath := filepath.Join(dir, "control_proxy.yml")
	err = ioutil.WriteFile(cfgFilePath, []byte(testCPConfigYaml), os.ModePerm)
	assert.NoError(t, err)

	platform_cfg.SetConfigDirectories(dir, dir+"/foo", dir+"/bar")
	cfg := config.GetControlProxyConfigs()
	assert.Equal(t, "/var/opt/magma/secrets/certs/gateway.crt", cfg.GwCertFile)
	assert.Equal(t, "/var/opt/magma/certs/rootCA.pem", cfg.RootCaFile)
	assert.Equal(t, "controller.magma.foobar.com", cfg.CloudAddr)
	assert.Equal(t, 9999, cfg.CloudPort)
	assert.Equal(t, false, cfg.ProxyCloudConnection)
}
