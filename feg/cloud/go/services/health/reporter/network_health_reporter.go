/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package reporter

import (
	"time"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health/metrics"
	"magma/feg/cloud/go/services/health/servicers"
	"magma/feg/cloud/go/services/health/storage"
	orc8rcfg "magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/magmad"

	"github.com/golang/glog"
)

type NetworkHealthStatusReporter struct {
	Store storage.HealthStorage
}

func NewNetworkHealthStatusReporter(store storage.HealthStorage) (*NetworkHealthStatusReporter, error) {
	return &NetworkHealthStatusReporter{Store: store}, nil
}

func (reporter *NetworkHealthStatusReporter) ReportHealthStatus(dur time.Duration) {
	for _ = range time.Tick(dur) {
		err := reporter.reportHealthStatus()
		if err != nil {
			glog.Errorf("err in reportHealthStatus: %v\n", err)
		}
	}
}

func (reporter *NetworkHealthStatusReporter) reportHealthStatus() error {
	networks, err := magmad.ListNetworks()
	if err != nil {
		return err
	}
	for _, nw := range networks {
		// Consider a FeG network to be only those that have FeG Network configs defined
		config, err := orc8rcfg.GetConfig(nw, feg.FegNetworkType, nw)
		if err != nil || config == nil {
			continue
		}
		gateways, err := magmad.ListGateways(nw)
		if err != nil {
			glog.Errorf("error getting gateways for network %v: %v\n", nw, err)
			continue
		}
		healthyGateways := 0
		for _, gw := range gateways {
			healthStatus, err := reporter.Store.GetHealth(nw, gw)
			if err != nil {
				glog.V(2).Infof("error getting health for network %s, gateway %s: %v\n", nw, gw, err)
				continue
			}
			status, _, err := servicers.AnalyzeHealthStats(healthStatus, nw)
			if err != nil {
				glog.V(2).Infof("error analyzing health stats for network %s, gateway %s: %v", nw, gw, err)
			}
			if status == protos.HealthStatus_HEALTHY {
				healthyGateways++
			}
		}
		metrics.TotalGatewayCount.WithLabelValues(nw).Set(float64(len(gateways)))
		metrics.HealthyGatewayCount.WithLabelValues(nw).Set(float64(healthyGateways))
	}
	return nil
}
