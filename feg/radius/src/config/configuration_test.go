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
	conf, err := Read("../lb.config.json")
	require.Nil(t, err)
	require.NotNil(t, conf)
	require.NotNil(t, conf.Server)
	require.NotEmpty(t, conf.Server.LoadBalance.ServiceTiers)
	require.NotEmpty(t, conf.Server.LoadBalance.LiveTier)
	require.NotEmpty(t, conf.Server.LoadBalance.Canaries)
}
