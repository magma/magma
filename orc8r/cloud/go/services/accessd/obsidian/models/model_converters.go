/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models

import (
	"magma/orc8r/cloud/go/identity"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes/duration"
)

var formatsRegistry = strfmt.NewFormats()

func PermissionsMaskToProto(mask PermissionsMask) accessprotos.AccessControl_Permission {
	permissions := accessprotos.AccessControl_Permission_value[string(mask[0])]
	permissions |= accessprotos.AccessControl_Permission_value[string(mask[1])]
	return accessprotos.AccessControl_Permission(permissions)
}

func PermissionsMaskFromProto(permissions accessprotos.AccessControl_Permission) PermissionsMask {
	return PermissionsMask{
		PermissionType(accessprotos.AccessControl_Permission_name[int32(permissions&accessprotos.AccessControl_READ)]),
		PermissionType(accessprotos.AccessControl_Permission_name[int32(permissions&accessprotos.AccessControl_WRITE)]),
	}
}

func ACLEntityToProto(entity *ACLEntity) *accessprotos.AccessControl_Entity {
	permissions := PermissionsMaskToProto(entity.Permissions)
	id := identityFromEntity(entity)
	accessControlEntity := &accessprotos.AccessControl_Entity{
		Id:          id,
		Permissions: permissions,
	}
	return accessControlEntity
}

func ACLEntityFromProto(accessControlEntity *accessprotos.AccessControl_Entity) *ACLEntity {
	aclEntity := &ACLEntity{
		Permissions: PermissionsMaskFromProto(accessControlEntity.Permissions),
	}
	if wildcard, ok := accessControlEntity.Id.Value.(*protos.Identity_Wildcard_); ok {
		switch wildcard.Wildcard.GetType() {
		case protos.Identity_Wildcard_Operator:
			aclEntity.EntityType = ACLEntityEntityTypeOPERATORWILDCARD
		case protos.Identity_Wildcard_Network:
			aclEntity.EntityType = ACLEntityEntityTypeNETWORKWILDCARD
		}
	} else if operator, ok := accessControlEntity.Id.Value.(*protos.Identity_Operator); ok {
		aclEntity.EntityType = ACLEntityEntityTypeOPERATOR
		aclEntity.OperatorID = OperatorID(*operator.ToCommonName())
	} else if network, ok := accessControlEntity.Id.Value.(*protos.Identity_Network); ok {
		aclEntity.EntityType = ACLEntityEntityTypeNETWORK
		aclEntity.NetworkID = NetworkID(*network.ToCommonName())
	}
	return aclEntity
}

func ACLToProto(acl ACLType) []*accessprotos.AccessControl_Entity {
	aclEntities := []*ACLEntity(acl)
	accessControlList := make([]*accessprotos.AccessControl_Entity, len(aclEntities))
	for i, entity := range aclEntities {
		accessControlList[i] = ACLEntityToProto(entity)
	}
	return accessControlList
}

func ACLFromProto(accessControlList []*accessprotos.AccessControl_Entity) ACLType {
	aclEntities := make([]*ACLEntity, len(accessControlList))
	for i, accessControlEntity := range accessControlList {
		aclEntities[i] = ACLEntityFromProto(accessControlEntity)
	}
	return ACLType(aclEntities)
}

func CSRToProto(csr *CsrType, operator *protos.Identity) *protos.CSR {
	return &protos.CSR{
		Id: operator,
		ValidTime: &duration.Duration{
			Seconds: csr.Duration.Seconds,
			Nanos:   csr.Duration.Nanoseconds,
		},
		CsrDer: []byte(csr.CsrDer),
	}
}

func CSRFromProto(csr *protos.CSR) *CsrType {
	return &CsrType{
		CsrDer: CsrDer(csr.CsrDer),
		Duration: &DurationType{
			Seconds:     csr.ValidTime.GetSeconds(),
			Nanoseconds: csr.ValidTime.GetNanos(),
		},
	}
}

func init() {
	// Echo encodes/decodes base64 encoded byte arrays, no verification needed
	b64 := strfmt.Base64([]byte(nil))
	formatsRegistry.Add("byte", &b64, func(_ string) bool { return true })
}

func identityFromEntity(entity *ACLEntity) *protos.Identity {
	switch entity.EntityType {
	case ACLEntityEntityTypeNETWORK:
		return identity.NewNetwork(string(entity.NetworkID))
	case ACLEntityEntityTypeNETWORKWILDCARD:
		return identity.NewNetworkWildcard()
	case ACLEntityEntityTypeOPERATOR:
		return identity.NewOperator(string(entity.OperatorID))
	case ACLEntityEntityTypeOPERATORWILDCARD:
		return identity.NewOperatorWildcard()
	}
	return nil
}
