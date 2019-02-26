/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"encoding/json"

	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/magmad"
)

type Conf1 struct {
	Value1 int
	Value2 string
	Value3 []byte
}

type Conf1Manager struct{}

func NewConfig1Manager() registry.ConfigManager {
	return &Conf1Manager{}
}

func (manager *Conf1Manager) GetConfigType() string {
	return "Config1"
}

func (manager *Conf1Manager) GetGatewayIdsForConfig(networkId string, configKey string) ([]string, error) {
	return magmad.ListGateways(networkId)
}

func (manager *Conf1Manager) MarshalConfig(config interface{}) ([]byte, error) {
	return json.Marshal(config)
}

func (manager *Conf1Manager) UnmarshalConfig(message []byte) (interface{}, error) {
	var out Conf1
	err := json.Unmarshal(message, &out)
	return &out, err
}

type Conf2 struct {
	Value1 []string
	Value2 int
}

type Conf2Manager struct{}

func NewConfig2Manager() registry.ConfigManager {
	return &Conf2Manager{}
}

func (manager *Conf2Manager) GetConfigType() string {
	return "Config2"
}

func (manager *Conf2Manager) GetGatewayIdsForConfig(networkId string, configKey string) ([]string, error) {
	return []string{configKey}, nil
}

func (manager *Conf2Manager) MarshalConfig(config interface{}) ([]byte, error) {
	return json.Marshal(config)
}

func (manager *Conf2Manager) UnmarshalConfig(message []byte) (interface{}, error) {
	var out Conf2
	err := json.Unmarshal(message, &out)
	return &out, err
}
