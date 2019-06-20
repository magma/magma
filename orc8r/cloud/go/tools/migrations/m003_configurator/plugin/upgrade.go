/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations/m003_configurator/migration"
	"magma/orc8r/cloud/go/tools/migrations/m003_configurator/plugin/types"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func migrateUpgradeTiers(
	sc *squirrel.StmtCache,
	builder sqorc.StatementBuilder,
	migratedGatewayMetasByNetwork map[string]map[string]migration.MigratedGatewayMeta,
) error {
	for nid, metas := range migratedGatewayMetasByNetwork {
		if err := migrateUpgradeTiersForNetwork(sc, builder, nid, metas); err != nil {
			return errors.Wrapf(err, "failed to migrate upgrade tiers for network %s", nid)
		}
	}
	return nil
}

func migrateUpgradeTiersForNetwork(
	sc *squirrel.StmtCache,
	builder sqorc.StatementBuilder,
	networkID string,
	gatewayMetas map[string]migration.MigratedGatewayMeta,
) error {
	// load tiers from upgrade service table and migrate them
	tiers, err := migration.UnmarshalJSONPBProtosFromDatastore(sc, builder, networkID, "tierVersions", &types.TierInfo{})
	if err != nil {
		return errors.WithStack(err)
	}
	if funk.IsEmpty(tiers) {
		return nil
	}
	migratedTiers := map[string][]byte{}
	for tierID, oldTier := range tiers {
		tierInfo := oldTier.(*types.TierInfo)
		var imgList []*types.TierImagesItems0
		for _, elem := range tierInfo.GetImages() {
			imgList = append(imgList, &types.TierImagesItems0{
				Name:  elem.Name,
				Order: elem.Order},
			)
		}
		migratedTier := &types.Tier{
			ID:      tierID,
			Name:    tierInfo.Name,
			Version: tierInfo.Version,
			Images:  imgList,
		}
		marshaledTier, err := migratedTier.MarshalBinary()
		if err != nil {
			return errors.Wrap(err, "failed to marshal migrated tier")
		}
		migratedTiers[tierID] = marshaledTier
	}

	// load magmad configs from config service so we know how gateways map to
	// tiers
	rows, err := builder.Select(migration.ConfigKeyCol, migration.ConfigValCol).
		From(migration.GetLegacyTableName(networkID, migration.ConfigTable)).
		Where(squirrel.Eq{migration.ConfigTypeCol: "magmad_gateway"}).
		RunWith(sc).
		Query()
	if err != nil {
		return errors.Wrap(err, "failed to load magmad configs")
	}
	defer rows.Close()

	gwMetasByTier := map[string][]migration.MigratedGatewayMeta{}
	for rows.Next() {
		var k string
		var v []byte

		err := rows.Scan(&k, &v)
		if err != nil {
			return errors.Wrap(err, "failed to scan magmad config")
		}

		conf := &types.OldMagmadGatewayConfig{}
		err = migration.Unmarshal(v, conf)
		if err != nil {
			return errors.Wrapf(err, "could not unmarshal magmad config %s", k)
		}
		// it's possible for the config service tables to hold configs for
		// gateways which don't exist anymore
		meta, found := gatewayMetas[k]
		if found {
			gwMetasByTier[conf.Tier] = append(gwMetasByTier[conf.Tier], meta)
		}
	}

	graphIDsToChangeByTierGraphID := map[string][]string{}
	tierInsertBuilder := builder.Insert(migration.EntityTable).
		Columns(migration.EntPkCol, migration.EntNidCol, migration.EntTypeCol, migration.EntKeyCol, migration.EntNameCol, migration.EntConfCol, migration.EntGidCol).
		RunWith(sc)
	assocInsertBuilder := builder.Insert(migration.EntityAssocTable).
		Columns(migration.AFrCol, migration.AToCol).
		RunWith(sc)
	for tierID, tierMsg := range migratedTiers {
		pk, gid := uuid.New().String(), uuid.New().String()
		name := tiers[tierID].(*types.TierInfo).Name

		tierInsertBuilder = tierInsertBuilder.Values(pk, networkID, "upgrade_tier", tierID, name, tierMsg, gid)
		for _, gwMeta := range gwMetasByTier[tierID] {
			graphIDsToChangeByTierGraphID[gid] = append(graphIDsToChangeByTierGraphID[gid], gwMeta.GraphID)
			assocInsertBuilder = assocInsertBuilder.Values(pk, gwMeta.Pk)
		}
	}

	_, err = tierInsertBuilder.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to insert new tiers")
	}

	_, err = assocInsertBuilder.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create tier assocs")
	}

	for newGid, oldGids := range graphIDsToChangeByTierGraphID {
		_, err := builder.Update(migration.EntityTable).
			Set(migration.EntGidCol, newGid).
			Where(squirrel.Eq{migration.EntGidCol: oldGids}).
			RunWith(sc).
			Exec()
		if err != nil {
			return errors.Wrap(err, "failed to update graph IDs after tier insertion")
		}
	}

	return nil
}

func migrateReleaseChannels(sc *squirrel.StmtCache, builder sqorc.StatementBuilder) error {
	rows, err := builder.Select(migration.DatastoreKeyCol, migration.DatastoreValCol).
		From("releases").
		RunWith(sc).
		Query()
	if err != nil {
		return errors.Wrap(err, "failed to load release channels")
	}
	defer rows.Close()

	channels := map[string][]byte{}
	for rows.Next() {
		var k string
		var v []byte
		err := rows.Scan(&k, &v)
		if err != nil {
			return errors.Wrap(err, "failed to scan channel row")
		}

		legacyChannel := &types.LegacyReleaseChannel{}
		err = migration.Unmarshal(v, legacyChannel)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal channel %s", k)
		}
		newChannel := &types.ReleaseChannel{
			SupportedVersions: legacyChannel.SupportedVersions,
		}
		marshaledChannel, err := newChannel.MarshalBinary()
		if err != nil {
			return errors.Wrapf(err, "failed to marshal channel %s", k)
		}
		channels[k] = marshaledChannel
	}

	channelInsertBuilder := builder.Insert(migration.EntityTable).
		Columns(migration.EntPkCol, migration.EntNidCol, migration.EntTypeCol, migration.EntKeyCol, migration.EntNameCol, migration.EntConfCol, migration.EntGidCol).
		RunWith(sc)
	for channelName, channel := range channels {
		pk, gid := uuid.New().String(), uuid.New().String()
		channelInsertBuilder = channelInsertBuilder.Values(pk, migration.InternalNetworkID, "upgrade_release_channel", channelName, channelName, channel, gid)
	}
	_, err = channelInsertBuilder.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to insert migrated release channels")
	}

	return nil
}
