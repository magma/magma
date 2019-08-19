/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package models includes definitions of swagger generated REST API model Go structures
package models

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/strfmt"

	"magma/feg/cloud/go/protos/mconfig"
	fegprotos "magma/feg/cloud/go/services/controller/protos"
	"magma/orc8r/cloud/go/protos"
)

// ConvertibleConfig for user-facing swagger models
var formatsRegistry = strfmt.NewFormats()

func (m *NetworkFederationConfigs) ValidateModel() error {
	return m.Validate(formatsRegistry)
}

func (m *NetworkFederationConfigs) ToServiceModel() (interface{}, error) {
	magmadConfig := &fegprotos.Config{
		S6A: &fegprotos.S6AConfig{Server: &fegprotos.DiamClientConfig{}},
		Hss: &fegprotos.HSSConfig{
			Server:            &fegprotos.DiamServerConfig{},
			LteAuthOp:         []byte{},
			LteAuthAmf:        []byte{},
			DefaultSubProfile: &fegprotos.HSSConfig_SubscriptionProfile{},
			SubProfiles:       make(map[string]*fegprotos.HSSConfig_SubscriptionProfile),
		},
		Gx:               &fegprotos.GxConfig{Server: &fegprotos.DiamClientConfig{}},
		Gy:               &fegprotos.GyConfig{Server: &fegprotos.DiamClientConfig{}},
		Swx:              &fegprotos.SwxConfig{Server: &fegprotos.DiamClientConfig{}},
		EapAka:           &fegprotos.EapAkaConfig{},
		AaaServer:        &fegprotos.AAAConfig{},
		ServedNetworkIds: []string{},
		Health:           &fegprotos.HealthConfig{},
	}
	protos.FillIn(m, magmadConfig)
	protos.FillIn(m.S6a, magmadConfig.S6A)
	protos.FillIn(m.Hss, magmadConfig.Hss)
	protos.FillIn(m.Gx, magmadConfig.Gx)
	protos.FillIn(m.Gy, magmadConfig.Gy)
	protos.FillIn(m.Swx, magmadConfig.Swx)
	protos.FillIn(m.Health, magmadConfig.Health)
	protos.FillIn(m.EapAka, magmadConfig.EapAka)
	protos.FillIn(m.AaaServer, magmadConfig.AaaServer)
	if err := fegprotos.ValidateNetworkConfig(magmadConfig); err != nil {
		return nil, err
	}
	return magmadConfig, nil
}

func (m *NetworkFederationConfigs) FromServiceModel(magmadModel interface{}) error {
	magmadConfig, ok := magmadModel.(*fegprotos.Config)
	if !ok {
		return fmt.Errorf(
			"Invalid magmad config type to convert to. Expected *Config but got %s",
			reflect.TypeOf(magmadModel),
		)
	}
	protos.FillIn(magmadConfig, m)
	if m.S6a == nil {
		m.S6a = &NetworkFederationConfigsS6a{Server: &DiameterClientConfigs{}}
	} else if m.S6a.Server == nil {
		m.S6a.Server = &DiameterClientConfigs{}
	}
	if m.Hss == nil {
		m.Hss = &NetworkFederationConfigsHss{
			DefaultSubProfile: &SubscriptionProfile{},
			Server:            &DiameterServerConfigs{},
			SubProfiles:       make(map[string]SubscriptionProfile),
		}
	} else {
		if m.Hss.DefaultSubProfile == nil {
			m.Hss.DefaultSubProfile = &SubscriptionProfile{}
		}
		if m.Hss.Server == nil {
			m.Hss.Server = &DiameterServerConfigs{}
		}
		if m.Hss.SubProfiles == nil {
			m.Hss.SubProfiles = make(map[string]SubscriptionProfile)
		}
	}
	if m.Gx == nil {
		m.Gx = &NetworkFederationConfigsGx{Server: &DiameterClientConfigs{}}
	} else if m.Gx.Server == nil {
		m.Gx.Server = &DiameterClientConfigs{}
	}
	if m.Gy == nil {
		m.Gy = &NetworkFederationConfigsGy{Server: &DiameterClientConfigs{}}
	} else if m.Gy.Server == nil {
		m.Gy.Server = &DiameterClientConfigs{}
	}
	if m.Swx == nil {
		m.Swx = &NetworkFederationConfigsSwx{Server: &DiameterClientConfigs{}}
	} else if m.Swx.Server == nil {
		m.Swx.Server = &DiameterClientConfigs{}
	}
	if m.Health == nil {
		m.Health = &NetworkFederationConfigsHealth{}
	}
	if m.EapAka == nil {
		m.EapAka = &NetworkFederationConfigsEapAka{}
	}
	if m.AaaServer == nil {
		m.AaaServer = &NetworkFederationConfigsAaaServer{}
	}
	protos.FillIn(magmadConfig.S6A, m.S6a)
	protos.FillIn(magmadConfig.Hss, m.Hss)
	protos.FillIn(magmadConfig.Gx, m.Gx)
	protos.FillIn(magmadConfig.Gy, m.Gy)
	protos.FillIn(magmadConfig.Swx, m.Swx)
	protos.FillIn(magmadConfig.Health, m.Health)
	protos.FillIn(magmadConfig.EapAka, m.EapAka)
	protos.FillIn(magmadConfig.AaaServer, m.AaaServer)
	if m.ServedNetworkIds == nil {
		m.ServedNetworkIds = []string{}
	}
	return nil
}

func (m *GatewayFegConfigs) ValidateModel() error {
	return m.Validate(formatsRegistry)
}

func (m *GatewayFegConfigs) ToServiceModel() (interface{}, error) {
	magmadConfig := &fegprotos.Config{
		S6A: &fegprotos.S6AConfig{Server: &fegprotos.DiamClientConfig{}},
		Hss: &fegprotos.HSSConfig{
			Server:            &fegprotos.DiamServerConfig{},
			LteAuthOp:         []byte{},
			LteAuthAmf:        []byte{},
			DefaultSubProfile: &fegprotos.HSSConfig_SubscriptionProfile{},
			SubProfiles:       make(map[string]*fegprotos.HSSConfig_SubscriptionProfile),
		},
		Gx:               &fegprotos.GxConfig{Server: &fegprotos.DiamClientConfig{}},
		Gy:               &fegprotos.GyConfig{Server: &fegprotos.DiamClientConfig{}},
		Swx:              &fegprotos.SwxConfig{Server: &fegprotos.DiamClientConfig{}},
		ServedNetworkIds: []string{},
		Health:           &fegprotos.HealthConfig{},
		EapAka:           &fegprotos.EapAkaConfig{},
		AaaServer:        &fegprotos.AAAConfig{},
	}

	protos.FillIn(m, magmadConfig)
	protos.FillIn(m.S6a, magmadConfig.S6A)
	protos.FillIn(m.Hss, magmadConfig.Hss)
	protos.FillIn(m.Gx, magmadConfig.Gx)
	protos.FillIn(m.Gy, magmadConfig.Gy)
	protos.FillIn(m.Swx, magmadConfig.Swx)
	protos.FillIn(m.Health, magmadConfig.Health)
	protos.FillIn(m.EapAka, magmadConfig.EapAka)
	protos.FillIn(m.AaaServer, magmadConfig.AaaServer)
	if err := fegprotos.ValidateGatewayConfig(magmadConfig); err != nil {
		return nil, err
	}
	return magmadConfig, nil
}

func (m *GatewayFegConfigs) FromServiceModel(magmadModel interface{}) error {
	magmadConfig, ok := magmadModel.(*fegprotos.Config)
	if !ok {
		return fmt.Errorf(
			"Invalid magmad config type to convert to. Expected *Config but got %s",
			reflect.TypeOf(magmadModel),
		)
	}
	protos.FillIn(magmadModel, m)
	if m.S6a == nil {
		m.S6a = &NetworkFederationConfigsS6a{&DiameterClientConfigs{}}
	} else if m.S6a.Server == nil {
		m.S6a.Server = &DiameterClientConfigs{}
	}
	if m.Hss == nil {
		m.Hss = &NetworkFederationConfigsHss{
			DefaultSubProfile: &SubscriptionProfile{},
			Server:            &DiameterServerConfigs{},
			SubProfiles:       make(map[string]SubscriptionProfile),
		}
	} else {
		if m.Hss.DefaultSubProfile == nil {
			m.Hss.DefaultSubProfile = &SubscriptionProfile{}
		}
		if m.Hss.Server == nil {
			m.Hss.Server = &DiameterServerConfigs{}
		}
		if m.Hss.SubProfiles == nil {
			m.Hss.SubProfiles = make(map[string]SubscriptionProfile)
		}
	}
	if m.Gx == nil {
		m.Gx = &NetworkFederationConfigsGx{Server: &DiameterClientConfigs{}}
	} else if m.Gx.Server == nil {
		m.Gx.Server = &DiameterClientConfigs{}
	}
	if m.Gy == nil {
		m.Gy = &NetworkFederationConfigsGy{Server: &DiameterClientConfigs{}}
	} else if m.Gy.Server == nil {
		m.Gy.Server = &DiameterClientConfigs{}
	}
	if m.Swx == nil {
		m.Swx = &NetworkFederationConfigsSwx{Server: &DiameterClientConfigs{}}
	} else if m.Swx == nil {
		m.Swx.Server = &DiameterClientConfigs{}
	}
	if m.Health == nil {
		m.Health = &NetworkFederationConfigsHealth{}
	}
	if m.EapAka == nil {
		m.EapAka = &NetworkFederationConfigsEapAka{}
	}
	if m.AaaServer == nil {
		m.AaaServer = &NetworkFederationConfigsAaaServer{}
	}
	protos.FillIn(magmadConfig.S6A, m.S6a)
	protos.FillIn(magmadConfig.Hss, m.Hss)
	protos.FillIn(magmadConfig.Gx, m.Gx)
	protos.FillIn(magmadConfig.Gy, m.Gy)
	protos.FillIn(magmadConfig.Swx, m.Swx)
	protos.FillIn(magmadConfig.Health, m.Health)
	protos.FillIn(magmadConfig.EapAka, m.EapAka)
	protos.FillIn(magmadConfig.AaaServer, m.AaaServer)
	if m.ServedNetworkIds == nil {
		m.ServedNetworkIds = []string{}
	}
	return nil
}

func (config *DiameterClientConfigs) ToMconfig() *mconfig.DiamClientConfig {
	res := &mconfig.DiamClientConfig{}
	protos.FillIn(config, res)
	return res
}

func (config *DiameterServerConfigs) ToMconfig() *mconfig.DiamServerConfig {
	res := &mconfig.DiamServerConfig{}
	protos.FillIn(config, res)
	return res
}

func (profile *SubscriptionProfile) ToMconfig() *mconfig.HSSConfig_SubscriptionProfile {
	res := &mconfig.HSSConfig_SubscriptionProfile{}
	protos.FillIn(profile, res)
	return res
}
