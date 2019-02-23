/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package blacklist

// These are config types which have since been deleted from the system. When
// we see these being read from a snapshot, ignore them.
// Ideally, we can run some cleanup migration to purge these from all the
// config tables at some point.
var blacklistedConfigTypes = map[string]struct{}{
	"dnsd_gateway":       {},
	"magmad_network":     {},
	"terrgraph_network":  {},
	"terrgraph_gateway":  {},
	"federation_gateway": {},
}

// IsConfigBlacklisted returns whether this is a config type that has been deleted from the system
func IsConfigBlacklisted(configType string) bool {
	_, isBlacklisted := blacklistedConfigTypes[configType]
	return isBlacklisted
}
