/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// package cloud_registry provides CloudRegistry interface for Go based gateways
package cloud_registry

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"

	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

const (
	grpcMaxTimeoutSec       = 60
	grpcMaxDelaySec         = 20
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
//        and optional proxy_cloud_connections, cloud_port, rootca_cert, gateway_cert/key fields if direct
//        cloud connection is needed
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
	useProxy := getUseProxyCloudConnection(serviceConfig)
	var addr string
	if useProxy {
		addr, err = getProxyAddress(serviceConfig)
	} else {
		addr, err = getCloudServiceAddress(serviceConfig)
	}
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcMaxTimeoutSec*time.Second)
	defer cancel()

	opts, err := getDialOptions(serviceConfig, authority, useProxy)
	if err != nil {
		return nil, err
	}
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

func getCloudServiceAddress(controlProxyConfig *config.ConfigMap) (string, error) {
	cloudAddr, err := controlProxyConfig.GetStringParam("cloud_address")
	if err != nil {
		return "", err
	}
	addrPieces := strings.Split(cloudAddr, ":")
	port, err := controlProxyConfig.GetIntParam("cloud_port")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", addrPieces[0], port), nil
}

func getUseProxyCloudConnection(serviceConfig *config.ConfigMap) bool {
	if proxied, err := serviceConfig.GetBoolParam("proxy_cloud_connections"); err == nil {
		return proxied
	}
	return true // Note: default is True -> proxy cloud connection
}

func getDialOptions(serviceConfig *config.ConfigMap, authority string, useProxy bool) ([]grpc.DialOption, error) {
	bckoff := backoff.DefaultConfig
	bckoff.MaxDelay = grpcMaxDelaySec * time.Second
	var opts = []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           bckoff,
			MinConnectTimeout: grpcMaxTimeoutSec * time.Second,
		}),
		grpc.WithBlock(),
	}
	if useProxy {
		opts = append(opts, grpc.WithInsecure(), grpc.WithAuthority(authority))
	} else {
		// always try to add OS certs
		certPool, err := x509.SystemCertPool()
		if err != nil {
			log.Printf("OS Cert Pool initialization error: %v", err)
			certPool = x509.NewCertPool()
		}
		if rootCaFile, err := serviceConfig.GetStringParam("rootca_cert"); err == nil && len(rootCaFile) > 0 {
			// Add magma RootCA
			if rootCa, err := ioutil.ReadFile(rootCaFile); err == nil {
				if !certPool.AppendCertsFromPEM(rootCa) {
					log.Printf("Failed to append certificates from %s", rootCaFile)
				}
			} else {
				log.Printf("Cannot load Root CA from '%s': %v", rootCaFile, err)
			}
		}
		tlsCfg := &tls.Config{ServerName: authority}
		if len(certPool.Subjects()) > 0 {
			tlsCfg.RootCAs = certPool
		} else {
			tlsCfg.InsecureSkipVerify = true
		}
		if clientCaFile, err := serviceConfig.GetStringParam("gateway_cert"); err == nil && len(clientCaFile) > 0 {
			if clientKeyFile, err := serviceConfig.GetStringParam("gateway_key"); err == nil && len(clientKeyFile) > 0 {
				clientCert, err := tls.LoadX509KeyPair(clientCaFile, clientKeyFile)
				if err == nil {
					tlsCfg.Certificates = []tls.Certificate{clientCert}
				} else {
					log.Printf("failed to load Client Certificate/Key from '%s', '%s': %v",
						clientCaFile, clientKeyFile, err)
				}
			}
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	}
	return opts, nil
}
