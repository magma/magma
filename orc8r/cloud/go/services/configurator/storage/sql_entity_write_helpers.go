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

	"magma/orc8r/cloud/go/sql_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/thoas/go-funk"
)

type entWithPk struct {
	pk string
	NetworkEntity
}

func (store *sqlConfiguratorStorage) insertIntoEntityTable(networkID string, entity NetworkEntity) (entWithPk, error) {
	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (pk, network_id, type, key, graph_id, name, description, physical_id, config)
		VALUES %s
	`, entityTable, sql_utils.GetPlaceholderArgList(1, 9))
	aclInsertQuery := fmt.Sprintf("INSERT INTO %s (id, entity_pk, scope, permission, type, id_filter) VALUES %s", entityAclTable, sql_utils.GetPlaceholderArgList(1, 6))
	stmts, err := sql_utils.PrepareStatements(store.tx, []string{insertQuery, aclInsertQuery})
	if err != nil {
		return entWithPk{}, err
	}
	defer sql_utils.GetCloseStatementsDeferFunc(stmts, "CreateEntities")()
	insertStmt, aclStmt := stmts[0], stmts[1]

	pk := store.idGenerator.New()
	// On create, we'll generate a new graph ID for the entity temporarily
	graphID := store.idGenerator.New()
	entity.GraphID = graphID

	_, err = insertStmt.Exec(pk, networkID, entity.Type, entity.Key, entity.GraphID, toNullable(entity.Name), toNullable(entity.Description), toNullable(entity.PhysicalID), toNullable(entity.Config))
	if err != nil {
		return entWithPk{}, fmt.Errorf("failed to create entity %s: %s", entity.GetTypeAndKey(), err)
	}

	// Create ACLs
	for i, acl := range entity.Permissions {
		aclID := store.idGenerator.New()
		scopeVal, err := serializeACLScope(acl.Scope)
		if err != nil {
			return entWithPk{}, fmt.Errorf("failed to create entity %s: %s", entity.GetTypeAndKey(), err)
		}
		typeVal, err := serializeACLType(acl.Type)
		if err != nil {
			return entWithPk{}, fmt.Errorf("failed to create entity %s: %s", entity.GetTypeAndKey(), err)
		}

		_, err = aclStmt.Exec(aclID, pk, scopeVal, acl.Permission, typeVal, serializeACLIDFilter(acl.IDFilter))
		if err != nil {
			return entWithPk{}, fmt.Errorf("failed to create permissions for entity %s: %s", entity.GetTypeAndKey(), err)
		}
		// `acl` in this context is a new variable allocation (copy-on-write),
		// so use the array index to modify the permission
		entity.Permissions[i].ID = aclID
	}

	return entWithPk{pk: pk, NetworkEntity: entity}, nil
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

	tksToLoad := []storage.TypeAndKey{}
	for _, edge := range targetEntity.Associations {
		tksToLoad = append(tksToLoad, edge)
	}
	tksToLoad = funk.Uniq(tksToLoad).([]storage.TypeAndKey)
	if funk.IsEmpty(tksToLoad) {
		return ret, nil
	}

	loadedEntsByPk, err := store.loadSpecificEntities(networkID, EntityLoadFilter{IDs: tksToLoad}, EntityLoadCriteria{})
	if err != nil {
		return ret, err
	}
	assocsNotFound := calculateEntitiesNotFound(loadedEntsByPk, tksToLoad)
	if !funk.IsEmpty(assocsNotFound) {
		return ret, fmt.Errorf("could not find entities for assocs to: %v", assocsNotFound)
	}

	for pk, ent := range loadedEntsByPk {
		ret[ent.GetTypeAndKey()] = entWithPk{pk: pk, NetworkEntity: *ent}
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
	adjacentGraphs := funk.Map(createdEntity.Associations, func(tk storage.TypeAndKey) string { return allAssociatedEntsByTk[tk].GraphID }).([]string)
	adjacentGraphs = funk.UniqString(adjacentGraphs)
	if funk.IsEmpty(adjacentGraphs) {
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
