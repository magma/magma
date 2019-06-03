/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package protos

import (
	"magma/orc8r/cloud/go/services/configurator/storage"
	commonStorage "magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/ptypes/wrappers"
)

// ToNetwork translates protobuf struct to corresponding storage struct
func (network *Network) ToNetwork() storage.Network {
	return storage.Network{
		ID:          network.Id,
		Name:        network.Name,
		Description: network.Description,
		Configs:     network.Configs,
	}
}

// ToNetworkEntity translates protobuf struct to corresponding storage struct
func (entity *NetworkEntity) ToNetworkEntity() storage.NetworkEntity {
	return storage.NetworkEntity{
		Key:                entity.Id,
		Type:               entity.Type,
		Name:               entity.Name,
		Description:        entity.Description,
		PhysicalID:         entity.PhysicalId,
		Config:             entity.Config,
		GraphID:            entity.GraphID,
		Associations:       ToTypeAndKeys(entity.Assocs),
		ParentAssociations: ToTypeAndKeys(entity.ParentAssocs),
		Permissions:        toStorageACLs(entity.Permissions),
	}
}

// ToTypeAndKey translates protobuf struct to corresponding storage struct
func (id *EntityID) ToTypeAndKey() commonStorage.TypeAndKey {
	return commonStorage.TypeAndKey{
		Type: id.Type,
		Key:  id.Id,
	}
}

// ToTypeAndKeys applies ToTypeAndKey to a slice of protobuf EntityIDs
func ToTypeAndKeys(entityIDs []*EntityID) []commonStorage.TypeAndKey {
	tks := []commonStorage.TypeAndKey{}
	for _, entity := range entityIDs {
		tks = append(tks, entity.ToTypeAndKey())
	}
	return tks
}

// ToEntityLoadCriteria translates protobuf struct to corresponding storage struct
func (criteria *EntityLoadCriteria) ToEntityLoadCriteria() storage.EntityLoadCriteria {
	return storage.EntityLoadCriteria{
		LoadMetadata:       criteria.LoadMetadata,
		LoadConfig:         criteria.LoadConfig,
		LoadAssocsToThis:   criteria.LoadAssocsTo,
		LoadAssocsFromThis: criteria.LoadAssocsFrom,
		LoadPermissions:    criteria.LoadPermissions,
	}
}

// ToNetworkLoadCriteria translates protobuf struct to corresponding storage struct
func (criteria *NetworkLoadCriteria) ToNetworkLoadCriteria() storage.NetworkLoadCriteria {
	return storage.NetworkLoadCriteria{
		LoadMetadata: criteria.LoadMetadata,
		LoadConfigs:  criteria.LoadConfigs,
	}
}

// ToEntityLoadFilter translates protobuf struct to corresponding storage struct
func ToEntityLoadFilter(typeFilter *wrappers.StringValue, keyFilter *wrappers.StringValue, ids []*EntityID) storage.EntityLoadFilter {
	entityLoadFilter := storage.EntityLoadFilter{
		TypeFilter: getStringPointer(typeFilter),
		KeyFilter:  getStringPointer(keyFilter),
		IDs:        ToTypeAndKeys(ids),
	}
	return entityLoadFilter
}

// ToNetworkUpdateCriteria translates protobuf struct to corresponding storage struct
func (criteria *NetworkUpdateCriteria) ToNetworkUpdateCriteria() storage.NetworkUpdateCriteria {
	return storage.NetworkUpdateCriteria{
		ID:                   criteria.Id,
		DeleteNetwork:        false,
		NewName:              getStringPointer(criteria.NewName),
		NewDescription:       getStringPointer(criteria.NewDescription),
		ConfigsToAddOrUpdate: criteria.ConfigsToAddOrUpdate,
		ConfigsToDelete:      criteria.ConfigsToDelete,
	}
}

// ToEntityUpdateCriteria translates protobuf struct to corresponding storage struct
func (criteria *EntityUpdateCriteria) ToEntityUpdateCriteria() storage.EntityUpdateCriteria {
	return storage.EntityUpdateCriteria{
		Type:                 criteria.Type,
		Key:                  criteria.Key,
		NewName:              getStringPointer(criteria.NewName),
		NewDescription:       getStringPointer(criteria.NewDescription),
		NewPhysicalID:        getStringPointer(criteria.NewPhysicalID),
		NewConfig:            getBytesPointer(criteria.NewConfig),
		AssociationsToAdd:    ToTypeAndKeys(criteria.AssociationsToAdd),
		AssociationsToDelete: ToTypeAndKeys(criteria.AssociationsToDelete),
		PermissionsToCreate:  toStorageACLs(criteria.PermissionsToCreate),
		PermissionsToUpdate:  toStorageACLs(criteria.PermissionsToUpdate),
		PermissionsToDelete:  criteria.PermissionsToDelete,
	}
}

// FromStorageNetwork translates storage struct to corresponding protobuf struct
func FromStorageNetwork(network storage.Network) *Network {
	return &Network{
		Id:          network.ID,
		Name:        network.Name,
		Description: network.Description,
		Configs:     network.Configs,
	}
}

// FromTKs translates storage struct to corresponding protobuf struct
func FromTKs(tks []commonStorage.TypeAndKey) []*EntityID {
	ids := []*EntityID{}
	for _, tk := range tks {
		ids = append(ids, &EntityID{Type: tk.Type, Id: tk.Key})
	}
	return ids
}

// FromStorageNetworkEntity translates storage struct to corresponding protobuf struct
func FromStorageNetworkEntity(entity storage.NetworkEntity) *NetworkEntity {
	return &NetworkEntity{
		Type:         entity.Type,
		Id:           entity.Key,
		Name:         entity.Name,
		Description:  entity.Description,
		PhysicalId:   entity.PhysicalID,
		Config:       entity.Config,
		GraphID:      entity.GraphID,
		Assocs:       FromTKs(entity.Associations),
		ParentAssocs: FromTKs(entity.ParentAssociations),
		Permissions:  fromStorageACLs(entity.Permissions),
	}
}

// FromStorageNetworkEntities translates storage struct to corresponding protobuf struct
func FromStorageNetworkEntities(entities []storage.NetworkEntity) []*NetworkEntity {
	pEntities := []*NetworkEntity{}
	for _, entity := range entities {
		pEntities = append(pEntities, FromStorageNetworkEntity(entity))
	}
	return pEntities
}

// GetStringWrapper wraps a pointer string value into protobuf StringValue
func GetStringWrapper(pStr *string) *wrappers.StringValue {
	if pStr == nil {
		return nil
	}
	return &wrappers.StringValue{Value: *pStr}
}

// GetStringWrapper wraps a pointer string value into protobuf StringValue
func GetBytesWrapper(bytes []byte) *wrappers.BytesValue {
	if bytes == nil {
		return nil
	}
	return &wrappers.BytesValue{Value: bytes}
}

func (acl *ACL) toACL() storage.ACL {
	return storage.ACL{
		ID:         acl.Id,
		Scope:      acl.getACLScope(),
		Permission: storage.ACLPermission(acl.Permission),
		Type:       acl.getACLType(),
		IDFilter:   acl.IdFilter,
	}
}

func fromStorageACLs(acls []storage.ACL) []*ACL {
	pACLs := []*ACL{}
	for _, acl := range acls {
		pACL := &ACL{
			Id:         acl.ID,
			Permission: ACL_Permission(acl.Permission),
			IdFilter:   acl.IDFilter,
		}
		if acl.Scope.NetworkIDs != nil {
			pACL.Scope = &ACL_NetworkIds{NetworkIds: &ACL_NetworkIDs{Ids: acl.Scope.NetworkIDs}}
		} else {
			pACL.Scope = &ACL_ScopeWildcard{ACL_Wildcard(acl.Scope.Wildcard)}
		}
		if acl.Type.EntityType != "" {
			pACL.Type = &ACL_EntityType{acl.Type.EntityType}
		} else {
			pACL.Type = &ACL_TypeWildcard{ACL_Wildcard(acl.Type.Wildcard)}
		}
		pACLs = append(pACLs, pACL)
	}
	return pACLs
}

func (acl *ACL) getACLScope() storage.ACLScope {
	aclScope := storage.ACLScope{}
	switch acl.Scope.(type) {
	case *ACL_NetworkIds:
		aclScope.NetworkIDs = acl.GetNetworkIds().Ids
	case *ACL_ScopeWildcard:
		aclScope.Wildcard = storage.ACLWildcard(acl.GetScopeWildcard())
	}
	return aclScope
}

func (acl *ACL) getACLType() storage.ACLType {
	aclType := storage.ACLType{}
	switch acl.Type.(type) {
	case *ACL_EntityType:
		aclType.EntityType = acl.GetEntityType()
	case *ACL_TypeWildcard:
		aclType.Wildcard = storage.ACLWildcard(acl.GetTypeWildcard())
	}
	return aclType
}

func toStorageACLs(acls []*ACL) []storage.ACL {
	sACLs := []storage.ACL{}
	for _, acl := range acls {
		sACLs = append(sACLs, acl.toACL())
	}
	return sACLs
}

func getStringPointer(strWrapper *wrappers.StringValue) *string {
	if strWrapper != nil {
		return &(strWrapper.Value)
	}
	return nil
}

func getBytesPointer(bytesWrapper *wrappers.BytesValue) *[]byte {
	if bytesWrapper != nil {
		return &(bytesWrapper.Value)
	}
	return nil
}
