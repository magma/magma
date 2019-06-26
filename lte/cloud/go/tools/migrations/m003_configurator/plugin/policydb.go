/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/lte/cloud/go/tools/migrations/m003_configurator/plugin/types"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations/m003_configurator/migration"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func migratePolicydb(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkIDs []string) error {
	for _, nid := range networkIDs {
		// first migrate the rules - we need the PKs so we can associate base
		// names with the rules later
		rulePKsByID, err := migratePolicydbRulesForNetwork(sc, builder, nid)
		if err != nil {
			return errors.WithStack(err)
		}

		err = migratePolicydbBaseNamesForNetwork(sc, builder, nid, rulePKsByID)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// returns pks by rule id
func migratePolicydbRulesForNetwork(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string) (map[string]string, error) {
	oldRules, err := migration.UnmarshalProtoMessagesFromDatastore(sc, builder, networkID, policyTable, &types.LegacyPolicyRule{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load rule names for network %s", networkID)
	}
	if funk.IsEmpty(oldRules) {
		return map[string]string{}, nil
	}

	ret := map[string]string{}
	insertBuilder := builder.Insert(migration.EntityTable).
		Columns(migration.EntPkCol, migration.EntNidCol, migration.EntTypeCol, migration.EntKeyCol, migration.EntConfCol, migration.EntGidCol).
		RunWith(sc)
	for ruleID, oldRule := range oldRules {
		newRule := migratePolicyRule(oldRule.(*types.LegacyPolicyRule))
		marshaledRule, err := newRule.MarshalBinary()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal rule (%s, %s)", networkID, ruleID)
		}

		pk, gid := uuid.New().String(), uuid.New().String()
		insertBuilder = insertBuilder.Values(pk, networkID, "policy", ruleID, marshaledRule, gid)
		ret[ruleID] = pk
	}
	_, err = insertBuilder.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert migrated policy rules")
	}

	return ret, nil
}

func migratePolicyRule(oldRule *types.LegacyPolicyRule) *types.PolicyRule {
	ret := &types.PolicyRule{}
	migration.FillIn(oldRule, ret)
	if oldRule.FlowList != nil {
		for i, flow := range oldRule.FlowList {
			ret.FlowList = append(ret.FlowList, new(types.FlowDescription))
			migration.FillIn(flow, ret.FlowList[i])
			match := flow.Match
			migration.FillIn(match, ret.FlowList[i].Match)
			protoName, ok := types.LegacyFlowMatch_IPProto_name[int32(match.IpProto)]
			if ok {
				ret.FlowList[i].Match.IPProto = &protoName
			}
			directionName, ok := types.LegacyFlowMatch_Direction_name[int32(match.Direction)]
			if ok {
				ret.FlowList[i].Match.Direction = directionName
			}
			actionName, ok := types.LegacyFlowDescription_Action_name[int32(flow.Action)]
			if ok {
				ret.FlowList[i].Action = &actionName
			}
		}
	}
	ret.Priority = &oldRule.Priority
	trackingName, ok := types.LegacyPolicyRule_TrackingType_name[int32(oldRule.TrackingType)]
	if ok {
		ret.TrackingType = trackingName
	}
	if oldRule.Redirect != nil {
		modelInfo := &types.RedirectInformation{}
		migration.FillIn(oldRule.Redirect, modelInfo)
		supportName, ok := types.LegacyRedirectInformation_Support_name[int32(oldRule.Redirect.Support)]
		if ok {
			modelInfo.Support = supportName
		}
		addrTypeName, ok := types.LegacyRedirectInformation_AddressType_name[int32(oldRule.Redirect.AddressType)]
		if ok {
			modelInfo.AddressType = addrTypeName
		}
	}
	ret.MonitoringKey = &oldRule.MonitoringKey
	ret.RatingGroup = &oldRule.RatingGroup
	return ret
}

func migratePolicydbBaseNamesForNetwork(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string, rulePKsByID map[string]string) error {
	// see comment on rule migration function
	oldBaseNames, err := migration.UnmarshalProtoMessagesFromDatastore(sc, builder, networkID, baseNameTable, &types.ChargingRuleNameSet{})
	if funk.IsEmpty(oldBaseNames) {
		return nil
	}

	// we'll need to keep track of the assocs from bn -> rule we need to create
	// also all rules attached to each bn need to have the same graph ID as
	// the created bn
	// we're making a critical assumption here that no rule is associated back
	// to more than one bn in production data, otherwise shit gets hairy
	assocsToCreate := [][2]string{} // [from PK, to PK]
	rulePKsByNewGraphID := map[string][]string{}

	// build insert statement just for all the BNs, and gather the graph info
	bnInsertBuilder := builder.Insert(migration.EntityTable).
		Columns(migration.EntPkCol, migration.EntNidCol, migration.EntTypeCol, migration.EntKeyCol, migration.EntConfCol, migration.EntGidCol).
		RunWith(sc)
	for bnID, oldBN := range oldBaseNames {
		newBN := &types.BaseNameRecord{
			RuleNames: oldBN.(*types.ChargingRuleNameSet).RuleNames,
		}
		marshaledBN, err := newBN.MarshalBinary()
		if err != nil {
			return errors.Wrapf(err, "failed to marshal new base name (%s, %s)", networkID, bnID)
		}

		bnPK, bnGid := uuid.New().String(), uuid.New().String()
		bnInsertBuilder = bnInsertBuilder.Values(bnPK, networkID, "base_name", bnID, marshaledBN, bnGid)
		for _, rn := range newBN.RuleNames {
			assocsToCreate = append(assocsToCreate, [2]string{bnPK, rulePKsByID[rn]})
			rulePKsByNewGraphID[bnGid] = append(rulePKsByNewGraphID[bnGid], rulePKsByID[rn])
		}
	}
	_, err = bnInsertBuilder.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to insert new base names")
	}

	// create insert statement for all assocs
	assocInsertBuilder := builder.Insert(migration.EntityAssocTable).
		Columns(migration.AFrCol, migration.AToCol).
		RunWith(sc)
	for _, assoc := range assocsToCreate {
		assocInsertBuilder = assocInsertBuilder.Values(assoc[0], assoc[1])
	}
	_, err = assocInsertBuilder.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to insert base name -> rule assocs")
	}

	// update graph IDs of rules
	for newGid, rulePks := range rulePKsByNewGraphID {
		_, err := builder.Update(migration.EntityTable).
			Set(migration.EntGidCol, newGid).
			Where(squirrel.Eq{migration.EntPkCol: rulePks}).
			RunWith(sc).
			Exec()
		if err != nil {
			return errors.Wrap(err, "failed to update policy rule graph IDs")
		}
	}

	return nil
}
