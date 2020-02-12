/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package storage

import (
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/lib/go/protos"
)

// AccessdStorage provides storage functionality for mapping identities to ACLs.
// Methods return errors with relevant grpc/codes/Code codes embedded in the error string,
// including guards against nil arguments.
type AccessdStorage interface {
	// ListAllIdentity returns all identities tracked as part of an ACL.
	ListAllIdentity() ([]*protos.Identity, error)

	// GetACL returns the ACL for the passed identity.
	// If not found, returns wrapped codes.NotFound.
	GetACL(id *protos.Identity) (*accessprotos.AccessControl_List, error)

	// GetManyACL returns a list of ACLs for the passed identities.
	GetManyACL(id []*protos.Identity) ([]*accessprotos.AccessControl_List, error)

	// PutACL associates an ACL with an identity, overwriting any previous ACL.
	PutACL(id *protos.Identity, acl *accessprotos.AccessControl_List) error

	// UpdateACLWithEntities updates the ID's ACL with additional entities.
	UpdateACLWithEntities(id *protos.Identity, entities []*accessprotos.AccessControl_Entity) error

	// DeleteACL removes the ACL associated with the passed identity.
	DeleteACL(id *protos.Identity) error
}
