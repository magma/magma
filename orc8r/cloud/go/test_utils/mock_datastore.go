/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"sync"

	"magma/orc8r/cloud/go/datastore"
)

type mockDatastoreTable map[string][]byte

// Datastore backed by a golang map
type MockDatastore struct {
	store map[string]mockDatastoreTable
}

var instance *MockDatastore
var once sync.Once

// Get a singleton mock datastore instance for tests that require multiple
// test services
func GetMockDatastoreInstance() *MockDatastore {
	once.Do(func() {
		instance = NewMockDatastore()
	})
	return instance
}

func NewMockDatastore() *MockDatastore {
	ds := new(MockDatastore)
	ds.store = make(map[string]mockDatastoreTable, 0)
	return ds
}

func (m *MockDatastore) initTable(table string) {
	if _, ok := m.store[table]; !ok {
		m.store[table] = make(map[string][]byte, 0)
	}
}

func (m *MockDatastore) Put(table string, key string, value []byte) error {
	m.initTable(table)
	m.store[table][key] = value
	return nil
}

func (m *MockDatastore) PutMany(table string, valuesToPut map[string][]byte) (map[string]error, error) {
	m.initTable(table)
	for k, v := range valuesToPut {
		m.store[table][k] = v
	}
	return map[string]error{}, nil
}

func (m *MockDatastore) Get(table string, key string) ([]byte, uint64, error) {
	m.initTable(table)
	value, ok := m.store[table][key]
	if ok {
		return value, 0, nil
	}
	return nil, 0, datastore.ErrNotFound
}

func (m *MockDatastore) GetMany(table string, keys []string) (map[string]datastore.ValueWrapper, error) {
	m.initTable(table)
	ret := make(map[string]datastore.ValueWrapper, len(keys))
	for _, k := range keys {
		val, ok := m.store[table][k]
		if ok {
			ret[k] = datastore.ValueWrapper{
				Value:      val,
				Generation: 0,
			}
		}
	}
	return ret, nil
}

func (m *MockDatastore) Delete(table string, key string) error {
	m.initTable(table)

	delete(m.store[table], key)
	return nil
}

func (m *MockDatastore) DeleteMany(table string, keys []string) (map[string]error, error) {
	m.initTable(table)
	for _, k := range keys {
		delete(m.store[table], k)
	}
	return map[string]error{}, nil
}

func (m *MockDatastore) ListKeys(table string) ([]string, error) {
	m.initTable(table)
	keys := make([]string, 0, len(m.store[table]))
	for key := range m.store[table] {
		keys = append(keys, key)
	}
	return keys, nil
}

func (m *MockDatastore) DeleteTable(table string) error {
	m.initTable(table)
	delete(m.store, table)
	return nil
}

func (m *MockDatastore) DoesKeyExist(table string, key string) (bool, error) {
	m.initTable(table)
	_, ok := m.store[table][key]
	if ok {
		return true, nil
	} else {
		return false, nil
	}
}
