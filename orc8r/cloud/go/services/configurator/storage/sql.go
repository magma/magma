/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"magma/orc8r/cloud/go/sql_utils"
	"magma/orc8r/cloud/go/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	networksTable      = "cfg_networks"
	networkConfigTable = "cfg_network_configs"

	entityTable      = "cfg_entities"
	entityAssocTable = "cfg_assocs"
	entityAclTable   = "cfg_acls"
)

const (
	wildcardAllString = "*"
)

type IDGenerator interface {
	New() string
}

type DefaultIDGenerator struct{}

func (*DefaultIDGenerator) New() string {
	return uuid.New().String()
}

// NewSQLConfiguratorStorageFactory returns a ConfiguratorStorageFactory
// implementation backed by a SQL database.
func NewSQLConfiguratorStorageFactory(db *sql.DB, generator IDGenerator, sqlBuilder sql_utils.StatementBuilder) ConfiguratorStorageFactory {
	return &sqlConfiguratorStorageFactory{db: db, idGenerator: generator, builder: sqlBuilder}
}

type sqlConfiguratorStorageFactory struct {
	db          *sql.DB
	idGenerator IDGenerator
	builder     sql_utils.StatementBuilder
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

	networksTableExec := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			version INTEGER NOT NULL DEFAULT 0
		)
	`, networksTable)

	networksConfigTableExec := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			network_id TEXT REFERENCES %s (id) ON DELETE CASCADE,
			type TEXT NOT NULL,
			value BYTEA,

			PRIMARY KEY (network_id, type)
		)
	`, networkConfigTable, networksTable)

	// Create an internal-only primary key (UUID) for entities.
	// This keeps index size in control and supporting table schemas simpler.
	entityTableExec := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			pk TEXT PRIMARY KEY,

			network_id TEXT REFERENCES %s (id) ON DELETE CASCADE,
			type TEXT NOT NULL,
			key TEXT NOT NULL,

			graph_id TEXT NOT NULL,

			name TEXT,
			description TEXT,
			physical_id TEXT,
			config BYTEA,

			version INTEGER NOT NULL DEFAULT 0,

			UNIQUE (network_id, key, type)			
		)
	`, entityTable, networksTable)

	entityAssocTableExec := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			from_pk TEXT REFERENCES %s (pk) ON DELETE CASCADE,
			to_pk TEXT REFERENCES %s (pk) ON DELETE CASCADE,

			PRIMARY KEY (from_pk, to_pk)
		)
	`, entityAssocTable, entityTable, entityTable)

	entityAclTableExec := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s
		(
			id TEXT PRIMARY KEY,
			entity_pk TEXT REFERENCES %s (pk) ON DELETE CASCADE,

			scope text NOT NULL,
			permission INTEGER NOT NULL,
			type text NOT NULL,
			id_filter TEXT,

			version INTEGER NOT NULL DEFAULT 0
		)
	`, entityAclTable, entityTable)

	// Named return value for err so we can automatically decide to
	// commit/rollback
	tablesToCreate := []string{
		networksTableExec,
		networksConfigTableExec,
		entityTableExec,
		entityAssocTableExec,
		entityAclTableExec,
	}
	for _, execQuery := range tablesToCreate {
		_, err = tx.Exec(execQuery)
		if err != nil {
			return
		}
	}

	// Create indexes (index is not implicitly created on a referencing FK)
	indexesToCreate := []string{
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS graph_id_idx ON %s (graph_id)", entityTable),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS acl_ent_pk_idx ON %s (entity_pk)", entityAclTable),
	}
	for _, execQuery := range indexesToCreate {
		_, err = tx.Exec(execQuery)
		if err != nil {
			return
		}
	}

	// Create internal network(s)
	_, err = fact.builder.Insert(networksTable).
		Columns("id", "name", "description").
		Values(InternalNetworkID, internalNetworkName, internalNetworkDescription).
		OnConflict(nil, "id").
		RunWith(tx).
		Exec()
	if err != nil {
		err = fmt.Errorf("error creating internal networks: %s", err)
		return
	}

	return
}

func (fact *sqlConfiguratorStorageFactory) StartTransaction(ctx context.Context, opts *TxOptions) (ConfiguratorStorage, error) {
	tx, err := fact.db.BeginTx(ctx, getSqlOpts(opts))
	if err != nil {
		return nil, err
	}
	return &sqlConfiguratorStorage{tx: tx, idGenerator: fact.idGenerator, builder: fact.builder}, nil
}

func getSqlOpts(opts *TxOptions) *sql.TxOptions {
	if opts == nil {
		return nil
	}
	return &sql.TxOptions{ReadOnly: opts.ReadOnly}
}

type sqlConfiguratorStorage struct {
	tx          *sql.Tx
	idGenerator IDGenerator
	builder     sql_utils.StatementBuilder
}

func (store *sqlConfiguratorStorage) Commit() error {
	return store.tx.Commit()
}

func (store *sqlConfiguratorStorage) Rollback() error {
	return store.tx.Rollback()
}

func (store *sqlConfiguratorStorage) LoadNetworks(ids []string, loadCriteria NetworkLoadCriteria) (NetworkLoadResult, error) {
	emptyRet := NetworkLoadResult{NetworkIDsNotFound: []string{}, Networks: []Network{}}
	if len(ids) == 0 {
		return emptyRet, nil
	}

	selectBuilder := store.builder.Select(getNetworkQueryColumns(loadCriteria)...).
		From(networksTable).
		Where(sq.Eq{
			fmt.Sprintf("%s.id", networksTable): ids,
		})
	if loadCriteria.LoadConfigs {
		selectBuilder = selectBuilder.LeftJoin(
			fmt.Sprintf(
				"%s ON %s.network_id = %s.id",
				networkConfigTable, networkConfigTable, networksTable,
			),
		)
	}
	rows, err := selectBuilder.RunWith(store.tx).Query()
	if err != nil {
		return emptyRet, fmt.Errorf("error querying for networks: %s", err)
	}
	defer sql_utils.CloseRowsLogOnError(rows, "LoadNetworks")

	loadedNetworksByID, loadedNetworkIDs, err := scanNetworkRows(rows, loadCriteria)
	if err != nil {
		return emptyRet, err
	}

	ret := NetworkLoadResult{
		NetworkIDsNotFound: getNetworkIDsNotFound(loadedNetworksByID, ids),
		Networks:           make([]Network, 0, len(loadedNetworksByID)),
	}
	for _, nid := range loadedNetworkIDs {
		ret.Networks = append(ret.Networks, *loadedNetworksByID[nid])
	}
	return ret, nil
}

func (store *sqlConfiguratorStorage) LoadAllNetworks(loadCriteria NetworkLoadCriteria) ([]Network, error) {
	emptyNetworks := []Network{}
	idsToExclude := []string{InternalNetworkID}

	selectBuilder := store.builder.Select(getNetworkQueryColumns(loadCriteria)...).
		From(networksTable).
		Where(sq.NotEq{
			fmt.Sprintf("%s.id", networksTable): idsToExclude,
		})
	if loadCriteria.LoadConfigs {
		selectBuilder = selectBuilder.LeftJoin(
			fmt.Sprintf(
				"%s ON %s.network_id = %s.id",
				networkConfigTable, networkConfigTable, networksTable,
			),
		)
	}
	rows, err := selectBuilder.RunWith(store.tx).Query()
	if err != nil {
		return emptyNetworks, fmt.Errorf("error querying for networks: %s", err)
	}
	defer sql_utils.CloseRowsLogOnError(rows, "LoadAllNetworks")

	loadedNetworksByID, loadedNetworkIDs, err := scanNetworkRows(rows, loadCriteria)
	if err != nil {
		return emptyNetworks, err
	}

	networks := make([]Network, 0, len(loadedNetworksByID))
	for _, nid := range loadedNetworkIDs {
		networks = append(networks, *loadedNetworksByID[nid])
	}
	return networks, nil
}

func (store *sqlConfiguratorStorage) CreateNetwork(network Network) (Network, error) {
	exists, err := store.doesNetworkExist(network.ID)
	if err != nil {
		return network, err
	}
	if exists {
		return network, fmt.Errorf("a network with ID %s already exists", network.ID)
	}

	_, err = store.builder.Insert(networksTable).
		Columns("id", "name", "description").
		Values(network.ID, network.Name, network.Description).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return network, fmt.Errorf("error inserting network: %s", err)
	}

	if funk.IsEmpty(network.Configs) {
		return network, nil
	}

	// Sort config keys for deterministic behavior
	configKeys := funk.Keys(network.Configs).([]string)
	sort.Strings(configKeys)
	insertBuilder := store.builder.Insert(networkConfigTable).
		Columns("network_id", "type", "value")
	for _, configKey := range configKeys {
		insertBuilder = insertBuilder.Values(network.ID, configKey, network.Configs[configKey])
	}
	_, err = insertBuilder.RunWith(store.tx).Exec()
	if err != nil {
		return network, errors.Wrap(err, "error inserting network configs")
	}

	return network, nil
}

func (store *sqlConfiguratorStorage) UpdateNetworks(updates []NetworkUpdateCriteria) error {
	if err := validateNetworkUpdates(updates); err != nil {
		return err
	}

	networksToDelete := []string{}
	networksToUpdate := []NetworkUpdateCriteria{}
	for _, update := range updates {
		if update.DeleteNetwork {
			networksToDelete = append(networksToDelete, update.ID)
		} else {
			networksToUpdate = append(networksToUpdate, update)
		}
	}

	stmtCache := sq.NewStmtCache(store.tx)
	defer sql_utils.ClearStatementCacheLogOnError(stmtCache)

	// Update networks first
	for _, update := range networksToUpdate {
		err := store.updateNetwork(update, stmtCache)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// Then delete all networks requested for deletion
	_, err := store.builder.Delete(networksTable).Where(sq.Eq{"id": networksToDelete}).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete networks")
	}
	return nil
}

func (store *sqlConfiguratorStorage) LoadEntities(networkID string, filter EntityLoadFilter, loadCriteria EntityLoadCriteria) (EntityLoadResult, error) {
	ret := EntityLoadResult{Entities: []NetworkEntity{}, EntitiesNotFound: []storage.TypeAndKey{}}

	// We load the requested entities in 3 steps:
	// First, we load the entities and their ACLs
	// Then, we load assocs if requested by the load criteria. Note that the
	// load criteria can specify to load edges to and/or from the requested
	// entities.
	// For each loaded edge, we need to load the (type, key) corresponding to
	// to the PK pair that an edge is represented as. These may be already
	// loaded as part of the first load from the entities table, so we can
	// be smart here and only load (type, key) for PKs which we don't know.
	// Finally, we will update the entity objects to return with their edges.

	entsByPk, err := store.loadFromEntitiesTable(networkID, filter, loadCriteria)
	if err != nil {
		return ret, err
	}
	assocs, allAssocPks, err := store.loadFromAssocsTable(filter, loadCriteria, entsByPk)
	if err != nil {
		return ret, err
	}
	entTksByPk, err := store.loadEntityTypeAndKeys(allAssocPks, entsByPk)
	if err != nil {
		return ret, err
	}

	entsByPk, _, err = updateEntitiesWithAssocs(entsByPk, assocs, entTksByPk, loadCriteria)
	if err != nil {
		return ret, err
	}

	for _, ent := range entsByPk {
		ret.Entities = append(ret.Entities, *ent)
	}
	ret.EntitiesNotFound = calculateEntitiesNotFound(entsByPk, filter.IDs)

	// Sort entities for deterministic returns
	entComparator := func(a, b NetworkEntity) bool {
		return storage.TypeAndKey{Type: a.Type, Key: a.Key}.String() < storage.TypeAndKey{Type: b.Type, Key: b.Key}.String()
	}
	sort.Slice(ret.Entities, func(i, j int) bool { return entComparator(ret.Entities[i], ret.Entities[j]) })

	return ret, nil
}

func (store *sqlConfiguratorStorage) CreateEntity(networkID string, entity NetworkEntity) (NetworkEntity, error) {
	exists, err := store.doesEntExist(networkID, entity.GetTypeAndKey())
	if err != nil {
		return NetworkEntity{}, err
	}
	if exists {
		return NetworkEntity{}, fmt.Errorf("an entity (%s) already exists", entity.GetTypeAndKey())
	}

	// First, we insert the entity and its ACLs. We do this first so we have a
	// pk for the entity to reference in edge creation.
	// Then we insert the associations as graph edges. This step involves a
	// lookup of the associated entities to retrieve their PKs (since we don't
	// expose PK to the world).
	// Finally, if the created entity "bridges" 1 or more graphs, we merge
	// those graphs into a single graph.
	// For simplicity, we don't do any cycle detection at the moment. This
	// shouldn't be a problem on the load side because we load graphs via
	// graph ID, not by traversing edges.

	createdEntWithPk, err := store.insertIntoEntityTable(networkID, entity)
	if err != nil {
		return NetworkEntity{}, err
	}

	err = store.createPermissions(networkID, createdEntWithPk.pk, createdEntWithPk.Permissions)
	if err != nil {
		return NetworkEntity{}, err
	}

	allAssociatedEntsByTk, err := store.createEdges(networkID, createdEntWithPk)
	if err != nil {
		return NetworkEntity{}, err
	}

	newGraphID, err := store.mergeGraphs(createdEntWithPk, allAssociatedEntsByTk)
	if err != nil {
		return NetworkEntity{}, err
	}
	createdEntWithPk.GraphID = newGraphID

	// If we were given duplicate edges, get rid of those
	createdEntWithPk.Associations = funk.Uniq(createdEntWithPk.Associations).([]storage.TypeAndKey)

	return createdEntWithPk.NetworkEntity, nil
}

func (store *sqlConfiguratorStorage) UpdateEntity(networkID string, update EntityUpdateCriteria) (NetworkEntity, error) {
	emptyRet := NetworkEntity{Type: update.Type, Key: update.Key}
	entToUpdate, err := store.loadEntToUpdate(networkID, update)
	if err != nil {
		return emptyRet, errors.Wrap(err, "failed to load entity being updated")
	}

	if update.DeleteEntity {
		// Cascading FK relations in the schema will handle the other tables
		exec := fmt.Sprintf("DELETE FROM %s WHERE (network_id, type, key) = ($1, $2, $3)", entityTable)
		_, err := store.tx.Exec(exec, networkID, update.Type, update.Key)
		if err != nil {
			return emptyRet, errors.Wrapf(err, "failed to delete entity (%s, %s)", update.Type, update.Key)
		}

		// Deleting a node could partition its graph
		err = store.fixGraph(networkID, entToUpdate.GraphID, &entToUpdate)
		if err != nil {
			return emptyRet, errors.Wrap(err, "failed to fix entity graph after deletion")
		}

		return emptyRet, nil
	}

	// Then, update the fields on the entity table
	err = store.processEntityFieldsUpdate(entToUpdate.pk, update, &entToUpdate.NetworkEntity)
	if err != nil {
		return entToUpdate.NetworkEntity, errors.WithStack(err)
	}

	// Next, update permissions
	err = store.processPermissionUpdates(networkID, entToUpdate.pk, update, &entToUpdate.NetworkEntity)
	if err != nil {
		return entToUpdate.NetworkEntity, errors.WithStack(err)
	}

	// Finally, process edge updates for the graph
	err = store.processEdgeUpdates(networkID, update, &entToUpdate)
	if err != nil {
		return entToUpdate.NetworkEntity, errors.WithStack(err)
	}

	return entToUpdate.NetworkEntity, nil
}

func (store *sqlConfiguratorStorage) LoadGraphForEntity(networkID string, entityID storage.TypeAndKey, loadCriteria EntityLoadCriteria) (EntityGraph, error) {
	// Technically you could do this in one DB query with a subquery in the
	// WHERE when selecting from the entity table.
	// But until we hit some kind of scaling limit, let's keep the code simple
	// and delegate to LoadGraph after loading the requested entity.

	// We just care about getting the graph ID off this entity so use an empty
	// load criteria
	loadResult, err := store.loadSpecificEntities(networkID, EntityLoadFilter{IDs: []storage.TypeAndKey{entityID}}, EntityLoadCriteria{})
	if err != nil {
		return EntityGraph{}, errors.Wrap(err, "failed to load entity for graph query")
	}

	var loadedEnt *NetworkEntity
	for _, ent := range loadResult {
		loadedEnt = ent
	}
	if loadedEnt == nil {
		return EntityGraph{}, errors.Errorf("could not find requested entity (%s) for graph query", entityID)
	}

	internalGraph, err := store.loadGraphInternal(networkID, loadedEnt.GraphID, loadCriteria)
	if err != nil {
		return EntityGraph{}, errors.WithStack(err)
	}

	rootPks := findRootNodes(internalGraph)
	if funk.IsEmpty(rootPks) {
		return EntityGraph{}, errors.Errorf("graph does not have root nodes because it is a ring")
	}

	// Fill entities with assocs. We will always fill out both directions of
	// associations so we'll alter the load criteria for the helper function.
	entTksByPk := funk.Map(
		internalGraph.entsByPk,
		func(pk string, ent *NetworkEntity) (string, storage.TypeAndKey) { return pk, ent.GetTypeAndKey() },
	).(map[string]storage.TypeAndKey)
	loadCriteria.LoadAssocsToThis, loadCriteria.LoadAssocsFromThis = true, true
	_, edges, err := updateEntitiesWithAssocs(internalGraph.entsByPk, internalGraph.edges, entTksByPk, loadCriteria)
	if err != nil {
		return EntityGraph{}, errors.Wrap(err, "failed to construct graph after loading")
	}

	// To make testing easier, we'll order the returned entities by TK
	retEnts := funk.Map(internalGraph.entsByPk, func(_ string, ent *NetworkEntity) NetworkEntity { return *ent }).([]NetworkEntity)
	retRoots := funk.Map(rootPks, func(pk string) storage.TypeAndKey { return entTksByPk[pk] }).([]storage.TypeAndKey)
	sort.Slice(retEnts, func(i, j int) bool {
		return storage.IsTKLessThan(retEnts[i].GetTypeAndKey(), retEnts[j].GetTypeAndKey())
	})
	sort.Slice(retRoots, func(i, j int) bool { return storage.IsTKLessThan(retRoots[i], retRoots[j]) })

	return EntityGraph{
		Entities:     retEnts,
		RootEntities: retRoots,
		Edges:        edges,
	}, nil
}
