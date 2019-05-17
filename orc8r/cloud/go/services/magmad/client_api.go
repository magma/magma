/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package magmad

import (
	"errors"
	"fmt"

	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/config"
	mdprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const ServiceName = "MAGMAD"

// getMagmadClient is a utility function to get a RPC connection to the
// magmad service
func getMagmadClient() (mdprotos.MagmadConfiguratorClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return mdprotos.NewMagmadConfiguratorClient(conn), err
}

// ListNetworks returns an array of all registered network IDs
func ListNetworks() ([]string, error) {
	md, err := getMagmadClient()
	if err != nil {
		return nil, err
	}
	idsList, err := md.ListNetworks(context.Background(), &protos.Void{})
	ids := idsList.GetList()
	res := make([]string, len(ids))
	for i, id := range ids {
		res[i] = id.GetNetwork()
	}
	return res, nil
}

// Returns the network record for a network ID
func GetNetwork(networkId string) (*mdprotos.MagmadNetworkRecord, error) {
	md, err := getMagmadClient()
	if err != nil {
		return nil, err
	}
	return md.GetNetwork(context.Background(), identity.NewNetwork(networkId))
}

// Update a network record
func UpdateNetwork(networkId string, record *mdprotos.MagmadNetworkRecord) error {
	md, err := getMagmadClient()
	if err != nil {
		return err
	}
	_, err = md.UpdateNetwork(context.Background(), &mdprotos.NetworkRecordRequest{Id: networkId, Record: record})
	return err
}

// Registers new network with ID specified in requestedId,
// will return error if the network already exist
// returns a new unique network ID
func RegisterNetwork(record *mdprotos.MagmadNetworkRecord, requestedId string) (string, error) {
	md, err := getMagmadClient()
	if err != nil {
		return "", err
	}
	req := &mdprotos.NetworkRecordRequest{Record: record, Id: requestedId}
	id, err := md.RegisterNetwork(context.Background(), req)
	return id.GetNetwork(), err
}

// Deletes given network if its Gateway & Subscriber tables are empty
func RemoveNetwork(networkId string) error {
	md, err := getMagmadClient()
	if err != nil {
		return err
	}
	_, err = md.RemoveNetwork(context.Background(), identity.NewNetwork(networkId))
	if err != nil {
		return err
	}

	// Also delete all configs for this network
	return config.DeleteAllNetworkConfigs(networkId)
}

// Deletes given network and its Gateway & Subscriber tables
func ForceRemoveNetwork(networkId string) error {
	md, err := getMagmadClient()
	if err != nil {
		return err
	}
	_, err = md.ForceRemoveNetwork(
		context.Background(), identity.NewNetwork(networkId))
	return err
}

// RegisterGateway adds new records for AG with HW ID specified by 'hwId'
// into all corresponding table, allocates a new Logical AG ID (first
// return parameter), initializes the newly created AG's configs to defaults,
// Will return error if a device with the given HW ID is already registered
func RegisterGateway(networkId string, record *mdprotos.AccessGatewayRecord) (string, error) {
	return RegisterGatewayWithId(networkId, record, "")
}

// Add new records for AG with id specified by requestedId
func RegisterGatewayWithId(networkId string, record *mdprotos.AccessGatewayRecord, requestedId string) (string, error) {
	if record == nil {
		return "", errors.New("Invalid Request: Nil Gateway Record")
	}
	md, err := getMagmadClient()
	if err != nil {
		return "", err
	}
	req := &mdprotos.GatewayRecordRequest{
		GatewayId: identity.NewGateway(record.GetHwId().GetId(), networkId, requestedId),
		Record:    record}

	gatewayId, err := md.RegisterGateway(context.Background(), req)
	if err != nil {
		return "", err
	}
	return gatewayId.GetGateway().GetLogicalId(), nil
}

// Lists all registered logical device IDs
func ListGateways(networkId string) ([]string, error) {
	md, err := getMagmadClient()
	if err != nil {
		return nil, err
	}
	idsList, err :=
		md.ListGateways(context.Background(), identity.NewNetwork(networkId))
	if err != nil {
		return nil, err
	}
	ids := idsList.GetList()
	gatewayIds := make([]string, len(ids))
	for i, id := range ids {
		gatewayIds[i] = id.GetGateway().GetLogicalId()
	}
	return gatewayIds, nil
}

// FindGatewayId returns logical AG Id for the given registered HW Id
func FindGatewayId(networkId string, hwId string) (string, error) {
	md, err := getMagmadClient()
	if err != nil {
		return "", err
	}
	gwId := identity.NewGateway(hwId, networkId, "")

	gwId, err = md.FindGatewayId(context.Background(), gwId)
	if err != nil {
		return "", err
	}
	return gwId.GetGateway().GetLogicalId(), nil
}

// FindGatewayRecord returns AG Record for a given registered logical ID
func FindGatewayRecord(networkId string, gatewayId string) (*mdprotos.AccessGatewayRecord, error) {
	md, err := getMagmadClient()
	if err != nil {
		return nil, err
	}
	gwId := identity.NewGateway("", networkId, gatewayId)
	return md.FindGatewayRecord(context.Background(), gwId)
}

// return the AccessGatewayRecord given hwId
func FindGatewayRecordWithHwId(hwId string) (*mdprotos.AccessGatewayRecord, error) {
	networkId, err := FindGatewayNetworkId(hwId)
	if err != nil {
		return nil, fmt.Errorf("Network ID Lookup Error for hwId %s: %s", hwId, err)
	}

	logicalId, err := FindGatewayId(networkId, hwId)
	if err != nil {
		return nil, fmt.Errorf("Logical ID Lookup Error for  hwId %s: %s", hwId, err)
	}

	gatewayRecord, err := FindGatewayRecord(networkId, logicalId)
	if err != nil {
		return nil, fmt.Errorf("GatewayRecord Lookup Error for hwId %s: %s", hwId, err)
	}
	return gatewayRecord, nil
}

// Finds and Updates the GW record, the record's HwId must be either omitted
// or must match the GW's registered HW ID, the HwId is not mutable
func UpdateGatewayRecord(networkId string, gatewayId string, record *mdprotos.AccessGatewayRecord) error {
	md, err := getMagmadClient()
	if err != nil {
		return err
	}
	req := &mdprotos.GatewayRecordRequest{
		GatewayId: identity.NewGateway("", networkId, gatewayId),
		Record:    record}
	_, err = md.UpdateGatewayRecord(context.Background(), req)
	return err
}

// FindGatewayNetworkId returns Network Id of the network, the Gatway HW ID
// is registered on
func FindGatewayNetworkId(hwId string) (string, error) {
	md, err := getMagmadClient()
	if err != nil {
		return "", err
	}
	netIdentity, err := md.FindGatewayNetworkId(
		context.Background(), identity.NewGateway(hwId, "", ""))
	return netIdentity.GetNetwork(), err
}

// RemoveGateway deletes all logical device & corresponding HW ID records &
// configs and effectively performs de-registration of the AG with the cloud
func RemoveGateway(networkId string, gatewayId string) error {
	md, err := getMagmadClient()
	if err != nil {
		return err
	}
	_, err = md.RemoveGateway(context.Background(), identity.NewGateway("", networkId, gatewayId))
	if err != nil {
		return err
	}

	// Delete all configs associated with this gateway
	return config.DeleteConfigsByKey(networkId, gatewayId)
}
