/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package servicers implements GRPC MagmadConfigurator service for controlling,
// managing and configuring Networks & Access Gateways associated with the cloud
// MagmadConfigurator relies on three datastore based key/value tables:
//
// 1. AG Logical ID ==> Managed AG Config:
// 	Serialized configs are stored in a datastore as marshaled protobuf JSON,
// 	keyed by AccessGateway Logical ID
// 2. AG Logical ID ==> AG Info Record:
// 	AG Information records are internal to cloud stored as Go globs keyed by
// 	AG Logical IDs
// 3. AG Hardware ID ==> AG Logical ID:
// 	Logical IDs are stored directly as byte slice and keyed by AG HW IDs
package servicers

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/protos"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
)

// Magmad Store
type MagmadConfigurator struct {
	Store datastore.Api
}

func NewMagmadConfigurator(store datastore.Api) *MagmadConfigurator {
	return &MagmadConfigurator{store}
}

// Magmad Store consists of three tables:
//
// 1. AG Logical ID ==> Managed AG Config:
// 	Serialized configs are stored in a datastore as marshaled protobuf JSON,
// 	keyed by AccessGateway Logical ID
// 2. AG Logical ID ==> AG Info Record:
// 	AG Information records are internal to cloud stored as Go globs keyed by
// 	AG Logical IDs
// 3. AG Hardware ID ==> AG Logical ID:
// 	Logical IDs are stored directly as byte slice and keyed by AG HW IDs

// Supported Magmad Configurator Tables
const (
	AgRecordTableName = "gatewayRecords"
	HwIdTableName     = "hwIds"
	NetworksTableName = "networks"
	GatewaysTableName = "gateways" // HW Gw Id => Network Id table

	// the following tables are currently used for cleanup (network delete)
	// until we have cloud service registry to manage inter service connections
	SubscribersTableName    = "subscriberdb"
	GatewaysStatusTableName = "gwstatus"
	TierTableName           = "tierVersions"
)

func (md MagmadConfigurator) deleteIfExists(tableName string, key string) error {
	exists, err := md.Store.DoesKeyExist(tableName, key)
	if err != nil {
		return err
	}

	if exists {
		return md.Store.Delete(tableName, key)
	}
	return nil
}
func getTablesToDropForNetworkDeletion(networkId string) []string {
	return []string{
		datastore.GetTableName(networkId, AgRecordTableName),
		datastore.GetTableName(networkId, HwIdTableName),
		datastore.GetTableName(networkId, SubscribersTableName),
		datastore.GetTableName(networkId, GatewaysStatusTableName),
		datastore.GetTableName(networkId, TierTableName),
	}
}

func (md MagmadConfigurator) getAllHWIdsInNetwork(networkId string) ([]string, error) {
	tableName := datastore.GetTableName(networkId, AgRecordTableName)
	allGwIds, err := md.Store.ListKeys(tableName)
	if err != nil {
		return nil, err
	}
	allAgRecordsMarshaled, err := md.Store.GetMany(tableName, allGwIds)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(allGwIds))
	for _, marshaledAgRecord := range allAgRecordsMarshaled {
		agRecord := &magmadprotos.AccessGatewayRecord{}
		err = protos.Unmarshal(marshaledAgRecord.Value, agRecord)
		if err != nil {
			return nil, err
		}
		ret = append(ret, agRecord.GetHwId().GetId())
	}
	return ret, nil
}

func (md MagmadConfigurator) validateNetworkExists(networkId string) error {
	exist, err := md.Store.DoesKeyExist(NetworksTableName, networkId)
	if err != nil {
		return fmt.Errorf(
			"Could not validate if %s exists", networkId,
		)
	}
	if !exist {
		return fmt.Errorf(
			"Invalid Network ID '%s'; network may not have been registered",
			networkId)
	}
	return err
}

func (md MagmadConfigurator) validateHwIdNotRegisteredInNetwork(networkId string, hwId string) error {
	exist, err := md.Store.DoesKeyExist(datastore.GetTableName(networkId, HwIdTableName), hwId)
	if err != nil {
		return fmt.Errorf("Could not validate if hwId %s is registered", hwId)
	}
	if exist {
		return fmt.Errorf(
			"Hardware ID '%s' Already Registered",
			hwId,
		)
	} else {
		return nil
	}
}

func (md MagmadConfigurator) validateUniqueLogicalId(networkId string, logicalId string) error {
	exist, err := md.Store.DoesKeyExist(datastore.GetTableName(networkId, AgRecordTableName), logicalId)
	if err != nil {
		return fmt.Errorf("Could not determine if %s exists already", logicalId)
	}
	if exist {
		return fmt.Errorf("Requested ID '%s' Already Registered", logicalId)

	} else {
		return nil
	}
}

func (md MagmadConfigurator) validateHwIdNotRegistered(hwId string) error {
	gwNetworkIdBytes, _, err := md.Store.Get(GatewaysTableName, hwId)
	if err != nil {
		return nil
	} else {
		return fmt.Errorf(
			"Gateway Hardware ID '%s' is already Registered on '%s' network",
			hwId,
			string(gwNetworkIdBytes),
		)
	}
}

// getGatewayWithNetwork verifies & returns valid Gateway Identity with network ID
func getGatewayWithNetwork(id *protos.Identity) (*protos.Identity_Gateway, error) {
	res := &protos.Identity_Gateway{}
	if id == nil {
		return res, status.Errorf(codes.InvalidArgument, "Nil Gateway Identity")
	}
	gwId := id.GetGateway()
	if gwId == nil {
		return res, status.Errorf(codes.InvalidArgument, "Invalid Gateway Identity")
	}
	res = gwId
	if len(res.NetworkId) == 0 {
		return res, status.Errorf(codes.InvalidArgument, "Empty Gateway Network ID")
	}
	return res, nil
}

func getGatewayWithNetworkAndHwId(id *protos.Identity) (*protos.Identity_Gateway, error) {
	res, err := getGatewayWithNetwork(id)
	if err == nil {
		return res, err
	}
	if len(res.HardwareId) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Empty Gateway Hardware ID")
	}
	return res, nil
}

func getGatewayWithLogicalId(id *protos.Identity) (*protos.Identity_Gateway, error) {
	res, err := getGatewayWithNetwork(id)
	if err == nil {
		return res, err
	}
	if len(res.LogicalId) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Empty Gateway Logical ID")
	}
	return res, nil
}

// getGatewayWithHardwareId verifies & returns valid Gateway Identity with HW ID
func getGatewayWithHardwareId(id *protos.Identity) (*protos.Identity_Gateway, error) {
	if id == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Nil Gateway Identity")
	}
	res := id.GetGateway()
	if res == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Gateway Identity")
	}
	if len(res.HardwareId) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Empty Gateway Hardware ID")
	}
	return res, nil
}
