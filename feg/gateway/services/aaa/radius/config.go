/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
package radius

import (
	"magma/feg/cloud/go/protos/mconfig"
)

const (
	defaultNetwork  = "udp"
	defaultAuthAddr = ":1812"
	defaultAcctAddr = ":1813"
	defaultSecret   = "123456"
)

var defaultConfigs = &mconfig.RadiusConfig{
	Secret:   []byte(defaultSecret),
	Network:  defaultNetwork,
	AuthAddr: defaultAuthAddr,
	AcctAddr: defaultAcctAddr,
}

func validateConfigs(cfg *mconfig.RadiusConfig) *mconfig.RadiusConfig {
	res := &mconfig.RadiusConfig{}
	if cfg != nil {
		*res = *cfg
	}
	if len(res.Secret) == 0 {
		res.Secret = []byte(defaultSecret)
	}
	if len(res.Network) == 0 {
		res.Network = defaultNetwork
	}
	if len(res.AuthAddr) == 0 {
		res.AuthAddr = defaultAuthAddr
	}
	if len(res.AcctAddr) == 0 {
		res.AcctAddr = defaultAcctAddr
	}
	return res
}
