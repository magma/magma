/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service implements the core of bootstrapper
package service

import (
	"log"

	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/service/config"
)

const (
	BootstrapCfgKey          = "bootstrap_config"
	BootstrapCfgChallengeKey = "challenge_key"
)

// ControlProxyCfg represents control_proxy.yml configuration struct
type ControlProxyCfg struct {
	NghttpxConfigLocation string `yaml:"nghttpx_config_location"`

	// GW Certificate params
	RootCaFile    string `yaml:"rootca_cert"`
	GwCertFile    string `yaml:"gateway_cert"`
	GwCertKeyFile string `yaml:"gateway_key"`

	LocalPort            int    `yaml:"local_port"`
	CloudAddr            string `yaml:"cloud_address"`
	CloudPort            int    `yaml:"cloud_port"`
	BootstrapAddr        string `yaml:"bootstrap_address"`
	BootstrapPort        int    `yaml:"bootstrap_port"`
	ProxyCloudConnection bool   `yaml:"proxy_cloud_connections"`
}

// NewDefaultBootsrapper returns new Bootstrapper struct with default configuration
func NewDefaultBootsrapper() *Bootstrapper {
	return &Bootstrapper{
		ChallengeKeyFile: "/var/opt/magma/certs/gw_challenge.key",
		CpConfig: ControlProxyCfg{
			NghttpxConfigLocation: "/var/tmp/nghttpx.conf",
			RootCaFile:            "/var/opt/magma/certs/rootCA.pem",
			GwCertFile:            "/var/opt/magma/certs/gateway.crt",
			GwCertKeyFile:         "/var/opt/magma/certs/gateway.key",
			LocalPort:             8443,
			CloudAddr:             "",
			CloudPort:             443,
			BootstrapAddr:         "",
			BootstrapPort:         443,
			ProxyCloudConnection:  false,
		},
	}
}

func (b *Bootstrapper) updateBootstrapperKeyCfg() *Bootstrapper {
	cfg, err := config.GetServiceConfig("", definitions.MagmadServiceName)
	if err != nil {
		log.Printf("Error Getting Bootstrapper Key Configs: %v", err)
		return b
	}
	bootstrCfg, err := cfg.GetMapParam(BootstrapCfgKey)
	if err == nil && bootstrCfg != nil {
		if param, ok := bootstrCfg[BootstrapCfgChallengeKey]; ok {
			if challengeKeyFilePath, ok := param.(string); ok && len(challengeKeyFilePath) > 0 {
				b.ChallengeKeyFile = challengeKeyFilePath
			} else {
				log.Printf("Bootstrapper Challenge File Path %v (%T) is not a string", param, param)
			}
		} else {
			log.Printf("Could not find Bootstrapper Challenge File Key %s", BootstrapCfgChallengeKey)
		}
	} else {
		log.Printf("Could not find Bootstrapper Config Key %s", BootstrapCfgKey)
	}
	return b
}

func (b *Bootstrapper) updateFromControlProxyCfg() *Bootstrapper {
	newCfg := b.CpConfig // copy current configs
	err := config.GetStructuredServiceConfig("", definitions.ControlProxyServiceName, &newCfg)
	if err != nil {
		log.Printf("Error Getting Bootstrapper Control Proxy Configs: %v,\n\tcontinue using old configs: %+v",
			err, b.CpConfig)
	} else {
		// success, copy over the new configs
		b.CpConfig = newCfg
	}
	return b
}
