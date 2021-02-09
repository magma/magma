package servicers_test

import (
	"magma/feg/gateway/services/s8_proxy/servicers"
	"magma/gateway/mconfig"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetS8ProxyConfig(t *testing.T) {
	conf := generateS8Mconfig(t, config)
	assert.Equal(t, ":1", conf.ClientAddr)
	assert.Equal(t, "10.0.0.1:0", conf.ServerAddr)
}

func generateS8Mconfig(t *testing.T, configString string) *servicers.S8ProxyConfig {
	err := mconfig.CreateLoadTempConfig(configString)
	assert.NoError(t, err)
	return servicers.GetS8ProxyConfig()
}

var (
	config = `{
		"configsByKey": {
			"s8_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.S8Config",
				"logLevel": "INFO",
				"local_address": ":1",
				"pgw_address": "10.0.0.1:0"
			}
		}
	}`
)
