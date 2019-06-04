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

	// Named return values below so we can automatically decide tx commit/
	// rollback in deferred function

	_, err = fact.builder.CreateTable(networksTable).
		IfNotExists().
		Column("id").Type(sql_utils.ColumnTypeText).PrimaryKey().EndColumn().
		Column("name").Type(sql_utils.ColumnTypeText).EndColumn().
		Column("description").Type(sql_utils.ColumnTypeText).EndColumn().
		Column("version").Type(sql_utils.ColumnTypeInt).NotNull().Default(0).EndColumn().
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create networks table")
		return
	}

	_, err = fact.builder.CreateTable(networkConfigTable).
		IfNotExists().
		Column("network_id").Type(sql_utils.ColumnTypeText).References(networksTable, "id").OnDelete(sql_utils.ColumnOnDeleteCascade).EndColumn().
		Column("type").Type(sql_utils.ColumnTypeText).NotNull().EndColumn().
		Column("value").Type(sql_utils.ColumnTypeBytes).EndColumn().
		PrimaryKey("network_id", "type").
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create network configs table")
		return
	}

	// Create an internal-only primary key (UUID) for entities.
	// This keeps index size in control and supporting table schemas simpler.
	_, err = fact.builder.CreateTable(entityTable).
		IfNotExists().
		Column("pk").Type(sql_utils.ColumnTypeText).PrimaryKey().EndColumn().
		Column("network_id").Type(sql_utils.ColumnTypeText).References(networksTable, "id").OnDelete(sql_utils.ColumnOnDeleteCascade).EndColumn().
		Column("type").Type(sql_utils.ColumnTypeText).NotNull().EndColumn().
		Column("key").Type(sql_utils.ColumnTypeText).NotNull().EndColumn().
		Column("graph_id").Type(sql_utils.ColumnTypeText).NotNull().EndColumn().
		Column("name").Type(sql_utils.ColumnTypeText).EndColumn().
		Column("description").Type(sql_utils.ColumnTypeText).EndColumn().
		Column("physical_id").Type(sql_utils.ColumnTypeText).EndColumn().
		Column("config").Type(sql_utils.ColumnTypeBytes).EndColumn().
		Column("version").Type(sql_utils.ColumnTypeInt).NotNull().Default(0).EndColumn().
		Unique("network_id", "key", "type").
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create entities table")
		return
	}

	_, err = fact.builder.CreateTable(entityAssocTable).
		Column("from_pk").Type(sql_utils.ColumnTypeText).References(entityTable, "pk").OnDelete(sql_utils.ColumnOnDeleteCascade).EndColumn().
		Column("to_pk").Type(sql_utils.ColumnTypeText).References(entityTable, "pk").OnDelete(sql_utils.ColumnOnDeleteCascade).EndColumn().
		PrimaryKey("from_pk", "to_pk").
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create entity assoc table")
		return
	}

	_, err = fact.builder.CreateTable(entityAclTable).
		Column("id").Type(sql_utils.ColumnTypeText).PrimaryKey().EndColumn().
		Column("entity_pk").Type(sql_utils.ColumnTypeText).References(entityTable, "pk").OnDelete(sql_utils.ColumnOnDeleteCascade).EndColumn().
		Column("scope").Type(sql_utils.ColumnTypeText).NotNull().EndColumn().
		Column("permission").Type(sql_utils.ColumnTypeInt).NotNull().EndColumn().
		Column("type").Type(sql_utils.ColumnTypeText).NotNull().EndColumn().
		Column("id_filter").Type(sql_utils.ColumnTypeText).EndColumn().
		Column("version").Type(sql_utils.ColumnTypeInt).NotNull().Default(0).EndColumn().
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create entity acl table")
		return
	}

	// Create indexes (index is not implicitly created on a referencing FK)
	_, err = fact.builder.CreateIndex("graph_id_idx").
		IfNotExists().
		On(entityTable).
		Columns("graph_id").
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create graph ID index")
		return
	}

	_, err = fact.builder.CreateIndex("acl_ent_pk_idx").
		IfNotExists().
		On(entityAclTable).
		Columns("entity_pk").
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create acl ent PK index")
		return
	}

	// Create internal network(s)
	_, err = fact.builder.Insert(networksTable).
		Columns("id", "name", "description").
		Values(InternalNetworkID, internalNetworkName, internalNetworkDescription).
		OnConflict(nil, "id").
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "error creating internal networks")
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
	defer sql_utils.ClearStatementCacheLogOnError(stmtCache, "UpdateNetworks")

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
		_, err := store.builder.Delete(entityTable).
			Where(sq.And{
				sq.Eq{"network_id": networkID},
				sq.Eq{"type": update.Type},
				sq.Eq{"key": update.Key},
			}).
			RunWith(store.tx).
			Exec()
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
	loadResult, err := store.loadFromEntitiesTable(networkID, EntityLoadFilter{IDs: []storage.TypeAndKey{entityID}}, EntityLoadCriteria{})
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
