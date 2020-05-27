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
	"fmt"
	"log"
	"strings"
	"sync"

	"magma/gateway/service_registry"
	"magma/gateway/status"
	"magma/orc8r/lib/go/protos"
)

type serviceInfoClientState struct {
	queryInFlight bool
	states        []*protos.State
}

var (
	pollerMu     sync.Mutex
	clientStates map[string]*serviceInfoClientState = map[string]*serviceInfoClientState{}
)

func startServiceQuery(service string, queryMetrics bool) error {
	conn, err := service_registry.Get().GetConnection(strings.ToUpper(service))
	if err != nil {
		return err
	}
	client := protos.NewService303Client(conn)

	pollerMu.Lock()
	defer pollerMu.Unlock()

	clState, ok := clientStates[service]
	if !ok || clState == nil {
		clState = &serviceInfoClientState{}
		clientStates[service] = clState
	}
	if clState.queryInFlight {
		return fmt.Errorf("'%s' fb303 client is still in use", service)
	}
	clState.queryInFlight = true

	go func() {
		var (
			serviceMetrics *protos.MetricsContainer
			merr           error
		)

		ctx := context.Background()
		statesResp, err := client.GetOperationalStates(ctx, &protos.Void{})
		if err != nil {
			log.Printf("service '%s' GetServiceInfo error: %v", service, err)
		}
		if queryMetrics {
			serviceMetrics, merr = client.GetMetrics(ctx, &protos.Void{})
			if merr != nil {
				log.Printf("service '%s' GetMetrics error: %v", service, err)
			}
		}
		pollerMu.Lock()
		clState.queryInFlight = false
		if err == nil {
			clState.states = statesResp.States
		}
		pollerMu.Unlock()

		if queryMetrics && merr == nil {
			enqueueMetrics(service, serviceMetrics)
		}
	}()
	return nil
}

func collect() *protos.ReportStatesRequest {
	marshaledGwState, err := json.Marshal(status.GetGatewayStatus())
	if err != nil {
		marshaledGwState, _ = json.Marshal(&status.GatewayStatus{
			Meta: map[string]string{"error": err.Error()},
		})
	}
	gwState := &protos.State{
		Type:     "gw_state",
		DeviceID: status.GetHwId(),
		Value:    marshaledGwState,
	}

	pollerMu.Lock()
	states := make([]*protos.State, 0, 1)
	states = append(states, gwState)
	for _, clientState := range clientStates {
		if clientState != nil && len(clientState.states) > 0 {
			states = append(states, clientState.states...)
			clientState.states = nil
		}
	}
	pollerMu.Unlock()

	return &protos.ReportStatesRequest{States: states}
}
