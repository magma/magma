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
	"sort"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	sql_protos "magma/orc8r/cloud/go/services/materializer/gateways/storage/sql/protos"
	"magma/orc8r/cloud/go/sql_utils"
)

// Returned map is guaranteed to have all the keys in updates param
// Gateways without existing views will have an empty GatewayState populated.
// 2nd returned argument is the sorted gwIDs from the updates param
func getGatewaysFromUpdateParams(tx *sql.Tx, networkID string, updates map[string]*storage.GatewayUpdateParams) (map[string]*storage.GatewayState, []string, error) {
	gwIDs := make([]string, 0, len(updates))
	for k := range updates {
		gwIDs = append(gwIDs, k)
	}
	// sort IDs so we get deterministic ordering for tests
	sort.Strings(gwIDs)

	existingViews, err := getGatewayStates(tx, networkID, gwIDs)
	if err != nil {
		return map[string]*storage.GatewayState{}, gwIDs, err
	}

	// Populate missing existing views with stub struct
	for id := range updates {
		if _, ok := existingViews[id]; !ok {
			existingViews[id] = &storage.GatewayState{GatewayID: id}
		}
	}
	return existingViews, gwIDs, nil
}

// Update a single gateway view
func updateGatewayView(tx *sql.Tx, networkID string, gatewayID string, update *storage.GatewayUpdateParams, existingView *storage.GatewayState) error {
	query, args, err := getUpsertQueryAndArgs(networkID, gatewayID, existingView.Config, update)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, args...)
	return err
}

// Construct the upsert query
func getUpsertQueryAndArgs(networkID string, gatewayID string, existingConfigs map[string]interface{}, update *storage.GatewayUpdateParams) (string, []interface{}, error) {
	updateArgList, updateArgValues, err := getUpdateArgListAndValues(existingConfigs, update)
	if err != nil {
		return "", []interface{}{}, err
	}

	// INSERT clause needs gateway_id argument, while UPDATE clause excludes it
	insertArgList := []string{"gateway_id"}
	insertArgList = append(insertArgList, updateArgList...)
	insertArgValues := []interface{}{gatewayID}
	insertArgValues = append(insertArgValues, updateArgValues...)

	// INSERT INTO table_name (gateway_id, col1, col2) VALUES ($1, $2, $3)
	// ON CONFLICT DO UPDATE SET col1 = $4, col2 = $5 WHERE gateway_id = $6
	query := fmt.Sprintf(
		"INSERT INTO %s %s VALUES %s ON CONFLICT(gateway_id) DO UPDATE SET %s WHERE %s.gateway_id = $%d",
		GetTableName(networkID),
		sql_utils.GetInsertArgListString(insertArgList...),
		sql_utils.GetPlaceholderArgList(1, len(insertArgList)),
		sql_utils.GetUpdateClauseString(len(insertArgList)+1, updateArgList...),
		GetTableName(networkID),
		len(insertArgList)+len(updateArgList)+1,
	)

	var allArgs []interface{}
	allArgs = append(allArgs, insertArgValues...)
	allArgs = append(allArgs, updateArgValues...)
	allArgs = append(allArgs, gatewayID)
	return query, allArgs, nil
}

// Get column names and corresponding values for an UPDATE
func getUpdateArgListAndValues(existingConfigs map[string]interface{}, update *storage.GatewayUpdateParams) ([]string, []interface{}, error) {
	var updateArgList []string
	var updateArgValues []interface{}

	if update.NewStatus != nil {
		marshaledStatus, err := protos.MarshalIntern(update.NewStatus)
		if err != nil {
			return []string{}, []interface{}{}, fmt.Errorf("Failed to marshal updated status: %s", err)
		}

		updateArgList = append(updateArgList, "status")
		updateArgValues = append(updateArgValues, marshaledStatus)
	}

	if update.NewRecord != nil {
		marshaledRecord, err := protos.MarshalIntern(update.NewRecord)
		if err != nil {
			return []string{}, []interface{}{}, fmt.Errorf("Failed to marshal updated record: %s", err)
		}

		updateArgList = append(updateArgList, "record")
		updateArgValues = append(updateArgValues, marshaledRecord)
	}

	shouldUpdateConfig := update.NewConfig != nil || update.ConfigsToDelete != nil
	if shouldUpdateConfig {
		newConfigs, err := getNewConfigs(existingConfigs, update.NewConfig, update.ConfigsToDelete)
		if err != nil {
			return []string{}, []interface{}{}, err
		}

		newConfigsMarshaled, err := protos.MarshalIntern(newConfigs)
		if err != nil {
			return []string{}, []interface{}{}, fmt.Errorf("Failed to marshal new configs: %s", err)
		}

		updateArgList = append(updateArgList, "configs")
		updateArgValues = append(updateArgValues, newConfigsMarshaled)
	}

	updateArgList = append(updateArgList, "\"offset\"")
	updateArgValues = append(updateArgValues, update.Offset)

	return updateArgList, updateArgValues, nil
}

func getNewConfigs(existingConfigs map[string]interface{}, configsToAddOrUpdate map[string]interface{}, configsToDelete []string) (*sql_protos.ViewConfigs, error) {
	newConfigs := map[string]interface{}{}
	for k, v := range existingConfigs {
		newConfigs[k] = v
	}
	for k, v := range configsToAddOrUpdate {
		newConfigs[k] = v
	}
	for _, k := range configsToDelete {
		delete(newConfigs, k)
	}

	marshaledConfigs, err := marshalConfigs(newConfigs)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling updated configs: %s", err)
	}
	return &sql_protos.ViewConfigs{Configs: marshaledConfigs}, nil
}

func marshalConfigs(configs map[string]interface{}) (map[string][]byte, error) {
	ret := make(map[string][]byte, len(configs))
	for k, v := range configs {
		marshaledV, err := registry.MarshalConfig(k, v)
		if err != nil {
			return map[string][]byte{}, err
		}
		ret[k] = marshaledV
	}
	return ret, nil
}
