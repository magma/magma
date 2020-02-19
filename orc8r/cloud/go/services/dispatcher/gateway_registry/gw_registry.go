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

	"magma/orc8r/cloud/go/services/dispatcher"
	"magma/orc8r/lib/go/registry"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// Gateway services
	GwMobilityd           GwServiceType = "mobilityd"
	GwMagmad              GwServiceType = "magmad"
	GwEnodebd             GwServiceType = "enodebd"
	GwPipelined           GwServiceType = "pipelined"
	GwSubscriberDB        GwServiceType = "subscriberdb"
	GwS6aService          GwServiceType = "s6a_service"
	GwSgsService          GwServiceType = "sgs_service"
	GwSessiondService     GwServiceType = "sessiond"
	GwSpgwService         GwServiceType = "spgw_service"
	GwAbortSessionService GwServiceType = "abort_session_service"
	GwAAAService          GwServiceType = "aaa_server"

	// SyncRPC gateway header key
	GatewayIdHeaderKey = "Gatewayid"

	HttpServerAddressPort = 9080
)

type GwServiceType string

type httpServerConfig struct {
	port int
	*sync.RWMutex
}

var services = []GwServiceType{
	GwMobilityd,
	GwMagmad,
	GwEnodebd,
	GwPipelined,
	GwSubscriberDB,
	GwS6aService,
	GwSgsService,
	GwSessiondService,
	GwSpgwService,
	GwAbortSessionService,
	GwAAAService,
}

var config = httpServerConfig{HttpServerAddressPort, &sync.RWMutex{}}

// SetPort sets the port of http_server.
// If a port is already set, this overrides the previous setting.
func SetPort(port int) error {
	config.Lock()
	config.port = port
	config.Unlock()
	return nil
}

// GetServiceAddressForGateway returns the ip addr for the
// SyncRPCHTTPServer instance, which is in the same process
// of the Dispatcher grpc server who has an open bidirectional
// stream with the gateway with hwId.
func GetServiceAddressForGateway(hwId string) (string, error) {
	hostName, err := dispatcher.GetHostnameForHwid(hwId)
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

// GetGatewayConnection gets a connection to the SyncRPC HTTP server
// who can forward the message to the corresponding gateway.
//
// Returns a connection and a context that should be based on for rpc calls on this connection.
// The context will put the Gatewayid in its metadata, which will be surfaced as HTTP/2 headers.
func GetGatewayConnection(service GwServiceType, hwId string) (*grpc.ClientConn, context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), registry.GrpcMaxTimeoutSec*time.Second)
	defer cancel()
	addr, err := GetServiceAddressForGateway(hwId)
	if err != nil {
		return nil, nil, err
	}
	conn, err := registry.GetClientConnection(
		ctx,
		addr,
		grpc.WithBackoffMaxDelay(registry.GrpcMaxDelaySec*time.Second),
		grpc.WithBlock(),
		grpc.WithAuthority(string(service)),
	)
	if err != nil {
		err = fmt.Errorf("Service %v connection error: %v", service, err)
		return nil, nil, err
	}
	customHeader := metadata.New(map[string]string{GatewayIdHeaderKey: hwId})
	ctxToRet := metadata.NewOutgoingContext(context.Background(), customHeader)
	return conn, ctxToRet, nil

}

func ListAllGwServices() []GwServiceType {
	return services
}
