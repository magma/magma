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
	"os"
	"sort"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	sq "github.com/Masterminds/squirrel"
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

	aclIdCol       = "id"
	aclEntCol      = "entity_pk"
	aclScopeCol    = "scope"
	aclPermCol     = "permission"
	aclTypeCol     = "type"
	aclIdFilterCol = "id_filter"
	aclVerCol      = "version"
)

// NewSQLConfiguratorStorageFactory returns a ConfiguratorStorageFactory
// implementation backed by a SQL database.
func NewSQLConfiguratorStorageFactory(db *sql.DB, generator storage.IDGenerator, sqlBuilder sqorc.StatementBuilder) ConfiguratorStorageFactory {
	return &sqlConfiguratorStorageFactory{db: db, idGenerator: generator, builder: sqlBuilder}
}

type sqlConfiguratorStorageFactory struct {
	db          *sql.DB
	idGenerator storage.IDGenerator
	builder     sqorc.StatementBuilder
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
		err = errors.Wrap(err, "failed to create networks table")
		return
	}

	// Adding a type column if it doesn't exist already. This will ensure network
	// tables that are already created will also have the type column.
	// TODO Remove after 1-2 months to ensure service isn't disrupted
	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s text", networksTable, nwTypeCol))
	// special case sqlite3 because ADD COLUMN IF NOT EXISTS is not supported
	// and we only run sqlite3 for unit tests
	if err != nil && os.Getenv("SQL_DRIVER") != "sqlite3" {
		err = errors.Wrap(err, "failed to add 'type' field to networks table")
	}

	_, err = fact.builder.CreateIndex("type_idx").
		IfNotExists().
		On(networksTable).
		Columns(nwTypeCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create network type index")
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
		err = errors.Wrap(err, "failed to create network configs table")
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
		err = errors.Wrap(err, "failed to create entities table")
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
		err = errors.Wrap(err, "failed to create entity assoc table")
		return
	}

	_, err = fact.builder.CreateTable(entityAclTable).
		IfNotExists().
		Column(aclIdCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(aclEntCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(aclScopeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(aclPermCol).Type(sqorc.ColumnTypeInt).NotNull().EndColumn().
		Column(aclTypeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(aclIdFilterCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(aclVerCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		ForeignKey(entityTable, map[string]string{aclEntCol: entPkCol}, sqorc.ColumnOnDeleteCascade).
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
		Columns(entGidCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create graph ID index")
		return
	}

	_, err = fact.builder.CreateIndex("acl_ent_pk_idx").
		IfNotExists().
		On(entityAclTable).
		Columns(aclEntCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create acl ent PK index")
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
		err = errors.Wrap(err, "error creating internal networks")
		return
	}

	return
}

func (fact *sqlConfiguratorStorageFactory) StartTransaction(ctx context.Context, opts *storage.TxOptions) (ConfiguratorStorage, error) {
	tx, err := fact.db.BeginTx(ctx, getSqlOpts(opts))
	if err != nil {
		return nil, err
	}
	return &sqlConfiguratorStorage{tx: tx, idGenerator: fact.idGenerator, builder: fact.builder}, nil
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
	tx          *sql.Tx
	idGenerator storage.IDGenerator
	builder     sqorc.StatementBuilder
}

func (store *sqlConfiguratorStorage) Commit() error {
	return store.tx.Commit()
}

func (store *sqlConfiguratorStorage) Rollback() error {
	return store.tx.Rollback()
}

func (store *sqlConfiguratorStorage) LoadNetworks(filter NetworkLoadFilter, loadCriteria NetworkLoadCriteria) (NetworkLoadResult, error) {
	emptyRet := NetworkLoadResult{NetworkIDsNotFound: []string{}, Networks: []*Network{}}
	if funk.IsEmpty(filter.Ids) && funk.IsEmpty(filter.TypeFilter) {
		return emptyRet, nil
	}

	selectBuilder := store.getLoadNetworksSelectBuilder(filter, loadCriteria)
	if loadCriteria.LoadConfigs {
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

	loadedNetworksByID, loadedNetworkIDs, err := scanNetworkRows(rows, loadCriteria)
	if err != nil {
		return emptyRet, err
	}

	ret := NetworkLoadResult{
		NetworkIDsNotFound: getNetworkIDsNotFound(loadedNetworksByID, filter.Ids),
		Networks:           make([]*Network, 0, len(loadedNetworksByID)),
	}
	for _, nid := range loadedNetworkIDs {
		ret.Networks = append(ret.Networks, loadedNetworksByID[nid])
	}
	return ret, nil
}

func (store *sqlConfiguratorStorage) LoadAllNetworks(loadCriteria NetworkLoadCriteria) ([]Network, error) {
	emptyNetworks := []Network{}
	idsToExclude := []string{InternalNetworkID}

	selectBuilder := store.builder.Select(getNetworkQueryColumns(loadCriteria)...).
		From(networksTable).
		Where(sq.NotEq{
			fmt.Sprintf("%s.%s", networksTable, nwIDCol): idsToExclude,
		})
	if loadCriteria.LoadConfigs {
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
		Columns(nwIDCol, nwTypeCol, nwNameCol, nwDescCol).
		Values(network.ID, network.Type, network.Name, network.Description).
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
		Columns(nwcIDCol, nwcTypeCol, nwcValCol)
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
	defer sqorc.ClearStatementCacheLogOnError(stmtCache, "UpdateNetworks")

	// Update networks first
	for _, update := range networksToUpdate {
		err := store.updateNetwork(update, stmtCache)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	_, err := store.builder.Delete(networkConfigTable).Where(sq.Eq{nwcIDCol: networksToDelete}).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete configs associated with networks")
	}
	_, err = store.builder.Delete(networksTable).Where(sq.Eq{nwIDCol: networksToDelete}).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete networks")
	}
	return nil
}

func (store *sqlConfiguratorStorage) LoadEntities(networkID string, filter EntityLoadFilter, loadCriteria EntityLoadCriteria) (EntityLoadResult, error) {
	ret := EntityLoadResult{Entities: []*NetworkEntity{}, EntitiesNotFound: []*EntityID{}}

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
		ret.Entities = append(ret.Entities, ent)
	}
	ret.EntitiesNotFound = calculateEntitiesNotFound(entsByPk, filter.IDs)

	// Sort entities for deterministic returns
	entComparator := func(a, b *NetworkEntity) bool {
		return a.GetTypeAndKey().String() < b.GetTypeAndKey().String()
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
	if !funk.IsEmpty(createdEntWithPk.Associations) {
		createdEntWithPk.Associations = funk.Chain(createdEntWithPk.Associations).
			Map(func(id *EntityID) storage.TypeAndKey { return id.ToTypeAndKey() }).
			Uniq().
			Map(func(tk storage.TypeAndKey) *EntityID { return (&EntityID{}).FromTypeAndKey(tk) }).
			Value().([]*EntityID)
	}

	createdEntWithPk.NetworkID = networkID
	return createdEntWithPk.NetworkEntity, nil
}

func (store *sqlConfiguratorStorage) UpdateEntity(networkID string, update EntityUpdateCriteria) (NetworkEntity, error) {
	emptyRet := NetworkEntity{Type: update.Type, Key: update.Key}
	entToUpdate, err := store.loadEntToUpdate(networkID, update)
	if err != nil && !update.DeleteEntity {
		return emptyRet, errors.Wrap(err, "failed to load entity being updated")
	}
	if entToUpdate == nil {
		return emptyRet, nil
	}

	if update.DeleteEntity {
		// Cascading FK relations in the schema will handle the other tables
		_, err := store.builder.Delete(entityTable).
			Where(sq.And{
				sq.Eq{entNidCol: networkID},
				sq.Eq{entTypeCol: update.Type},
				sq.Eq{entKeyCol: update.Key},
			}).
			RunWith(store.tx).
			Exec()
		if err != nil {
			return emptyRet, errors.Wrapf(err, "failed to delete entity (%s, %s)", update.Type, update.Key)
		}

		// Deleting a node could partition its graph
		err = store.fixGraph(networkID, entToUpdate.GraphID, entToUpdate)
		if err != nil {
			return emptyRet, errors.Wrap(err, "failed to fix entity graph after deletion")
		}

		return emptyRet, nil
	}

	// Then, update the fields on the entity table
	entToUpdate.NetworkID = networkID
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
	err = store.processEdgeUpdates(networkID, update, entToUpdate)
	if err != nil {
		return entToUpdate.NetworkEntity, errors.WithStack(err)
	}

	return entToUpdate.NetworkEntity, nil
}

func (store *sqlConfiguratorStorage) LoadGraphForEntity(networkID string, entityID EntityID, loadCriteria EntityLoadCriteria) (EntityGraph, error) {
	// Technically you could do this in one DB query with a subquery in the
	// WHERE when selecting from the entity table.
	// But until we hit some kind of scaling limit, let's keep the code simple
	// and delegate to LoadGraph after loading the requested entity.

	// We just care about getting the graph ID off this entity so use an empty
	// load criteria
	loadResult, err := store.loadFromEntitiesTable(networkID, EntityLoadFilter{IDs: []*EntityID{&entityID}}, EntityLoadCriteria{})
	if err != nil {
		return EntityGraph{}, errors.Wrap(err, "failed to load entity for graph query")
	}

	var loadedEnt *NetworkEntity
	for _, ent := range loadResult {
		loadedEnt = ent
	}
	if loadedEnt == nil {
		return EntityGraph{}, errors.Errorf("could not find requested entity (%s) for graph query", entityID.String())
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
	retEnts := funk.Map(internalGraph.entsByPk, func(_ string, ent *NetworkEntity) *NetworkEntity { return ent }).([]*NetworkEntity)
	retRoots := funk.Map(rootPks, func(pk string) *EntityID { return &EntityID{Type: entTksByPk[pk].Type, Key: entTksByPk[pk].Key} }).([]*EntityID)
	sort.Slice(retEnts, func(i, j int) bool {
		return storage.IsTKLessThan(retEnts[i].GetTypeAndKey(), retEnts[j].GetTypeAndKey())
	})
	sort.Slice(retRoots, func(i, j int) bool {
		return storage.IsTKLessThan(retRoots[i].ToTypeAndKey(), retRoots[j].ToTypeAndKey())
	})

	return EntityGraph{
		Entities:     retEnts,
		RootEntities: retRoots,
		Edges:        edges,
	}, nil
}
