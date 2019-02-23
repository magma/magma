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

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/registry"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	sql_protos "magma/orc8r/cloud/go/services/materializer/gateways/storage/sql/protos"
	"magma/orc8r/cloud/go/sql_utils"
)

func getGatewayStates(tx *sql.Tx, networkID string, gatewayIDs []string) (map[string]*storage.GatewayState, error) {
	query, args := getSelectQueryAndArgs(networkID, gatewayIDs)
	rows, err := tx.Query(query, args...)
	if err != nil {
		return map[string]*storage.GatewayState{}, fmt.Errorf("Storage query error: %s", err)
	}
	defer rows.Close()

	scannedRows := map[string]*storage.GatewayState{}
	for rows.Next() {
		var id string
		var status, record, configs []byte
		var offset uint64

		err = rows.Scan(&id, &status, &record, &configs, &offset)
		if err != nil {
			return map[string]*storage.GatewayState{}, fmt.Errorf("Storage read error: %s", err)
		}
		gwState, err := createGatewayState(id, status, record, configs, offset)
		if err != nil {
			return map[string]*storage.GatewayState{}, fmt.Errorf("Could not unmarshal gateway %s: %s", id, err)
		}
		scannedRows[id] = gwState
	}
	return scannedRows, nil
}

// empty or nil gatewayIDs arg will select all rows
func getSelectQueryAndArgs(networkID string, gatewayIDs []string) (string, []interface{}) {
	selectClause := fmt.Sprintf("SELECT gateway_id, status, record, configs, \"offset\" FROM %s", GetTableName(networkID))
	if gatewayIDs == nil || len(gatewayIDs) == 0 {
		return selectClause, []interface{}{}
	}

	inList := sql_utils.GetPlaceholderArgList(1, len(gatewayIDs))
	return fmt.Sprintf("%s WHERE gateway_id IN %s", selectClause, inList), getSqlGatewayIdArgs(gatewayIDs)
}

func getSqlGatewayIdArgs(gatewayIDs []string) []interface{} {
	ret := make([]interface{}, 0, len(gatewayIDs))
	for _, id := range gatewayIDs {
		ret = append(ret, id)
	}
	return ret
}

func createGatewayState(id string, status []byte, record []byte, configs []byte, offset uint64) (*storage.GatewayState, error) {
	ret := &storage.GatewayState{GatewayID: id, Offset: int64(offset)}

	if status != nil && len(status) > 0 {
		statusProto := &protos.GatewayStatus{}
		err := protos.Unmarshal(status, statusProto)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling gateway status: %s", err)
		}
		ret.Status = statusProto
	}

	if record != nil && len(record) > 0 {
		recordProto := &magmadprotos.AccessGatewayRecord{}
		err := protos.Unmarshal(record, recordProto)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshaling gateway record: %s", err)
		}
		ret.Record = recordProto
	}

	unmarshaledConfigs, err := unmarshalConfigs(configs)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling gateway configs: %s", err)
	}
	ret.Config = unmarshaledConfigs
	return ret, nil
}

func unmarshalConfigs(configs []byte) (map[string]interface{}, error) {
	if configs == nil || len(configs) == 0 {
		return map[string]interface{}{}, nil
	}

	configProto := &sql_protos.ViewConfigs{}
	err := protos.Unmarshal(configs, configProto)
	if err != nil {
		return map[string]interface{}{}, err
	}

	ret := make(map[string]interface{}, len(configProto.Configs))
	for configType, configVal := range configProto.Configs {
		unmarshaledConfig, err := registry.UnmarshalConfig(configType, configVal)
		if err != nil {
			return map[string]interface{}{}, err
		}
		ret[configType] = unmarshaledConfig
	}
	return ret, nil
}
