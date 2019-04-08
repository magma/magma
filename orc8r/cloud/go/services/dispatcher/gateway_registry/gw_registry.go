/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gateway_registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/directoryd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// gateway services
	GWMOBILITYD       GwServiceType = "mobilityd"
	GWMAGMAD          GwServiceType = "magmad"
	GWENODEBD         GwServiceType = "enodebd"
	GWPIPELINED       GwServiceType = "pipelined"
	GWSUBSCRIBERDB    GwServiceType = "subscriberdb"
	GWS6ASERVICE      GwServiceType = "s6a_service"
	GWSGSSERVICE      GwServiceType = "sgs_service"
	GWSESSIONDSERVICE GwServiceType = "sessiond"

	// syncRPC gateway header key
	GATEWAYID_HEADER_KEY = "Gatewayid"

	HTTP_SERVER_ADDRESS_PORT = 9080
)

type GwServiceType string

var services []GwServiceType

func init() {
	services = append(services, GWMOBILITYD)
	services = append(services, GWMAGMAD)
	services = append(services, GWENODEBD)
	services = append(services, GWPIPELINED)
	services = append(services, GWSUBSCRIBERDB)
	services = append(services, GWS6ASERVICE)
	services = append(services, GWSGSSERVICE)
	services = append(services, GWSESSIONDSERVICE)
}

type httpServerConfig struct {
	port int
	*sync.RWMutex
}

var (
	config = httpServerConfig{HTTP_SERVER_ADDRESS_PORT,
		&sync.RWMutex{}}
)

// SetPort sets the port of http_server.
// If a port is already set, this overrides the previous setting.
func SetPort(port int) error {
	config.Lock()
	config.port = port
	config.Unlock()
	return nil
}

// Returns the ip addr for the SyncRPCHTTPServer instance, which is
// in the same process of the Dispatcher grpc server who has an open
// bidirectional stream with the gateway with hwId.
func GetServiceAddressForGateway(hwId string) (string, error) {
	hostName, err := directoryd.GetHostNameByIMSI(hwId)
	if err != nil {
		fmt.Printf("err getting hostName in GetServiceAddressForGateway for hwId %v: %v\n", hwId, err)
		return "", err
	}
	config.RLock()
	port := config.port
	config.RUnlock()
	addr := fmt.Sprintf("%s:%v", hostName, port)
	return addr, nil
}

// get a connection to the SYNCRPCHTTPSERVER who can forward the message
// to the corresponding gateway
// return a connection, and a context that should be based on for rpc calls on this connection.
// The context will put the Gatewayid in its metadata, which will be surfaced as http/2 headers
func GetGatewayConnection(service GwServiceType, hwId string) (*grpc.ClientConn, context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), registry.GrpxMaxTimeoutSec*time.Second)
	defer cancel()
	addr, err := GetServiceAddressForGateway(hwId)
	if err != nil {
		return nil, nil, err
	}
	conn, err := registry.GetClientConnection(ctx, addr, grpc.WithBackoffMaxDelay(registry.GrpcMaxDelaySec*time.Second),
		grpc.WithBlock(), grpc.WithAuthority(string(service)))
	if err != nil {
		err = fmt.Errorf("Service %v connection error: %v", service, err)
		return nil, nil, err
	}
	customHeader := metadata.New(map[string]string{GATEWAYID_HEADER_KEY: hwId})
	ctxToRet := metadata.NewOutgoingContext(context.Background(), customHeader)
	return conn, ctxToRet, nil

}

func ListAllGwServices() []GwServiceType {
	return services
}
