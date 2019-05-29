/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"magma/orc8r/cloud/go/sql_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type entWithPk struct {
	pk string
	NetworkEntity
}

func (store *sqlConfiguratorStorage) doesEntExist(networkID string, tk storage.TypeAndKey) (bool, error) {
	query := fmt.Sprintf("SELECT count(1) FROM %s WHERE (network_id, type, key) = ($1, $2, $3)", entityTable)
	row := store.tx.QueryRow(query, networkID, tk.Type, tk.Key)
	var count uint64
	err := row.Scan(&count)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check for existence of entity %s: %s", tk, err)
	}

	return count > 0, nil
}

func (store *sqlConfiguratorStorage) insertIntoEntityTable(networkID string, entity NetworkEntity) (entWithPk, error) {
	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (pk, network_id, type, key, graph_id, name, description, physical_id, config)
		VALUES %s
	`, entityTable, sql_utils.GetPlaceholderArgList(1, 9))
	stmts, err := sql_utils.PrepareStatements(store.tx, []string{insertQuery})
	if err != nil {
		return entWithPk{}, err
	}
	defer sql_utils.GetCloseStatementsDeferFunc(stmts, "CreateEntities")()
	insertStmt := stmts[0]

	pk := store.idGenerator.New()
	// On create, we'll generate a new graph ID for the entity temporarily
	graphID := store.idGenerator.New()
	entity.GraphID = graphID

	_, err = insertStmt.Exec(pk, networkID, entity.Type, entity.Key, entity.GraphID, toNullable(entity.Name), toNullable(entity.Description), toNullable(entity.PhysicalID), toNullable(entity.Config))
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

	aclInsertQuery := fmt.Sprintf("INSERT INTO %s (id, entity_pk, scope, permission, type, id_filter) VALUES %s", entityAclTable, sql_utils.GetPlaceholderArgList(1, 6))
	stmts, err := sql_utils.PrepareStatements(store.tx, []string{aclInsertQuery})
	if err != nil {
		return errors.Wrap(err, "failed to prepare permission insert")
	}
	defer sql_utils.GetCloseStatementsDeferFunc(stmts, "createPermissions")()

	aclStmt := stmts[0]
	for i, acl := range acls {
		aclID := store.idGenerator.New()
		scopeVal, err := serializeACLScope(acl.Scope)
		if err != nil {
			return err
		}
		typeVal, err := serializeACLType(acl.Type)
		if err != nil {
			return err
		}

		_, err = aclStmt.Exec(aclID, pk, scopeVal, acl.Permission, typeVal, serializeACLIDFilter(acl.IDFilter))
		if err != nil {
			return errors.Wrapf(err, "failed to create permission %s", aclID)
		}
		// `acl` in this context is a new variable allocation (copy-on-write),
		// so use the array index to modify the permission
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

	assocInsertQuery := fmt.Sprintf("INSERT INTO %s (from_pk, to_pk) VALUES ($1, $2) ON CONFLICT DO NOTHING", entityAssocTable)
	stmts, err := sql_utils.PrepareStatements(store.tx, []string{assocInsertQuery})
	if err != nil {
		return entsByTk, err
	}
	defer sql_utils.GetCloseStatementsDeferFunc(stmts, "CreateEntities")()
	assocStmt := stmts[0]

	for _, edge := range entity.GetGraphEdges() {
		fromPk := entsByTk[edge.From].pk
		toPk := entsByTk[edge.To].pk

		_, err := assocStmt.Exec(fromPk, toPk)
		if err != nil {
			return entsByTk, fmt.Errorf("error creating assoc (%s, %s): %s", edge.From, edge.To, err)
		}
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

	graphUpdateQuery := fmt.Sprintf("UPDATE %s SET graph_id = $1 WHERE graph_id = $2", entityTable)
	stmts, err := sql_utils.PrepareStatements(store.tx, []string{graphUpdateQuery})
	if err != nil {
		return "", err
	}
	defer sql_utils.GetCloseStatementsDeferFunc(stmts, "CreateEntity")()

	for _, oldGraphID := range graphIDsToChange {
		_, err := stmts[0].Exec(targetGraphID, oldGraphID)
		if err != nil {
			return "", fmt.Errorf("error updating entity graphs: %s", err)
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
	exec, args, err := getUpdateEntityExec(pk, update)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = store.tx.Exec(exec, args...)
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

type updateEntityExecTemplateArgs struct {
	TableName, Fields, FieldsPlaceholder, ConditionPlaceholder string
}

func getUpdateEntityExec(pk string, update EntityUpdateCriteria) (string, []interface{}, error) {
	// UPDATE cfg_entities SET (name, description, physical_id, config, version) = ($1, $2, $3, $4, cfg_entities.version + 1)
	// WHERE pk = $5
	tmpl := template.Must(template.New("update_ent_exec").Parse(`
		UPDATE {{.TableName}} SET {{.Fields}} = {{.FieldsPlaceholder}}
		WHERE pk = {{.ConditionPlaceholder}}
	`))
	tmplArgs, sqlArgs := getUpdateEntityExecTemplateArgsAndSQLArgs(pk, update)

	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, tmplArgs)
	if err != nil {
		return "", []interface{}{}, errors.Wrap(err, "failed to format entity update query")
	}
	return buf.String(), sqlArgs, nil
}

func getUpdateEntityExecTemplateArgsAndSQLArgs(pk string, update EntityUpdateCriteria) (updateEntityExecTemplateArgs, []interface{}) {
	tmplArgs := updateEntityExecTemplateArgs{
		TableName: entityTable,
	}
	args := []interface{}{}

	fields := []string{}
	if update.NewName != nil {
		fields = append(fields, "name")
		args = append(args, *update.NewName)
	}
	if update.NewDescription != nil {
		fields = append(fields, "description")
		args = append(args, *update.NewDescription)
	}
	if update.NewPhysicalID != nil {
		fields = append(fields, "physical_id")
		args = append(args, *update.NewPhysicalID)
	}
	if update.NewConfig != nil {
		fields = append(fields, "config")
		args = append(args, *update.NewConfig)
	}
	fields = append(fields, "version")

	tmplArgs.Fields = fmt.Sprintf("(%s)", strings.Join(fields, ", "))
	tmplArgs.FieldsPlaceholder = sql_utils.GetPlaceholderArgListWithSuffix(
		1,
		// -1 here because version is set in-place
		len(fields)-1,
		"version + 1",
	)
	tmplArgs.ConditionPlaceholder = fmt.Sprintf("$%d", len(fields))

	args = append(args, pk)
	return tmplArgs, args
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

	updateExec := fmt.Sprintf(`
		UPDATE %s SET (scope, permission, type, id_filter, version) = ($1, $2, $3, $4, %s.version + 1)
		WHERE %s.id = $5
	`, entityAclTable, entityAclTable, entityAclTable)
	stmts, err := sql_utils.PrepareStatements(store.tx, []string{updateExec})
	if err != nil {
		return errors.Wrap(err, "failed to prepare ACl update statement")
	}
	defer sql_utils.GetCloseStatementsDeferFunc(stmts, "updatePermissions")()

	updateStmt := stmts[0]
	for _, acl := range permissions {
		scopeVal, err := serializeACLScope(acl.Scope)
		if err != nil {
			return err
		}
		typeVal, err := serializeACLType(acl.Type)
		if err != nil {
			return err
		}

		_, err = updateStmt.Exec(scopeVal, acl.Permission, typeVal, serializeACLIDFilter(acl.IDFilter), acl.ID)
		if err != nil {
			return errors.Wrapf(err, "failed to update permission %s", acl.ID)
		}
	}

	entOut.Permissions = append(entOut.Permissions, permissions...)
	return nil
}

func (store *sqlConfiguratorStorage) doAllACLsExist(acls []ACL) (bool, error) {
	countPermsQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE id IN %s",
		entityAclTable,
		sql_utils.GetPlaceholderArgList(1, len(acls)),
	)
	aclIDs := funk.Map(acls, func(acl ACL) interface{} { return acl.ID }).([]interface{})
	row := store.tx.QueryRow(countPermsQuery, aclIDs...)

	var count uint64
	err := row.Scan(&count)
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

	delExec := fmt.Sprintf("DELETE FROM %s WHERE id IN %s", entityAclTable, sql_utils.GetPlaceholderArgList(1, len(ids)))
	_, err := store.tx.Exec(delExec, ids...)
	if err != nil {
		return errors.Wrap(err, "failed to delete ACLs")
	}

	if funk.IsEmpty(entOut.Permissions) {
		return nil
	}

	idsSet := funk.Map(ids, func(i interface{}) (string, struct{}) { return i.(string), struct{}{} }).(map[string]struct{})
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

	deleteEdgeExec := fmt.Sprintf("DELETE FROM %s WHERE (from_pk, to_pk) = ($1, $2)", entityAssocTable)
	stmts, err := sql_utils.PrepareStatements(store.tx, []string{deleteEdgeExec})
	if err != nil {
		return errors.Wrap(err, "failed to prepare assoc delete statement")
	}
	defer sql_utils.GetCloseStatementsDeferFunc(stmts, "deleteEdges")()

	deleteStmt := stmts[0]
	for _, edge := range edgesToDelete {
		_, err = deleteStmt.Exec(entToUpdateOut.pk, loadedEntsByTk[edge].pk)
		if err != nil {
			return errors.Wrapf(err, "failed to delete assoc to %s", edge)
		}
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
