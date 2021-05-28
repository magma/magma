package servicers_test

import (
	"testing"

	"magma/feg/gateway/services/s8_proxy/servicers"
	"magma/gateway/mconfig"

	"github.com/stretchr/testify/assert"
)

func TestGetS8ProxyConfig(t *testing.T) {
	conf := generateS8Mconfig(t, config)
	assert.Equal(t, ":2222", conf.ClientAddr)
	assert.Equal(t, "10.0.0.1:9999", conf.ServerAddr.String())
}

func TestGetS8ProxyConfig_noPgw(t *testing.T) {
	conf := generateS8Mconfig(t, config_noPGW)
	assert.Equal(t, ":3333", conf.ClientAddr)
	assert.Nil(t, conf.ServerAddr)
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
				"localAddress": ":2222",
				"pgwAddress": "10.0.0.1:9999"
			}
		}
	}`

	config_noPGW = `{
		"configsByKey": {
			"s8_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.S8Config",
				"logLevel": "INFO",
				"localAddress": ":3333",
				"pgwAddress": ""
			}
		}
	}`
)
