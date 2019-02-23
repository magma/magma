/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/magmad"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/stretchr/testify/assert"
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

func GetMockStatus(t *testing.T, gwID string) (*protos.GatewayStatus, []byte) {
	ret := &protos.GatewayStatus{
		Time: 10,
		Checkin: &protos.CheckinRequest{
			GatewayId:       gwID,
			MagmaPkgVersion: "v1",
			Status: &protos.ServiceStatus{
				Meta: map[string]string{
					"a": "b",
					"c": "d",
				},
			},
			SystemStatus: &protos.SystemStatus{
				Time:    100,
				CpuUser: 300,
				MemFree: 1024,
			},
			KernelVersionsInstalled: []string{},
		},
	}
	retBytes, err := protos.MarshalIntern(ret)
	assert.NoError(t, err)
	return ret, retBytes
}

func GetDefaultStatus(t *testing.T) *protos.GatewayStatus {
	ret, _ := GetMockStatus(t, "gw1")
	return ret
}

func GetMockRecord(t *testing.T, hwID string) (*magmadprotos.AccessGatewayRecord, []byte) {
	ret := &magmadprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{
			Id: hwID,
		},
		Name: "Gateway 1",
		Key: &protos.ChallengeKey{
			KeyType: protos.ChallengeKey_ECHO,
			Key:     []byte{1, 2, 3, 4},
		},
		Ip:   "127.0.0.1",
		Port: 8000,
	}
	retBytes, err := protos.MarshalIntern(ret)
	assert.NoError(t, err)
	return ret, retBytes
}

func GetDefaultRecord(t *testing.T) *magmadprotos.AccessGatewayRecord {
	ret, _ := GetMockRecord(t, "hwid1")
	return ret
}
