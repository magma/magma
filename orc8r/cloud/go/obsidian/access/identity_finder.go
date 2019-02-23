/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package access

import (
	"strings"

	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/protos"
)

const (
	MAGMA_ROOT_PART     = handlers.REST_ROOT + handlers.URL_SEP
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
			parts := strings.Split(path[MAGMA_ROOT_PART_LEN:], handlers.URL_SEP)
			if len(parts) > 0 {
				finder, ok := finderRegistry[parts[0]]
				if ok {
					return finder(c)
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
