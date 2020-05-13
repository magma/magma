/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metricsd

import (
	"fmt"
	"sync"
)

type mpRegistry struct {
	sync.Mutex
	profiles map[string]MetricsProfile
}

// registry is a global registry of available MetricsProfiles. This will be
// populated by MagmaPlugins eventually.
var registry = &mpRegistry{profiles: map[string]MetricsProfile{}}

// RegisterMetricsProfiles registers a collection of MetricsProfiles with
// metricsd. If any profile fails to register, changes will be rolled back.
// This function is thread-safe.
func RegisterMetricsProfiles(profiles ...MetricsProfile) error {
	registry.Lock()
	defer registry.Unlock()
	for i, profile := range profiles {
		if err := registerUnsafe(profile); err != nil {
			unregisterUnsafe(profiles[:i])
			return err
		}
	}
	return nil
}

func registerUnsafe(profile MetricsProfile) error {
	if _, nameExists := registry.profiles[profile.Name]; nameExists {
		return fmt.Errorf("A metrics profile with the name %s already exists", profile.Name)
	}
	registry.profiles[profile.Name] = profile
	return nil
}

func unregisterUnsafe(profiles []MetricsProfile) {
	for _, profile := range profiles {
		delete(registry.profiles, profile.Name)
	}
}

// GetMetricsProfile returns a registered MetricsProfile by name. Will return
// an error if no profile with the given name is found.
func GetMetricsProfile(name string) (MetricsProfile, error) {
	registry.Lock()
	defer registry.Unlock()

	profile, exists := registry.profiles[name]
	if !exists {
		return MetricsProfile{}, fmt.Errorf("no metrics profile with name %s is registered", name)
	}
	return profile, nil
}
