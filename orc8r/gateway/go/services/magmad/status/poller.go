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

type serviceInfoClient struct {
	client protos.Service303Client
	states []*protos.State
}

var (
	pollerMu sync.Mutex
	clients  map[string]*serviceInfoClient = map[string]*serviceInfoClient{}
)

func startServiceQuery(service string) error {
	conn, err := service_registry.Get().GetConnection(strings.ToUpper(service))
	if err != nil {
		return err
	}
	client := protos.NewService303Client(conn)

	pollerMu.Lock()
	defer pollerMu.Unlock()
	cl, ok := clients[service]
	if !ok || cl == nil {
		cl = &serviceInfoClient{}
		clients[service] = cl
	}
	if cl.client != nil {
		return fmt.Errorf("'%s' fb303 client is still in use", service)
	}
	cl.client = client
	go func() {
		statesResp, err := cl.client.GetOperationalStates(context.Background(), &protos.Void{})

		pollerMu.Lock()
		defer pollerMu.Unlock()

		cl.client = nil
		if err != nil {
			log.Printf("service '%s' GetServiceInfo error: %v", service, err)
		} else {
			cl.states = statesResp.States
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
	states := make([]*protos.State, 0, len(clients)+1)
	states = append(states, gwState)
	for _, client := range clients {
		if client != nil && client.states != nil {
			states = append(states, client.states...)
			client.states = nil
		}
	}
	pollerMu.Unlock()

	return &protos.ReportStatesRequest{States: states}
}
