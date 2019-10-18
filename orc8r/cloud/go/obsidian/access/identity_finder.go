/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package access

import (
	"strings"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/protos"

	"github.com/labstack/echo"
)

const (
	MAGMA_ROOT_PART     = obsidian.RestRoot + obsidian.UrlSep
	MAGMA_ROOT_PART_LEN = len(MAGMA_ROOT_PART)
)

// FindRequestedIdentities examines the request URL and finds Identities of
// all Entities, the request needs to have access to.
//
// If FindRequestedIdentities cannot determine the entities from the request
// OR the URL is malformed OR request context is invalid - it will return a list
// of "supervisor's wildcards" - a list all known entity type wildcards
// which would correspond to an ACL typical for a supervisor/admin "can do all"
// operators
func FindRequestedIdentities(c echo.Context) []*protos.Identity {
	if c != nil {
		path := c.Path()
		if strings.HasPrefix(path, MAGMA_ROOT_PART) {
			parts := strings.Split(path[MAGMA_ROOT_PART_LEN:], obsidian.UrlSep)
			if len(parts) > 0 {
				p := parts[0]
				registry, ok := finderRegistries[p]
				if ok && len(parts) > 1 {
					p = parts[1]
				} else {
					// fall back to "versionless" V0
					registry, ok = finderRegistries[obsidian.V0]
				}
				fr, ok := registry[p]
				if ok {
					return fr.finder(c)
				}
			}
		}
	}
	return SupervisorWildcards()
}

// SupervisorWildcards returns a newly created list of "supervisor's wildcards":
// 	a list all known entity type wildcards which would correspond to an ACL
//  typical to a supervisor/admin "can do all" operators
func SupervisorWildcards() []*protos.Identity {
	return []*protos.Identity{
		identity.NewNetworkWildcard(),
		identity.NewOperatorWildcard(),
		identity.NewGatewayWildcard(),
	}
}
