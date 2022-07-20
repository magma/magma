package servicers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/testcore/hss/servicers"
	gwMconfig "magma/gateway/mconfig"
)

func TestGetHssProxyConfig(t *testing.T) {
	conf := generateHssMconfig(t, test_config)
	assert.Equal(t, "magma-fedgw.magma.com", conf.Server.DestHost)
}

func TestValidateConfig(t *testing.T) {
	conf := generateHssMconfig(t, test_config)
	err := servicers.ValidateConfig(conf)
	assert.NoError(t, err)
}

func TestValidateHostWithPortString(t *testing.T) {
	err := servicers.ValidateHostWithPort("")
	assert.Error(t, err)
	err = servicers.ValidateHostWithPort("localhost:3768")
	assert.NoError(t, err)
	err = servicers.ValidateHostWithPort("127.0.0.1:3768")
	assert.NoError(t, err)
	err = servicers.ValidateHostWithPort("127.0.0.1")
	assert.NoError(t, err)
	err = servicers.ValidateHostWithPort("localhost")
	assert.NoError(t, err)
}

func generateHssMconfig(t *testing.T, configString string) *mconfig.HSSConfig {
	err := gwMconfig.CreateLoadTempConfig(configString)
	assert.NoError(t, err)
	config, err := servicers.GetHSSConfig()
	assert.NoError(t, err)
	return config
}

var (
	test_config = `{
	  "configsByKey": {
		"hss": {
		  "@type": "type.googleapis.com/magma.mconfig.HSSConfig",
		  "server": {
			"protocol": "tcp",
			"address": "localhost:3768",
			"local_address": "localhost:3767",
			"dest_host": "magma-fedgw.magma.com",
			"dest_realm": "magma.com"
		  },
		  "lte_auth_op": "",
		  "lte_auth_amf": "",
		  "default_sub_profile": {
			"max_ul_bit_rate": 10000000,
			"max_dl_bit_rate": 20000000
		  }
		}
	  }
	}`
)
