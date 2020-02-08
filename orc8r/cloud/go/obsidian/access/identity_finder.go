/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package access

import (
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/lib/go/protos"

	"github.com/labstack/echo"
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
	finder := GetIdentityFinder(c)
	if finder != nil {
		return finder(c)
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
