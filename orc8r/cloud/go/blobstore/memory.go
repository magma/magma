/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package blobstore

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"magma/orc8r/cloud/go/storage"
	magmaerrors "magma/orc8r/lib/go/errors"

	"github.com/thoas/go-funk"
)

type changeType int
type tNetworkID = string

const (
	CreateOrUpdate changeType = 1
	Delete         changeType = 2
)

type change struct {
	cType changeType
	blob  Blob
}

type changesByID map[storage.TypeAndKey]change

// transactionTable maps networkIDs to a map of updates
type transactionTable map[tNetworkID]changesByID

type blobsByID map[storage.TypeAndKey]Blob

// blobTable maps networkIDs to a map of Blobs
type blobTable map[tNetworkID]blobsByID

type keySet map[string]interface{}

type memoryBlobStorage struct {
	// guards both transactionExists flag and changes
	sync.RWMutex
	// transactionExists indicates whether there exists an on-going transaction
	transactionExists bool
	// changes stores changes during a transaction
	changes transactionTable

	// stores everything needed to access the shared map
	shared sharedMemoryBlobTables
}

type sharedMemoryBlobTables struct {
	*sync.RWMutex
	table blobTable
}

type memoryBlobStoreFactory struct {
	sync.RWMutex
	table blobTable
}

type networkIDAndTK struct {
	networkID  string
	typeAndKey storage.TypeAndKey
}

type nIDAndTKSet map[networkIDAndTK]interface{}

// NewMemoryBlobStorageFactory returns a BlobStorageFactory implementation
// which will return storage APIs backed by an in-memory map.
func NewMemoryBlobStorageFactory() BlobStorageFactory {
	return &memoryBlobStoreFactory{table: blobTable{}}
}

func (fact *memoryBlobStoreFactory) StartTransaction(opts *storage.TxOptions) (TransactionalBlobStorage, error) {
	return &memoryBlobStorage{
		shared:            sharedMemoryBlobTables{RWMutex: &fact.RWMutex, table: fact.table},
		transactionExists: true,
		changes:           transactionTable{}}, nil
}

func (fact *memoryBlobStoreFactory) InitializeFactory() error {
	return nil
}

func (store *memoryBlobStorage) Commit() error {
	store.Lock()
	defer store.Unlock()

	if err := store.validateTx(); err != nil {
		return err
	}

	store.shared.Lock()
	store.applyChangesToShared()
	store.shared.Unlock()

	store.resetTransaction()
	return nil
}

func (store *memoryBlobStorage) Rollback() error {
	store.Lock()
	defer store.Unlock()

	if err := store.validateTx(); err != nil {
		return err
	}
	store.resetTransaction()
	return nil
}

// ListKeys grabs keys from the shared map first, and then updates the keys
// with changes from the ongoing transaction
func (store *memoryBlobStorage) ListKeys(networkID string, typeVal string) ([]string, error) {
	store.RLock()
	defer store.RUnlock()

	if err := store.validateTx(); err != nil {
		return nil, err
	}

	store.shared.RLock()
	keySet := store.listKeysFromShared(networkID, typeVal)
	store.shared.RUnlock()

	return store.updateKeysWithLocalChangesUnsafe(networkID, typeVal, keySet)
}

func (store *memoryBlobStorage) Get(networkID string, id storage.TypeAndKey) (Blob, error) {
	multiRet, err := store.GetMany(networkID, []storage.TypeAndKey{id})
	if err != nil {
		return Blob{}, err
	}
	if len(multiRet) == 0 {
		return Blob{}, magmaerrors.ErrNotFound
	}
	return multiRet[0], nil
}

// GetMany grabs blobs corresponding to the ids from the shared map, then
// updates the blobs with changes from the ongoing transaction
func (store *memoryBlobStorage) GetMany(networkID string, ids []storage.TypeAndKey) ([]Blob, error) {
	store.RLock()
	defer store.RUnlock()

	if err := store.validateTx(); err != nil {
		return nil, err
	}

	store.shared.RLock()
	sharedBlobs := store.getManyFromShared(networkID, ids)
	store.shared.RUnlock()

	return store.updateBlobsWithLocalChangesUnsafe(networkID, ids, sharedBlobs)
}

func (store *memoryBlobStorage) CreateOrUpdate(networkID string, blobs []Blob) error {
	store.Lock()
	defer store.Unlock()

	if err := store.validateTx(); err != nil {
		return err
	}

	ids := blobsToIDs(blobs)
	store.shared.RLock()
	sharedBlobSet := store.getManyFromShared(networkID, ids)
	store.shared.RUnlock()

	// check shared first and grab existing versions
	for i, blob := range blobs {
		id := blob.toID()
		sharedBlob, ok := sharedBlobSet[id]
		// increment version if it isn't set in the update
		if ok && blob.Version == 0 {
			blobs[i].Version = sharedBlob.Version + 1
		}
	}

	store.changes.initializeNetworkTable(networkID)
	perNetworkLocalMap := store.changes[networkID]
	for _, blob := range blobs {
		id := blob.toID()
		storedChange, exists := perNetworkLocalMap[id]
		if exists && storedChange.cType == CreateOrUpdate && blob.Version == 0 {
			blob.Version = storedChange.blob.Version + 1
		}
		perNetworkLocalMap[id] = change{cType: CreateOrUpdate, blob: blob}
	}
	return nil
}

func (store *memoryBlobStorage) Delete(networkID string, ids []storage.TypeAndKey) error {
	store.Lock()
	defer store.Unlock()

	if err := store.validateTx(); err != nil {
		return err
	}

	store.changes.initializeNetworkTable(networkID)
	for _, id := range ids {
		store.changes[networkID][id] = change{cType: Delete}
	}
	return nil
}

func (store *memoryBlobStorage) GetExistingKeys(keys []string, filter SearchFilter) ([]string, error) {
	store.Lock()
	defer store.Unlock()

	if err := store.validateTx(); err != nil {
		return nil, err
	}
	keySet := funk.Map(keys, func(k string) (string, interface{}) { return k, nil }).(map[string]interface{})
	if funk.NotEmpty(filter.NetworkID) {
		return store.getExistingKeysInNetwork(*filter.NetworkID, keySet)
	}
	return store.getExistingKeysAllNetworks(keySet)
}

func (store *memoryBlobStorage) IncrementVersion(networkID string, id storage.TypeAndKey) error {
	store.Lock()
	defer store.Unlock()

	if err := store.validateTx(); err != nil {
		return err
	}

	blob := Blob{
		Type:    id.Type,
		Key:     id.Key,
		Version: 1,
	}

	store.shared.RLock()
	master, ok := store.shared.table[networkID]
	if ok {
		sharedBlob, ok := master[id]
		if ok {
			blob.Version = sharedBlob.Version + 1
			blob.Value = sharedBlob.Value
		}
	}
	store.shared.RUnlock()

	store.changes.initializeNetworkTable(networkID)
	perNetworkLocalMap := store.changes[networkID]
	storedChange, exists := perNetworkLocalMap[id]
	if exists && storedChange.cType == CreateOrUpdate {
		blob.Version = storedChange.blob.Version + 1
		blob.Value = storedChange.blob.Value
	}
	perNetworkLocalMap[id] = change{cType: CreateOrUpdate, blob: blob}

	return nil
}

func (store *memoryBlobStorage) getExistingKeysInNetwork(networkID string, keySet keySet) ([]string, error) {
	store.shared.RLock()
	_, ok := store.shared.table[networkID]
	if !ok {
		store.shared.RUnlock()
		return nil, fmt.Errorf("Network %s does not exist", networkID)
	}
	existingKeysSet := nIDAndTKSet{}
	store.getExistingKeysAllNetworksFromShared(networkID, keySet, existingKeysSet)
	store.shared.RUnlock()

	_, ok = store.changes[networkID]
	if ok {
		store.updateSearchedKeysWithLocalChanges(networkID, keySet, existingKeysSet)
	}
	return existingKeysSet.sortAndRemoveDuplicate(), nil
}

func (store *memoryBlobStorage) getExistingKeysAllNetworks(keys keySet) ([]string, error) {
	store.shared.RLock()
	existingKeysSet := store.getExistingKeysFromShared(keys)
	store.shared.RUnlock()

	for networkID := range store.changes {
		store.updateSearchedKeysWithLocalChanges(networkID, keys, existingKeysSet)
	}
	return existingKeysSet.sortAndRemoveDuplicate(), nil
}

func (store *memoryBlobStorage) getExistingKeysFromShared(keySet keySet) nIDAndTKSet {
	existingKeysSet := nIDAndTKSet{}
	for networkID := range store.shared.table {
		store.getExistingKeysAllNetworksFromShared(networkID, keySet, existingKeysSet)
	}
	return existingKeysSet
}

// existingKeysSet is also an outputting parameter
func (store *memoryBlobStorage) getExistingKeysAllNetworksFromShared(networkID string, keySet keySet, existingKeysSet nIDAndTKSet) {
	for tk := range store.shared.table[networkID] {
		if _, exists := keySet[tk.Key]; exists {
			existingKeysSet[networkIDAndTK{networkID: networkID, typeAndKey: tk}] = nil
		}
	}
}

// existingKeysSet is also an outputting parameter
func (store *memoryBlobStorage) updateSearchedKeysWithLocalChanges(networkID string, keySet keySet, existingKeysSet nIDAndTKSet) error {
	for tk, change := range store.changes[networkID] {
		if _, exists := keySet[tk.Key]; exists {
			id := networkIDAndTK{networkID: networkID, typeAndKey: tk}
			switch change.cType {
			case Delete:
				delete(existingKeysSet, id)
			case CreateOrUpdate:
				existingKeysSet[id] = nil
			default:
				return fmt.Errorf("This transaction contains ill-formatted changes.")
			}
		}
	}
	return nil
}

func (set *nIDAndTKSet) sortAndRemoveDuplicate() []string {
	deduped := funk.Map(*set, func(id networkIDAndTK, _ interface{}) (string, interface{}) { return id.typeAndKey.Key, nil }).(map[string]interface{})
	keys := funk.Map(deduped, func(key string, _ interface{}) string { return key }).([]string)
	sort.Strings(keys)
	return keys
}

// Must be called with read lock on change map.
func (store *memoryBlobStorage) validateTx() error {
	if store.transactionExists == false {
		return errors.New("No transaction is available")
	}
	return nil
}

// Traverse through the changes from the transaction and put them into the
// shared map. Must be called with write lock on both local and shared maps.
func (store *memoryBlobStorage) applyChangesToShared() error {
	fact := store.shared.table
	for networkID, perNetworkChangeMap := range store.changes {
		for id, change := range perNetworkChangeMap {
			switch change.cType {
			case Delete:
				delete(fact[networkID], id)
			case CreateOrUpdate:
				fact.initializeNetworkTable(networkID)
				fact[networkID][id] = change.blob
			default:
				return fmt.Errorf("This transcaction contains ill-formatted changes.")
			}
		}
	}
	return nil
}

// Must be called with write lock on change map.
func (store *memoryBlobStorage) resetTransaction() {
	store.transactionExists = false
	store.changes = nil
}

// Given a networkID and a type this function looks in the shared map and
// returns a set of keys that match the given type. Must be called with read
// lock on shared map.
func (store *memoryBlobStorage) listKeysFromShared(networkID string, typeVal string) map[string]struct{} {
	keySet := map[string]struct{}{}

	table, ok := store.shared.table[networkID]
	if !ok {
		return keySet
	}
	for _, blob := range table {
		if blob.Type == typeVal {
			keySet[blob.Key] = struct{}{}
		}
	}
	return keySet
}

// Given a networkID, a type, and a map of keys found from the shared map, this
// function looks through the local map of changes and applies them onto the
// keys. Must be called with lock on local map.
func (store *memoryBlobStorage) updateKeysWithLocalChangesUnsafe(networkID string, typeToQuery string, keySetFromShared map[string]struct{}) ([]string, error) {
	networkMap, ok := store.changes[networkID]
	if !ok {
		return fromKeySet(keySetFromShared), nil
	}

	for id, change := range networkMap {
		if id.Type == typeToQuery {
			switch change.cType {
			case Delete:
				delete(keySetFromShared, id.Key)
			case CreateOrUpdate:
				keySetFromShared[id.Key] = struct{}{}
			default:
				return nil, fmt.Errorf("This transcaction contains ill-formatted changes.")
			}
		}
	}
	return fromKeySet(keySetFromShared), nil
}

// Given a networkID and a list of ids this function looks in the shared map
// and returns a map of id:blob that match the given ids. Must be called with
// read lock on the shared table.
func (store *memoryBlobStorage) getManyFromShared(networkID string, ids []storage.TypeAndKey) blobsByID {
	blobSet := blobsByID{}

	master, ok := store.shared.table[networkID]
	if !ok {
		return blobSet
	}

	for _, id := range ids {
		blob, ok := master[id]
		if ok {
			blobSet[id] = blob
		}
	}
	return blobSet
}

// Given a networkID, a list of ids, and a map of id:blob gathered from
// getManyFromShared, this function looks through items in the local map that
// match the given ids and applies the changes onto the blobs. This function
// returns a list of blobs from the modified map.
// Must be called with read lock on change map.
func (store *memoryBlobStorage) updateBlobsWithLocalChangesUnsafe(networkID string, idsToQuery []storage.TypeAndKey, blobsByID blobsByID) ([]Blob, error) {
	networkMap, existsInLocal := store.changes[networkID]
	if !existsInLocal {
		return blobsByID.toBlobList(), nil
	}

	for _, id := range idsToQuery {
		change, exists := networkMap[id]
		if !exists {
			continue
		}
		switch {
		case change.cType == Delete:
			delete(blobsByID, id)
		case change.cType == CreateOrUpdate:
			blobsByID[id] = change.blob
		default:
			return nil, fmt.Errorf("This transcaction contains ill-formatted changes.")
		}
	}
	return blobsByID.toBlobList(), nil
}

// Adds a field if it doesn't exist already.
func (table blobTable) initializeNetworkTable(networkID tNetworkID) {
	if _, ok := table[networkID]; !ok {
		table[networkID] = blobsByID{}
	}
}

// Adds a field if it doesn't exist already.
func (table transactionTable) initializeNetworkTable(networkID tNetworkID) {
	if _, ok := table[networkID]; !ok {
		table[networkID] = changesByID{}
	}
}

func (blob *Blob) toID() storage.TypeAndKey {
	return storage.TypeAndKey{Type: blob.Type, Key: blob.Key}
}

func blobsToIDs(blobs []Blob) []storage.TypeAndKey {
	ids := []storage.TypeAndKey{}
	for _, blob := range blobs {
		ids = append(ids, blob.toID())
	}
	return ids
}

func (blobSet blobsByID) toBlobList() []Blob {
	blobs := []Blob{}
	for _, blob := range blobSet {
		blobs = append(blobs, blob)
	}
	return blobs
}

func fromKeySet(keySet map[string]struct{}) []string {
	keys := []string{}
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
