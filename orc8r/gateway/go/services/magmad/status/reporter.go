/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package status implements magmad status collector & reporter
package status

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"magma/gateway/config"
	"magma/gateway/mconfig"
	"magma/gateway/service_registry"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
	mconfig_proto "magma/orc8r/lib/go/protos/mconfig"
)

const (
	DefaultCheckinIntervalSeconds int32 = 60
	MinCheckinIntervalSeconds           = 30

	serviceCollectDelay = time.Second * 10
)

// StartReporter starts state collection & reporting loop
// StartReporter never returns, it'll log errors if any and continue
func StartReporter() {
	timer := time.NewTimer(time.Second * time.Duration(DefaultCheckinIntervalSeconds))
	for {
		mdc := config.GetMagmadConfigs()
		fb303services := mdc.MagmaServices
		nonFb303Services := map[string]struct{}{}
		for _, s := range mdc.NonService303Services {
			nonFb303Services[s] = struct{}{}
		}
		for _, fb303service := range fb303services {
			if _, nonFb303 := nonFb303Services[fb303service]; !nonFb303 {
				if err := startServiceQuery(fb303service); err != nil {
					log.Printf("error querying service '%s' state: %v", fb303service, err)
				}
			}
		}
		time.Sleep(serviceCollectDelay)

		stateConn, err := service_registry.Get().GetSharedCloudConnection(definitions.StateServiceName)
		if err != nil {
			log.Printf("failed to connect to state reporting service: %v", err)
		} else {
			res, err := protos.NewStateServiceClient(stateConn).ReportStates(context.Background(), collect())
			if err != nil {
				log.Printf("ReportStates error: %v", err)
			} else if len(res.GetUnreportedStates()) > 0 {
				resStr, _ := json.Marshal(res.GetUnreportedStates())
				log.Printf("status unreported states: %s", resStr)
			}
		}
		<-timer.C // wait for timer

		// update timer based on the latest configs
		intervalSeconds := DefaultCheckinIntervalSeconds
		mconf := &mconfig_proto.MagmaD{}
		err = mconfig.GetServiceConfigs(definitions.MagmadServiceName, mconf)
		if err == nil && mconf.CheckinInterval != 0 {
			intervalSeconds = mconf.CheckinInterval
		}
		if intervalSeconds < MinCheckinIntervalSeconds {
			intervalSeconds = MinCheckinIntervalSeconds
		}
		timer.Reset(time.Second * time.Duration(intervalSeconds))
	}
}
