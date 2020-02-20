/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service_test
package service_test

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"magma/gateway/services/bootstrapper/service"
	"magma/orc8r/lib/go/service/config"

	"github.com/stretchr/testify/assert"
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
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(testConfigYaml)); err != nil {
		t.Fatal(err)
	}
	dir := path.Dir(tmpfile.Name())
	serviceName := strings.TrimSuffix(path.Base(tmpfile.Name()), ".yml")
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	b := service.NewDefaultBootsrapper()
	if err := config.GetStructuredServiceConfigExt(
		"", serviceName, dir, "", "", &(b.CpConfig)); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "/var/opt/magma/secrets/certs/gateway.crt", b.CpConfig.GwCertFile)
	assert.Equal(t, "/var/opt/magma/certs/rootCA.pem", b.CpConfig.RootCaFile)
	assert.Equal(t, "controller.magma.foobar.com", b.CpConfig.CloudAddr)
	assert.Equal(t, 9999, b.CpConfig.CloudPort)
	assert.Equal(t, false, b.CpConfig.ProxyCloudConnection)
}
