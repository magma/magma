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
	"database/sql"
	"fmt"
	"os"
	"sort"

	sq "github.com/Masterminds/squirrel"
	"github.com/thoas/go-funk"
	"google.golang.org/protobuf/proto"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
)

const (
	networksTable      = "cfg_networks"
	networkConfigTable = "cfg_network_configs"

	entityTable      = "cfg_entities"
	entityAssocTable = "cfg_assocs"
)

const (
	nwIDCol   = "id"
	nwTypeCol = "type"
	nwNameCol = "name"
	nwDescCol = "description"
	nwVerCol  = "version"

	nwcIDCol   = "network_id"
	nwcTypeCol = "type"
	nwcValCol  = "value"

	entPkCol   = "pk"
	entNidCol  = "network_id"
	entTypeCol = "type"
	entKeyCol  = "\"key\""
	entGidCol  = "graph_id"
	entNameCol = "name"
	entDescCol = "description"
	entPidCol  = "physical_id"
	entConfCol = "config"
	entVerCol  = "version"

	aFrCol = "from_pk"
	aToCol = "to_pk"
)

// NewSQLConfiguratorStorageFactory returns a ConfiguratorStorageFactory
// implementation backed by a SQL database.
func NewSQLConfiguratorStorageFactory(db *sql.DB, generator storage.IDGenerator, sqlBuilder sqorc.StatementBuilder, maxEntityLoadSize uint32) ConfiguratorStorageFactory {
	return &sqlConfiguratorStorageFactory{db: db, idGenerator: generator, builder: sqlBuilder, maxEntityLoadSize: maxEntityLoadSize}
}

type sqlConfiguratorStorageFactory struct {
	db                *sql.DB
	idGenerator       storage.IDGenerator
	builder           sqorc.StatementBuilder
	maxEntityLoadSize uint32
}

func (fact *sqlConfiguratorStorageFactory) InitializeServiceStorage() (err error) {
	tx, err := fact.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("%s; rollback error: %s", err, rollbackErr)
			}
		}
	}()

	// Named return values below so we can automatically decide tx commit/
	// rollback in deferred function

	_, err = fact.builder.CreateTable(networksTable).
		IfNotExists().
		Column(nwIDCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(nwTypeCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(nwNameCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(nwDescCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(nwVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		RunWith(tx).
		Exec()
	if err != nil {
		err = fmt.Errorf("failed to create networks table: %w", err)
		return
	}

	// Adding a type column if it doesn't exist already. This will ensure network
	// tables that are already created will also have the type column.
	// TODO Remove after 1-2 months to ensure service isn't disrupted
	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s text", networksTable, nwTypeCol))
	// special case sqlite3 because ADD COLUMN IF NOT EXISTS is not supported
	// and we only run sqlite3 for unit tests
	if err != nil && os.Getenv("SQL_DRIVER") != "sqlite3" {
		err = fmt.Errorf("failed to add 'type' field to networks table: %w", err)
	}

	_, err = fact.builder.CreateIndex("type_idx").
		IfNotExists().
		On(networksTable).
		Columns(nwTypeCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = fmt.Errorf("failed to create network type index: %w", err)
		return
	}

	_, err = fact.builder.CreateTable(networkConfigTable).
		IfNotExists().
		Column(nwcIDCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(nwcTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(nwcValCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		PrimaryKey(nwcIDCol, nwcTypeCol).
		ForeignKey(networksTable, map[string]string{nwcIDCol: nwIDCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		err = fmt.Errorf("failed to create network configs table: %w", err)
		return
	}

	// Create an internal-only primary key (UUID) for entities.
	// This keeps index size in control and supporting table schemas simpler.
	_, err = fact.builder.CreateTable(entityTable).
		IfNotExists().
		Column(entPkCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(entNidCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(entTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(entKeyCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(entGidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(entNameCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(entDescCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(entPidCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(entConfCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(entVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		Unique(entNidCol, entKeyCol, entTypeCol).
		Unique(entPidCol).
		ForeignKey(networksTable, map[string]string{entNidCol: nwIDCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		err = fmt.Errorf("failed to create entities table: %w", err)
		return
	}

	_, err = fact.builder.CreateTable(entityAssocTable).
		IfNotExists().
		Column(aFrCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(aToCol).Type(sqorc.ColumnTypeText).EndColumn().
		PrimaryKey(aFrCol, aToCol).
		ForeignKey(entityTable, map[string]string{aFrCol: entPkCol}, sqorc.ColumnOnDeleteCascade).
		ForeignKey(entityTable, map[string]string{aToCol: entPkCol}, sqorc.ColumnOnDeleteCascade).
		RunWith(tx).
		Exec()
	if err != nil {
		err = fmt.Errorf("failed to create entity assoc table: %w", err)
		return
	}

	// Create indexes (index is not implicitly created on a referencing FK)
	_, err = fact.builder.CreateIndex("graph_id_idx").
		IfNotExists().
		On(entityTable).
		Columns(entGidCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = fmt.Errorf("failed to create graph ID index: %w", err)
		return
	}

	// Create internal network(s)
	_, err = fact.builder.Insert(networksTable).
		Columns(nwIDCol, nwTypeCol, nwNameCol, nwDescCol).
		Values(InternalNetworkID, internalNetworkType, internalNetworkName, internalNetworkDescription).
		OnConflict(nil, nwIDCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = fmt.Errorf("error creating internal networks: %w", err)
		return
	}

	return
}

func (fact *sqlConfiguratorStorageFactory) StartTransaction(ctx context.Context, opts *storage.TxOptions) (ConfiguratorStorage, error) {
	tx, err := fact.db.BeginTx(ctx, getSqlOpts(opts))
	if err != nil {
		return nil, err
	}
	return &sqlConfiguratorStorage{tx: tx, idGenerator: fact.idGenerator, builder: fact.builder, maxEntityLoadSize: fact.maxEntityLoadSize}, nil
}

func getSqlOpts(opts *storage.TxOptions) *sql.TxOptions {
	if opts == nil {
		return nil
	}
	if opts.Isolation == 0 {
		return &sql.TxOptions{ReadOnly: opts.ReadOnly}
	}
	return &sql.TxOptions{ReadOnly: opts.ReadOnly, Isolation: sql.IsolationLevel(opts.Isolation)}
}

type sqlConfiguratorStorage struct {
	tx                *sql.Tx
	idGenerator       storage.IDGenerator
	builder           sqorc.StatementBuilder
	maxEntityLoadSize uint32
}

func (store *sqlConfiguratorStorage) Commit() error {
	return store.tx.Commit()
}

func (store *sqlConfiguratorStorage) Rollback() error {
	return store.tx.Rollback()
}

func (store *sqlConfiguratorStorage) LoadNetworks(filter *NetworkLoadFilter, loadCriteria *NetworkLoadCriteria) (*NetworkLoadResult, error) {
	filterCopy := proto.Clone(filter).(*NetworkLoadFilter)
	loadCriteriaCopy := proto.Clone(loadCriteria).(*NetworkLoadCriteria)
	emptyRet := &NetworkLoadResult{NetworkIDsNotFound: []string{}, Networks: []*Network{}}
	if funk.IsEmpty(filterCopy.Ids) && funk.IsEmpty(filterCopy.TypeFilter) {
		return emptyRet, nil
	}

	selectBuilder := store.getLoadNetworksSelectBuilder(filterCopy, loadCriteriaCopy)
	if loadCriteriaCopy.LoadConfigs {
		selectBuilder = selectBuilder.LeftJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				networkConfigTable, networkConfigTable, nwcIDCol, networksTable, nwIDCol,
			),
		)
	}
	rows, err := selectBuilder.RunWith(store.tx).Query()
	if err != nil {
		return emptyRet, fmt.Errorf("error querying for networks: %s", err)
	}
	defer sqorc.CloseRowsLogOnError(rows, "LoadNetworks")

	loadedNetworksByID, loadedNetworkIDs, err := scanNetworkRows(rows, loadCriteriaCopy)
	if err != nil {
		return emptyRet, err
	}

	ret := &NetworkLoadResult{
		NetworkIDsNotFound: getNetworkIDsNotFound(loadedNetworksByID, filterCopy.Ids),
		Networks:           make([]*Network, 0, len(loadedNetworksByID)),
	}
	for _, nid := range loadedNetworkIDs {
		ret.Networks = append(ret.Networks, loadedNetworksByID[nid])
	}
	return ret, nil
}

func (store *sqlConfiguratorStorage) LoadAllNetworks(loadCriteria *NetworkLoadCriteria) ([]*Network, error) {
	emptyNetworks := []*Network{}
	idsToExclude := []string{InternalNetworkID}

	loadCriteriaCopy := proto.Clone(loadCriteria).(*NetworkLoadCriteria)
	selectBuilder := store.builder.Select(getNetworkQueryColumns(loadCriteriaCopy)...).
		From(networksTable).
		Where(sq.NotEq{
			fmt.Sprintf("%s.%s", networksTable, nwIDCol): idsToExclude,
		})
	if loadCriteriaCopy.LoadConfigs {
		selectBuilder = selectBuilder.LeftJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				networkConfigTable, networkConfigTable, nwcIDCol, networksTable, nwIDCol,
			),
		)
	}
	rows, err := selectBuilder.RunWith(store.tx).Query()
	if err != nil {
		return emptyNetworks, fmt.Errorf("error querying for networks: %s", err)
	}
	defer sqorc.CloseRowsLogOnError(rows, "LoadAllNetworks")

	loadedNetworksByID, loadedNetworkIDs, err := scanNetworkRows(rows, loadCriteriaCopy)
	if err != nil {
		return emptyNetworks, err
	}

	networks := make([]*Network, 0, len(loadedNetworksByID))
	for _, nid := range loadedNetworkIDs {
		networks = append(networks, loadedNetworksByID[nid])
	}
	return networks, nil
}

func (store *sqlConfiguratorStorage) CreateNetwork(network *Network) (*Network, error) {
	networkCopy := proto.Clone(network).(*Network)
	exists, err := store.doesNetworkExist(networkCopy.ID)
	if err != nil {
		return &Network{}, err
	}
	if exists {
		return &Network{}, fmt.Errorf("a network with ID %s already exists", networkCopy.ID)
	}

	_, err = store.builder.Insert(networksTable).
		Columns(nwIDCol, nwTypeCol, nwNameCol, nwDescCol).
		Values(networkCopy.ID, networkCopy.Type, networkCopy.Name, networkCopy.Description).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return &Network{}, fmt.Errorf("error inserting network: %s", err)
	}

	if funk.IsEmpty(networkCopy.Configs) {
		return networkCopy, nil
	}

	// Sort config keys for deterministic behavior
	configKeys := funk.Keys(networkCopy.Configs).([]string)
	sort.Strings(configKeys)
	insertBuilder := store.builder.Insert(networkConfigTable).
		Columns(nwcIDCol, nwcTypeCol, nwcValCol)
	for _, configKey := range configKeys {
		insertBuilder = insertBuilder.Values(networkCopy.ID, configKey, networkCopy.Configs[configKey])
	}
	_, err = insertBuilder.RunWith(store.tx).Exec()
	if err != nil {
		return &Network{}, fmt.Errorf("error inserting network configs: %w", err)
	}

	return networkCopy, nil
}

func (store *sqlConfiguratorStorage) UpdateNetworks(updates []*NetworkUpdateCriteria) error {
	if err := validateNetworkUpdates(updates); err != nil {
		return err
	}

	networksToDelete := []string{}
	networksToUpdate := []*NetworkUpdateCriteria{}
	for _, update := range updates {
		if update.DeleteNetwork {
			networksToDelete = append(networksToDelete, update.ID)
		} else {
			networksToUpdate = append(networksToUpdate, update)
		}
	}

	stmtCache := sq.NewStmtCache(store.tx)
	defer sqorc.ClearStatementCacheLogOnError(stmtCache, "UpdateNetworks")

	// Update networks first
	for _, update := range networksToUpdate {
		err := store.updateNetwork(update, stmtCache)
		if err != nil {
			return err
		}
	}

	_, err := store.builder.Delete(networkConfigTable).Where(sq.Eq{nwcIDCol: networksToDelete}).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return fmt.Errorf("failed to delete configs associated with networks: %w", err)
	}
	_, err = store.builder.Delete(networksTable).Where(sq.Eq{nwIDCol: networksToDelete}).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return fmt.Errorf("failed to delete networks: %w", err)
	}
	return nil
}

func (store *sqlConfiguratorStorage) CountEntities(networkID string, filter *EntityLoadFilter) (*EntityCountResult, error) {
	filterCopy := proto.Clone(filter).(*EntityLoadFilter)
	ret := &EntityCountResult{Count: 0}
	count, err := store.countEntities(networkID, filterCopy)
	if err != nil {
		return ret, err
	}
	ret.Count = count
	return ret, nil
}

func (store *sqlConfiguratorStorage) LoadEntities(networkID string, filter *EntityLoadFilter, criteria *EntityLoadCriteria) (*EntityLoadResult, error) {
	filterCopy := proto.Clone(filter).(*EntityLoadFilter)
	criteriaCopy := proto.Clone(criteria).(*EntityLoadCriteria)
	if err := validatePaginatedLoadParameters(filterCopy, criteriaCopy); err != nil {
		return &EntityLoadResult{}, err
	}

	entsByTK, err := store.loadEntities(networkID, filterCopy, criteriaCopy)
	if err != nil {
		return &EntityLoadResult{}, err
	}

	if criteria.LoadAssocsFromThis {
		assocs, err := store.loadAssocs(networkID, filterCopy, criteriaCopy, loadChildren)
		if err != nil {
			return &EntityLoadResult{}, err
		}
		for _, assoc := range assocs {
			// Assoc may be from a not-loaded ent
			e, ok := entsByTK[assoc.fromTK]
			if ok {
				e.Associations = append(e.Associations, assoc.getToID())
			}
		}
		for _, ent := range entsByTK {
			SortIDs(ent.Associations) // for deterministic return
		}
	}

	if criteria.LoadAssocsToThis {
		parentAssocs, err := store.loadAssocs(networkID, filterCopy, criteriaCopy, loadParents)
		if err != nil {
			return &EntityLoadResult{}, err
		}
		for _, parentAssoc := range parentAssocs {
			// Assoc may be to a not-loaded ent
			e, ok := entsByTK[parentAssoc.toTK]
			if ok {
				e.ParentAssociations = append(e.ParentAssociations, parentAssoc.getFromID())
			}
		}
		for _, ent := range entsByTK {
			SortIDs(ent.ParentAssociations) // for deterministic return
		}
	}

	res := &EntityLoadResult{}

	for _, ent := range entsByTK {
		res.Entities = append(res.Entities, ent)
	}
	SortEntities(res.Entities) // for deterministic return
	res.EntitiesNotFound = calculateIDsNotFound(entsByTK, filterCopy.IDs)

	// Set next page token when there may be more pages to return
	if len(res.Entities) == store.getEntityLoadPageSize(criteriaCopy) {
		res.NextPageToken, err = getNextPageToken(res.Entities)
		if err != nil {
			return &EntityLoadResult{}, err
		}
	}

	return res, nil
}

func (store *sqlConfiguratorStorage) CreateEntity(networkID string, entity *NetworkEntity) (*NetworkEntity, error) {
	entityCopy := proto.Clone(entity).(*NetworkEntity)
	exists, err := store.doesEntExist(networkID, entityCopy.GetTK())
	if err != nil {
		return &NetworkEntity{}, err
	}
	if exists {
		return &NetworkEntity{}, fmt.Errorf("an entity '%s' already exists", entityCopy.GetTK())
	}

	// Physical ID must be unique across all networks, since we use a gateway's
	// physical ID to search for its network (and ent)
	physicalIDExists, err := store.doesPhysicalIDExist(entityCopy.GetPhysicalID())
	if err != nil {
		return &NetworkEntity{}, err
	}
	if physicalIDExists {
		return &NetworkEntity{}, fmt.Errorf("an entity with physical ID '%s' already exists", entityCopy.GetPhysicalID())
	}

	// First insert the associations as graph edges. This step involves a
	// lookup of the associated entities to retrieve their PKs (since we don't
	// trust the provided PKs).
	// Finally, if the created entity "bridges" 1 or more graphs, we merge
	// those graphs into a single graph.
	// For simplicity, we don't do any cycle detection at the moment. This
	// shouldn't be a problem on the load side because we load graphs via
	// graph ID, not by traversing edges.

	createdEnt, err := store.insertIntoEntityTable(networkID, entityCopy)
	if err != nil {
		return &NetworkEntity{}, err
	}

	allAssociatedEntsByTk, err := store.createEdges(networkID, createdEnt)
	if err != nil {
		return &NetworkEntity{}, err
	}

	newGraphID, err := store.mergeGraphs(createdEnt, allAssociatedEntsByTk)
	if err != nil {
		return &NetworkEntity{}, err
	}
	createdEnt.GraphID = newGraphID

	// If we were given duplicate edges, get rid of those
	if funk.NotEmpty(createdEnt.Associations) {
		createdEnt.Associations = funk.Chain(createdEnt.Associations).
			Map(func(id *EntityID) storage.TK { return id.ToTK() }).
			Uniq().
			Map(func(tk storage.TK) *EntityID { return (&EntityID{}).FromTK(tk) }).
			Value().([]*EntityID)
	}

	createdEnt.NetworkID = networkID
	return createdEnt, nil
}

func (store *sqlConfiguratorStorage) UpdateEntity(networkID string, update *EntityUpdateCriteria) (*NetworkEntity, error) {
	updateCopy := proto.Clone(update).(*EntityUpdateCriteria)
	emptyRet := &NetworkEntity{Type: update.Type, Key: update.Key}
	entToUpdate, err := store.loadEntToUpdate(networkID, updateCopy)
	if err != nil && !updateCopy.DeleteEntity {
		return emptyRet, fmt.Errorf("failed to load entity being updated: %w", err)
	}
	if entToUpdate == nil {
		return emptyRet, nil
	}

	if updateCopy.DeleteEntity {
		// Cascading FK relations in the schema will handle the other tables
		_, err := store.builder.Delete(entityTable).
			Where(sq.And{
				sq.Eq{entNidCol: networkID},
				sq.Eq{entTypeCol: updateCopy.Type},
				sq.Eq{entKeyCol: updateCopy.Key},
			}).
			RunWith(store.tx).
			Exec()
		if err != nil {
			return emptyRet, fmt.Errorf("failed to delete entity (%s, %s): %w", updateCopy.Type, updateCopy.Key, err)
		}

		// Deleting a node could partition its graph
		err = store.fixGraph(networkID, entToUpdate.GraphID, entToUpdate)
		if err != nil {
			return emptyRet, fmt.Errorf("failed to fix entity graph after deletion: %w", err)
		}

		return emptyRet, nil
	}

	// Then, update the fields on the entity table
	entToUpdate.NetworkID = networkID
	err = store.processEntityFieldsUpdate(entToUpdate.Pk, updateCopy, entToUpdate)
	if err != nil {
		return entToUpdate, err
	}

	// Finally, process edge updates for the graph
	err = store.processEdgeUpdates(networkID, updateCopy, entToUpdate)
	if err != nil {
		return entToUpdate, err
	}

	return entToUpdate, nil
}

func (store *sqlConfiguratorStorage) LoadGraphForEntity(networkID string, entityID *EntityID, loadCriteria *EntityLoadCriteria) (*EntityGraph, error) {
	entityIDCopy := proto.Clone(entityID).(*EntityID)
	loadCriteriaCopy := proto.Clone(loadCriteria).(*EntityLoadCriteria)
	// We just care about getting the graph ID off this entity so use an empty
	// load criteria
	singleEnt, err := store.loadEntities(networkID, &EntityLoadFilter{IDs: []*EntityID{entityIDCopy}}, &EntityLoadCriteria{})
	if err != nil {
		return &EntityGraph{}, fmt.Errorf("failed to load entity for graph query: %w", err)
	}

	var ent *NetworkEntity
	for _, e := range singleEnt {
		ent = e
	}
	if ent == nil {
		return &EntityGraph{}, fmt.Errorf("could not find requested entity (%s) for graph query", entityIDCopy.String())
	}

	internalGraph, err := store.loadGraphInternal(networkID, ent.GraphID, loadCriteriaCopy)
	if err != nil {
		return &EntityGraph{}, err
	}

	rootPKs := findRootNodes(internalGraph)
	if funk.IsEmpty(rootPKs) {
		return &EntityGraph{}, fmt.Errorf("graph does not have root nodes")
	}

	edges, err := updateEntitiesWithAssocs(internalGraph.entsByTK, internalGraph.edges)
	if err != nil {
		return &EntityGraph{}, fmt.Errorf("failed to construct graph after loading: %w", err)
	}

	// To make testing easier, we'll order the returned entities by TK
	entsByPK := internalGraph.entsByTK.ByPK()
	retEnts := internalGraph.entsByTK.Ents()
	retRoots := funk.Map(rootPKs, func(pk string) *EntityID { return &EntityID{Type: entsByPK[pk].Type, Key: entsByPK[pk].Key} }).([]*EntityID)
	SortEntities(retEnts)
	SortIDs(retRoots)

	return &EntityGraph{
		Entities:     retEnts,
		RootEntities: retRoots,
		Edges:        edges,
	}, nil
}
