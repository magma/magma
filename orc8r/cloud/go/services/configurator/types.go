/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package configurator

import (
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator/storage"
	storage2 "magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// A network represents a tenant. Networks can be configured in a hierarchical
// manner - network-level configurations are assumed to apply across multiple
// entities within the network.
type Network struct {
	ID string
	// Specifies the type of network. (lte, wifi, etc)
	Type string

	Name        string
	Description string

	// Configs maps between a type value and a generic configuration blob.
	// The map key will point to the Serde implementation which can
	// serialize the associated value.
	Configs map[string]interface{}

	Version uint64
}

func (n Network) ToStorageProto() (*storage.Network, error) {
	ret := &storage.Network{
		ID:          n.ID,
		Type:        n.Type,
		Name:        n.Name,
		Description: n.Description,
	}
	bConfigs, err := marshalConfigs(n.Configs, NetworkConfigSerdeDomain)
	if err != nil {
		return nil, errors.Wrapf(err, "error serializing network %s", n.ID)
	}
	ret.Configs = bConfigs
	return ret, nil
}

func (n Network) FromStorageProto(protoNet *storage.Network) (Network, error) {
	iConfigs, err := unmarshalConfigs(protoNet.Configs, NetworkConfigSerdeDomain)
	if err != nil {
		return n, errors.Wrapf(err, "error deserializing network %s", protoNet.ID)
	}

	n.ID = protoNet.ID
	n.Type = protoNet.Type
	n.Name = protoNet.Name
	n.Description = protoNet.Description
	n.Version = protoNet.Version
	n.Configs = iConfigs
	return n, nil
}

// NetworkUpdateCriteria specifies how to update a network
type NetworkUpdateCriteria struct {
	// ID of the network to update
	ID string

	// Set DeleteNetwork to true to delete the network
	DeleteNetwork bool

	// Set NewType, NewName or NewDescription to nil to indicate that no update is
	// desired. To clear the value of name or description, set these fields to
	// a pointer to an empty string.
	NewType        *string
	NewName        *string
	NewDescription *string

	// New config values to add or existing ones to update
	ConfigsToAddOrUpdate map[string]interface{}

	// Config values to delete
	ConfigsToDelete []string
}

func (nuc NetworkUpdateCriteria) toStorageProto() (*storage.NetworkUpdateCriteria, error) {
	bConfigs, err := marshalConfigs(nuc.ConfigsToAddOrUpdate, NetworkConfigSerdeDomain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal update")
	}

	ret := &storage.NetworkUpdateCriteria{
		ID: nuc.ID,

		DeleteNetwork: nuc.DeleteNetwork,

		NewName:        strPtrToWrapper(nuc.NewName),
		NewDescription: strPtrToWrapper(nuc.NewDescription),
		NewType:        strPtrToWrapper(nuc.NewType),

		ConfigsToAddOrUpdate: bConfigs,
		ConfigsToDelete:      nuc.ConfigsToDelete,
	}
	return ret, nil
}

// NetworkEntity is the storage representation of a logical component of a
// network. Networks are partitioned into DAGs of entities.
type NetworkEntity struct {
	// Network that the entity belongs to. This is a READ-ONLY field and will
	// be ignored if provided to write APIs.
	NetworkID string

	// (Type, Key) forms a unique identifier for the network entity within its
	// network.
	Type string
	Key  string

	Name        string
	Description string

	// PhysicalID will be non-empty if the entity corresponds to a physical
	// asset.
	PhysicalID string

	Config interface{}

	// GraphID is a mostly-internal field to designate the DAG that this
	// network entity belongs to. This field is system-generated and will be
	// ignored if set during entity creation.
	GraphID string

	// Associations are the directed edges originating from this entity.
	Associations []storage2.TypeAndKey

	// ParentAssociations are the directed edges ending at this entity.
	// This is a read-only field and will be ignored if set during entity
	// creation.
	ParentAssociations []storage2.TypeAndKey

	// Note that we are not exposing permissions in the client API at this
	// time

	Version uint64
}

func (ent NetworkEntity) toStorageProto() (*storage.NetworkEntity, error) {
	ret := &storage.NetworkEntity{
		// NetworkID is a read-only field so we won't fill it in
		Type:        ent.Type,
		Key:         ent.Key,
		Name:        ent.Name,
		Description: ent.Description,
		PhysicalID:  ent.PhysicalID,

		Associations: tksToEntIDs(ent.Associations),

		// don't set graphID, parent assocs, or version because those are
		// read-only fields
	}

	if ent.Config != nil {
		bConfigs, err := serde.Serialize(NetworkEntitySerdeDomain, ent.Type, ent.Config)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to serialize entity %s", ent.GetTypeAndKey())
		}
		ret.Config = bConfigs
	}

	return ret, nil
}

func (ent NetworkEntity) fromStorageProto(protoEnt *storage.NetworkEntity) (NetworkEntity, error) {
	ent.NetworkID = protoEnt.NetworkID
	ent.Type = protoEnt.Type
	ent.Key = protoEnt.Key
	ent.Name = protoEnt.Name
	ent.Description = protoEnt.Description
	ent.PhysicalID = protoEnt.PhysicalID
	ent.GraphID = protoEnt.GraphID
	ent.Associations = entIDsToTKs(protoEnt.Associations)
	ent.ParentAssociations = entIDsToTKs(protoEnt.ParentAssociations)
	ent.Version = protoEnt.Version

	if !funk.IsEmpty(protoEnt.Config) {
		iConfig, err := serde.Deserialize(NetworkEntitySerdeDomain, ent.Type, protoEnt.Config)
		if err != nil {
			return ent, errors.Wrapf(err, "failed to deserialize entity %s", ent.GetTypeAndKey())
		}
		ent.Config = iConfig
	}

	return ent, nil
}

func (ent NetworkEntity) GetTypeAndKey() storage2.TypeAndKey {
	return storage2.TypeAndKey{Type: ent.Type, Key: ent.Key}
}

func (ent NetworkEntity) isEntityWriteOperation() {}

type NetworkEntities []NetworkEntity

func (ne NetworkEntities) ToEntitiesByID() map[storage2.TypeAndKey]NetworkEntity {
	ret := make(map[storage2.TypeAndKey]NetworkEntity, len(ne))
	for _, ent := range ne {
		ret[ent.GetTypeAndKey()] = ent
	}
	return ret
}

// EntityGraph represents a DAG of associated network entities
type EntityGraph struct {
	Entities     []NetworkEntity
	RootEntities []storage2.TypeAndKey
	Edges        []GraphEdge

	// unexported fields for caching intermediate graph operations
	entsByTK         map[storage2.TypeAndKey]NetworkEntity
	edgesByTK        map[storage2.TypeAndKey][]storage2.TypeAndKey
	reverseEdgesByTK map[storage2.TypeAndKey][]storage2.TypeAndKey
}

func (eg EntityGraph) FromStorageProto(protoGraph *storage.EntityGraph) (EntityGraph, error) {
	eg.Entities = make([]NetworkEntity, 0, len(protoGraph.Entities))
	for _, protoEnt := range protoGraph.Entities {
		ent, err := (NetworkEntity{}).fromStorageProto(protoEnt)
		if err != nil {
			return eg, errors.Wrapf(err, "failed to deserialize entity %s", protoEnt.GetTypeAndKey())
		}
		eg.Entities = append(eg.Entities, ent)
	}

	eg.RootEntities = entIDsToTKs(protoGraph.RootEntities)

	eg.Edges = make([]GraphEdge, 0, len(protoGraph.Edges))
	for _, protoEdge := range protoGraph.Edges {
		eg.Edges = append(eg.Edges, (GraphEdge{}).fromStorageProto(protoEdge))
	}
	return eg, nil
}

func (eg EntityGraph) ToStorageProto() (*storage.EntityGraph, error) {
	protoGraph := &storage.EntityGraph{}
	for _, ent := range eg.Entities {
		protoEnt, err := ent.toStorageProto()
		if err != nil {
			return protoGraph, errors.Wrapf(err, "failed to convert entity %s to storage proto", ent.GetTypeAndKey())
		}
		protoGraph.Entities = append(protoGraph.Entities, protoEnt)
	}
	protoGraph.RootEntities = tksToEntIDs(eg.RootEntities)

	for _, edge := range eg.Edges {
		protoGraph.Edges = append(protoGraph.Edges, edge.toStorageProto())
	}
	return protoGraph, nil
}

type GraphEdge struct {
	From storage2.TypeAndKey
	To   storage2.TypeAndKey
}

func (ge GraphEdge) fromStorageProto(protoEdge *storage.GraphEdge) GraphEdge {
	ge.From = protoEdge.From.ToTypeAndKey()
	ge.To = protoEdge.To.ToTypeAndKey()
	return ge
}

func (ge GraphEdge) toStorageProto() *storage.GraphEdge {
	protoTo := (&storage.EntityID{}).FromTypeAndKey(ge.To)
	protoFrom := (&storage.EntityID{}).FromTypeAndKey(ge.From)
	return &storage.GraphEdge{
		To:   protoTo,
		From: protoFrom,
	}
}

// EntityLoadFilter specifies which entities to load from storage
type EntityLoadFilter struct {
	// If TypeFilter is provided, the query will return all entities matching
	// the given type.
	TypeFilter *string

	// If KeyFilter is provided, the query will return all entities matching the
	// given ID.
	KeyFilter *string

	// If IDs is provided, the query will return all entities matching the
	// provided TypeAndKeys. TypeFilter and KeyFilter are ignored if IDs is
	// provided.
	IDs []storage2.TypeAndKey

	// If PhysicalID is provided, the query will return all entities matching
	// the provided ID value.
	PhysicalID *string
}

// EntityLoadCriteria specifies how much of an entity to load
type EntityLoadCriteria struct {
	// Set LoadMetadata to true to load the metadata fields (name, description)
	LoadMetadata bool

	LoadConfig bool

	LoadAssocsToThis   bool
	LoadAssocsFromThis bool
}

func (elc EntityLoadCriteria) toStorageProto() *storage.EntityLoadCriteria {
	return &storage.EntityLoadCriteria{
		LoadMetadata:       elc.LoadMetadata,
		LoadConfig:         elc.LoadConfig,
		LoadAssocsToThis:   elc.LoadAssocsToThis,
		LoadAssocsFromThis: elc.LoadAssocsFromThis,
	}
}

// FullEntityLoadCriteria returns an EntityLoadCriteria that loads everything
// possible on an entity
func FullEntityLoadCriteria() EntityLoadCriteria {
	return EntityLoadCriteria{
		LoadMetadata:       true,
		LoadConfig:         true,
		LoadAssocsToThis:   true,
		LoadAssocsFromThis: true,
	}
}

// EntityLoadResult encapsulates the result of a LoadEntities call
type EntityLoadResult struct {
	// Loaded entities
	Entities []NetworkEntity
	// Entities which were not found
	EntitiesNotFound []storage2.TypeAndKey
}

// EntityWriteOperation is an interface around entity creation/update for the
// generic multi-operation configurator endpoint.
type EntityWriteOperation interface {
	isEntityWriteOperation()
}

// EntityUpdateCriteria specifies a patch operation on a network entity.
type EntityUpdateCriteria struct {
	// (Type, Key) of the entity to update
	Type string
	Key  string

	// Set DeleteEntity to true to mark the entity for deletion
	DeleteEntity bool

	NewName        *string
	NewDescription *string

	NewPhysicalID *string

	// A nil value here indicates no update.
	NewConfig interface{}

	// Set to true to clear the entity's config
	DeleteConfig bool

	// IMPORTANT: Setting AssociationsToSet to an empty, non-nil array value
	// specifies an intent to clear all associations originating from this
	// entity.
	// A nil field value will be ignored.
	AssociationsToSet    []storage2.TypeAndKey
	AssociationsToAdd    []storage2.TypeAndKey
	AssociationsToDelete []storage2.TypeAndKey
}

func (euc EntityUpdateCriteria) toStorageProto() (*storage.EntityUpdateCriteria, error) {
	ret := &storage.EntityUpdateCriteria{
		Type:                 euc.Type,
		Key:                  euc.Key,
		DeleteEntity:         euc.DeleteEntity,
		NewName:              strPtrToWrapper(euc.NewName),
		NewDescription:       strPtrToWrapper(euc.NewDescription),
		NewPhysicalID:        strPtrToWrapper(euc.NewPhysicalID),
		AssociationsToAdd:    tksToEntIDs(euc.AssociationsToAdd),
		AssociationsToDelete: tksToEntIDs(euc.AssociationsToDelete),
	}

	if euc.AssociationsToSet != nil {
		ret.AssociationsToSet = &storage.EntityAssociationsToSet{
			AssociationsToSet: tksToEntIDs(euc.AssociationsToSet),
		}
	}

	if euc.NewConfig != nil {
		bConfig, err := serde.Serialize(NetworkEntitySerdeDomain, euc.Type, euc.NewConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to serialize update %s", euc.GetTypeAndKey())
		}
		ret.NewConfig = &wrappers.BytesValue{Value: bConfig}
	}
	if euc.DeleteConfig {
		ret.NewConfig = &wrappers.BytesValue{Value: []byte{}}
	}

	return ret, nil
}

func (euc EntityUpdateCriteria) GetTypeAndKey() storage2.TypeAndKey {
	return storage2.TypeAndKey{Type: euc.Type, Key: euc.Key}
}

func (euc EntityUpdateCriteria) isEntityWriteOperation() {}

func marshalConfigs(configs map[string]interface{}, domain string) (map[string][]byte, error) {
	ret := map[string][]byte{}
	for configType, iConfig := range configs {
		sConfig, err := serde.Serialize(domain, configType, iConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to serialize config %s", configType)
		}
		ret[configType] = sConfig
	}
	return ret, nil
}

func unmarshalConfigs(configs map[string][]byte, domain string) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	for configType, sConfig := range configs {
		iConfig, err := serde.Deserialize(domain, configType, sConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to deserialize config %s", configType)
		}
		ret[configType] = iConfig
	}
	return ret, nil
}

func strPtrToWrapper(in *string) *wrappers.StringValue {
	if in == nil {
		return nil
	}
	return &wrappers.StringValue{Value: *in}
}

func tksToEntIDs(tks []storage2.TypeAndKey) []*storage.EntityID {
	if funk.IsEmpty(tks) {
		return nil
	}

	return funk.Map(
		tks,
		func(tk storage2.TypeAndKey) *storage.EntityID { return (&storage.EntityID{}).FromTypeAndKey(tk) }).([]*storage.EntityID)
}

func entIDsToTKs(ids []*storage.EntityID) []storage2.TypeAndKey {
	if funk.IsEmpty(ids) {
		return nil
	}

	return funk.Map(
		ids,
		func(id *storage.EntityID) storage2.TypeAndKey { return id.ToTypeAndKey() },
	).([]storage2.TypeAndKey)
}
