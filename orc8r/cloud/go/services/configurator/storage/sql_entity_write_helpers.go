/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import (
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type entWithPk struct {
	pk string
	NetworkEntity
}

func (store *sqlConfiguratorStorage) doesEntExist(networkID string, tk storage.TypeAndKey) (bool, error) {
	var count uint64
	err := store.builder.Select("COUNT(1)").
		From(entityTable).
		Where(sq.And{
			sq.Eq{"network_id": networkID},
			sq.Eq{"type": tk.Type},
			sq.Eq{"key": tk.Key},
		}).
		RunWith(store.tx).
		QueryRow().Scan(&count)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check for existence of entity %s: %s", tk, err)
	}

	return count > 0, nil
}

func (store *sqlConfiguratorStorage) insertIntoEntityTable(networkID string, entity NetworkEntity) (entWithPk, error) {
	pk := store.idGenerator.New()
	// On create, we'll generate a new graph ID for the entity temporarily
	graphID := store.idGenerator.New()
	entity.GraphID = graphID

	_, err := store.builder.Insert(entityTable).
		Columns("pk", "network_id", "type", "key", "graph_id", "name", "description", "physical_id", "config").
		Values(pk, networkID, entity.Type, entity.Key, entity.GraphID, toNullable(entity.Name), toNullable(entity.Description), toNullable(entity.PhysicalID), toNullable(entity.Config)).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return entWithPk{}, fmt.Errorf("failed to create entity %s: %s", entity.GetTypeAndKey(), err)
	}
	return entWithPk{pk: pk, NetworkEntity: entity}, nil
}

// acls is an output parameter - entries will be updated in-place with
// system-generated IDs.
func (store *sqlConfiguratorStorage) createPermissions(networkID string, pk string, acls []ACL) error {
	if funk.IsEmpty(acls) {
		return nil
	}

	insertBuilder := store.builder.Insert(entityAclTable).
		Columns("id", "entity_pk", "scope", "permission", "type", "id_filter")

	aclIDs := make([]string, 0, len(acls))
	for _, acl := range acls {
		aclID := store.idGenerator.New()
		scopeVal, err := serializeACLScope(acl.Scope)
		if err != nil {
			return err
		}
		typeVal, err := serializeACLType(acl.Type)
		if err != nil {
			return err
		}

		insertBuilder = insertBuilder.Values(aclID, pk, scopeVal, acl.Permission, typeVal, serializeACLIDFilter(acl.IDFilter))
		aclIDs = append(aclIDs, aclID)
	}

	_, err := insertBuilder.RunWith(store.tx).Exec()
	if err != nil {
		return errors.Wrap(err, "failed to create permissions")
	}

	for i, aclID := range aclIDs {
		acls[i].ID = aclID
	}
	return nil
}

func (store *sqlConfiguratorStorage) createEdges(networkID string, entity entWithPk) (map[storage.TypeAndKey]entWithPk, error) {
	// Load the associated entities first because we need to know PKs
	// This will also load graph ID on the entity because creating an edge can
	// involve merging previously disjoint graphs.
	entsByTk, err := store.loadEntsFromEdges(networkID, entity)
	if err != nil {
		return entsByTk, err
	}
	if funk.IsEmpty(entity.GetGraphEdges()) {
		return entsByTk, err
	}

	insertBuilder := store.builder.Insert(entityAssocTable).
		Columns("from_pk", "to_pk").
		OnConflict(nil, "from_pk", "to_pk")
	for _, edge := range entity.GetGraphEdges() {
		fromPk := entsByTk[edge.From].pk
		toPk := entsByTk[edge.To].pk
		insertBuilder = insertBuilder.Values(fromPk, toPk)
	}
	_, err = insertBuilder.RunWith(store.tx).Exec()
	if err != nil {
		return entsByTk, errors.Wrap(err, "error creating assocs")
	}
	return entsByTk, nil
}

func (store *sqlConfiguratorStorage) loadEntsFromEdges(networkID string, targetEntity entWithPk) (map[storage.TypeAndKey]entWithPk, error) {
	ret := map[storage.TypeAndKey]entWithPk{targetEntity.GetTypeAndKey(): targetEntity}

	loadedEntsByTk, err := store.loadEntsWithPksByTK(networkID, targetEntity.Associations)
	if err != nil {
		return ret, errors.WithStack(err)
	}
	loadedEntsByTk[targetEntity.GetTypeAndKey()] = targetEntity
	return loadedEntsByTk, nil
}

func (store *sqlConfiguratorStorage) loadEntsWithPksByTK(networkID string, tksToLoad []storage.TypeAndKey) (map[storage.TypeAndKey]entWithPk, error) {
	ret := make(map[storage.TypeAndKey]entWithPk, len(tksToLoad)+1)
	if funk.IsEmpty(tksToLoad) {
		return ret, nil
	}

	uniqTksToLoad := funk.Uniq(tksToLoad).([]storage.TypeAndKey)
	loadedEntsByPk, err := store.loadFromEntitiesTable(networkID, EntityLoadFilter{IDs: uniqTksToLoad}, EntityLoadCriteria{})
	if err != nil {
		return ret, errors.WithStack(err)
	}
	for pk, ent := range loadedEntsByPk {
		ret[ent.GetTypeAndKey()] = entWithPk{pk: pk, NetworkEntity: *ent}
	}

	entsNotFound := calculateEntitiesNotFound(loadedEntsByPk, tksToLoad)
	if !funk.IsEmpty(entsNotFound) {
		return ret, errors.Errorf("could not find entities matching %v", entsNotFound)
	}

	return ret, nil
}

func (store *sqlConfiguratorStorage) mergeGraphs(createdEntity entWithPk, allAssociatedEntsByTk map[storage.TypeAndKey]entWithPk) (string, error) {
	// If we create a node which bridges 2 previously disjoint graphs, then
	// we need to change the ID of one of the graphs to the joined one.

	// If we associate to no graphs, then no-op - we'll use the
	// system-generated graph ID for this single-node graph.

	// If we associate to only 1 graph, then we'll overwrite this node's
	// graph ID with the ID of that graph.

	// If we associate to 2+ graphs, this means that we need to merge all of
	// them into a single graph. Pick the lexicographically smallest graph ID
	// to use as the ID for the final graph
	adjacentGraphs := funk.Chain(createdEntity.Associations).
		Map(func(tk storage.TypeAndKey) string { return allAssociatedEntsByTk[tk].GraphID }).
		Uniq().
		Value().([]string)
	noMergeNecessary := funk.IsEmpty(adjacentGraphs) || (len(adjacentGraphs) == 1 && adjacentGraphs[0] == createdEntity.GraphID)
	if noMergeNecessary {
		return createdEntity.GraphID, nil
	}

	sort.Strings(adjacentGraphs)
	targetGraphID := adjacentGraphs[0]
	graphIDsToChange := []string{createdEntity.GraphID}
	for _, oldGraphID := range adjacentGraphs[1:] {
		graphIDsToChange = append(graphIDsToChange, oldGraphID)
	}

	// let squirrel cache prepared statements for us (there should only be 1)
	sc := sq.NewStmtCache(store.tx)
	defer sqorc.ClearStatementCacheLogOnError(sc, "mergeGraphs")

	for _, oldGraphID := range graphIDsToChange {
		_, err := store.builder.Update(entityTable).
			Set("graph_id", targetGraphID).
			Where(sq.Eq{"graph_id": oldGraphID}).
			RunWith(sc).
			Exec()
		if err != nil {
			return "", errors.Wrap(err, "error updating entity graphs")
		}
	}

	return targetGraphID, nil
}

func (store *sqlConfiguratorStorage) loadEntToUpdate(networkID string, update EntityUpdateCriteria) (entWithPk, error) {
	loadedEntByPk, err := store.loadFromEntitiesTable(
		networkID,
		EntityLoadFilter{IDs: []storage.TypeAndKey{update.GetTypeAndKey()}},
		EntityLoadCriteria{},
	)
	if err != nil {
		return entWithPk{}, errors.Wrap(err, "failed to load entity to update")
	}
	if len(loadedEntByPk) != 1 {
		return entWithPk{}, errors.Errorf("expected to load 1 ent for update, got %d", len(loadedEntByPk))
	}

	return funk.Chain(loadedEntByPk).
		Map(func(pk string, ent *NetworkEntity) entWithPk { return entWithPk{pk: pk, NetworkEntity: *ent} }).
		Head().(entWithPk), nil
}

// entOut is an output parameter
func (store *sqlConfiguratorStorage) processEntityFieldsUpdate(pk string, update EntityUpdateCriteria, entOut *NetworkEntity) error {
	_, err := store.getEntityUpdateQueryBuilder(pk, update).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to update entity fields")
	}

	if update.NewName != nil {
		entOut.Name = *update.NewName
	}
	if update.NewDescription != nil {
		entOut.Description = *update.NewDescription
	}
	if update.NewPhysicalID != nil {
		entOut.PhysicalID = *update.NewPhysicalID
	}
	if update.NewConfig != nil {
		entOut.Config = *update.NewConfig
	}
	entOut.Version++

	return nil
}

// entOut is an output parameter
func (store *sqlConfiguratorStorage) processPermissionUpdates(networkID string, entPk string, update EntityUpdateCriteria, entOut *NetworkEntity) error {
	if !funk.IsEmpty(update.PermissionsToCreate) {
		err := store.createPermissions(networkID, entPk, update.PermissionsToCreate)
		if err != nil {
			return errors.WithStack(err)
		}
		entOut.Permissions = append(entOut.Permissions, update.PermissionsToCreate...)
	}

	if !funk.IsEmpty(update.PermissionsToUpdate) {
		err := store.updatePermissions(entPk, update.PermissionsToUpdate, entOut)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if !funk.IsEmpty(update.PermissionsToDelete) {
		err := store.deletePermissions(update.PermissionsToDelete, entOut)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// entToUpdateOut is an output parameter
func (store *sqlConfiguratorStorage) processEdgeUpdates(networkID string, update EntityUpdateCriteria, entToUpdateOut *entWithPk) error {
	if funk.IsEmpty(update.AssociationsToAdd) && funk.IsEmpty(update.AssociationsToDelete) {
		return nil
	}

	// First, create edges. Because createEdges expects an entWithPk,
	// we'll just make the ent's Associations the edges we want to create
	entToUpdateOut.Associations = update.AssociationsToAdd
	newlyAssociatedEntsByTk, err := store.createEdges(networkID, *entToUpdateOut)
	if err != nil {
		entToUpdateOut.Associations = nil
		return errors.WithStack(err)
	}

	// Just like entity creation, we might need to merge graphs after adding
	newGraphID, err := store.mergeGraphs(*entToUpdateOut, newlyAssociatedEntsByTk)
	if err != nil {
		return errors.WithStack(err)
	}
	entToUpdateOut.GraphID = newGraphID

	// Now delete edges.
	err = store.deleteEdges(networkID, update.AssociationsToDelete, entToUpdateOut)
	if err != nil {
		return errors.WithStack(err)
	}

	// Finally, we need to load the entire graph corresponding to this node's
	// graph ID (which may or may not have changed during the update).
	// If we only created edges, we know that this graph is correct.
	// However, if we deleted edges, we could have partitioned this graph, so
	// we need to do a connected component search. If we come up with multiple
	// components, then each new component needs to be updated with a new
	// graph ID.
	if funk.IsEmpty(update.AssociationsToDelete) {
		return nil
	}

	err = store.fixGraph(networkID, newGraphID, entToUpdateOut)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (store *sqlConfiguratorStorage) getEntityUpdateQueryBuilder(pk string, update EntityUpdateCriteria) sq.UpdateBuilder {
	// UPDATE cfg_entities SET (name, description, physical_id, config, version) = ($1, $2, $3, $4, cfg_entities.version + 1)
	// WHERE pk = $5
	updateBuilder := store.builder.Update(entityTable).Where(sq.Eq{"pk": pk})
	if update.NewName != nil {
		updateBuilder = updateBuilder.Set("name", *update.NewName)
	}
	if update.NewDescription != nil {
		updateBuilder = updateBuilder.Set("description", *update.NewDescription)
	}
	if update.NewPhysicalID != nil {
		updateBuilder = updateBuilder.Set("physical_id", *update.NewPhysicalID)
	}
	if update.NewConfig != nil {
		updateBuilder = updateBuilder.Set("config", *update.NewConfig)
	}
	updateBuilder = updateBuilder.Set("version", sq.Expr("version+1"))
	return updateBuilder
}

// entOut is an output parameter
func (store *sqlConfiguratorStorage) updatePermissions(entPk string, permissions []ACL, entOut *NetworkEntity) error {
	aclsExist, err := store.doAllACLsExist(permissions)
	if err != nil {
		return err
	}
	if !aclsExist {
		return errors.New("not all ACLs being updated exist")
	}

	// We'll let squirrel cache prepared statements for us (there should only
	// be 1 in the current impl because update is all-or-nothing)
	sc := sq.NewStmtCache(store.tx)
	defer sqorc.ClearStatementCacheLogOnError(sc, "updatePermissions")

	for _, acl := range permissions {
		scopeVal, err := serializeACLScope(acl.Scope)
		if err != nil {
			return err
		}
		typeVal, err := serializeACLType(acl.Type)
		if err != nil {
			return err
		}

		// UPDATE cfg_acls SET (scope, permission, type, id_filter, version) = ($1, $2, $3, $4, cfg_acls.version+1)
		// WHERE cfg_acls.id = $5
		_, err = store.builder.Update(entityAclTable).
			Set("scope", scopeVal).
			Set("permission", acl.Permission).
			Set("type", typeVal).
			Set("id_filter", serializeACLIDFilter(acl.IDFilter)).
			Set("version", sq.Expr("version+1")).
			Where(sq.Eq{"id": acl.ID}).
			RunWith(sc).
			Exec()
		if err != nil {
			return errors.Wrapf(err, "failed to update permission %s", acl.ID)
		}
	}

	entOut.Permissions = append(entOut.Permissions, permissions...)
	return nil
}

func (store *sqlConfiguratorStorage) doAllACLsExist(acls []ACL) (bool, error) {
	aclIDs := funk.Map(acls, func(acl ACL) interface{} { return acl.ID }).([]interface{})
	var count uint64

	err := store.builder.Select("COUNT(*)").
		From(entityAclTable).
		Where(sq.Eq{"id": aclIDs}).
		RunWith(store.tx).
		QueryRow().Scan(&count)
	if err == sql.ErrNoRows {
		return false, errors.New("no ACLs found matching ACLs to update")
	}
	if err != nil {
		return false, errors.Wrap(err, "failed to query for ACLs matching ACLs to update")
	}
	return count == uint64(len(acls)), nil
}

// entOut is an output parameter
func (store *sqlConfiguratorStorage) deletePermissions(aclIDs []string, entOut *NetworkEntity) error {
	ids := make([]interface{}, 0, len(aclIDs))
	funk.ConvertSlice(funk.UniqString(aclIDs), &ids)

	_, err := store.builder.Delete(entityAclTable).
		Where(sq.Eq{"id": aclIDs}).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete ACLs")
	}

	if funk.IsEmpty(entOut.Permissions) {
		return nil
	}

	idsSet := funk.Map(aclIDs, func(i string) (string, bool) { return i, true }).(map[string]bool)
	entOut.Permissions = funk.Filter(entOut.Permissions, func(acl ACL) bool {
		_, deleted := idsSet[acl.ID]
		return !deleted
	}).([]ACL)

	return nil
}

// entToUpdateOut is an output parameter
func (store *sqlConfiguratorStorage) deleteEdges(networkID string, edgesToDelete []storage.TypeAndKey, entToUpdateOut *entWithPk) error {
	if funk.IsEmpty(edgesToDelete) {
		return nil
	}

	loadedEntsByTk, err := store.loadEntsWithPksByTK(networkID, edgesToDelete)
	if err != nil {
		return errors.Wrap(err, "could not load entities matching associations to delete")
	}

	orClause := make(sq.Or, 0, len(edgesToDelete))
	for _, edge := range edgesToDelete {
		orClause = append(orClause, sq.And{
			sq.Eq{"from_pk": entToUpdateOut.pk},
			sq.Eq{"to_pk": loadedEntsByTk[edge].pk},
		})
	}

	_, err = store.builder.Delete(entityAssocTable).
		Where(orClause).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to delete assocs")
	}

	if funk.IsEmpty(entToUpdateOut.Associations) {
		return nil
	}

	// Remove deleted edges from the passed in ent
	edgesToDeleteSet := funk.Map(edgesToDelete, func(tk storage.TypeAndKey) (storage.TypeAndKey, struct{}) { return tk, struct{}{} }).(map[storage.TypeAndKey]struct{})
	entToUpdateOut.Associations = funk.Filter(entToUpdateOut.Associations, func(tk storage.TypeAndKey) bool {
		_, wasDeleted := edgesToDeleteSet[tk]
		return !wasDeleted
	}).([]storage.TypeAndKey)

	return nil
}

func serializeACLScope(scope ACLScope) (string, error) {
	switch scope.Wildcard {
	case NoWildcard:
		return strings.Join(scope.NetworkIDs, ","), nil
	case WildcardAll:
		return wildcardAllString, nil
	default:
		return "", fmt.Errorf("unrecognized ACL scope wildcard %v", scope.Wildcard)
	}
}

func serializeACLType(t ACLType) (string, error) {
	switch t.Wildcard {
	case NoWildcard:
		return t.EntityType, nil
	case WildcardAll:
		return wildcardAllString, nil
	default:
		return "", fmt.Errorf("unrecognized ACL type wildcard %v", t.Wildcard)
	}
}

func serializeACLIDFilter(filter []string) sql.NullString {
	if funk.IsEmpty(filter) {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: strings.Join(filter, ",")}
}

func toNullable(field interface{}) interface{} {
	t := reflect.TypeOf(field)
	switch t.Kind() {
	case reflect.String:
		if field.(string) == "" {
			return nil
		} else {
			return field
		}
	case reflect.Array, reflect.Slice:
		if funk.IsEmpty(field) {
			return nil
		} else {
			return field
		}
	default:
		return field
	}
}
