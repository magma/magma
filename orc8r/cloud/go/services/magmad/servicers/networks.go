/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"
	"strings"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/protos"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//==============================================================================
// Network APIs
//==============================================================================

func (md MagmadConfigurator) GetNetwork(ctx context.Context, nid *protos.Identity) (*magmad_protos.MagmadNetworkRecord, error) {
	networkId := nid.GetNetwork()
	if nid == nil || len(networkId) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request - empty network ID")
	}

	marshaledRecord, _, err := md.Store.Get(NetworksTableName, networkId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error retrieving network record: %s", err)
	}
	ret := &magmad_protos.MagmadNetworkRecord{}
	err = protos.Unmarshal(marshaledRecord, ret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error unmarshaling network record: %s", err)
	}
	return ret, nil
}

func (md MagmadConfigurator) UpdateNetwork(ctx context.Context, req *magmad_protos.NetworkRecordRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if req == nil || req.Record == nil || len(req.Id) == 0 {
		return ret, status.Error(codes.InvalidArgument, "Invalid request")
	}

	marshaledRecord, err := protos.MarshalIntern(req.Record)
	if err != nil {
		return ret, status.Errorf(codes.Internal, "Error marshaling network record: %s", err)
	}
	err = md.Store.Put(NetworksTableName, req.Id, marshaledRecord)
	if err != nil {
		return ret, status.Errorf(codes.Internal, "Error persisting network record: %s", err)
	}
	return ret, nil
}

func (md MagmadConfigurator) RegisterNetwork(ctx context.Context, req *magmad_protos.NetworkRecordRequest) (*protos.Identity, error) {
	if req == nil || req.Record == nil {
		return &protos.Identity{}, status.Errorf(codes.InvalidArgument, "Invalid Request")
	}

	networkId, err := md.getNetworkId(req.Id)
	if err != nil {
		return &protos.Identity{}, err
	}
	res := identity.NewNetwork(networkId)

	marshaledNetworkRecord, err := protos.MarshalIntern(req.Record)
	if err != nil {
		return res, status.Errorf(codes.Internal, "Error marshaling network record: %s", err)
	}
	err = md.Store.Put(NetworksTableName, networkId, marshaledNetworkRecord)
	if err != nil {
		return res, status.Errorf(codes.Internal, err.Error())
	}
	return res, nil
}

func (md MagmadConfigurator) ListNetworks(context.Context, *protos.Void) (*protos.Identity_List, error) {
	res := &protos.Identity_List{}
	nids, err := md.Store.ListKeys(NetworksTableName)
	if err != nil {
		return res, status.Errorf(codes.Internal, err.Error())
	}
	res.List = make([]*protos.Identity, len(nids))
	for i, nid := range nids {
		res.List[i] = identity.NewNetwork(nid)
	}
	return res, nil
}

func (md MagmadConfigurator) RemoveNetwork(ctx context.Context, nid *protos.Identity) (*protos.Void, error) {
	res := &protos.Void{}
	networkId := nid.GetNetwork()
	if nid == nil || len(networkId) == 0 {
		return res, status.Errorf(codes.InvalidArgument, "Invalid Request")
	}

	exist, err := md.Store.DoesKeyExist(NetworksTableName, networkId)
	if err != nil {
		return res, status.Errorf(codes.Internal, "Failed to determine if key %s exists", networkId)
	}
	if !exist {
		return res, status.Errorf(codes.NotFound, "Network ID %s is not found", networkId)
	}

	tablesToDrop := getTablesToDropForNetworkDeletion(networkId)

	for _, tableName := range tablesToDrop {
		keys, err := md.Store.ListKeys(tableName)
		if err == nil && keys != nil && len(keys) != 0 {
			return res, status.Errorf(codes.FailedPrecondition,
				"Table %s of network %s is not empty. "+
					"All Gateways and Subscribers must be removed first.",
				tableName, networkId)
		}
	}

	for _, tableName := range tablesToDrop {
		err = md.Store.DeleteTable(tableName)
		if err != nil {
			return res, status.Errorf(codes.Internal, "Failed to delete %s table. Error: %s", tableName, err)
		}
	}

	err = md.Store.Delete(NetworksTableName, networkId)
	if err != nil {
		return res, status.Errorf(codes.Internal, err.Error())
	}
	return res, nil
}

func (md MagmadConfigurator) ForceRemoveNetwork(ctx context.Context, nid *protos.Identity) (*protos.Void, error) {
	res := &protos.Void{}
	networkId := nid.GetNetwork()
	if nid == nil || len(networkId) == 0 {
		return res, status.Errorf(codes.InvalidArgument, "Invalid Request")
	}
	exist, err := md.Store.DoesKeyExist(NetworksTableName, networkId)
	if err != nil {
		return res, status.Errorf(codes.Internal, "Failed to determine if key %s exists", networkId)
	}
	if !exist {
		return res, status.Errorf(codes.NotFound, "Network ID %s is not found", networkId)
	}

	// Exit early if we can't query for hardware IDs to force user retry
	allHwIds, err := md.getAllHWIdsInNetwork(networkId)
	if err != nil {
		return res, status.Errorf(codes.Internal,
			"Failed to query gateway hardware IDs in network, "+
				"exiting before performing any deletions. Error: %s",
			err)
	}

	// Accumulate all errors that don't require an early exit for retry
	var allOperationErrors []string

	// Clear hardware IDs
	for _, hwId := range allHwIds {
		// Double check that the HW ID exists
		hwIdExists, err := md.Store.DoesKeyExist(GatewaysTableName, hwId)
		if err != nil {
			msg := fmt.Sprintf("Error while checking if hardware ID %s exists: %s", hwId, err)
			glog.Error(msg)
			allOperationErrors = append(allOperationErrors, msg)
			continue
		}

		// Only delete if it does exist (don't pollute error collection)
		if hwIdExists {
			err = md.Store.Delete(GatewaysTableName, hwId)
			if err != nil {
				msg := fmt.Sprintf("Error while deleting hardware ID %s: %s", hwId, err)
				glog.Error(msg)
				allOperationErrors = append(allOperationErrors, msg)
			}
		}
	}

	// Exit early if any failures happened above so that this can be retry-able.
	// Otherwise, the hardware IDs will disappear when the table is dropped next.
	if len(allOperationErrors) > 0 {
		return res, status.Errorf(codes.Internal,
			"Encountered the following errors while clearing hardware IDs:\n"+
				"\t%s\nPlease retry the operation.",
			strings.Join(allOperationErrors, "\n\t"))
	}

	// Drop all tables
	tablesToDrop := getTablesToDropForNetworkDeletion(networkId)
	for _, tableName := range tablesToDrop {
		err = md.Store.DeleteTable(tableName)
		if err != nil {
			msg := fmt.Sprintf("Error while deleting table %s: %s", tableName, err)
			glog.Error(msg)
			allOperationErrors = append(allOperationErrors, msg)
		}
	}

	// Only clear the network if no errors were encountered - this way,
	// users can retry the operation from the frontend because the network ID
	// will still show up in a LIST request
	if len(allOperationErrors) == 0 {
		err = md.Store.Delete(NetworksTableName, networkId)
		if err != nil {
			return res, status.Errorf(codes.Internal, "Error while deleting network record %s: %s", networkId, err)
		}
		return res, nil
	} else {
		return res, status.Errorf(codes.Internal,
			"Encountered the following errors while deleting the network:\n"+
				"\t%s\nPlease address the issues and retry the operation.",
			strings.Join(allOperationErrors, "\n\t"),
		)
	}
}

func (md MagmadConfigurator) getNetworkId(requestedId string) (string, error) {
	if len(requestedId) == 0 {
		return "", status.Errorf(codes.Internal, "No network ID was provided")
	}

	requestedLower := strings.ToLower(requestedId)
	exists, err := md.Store.DoesKeyExist(NetworksTableName, requestedLower)
	if err != nil {
		return "", status.Errorf(codes.Internal, "Error checking if requested network ID already exists: %s", err)
	}

	if exists {
		return "", status.Errorf(codes.AlreadyExists, "Network ID %s already exists", requestedLower)
	}
	return requestedLower, nil
}
