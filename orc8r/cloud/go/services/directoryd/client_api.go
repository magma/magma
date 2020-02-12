/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package directoryd provides a client API for interacting with the
// directory service, which manages the UE location information
package directoryd

import (
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const ServiceName = "DIRECTORYD"

// Get a thin RPC client to the directory service.
func GetDirectorydClient() (protos.DirectoryServiceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewDirectoryServiceClient(conn), err
}

func GetHardwareIdByIMSI(imsi string, networkId string) (string, error) {
	return getLocation(protos.TableID_IMSI_TO_HWID, imsi, networkId)
}

func GetHostNameByIMSI(hwId string) (string, error) {
	return getLocation(protos.TableID_HWID_TO_HOSTNAME, hwId, "")
}

func getLocation(tableId protos.TableID, recordId string, networkId string) (string, error) {
	client, err := GetDirectorydClient()
	if err != nil {
		return "", err
	}
	req := &protos.GetLocationRequest{
		Table:     tableId,
		Id:        recordId,
		NetworkID: networkId,
	}
	ctx := context.Background()
	record, err := client.GetLocation(ctx, req)
	if err != nil {
		return "", err
	}
	return record.Location, nil
}

func UpdateHostNameByHwId(hwId string, hostName string) error {
	return updateLocation(protos.TableID_HWID_TO_HOSTNAME, hwId, hostName)
}

func updateLocation(tableId protos.TableID, recordId string, location string) error {
	client, err := GetDirectorydClient()
	if err != nil {
		return err
	}
	req := &protos.UpdateDirectoryLocationRequest{
		Table:  tableId,
		Id:     recordId,
		Record: &protos.LocationRecord{Location: location},
	}
	ctx := context.Background()
	_, err = client.UpdateLocation(ctx, req)
	return err
}

func DeleteHardwareIdByIMSI(imsi string, networkID string) error {
	return deleteLocation(protos.TableID_IMSI_TO_HWID, imsi, networkID)
}

func DeleteHostNameByIMSI(hwId string) error {
	return deleteLocation(protos.TableID_HWID_TO_HOSTNAME, hwId, "")
}

func deleteLocation(tableId protos.TableID, recordId string, networkId string) error {
	client, err := GetDirectorydClient()
	if err != nil {
		return err
	}
	req := &protos.DeleteLocationRequest{
		Table:     tableId,
		Id:        recordId,
		NetworkID: networkId,
	}
	ctx := context.Background()
	_, err = client.DeleteLocation(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
