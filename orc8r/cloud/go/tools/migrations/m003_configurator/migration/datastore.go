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
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func UnmarshalJSONPBProtosFromDatastore(rows *sql.Rows, msgInstance proto.Message) (map[string]proto.Message, error) {
	ret := map[string]proto.Message{}
	for rows.Next() {
		k, v, err := scanNextDatastoreRow(rows)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		newMsg := reflect.New(reflect.TypeOf(msgInstance).Elem()).Interface().(proto.Message)
		err = Unmarshal(v, newMsg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal datastore row %s:\n%s", k, string(v))
		}
		ret[k] = newMsg
	}
	return ret, nil
}

func UnmarshalProtoMessagesFromDatastore(rows *sql.Rows, msgInstance proto.Message) (map[string]proto.Message, error) {
	ret := map[string]proto.Message{}
	for rows.Next() {
		k, v, err := scanNextDatastoreRow(rows)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		newMsg := reflect.New(reflect.TypeOf(msgInstance).Elem()).Interface().(proto.Message)
		err = proto.Unmarshal(v, newMsg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal datastore row %s:\n%s", k, string(v))
		}
		ret[k] = newMsg
	}
	return ret, nil
}

func scanNextDatastoreRow(rows *sql.Rows) (string, []byte, error) {
	var k string
	var v []byte

	err := rows.Scan(&k, &v)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to scan datastore row")
	}
	return k, v, nil
}
