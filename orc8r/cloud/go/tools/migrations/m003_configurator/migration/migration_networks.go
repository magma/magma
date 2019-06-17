/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package migration

import (
	"database/sql"
	"encoding/json"

	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func MigrateNetworks(tx *sql.Tx, builder sqorc.StatementBuilder) ([]string, error) {
	networkRecords, err := loadAllLegacyNetworks(tx, builder)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// create networks in configurator
	nwInsertBuilder := builder.Insert(networksTable).
		Columns(nwIDCol, nwNameCol).
		RunWith(tx)
	for nid, nr := range networkRecords {
		nwInsertBuilder = nwInsertBuilder.Values(nid, nr.Name)
	}
	_, err = nwInsertBuilder.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new networks")
	}

	// migrate network configs
	migratedNwConfigs, err := getNewNetworkConfigValues(networkRecords, tx, builder)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	nwcInsertBuilder := builder.Insert(networkConfigTable).
		Columns(nwcIDCol, nwcTypeCol, nwcValCol).
		RunWith(tx)
	for nid, newConfigs := range migratedNwConfigs {
		for t, v := range newConfigs {
			nwcInsertBuilder = nwcInsertBuilder.Values(nid, t, v)
		}
	}
	_, err = nwcInsertBuilder.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert network configs")
	}

	return funk.Keys(networkRecords).([]string), nil
}

func loadAllLegacyNetworks(tx *sql.Tx, builder sqorc.StatementBuilder) (map[string]legacyNetworkRecord, error) {
	rows, err := builder.Select(datastoreKeyCol, datastoreValCol).
		From(NetworksTableName).
		RunWith(tx).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for networks")
	}
	defer rows.Close()

	networkRecords := map[string]legacyNetworkRecord{}
	for rows.Next() {
		var k string
		var v []byte
		err := rows.Scan(&k, &v)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan network")
		}

		nr := &legacyNetworkRecord{}
		err = json.Unmarshal(v, nr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal network")
		}
		networkRecords[k] = *nr
	}
	return networkRecords, nil
}

func getNewNetworkConfigValues(networkRecords map[string]legacyNetworkRecord, tx *sql.Tx, builder sqorc.StatementBuilder) (map[string]map[string][]byte, error) {
	// migrate network configs in configurator
	// first, grab all configs for each network and delegate to plugins to
	// migrate the binary values
	migratedNwConfigs := map[string]map[string][]byte{}
	for nid, nr := range networkRecords {
		rows, err := builder.Select(configTypeCol, configValCol).
			From(GetLegacyTableName(nid, configTable)).
			Where(squirrel.Eq{configKeyCol: nid}).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrap(err, "failed to query for network configs")
		}
		defer rows.Close()

		migratedNwConfigs[nid] = map[string][]byte{}
		for rows.Next() {
			var t string
			var v []byte
			err := rows.Scan(&t, &v)
			if err != nil {
				return nil, errors.Wrap(err, "failed to scan network config")
			}

			convertedVal, err := migrateConfig(t, v)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert network config %s: %s", t, string(v))
			}
			if convertedVal != nil {
				migratedNwConfigs[nid][t] = convertedVal
			}
		}

		// we also need to pull the "features" part of the network record to
		// a new network config type
		features := &featuresConfig{Features: nr.Features}
		featuresVal, err := json.Marshal(features)
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal network features")
		}
		migratedNwConfigs[nid][NetworkFeaturesConfig] = featuresVal
	}
	return migratedNwConfigs, nil
}

type legacyNetworkRecord struct {
	Name     string            `json:"name"`
	Features map[string]string `json:"features"`
}

type featuresConfig struct {
	Features map[string]string `json:"features,omitempty"`
}
