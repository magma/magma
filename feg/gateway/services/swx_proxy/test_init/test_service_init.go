/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"fmt"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/mconfig"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/swx_proxy/servicers/test"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) error {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SWX_PROXY)

	// Create tmp mconfig test file & load configs from it
	fegConfigFmt := `{
		"configsByKey": {
			"swx_proxy": {
				"@type": "type.googleapis.com/magma.mconfig.SwxConfig",
				"logLevel": "INFO",
				"server": {
					"protocol": "sctp",
					"address": "%s",
					"retransmits": 3,
					"watchdogInterval": 1,
					"retryCount": 5,
					"productName": "magma_test",
					"realm": "openair4G.eur",
					"host": "magma-oai.openair4G.eur"
				}
			}
		}
	}`

	err := mconfig.CreateLoadTempConfig(fmt.Sprintf(fegConfigFmt, "127.0.0.1:0"))
	if err != nil {
		return err
	}
	clientCfg, serverCfg := servicers.GetSwxProxyConfigs()
	serverAddr, err := test.StartTestSwxServer(serverCfg.Protocol, serverCfg.Addr)
	if err != nil {
		return err
	}
	// Update server config with chosen port of swx test server
	serverCfg.Addr = serverAddr
	service, err := servicers.NewSwxProxy(clientCfg, serverCfg)
	if err != nil {
		return err
	}
	protos.RegisterSwxProxyServer(srv.GrpcServer, service)
	go srv.RunTest(lis)
	return nil
}
