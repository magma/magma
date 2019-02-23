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

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/protos"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//
// MagmadConfigurator service implementation
//

// RegisterGateway adds new records for AG with HW ID specified by
// record.hw_id into all corresponding tables using id.gatewayId as the the
// newly registered gateway's logical ID, initializes the newly created AG's
// configs to defaults
// Will return error if a device with the given HW ID or requested by
// id.gatewayId logical id is already registered
// if id.gatewayId is empty - will register the network with a new unique ID
func (md MagmadConfigurator) RegisterGateway(ctx context.Context, req *magmadprotos.GatewayRecordRequest) (*protos.Identity, error) {
	res := &protos.Identity{}
	if err := validateGatewayRecordRequest(req); err != nil {
		return res, err
	}
	if req.Record.HwId == nil || len(req.Record.HwId.Id) == 0 {
		return res, status.Errorf(codes.InvalidArgument, "Invalid Gateway HwId")
	}

	gw, err := getGatewayWithNetworkAndHwId(req.GatewayId)
	if err != nil {
		return res, err
	}
	if len(gw.LogicalId) == 0 { // Caller did not request logicalID, use HW ID
		gw.LogicalId = gw.HardwareId
	}
	res.SetGateway(gw)

	if err := md.validateNewGateway(gw); err != nil {
		return res, err
	}
	req.Record.HwId.Id = gw.HardwareId

	// Persistence
	err = md.persistNewGatewayRollbackOnError(gw.NetworkId, gw.LogicalId, req.Record)
	if err != nil {
		glog.Error(err)
		return res, status.Errorf(codes.Internal, "Error registering access gateway: %s", err.Error())
	}
	return res, nil
}

// Lists all registered logical device IDs
func (md MagmadConfigurator) ListGateways(ctx context.Context, netId *protos.Identity) (*protos.Identity_List, error) {
	res := &protos.Identity_List{}
	nid := netId.GetNetwork()
	if len(nid) == 0 {
		return res, status.Errorf(codes.InvalidArgument, "Invalid Network ID")
	}
	ids, err := md.Store.ListKeys(datastore.GetTableName(nid, AgRecordTableName))
	if err != nil {
		return res, status.Errorf(codes.Internal, err.Error())
	}
	res.List = make([]*protos.Identity, len(ids))
	for i, id := range ids {
		res.List[i] = new(protos.Identity).SetGateway(&protos.Identity_Gateway{NetworkId: nid, LogicalId: id})
	}
	return res, nil
}

// FindGatewayId augments the given GW Identity including HW ID & Network ID
// with logical AG Id corresponding to given registered HW Id on the network
func (md MagmadConfigurator) FindGatewayId(ctx context.Context, gwId *protos.Identity) (*protos.Identity, error) {
	gw, err := getGatewayWithNetworkAndHwId(gwId)
	if err != nil {
		return gwId, err
	}
	gwId.SetGateway(gw)
	lId, _, err := md.Store.Get(
		datastore.GetTableName(gw.GetNetworkId(), HwIdTableName),
		gw.GetHardwareId())
	if err != nil {
		return gwId, status.Errorf(codes.NotFound, err.Error())
	}
	gw.LogicalId = string(lId)

	return gwId, nil
}

// FindGatewayRecord returns AG Record for a given registered logical ID
func (md MagmadConfigurator) FindGatewayRecord(ctx context.Context, gwId *protos.Identity) (*magmadprotos.AccessGatewayRecord, error) {
	res := &magmadprotos.AccessGatewayRecord{}
	gw, err := getGatewayWithLogicalId(gwId)
	if err != nil {
		return res, err
	}
	marshaledRecord, _, err := md.Store.Get(
		datastore.GetTableName(gw.GetNetworkId(), AgRecordTableName),
		gw.GetLogicalId())
	if err != nil {
		return res, status.Errorf(codes.NotFound, err.Error())
	}
	return res, protos.Unmarshal(marshaledRecord, res)
}

// UpdateGatewayRecord updates AG Record for a given registered gateway
func (md MagmadConfigurator) UpdateGatewayRecord(ctx context.Context, req *magmadprotos.GatewayRecordRequest) (*protos.Void, error) {
	res := &protos.Void{}
	err := validateGatewayRecordRequest(req)
	if err != nil {
		return res, err
	}
	gw, err := getGatewayWithLogicalId(req.GatewayId)
	if err != nil {
		return res, err
	}
	tableName := datastore.GetTableName(gw.GetNetworkId(), AgRecordTableName)
	marshaledRecord, _, err := md.Store.Get(tableName, gw.GetLogicalId())

	if err != nil {
		return res, status.Errorf(codes.NotFound, err.Error())
	}
	agRecord := &magmadprotos.AccessGatewayRecord{}
	err = protos.Unmarshal(marshaledRecord, agRecord)
	if err != nil {
		return res, status.Errorf(codes.Internal, "Unmarshal error: %s", err)
	}
	if req.Record.HwId != nil &&
		(agRecord.HwId == nil || req.Record.HwId.Id != agRecord.HwId.Id) {
		return res, status.Errorf(codes.InvalidArgument, "Cannot modify GW HwId (From %#v to %#v)", agRecord.HwId, req.Record.HwId)
	}
	// Update record, but preserve original HW ID
	tmp := agRecord.HwId
	*agRecord, agRecord.HwId = *req.Record, tmp

	marshaledRecord, err = protos.MarshalIntern(agRecord)
	if err != nil {
		return res, status.Errorf(codes.Internal, "MarshalGatewayRecord Error: %s", err)
	}
	err = md.Store.Put(tableName, gw.GetLogicalId(), marshaledRecord)
	if err != nil {
		return res, status.Errorf(codes.Internal, "AgRecords Table Put error: %s for ID: %s", err, gw.GetLogicalId())
	}
	return res, nil
}

// FindGatewayNetworkId returns Network Id of the network, the Gatway HW ID
// is registered on
func (md MagmadConfigurator) FindGatewayNetworkId(ctx context.Context, id *protos.Identity) (*protos.Identity, error) {
	gw, err := getGatewayWithHardwareId(id)
	if err != nil {
		return &protos.Identity{}, err
	}
	networkId, _, err := md.Store.Get(GatewaysTableName, gw.HardwareId)
	if err != nil {
		return &protos.Identity{}, status.Errorf(codes.NotFound, err.Error())
	}
	return identity.NewNetwork(string(networkId)), nil
}

// RemoveGateway deletes all logical device & corresponding HW ID records & configs
// and effectively performs de-registration of the AG with the cloud
func (md MagmadConfigurator) RemoveGateway(ctx context.Context, id *protos.Identity) (*protos.Void, error) {
	res := &protos.Void{}
	gw, err := getGatewayWithLogicalId(id)
	if err != nil {
		return res, err
	}
	networkId := gw.GetNetworkId()
	logicalId := gw.GetLogicalId()
	gwRecord, err := md.FindGatewayRecord(ctx, id)
	if err != nil {
		return res, err
	}
	hwId := gwRecord.GetHwId().GetId()

	var allOperationErrors []string

	// Delete if exists in case this is a retry of a partially failed removal
	if err := md.deleteIfExists(datastore.GetTableName(networkId, HwIdTableName), hwId); err != nil {
		msg := fmt.Sprintf("Failed to delete logical ID mapping. Error: %s", err)
		glog.Error(msg)
		allOperationErrors = append(allOperationErrors, msg)
	}
	if err := md.deleteIfExists(GatewaysTableName, hwId); err != nil {
		msg := fmt.Sprintf("Failed to delete hardware ID network mapping. Error: %s", err)
		glog.Error(msg)
		allOperationErrors = append(allOperationErrors, msg)
	}
	if err := md.deleteIfExists(datastore.GetTableName(networkId, GatewaysStatusTableName), logicalId); err != nil {
		msg := fmt.Sprintf("Failed to clean up gateway checkin statuses. Error: %s", err)
		glog.Error(msg)
		allOperationErrors = append(allOperationErrors, msg)
	}

	// Only remove the gateway record if all cleanup steps are successful
	// This way the gateway ID still shows up on a LIST request
	if len(allOperationErrors) == 0 {
		err := md.Store.Delete(datastore.GetTableName(networkId, AgRecordTableName), logicalId)
		if err != nil {
			msg := fmt.Sprintf("Failed to delete gateway record. Error: %s", err)
			glog.Error(msg)
			return res, status.Errorf(codes.Internal, msg)
		}
		return res, nil
	} else {
		allErrorsJoined := strings.Join(allOperationErrors, "\n\t")
		return res, status.Errorf(
			codes.Internal,
			"Encountered the following errors while removing the gateway:\n"+
				"\t%s\nPlease address the issues and retry the operation.",
			allErrorsJoined)
	}
}

func (md MagmadConfigurator) validateNewGateway(gw *protos.Identity_Gateway) error {
	if strings.ToLower(gw.LogicalId) == strings.ToLower(gw.NetworkId) {
		return status.Errorf(codes.FailedPrecondition, "Gateway ID must be different from network ID")
	}
	if err := md.validateNetworkExists(gw.NetworkId); err != nil {
		return status.Errorf(codes.FailedPrecondition, err.Error())
	}
	err := md.validateHwIdNotRegisteredInNetwork(gw.NetworkId, gw.HardwareId)
	if err != nil {
		return status.Errorf(codes.AlreadyExists, err.Error())
	}
	if err := md.validateUniqueLogicalId(gw.NetworkId, gw.LogicalId); err != nil {
		return status.Errorf(codes.AlreadyExists, err.Error())
	}
	if err := md.validateHwIdNotRegistered(gw.HardwareId); err != nil {
		return status.Errorf(codes.AlreadyExists, err.Error())
	}
	return nil
}

// Persist data pertaining to a newly registered gateway to all related tables.
// This function will write to tables in serial and attempt to rollback (delete)
// all previous writes on any errors encountered.
// Note that there is no strict guarantee that rollback will succeed given that
// this is run in a non-transactional context.
func (md MagmadConfigurator) persistNewGatewayRollbackOnError(networkId string, logicalId string, agRecord *magmadprotos.AccessGatewayRecord) error {
	// Serialize data
	hwId := agRecord.GetHwId().GetId()
	marshaledAgRecord, err := protos.MarshalIntern(agRecord)
	if err != nil {
		return fmt.Errorf("Error marshaling gateway record: %s", err)
	}

	// Persist new data
	err = md.Store.Put(GatewaysTableName, agRecord.HwId.Id, []byte(networkId))
	if err != nil {
		return fmt.Errorf("Gateways Table Put error: %s for HW ID: %s", err, hwId)
	}

	err = md.Store.Put(datastore.GetTableName(networkId, AgRecordTableName), logicalId, marshaledAgRecord)
	if err != nil {
		_ = md.Store.Delete(GatewaysTableName, agRecord.HwId.Id)
		return fmt.Errorf("AgRecords Table Put error: %s for ID: %s", err, logicalId)
	}

	err = md.Store.Put(datastore.GetTableName(networkId, HwIdTableName), agRecord.HwId.Id, []byte(logicalId))
	if err != nil {
		_ = md.Store.Delete(GatewaysTableName, agRecord.HwId.Id)
		_ = md.Store.Delete(datastore.GetTableName(networkId, AgRecordTableName), logicalId)
		return fmt.Errorf("HwId Table Put error: %s for HW ID: %s", err, agRecord.HwId.Id)
	}
	return nil
}

// *GatewayRecordRequest validator
func validateGatewayRecordRequest(req *magmadprotos.GatewayRecordRequest) error {
	if req == nil {
		return status.Errorf(codes.InvalidArgument, "Nil GatewayRecordRequest")
	}
	if req.Record == nil {
		return status.Errorf(codes.InvalidArgument, "Nil Record")
	}
	return nil
}
