/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer

import (
	"fmt"
	"sync"
	"testing"

	merrors "magma/orc8r/lib/go/errors"
)

type indexerRegistry struct {
	sync.RWMutex
	indexers map[string]Indexer
}

// registry is a singleton providing process-global access to registered indexers.
var registry = &indexerRegistry{
	indexers: map[string]Indexer{},
}

// RegisterAll registers indexers with the state service to be called
// on updates to synced state.
func RegisterAll(indexers ...Indexer) error {
	registry.Lock()
	defer registry.Unlock()
	for i, indexer := range indexers {
		if err := registerUnsafe(indexer); err != nil {
			unregisterUnsafe(indexers[:i])
			return err
		}
	}
	return nil
}

// GetIndexer returns the registered indexer with ID.
// If not found, returns ErrNotFound from magma/orc8r/lib/go/errors.
func GetIndexer(id string) (Indexer, error) {
	registry.Lock()
	defer registry.Unlock()

	indexer, exists := registry.indexers[id]
	if !exists {
		return nil, merrors.ErrNotFound
	}
	return indexer, nil
}

// GetAllIndexers returns all registered indexers.
func GetAllIndexers() []Indexer {
	registry.Lock()
	defer registry.Unlock()

	indexers := make([]Indexer, 0, len(registry.indexers))
	for _, indexer := range registry.indexers {
		indexers = append(indexers, indexer)
	}

	return indexers
}

// GetAllIndexerVersionsByID returns a map of registered indexer IDs to their registered ("desired") versions.
func GetAllIndexerVersionsByID() map[string]Version {
	registry.Lock()
	defer registry.Unlock()

	versions := make(map[string]Version, len(registry.indexers))
	for _, indexer := range registry.indexers {
		versions[indexer.GetID()] = indexer.GetVersion()
	}

	return versions
}

func registerUnsafe(indexer Indexer) error {
	id := indexer.GetID()
	if _, exists := registry.indexers[id]; exists {
		return fmt.Errorf("an indexer with the ID %s already exists", id)
	}
	registry.indexers[id] = indexer
	return nil
}

func unregisterUnsafe(indexers []Indexer) {
	for _, indexer := range indexers {
		delete(registry.indexers, indexer.GetID())
	}
}

// DeregisterAllForTest deregisters all previously-registered indexers.
// This should only be called by test code.
func DeregisterAllForTest(t *testing.T) {
	if t == nil {
		panic("for tests only")
	}
	registry.indexers = map[string]Indexer{}
}

// RegisterForTest sets an indexer in the registry.
// Overwrites any existing indexer with the same indexer ID.
// This should only be called by test code.
func RegisterForTest(t *testing.T, idx Indexer) {
	if t == nil {
		panic("for tests only")
	}
	registry.indexers[idx.GetID()] = idx
}
