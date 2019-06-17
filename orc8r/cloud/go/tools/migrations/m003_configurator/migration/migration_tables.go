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
		networksTable,
		networkConfigTable,
		entityTable,
		entityAssocTable,
		entityAclTable,
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
	_, err := builder.CreateTable(networksTable).
		IfNotExists().
		Column(nwIDCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(nwNameCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(nwDescCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(nwVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create networks table")
	}

	_, err = builder.CreateTable(networkConfigTable).
		IfNotExists().
		Column(nwcIDCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(nwcTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(nwcValCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		PrimaryKey(nwcIDCol, nwcTypeCol).
		ForeignKey(networksTable, map[string]string{nwcIDCol: nwIDCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create network configs table")
	}

	// Create an internal-only primary key (UUID) for entities.
	// This keeps index size in control and supporting table schemas simpler.
	_, err = builder.CreateTable(entityTable).
		IfNotExists().
		Column(entPkCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(entNidCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(entTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(entKeyCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(entGidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(entNameCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(entDescCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(entPidCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(entConfCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(entVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		Unique(entNidCol, entKeyCol, entTypeCol).
		ForeignKey(networksTable, map[string]string{entNidCol: nwIDCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create entities table")
	}

	_, err = builder.CreateTable(entityAssocTable).
		IfNotExists().
		Column(aFrCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(aToCol).Type(sqorc.ColumnTypeText).EndColumn().
		PrimaryKey(aFrCol, aToCol).
		ForeignKey(entityTable, map[string]string{aFrCol: entPkCol}, sqorc.ColumnOnDeleteCascade).
		ForeignKey(entityTable, map[string]string{aToCol: entPkCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create entity assoc table")
	}

	_, err = builder.CreateTable(entityAclTable).
		IfNotExists().
		Column(aclIdCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(aclEntCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(aclScopeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(aclPermCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
		Column(aclTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(aclIdFilterCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(aclVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		ForeignKey(entityTable, map[string]string{aclEntCol: entPkCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create entity acl table")
	}

	// Create indexes (index is not implicitly created on a referencing FK)
	_, err = builder.CreateIndex("graph_id_idx").
		IfNotExists().
		On(entityTable).
		Columns(entGidCol).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create graph ID index")
	}

	_, err = builder.CreateIndex("acl_ent_pk_idx").
		IfNotExists().
		On(entityAclTable).
		Columns(aclEntCol).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create acl ent PK index")
	}

	// Create internal network(s)
	_, err = builder.Insert(networksTable).
		Columns(nwIDCol, nwNameCol, nwDescCol).
		Values(InternalNetworkID, internalNetworkName, internalNetworkDescription).
		OnConflict(nil, nwIDCol).
		RunWith(tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "error creating internal networks")
	}
	return nil
}
