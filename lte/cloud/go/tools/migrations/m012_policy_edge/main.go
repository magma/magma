/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

/*
	m012_policy_edge ensures assoc directions between base names, subscribers,
	and policies are as intended.

    Enforces the following assoc directions by swapping when necessary:
        - base_name -> policy_rule
        - base_name -> subscriber
        - subscriber -> policy_rule

	This migration touches a lot of SQL objects within a serializable
	transaction. If the migration fails due to serialization failure, consider
	retrying.

	The migration executable imports main library code in two ways
		- sqorc	-- statement builders
	If these imports break this executable, recourse to the following
		- sqorc	-- fix breakage, or change statement builders to manual SQL
				   string execs

	In a single, serializable transaction, performs roughly the following,
    per enforced assoc direction:

    For the desired assoc direction from child -> parent, where it was
    previously parent -> child:

    Get all assocs with parent pointing to another entity

    Get all PKs of child entity

    Filter all assocs from parent entity so the only assocs remaining are
    parent -> child

    Swap all remaining assocs with update statements
*/

package main

import (
	"database/sql"
	"flag"
	"fmt"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	baseNameEntType   = "base_name"
	subscriberEntType = "subscriber"
	policyEntType     = "policy"
)

// Duplicated from configurator
const (
	entityTable      = "cfg_entities"
	entityAssocTable = "cfg_assocs"

	entPkCol   = "pk"
	entNidCol  = "network_id"
	entTypeCol = "type"

	aFrCol = "from_pk"
	aToCol = "to_pk"
)

var (
	dryRun bool
)

type assocPair struct {
	fromPk string
	toPk   string
}

func init() {
	flag.BoolVar(&dryRun, "dry", false, "don't commit changes")
	flag.Parse()
	_ = flag.Set("alsologtostderr", "true") // enable printing to console
}

// Migration without SQL error will result in `SUCCESS` printed to stderr.
func main() {
	defer glog.Flush()

	glog.Info("BEGIN MIGRATION")

	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		glog.Fatal(errors.Wrap(err, "could not open db connection"))
	}
	builder := sqorc.GetSqlBuilder()

	txFn := func(tx *sql.Tx) (interface{}, error) {
		glog.Info("RUNNING TRANSACTION")

		// Enforce assoc direction from policy_rule -> base_name
		bnToPolicyAssocs, err := getParentToChildAssocs(tx, builder, baseNameEntType, policyEntType)
		if err != nil {
			return nil, err
		}
		glog.Infof("Swapping %d assoc directions to enforce %s -> %s", len(bnToPolicyAssocs), policyEntType, baseNameEntType)
		err = swapAssocDirections(tx, builder, bnToPolicyAssocs)
		if err != nil {
			return nil, err
		}

		// Enforce assoc direction from subscriber -> base_name
		bnToSubAssocs, err := getParentToChildAssocs(tx, builder, baseNameEntType, subscriberEntType)
		if err != nil {
			return nil, err
		}
		glog.Infof("Swapping %d assoc directions to enforce %s -> %s", len(bnToSubAssocs), subscriberEntType, baseNameEntType)
		err = swapAssocDirections(tx, builder, bnToSubAssocs)
		if err != nil {
			return nil, err
		}

		// Enforce assoc direction from subscriber -> policy_rule
		policyToSubAssocs, err := getParentToChildAssocs(tx, builder, policyEntType, subscriberEntType)
		if err != nil {
			return nil, err
		}
		glog.Infof("Swapping %d assoc directions to enforce %s -> %s", len(policyToSubAssocs), subscriberEntType, policyEntType)
		err = swapAssocDirections(tx, builder, policyToSubAssocs)
		if err != nil {
			return nil, err
		}

		if dryRun {
			return nil, errors.New("throwing error instead of committing because -dry was specified")
		}

		return nil, nil
	}
	_, err = migrations.ExecInTx(db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, txFn)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Info("SUCCESS")
	glog.Info("END MIGRATION")
}

func getParentToChildAssocs(tx *sql.Tx, builder sqorc.StatementBuilder, parentType string, childType string) ([]assocPair, error) {
	parentAssocs, err := getEntAssocs(tx, builder, parentType)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("getParentToChildAssocs: get %s assocs", parentType))
	}
	childPks, err := getEntPks(tx, builder, childType)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("getParentToChildAssocs: get %s pks", childType))
	}
	// Filter parent type associations for entries pointing to child type
	parentToChildAssocs := []assocPair{}
	for i := range parentAssocs {
		toPk := parentAssocs[i].toPk
		_, isChildType := childPks[toPk]
		if isChildType {
			parentToChildAssocs = append(parentToChildAssocs, parentAssocs[i])
		}
	}
	return parentToChildAssocs, nil
}

func getEntAssocs(tx *sql.Tx, builder sqorc.StatementBuilder, entType string) ([]assocPair, error) {
	// SELECT assocs.from_pk, assocs.to_pk, ents.type
	// FROM cfg_assocs AS assocs
	// JOIN cfg_entities AS ents ON ents.type = 'entType' AND assocs.from_pk = ents.pk;
	rows, err := builder.
		Select(fmt.Sprintf("assocs.%s, assocs.%s, ents.%s", aFrCol, aToCol, entTypeCol)).
		From(fmt.Sprintf("%s AS assocs", entityAssocTable)).
		Join(fmt.Sprintf("%s AS ents ON ents.%s='%s' AND assocs.%s=ents.%s", entityTable, entTypeCol, entType, aFrCol, entPkCol)).
		RunWith(tx).Query()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("getEntAssocs: select %s assocs", entType))
	}

	assocs := []assocPair{}
	for rows.Next() {
		fromPk, toPk, entType := &sql.NullString{}, &sql.NullString{}, &sql.NullString{}
		err = rows.Scan(fromPk, toPk, entType)
		if err != nil {
			return nil, errors.Wrap(err, "getEntAssocs: scan assocs")
		}
		assocs = append(assocs, assocPair{fromPk: fromPk.String, toPk: toPk.String})
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "getEntAssocs: SQL rows error")
	}
	return assocs, nil
}

func getEntPks(tx *sql.Tx, builder sqorc.StatementBuilder, entType string) (map[string]bool, error) {
	// SELECT pk
	// FROM cfg_entities
	// WHERE type='entType';
	rows, err := builder.
		Select(entPkCol).
		From(entityTable).
		Where(fmt.Sprintf("%s='%s'", entTypeCol, entType)).
		RunWith(tx).Query()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("getEntPks: select %s pks", entType))
	}
	policyPks := map[string]bool{}
	for rows.Next() {
		pk := &sql.NullString{}
		err = rows.Scan(pk)
		if err != nil {
			return nil, errors.Wrap(err, "getEntPks: scan query results")
		}
		policyPks[pk.String] = true
	}
	return policyPks, nil
}

func swapAssocDirections(tx *sql.Tx, builder sqorc.StatementBuilder, assocs []assocPair) error {
	for _, assoc := range assocs {
		err := swapAssocDirection(tx, builder, assoc)
		if err != nil {
			return err
		}
	}
	return nil
}

func swapAssocDirection(tx *sql.Tx, builder sqorc.StatementBuilder, assoc assocPair) error {
	b := builder.
		Update(entityAssocTable).
		Set(aFrCol, assoc.toPk).
		Set(aToCol, assoc.fromPk).
		Where(squirrel.Eq{"from_pk": assoc.fromPk, "to_pk": assoc.toPk})
	sqlStr, args, _ := b.ToSql()
	glog.Infof("[RUN] %s %v", sqlStr, args)
	_, err := b.RunWith(tx).Exec()
	if err != nil {
		return errors.Wrap(err, "swapAssocDirection: update error")
	}
	return nil
}

func makeCol(table, col string) string {
	return fmt.Sprintf("%s.%s", table, col)
}
