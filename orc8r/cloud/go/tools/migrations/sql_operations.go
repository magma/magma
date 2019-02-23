/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package migrations

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/golang/glog"
)

func GetTableName(networkId string, baseName string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(networkId), baseName)
}

// If the table DNE, log and return empty map
func GetAllValuesFromTable(tx *sql.Tx, table string) (map[string][]byte, error) {
	// Not every network may have gateways or meshes, in which case the
	// corresponding tables won't exist. Check and return early if so.
	exists, err := doesTableExist(tx, table)
	if err != nil {
		return nil, fmt.Errorf("Error checking if table %s exists: %s", table, err)
	}
	if !exists {
		glog.Errorf("Table %s does not exist, returning empty result from get", table)
		return map[string][]byte{}, nil
	}

	query := fmt.Sprintf("SELECT key, value FROM %s", table)
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	ret := map[string][]byte{}
	for rows.Next() {
		var key string
		var val []byte

		err = rows.Scan(&key, &val)
		if err != nil {
			return nil, err
		}
		ret[key] = val
	}
	return ret, nil
}

// IMPORTANT: This is NOT portable, and ONLY works on postgres!
func doesTableExist(tx *sql.Tx, table string) (bool, error) {
	row := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name=$1)", table)
	ret := false
	err := row.Scan(&ret)
	return ret, err
}

func GetAllKeysFromTable(tx *sql.Tx, table string) ([]string, error) {
	query := fmt.Sprintf("SELECT key FROM %s", table)
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var ret []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, err
		}
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret, nil
}
