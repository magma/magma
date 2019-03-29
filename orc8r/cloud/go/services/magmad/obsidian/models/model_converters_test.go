/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models_test

import (
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/protos"
	checkind_models "magma/orc8r/cloud/go/services/checkind/obsidian/models"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
)

func TestGatewayStateToModel(t *testing.T) {
	state := &view_factory.GatewayState{
		GatewayID: "gw0",
		Config: map[string]interface{}{
			"Hello": "World!",
		},
		Status: &protos.GatewayStatus{
			Time: 12345,
			Checkin: &protos.CheckinRequest{
				GatewayId: "gw0",
				Status: &protos.ServiceStatus{
					Meta: map[string]string{
						"hello": "world",
					},
				},
				SystemStatus: &protos.SystemStatus{
					Time:       1495484735606,
					CpuUser:    31498,
					CpuSystem:  8361,
					CpuIdle:    1869111,
					MemTotal:   1016084,
					MemUsed:    54416,
					MemFree:    412772,
					UptimeSecs: 1234,
					SwapTotal:  1016081,
					SwapUsed:   54415,
					SwapFree:   412771,
					DiskPartitions: []*protos.DiskPartition{
						{
							Device:     "/dev/sda1",
							MountPoint: "/",
							Total:      1,
							Used:       2,
							Free:       3,
						},
					},
				},
				PlatformInfo: &protos.PlatformInfo{
					VpnIp: "facebook.com",
					Packages: []*protos.Package{
						{
							Name:    "magma",
							Version: "0.0.0.0",
						},
					},
					KernelVersion:           "42",
					KernelVersionsInstalled: []string{"42", "43"},
				},
				MachineInfo: &protos.MachineInfo{
					CpuInfo: &protos.CPUInfo{
						CoreCount:      4,
						ThreadsPerCore: 1,
						Architecture:   "x86_64",
						ModelName:      "Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz",
					},
					NetworkInfo: &protos.NetworkInfo{
						NetworkInterfaces: []*protos.NetworkInterface{
							{
								NetworkInterfaceId: "gtp_br0",
								Status:             protos.NetworkInterface_UP,
								MacAddress:         "08:00:27:1e:8a:32",
								IpAddresses:        []string{"10.10.10.1"},
								Ipv6Addresses:      []string{"fe80::a00:27ff:fe1e:8332"},
							},
						},
						RoutingTable: []*protos.Route{
							{
								DestinationIp:      "0.0.0.0",
								GatewayIp:          "10.10.10.1",
								Genmask:            "255.255.255.0",
								NetworkInterfaceId: "eth0",
							},
						},
					},
				},
			},
		},
		Record: &magmadprotos.AccessGatewayRecord{
			HwId: &protos.AccessGatewayID{Id: "gw0"},
			Name: "Gateway 0",
			Key: &protos.ChallengeKey{
				KeyType: protos.ChallengeKey_ECHO,
			},
		},
	}
	expectedModel := &magmad_models.GatewayStateType{
		Config: map[string]interface{}{
			"Hello": "World!",
		},
		GatewayID: "gw0",
		Status: &checkind_models.GatewayStatus{
			CheckinTime: 12345,
			Meta:        map[string]string{"hello": "world"},
			SystemStatus: &checkind_models.SystemStatus{
				Time:       1495484735606,
				CPUUser:    31498,
				CPUSystem:  8361,
				CPUIDLE:    1869111,
				MemTotal:   1016084,
				MemUsed:    54416,
				MemFree:    412772,
				UptimeSecs: 1234,
				SwapTotal:  1016081,
				SwapUsed:   54415,
				SwapFree:   412771,
				DiskPartitions: []*checkind_models.DiskPartition{
					{
						Device:     "/dev/sda1",
						MountPoint: "/",
						Total:      1,
						Used:       2,
						Free:       3,
					},
				},
			},
			PlatformInfo: &checkind_models.PlatformInfo{
				VpnIP: "facebook.com",
				Packages: []*checkind_models.Package{
					{
						Name:    "magma",
						Version: "0.0.0.0",
					},
				},
				KernelVersion:           "42",
				KernelVersionsInstalled: []string{"42", "43"},
			},
			MachineInfo: &checkind_models.MachineInfo{
				CPUInfo: &checkind_models.MachineInfoCPUInfo{
					CoreCount:      4,
					ThreadsPerCore: 1,
					Architecture:   "x86_64",
					ModelName:      "Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz",
				},
				NetworkInfo: &checkind_models.MachineInfoNetworkInfo{
					NetworkInterfaces: []*checkind_models.NetworkInterface{
						{
							NetworkInterfaceID: "gtp_br0",
							Status:             checkind_models.NetworkInterfaceStatusUP,
							MacAddress:         "08:00:27:1e:8a:32",
							IPAddresses:        []string{"10.10.10.1"},
							IPV6Addresses:      []string{"fe80::a00:27ff:fe1e:8332"},
						},
					},
					RoutingTable: []*checkind_models.Route{
						{
							DestinationIP:      "0.0.0.0",
							GatewayIP:          "10.10.10.1",
							Genmask:            "255.255.255.0",
							NetworkInterfaceID: "eth0",
						},
					},
				},
			},
			HardwareID:              "gw0",
			KernelVersion:           "42",
			KernelVersionsInstalled: []string{"42", "43"},
			Version:                 "0.0.0.0",
			VpnIP:                   "facebook.com",
		},
		Record: &magmad_models.AccessGatewayRecord{
			HwID: &magmad_models.HwGatewayID{ID: "gw0"},
			Key: &magmad_models.ChallengeKey{
				KeyType: magmad_models.ChallengeKeyKeyTypeECHO,
			},
			Name: "Gateway 0",
		},
	}
	actualModel, err := magmad_models.GatewayStateToModel(state)
	assert.NoError(t, err)
	assert.Equal(t, expectedModel, actualModel)
}

func TestGatewayStateToModelNilFields(t *testing.T) {
	state := &view_factory.GatewayState{
		GatewayID: "gw0",
		Config:    make(map[string]interface{}),
		Status:    nil,
		Record:    nil,
	}
	expectedModel := &magmad_models.GatewayStateType{
		Config:    make(map[string]interface{}),
		GatewayID: "gw0",
		Status:    nil,
		Record:    nil,
	}
	actualModel, err := magmad_models.GatewayStateToModel(state)
	assert.NoError(t, err)
	assert.Equal(t, expectedModel, actualModel)

	state.Status = &protos.GatewayStatus{
		Time:    12345,
		Checkin: nil,
	}
	expectedModel.Status = &checkind_models.GatewayStatus{
		CheckinTime:  12345,
		HardwareID:   "",
		Meta:         map[string]string{},
		SystemStatus: nil,
		Version:      "",
	}
	actualModel, err = magmad_models.GatewayStateToModel(state)
	assert.NoError(t, err)
	assert.Equal(t, expectedModel, actualModel)

	state.Record = &magmadprotos.AccessGatewayRecord{
		HwId: nil,
		Name: "gw0",
		Key:  nil,
	}
	expectedModel.Record = &magmad_models.AccessGatewayRecord{
		HwID: nil,
		Key:  nil,
		Name: "gw0",
	}
	actualModel, err = magmad_models.GatewayStateToModel(state)
	assert.NoError(t, err)
	assert.Equal(t, expectedModel, actualModel)
}

func TestJSONMapToProtobufStruct(t *testing.T) {
	jsonMap := map[string]interface{}{
		"nil":    nil,
		"number": 1.0,
		"string": "str",
		"struct": map[string]interface{}{
			"a": 2.0,
		},
		"list": []interface{}{1.0, "foo"},
	}
	marshaled, err := json.Marshal(jsonMap)
	assert.NoError(t, err)
	expectedProtobufStruct := &structpb.Struct{}
	err = jsonpb.UnmarshalString(string(marshaled), expectedProtobufStruct)
	assert.NoError(t, err)

	actualProtobufStruct, err := magmad_models.JSONMapToProtobufStruct(jsonMap)

	assert.NoError(t, err)
	assert.Equal(t, expectedProtobufStruct, actualProtobufStruct)
}

func TestProtobufStructToJSONMap(t *testing.T) {
	expectedJsonMap := map[string]interface{}{
		"nil":    nil,
		"number": 1.0,
		"string": "str",
		"struct": map[string]interface{}{
			"a": 2.0,
		},
		"list": []interface{}{1.0, "foo"},
	}
	marshaled, err := json.Marshal(expectedJsonMap)
	assert.NoError(t, err)
	protobufStruct := &structpb.Struct{}
	err = jsonpb.UnmarshalString(string(marshaled), protobufStruct)
	assert.NoError(t, err)

	actualJsonMap, err := magmad_models.ProtobufStructToJSONMap(protobufStruct)

	assert.NoError(t, err)
	assert.Equal(t, expectedJsonMap, actualJsonMap)
}
