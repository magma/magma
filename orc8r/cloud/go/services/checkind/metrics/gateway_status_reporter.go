/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metrics

import (
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/checkind/store"
	"magma/orc8r/cloud/go/services/magmad"

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
	networks, err := magmad.ListNetworks()
	if err != nil {
		return err
	}
	for _, nw := range networks {
		gateways, err := magmad.ListGateways(nw)
		if err != nil {
			glog.Errorf("error getting gateways for network %v: %v\n", nw, err)
			continue
		}
		numUpGateways := 0
		for _, gw := range gateways {
			req := protos.GatewayStatusRequest{NetworkId: nw, LogicalId: gw}
			status, err := reporter.Store.GetGatewayStatus(&req)
			if err != nil {
				glog.V(2).Infof("error getting status for gateway %v: %v\n", gw, err)
				continue
			}
			// last check in more than 5 minutes ago
			if (time.Now().Unix() - int64(status.Time)/1000) > 60*5 {
				gwCheckinStatus.WithLabelValues(nw, gw).Set(0)
			} else {
				gwCheckinStatus.WithLabelValues(nw, gw).Set(1)
				numUpGateways += 1
			}
			// report mconfig age
			mconfigCreatedAt := status.Checkin.GetPlatformInfo().GetConfigInfo().GetMconfigCreatedAt()
			if mconfigCreatedAt != 0 {
				gwMconfigAge.WithLabelValues(nw, gw).Set(float64(status.Time/1000 - mconfigCreatedAt))
			}
		}
		upGwCount.WithLabelValues(nw).Set(float64(numUpGateways))
		totalGwCount.WithLabelValues(nw).Set(float64(len(gateways)))
	}
	return nil
}
