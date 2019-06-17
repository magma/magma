/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package migration

import (
	"encoding/json"

	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func MigrateGateways(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkIDs []string) (map[string]map[string]MigratedGatewayMeta, error) {
	migratedGatewayMetaByNetwork := map[string]map[string]MigratedGatewayMeta{}
	for _, nid := range networkIDs {
		migratedGWs, err := migrateGatewaysForNetwork(sc, builder, nid)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		migratedGatewayMetaByNetwork[nid] = migratedGWs
	}
	return migratedGatewayMetaByNetwork, nil
}

func migrateGatewaysForNetwork(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string) (map[string]MigratedGatewayMeta, error) {
	gwIDs, migratedIDsByGw, err := migrateGatewayRecords(sc, builder, networkID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if funk.IsEmpty(gwIDs) {
		return map[string]MigratedGatewayMeta{}, nil
	}

	oldConfigsByID, err := loadAllOldGatewayConfigs(sc, builder, networkID, gwIDs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for gwID, legacyConfigs := range oldConfigsByID {
		for ctype, oldVal := range legacyConfigs {
			newVal, err := migrateConfig(ctype, oldVal)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to migrate config %s for gateway %s", ctype, gwID)
			}
			if newVal == nil {
				continue
			}

			switch ctype {
			case "magmad_gateway":
				_, err := builder.Update(entityTable).
					Set(entConfCol, newVal).
					Where(squirrel.Eq{
						entNidCol:  networkID,
						entTypeCol: "magmad_gateway",
						entKeyCol:  gwID,
					}).
					RunWith(sc).
					Exec()
				if err != nil {
					return nil, errors.Wrapf(err, "failed to update magmad gateway %s with migrated config", gwID)
				}
				break
			default:
				newEntPk := uuid.New().String()
				_, err := builder.Insert(entityTable).
					Columns(entPkCol, entNidCol, entTypeCol, entKeyCol, entConfCol, entGidCol).
					Values(newEntPk, networkID, ctype, gwID, newVal, migratedIDsByGw[gwID].GraphID).
					RunWith(sc).
					Exec()
				if err != nil {
					return nil, errors.Wrapf(err, "failed to create new entity for %s with key %s", ctype, gwID)
				}

				_, err = builder.Update(entityAssocTable).
					Set(aFrCol, migratedIDsByGw[gwID].Pk).
					Set(aToCol, newEntPk).
					RunWith(sc).
					Exec()
				if err != nil {
					return nil, errors.Wrapf(err, "failed to associate magmad gateway with new entity (%s, %s)", ctype, gwID)
				}
				break
			}
		}
	}

	return migratedIDsByGw, nil
}

type legacyGatewayConfigs map[string][]byte

func loadAllOldGatewayConfigs(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string, gwIDs []string) (map[string]legacyGatewayConfigs, error) {
	rows, err := builder.Select(configTypeCol, configKeyCol, configValCol).
		From(GetLegacyTableName(networkID, configTable)).
		Where(squirrel.Eq{configKeyCol: gwIDs}).
		RunWith(sc).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "failed to query for gateway configs")
	}
	defer rows.Close()

	ret := map[string]legacyGatewayConfigs{}
	for rows.Next() {
		var id, ctype string
		var oldVal []byte

		err := rows.Scan(&ctype, &id, &oldVal)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan gateway config row")
		}

		legacyConfigs, ok := ret[id]
		if !ok {
			legacyConfigs = legacyGatewayConfigs{}
			ret[id] = legacyConfigs
		}
		legacyConfigs[ctype] = oldVal
	}

	return ret, nil
}

// returns gw ids and map between gw id and graph id
type MigratedGatewayMeta struct {
	Pk, GraphID string
}

func migrateGatewayRecords(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string) ([]string, map[string]MigratedGatewayMeta, error) {
	rows, err := builder.Select(datastoreKeyCol, datastoreValCol).
		From(GetLegacyTableName(networkID, AgRecordTableName)).
		RunWith(sc).
		Query()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to query for gateways in network %s", networkID)
	}
	defer rows.Close()

	// First, load all gateway records
	gwRecords := map[string]*migratedGatewayRecord{}
	for rows.Next() {
		var k string
		var v []byte
		err := rows.Scan(&k, &v)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to scan gateway")
		}

		rec := &legacyGatewayRecord{}
		err = json.Unmarshal(v, rec)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to unmarshal gateway record %s", k)
		}

		newRec := &migratedGatewayRecord{
			HwID: &migratedHwID{ID: rec.HwID.ID},
			Key: &migratedChallengeKey{
				Key:     rec.Key.Key,
				KeyType: rec.Key.KeyType,
			},
			Name: rec.Name,
		}
		gwRecords[k] = newRec
	}
	if funk.IsEmpty(gwRecords) {
		return []string{}, map[string]MigratedGatewayMeta{}, nil
	}

	// We'll migrate the gateway records into the device service and create
	// placeholder entities for the logical access gateways in configurator.
	// We'll fill in the configs from magmad configs later
	migratedMetasByID := map[string]MigratedGatewayMeta{}
	recInsertBuilder := builder.Insert(deviceServiceTable).
		Columns(blobNidCol, blobTypeCol, blobKeyCol, blobValCol).
		RunWith(sc)
	entInsertBuilder := builder.Insert(entityTable).
		Columns(entPkCol, entNidCol, entTypeCol, entKeyCol, entPidCol, entGidCol).
		RunWith(sc)
	for logicalID, record := range gwRecords {
		marshaledRecord, err := json.Marshal(record)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to marshal migrated gateway record %s", logicalID)
		}

		graphID := uuid.New().String()
		pk := uuid.New().String()
		migratedMetasByID[logicalID] = MigratedGatewayMeta{GraphID: graphID, Pk: pk}

		recInsertBuilder = recInsertBuilder.Values(networkID, "access_gateway_record", record.HwID.ID, marshaledRecord)
		entInsertBuilder = entInsertBuilder.Values(
			pk,
			networkID,
			"magmad_gateway",
			logicalID,
			record.HwID.ID,
			graphID,
		)
	}

	_, err = recInsertBuilder.Exec()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to insert migrated gateway records")
	}
	_, err = entInsertBuilder.Exec()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to insert migrated access gateway entities")
	}

	return funk.Keys(gwRecords).([]string), migratedMetasByID, nil
}

type legacyGatewayRecord struct {
	HwID *legacyHwID         `json:"hwId"`
	Name string              `json:"name"`
	Key  *legacyChallengeKey `json:"key"`
}

type legacyHwID struct {
	ID string `json:"id"`
}

type legacyChallengeKey struct {
	KeyType string `json:"keyType"`
	Key     []byte `json:"key"`
}

type migratedGatewayRecord struct {
	HwID *migratedHwID         `json:"hw_id" magma_alt_name:"HwId"`
	Key  *migratedChallengeKey `json:"key"`
	Name string                `json:"name,omitempty"`
}

type migratedHwID struct {
	ID string `json:"id"`
}

type migratedChallengeKey struct {
	Key     []byte `json:"key"`
	KeyType string `json:"key_type"`
}
