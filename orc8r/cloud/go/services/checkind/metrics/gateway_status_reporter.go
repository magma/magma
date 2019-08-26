/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metrics

import (
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/checkind/store"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"

	"github.com/golang/glog"
)

type GatewayStatusReporter struct {
	Store *store.CheckinStore
}

func NewGatewayStatusReporter(store *store.CheckinStore) (*GatewayStatusReporter, error) {
	return &GatewayStatusReporter{Store: store}, nil
}

func (reporter *GatewayStatusReporter) ReportCheckinStatus(dur time.Duration) {
	for _ = range time.Tick(dur) {
		err := reporter.reportCheckinStatus()
		if err != nil {
			glog.Errorf("err in reportCheckinStatus: %v\n", err)
		}
	}
}

func (reporter *GatewayStatusReporter) reportCheckinStatus() error {
	networks, err := configurator.ListNetworkIDs()
	if err != nil {
		return err
	}
	for _, networkID := range networks {
		magmadGatewayTypeVal := orc8r.MagmadGatewayType
		gateways, _, err := configurator.LoadEntities(networkID, &magmadGatewayTypeVal, nil, nil, nil, configurator.EntityLoadCriteria{})
		if err != nil {
			glog.Errorf("error getting gateways for network %v: %v\n", networkID, err)
			continue
		}
		numUpGateways := 0
		for _, gatewayEntity := range gateways {
			gatewayID := gatewayEntity.Key
			status, err := state.GetGatewayStatus(networkID, gatewayEntity.PhysicalID)
			if err != nil {
				glog.Errorf("Failed to get state for nwID:%s gwID:%s deviceID:%s", networkID, gatewayID, gatewayEntity.PhysicalID)
				continue
			}
			// last check in more than 5 minutes ago
			if (time.Now().Unix() - int64(status.CheckinTime)/1000) > 60*5 {
				gwCheckinStatus.WithLabelValues(networkID, gatewayID).Set(0)
			} else {
				gwCheckinStatus.WithLabelValues(networkID, gatewayID).Set(1)
				numUpGateways += 1
			}

			// report mconfig age
			mconfigCreatedAt := status.PlatformInfo.ConfigInfo.MconfigCreatedAt
			if mconfigCreatedAt != 0 {
				gwMconfigAge.WithLabelValues(networkID, gatewayID).Set(float64(status.CheckinTime/1000 - mconfigCreatedAt))
			}
		}
		upGwCount.WithLabelValues(networkID).Set(float64(numUpGateways))
		totalGwCount.WithLabelValues(networkID).Set(float64(len(gateways)))
	}
	return nil
}
