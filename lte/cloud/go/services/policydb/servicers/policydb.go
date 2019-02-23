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

type PolicyDBServer struct {
	store datastore.Api
}

func NewPolicyDBServer(store datastore.Api) *PolicyDBServer {
	srv := new(PolicyDBServer)
	srv.store = store
	return srv
}

func (srv *PolicyDBServer) AddRule(
	ctx context.Context,
	ruleData *protos.PolicyRuleData,
) (*orcprotos.Void, error) {
	ruleID := ruleData.Rule.Id
	table := datastore.GetTableName(ruleData.NetworkId.Id, POLICY_TABLE)

	if _, _, err := srv.store.Get(table, ruleID); err == nil {
		glog.Errorf("Rule %s already exists ", ruleID)
		return &orcprotos.Void{}, status.Errorf(codes.AlreadyExists,
			"Rule already exists")
	}

	// Marshal the protobuf and store the byte stream in the Datastore
	value, err := proto.Marshal(ruleData.Rule)
	if err != nil {
		glog.Errorf("Error serializing rule %s: %s", ruleID, err)
		return &orcprotos.Void{}, status.Errorf(codes.Aborted, "Marshalling error")
	}

	// Add the rule to the Datastore
	if err = srv.store.Put(table, ruleID, value); err != nil {
		glog.Errorf("Error persisting rule %s: %s", ruleID, err)
		return &orcprotos.Void{}, status.Errorf(
			codes.Aborted, "Error adding rule: %s", err)
	}
	return &orcprotos.Void{}, nil
}

func (srv *PolicyDBServer) DeleteRule(
	ctx context.Context,
	lookup *protos.PolicyRuleLookup,
) (*orcprotos.Void, error) {
	ruleID := lookup.RuleId
	table := datastore.GetTableName(lookup.NetworkId.Id, POLICY_TABLE)

	if err := srv.store.Delete(table, ruleID); err != nil {
		glog.Errorf("Error deleting rule %s: %s", ruleID, err)
		return &orcprotos.Void{}, status.Errorf(codes.Aborted, "Deletion error!")
	}
	return &orcprotos.Void{}, nil
}

func (srv *PolicyDBServer) UpdateRule(
	ctx context.Context,
	ruleData *protos.PolicyRuleData,
) (*orcprotos.Void, error) {
	ruleId := ruleData.Rule.Id
	table := datastore.GetTableName(ruleData.NetworkId.Id, POLICY_TABLE)

	// Marshal the protobuf and store the byte stream in the Datastore
	value, err := proto.Marshal(ruleData.Rule)
	if err != nil {
		glog.Errorf("Error serializing rule %s: %s", ruleId, err)
		return &orcprotos.Void{}, status.Errorf(codes.Aborted, "Marshalling error")
	}

	// Update the rule in the Datastore
	if err = srv.store.Put(table, ruleId, value); err != nil {
		glog.Errorf("Error persisting rule %s: %s", ruleId, err)
		return &orcprotos.Void{}, status.Errorf(codes.Aborted, "Error updating rule")
	}
	return &orcprotos.Void{}, nil
}

func (srv *PolicyDBServer) GetRule(
	ctx context.Context,
	lookup *protos.PolicyRuleLookup,
) (*protos.PolicyRule, error) {
	ruleID := lookup.RuleId
	rule := protos.PolicyRule{}
	table := datastore.GetTableName(lookup.NetworkId.Id, POLICY_TABLE)

	value, _, err := srv.store.Get(table, ruleID)
	if err != nil {
		glog.Errorf("Error fetching rule %s: %s", &rule, err)
		return &rule, status.Errorf(codes.Aborted, "Error fetching rule")
	}
	if err = proto.Unmarshal(value, &rule); err != nil {
		glog.Errorf("Error parsing rule %s: %s", &rule, err)
		return &rule, status.Errorf(codes.Aborted, "Unmarshalling error")
	}
	return &rule, nil
}

func (srv *PolicyDBServer) ListRules(
	ctx context.Context,
	networkId *orcprotos.NetworkID,
) (*protos.PolicyRuleSet, error) {
	table := datastore.GetTableName(networkId.Id, POLICY_TABLE)
	keys, err := srv.store.ListKeys(table)
	if err != nil {
		glog.Error(err)
		return &protos.PolicyRuleSet{}, status.Errorf(
			codes.Aborted, "Error listing rules %s, for network: %s",
			err, networkId.Id)
	}

	rules := make([]*protos.PolicyRule, 0, len(keys))
	for _, key := range keys {
		value, _, err := srv.store.Get(table, key)
		if err != nil {
			glog.Errorf("Error fetching rule %s: %s", key, err)
			continue
		}
		rule := protos.PolicyRule{}
		err = proto.Unmarshal(value, &rule)
		if err != nil {
			glog.Errorf("Error parsing rule %s: %s", key, err)
			continue
		}
		rules = append(rules, &rule)
	}
	return &protos.PolicyRuleSet{Rules: rules}, nil
}
