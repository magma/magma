/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadRadiusConfig(t *testing.T) {
	conf, err := Read("../radius.config.json")
	require.Nil(t, err)
	require.NotNil(t, conf)
	require.NotNil(t, conf.Server)
	require.Equal(t, conf.Server.LoadBalance, LoadBalanceConfig{})
}

func TestLoadLBConfig(t *testing.T) {
	conf, err := Read("./samples/lb.config.json")
	require.Nil(t, err)
	require.NotNil(t, conf)
	require.NotNil(t, conf.Server)
	require.NotEmpty(t, conf.Server.LoadBalance.ServiceTiers)
	require.NotEmpty(t, conf.Server.LoadBalance.LiveTier)
	require.NotEmpty(t, conf.Server.LoadBalance.Canaries)
}
