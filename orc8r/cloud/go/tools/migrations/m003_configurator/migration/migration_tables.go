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
	"fmt"

	"magma/orc8r/cloud/go/sqorc"

	"github.com/pkg/errors"
)

func DropNewTables(tx *sql.Tx) error {
	tablesToDrop := []string{
		NetworksTable,
		NetworkConfigTable,
		EntityTable,
		EntityAssocTable,
		EntityAclTable,
		deviceServiceTable,
		StateServiceTable,
	}

	for _, tableName := range tablesToDrop {
		_, err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName))
		if err != nil {
			return errors.Wrapf(err, "failed to drop table %s", tableName)
		}
	}
	return nil
}

func SetupTables(tx *sql.Tx, builder sqorc.StatementBuilder) error {
	// device service tables
	_, err := builder.CreateTable(deviceServiceTable).
		IfNotExists().
		Column(BlobNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(BlobTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(BlobKeyCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(BlobValCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(BlobVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		PrimaryKey(BlobNidCol, BlobTypeCol, BlobKeyCol).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create devices table")
	}

	// state service tables
	_, err = builder.CreateTable(StateServiceTable).
		IfNotExists().
		Column(BlobNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(BlobTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(BlobKeyCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(BlobValCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(BlobVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		PrimaryKey(BlobNidCol, BlobTypeCol, BlobKeyCol).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create states table")
	}

	// configurator tables
	_, err = builder.CreateTable(NetworksTable).
		IfNotExists().
		Column(NwIDCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(NwNameCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(NwDescCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(NwVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create networks table")
	}

	_, err = builder.CreateTable(NetworkConfigTable).
		IfNotExists().
		Column(NwcIDCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(NwcTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(NwcValCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		PrimaryKey(NwcIDCol, NwcTypeCol).
		ForeignKey(NetworksTable, map[string]string{NwcIDCol: NwIDCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create network configs table")
	}

	// Create an internal-only primary key (UUID) for entities.
	// This keeps index size in control and supporting table schemas simpler.
	_, err = builder.CreateTable(EntityTable).
		IfNotExists().
		Column(EntPkCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(EntNidCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(EntTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(EntKeyCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(EntGidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(EntNameCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(EntDescCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(EntPidCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(EntConfCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(EntVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		Unique(EntNidCol, EntKeyCol, EntTypeCol).
		ForeignKey(NetworksTable, map[string]string{EntNidCol: NwIDCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create entities table")
	}

	_, err = builder.CreateTable(EntityAssocTable).
		IfNotExists().
		Column(AFrCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(AToCol).Type(sqorc.ColumnTypeText).EndColumn().
		PrimaryKey(AFrCol, AToCol).
		ForeignKey(EntityTable, map[string]string{AFrCol: EntPkCol}, sqorc.ColumnOnDeleteCascade).
		ForeignKey(EntityTable, map[string]string{AToCol: EntPkCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create entity assoc table")
	}

	_, err = builder.CreateTable(EntityAclTable).
		IfNotExists().
		Column(AclIdCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(AclEntCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(AclScopeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(AclPermCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
		Column(AclTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(AclIdFilterCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(AclVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		ForeignKey(EntityTable, map[string]string{AclEntCol: EntPkCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create entity acl table")
	}

	// Create indexes (index is not implicitly created on a referencing FK)
	_, err = builder.CreateIndex("graph_id_idx").
		IfNotExists().
		On(EntityTable).
		Columns(EntGidCol).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create graph ID index")
	}

	_, err = builder.CreateIndex("acl_ent_pk_idx").
		IfNotExists().
		On(EntityAclTable).
		Columns(AclEntCol).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create acl ent PK index")
	}

	_, err = builder.CreateIndex("phys_id_idx").
		IfNotExists().
		On(EntityTable).
		Columns(EntPidCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create physical ID index")
	}

	// Create internal network(s)
	_, err = builder.Insert(NetworksTable).
		Columns(NwIDCol, NwNameCol, NwDescCol).
		Values(InternalNetworkID, internalNetworkName, internalNetworkDescription).
		OnConflict(nil, NwIDCol).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "error creating internal networks")
	}
	return nil
}
