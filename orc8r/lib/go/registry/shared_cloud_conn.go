/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// package registry provides Registry interface for Go based gateways
// as well as cloud connection routines
package registry

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"magma/orc8r/lib/go/service/config"
)

// DefaultSharedCloudConnectionTTL - default duration to reuse the same connection
const DefaultSharedCloudConnectionTTL = time.Hour * 4

// sharedCloudConnectionTTL - default duration to reuse the same connection the TTL should be a fraction of gateway
// certificate renewal period to make sure, the connection is not reused past the certificate expiration
var sharedCloudConnectionTTL = DefaultSharedCloudConnectionTTL

// GetSharedCloudConnection returns a new GRPC service connection to the service in the cloud for a gateway
// either directly or via control proxy
// GetSharedCloudConnection will return an existing cached cloud connection if it's available and healthy,
// if not - it'll try to create, cache and return a new cloud connection
// Input: service - name of cloud service to connect to
//
// Output: *grpc.ClientConn with connection to cloud service
//         error if it exists
func (registry *ServiceRegistry) GetSharedCloudConnection(service string) (*grpc.ClientConn, error) {
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
	return registry.GetSharedCloudConnectionFromServiceConfig(cpc, service)
}

// GetSharedCloudConnectionFromServiceConfig returns a connection to the cloud
// using a specific control_proxy service config map. This map must contain the cloud_address
// and local_port params
// GetSharedCloudConnectionFromServiceConfig will return an existing cached cloud connection if it's available and
// healthy, if not - it'll try to create, cache and return a new cloud connection
// Input:  serviceConfig - ConfigMap containing cloud_address and local_port
//         and optional proxy_cloud_connections, cloud_port, rootca_cert, gateway_cert/key fields if direct
//         cloud connection is needed
//         service - name of cloud service to connect to
//
// Output: *grpc.ClientConn with connection to cloud service
//         error if it exists
// Note:   controlProxyConfig differences are ignored in cached connection mapping,
//         if an update to ConfigMap is required - use CleanupSharedCloudConnection() to flush the service conn cache
func (registry *ServiceRegistry) GetSharedCloudConnectionFromServiceConfig(
	controlProxyConfig *config.ConfigMap, service string) (*grpc.ClientConn, error) {

	// First try to get an existing connection with reader lock
	registry.cloudConnMu.RLock()
	conn, connExists := registry.cloudConnections[service]
	registry.cloudConnMu.RUnlock()

	timeNow := time.Now()
	if connExists && (conn.ClientConn != nil) {
		if conn.GetState() == connectivity.Ready && conn.expiration.After(timeNow) {
			return conn.ClientConn, nil // cached connection is good & current - return it
		}
	}

	// Attempt to connect outside of the lock
	grpcConn, connectErr := registry.GetCloudConnectionFromServiceConfig(controlProxyConfig, service)
	if connectErr != nil || grpcConn == nil || grpcConn.GetState() != connectivity.Ready {
		connectErr = fmt.Errorf("cloud service '%s' connection error: %v", service, connectErr)
	}

	newConnectionTTL := GetSharedCloudConnectionTTL()

	// GRPC connection is successfully established, update the cache
	registry.cloudConnMu.Lock() // LOCK

	// check if another thread already created/updated the shared connection for the service
	conn, connExists = registry.cloudConnections[service]
	if connExists && conn.ClientConn != nil {
		if connState := conn.GetState(); connState == connectivity.Ready {
			if conn.expiration.After(timeNow) {
				// another thread already created/updated the shared connection for the service,
				// return it & close just created grpcConn, it's no longer needed
				registry.cloudConnMu.Unlock() // UNLOCK

				if connectErr == nil {
					grpcConn.Close()
				}
				return conn.ClientConn, nil
			} else {
				if connectErr != nil {
					// we failed to create a new cloud connection, better to use expired but otherwise valid
					// connection and retry later then fail the call
					registry.cloudConnMu.Unlock() // UNLOCK

					log.Printf(
						"failed to create a new service '%s' cloud connection: %v; using expired",
						service, connectErr)
					return conn.ClientConn, nil
				}
				// connection is expired, but is still valid, give existing users time to complete
				// and defer + delay connection close call
				defer func() {
					go func() {
						time.Sleep(time.Second * grpcMaxDelaySec)
						conn.Close()
					}()
				}()
				log.Printf("service '%s' cloud connection is expired, will reconnect", service)
			}
		} else { // if connState := conn.GetState(); connState == connectivity.Ready ... else
			// connection is already broken, close on return without delay
			defer conn.Close()
			log.Printf("unhealthy state '%s' of service '%s' cloud connection", connState, service)
		}
	}
	if connectErr == nil {
		registry.cloudConnections[service] =
			cloudConnection{ClientConn: grpcConn, expiration: timeNow.Add(newConnectionTTL)}
	} else {
		delete(registry.cloudConnections, service)
	}

	registry.cloudConnMu.Unlock() // UNLOCK

	return grpcConn, connectErr
}

// CleanupSharedCloudConnection removes cached cloud connection for the service from cache and closes it
// Returns true if connection was cached
func (registry *ServiceRegistry) CleanupSharedCloudConnection(service string) bool {
	registry.cloudConnMu.Lock()
	conn, ok := registry.cloudConnections[service]
	if ok {
		delete(registry.cloudConnections, service)
		if ok = conn.ClientConn != nil; ok {
			defer conn.ClientConn.Close()
		}
	}
	registry.cloudConnMu.Unlock()
	return ok
}

// SetSharedCloudConnectionTTL atomically sets Shared Cloud Connection TTL
// Note: the new TTL will apply only to newly created connections,
// existing cached connections will not be affected
func SetSharedCloudConnectionTTL(ttl time.Duration) {
	atomic.StoreInt64((*int64)(&sharedCloudConnectionTTL), int64(ttl))
}

// GetSharedCloudConnectionTTL atomically gets and returns current Shared Cloud Connection TTL value
func GetSharedCloudConnectionTTL() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&sharedCloudConnectionTTL)))
}
