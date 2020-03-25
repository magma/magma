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

	platformregistry "magma/orc8r/lib/go/registry"

	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const (
	ServiceName = "DIRECTORYD"
	ImsiPrefix  = "IMSI"
)

// Get a thin RPC client to the gateway directory service.
func GetGatewayDirectorydClient() (protos.GatewayDirectoryServiceClient, error) {
	conn, err := platformregistry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewGatewayDirectoryServiceClient(conn), nil
}

// UpdateRecord updates the directory record for the provided ID with the calling
// GW's HW ID and any associated identifiers
func UpdateRecord(request *protos.UpdateRecordRequest) error {
	if len(request.GetId()) == 0 {
		return fmt.Errorf("Empty ID")
	}
	client, err := GetGatewayDirectorydClient()
	if err != nil {
		return err
	}
	request.Id = PrependImsiPrefix(request.GetId())
	_, err = client.UpdateRecord(context.Background(), request)
	if err != nil {
		glog.Error(err)
	}
	return err
}

// DeleteRecord deletes the directory record for the provided ID
func DeleteRecord(request *protos.DeleteRecordRequest) error {
	if len(request.GetId()) == 0 {
		return fmt.Errorf("Empty ID")
	}
	client, err := GetGatewayDirectorydClient()
	if err != nil {
		return err
	}
	request.Id = PrependImsiPrefix(request.GetId())
	_, err = client.DeleteRecord(context.Background(), request)
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
