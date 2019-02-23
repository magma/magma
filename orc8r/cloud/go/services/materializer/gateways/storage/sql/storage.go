/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sql

import (
	"database/sql"
	"fmt"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/sql_utils"
)

type sqlGatewayViewStorage struct {
	db *sql.DB
}

// NewSqlGatewayViewStorage returns a GatewayViewStorage implementation backed
// by SQL. The backing database must support upsert-style statements
// (INSERT ... ON CONFLICT)
func NewSqlGatewayViewStorage(db *sql.DB) storage.GatewayViewStorage {
	return &sqlGatewayViewStorage{db: db}
}

const tableName = "gateway_views"

func GetTableName(networkId string) string {
	return datastore.GetTableName(networkId, tableName)
}

func initTable(tx *sql.Tx, table string) error {
	queryFormat := `
		CREATE TABLE IF NOT EXISTS %s
		(
			gateway_id text PRIMARY KEY,
			status bytea,
			record bytea,
			configs bytea,
			"offset" INTEGER NOT NULL DEFAULT 0
		)
	`
	_, err := tx.Exec(fmt.Sprintf(queryFormat, table))
	return err
}

// No implementation - tables are created "on-demand"
func (s *sqlGatewayViewStorage) InitTables() error {
	return nil
}

func (s *sqlGatewayViewStorage) GetGatewayViewsForNetwork(networkID string) (map[string]*storage.GatewayState, error) {
	return s.GetGatewayViews(networkID, []string{})
}

func (s *sqlGatewayViewStorage) GetGatewayViews(networkID string, gatewayIDs []string) (map[string]*storage.GatewayState, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		return getGatewayStates(tx, networkID, gatewayIDs)
	}

	ret, err := sql_utils.ExecInTx(s.db, getInitFn(networkID), txFn)
	if err != nil {
		return map[string]*storage.GatewayState{}, err
	}
	return ret.(map[string]*storage.GatewayState), nil
}

func (s *sqlGatewayViewStorage) UpdateOrCreateGatewayViews(networkID string, updates map[string]*storage.GatewayUpdateParams) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		existingViews, gwIDs, err := getGatewaysFromUpdateParams(tx, networkID, updates)
		if err != nil {
			return nil, fmt.Errorf("Error loading existing gateway views: %s", err)
		}

		for _, gwID := range gwIDs {
			err = updateGatewayView(tx, networkID, gwID, updates[gwID], existingViews[gwID])
			if err != nil {
				return nil, fmt.Errorf("Error updating gateway %s: %s", gwID, err)
			}
		}
		return nil, nil
	}

	_, err := sql_utils.ExecInTx(s.db, getInitFn(networkID), txFn)
	return err
}

func (s *sqlGatewayViewStorage) DeleteGatewayViews(networkID string, gatewayIDs []string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		query := fmt.Sprintf(
			"DELETE FROM %s WHERE gateway_id IN %s",
			GetTableName(networkID),
			sql_utils.GetPlaceholderArgList(1, len(gatewayIDs)),
		)
		_, err := tx.Exec(query, getSqlGatewayIdArgs(gatewayIDs)...)
		if err != nil {
			return nil, fmt.Errorf("Error deleting gateway views: %s", err)
		}
		return nil, nil
	}

	_, err := sql_utils.ExecInTx(s.db, getInitFn(networkID), txFn)
	return err
}

func getInitFn(networkID string) func(*sql.Tx) error {
	return func(tx *sql.Tx) error {
		return initTable(tx, GetTableName(networkID))
	}
}
