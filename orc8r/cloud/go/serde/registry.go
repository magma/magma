/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package serde

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

type serdeRegistry struct {
	sync.RWMutex
	serdeRegistriesByDomain map[string]*serdes
}

var registry = &serdeRegistry{serdeRegistriesByDomain: map[string]*serdes{}}

var (
	// missingDomainsCalculatedCallback is an unexported hook for testing
	missingDomainsCalculatedCallback = func() {}

	// newDomainsCreatedCallback is an unexported hook for testing
	newDomainsCreatedCallback = func() {}
)

// RegisterSerdes will register a collection of Serde implementations with
// the global Serde registry. The semantics are all-or-nothing: if an error is
// encountered while registering any Serde, the registry will rollback all
// changes made. This function is thread-safe.
func RegisterSerdes(serdesToRegister ...Serde) error {
	serdesByDomain := getSerdesByDomain(serdesToRegister)
	missingDomains := getNewDomainsToCreate(serdesByDomain)
	missingDomainsCalculatedCallback()
	// missingDomains will always be a superset of the changes that we actually
	// need to write since there's no public API to delete a domain.
	// So we don't need to recalculate it, we can just acquire the write lock
	// and filter out the domains which were created in the meantime.
	if len(missingDomains) > 0 {
		createNewDomains(serdesByDomain)
	}
	newDomainsCreatedCallback()

	// Read lock here because we're not modifying the top-level registry
	registry.RLock()
	defer registry.RUnlock()

	domains := getSortedSerdeDomainKeys(serdesByDomain)
	for i, domain := range domains {
		subregistry := registry.serdeRegistriesByDomain[domain]
		err := subregistry.register(serdesByDomain[domain])
		if err != nil {
			// :i because the current subregistry will rollback on error
			for _, rollbackDomain := range domains[:i] {
				subregistry.unregister(serdesByDomain[rollbackDomain])
			}
			return fmt.Errorf("Error registering serdes: %s; registry has been rolled back", err)
		}
	}
	return nil
}

// Serialize serializes an object (`data`) by delegating to the appropriate
// Serde identified by the domain and typeVal. This function is thread-safe.
func Serialize(domain string, typeVal string, data interface{}) ([]byte, error) {
	registry.RLock()
	defer registry.RUnlock()
	subregistry, ok := registry.serdeRegistriesByDomain[domain]
	if !ok {
		return []byte{}, fmt.Errorf("No serdes registered for domain %s", domain)
	}
	return subregistry.serialize(typeVal, data)
}

// Deserialize deserializes a bytearray by delegating to the appropriate Serde
// identified by the domain and typeVal. This function is thread-safe.
// If the data parameter is nil or empty, this function will return nil, nil.
func Deserialize(domain string, typeVal string, data []byte) (interface{}, error) {
	if data == nil || len(data) == 0 {
		return nil, nil
	}

	registry.RLock()
	defer registry.RUnlock()
	subregistry, ok := registry.serdeRegistriesByDomain[domain]
	if !ok {
		return []byte{}, fmt.Errorf("No serdes registered for domain %s", domain)
	}
	return subregistry.deserialize(typeVal, data)
}

func getSerdesByDomain(serdesToGroup []Serde) map[string][]Serde {
	ret := map[string][]Serde{}
	for _, s := range serdesToGroup {
		domain := s.GetDomain()
		if _, ok := ret[domain]; !ok {
			ret[domain] = []Serde{}
		}
		ret[domain] = append(ret[domain], s)
	}
	return ret
}

func getNewDomainsToCreate(serdesByDomain map[string][]Serde) []string {
	registry.RLock()
	defer registry.RUnlock()

	var ret []string
	for domain := range serdesByDomain {
		if _, ok := registry.serdeRegistriesByDomain[domain]; !ok {
			ret = append(ret, domain)
		}
	}
	return ret
}

func createNewDomains(serdesByDomain map[string][]Serde) {
	registry.Lock()
	defer registry.Unlock()

	for domain := range serdesByDomain {
		// Check if we need to create a new entry in the registry because this
		// could have changed between releasing the read lock and acquiring
		// the write lock
		if _, ok := registry.serdeRegistriesByDomain[domain]; ok {
			continue
		}
		registry.serdeRegistriesByDomain[domain] = &serdes{serdesByKey: map[string]Serde{}}
	}
}

func getSortedSerdeDomainKeys(serdesByDomain map[string][]Serde) []string {
	ret := make([]string, 0, len(serdesByDomain))
	for k := range serdesByDomain {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

// serdes is a registry of Serdes for a single domain
type serdes struct {
	sync.RWMutex
	serdesByKey map[string]Serde
}

func (s *serdes) register(serdesToRegister []Serde) error {
	s.Lock()
	defer s.Unlock()

	for i, serde := range serdesToRegister {
		if err := s.registerUnsafe(serde); err != nil {
			s.unregisterUnsafe(serdesToRegister[:i])
			return err
		}
	}
	return nil
}

func (s *serdes) unregister(serdesToUnregister []Serde) {
	s.Lock()
	defer s.Unlock()
	s.unregisterUnsafe(serdesToUnregister)
}

func (s *serdes) registerUnsafe(serde Serde) error {
	if _, ok := s.serdesByKey[serde.GetType()]; ok {
		return fmt.Errorf("Serde with key %s is already registered", serde.GetType())
	}
	s.serdesByKey[serde.GetType()] = serde
	return nil
}

func (s *serdes) unregisterUnsafe(serdesToUnregister []Serde) {
	for _, serde := range serdesToUnregister {
		delete(s.serdesByKey, serde.GetType())
	}
}

func (s *serdes) serialize(t string, data interface{}) ([]byte, error) {
	s.RLock()
	defer s.RUnlock()
	serde, err := s.getSerdeUnsafe(t)
	if err != nil {
		return nil, err
	}
	return serde.Serialize(data)
}

func (s *serdes) deserialize(t string, data []byte) (interface{}, error) {
	s.RLock()
	defer s.RUnlock()
	serde, err := s.getSerdeUnsafe(t)
	if err != nil {
		return nil, err
	}
	return serde.Deserialize(data)
}

func (s *serdes) getSerdeUnsafe(t string) (Serde, error) {
	serde, ok := s.serdesByKey[t]
	if !ok {
		return nil, fmt.Errorf("No Serde found for type %s", t)
	}
	return serde, nil
}

// UnregisterAllSerdes should only be used in test code!!!!!
// DO NOT USE IN ANYTHING OTHER THAN TESTS
func UnregisterAllSerdes(t *testing.T) {
	if t == nil {
		panic("Nice try")
	}
	registry.Lock()
	defer registry.Unlock()
	registry.serdeRegistriesByDomain = map[string]*serdes{}
}

// UnregisterSerdesForDomain should only be used in test code!!!!
// DO NOT USE IN ANYTHING OTHER THAN TESTS
func UnregisterSerdesForDomain(t *testing.T, domain string) {
	if t == nil {
		panic("Nice try")
	}
	registry.Lock()
	defer registry.Unlock()
	registry.serdeRegistriesByDomain[domain] = &serdes{serdesByKey: map[string]Serde{}}
}
