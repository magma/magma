/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"magma/orc8r/cloud/go/sql_utils"

	"github.com/thoas/go-funk"
)

type networkQueryArgs struct {
	Fields, TableName, JoinQuery, IdList string
}

func getNetworkQuery(ids []string, criteria NetworkLoadCriteria) (string, error) {
	// SELECT cfg_networks.id, cfg_networks.name, cfg_networks.description cfg_network_configs.type, cfg_network_configs.value
	// FROM cfg_networks
	// [[ LEFT JOIN cfg_network_configs ON cfg_networks.id = cfg_network_configs.network_id ]]
	// WHERE cfg_networks.id IN (id list)
	queryTemplate := template.Must(template.New("nw_query").Parse(
		"SELECT {{.Fields}} FROM {{.TableName}} " +
			"{{.JoinQuery}} " +
			"WHERE {{.TableName}}.id IN {{.IdList}}",
	))
	args := getNetworkQueryArgs(ids, criteria)

	buf := new(bytes.Buffer)
	err := queryTemplate.Execute(buf, args)
	if err != nil {
		return "", fmt.Errorf("failed to format network query: %s", err)
	}
	return buf.String(), nil
}

func getNetworkQueryArgs(ids []string, criteria NetworkLoadCriteria) networkQueryArgs {
	ret := networkQueryArgs{
		TableName: networksTable,
		Fields:    fmt.Sprintf("%s.id", networksTable),
		IdList:    sql_utils.GetPlaceholderArgList(1, len(ids)),
	}

	if criteria.LoadMetadata {
		ret.Fields += fmt.Sprintf(", %s.name, %s.description", networksTable, networksTable)
	}

	if criteria.LoadConfigs {
		ret.Fields += fmt.Sprintf(", %s.type, %s.value", networkConfigTable, networkConfigTable)
		ret.JoinQuery = fmt.Sprintf("LEFT JOIN %s ON %s.network_id = %s.id", networkConfigTable, networkConfigTable, networksTable)
	}

	ret.Fields += fmt.Sprintf(", %s.version", networksTable)

	return ret
}

func scanNextNetworkRow(rows *sql.Rows, criteria NetworkLoadCriteria) (Network, error) {
	var id string
	var name, description sql.NullString
	var cfgType sql.NullString
	var cfgValue []byte
	var version uint64

	scanArgs := []interface{}{
		&id,
	}
	if criteria.LoadMetadata {
		scanArgs = append(scanArgs, &name, &description)
	}
	if criteria.LoadConfigs {
		scanArgs = append(scanArgs, &cfgType, &cfgValue)
	}
	scanArgs = append(scanArgs, &version)

	err := rows.Scan(scanArgs...)
	if err != nil {
		return Network{}, fmt.Errorf("error while scanning network row: %s", err)
	}

	ret := Network{ID: id, Name: nullStringToValue(name), Description: nullStringToValue(description), Configs: map[string][]byte{}, Version: version}
	if criteria.LoadConfigs && cfgType.Valid {
		ret.Configs[cfgType.String] = cfgValue
	}
	return ret, nil
}

func getNetworkIDsNotFound(networksByID map[string]*Network, queriedIDs []string) []string {
	ret := []string{}
	for _, id := range queriedIDs {
		if _, ok := networksByID[id]; !ok {
			ret = append(ret, id)
		}
	}
	sort.Strings(ret)
	return ret
}

func (store *sqlConfiguratorStorage) doesNetworkExist(id string) (bool, error) {
	query := fmt.Sprintf("SELECT count(1) FROM %s WHERE id = $1", networksTable)
	row := store.tx.QueryRow(query, id)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if network id %s exists: %s", id, err)
	}

	return count > 0, nil
}

func validateNetworkUpdates(updates []NetworkUpdateCriteria) error {
	updatesByID := funk.ToMap(updates, "ID").(map[string]NetworkUpdateCriteria)
	if len(updatesByID) < len(updates) {
		return errors.New("multiple updates for a single network are not allowed")
	}
	return nil
}

func (store *sqlConfiguratorStorage) updateNetwork(update NetworkUpdateCriteria, upsertConfigStmt *sql.Stmt, deleteConfigStmt *sql.Stmt) error {
	updNetworkQuery, updNetworkArgs, err := getNetworkUpdateExec(update)
	if err != nil {
		return err
	}

	_, err = store.tx.Exec(updNetworkQuery, updNetworkArgs...)
	if err != nil {
		return fmt.Errorf("error updating network %s: %s", update.ID, err)
	}

	// Sort config keys for deterministic behavior
	configUpdateTypes := funk.Keys(update.ConfigsToAddOrUpdate).([]string)
	sort.Strings(configUpdateTypes)
	for _, configType := range configUpdateTypes {
		configValue := update.ConfigsToAddOrUpdate[configType]
		_, err := upsertConfigStmt.Exec(update.ID, configType, configValue, configValue)
		if err != nil {
			return fmt.Errorf("error updating config %s on network %s: %s", configType, update.ID, err)
		}
	}
	for _, configType := range update.ConfigsToDelete {
		_, err := deleteConfigStmt.Exec(update.ID, configType)
		if err != nil {
			return fmt.Errorf("error deleting config %s on network %s: %s", configType, update.ID, err)
		}
	}

	return nil
}

func getNetworkUpdateExec(update NetworkUpdateCriteria) (string, []interface{}, error) {
	// UPDATE cfg_networks SET (name, description, version) = ($1, $2, cfg_networks.version + 1) WHERE id = $3
	networkExecTemplate := template.Must(template.New("nw_upd_exec").Parse(
		"UPDATE {{.TableName}} SET {{.Fields}} = {{.FieldsPlaceholder}} " +
			"WHERE {{.TableName}}.id = {{.IDPlaceholder}}",
	))
	templateArgs := getNetworkExecTemplateArgs(update)

	buf := new(bytes.Buffer)
	err := networkExecTemplate.Execute(buf, templateArgs)
	if err != nil {
		return "", []interface{}{}, fmt.Errorf("failed to format network update query: %s", err)
	}
	return buf.String(), getNetworkExecQueryArgs(update), nil
}

type networkExecTemplateArgs struct {
	TableName, Fields, FieldsPlaceholder, IDPlaceholder string
}

func getNetworkExecTemplateArgs(update NetworkUpdateCriteria) networkExecTemplateArgs {
	ret := networkExecTemplateArgs{
		TableName: networksTable,
	}

	fields := []string{}
	if update.NewName != nil {
		fields = append(fields, "name")
	}
	if update.NewDescription != nil {
		fields = append(fields, "description")
	}
	fields = append(fields, "version")

	ret.Fields = fmt.Sprintf("(%s)", strings.Join(fields, ", "))
	ret.FieldsPlaceholder = sql_utils.GetPlaceholderArgListWithSuffix(
		1,
		len(fields)-1,
		fmt.Sprintf("%s.version + 1", networksTable),
	)
	ret.IDPlaceholder = fmt.Sprintf("$%d", len(fields)+1)

	return ret
}

func getNetworkExecQueryArgs(update NetworkUpdateCriteria) []interface{} {
	ret := []interface{}{}
	if update.NewName != nil {
		ret = append(ret, stringPtrToVal(update.NewName))
	}
	if update.NewDescription != nil {
		ret = append(ret, stringPtrToVal(update.NewDescription))
	}
	ret = append(ret, update.ID)
	return ret
}

func stringPtrToVal(in *string) interface{} {
	if *in == "" {
		return nil
	}
	return *in
}

func nullStringToValue(in sql.NullString) string {
	if in.Valid {
		return in.String
	}
	return ""
}
