/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"encoding/json"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/configurator"
)

type Conf1 struct {
	Value1 int
	Value2 string
	Value3 []byte
}

type Conf1Manager struct{}

func NewConfig1Manager() serde.Serde {
	return &Conf1Manager{}
}

func (manager *Conf1Manager) GetDomain() string {
	return config.SerdeDomain
}

func (manager *Conf1Manager) GetType() string {
	return "Config1"
}

func (manager *Conf1Manager) Serialize(config interface{}) ([]byte, error) {
	return json.Marshal(config)
}

func (manager *Conf1Manager) Deserialize(message []byte) (interface{}, error) {
	var out Conf1
	err := json.Unmarshal(message, &out)
	return &out, err
}

type Conf2 struct {
	Value1 []string
	Value2 int
}

type Conf2Manager struct{}

func NewConfig2Manager() serde.Serde {
	return &Conf2Manager{}
}

func (manager *Conf2Manager) GetDomain() string {
	return config.SerdeDomain
}

func (manager *Conf2Manager) GetType() string {
	return "Config2"
}

func (manager *Conf2Manager) Serialize(config interface{}) ([]byte, error) {
	return json.Marshal(config)
}

func (manager *Conf2Manager) Deserialize(message []byte) (interface{}, error) {
	var out Conf2
	err := json.Unmarshal(message, &out)
	return &out, err
}

type Conf1ConfiguratorManager struct{}

func NewConfig1ConfiguratorManager() serde.Serde {
	return &Conf1ConfiguratorManager{}
}

func (manager *Conf1ConfiguratorManager) GetDomain() string {
	return configurator.NetworkEntitySerdeDomain
}

func (manager *Conf1ConfiguratorManager) GetType() string {
	return "Config1"
}

func (manager *Conf1ConfiguratorManager) Serialize(config interface{}) ([]byte, error) {
	return json.Marshal(config)
}

func (manager *Conf1ConfiguratorManager) Deserialize(message []byte) (interface{}, error) {
	var out Conf1
	err := json.Unmarshal(message, &out)
	return &out, err
}

type Conf2ConfiguratorManager struct{}

func NewConfig2ConfiguratorManager() serde.Serde {
	return &Conf2ConfiguratorManager{}
}

func (manager *Conf2ConfiguratorManager) GetDomain() string {
	return configurator.NetworkEntitySerdeDomain
}

func (manager *Conf2ConfiguratorManager) GetType() string {
	return "Config2"
}

func (manager *Conf2ConfiguratorManager) Serialize(config interface{}) ([]byte, error) {
	return json.Marshal(config)
}

func (manager *Conf2ConfiguratorManager) Deserialize(message []byte) (interface{}, error) {
	var out Conf2
	err := json.Unmarshal(message, &out)
	return &out, err
}
