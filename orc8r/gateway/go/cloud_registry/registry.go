/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// package cloud_registry provides CloudRegistry interface for Go based gateways
package cloud_registry

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

const (
	grpcMaxTimeoutSec       = 60
	grpcMaxDelaySec         = 10
	controlProxyServiceName = "CONTROL_PROXY"
)

// control proxy config map
// it'll be initialized on demand and used thereafter
// any changed to the control proxy service config file would require process restart to take effect
var controlProxyConfig atomic.Value

type CloudRegistry interface {
	GetCloudConnection(service string) (*grpc.ClientConn, error)
	GetCloudConnectionFromServiceConfig(serviceConfig *config.ConfigMap, service string) (*grpc.ClientConn, error)
}

type ProxiedCloudRegistry struct{}

func New() *ProxiedCloudRegistry {
	return &ProxiedCloudRegistry{}
}

// GetCloudConnection returns a connection to a cloud service through the control
// proxy setup
// Input: service - name of cloud service to connect to
//
// Output: *grpc.ClientConn with connection to cloud service
//         error if it exists
func (cr *ProxiedCloudRegistry) GetCloudConnection(service string) (*grpc.ClientConn, error) {
	cpc, ok := controlProxyConfig.Load().(*config.ConfigMap)
	if (!ok) || cpc == nil {
		var err error
		// moduleName is "" since all feg configs lie in /etc/magma/configs without a module name
		cpc, err = config.GetServiceConfig("", "control_proxy")
		if err != nil {
			return nil, err
		}
		controlProxyConfig.Store(cpc)
	}
	return cr.GetCloudConnectionFromServiceConfig(cpc, service)
}

// GetCloudConnectionFromServiceConfig returns a connection to the cloud
// using a specific service config map. This map must contain the cloud_address
// and local_port params
// Input: serviceConfig - ConfigMap containing cloud_address and local_port
//        service - name of cloud service to connect to
//
// Output: *grpc.ClientConn with connection to cloud service
//         error if it exists
func (*ProxiedCloudRegistry) GetCloudConnectionFromServiceConfig(
	serviceConfig *config.ConfigMap, service string) (*grpc.ClientConn, error) {

	authority, err := getAuthority(serviceConfig, service)
	if err != nil {
		return nil, err
	}
	addr, err := getProxyAddress(serviceConfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcMaxTimeoutSec*time.Second)
	defer cancel()

	opts := getDialOptions(authority)
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("Address: %s GRPC Dial error: %s", addr, err)
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return conn, nil
}

// Returns the RPC address of the service.
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	return platform_registry.GetServiceAddress(service)
}

func getAuthority(
	serviceConfig *config.ConfigMap,
	service string,
) (string, error) {
	cloudAddr, err := serviceConfig.GetStringParam("cloud_address")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", service, cloudAddr), nil
}

func getProxyAddress(serviceConfig *config.ConfigMap) (string, error) {
	localPort, err := serviceConfig.GetIntParam("local_port")
	if err != nil {
		return "", err
	}
	localAddress, err := GetServiceAddress(controlProxyServiceName)
	if err != nil {
		return "", err
	}
	addrPieces := strings.Split(localAddress, ":")
	return fmt.Sprintf("%s:%d", addrPieces[0], localPort), nil
}

func getDialOptions(authority string) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithBackoffMaxDelay(grpcMaxDelaySec * time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithAuthority(authority),
	}
}
