/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package configurator

import (
	"fmt"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator/storage"
	storage2 "magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// Network partitions a set of entities. Networks can be configured in a
// hierarchical manner - network-level configurations are assumed to apply
// across multiple entities within the network.
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

// ToProto converts a Network struct to its proto representation.
func (n Network) ToProto(serdes serde.Registry) (*storage.Network, error) {
	ret := &storage.Network{
		ID:          n.ID,
		Type:        n.Type,
		Name:        n.Name,
		Description: n.Description,
	}
	bConfigs, err := marshalConfigs(n.Configs, serdes)
	if err != nil {
		return nil, errors.Wrapf(err, "error serializing network %s", n.ID)
	}
	ret.Configs = bConfigs
	return ret, nil
}

// FromProto converts a Network proto to it's internal go struct format.
func (n Network) FromProto(protoNet *storage.Network, serdes serde.Registry) (Network, error) {
	iConfigs, err := unmarshalConfigs(protoNet.Configs, serdes)
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

func (nuc NetworkUpdateCriteria) toProto(serdes serde.Registry) (*storage.NetworkUpdateCriteria, error) {
	bConfigs, err := marshalConfigs(nuc.ConfigsToAddOrUpdate, serdes)
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

// NetworkEntity is a logical component of a network.
// Networks are partitioned into DAGs of entities. Read-only fields will be
// ignored if provided to write APIs.
type NetworkEntity struct {
	// Network the entity belongs to.
	// Read-only.
	NetworkID string

	// (Type, Key) forms a unique identifier for the network entity within its
	// network.
	Type string
	Key  string

	// Config value for the entity, deserialized.
	Config interface{}
	// isSerialized is true iff Config was set by this pkg and the contents
	// were not deserialized from byte array.
	isSerialized bool

	// Name and Description are metadata annotations for human consumption.
	Name        string
	Description string

	// PhysicalID is non-empty when the entity corresponds to a physical asset.
	PhysicalID string

	// GraphID is a mostly-internal field to designate the DAG that this
	// network entity belongs to.
	// Read-only.
	GraphID string

	// Associations are the directed edges originating from this entity.
	Associations storage2.TKs

	// ParentAssociations are the directed edges ending at this entity.
	// Read-only.
	ParentAssociations storage2.TKs

	// Version of the entity.
	// Read-only.
	Version uint64
}

func (ent NetworkEntity) toProto(serdes serde.Registry) (*storage.NetworkEntity, error) {
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
		bConfigs, err := serde.Serialize(ent.Config, ent.Type, serdes)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to serialize entity %s", ent.GetTypeAndKey())
		}
		ret.Config = bConfigs
	}

	return ret, nil
}

func (ent NetworkEntity) fromProto(p *storage.NetworkEntity, serdes serde.Registry) (NetworkEntity, error) {
	e := ent.fromProtoSerialized(p)

	if len(p.Config) != 0 {
		iConfig, err := serde.Deserialize(p.Config, e.Type, serdes)
		if err != nil {
			return ent, errors.Wrapf(err, "failed to deserialize entity %s", ent.GetTypeAndKey())
		}
		e.Config = iConfig
		e.isSerialized = false
	}

	return e, nil
}

// fromProtoSerialized is the same as fromProto, except it leaves the entity's
// config as serialized bytes.
func (ent NetworkEntity) fromProtoSerialized(p *storage.NetworkEntity) NetworkEntity {
	e := NetworkEntity{
		NetworkID:          p.NetworkID,
		Type:               p.Type,
		Key:                p.Key,
		Name:               p.Name,
		Description:        p.Description,
		PhysicalID:         p.PhysicalID,
		GraphID:            p.GraphID,
		Associations:       entIDsToTKs(p.Associations),
		ParentAssociations: entIDsToTKs(p.ParentAssociations),
		Version:            p.Version,
	}
	if len(p.Config) != 0 {
		e.Config = p.Config
		e.isSerialized = true
	}
	return e
}

// fromProtoWithDefault tries to return fromProto, defaulting to
// fromProtoSerialized when it encounters an error.
func (ent NetworkEntity) fromProtoWithDefault(p *storage.NetworkEntity, serdes serde.Registry) (NetworkEntity, error) {
	// Default to returning serialized when no serde found
	if !serde.HasSerde(serdes, p.Type) {
		return ent.fromProtoSerialized(p), nil
	}

	e, err := ent.fromProto(p, serdes)
	if err != nil {
		return NetworkEntity{}, err
	}
	return e, nil
}

// IsSerialized returns true iff this package created the ent and its state
// was not deserialized.
func (ent NetworkEntity) IsSerialized() bool {
	return ent.isSerialized
}

func (ent NetworkEntity) GetTypeAndKey() storage2.TypeAndKey {
	return storage2.TypeAndKey{Type: ent.Type, Key: ent.Key}
}

func (ent NetworkEntity) isEntityWriteOperation() {}

type NetworkEntities []NetworkEntity

func (ne NetworkEntities) MakeByTK() NetworkEntitiesByTK {
	ret := make(map[storage2.TypeAndKey]NetworkEntity, len(ne))
	for _, ent := range ne {
		ret[ent.GetTypeAndKey()] = ent
	}
	return ret
}

// MakeByParentTK returns the network entities, keyed by the TK of their parent
// associations, once per parent association.
func (ne NetworkEntities) MakeByParentTK() map[storage2.TypeAndKey]NetworkEntities {
	ret := map[storage2.TypeAndKey]NetworkEntities{}
	for _, ent := range ne {
		for _, parentTK := range ent.ParentAssociations {
			ret[parentTK] = append(ret[parentTK], ent)
		}
	}
	return ret
}

func (ne NetworkEntities) TKs() storage2.TKs {
	var tks storage2.TKs
	for _, e := range ne {
		tks = append(tks, e.GetTypeAndKey())
	}
	return tks
}

func (ne NetworkEntities) GetFirst(typ string) (NetworkEntity, error) {
	for _, e := range ne {
		if e.Type == typ {
			return e, nil
		}
	}
	return NetworkEntity{}, fmt.Errorf("no network entity of type %s found in %v", typ, ne)
}

func (ne NetworkEntities) fromProtos(protos []*storage.NetworkEntity, serdes serde.Registry) (NetworkEntities, error) {
	var ents NetworkEntities
	for _, p := range protos {
		ent, err := (&NetworkEntity{}).fromProto(p, serdes)
		if err != nil {
			return nil, err
		}
		ents = append(ents, ent)
	}
	return ents, nil
}

func (ne NetworkEntities) fromProtosSerialized(protos []*storage.NetworkEntity) NetworkEntities {
	var ents NetworkEntities
	for _, p := range protos {
		ents = append(ents, (&NetworkEntity{}).fromProtoSerialized(p))
	}
	return ents
}

type NetworkEntitiesByTK map[storage2.TypeAndKey]NetworkEntity

// Merge returns A+B when called as A.Merge(B).
// B overwrites any shared keys with A.
func (n NetworkEntitiesByTK) Merge(nn NetworkEntitiesByTK) NetworkEntitiesByTK {
	merged := NetworkEntitiesByTK{}
	for tk, ent := range n {
		merged[tk] = ent
	}
	for tk, ent := range nn {
		merged[tk] = ent
	}
	return merged
}

// Filter for TK type.
func (n NetworkEntitiesByTK) Filter(typ string) NetworkEntitiesByTK {
	filtered := NetworkEntitiesByTK{}
	for tk, ent := range n {
		if typ == tk.Type {
			filtered[tk] = ent
		}
	}
	return filtered
}

// MultiFilter for TK types.
func (n NetworkEntitiesByTK) MultiFilter(types ...string) NetworkEntitiesByTK {
	filtered := NetworkEntitiesByTK{}
	for tk, ent := range n {
		if funk.ContainsString(types, tk.Type) {
			filtered[tk] = ent
		}
	}
	return filtered
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

// FromProto converts a proto EntityGraph to its native counterpart.
func (eg EntityGraph) FromProto(protoGraph *storage.EntityGraph, serdes serde.Registry) (EntityGraph, error) {
	eg.Entities = make([]NetworkEntity, 0, len(protoGraph.Entities))
	for _, protoEnt := range protoGraph.Entities {
		ent, err := (NetworkEntity{}).fromProtoWithDefault(protoEnt, serdes)
		if err != nil {
			return eg, errors.Wrapf(err, "failed to deserialize entity %s", protoEnt.GetTypeAndKey())
		}
		eg.Entities = append(eg.Entities, ent)
	}

	eg.RootEntities = entIDsToTKs(protoGraph.RootEntities)

	eg.Edges = make([]GraphEdge, 0, len(protoGraph.Edges))
	for _, protoEdge := range protoGraph.Edges {
		eg.Edges = append(eg.Edges, (GraphEdge{}).fromProto(protoEdge))
	}
	return eg, nil
}

// ToProto converts an EntityGraph struct to it's proto representation.
func (eg EntityGraph) ToProto(serdes serde.Registry) (*storage.EntityGraph, error) {
	protoGraph := &storage.EntityGraph{}
	for _, ent := range eg.Entities {
		protoEnt, err := ent.toProto(serdes)
		if err != nil {
			return protoGraph, errors.Wrapf(err, "failed to convert entity %s to storage proto", ent.GetTypeAndKey())
		}
		protoGraph.Entities = append(protoGraph.Entities, protoEnt)
	}
	protoGraph.RootEntities = tksToEntIDs(eg.RootEntities)

	for _, edge := range eg.Edges {
		protoGraph.Edges = append(protoGraph.Edges, edge.toProto())
	}
	return protoGraph, nil
}

type GraphEdge struct {
	From storage2.TypeAndKey
	To   storage2.TypeAndKey
}

func (ge GraphEdge) fromProto(protoEdge *storage.GraphEdge) GraphEdge {
	ge.From = protoEdge.From.ToTypeAndKey()
	ge.To = protoEdge.To.ToTypeAndKey()
	return ge
}

func (ge GraphEdge) toProto() *storage.GraphEdge {
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

	// The following parameters allow for pagination of entity loads. These
	// load criteria parameters should be used in combination with one another.

	// PageSize is the maximum number of entities returned per load.
	PageSize uint32
	// NextPageToken is an opaque token provided to load the next page of
	// entities.
	PageToken string
}

func (elc EntityLoadCriteria) toProto() *storage.EntityLoadCriteria {
	return &storage.EntityLoadCriteria{
		LoadMetadata:       elc.LoadMetadata,
		LoadConfig:         elc.LoadConfig,
		LoadAssocsToThis:   elc.LoadAssocsToThis,
		LoadAssocsFromThis: elc.LoadAssocsFromThis,
		PageSize:           elc.PageSize,
		PageToken:          elc.PageToken,
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
	// NextPageToken is an opaque token provided to load the next page of
	// entities.
	NextPageToken string
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

func (euc EntityUpdateCriteria) toProto(serdes serde.Registry) (*storage.EntityUpdateCriteria, error) {
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
		// AssociationsToSet overrides AssociationsToAdd, so this check
		// prevents accidentally mixing the two fields
		if len(euc.AssociationsToAdd) != 0 {
			return nil, errors.New("cannot both add and set associations in the same EntityUpdateCriteria")
		}
		ret.AssociationsToSet = &storage.EntityAssociationsToSet{
			AssociationsToSet: tksToEntIDs(euc.AssociationsToSet),
		}
	}

	if euc.NewConfig != nil {
		bConfig, err := serde.Serialize(euc.NewConfig, euc.Type, serdes)
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

func marshalConfigs(configs map[string]interface{}, serdes serde.Registry) (map[string][]byte, error) {
	ret := map[string][]byte{}
	for configType, iConfig := range configs {
		if iConfig == nil {
			continue
		}

		sConfig, err := serde.Serialize(iConfig, configType, serdes)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to serialize config %s", configType)
		}
		ret[configType] = sConfig
	}
	return ret, nil
}

func unmarshalConfigs(configs map[string][]byte, serdes serde.Registry) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	for typ, config := range configs {
		// Skip unrecognized network config types, since one network can hold
		// configs from multiple network types
		if !serde.HasSerde(serdes, typ) {
			continue
		}
		iConfig, err := serde.Deserialize(config, typ, serdes)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to deserialize network config %s", typ)
		}
		ret[typ] = iConfig
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
