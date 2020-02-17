/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dispatcher

import (
	"context"

	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	platformregistry "magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

const ServiceName = "DISPATCHER"

// GetHostnameForHwid returns the controller hostname mapped for the hwid.
func GetHostnameForHwid(hwid string) (string, error) {
	client, err := getDispatcherClient()
	if err != nil {
		return "", err
	}
	hostname, err := client.GetHostnameForHwid(context.Background(), &protos.HardwareID{Hwid: hwid})
	return hostname.GetName(), err
}

// getDispatcherClient returns a new RPC client for the dispatcher service.
func getDispatcherClient() (protos.SyncRPCServiceClient, error) {
	conn, err := platformregistry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewSyncRPCServiceClient(conn), nil
}
