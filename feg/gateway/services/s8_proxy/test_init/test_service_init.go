package test_init

import (
	"context"
	"fmt"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/s8_proxy/servicers"
	"magma/feg/gateway/services/s8_proxy/servicers/mock_pgw"
	"magma/gateway/mconfig"
	"magma/orc8r/cloud/go/test_utils"
)

// StartS8AndPGWService start both S8 proxy service and PGW (GTP service) for testing
func StartS8AndPGWService(t *testing.T, clientAddr, serverAddr string) (*mock_pgw.MockPgw, error) {
	// Start pgw and get the server address and real port
	mockPgw, err := mock_pgw.NewStarted(context.Background(), serverAddr)
	if err != nil {
		return nil, err
	}
	// overwrite server Addrs to make sure we have the right port
	serverAddr = mockPgw.LocalAddr().String()

	// create config string with its proper values
	fegConfigFmt := `{
		"configsByKey": {
			"s8_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.S8Config",
				"logLevel": "INFO",
				"local_address": "%s",
				"pgw_address": "%s"
			}
		}
	}`
	configStr := fmt.Sprintf(fegConfigFmt, clientAddr, serverAddr)

	// load mconfig
	err = mconfig.CreateLoadTempConfig(configStr)
	if err != nil {
		return nil, err
	}
	config := servicers.GetS8ProxyConfig()

	// create and launch s8 Proxy
	s8service, err := servicers.NewS8Proxy(config)
	if err != nil {
		return nil, err
	}
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.S8_PROXY)
	protos.RegisterS8ProxyServer(srv.GrpcServer, s8service)
	go srv.RunTest(lis)
	return mockPgw, nil
}
