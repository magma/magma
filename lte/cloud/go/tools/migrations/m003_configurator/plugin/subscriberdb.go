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
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const subscriberTable = "subscriberdb"

func migrateSubscriberdb(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkIDs []string) error {
	for _, nid := range networkIDs {
		err := migrateSubscribersForNetwork(sc, builder, nid)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func migrateSubscribersForNetwork(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string) error {
	oldSubscribers, err := migration.UnmarshalProtoMessagesFromDatastore(sc, builder, networkID, subscriberTable, &types.SubscriberData{})
	if err != nil {
		return errors.WithStack(err)
	}
	if funk.IsEmpty(oldSubscribers) {
		return nil
	}

	// we need to shard out the existing subscriber data model to configurator
	// and state services
	statesExist := false
	stateInsertBuilder := builder.Insert(migration.StateServiceTable).
		Columns(migration.BlobNidCol, migration.BlobTypeCol, migration.BlobKeyCol, migration.BlobValCol).
		RunWith(sc)
	subInsertBuilder := builder.Insert(migration.EntityTable).
		Columns(migration.EntPkCol, migration.EntNidCol, migration.EntTypeCol, migration.EntKeyCol, migration.EntConfCol, migration.EntGidCol).
		RunWith(sc)
	for sid, oldMsg := range oldSubscribers {
		oldSub := oldMsg.(*types.SubscriberData)
		newSub, err := migrateSubscriber(oldSub)
		if err != nil {
			return errors.Wrapf(err, "failed to migrate subscriber (%s, %s)", networkID, sid)
		}
		pk, gid := uuid.New().String(), uuid.New().String()
		subInsertBuilder = subInsertBuilder.Values(pk, networkID, "subscriber", sid, newSub, gid)

		if oldSub.State != nil {
			newState := &types.SubscriberState{
				LteAuthNextSeq:          oldSub.State.LteAuthNextSeq,
				TgppAaaServerName:       oldSub.State.TgppAaaServerName,
				TgppAaaServerRegistered: oldSub.State.TgppAaaServerRegistered,
			}
			marshaledState, err := newState.MarshalBinary()
			if err != nil {
				return errors.Wrapf(err, "could not marshal state for subscriber (%s, %s)", networkID, sid)
			}

			stateInsertBuilder = stateInsertBuilder.Values(networkID, "subscriber", sid, marshaledState)
			statesExist = true
		}
	}

	_, err = subInsertBuilder.Exec()
	if err != nil {
		return errors.Wrapf(err, "failed to insert migrated subscribers for network %s", networkID)
	}

	if statesExist {
		_, err = stateInsertBuilder.Exec()
		if err != nil {
			return errors.Wrapf(err, "failed to insert migrated subscriber states for network %s", networkID)
		}
	}

	return nil
}

func migrateSubscriber(oldSub *types.SubscriberData) ([]byte, error) {
	newSub := &types.Subscriber{}
	migration.FillIn(oldSub, newSub)
	newSub.ID = types.SubscriberID("IMSI" + oldSub.Sid.Id)
	if newSub.Lte != nil && oldSub.Lte != nil {
		t, ok := types.LegacyLTESubscription_LTESubscriptionState_name[int32(oldSub.Lte.State)]
		if ok {
			newSub.Lte.State = t
		} else {
			newSub.Lte.State = "INACTIVE"
		}
		t, ok = types.LegacyLTESubscription_LTEAuthAlgo_name[int32(oldSub.Lte.AuthAlgo)]
		if ok {
			newSub.Lte.AuthAlgo = t
		} else {
			newSub.Lte.AuthAlgo = "MILENAGE"
		}
		if len(oldSub.Lte.AuthKey) > 0 {
			newSub.Lte.AuthKey = (*strfmt.Base64)(&oldSub.Lte.AuthKey)
		} else {
			newSub.Lte.AuthKey = nil
		}
		if len(oldSub.Lte.AuthOpc) > 0 {
			newSub.Lte.AuthOpc = (*strfmt.Base64)(&oldSub.Lte.AuthOpc)
		} else {
			newSub.Lte.AuthOpc = nil
		}
	}
	return newSub.MarshalBinary()
}
