/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package migration

import (
	"fmt"

	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// duplicated constants from configurator
const (
	networksTable      = "cfg_networks"
	networkConfigTable = "cfg_network_configs"

	entityTable      = "cfg_entities"
	entityAssocTable = "cfg_assocs"
	entityAclTable   = "cfg_acls"

	nwIDCol   = "id"
	nwNameCol = "name"
	nwDescCol = "description"
	nwVerCol  = "version"

	nwcIDCol   = "network_id"
	nwcTypeCol = "type"
	nwcValCol  = "value"

	entPkCol   = "pk"
	entNidCol  = "network_id"
	entTypeCol = "type"
	entKeyCol  = "\"key\""
	entGidCol  = "graph_id"
	entNameCol = "name"
	entDescCol = "description"
	entPidCol  = "physical_id"
	entConfCol = "config"
	entVerCol  = "version"

	aFrCol = "from_pk"
	aToCol = "to_pk"

	aclIdCol       = "id"
	aclEntCol      = "entity_pk"
	aclScopeCol    = "scope"
	aclPermCol     = "permission"
	aclTypeCol     = "type"
	aclIdFilterCol = "id_filter"
	aclVerCol      = "version"

	InternalNetworkID          = "network_magma_internal"
	internalNetworkName        = "Internal Magma Network"
	internalNetworkDescription = "Internal network to hold non-network entities"
)

// duplicated constants from blobstore
const (
	blobNidCol  = "network_id"
	blobTypeCol = "type"
	blobKeyCol  = "\"key\""
	blobValCol  = "value"
	blobVerCol  = "version"
)

const (
	deviceServiceTable = "device"
)

// duplicated constants from magmad
const (
	AgRecordTableName = "gatewayRecords"
	HwIdTableName     = "hwIds"
	NetworksTableName = "networks"
	GatewaysTableName = "gateways"
)

// duplicated constants from config
const (
	configTable   = "configurations"
	configTypeCol = "type"
	configKeyCol  = "\"key\""
	configValCol  = "value"
)

const (
	datastoreKeyCol = "\"key\""
	datastoreValCol = "value"
)

const (
	DnsdNetworkType     = "dnsd_network"
	CellularNetworkType = "cellular_network"
)

const (
	NetworkFeaturesConfig = "orc8r_features"
)

func GetLegacyTableName(networkID string, table string) string {
	return fmt.Sprintf("%s_%s", networkID, table)
}

func Migrate(dbDriver string, dbSource string) error {
	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		return errors.Wrap(err, "could not open db connection")
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "error opening tx")
	}

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "error setting tx isolation level")
	}

	// Start by dropping all the new tables so the migration is idempotent
	err = DropNewTables(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Set up the new tables from scratch
	err = SetupTables(tx, sqorc.GetSqlBuilder())
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	sc := squirrel.NewStmtCache(tx)
	defer sc.Clear()

	// migrate networks
	networkIDs, err := MigrateNetworks(sc, sqorc.GetSqlBuilder())
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "failed to migrate networks")
	}

	// migrate gateways
	_, err = MigrateGateways(sc, sqorc.GetSqlBuilder(), networkIDs)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "failed to migrate gateways")
	}

	// TODO: custom per-module migrations

	return tx.Commit()
}
