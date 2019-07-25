/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"magma/feg/cloud/go/feg"
	cfgprotos "magma/feg/cloud/go/services/controller/protos"
	orc8rcfg "magma/orc8r/cloud/go/services/config"

	"github.com/golang/glog"
)

const (
	defaultCpuUtilThreshold      = 0.90
	defaultMemAvailableThreshold = 0.90
	defaultStaleUpdateThreshold  = 30
)

var defaultServices = []string{"S6A_PROXY", "SESSION_PROXY"}

func GetHealthConfigForNetwork(networkID string) *healthConfig {
	defaultConfig := &healthConfig{
		services:              defaultServices,
		cpuUtilThreshold:      defaultCpuUtilThreshold,
		memAvailableThreshold: defaultMemAvailableThreshold,
		staleUpdateThreshold:  defaultStaleUpdateThreshold,
	}
	config, err := orc8rcfg.GetConfig(networkID, feg.FegNetworkType, networkID)
	if err != nil {
		glog.V(2).Infof("Using default health configuration for network %s; %s", networkID, err)
		return defaultConfig
	}
	cloudFegConfig, ok := config.(*cfgprotos.Config)
	if !ok {
		glog.V(2).Infof("Using default health configuration for network %s; Invalid config format", networkID)
		return defaultConfig
	}
	healthParams := cloudFegConfig.GetHealth()
	if healthParams == nil {
		glog.V(2).Infof("Using default health configuration for network %s; Health config not found", networkID)
		return defaultConfig
	}
	if healthParams.GetCpuUtilizationThreshold() == 0 {
		glog.V(2).Infof("Using default health configuration for network %s; Cpu utilization threshold cannot be 0", networkID)
		return defaultConfig
	}
	if healthParams.GetMemoryAvailableThreshold() == 0 {
		glog.V(2).Infof("Using default health configuration for network %s; Memory available threshold cannot be 0", networkID)
		return defaultConfig
	}
	staleUpdateThreshold := healthParams.GetUpdateFailureThreshold() * healthParams.GetUpdateIntervalSecs()
	if staleUpdateThreshold == 0 {
		glog.V(2).Infof("Using default health configuration for network %s; Stale update threshold cannot be 0", networkID)
		return defaultConfig
	}
	return &healthConfig{
		services:              healthParams.GetHealthServices(),
		cpuUtilThreshold:      healthParams.GetCpuUtilizationThreshold(),
		memAvailableThreshold: healthParams.GetMemoryAvailableThreshold(),
		staleUpdateThreshold:  staleUpdateThreshold,
	}
}
