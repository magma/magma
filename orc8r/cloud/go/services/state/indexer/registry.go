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
)

type indexerRegistry struct {
	sync.RWMutex
	indexers map[string]Indexer
}

var registry = &indexerRegistry{
	indexers: map[string]Indexer{},
}

// RegisterIndexers registers Indexers with the state service to be called
// on updates to synced state.
func RegisterIndexers(indexers ...Indexer) error {
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
