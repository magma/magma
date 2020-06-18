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

	"github.com/thoas/go-funk"
)

type remoteRegistry struct {
	sync.RWMutex
	// TODO(hcgatewood): convert to RemoteIndexer-only once supercontainer is broken up
	indexerByService map[string]Indexer
}

var (
	// reg is a singleton providing process-global access to registered indexers.
	reg = &remoteRegistry{indexerByService: map[string]Indexer{}}
)

// RegisterIndexers registers an indexer.
func RegisterIndexers(indexers ...Indexer) error {
	reg.Lock()
	defer reg.Unlock()

	for _, x := range indexers {
		id := x.GetID()
		if _, exists := reg.indexerByService[id]; exists {
			return fmt.Errorf("an indexer for the service %s already exists", id)
		}
		reg.indexerByService[id] = x
	}
	return nil
}

// GetIndexer returns the remote indexer for a desired service.
// Returns nil if not found.
func GetIndexer(serviceName string) Indexer {
	reg.Lock()
	defer reg.Unlock()

	return reg.indexerByService[serviceName]
}

// GetIndexers returns all registered indexers.
func GetIndexers() []Indexer {
	reg.Lock()
	defer reg.Unlock()

	var ret []Indexer
	for _, x := range reg.indexerByService {
		ret = append(ret, x)
	}
	return ret
}

// GetIndexersForState returns all registered indexers which handle the passed state type.
func GetIndexersForState(stateType string) []Indexer {
	reg.Lock()
	defer reg.Unlock()

	var filtered []Indexer
	for _, x := range reg.indexerByService {
		if funk.Contains(x.GetTypes(), stateType) {
			filtered = append(filtered, x)
		}

	}
	return filtered
}

// DeregisterAllForTest deregisters all previously-registered indexers.
// This should only be called by test code.
func DeregisterAllForTest(t *testing.T) {
	if t == nil {
		panic("for tests only")
	}
	reg.indexerByService = map[string]Indexer{}
}

// RegisterForTest sets an indexer in the registry.
// Overwrites any existing indexer with the same indexer ID.
// This should only be called by test code.
func RegisterForTest(t *testing.T, idx Indexer) {
	if t == nil {
		panic("for tests only")
	}
	reg.indexerByService[idx.GetID()] = idx
}
