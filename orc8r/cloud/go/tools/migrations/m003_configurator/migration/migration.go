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
	NetworksTable      = "cfg_networks"
	NetworkConfigTable = "cfg_network_configs"

	EntityTable      = "cfg_entities"
	EntityAssocTable = "cfg_assocs"
	EntityAclTable   = "cfg_acls"

	NwIDCol   = "id"
	NwNameCol = "name"
	NwDescCol = "description"
	NwVerCol  = "version"

	NwcIDCol   = "network_id"
	NwcTypeCol = "type"
	NwcValCol  = "value"

	EntPkCol   = "pk"
	EntNidCol  = "network_id"
	EntTypeCol = "type"
	EntKeyCol  = "\"key\""
	EntGidCol  = "graph_id"
	EntNameCol = "name"
	EntDescCol = "description"
	EntPidCol  = "physical_id"
	EntConfCol = "config"
	EntVerCol  = "version"

	AFrCol = "from_pk"
	AToCol = "to_pk"

	AclIdCol       = "id"
	AclEntCol      = "entity_pk"
	AclScopeCol    = "scope"
	AclPermCol     = "permission"
	AclTypeCol     = "type"
	AclIdFilterCol = "id_filter"
	AclVerCol      = "version"

	InternalNetworkID          = "network_magma_internal"
	internalNetworkName        = "Internal Magma Network"
	internalNetworkDescription = "Internal network to hold non-network entities"
)

// duplicated constants from blobstore
const (
	BlobNidCol  = "network_id"
	BlobTypeCol = "type"
	BlobKeyCol  = "\"key\""
	BlobValCol  = "value"
	BlobVerCol  = "version"
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
	ConfigTable   = "configurations"
	ConfigTypeCol = "type"
	ConfigKeyCol  = "\"key\""
	ConfigValCol  = "value"
)

const (
	DatastoreKeyCol     = "\"key\""
	DatastoreValCol     = "value"
	DatastoreGenCol     = "generation_number"
	DatastoreDeletedCol = "deleted"
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
	migratedGatewayMetas, err := MigrateGateways(sc, sqorc.GetSqlBuilder(), networkIDs)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "failed to migrate gateways")
	}

	err = RunCustomPluginMigrations(sc, sqorc.GetSqlBuilder(), migratedGatewayMetas)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "failed to run custom plugin migration jobs")
	}

	return tx.Commit()
}
