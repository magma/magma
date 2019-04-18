/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"magma/orc8r/cloud/go/sql_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
	"github.com/thoas/go-funk"
)

const (
	networksTable      = "cfg_networks"
	networkConfigTable = "cfg_network_configs"
)

// NewSQLConfiguratorStorageFactory returns a ConfiguratorStorageFactory
// implementation backed by a SQL database.
func NewSQLConfiguratorStorageFactory(db *sql.DB) ConfiguratorStorageFactory {
	return &sqlConfiguratorStorageFactory{db}
}

type sqlConfiguratorStorageFactory struct {
	db *sql.DB
}

func (fact *sqlConfiguratorStorageFactory) InitializeServiceStorage() (err error) {
	tx, err := fact.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("%s; rollback error: %s", err, rollbackErr)
			}
		}
	}()

	networksTableExec := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			id text PRIMARY KEY,
			name text,
			description text,
			version INTEGER NOT NULL DEFAULT 0
		)
	`, networksTable)

	networksConfigTableExec := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			network_id text REFERENCES %s (id) ON DELETE CASCADE,
			type text NOT NULL,
			value bytea,

			PRIMARY KEY (network_id, type)
		)
	`, networkConfigTable, networksTable)

	// Named return value for err so we can automatically decide to
	// commit/rollback
	tablesToCreate := []string{
		networksTableExec,
		networksConfigTableExec,
	}
	for _, execQuery := range tablesToCreate {
		_, err = tx.Exec(execQuery)
		if err != nil {
			return
		}
	}

	// Create internal network(s)
	_, err = tx.Exec(
		fmt.Sprintf("INSERT INTO %s (id, name, description) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING", networksTable),
		InternalNetworkID, internalNetworkName, internalNetworkDescription,
	)
	if err != nil {
		err = fmt.Errorf("error creating internal networks: %s", err)
		return
	}

	return
}

func (fact *sqlConfiguratorStorageFactory) StartTransaction(ctx context.Context, opts *TxOptions) (ConfiguratorStorage, error) {
	tx, err := fact.db.BeginTx(ctx, getSqlOpts(opts))
	if err != nil {
		return nil, err
	}
	return &sqlConfiguratorStorage{tx: tx}, nil
}

func getSqlOpts(opts *TxOptions) *sql.TxOptions {
	if opts == nil {
		return nil
	}
	return &sql.TxOptions{ReadOnly: opts.ReadOnly}
}

type sqlConfiguratorStorage struct {
	tx *sql.Tx
}

func (store *sqlConfiguratorStorage) Commit() error {
	return store.tx.Commit()
}

func (store *sqlConfiguratorStorage) Rollback() error {
	return store.tx.Rollback()
}

func (store *sqlConfiguratorStorage) LoadNetworks(ids []string, loadCriteria NetworkLoadCriteria) (NetworkLoadResult, error) {
	emptyRet := NetworkLoadResult{NetworkIDsNotFound: []string{}, Networks: []Network{}}
	if len(ids) == 0 {
		return emptyRet, nil
	}

	query, err := getNetworkQuery(ids, loadCriteria)
	if err != nil {
		return emptyRet, err
	}
	queryArgs := make([]interface{}, 0, len(ids))
	funk.ConvertSlice(ids, &queryArgs)

	rows, err := store.tx.Query(query, queryArgs...)
	if err != nil {
		return emptyRet, fmt.Errorf("error querying for networks: %s", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			glog.Warningf("error while closing *Rows object in LoadNetworks: %s", err)
		}
	}()

	// Pointer values because we're modifying .Config in-place
	loadedNetworksByID := map[string]*Network{}
	for rows.Next() {
		nwResult, err := scanNextNetworkRow(rows, loadCriteria)
		if err != nil {
			return emptyRet, err
		}

		existingNetwork, wasLoaded := loadedNetworksByID[nwResult.ID]
		if wasLoaded {
			for k, v := range nwResult.Configs {
				existingNetwork.Configs[k] = v
			}
		} else {
			loadedNetworksByID[nwResult.ID] = &nwResult
		}
	}

	// Sort map keys so we return deterministically
	loadedNetworkIDs := funk.Keys(loadedNetworksByID).([]string)
	sort.Strings(loadedNetworkIDs)

	ret := NetworkLoadResult{
		NetworkIDsNotFound: getNetworkIDsNotFound(loadedNetworksByID, ids),
		Networks:           make([]Network, 0, len(loadedNetworksByID)),
	}
	for _, nid := range loadedNetworkIDs {
		ret.Networks = append(ret.Networks, *loadedNetworksByID[nid])
	}
	return ret, nil
}

func (store *sqlConfiguratorStorage) CreateNetwork(network Network) (Network, error) {
	exists, err := store.doesNetworkExist(network.ID)
	if err != nil {
		return network, err
	}
	if exists {
		return network, fmt.Errorf("a network with ID %s already exists", network.ID)
	}

	// Insert the network, then insert its configs
	exec := fmt.Sprintf("INSERT INTO %s (id, name, description) VALUES ($1, $2, $3)", networksTable)
	_, err = store.tx.Exec(exec, network.ID, network.Name, network.Description)
	if err != nil {
		return network, fmt.Errorf("error inserting network: %s", err)
	}

	if funk.IsEmpty(network.Configs) {
		return network, nil
	}

	configExec := fmt.Sprintf("INSERT INTO %s (network_id, type, value) VALUES ($1, $2, $3)", networkConfigTable)
	configInsertStatement, err := store.tx.Prepare(configExec)
	if err != nil {
		return network, fmt.Errorf("error preparing network configuration insert: %s", err)
	}
	defer sql_utils.GetCloseStatementsDeferFunc([]*sql.Stmt{configInsertStatement}, "CreateNetwork")()

	// Sort config keys for deterministic behavior
	configKeys := funk.Keys(network.Configs).([]string)
	sort.Strings(configKeys)
	for _, configKey := range configKeys {
		_, err := configInsertStatement.Exec(network.ID, configKey, network.Configs[configKey])
		if err != nil {
			return network, fmt.Errorf("error inserting config %s: %s", configKey, err)
		}
	}

	return network, nil
}

func (store *sqlConfiguratorStorage) UpdateNetworks(updates []NetworkUpdateCriteria) (FailedOperations, error) {
	failures := FailedOperations{}
	if err := validateNetworkUpdates(updates); err != nil {
		return failures, err
	}

	// Prepare statements
	deleteExec := fmt.Sprintf("DELETE FROM %s WHERE id = $1", networksTable)
	upsertConfigExec := fmt.Sprintf(`
		INSERT INTO %s (network_id, type, value) VALUES ($1, $2, $3)
		ON CONFLICT (network_id, type) DO UPDATE SET value = $4
	`, networkConfigTable)
	deleteConfigExec := fmt.Sprintf("DELETE FROM %s WHERE (network_id, type) = ($1, $2)", networkConfigTable)
	statements, err := sql_utils.PrepareStatements(store.tx, []string{deleteExec, upsertConfigExec, deleteConfigExec})
	if err != nil {
		return failures, err
	}
	defer sql_utils.GetCloseStatementsDeferFunc(statements, "UpdateNetworks")()

	deleteStmt, upsertConfigStmt, deleteConfigStmt := statements[0], statements[1], statements[2]
	for _, update := range updates {
		if update.DeleteNetwork {
			_, err := deleteStmt.Exec(update.ID)
			if err != nil {
				failures[update.ID] = fmt.Errorf("error deleting network %s: %s", update.ID, err)
			}
		} else {
			err := store.updateNetwork(update, upsertConfigStmt, deleteConfigStmt)
			if err != nil {
				failures[update.ID] = err
			}
		}
	}

	if funk.IsEmpty(failures) {
		return failures, nil
	}
	return failures, errors.New("some errors were encountered, see return value for details")
}

func (store *sqlConfiguratorStorage) LoadEntities(ids []storage.TypeAndKey, loadCriteria EntityLoadCriteria) (EntityLoadResult, error) {
	panic("implement me")
}

func (store *sqlConfiguratorStorage) CreateEntities(entities []NetworkEntity) (EntityCreationResult, error) {
	panic("implement me")
}

func (store *sqlConfiguratorStorage) UpdateEntities(updates []EntityUpdateCriteria) (FailedOperations, error) {
	panic("implement me")
}

func (store *sqlConfiguratorStorage) LoadGraphForEntity(entityID storage.TypeAndKey, loadCriteria EntityLoadCriteria) (EntityGraph, error) {
	panic("implement me")
}
