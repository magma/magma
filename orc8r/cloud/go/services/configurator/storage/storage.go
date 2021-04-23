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

package storage

import (
	"context"
	"fmt"
	"sort"

	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// ConfiguratorStorageFactory creates ConfiguratorStorage implementations bound
// to transactions.
type ConfiguratorStorageFactory interface {
	// InitializeServiceStorage should be called on service start to initialize
	// the tables that configurator storage implementations will depend on.
	InitializeServiceStorage() error

	// StartTransaction returns a ConfiguratorStorage implementation bound to
	// a transaction. Transaction options can be optionally provided.
	StartTransaction(ctx context.Context, opts *storage.TxOptions) (ConfiguratorStorage, error)
}

// ConfiguratorStorage is the interface for the configurator service's access
// to its data storage layer. Each ConfiguratorStorage instance is tied to a
// transaction within which all function calls operate.
type ConfiguratorStorage interface {

	// Commit commits the underlying transaction
	Commit() error

	// Rollback rolls back the underlying transaction
	Rollback() error

	// =======================================================================
	// Network Operations
	// =======================================================================

	// LoadNetworks returns a set of networks corresponding to the provided
	// load criteria. Any networks which aren't found are excluded from the
	// returned value.
	LoadNetworks(filter NetworkLoadFilter, loadCriteria NetworkLoadCriteria) (NetworkLoadResult, error)

	// LoadAllNetworks returns all networks registered
	LoadAllNetworks(loadCriteria NetworkLoadCriteria) ([]Network, error)

	// CreateNetwork creates a new network. The created network is returned.
	CreateNetwork(network Network) (Network, error)

	// UpdateNetworks updates a set of networks.
	UpdateNetworks(updates []NetworkUpdateCriteria) error

	// =======================================================================
	// Entity Operations
	// =======================================================================

	// LoadEntities returns a set of entities corresponding to the provided
	// load criteria. Any entities which aren't found are excluded from the
	// returned value.

	// Loads can be paginated by specifying a page size and token in the entity
	// load criteria. To exhaustively read all pages, clients must continue
	// querying until an empty page token is received in the load result.
	LoadEntities(networkID string, filter EntityLoadFilter, loadCriteria EntityLoadCriteria) (EntityLoadResult, error)

	// CountEntities returns the count of entities corresponding to the provided
	// load criteria.
	CountEntities(networkID string, filter EntityLoadFilter, loadCriteria EntityLoadCriteria) (EntityCountResult, error)

	// CreateEntity creates a new entity. The created entity is returned
	// with system-generated fields filled in.
	CreateEntity(networkID string, entity NetworkEntity) (NetworkEntity, error)

	// UpdateEntity updates an entity.
	// The updates to the specified entity will be returned as a NetworkEntity
	// object. Apart from identity fields, only fields which were updated will
	// be filled out, with system-generated IDs included.
	UpdateEntity(networkID string, update EntityUpdateCriteria) (NetworkEntity, error)

	// =======================================================================
	// Graph Operations
	// =======================================================================

	// LoadGraphForEntity returns the full DAG which contains the requested
	// entity. The load criteria fields on associations are ignored, and the
	// returned entities will always have both association fields filled out.
	LoadGraphForEntity(networkID string, entityID EntityID, loadCriteria EntityLoadCriteria) (EntityGraph, error)
}

// RollbackLogOnError calls Rollback on the provided ConfiguratorStorage and
// logs if Rollback resulted in an error.
func RollbackLogOnError(store ConfiguratorStorage) {
	err := store.Rollback()
	if err != nil {
		glog.Errorf("Error while rolling back tx: %+v", errors.WithStack(err))
	}
}

// CommitLogOnError calls Commit on the provided ConfiguratorStorage and logs
// if Commit resulted in an error.
func CommitLogOnError(store ConfiguratorStorage) {
	err := store.Commit()
	if err != nil {
		glog.Errorf("Error while committing tx: %+v", errors.WithStack(err))
	}
}

// InternalNetworkID is the ID of the network under which all non-tenant
// entities are organized under.
const InternalNetworkID = "network_magma_internal"
const internalNetworkType = "Internal"
const internalNetworkName = "Internal Magma Network"
const internalNetworkDescription = "Internal network to hold non-network entities"

// FullNetworkLoadCriteria is a utility variable to specify a full network load
var FullNetworkLoadCriteria = NetworkLoadCriteria{LoadMetadata: true, LoadConfigs: true}

func (m *EntityID) ToTypeAndKey() storage.TypeAndKey {
	return storage.TypeAndKey{Type: m.Type, Key: m.Key}
}

func (m *EntityID) FromTypeAndKey(tk storage.TypeAndKey) *EntityID {
	m.Type = tk.Type
	m.Key = tk.Key
	return m
}

func SortIDs(ids []*EntityID) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].ToTypeAndKey().IsLessThan(ids[j].ToTypeAndKey())
	})
}

func SortEntities(ents []*NetworkEntity) {
	sort.Slice(ents, func(i, j int) bool {
		return ents[i].GetTypeAndKey().String() < ents[j].GetTypeAndKey().String()
	})
}

func (m *NetworkEntity) GetID() *EntityID {
	return &EntityID{Type: m.Type, Key: m.Key}
}

func (m *NetworkEntity) GetTypeAndKey() storage.TypeAndKey {
	return m.GetID().ToTypeAndKey()
}

func (m NetworkEntity) GetGraphEdges() []*GraphEdge {
	myID := m.GetID()
	existingAssocs := map[storage.TypeAndKey]bool{}

	edges := make([]*GraphEdge, 0, len(m.Associations))
	for _, assoc := range m.Associations {
		if _, exists := existingAssocs[assoc.ToTypeAndKey()]; exists {
			continue
		}
		edges = append(edges, &GraphEdge{From: myID, To: assoc})
		existingAssocs[assoc.ToTypeAndKey()] = true
	}

	return edges
}

type EntitiesByPK map[string]*NetworkEntity

type EntitiesByTK map[storage.TypeAndKey]*NetworkEntity

func (e EntitiesByTK) ByPK() EntitiesByPK {
	byPK := make(map[string]*NetworkEntity, len(e))
	for _, ent := range e {
		byPK[ent.Pk] = ent
	}
	return byPK
}

func (e EntitiesByTK) PKs() []string {
	pks := make([]string, 0, len(e))
	for _, ent := range e {
		pks = append(pks, ent.Pk)
	}
	return pks
}

func (e EntitiesByTK) Ents() []*NetworkEntity {
	ents := make([]*NetworkEntity, 0, len(e))
	for _, ent := range e {
		ents = append(ents, ent)
	}
	return ents
}

// IsLoadAllEntities return true if the EntityLoadFilter is specifying to load
// all entities in a network, false if there are any filter conditions.
func (m *EntityLoadFilter) IsLoadAllEntities() bool {
	return m.TypeFilter == nil && m.KeyFilter == nil && m.GraphID == nil && funk.IsEmpty(m.IDs)
}

// FullEntityLoadCriteria is an EntityLoadCriteria which loads everything
var FullEntityLoadCriteria = EntityLoadCriteria{
	LoadMetadata:       true,
	LoadConfig:         true,
	LoadAssocsToThis:   true,
	LoadAssocsFromThis: true,
}

func (m *EntityUpdateCriteria) GetID() *EntityID {
	return &EntityID{Type: m.Type, Key: m.Key}
}

func (m *EntityUpdateCriteria) GetTypeAndKey() storage.TypeAndKey {
	return storage.TypeAndKey{Type: m.Type, Key: m.Key}
}

func (m *EntityUpdateCriteria) getEdgesToCreate() []*EntityID {
	if m.AssociationsToSet != nil {
		return m.AssociationsToSet.AssociationsToSet
	}
	return m.AssociationsToAdd
}

func (m *GraphEdge) ToString() string {
	return fmt.Sprintf("%s, %s", m.From.ToTypeAndKey(), m.To.ToTypeAndKey())
}
