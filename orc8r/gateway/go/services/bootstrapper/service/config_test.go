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
package service_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/gateway/config"
	platform_config "magma/orc8r/lib/go/service/config"
)

const testConfigYaml = `
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

func TestConfig(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "test_control_proxy*.yml")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(testConfigYaml)); err != nil {
		t.Fatal(err)
	}
	dir := path.Dir(tmpfile.Name())
	serviceName := strings.TrimSuffix(path.Base(tmpfile.Name()), ".yml")
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	cfg := config.NewDefaultControlProxyCfg()
	file1, file2, err := platform_config.GetStructuredServiceConfigExt(
		"", serviceName, dir, "", "", cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, tmpfile.Name(), file1)
	assert.Equal(t, "", file2)
	assert.Equal(t, "/var/opt/magma/secrets/certs/gateway.crt", cfg.GwCertFile)
	assert.Equal(t, "/var/opt/magma/certs/rootCA.pem", cfg.RootCaFile)
	assert.Equal(t, "controller.magma.foobar.com", cfg.CloudAddr)
	assert.Equal(t, 9999, cfg.CloudPort)
	assert.Equal(t, false, cfg.ProxyCloudConnection)
}
