/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"encoding/hex"
	"net"

	"magma/cwf/gateway/registry"
	"magma/orc8r/cloud/go/service/config"

	"github.com/golang/glog"
)

const (
	defaultRadiusAddress = "192.168.70.101:1812"
	defaultRadiusSecret  = "123456"
	defaultCwagTestBr    = "cwag_test_br0"
	defaultBrMac         = "76-02-5B-80-EC-44"
)

var (
	defaultAmf = []byte("\x67\x41")
	defaultOp  = []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11")
)

func GetUESimConfig() (*UESimConfig, error) {
	uecfg, err := config.GetServiceConfig("", registry.UeSim)
	if err != nil {
		glog.Error(err)
		return &UESimConfig{
			op:            defaultOp,
			amf:           defaultAmf,
			radiusAddress: string(defaultRadiusAddress),
			radiusSecret:  string(defaultRadiusSecret),
			brMac:         defaultBrMac,
		}, nil
	}
	addr, err := uecfg.GetStringParam("radius_address")
	if err != nil {
		addr = defaultRadiusAddress
	}
	secret, err := uecfg.GetStringParam("radius_secret")
	if err != nil {
		secret = defaultRadiusSecret
	}
	brName, err := uecfg.GetStringParam("ue_bridge")
	if err != nil {
		brName = defaultCwagTestBr
	}
	brMac := getBridgeMac(brName)
	amfBytes := getHexParam(uecfg, "amf", defaultAmf)
	opBytes := getHexParam(uecfg, "op", defaultOp)
	glog.Infof("UE SIM Config - OP: %x, AMF: %x, RADIUS Endpoint: %s, RADIUS Secret: %s",
		opBytes, amfBytes, string(addr), string(secret))
	return &UESimConfig{
		op:            opBytes,
		amf:           amfBytes,
		radiusAddress: string(addr),
		radiusSecret:  string(secret),
		brMac:         brMac,
	}, nil
}

// TODO: Store UE MAC and add necessary OVS flows to allow traffic to use
// the stored UE MAC as eth src
func getBridgeMac(br string) string {
	brInterface, err := net.InterfaceByName(br)
	if err != nil {
		glog.Errorf("No bridge named %s exists. Using default: %s as bridge MAC", br, defaultBrMac)
		return defaultBrMac
	}
	return brInterface.HardwareAddr.String()
}

func getHexParam(cfg *config.ConfigMap, param string, defaultBytes []byte) []byte {
	param, err := cfg.GetStringParam(param)
	if err != nil {
		return defaultBytes
	}
	paramBytes, err := hex.DecodeString(param)
	if err != nil {
		return defaultBytes
	}
	return paramBytes
}
