/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package datastore

import "sync"

type SyncStore struct {
	store Api
	lock  *sync.Mutex
}

func NewSyncStore(store Api) *SyncStore {
	return &SyncStore{store, &sync.Mutex{}}
}

func (s *SyncStore) Put(table string, key string, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store.Put(table, key, value)
}

func (s *SyncStore) PutMany(table string, valuesToPut map[string][]byte) (map[string]error, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store.PutMany(table, valuesToPut)
}

func (s *SyncStore) Get(table string, key string) ([]byte, uint64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store.Get(table, key)
}

func (s *SyncStore) GetMany(table string, keys []string) (map[string]ValueWrapper, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store.GetMany(table, keys)
}

func (s *SyncStore) Delete(table string, key string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store.Delete(table, key)
}

func (s *SyncStore) DeleteMany(table string, keys []string) (map[string]error, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store.DeleteMany(table, keys)
}

func (s *SyncStore) ListKeys(table string) ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store.ListKeys(table)
}

func (s *SyncStore) DeleteTable(table string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.store.DeleteTable(table)
}
