/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
PolicyDB servicer provides the gRPC interface for the REST and
services to interact with the PolicyRule and Base Names data.

The servicer require a backing Datastore (which is typically Postgres)
for storing and retrieving the data.
*/
package servicers

import (
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/datastore"
	orcprotos "magma/orc8r/cloud/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddBaseName adds new Charging Rule Base Name Record (list of corresponding rule names)
// or Updates an existing Record corresponding to the given network & base name
// Returns the the existing base name if present
func (srv *PolicyDBServer) AddBaseName(
	ctx context.Context, bnData *protos.ChargingRuleBaseNameRequest,
) (*protos.ChargingRuleNameSet, error) {

	lookup := bnData.GetLookup()
	table := datastore.GetTableName(lookup.GetNetworkID().GetId(), CHARGING_RULE_BASE_NAME_TABLE)

	// Old record if existed
	res := new(protos.ChargingRuleNameSet)
	marshaled, _, err := srv.store.Get(table, lookup.GetName())
	if err == nil {
		if err = proto.Unmarshal(marshaled, res); err != nil {
			glog.Errorf("Error Unmarshalling Base Name %s -> %v: %s", lookup.GetName(), *res, err)
			return res, status.Errorf(codes.Aborted, "Unmarshalling error")
		}
	}

	// Dedup Rule Names before storing them
	record := bnData.GetRecord()
	ruleNamesMap := map[string]struct{}{}
	newRuleNames := record.GetRuleNames()
	for i, ln := 0, len(newRuleNames); i < ln; {
		ruleName := newRuleNames[i]
		if _, exist := ruleNamesMap[ruleName]; exist {
			glog.Errorf("Duplicate Rule Name '%s' in Base Name Set '%s[%d]'", ruleName, lookup.GetName(), i)
			newRuleNames = append(newRuleNames[:i], newRuleNames[i+1:]...)
			ln--
		} else {
			ruleNamesMap[ruleName] = struct{}{}
			i++
		}
		record.RuleNames = newRuleNames
	}

	// Marshal the protobuf and store the byte stream in the Datastore
	marshaled, err = proto.Marshal(record)
	if err != nil {
		glog.Errorf("Error serializing Base Name %s: %s", lookup.GetName(), err)
		return res, status.Errorf(codes.Aborted, "Marshalling error")
	}

	// Add the Base Name Record to the Datastore
	if err = srv.store.Put(table, lookup.GetName(), marshaled); err != nil {
		glog.Errorf("Error persisting Base Name %s: %s", lookup.GetName(), err)
		return res, status.Errorf(codes.Aborted, "Error adding Base Name: %s", err)
	}
	return res, nil
}

// DeleteBaseName deletes an existing Charging Rule Base Name and its Record
func (srv *PolicyDBServer) DeleteBaseName(
	ctx context.Context, lookup *protos.ChargingRuleBaseNameLookup,
) (*orcprotos.Void, error) {

	table := datastore.GetTableName(lookup.GetNetworkID().GetId(), CHARGING_RULE_BASE_NAME_TABLE)
	if err := srv.store.Delete(table, lookup.GetName()); err != nil {
		glog.Errorf("Error deleting rule %s: %s", lookup.GetName(), err)
		return &orcprotos.Void{}, status.Errorf(codes.Aborted, "Deletion error!")
	}
	return &orcprotos.Void{}, nil
}

// GetBaseName returns the ChargingRuleBaseNameRecord given the base name and the network.
func (srv *PolicyDBServer) GetBaseName(
	ctx context.Context,
	lookup *protos.ChargingRuleBaseNameLookup,
) (*protos.ChargingRuleNameSet, error) {

	table := datastore.GetTableName(lookup.GetNetworkID().GetId(), CHARGING_RULE_BASE_NAME_TABLE)
	res := new(protos.ChargingRuleNameSet)
	marshaled, _, err := srv.store.Get(table, lookup.GetName())
	if err != nil {
		glog.Errorf("Error fetching Base Name %s: %s", lookup.GetName(), err)
		return res, status.Errorf(codes.Aborted, "Error fetching rule")
	}
	if err = proto.Unmarshal(marshaled, res); err != nil {
		glog.Errorf("Error parsing Base Name %s: %s", lookup.GetName(), err)
		return res, status.Errorf(codes.Aborted, "Unmarshalling error")
	}
	return res, nil
}

// ListBaseNames returns a list of all Base Names for the given Network, the Rule Name Lists
// associated with each base name can be retrieved using separate GetBaseName() call
func (srv *PolicyDBServer) ListBaseNames(
	ctx context.Context,
	networkId *orcprotos.NetworkID,
) (*protos.ChargingRuleNameSet, error) {
	table := datastore.GetTableName(networkId.GetId(), CHARGING_RULE_BASE_NAME_TABLE)
	keys, err := srv.store.ListKeys(table)
	if err != nil {
		glog.Error(err)
		return &protos.ChargingRuleNameSet{}, status.Errorf(
			codes.Aborted, "Error listing Base Names %s, for network: %s", err, networkId.GetId())
	}
	return &protos.ChargingRuleNameSet{RuleNames: keys}, nil
}
