/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package migration

import (
	"reflect"

	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func UnmarshalJSONPBProtosFromDatastore(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string, table string, msgInstance proto.Message) (map[string]proto.Message, error) {
	return unmarshalProtosFromDatastore(sc, builder, networkID, table, msgInstance, Unmarshal)
}

func UnmarshalProtoMessagesFromDatastore(sc *squirrel.StmtCache, builder sqorc.StatementBuilder, networkID string, table string, msgInstance proto.Message) (map[string]proto.Message, error) {
	return unmarshalProtosFromDatastore(sc, builder, networkID, table, msgInstance, proto.Unmarshal)
}

func unmarshalProtosFromDatastore(
	sc *squirrel.StmtCache,
	builder sqorc.StatementBuilder,
	networkID string,
	table string,
	msgInstance proto.Message,
	unmarshaler func([]byte, proto.Message) error,
) (map[string]proto.Message, error) {
	_, err := builder.CreateTable(GetLegacyTableName(networkID, table)).
		IfNotExists().
		Column(DatastoreKeyCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(DatastoreValCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(DatastoreGenCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		Column(DatastoreDeletedCol).Type(sqorc.ColumnTypeBool).NotNull().Default("FALSE").EndColumn().
		RunWith(sc).
		Exec()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to provisionally create %s table for network %s", table, networkID)
	}

	rows, err := builder.Select(DatastoreKeyCol, DatastoreValCol).
		From(GetLegacyTableName(networkID, table)).
		RunWith(sc).
		Query()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query %s for network %s", table, networkID)
	}
	defer rows.Close()

	ret := map[string]proto.Message{}
	for rows.Next() {
		var k string
		var v []byte
		err := rows.Scan(&k, &v)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan datastore row")
		}

		newMsg := reflect.New(reflect.TypeOf(msgInstance).Elem()).Interface().(proto.Message)
		err = unmarshaler(v, newMsg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal datastore row %s:\n%s", k, string(v))
		}
		ret[k] = newMsg
	}
	return ret, nil
}
