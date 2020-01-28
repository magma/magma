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
	"github.com/pkg/errors"
)

const (
	defaultRadiusAuthAddress = "192.168.70.101:1812"
	defaultRadiusAcctAddress = "192.168.70.101:1813"
	defaultRadiusSecret      = "123456"
	defaultCwagTestBr        = "cwag_test_br0"
	defaultBrMac             = "76-02-5B-80-EC-44"
)

var (
	defaultAmf = []byte("\x67\x41")
	defaultOp  = []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11")
)

func getUESimConfig() (*UESimConfig, error) {
	uecfg, err := config.GetServiceConfig("", registry.UeSim)
	if err != nil {
		glog.Error(errors.Wrap(err, "No service config found, using default config"))
		return getDefaultUESimConfig(), nil
	}
	authAddr, err := uecfg.GetStringParam("radius_auth_address")
	if err != nil {
		authAddr = defaultRadiusAuthAddress
	}
	acctAddr, err := uecfg.GetStringParam("radius_acct_address")
	if err != nil {
		acctAddr = defaultRadiusAcctAddress
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
		opBytes, amfBytes, authAddr, secret)
	return &UESimConfig{
		op:                opBytes,
		amf:               amfBytes,
		radiusAuthAddress: authAddr,
		radiusAcctAddress: acctAddr,
		radiusSecret:      secret,
		brMac:             brMac,
	}, nil
}

func getDefaultUESimConfig() *UESimConfig {
	return &UESimConfig{
		op:                defaultOp,
		amf:               defaultAmf,
		radiusAuthAddress: defaultRadiusAuthAddress,
		radiusAcctAddress: defaultRadiusAcctAddress,
		radiusSecret:      defaultRadiusSecret,
		brMac:             defaultBrMac,
	}
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
