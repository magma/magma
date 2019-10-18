/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package access

import (
	"strings"

	"magma/orc8r/cloud/go/obsidian"

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/protos"
)

type RequestIdentityFinder func(c echo.Context) []*protos.Identity

type finderRegistryType map[string]struct {
	finder       RequestIdentityFinder
	identityRoot string
}

var finderRegistries = map[string]finderRegistryType{
	obsidian.V0: makeFinderRegistry(obsidian.V0),
	obsidian.V1: makeFinderRegistry(obsidian.V1),
}

func makeFinderRegistry(version string) finderRegistryType {
	networkRoot := makeVersionedRoot(version, obsidian.MagmaNetworksUrlPart)
	operatorRoot := makeVersionedRoot(version, obsidian.MagmaOperatorsUrlPart)
	return finderRegistryType{
		obsidian.MagmaNetworksUrlPart: {
			finder:       func(c echo.Context) []*protos.Identity { return getNetworkIdentity(c, version, networkRoot) },
			identityRoot: networkRoot,
		},
		obsidian.MagmaOperatorsUrlPart: {
			finder:       func(c echo.Context) []*protos.Identity { return getOperatorIdentity(c, version, operatorRoot) },
			identityRoot: operatorRoot,
		},
	}
}

func makeVersionedRoot(version, part string) string {
	if len(version) > 0 {
		return obsidian.RestRoot + obsidian.UrlSep + version + obsidian.UrlSep + part
	} else {
		return obsidian.RestRoot + obsidian.UrlSep + part
	}
}

// Network Identity Finder
func getNetworkIdentity(c echo.Context, version, identityRoot string) []*protos.Identity {
	if c != nil && strings.HasPrefix(c.Path(), identityRoot) {
		nid, err := obsidian.GetNetworkId(c)
		if err == nil && len(nid) > 0 {
			// All checks pass - return a Network Identity
			return []*protos.Identity{identity.NewNetwork(nid)}
		}
		// No network ID -> requires wildcard access
		return []*protos.Identity{identity.NewNetworkWildcard()}
	}
	// We don't really know what resource is being requested - request all wildcards
	return SupervisorWildcards()
}

// Operator Identity Finder
func getOperatorIdentity(c echo.Context, version, identityRoot string) []*protos.Identity {
	if c != nil && strings.HasPrefix(c.Path(), identityRoot) {
		oid, err := obsidian.GetOperatorId(c)
		if err == nil && len(oid) > 0 {
			// All checks pass - return a Network Identity
			return []*protos.Identity{identity.NewOperator(oid)}
		}
		// No network ID -> requires wildcard access
		return []*protos.Identity{identity.NewOperatorWildcard()}
	}
	// We don't really know what resource is being requested - request all wildcards
	return SupervisorWildcards()
}
