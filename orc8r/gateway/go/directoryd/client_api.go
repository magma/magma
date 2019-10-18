/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package directoryd provides a client API for interacting with the
// directory cloud service, which manages the UE location information
package directoryd

import (
	"fmt"
	"strings"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/gateway/cloud_registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const (
	ServiceName = "DIRECTORYD"
	ImsiPrefix = "IMSI"
)

// Get a thin RPC client to the directory service.
func GetDirectorydClient() (protos.DirectoryServiceClient, error) {
	conn, err := cloud_registry.New().GetCloudConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewDirectoryServiceClient(conn), nil
}

// AddIMSI associates given IMSI (UE) with the calling GW's HW ID
func AddIMSI(imsi string) error {
	if len(imsi) == 0 {
		return fmt.Errorf("Empty IMSI")
	}
	client, err := GetDirectorydClient()
	if err != nil {
		return err
	}

	req := &protos.UpdateDirectoryLocationRequest{
		Table:  protos.TableID_IMSI_TO_HWID,
		Id:     PrependImsiPrefix(imsi),
		// request Record will be populated by directoryd cloud service
	}
	ctx := context.Background()
	_, err = client.UpdateLocation(ctx, req)
	if err != nil {
		glog.Error(err)
	}
	return err
}

// RemoveIMSI disassociates given IMSI (UE) from the calling GW's HW ID
func RemoveIMSI(imsi string) error {
	if len(imsi) == 0 {
		return fmt.Errorf("Empty IMSI")
	}
	client, err := GetDirectorydClient()
	if err != nil {
		return err
	}
	req := &protos.DeleteLocationRequest{
		Table: protos.TableID_IMSI_TO_HWID,
		Id:    PrependImsiPrefix(imsi),
	}
	ctx := context.Background()
	_, err = client.DeleteLocation(ctx, req)
	if err != nil {
		glog.Error(err)
	}
	return err
}

func PrependImsiPrefix(imsi string) string {
	if !strings.HasPrefix(imsi, ImsiPrefix) {
		imsi = ImsiPrefix + imsi
	}
	return imsi
}