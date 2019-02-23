/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Datastore provides a key-value pair interface for the cloud services
package datastore

import "errors"

type ValueWrapper struct {
	Value      []byte
	Generation uint64
}

// ErrNotFound is returned by GET when no record for the key is found
var ErrNotFound = errors.New("No record for query")

type Api interface {
	Put(table string, key string, value []byte) error
	PutMany(table string, valuesToPut map[string][]byte) (map[string]error, error)
	Get(table string, key string) ([]byte, uint64, error)
	GetMany(table string, keys []string) (map[string]ValueWrapper, error)
	Delete(table string, key string) error
	DeleteMany(table string, keys []string) (map[string]error, error)
	ListKeys(table string) ([]string, error)
	DeleteTable(table string) error
	DoesKeyExist(table string, key string) (bool, error)
}
