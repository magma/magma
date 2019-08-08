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
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func migrateEnodebs(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, migratedGatewayMetasByNetwork map[string]map[string]migration.MigratedGatewayMeta) error {
	for nid, metas := range migratedGatewayMetasByNetwork {
		if err := migrateEnodebsForNetwork(sc, builder, nid, metas); err != nil {
			return errors.Wrapf(err, "failed to migrate enodebs for network %s", nid)
		}
	}
	return nil
}

func migrateEnodebsForNetwork(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string, gatewayMetas map[string]migration.MigratedGatewayMeta) error {
	// load enodeb and cellular gateway configs from config service
	rows, err := builder.Select(migration.ConfigTypeCol, migration.ConfigKeyCol, migration.ConfigValCol).
		From(migration.GetLegacyTableName(networkID, migration.ConfigTable)).
		Where(squirrel.Eq{migration.ConfigTypeCol: []string{"cellular_enodeb", "cellular_gateway"}}).
		RunWith(sc).
		Query()
	if err != nil {
		return errors.Wrap(err, "failed to load enodeb configs")
	}
	defer rows.Close()

	migratedEnodebs := map[string][]byte{}
	gwIDByEnodebID := map[string]string{}
	for rows.Next() {
		var t, k string
		var v []byte
		err := rows.Scan(&t, &k, &v)
		if err != nil {
			return errors.Wrap(err, "failed to scan datastore row")
		}

		switch t {
		case "cellular_enodeb":
			oldEnodeb := &types.CellularEnodebConfig{}
			err = migration.Unmarshal(v, oldEnodeb)
			if err != nil {
				return errors.Wrapf(err, "could not unmarshal enodeb config %s", k)
			}

			newEnodeb := &types.NetworkEnodebConfigs{}
			migration.FillIn(oldEnodeb, newEnodeb)
			newV, err := newEnodeb.MarshalBinary()
			if err != nil {
				return errors.Wrapf(err, "could not marshal new enodeb config %s", k)
			}
			migratedEnodebs[k] = newV
		case "cellular_gateway":
			cfg := &types.CellularGatewayConfig{}
			err = migration.Unmarshal(v, cfg)
			if err != nil {
				return errors.Wrapf(err, "could not unmarshal cellular gw config %s for enodebs", k)
			}
			for _, enodebID := range cfg.AttachedEnodebSerials {
				gwIDByEnodebID[enodebID] = k
			}
		default:
			glog.Errorf("unexpected config type %s", t)
		}

	}
	if funk.IsEmpty(migratedEnodebs) {
		return nil
	}

	enodebPKByKey := map[string]string{}
	enodebInsertBuilder := builder.Insert(migration.EntityTable).
		Columns(migration.EntPkCol, migration.EntNidCol, migration.EntTypeCol, migration.EntKeyCol, migration.EntConfCol, migration.EntGidCol).
		RunWith(sc)
	for k, v := range migratedEnodebs {
		// enodebs need to be created in the same graph as their associated
		// cellular gateway
		// init gid here so that if the enodeb isn't associated to any gateways
		// we still give it a graph ID
		pk, gid := uuid.New().String(), uuid.New().String()
		gwKey, found := gwIDByEnodebID[k]
		if found {
			gwMeta, metaFound := gatewayMetas[gwKey]
			if !metaFound {
				glog.Errorf("enodeb %s is associated to unmigrated gateway %s", k, gwKey)
			} else {
				gid = gwMeta.GraphID
			}
		} else {
			glog.Infof("enodeb %s isn't associated to any gateways", k)
		}

		enodebPKByKey[k] = pk
		enodebInsertBuilder = enodebInsertBuilder.Values(pk, networkID, "cellular_enodeb", k, v, gid)
	}
	_, err = enodebInsertBuilder.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to insert migrated enodebs")
	}

	// we need to grab the PKs of cellular gateways, not magmad gateways
	// load these from configurator entity table
	pkRows, err := builder.Select(migration.EntPkCol, migration.EntKeyCol).
		From(migration.EntityTable).
		Where(squirrel.Eq{migration.EntTypeCol: "cellular_gateway"}).
		RunWith(sc).
		Query()
	if err != nil {
		return errors.Wrap(err, "failed to load cellular gateway PKs")
	}
	defer pkRows.Close()
	cellularPKsByKey := map[string]string{}
	for pkRows.Next() {
		var pk, k string
		err = pkRows.Scan(&pk, &k)
		if err != nil {
			return errors.Wrap(err, "failed to scan row for cellular gateway PK")
		}
		cellularPKsByKey[k] = pk
	}

	if funk.IsEmpty(gwIDByEnodebID) {
		return nil
	}
	enodebAssocBuilder := builder.Insert(migration.EntityAssocTable).
		Columns(migration.AFrCol, migration.AToCol).
		RunWith(sc)
	for enbKey, gwKey := range gwIDByEnodebID {
		enbPK, enbPKFound := enodebPKByKey[enbKey]
		gwPk, gwPKFound := cellularPKsByKey[gwKey]
		if !enbPKFound || !gwPKFound {
			glog.Errorf("did not find PK for cellular gateway %s", gwKey)
			continue
		}
		enodebAssocBuilder = enodebAssocBuilder.Values(gwPk, enbPK)
	}
	_, err = enodebAssocBuilder.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to insert gw <-> enodeb assocs")
	}

	return nil
}
