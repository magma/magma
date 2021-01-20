/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// access_helper provides ToString() receiver for AccessControl_Permission mask
package protos

import (
	"magma/orc8r/lib/go/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ACCESS_CONTROL_ALL_PERMISSIONS is a bitmask for all existing permissions
// unfortunately, it cannot be const since it has to be 'built' by package's
// init to simplify future maintenance
var ACCESS_CONTROL_ALL_PERMISSIONS AccessControl_Permission

func init() {
	ACCESS_CONTROL_ALL_PERMISSIONS = AccessControl_NONE
	for _, val := range AccessControl_Permission_value {
		ACCESS_CONTROL_ALL_PERMISSIONS |= AccessControl_Permission(val)
	}
}

// ToString returns a string representation of AccessControl_Permission as a mask
// protoc generated String() receiver treats it as enum and does not represent
// the 'mask' use case
func (p AccessControl_Permission) ToString() string {
	res := ""
	for mask, name := range AccessControl_Permission_name {
		if int32(p)&mask != 0 {
			if len(res) == 0 {
				res = name
			} else {
				res += "|" + name
			}
		}
	}
	if len(res) == 0 {
		res = AccessControl_Permission_name[0]
	}
	return res
}

// GetHashToACL converts the passed slice to a map, whose keys are the hash strings of each ACL's operator.
func GetHashToACL(acls []*AccessControl_List) map[string]*AccessControl_List {
	ret := make(map[string]*AccessControl_List)
	for _, acl := range acls {
		ret[acl.Operator.HashString()] = acl
	}
	return ret
}

// AddToACL adds slice of Entities to the acl of Operator 'oper'.
// If an entity with the same Id is already in the ACL, it'll be updated.
func AddToACL(acl *AccessControl_List, entities []*AccessControl_Entity) error {
	if acl == nil || len(entities) == 0 {
		return nil
	}
	if acl.Entities == nil {
		acl.Entities = map[string]*AccessControl_Entity{}
	}

	for i, ent := range entities {
		if ent == nil || ent.Id == nil {
			return status.Errorf(
				codes.InvalidArgument, "Invalid Entity @ index: %d ", i)
		}
	}
	for _, ent := range entities {
		entHashStr := ent.Id.HashString()
		acl.Entities[entHashStr] = ent
	}
	return nil
}

// GetEntityPermissions returns the aggregated ACL's permissions for a given
// entity.
// Aggregated permissions are calculated by ORing permissions for a wildcard
// of the entity type (if present in the ACL) with permissions for the entity's
// exact Identity match (if present):
//     perm = permissions[Id Type Wildcard] | permissions[Id Of Entity]
//
// getEntityPermissions will return AccessControl_NONE if the entity's identity
// is not in the list and the list doesn't have a corresponding to the entity
// type wildcard.
func GetEntityPermissions(
	acl *AccessControl_List,
	entity *protos.Identity,
) AccessControl_Permission {
	res := AccessControl_NONE
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

// CheckEntitiesPermissions verifies permissions for given entList with given
// ACL.
// Returns nil if all entities from entList have at least requested permissions
// in the ACL, error otherwise
func CheckEntitiesPermissions(
	acl *AccessControl_List,
	entList []*AccessControl_Entity,
) error {
	for _, ent := range entList {
		if ent != nil {
			reqPerm := ent.Permissions // Requested permissions for entity
			aclPerm := GetEntityPermissions(acl, ent.Id)
			if reqPerm&aclPerm != reqPerm {
				return protos.Errorf(
					codes.PermissionDenied,
					"Unsatisfied permissions, need: b%08b, got: b%08b for %s",
					reqPerm, aclPerm, ent.Id.HashString())
			}
		}
	}
	return nil
}

// VerifyPermissionsRequest is a helper function which checks validity of
// AccessControl_PermissionsRequest.
func VerifyPermissionsRequest(req *AccessControl_PermissionsRequest) error {
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

func VerifyACLRequest(req *AccessControl_ListRequest) error {
	if req == nil {
		return protos.Errorf(codes.InvalidArgument, "Nil AccessControl_ListRequest")
	}
	if req.Operator == nil {
		return protos.Errorf(codes.InvalidArgument, "Nil Operator")
	}
	return nil
}
