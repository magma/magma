/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

// Internal accessd related utility functions
import (
	"magma/orc8r/cloud/go/protos"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
)

// storeParamsPair returns a pair of strings to aid querying Access Control
// related datastore tables - (key, table)
// key: id's Identity Hash String to be used as the query key and
// table: currently is ACCESS_TABLE
func getKeyTablePair(id *protos.Identity) (string, string) {
	if id != nil {
		table := ACCESS_TABLE
		return id.HashString(), table
	}
	return "", ""
}

// getACL fetches Operator's ACL from srv store, verifies and returns it
func (srv *AccessControlServer) getACL(
	oper *protos.Identity,
) (*accessprotos.AccessControl_List, error) {

	if oper == nil {
		return &accessprotos.AccessControl_List{}, protos.Errorf(codes.InvalidArgument, "Nil Operator")
	}
	opkey, table := getKeyTablePair(oper)
	return srv.getACLForKey(table, opkey)
}

func (srv *AccessControlServer) getACLForKey(
	table string, opkey string,
) (*accessprotos.AccessControl_List, error) {
	acl := &accessprotos.AccessControl_List{}
	marshaledAcl, _, err := srv.store.Get(table, opkey)
	if err != nil {
		return acl,
			protos.Errorf(codes.NotFound,
				"Get ACL error '%s' for Operator %s, table %s",
				err, opkey, table)
	}
	err = proto.Unmarshal(marshaledAcl, acl)
	if err != nil {
		return acl,
			protos.Errorf(codes.Unknown,
				"ACL Unmarshal error '%s' for Operator %s from table %s",
				err, opkey, table)
	}
	if acl.Operator == nil {
		// in the case of unlikely table corruption - return bad ACL as well
		// to help caller to diagnose the issue
		return acl,
			protos.Errorf(codes.Unknown,
				"Nil Operator for ACL Key %s; Table %s may be corupted",
				opkey, table)
	}
	aclOpHashStr := acl.Operator.HashString()
	if aclOpHashStr != opkey { // verify, operator is the same as ACL's owner
		// in the case of unlikely table corruption - return bad ACL as well
		// to help caller to diagnose the issue
		return acl, protos.Errorf(codes.Unknown,
			"Corrupted ACL in table %s: ACL Operator %s and key %s mismatch",
			table, aclOpHashStr, opkey)
	}

	return acl, nil
}

// getACLs fetches Operators' ACLs from srv store, verifies and returns them
func (srv *AccessControlServer) getACLs(
	opers []*protos.Identity,
) ([]*accessprotos.AccessControl_List, error) {
	if opers == nil {
		return []*accessprotos.AccessControl_List{},
			protos.Errorf(codes.InvalidArgument, "Nil Operator List")
	}
	keys := make([]string, 0, len(opers))
	for _, oper := range opers {
		if oper != nil {
			keys = append(keys, oper.HashString())
		}
	}
	return srv.getACLsForKeys(ACCESS_TABLE, keys)
}

func (srv *AccessControlServer) getACLsForKeys(
	table string, opkeys []string,
) ([]*accessprotos.AccessControl_List, error) {
	marshaledAclValues, err := srv.store.GetMany(table, opkeys)
	if err != nil {
		return nil,
			protos.Errorf(codes.NotFound,
				"Get ACLs error '%s' for Operators %v, table %s",
				err, opkeys, table)
	}
	marshaledAcls := make([]*accessprotos.AccessControl_List, 0, len(marshaledAclValues))
	for opkey, marshaledAclVal := range marshaledAclValues {
		acl := &accessprotos.AccessControl_List{}
		err = proto.Unmarshal(marshaledAclVal.Value, acl)
		if err != nil {
			return marshaledAcls,
				protos.Errorf(codes.Unknown,
					"ACLs Unmarshal error '%s' for Operator %s from table %s",
					err, opkey, table)
		}
		if acl.Operator == nil {
			// in the case of unlikely table corruption - return bad ACL as well
			// to help caller to diagnose the issue
			return marshaledAcls,
				protos.Errorf(codes.Unknown,
					"Nil Operator for Operator key %s, ACL %v; Table %s may be corupted",
					opkey, acl, table)
		}
		aclOpHashStr := acl.Operator.HashString()
		if aclOpHashStr != opkey { // verify, operator is the same as ACL's owner
			// in the case of unlikely table corruption - return bad ACL as well
			// to help caller to diagnose the issue
			return marshaledAcls, protos.Errorf(codes.Unknown,
				"Corrupted ACL in table %s: ACL Operator %s and Operator key %s mismatch",
				table, aclOpHashStr, opkey)
		}
		marshaledAcls = append(marshaledAcls, acl)
	}
	return marshaledAcls, nil
}

// putACL writes Operator's ACL to srv store
func (srv *AccessControlServer) putACL(oper *protos.Identity, acl *accessprotos.AccessControl_List) error {
	if oper == nil {
		return protos.Errorf(codes.InvalidArgument, "Nil Operator for ACL")
	}
	opkey, table := getKeyTablePair(oper)
	if acl == nil {
		return protos.Errorf(
			codes.InvalidArgument, "Nil ACL for Operator %s", opkey)
	}

	marshaledAcl, err := proto.Marshal(acl)
	if err != nil {
		return protos.Errorf(codes.Unknown,
			"ACL Marshal error '%s' for Operator %s", err, opkey)
	}
	err = srv.store.Put(table, opkey, marshaledAcl)
	if err != nil {
		return protos.Errorf(
			codes.Unknown,
			"ACL PUT error '%s' for Operator %s, table %s",
			err, opkey, table)
	}
	return nil
}

// Add slice of Entities to the acl of Operator 'oper'
// If an entity with the same Id is already in the ACL, it'll be updated
func addToACL(
	oper *protos.Identity,
	acl *accessprotos.AccessControl_List,
	entities []*accessprotos.AccessControl_Entity,
) error {
	if oper == nil || acl == nil || len(entities) == 0 {
		return nil
	}
	if acl.Entities == nil {
		acl.Entities = map[string]*accessprotos.AccessControl_Entity{}
	}

	opkey := oper.HashString()

	aclOpHashStr := acl.Operator.HashString()
	if aclOpHashStr != opkey { // verify, operator is the same as ACL's owner
		return protos.Errorf(codes.Unknown,
			"Corrupted ACL: Operator hash '%s' and key '%s' mismatch",
			aclOpHashStr, opkey)
	}

	for i, ent := range entities {
		if ent == nil || ent.Id == nil {
			return protos.Errorf(
				codes.InvalidArgument, "Invalid Entity @ index: %d ", i)
		}
	}
	for _, ent := range entities {
		entHashStr := ent.Id.HashString()
		acl.Entities[entHashStr] = ent
	}
	return nil
}

// Returns the managing Identity's ACL Entity for given entity identity
func (srv *AccessControlServer) getACLEntity(
	req *accessprotos.AccessControl_PermissionsRequest,
) (*accessprotos.AccessControl_Entity, error) {
	res := &accessprotos.AccessControl_Entity{}
	err := verifyPermissionsRequest(req)
	if err != nil {
		return res, err
	}
	acl, err := srv.getACL(req.Operator)
	if err != nil {
		return res, err
	}
	res.Id = req.Entity
	res.Permissions = getEntityPermissions(acl, res.Id) // Aggregated entity permissions
	return res, nil
}

// getEntityPermissions returns the aggregated ACL's permissions for a given
// entity.
// Aggregated permissions are calculated by ORing permissions for a wildcard
// of the entity type (if present in the ACL) with permissions for the entity's
// exact Identity match (if present):
//     perm = permissions[Id Type Wildcard] | permissions[Id Of Entity]
//
// getEntityPermissions will return AccessControl_NONE if the entity's identity
// is not in the list and the list doesn't have a corresponding to the entity
// type wildcard.
func getEntityPermissions(
	acl *accessprotos.AccessControl_List,
	entity *protos.Identity,
) accessprotos.AccessControl_Permission {
	res := accessprotos.AccessControl_NONE
	if acl != nil && entity != nil {
		if wc := entity.GetWildcardForIdentity(); wc != nil {
			hash := wc.HashString()
			ent, ok := acl.Entities[hash]
			if ok && ent.Id.Match(entity) {
				res = ent.Permissions
			}
		}
		hash := entity.HashString()
		ent, ok := acl.Entities[hash]
		if ok {
			res |= ent.Permissions
		}
	}
	return res
}

// checkEntitiesPermissions verifies permissions for given entList with given
// ACL.
// Returns nil if all entities from entList have at least requested permissions
// in the ACL, error otherwise
func checkEntitiesPermissions(
	acl *accessprotos.AccessControl_List,
	entList []*accessprotos.AccessControl_Entity,
) error {
	if entList != nil {
		for _, ent := range entList {
			if ent != nil {
				reqPerm := ent.Permissions // Requested permissions for entity
				aclPerm := getEntityPermissions(acl, ent.Id)
				if reqPerm&aclPerm != reqPerm {
					return protos.Errorf(
						codes.PermissionDenied,
						"Unsatisfied permissions, need: b%08b, got: b%08b for %s",
						reqPerm, aclPerm, ent.Id.HashString())
				}
			}
		}

	}
	return nil
}

// verifyPermissionsRequest is a helper function which checks validity of
// AccessControl_PermissionsRequest
func verifyPermissionsRequest(req *accessprotos.AccessControl_PermissionsRequest) error {
	if req == nil {
		return protos.Errorf(codes.InvalidArgument, "Nil PermissionsRequest")
	}
	if req.Operator == nil {
		return protos.Errorf(codes.InvalidArgument, "Nil PermissionsRequest Operator")
	}
	if req.Entity == nil {
		return protos.Errorf(codes.InvalidArgument, "Nil PermissionsRequest Entity")
	}
	return nil
}

func verifyACLRequest(req *accessprotos.AccessControl_ListRequest) error {
	if req == nil {
		return protos.Errorf(codes.InvalidArgument, "Nil AccessControl_ListRequest")
	}
	if req.Operator == nil {
		return protos.Errorf(codes.InvalidArgument, "Nil Operator")
	}
	return nil
}
