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

var finderRegistry = map[string]RequestIdentityFinder{
	obsidian.MagmaNetworksUrlPart:  getNetworkIdentity,
	obsidian.MagmaOperatorsUrlPart: getOperatorIdentity,
}

// Network Identity Finder
func getNetworkIdentity(c echo.Context) []*protos.Identity {
	if c != nil && strings.HasPrefix(c.Path(), obsidian.NetworksRoot) {
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
func getOperatorIdentity(c echo.Context) []*protos.Identity {
	if c != nil && strings.HasPrefix(c.Path(), obsidian.OperatorsRoot) {
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
