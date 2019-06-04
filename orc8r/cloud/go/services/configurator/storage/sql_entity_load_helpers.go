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
	"sort"
	"strings"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func (store *sqlConfiguratorStorage) loadFromEntitiesTable(networkID string, filter EntityLoadFilter, criteria EntityLoadCriteria) (map[string]*NetworkEntity, error) {
	// Pointer values because we're modifying entities in-place with ACLs (LEFT JOIN)
	entsByPk := map[string]*NetworkEntity{}

	selectBuilder := store.getLoadEntitiesSelectBuilder(networkID, filter, criteria)
	rows, err := selectBuilder.RunWith(store.tx).Query()
	if err != nil {
		return entsByPk, errors.Wrap(err, "error querying for entities")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			glog.Errorf("error closing *Rows in LoadEntities: %s", err)
		}
	}()

	for rows.Next() {
		err = scanNextEntityRow(rows, criteria, entsByPk)
		if err != nil {
			return entsByPk, err
		}
	}
	return entsByPk, nil
}

func (store *sqlConfiguratorStorage) getLoadEntitiesSelectBuilder(networkID string, filter EntityLoadFilter, criteria EntityLoadCriteria) sq.SelectBuilder {
	// SELECT ent.pk, ent.key, ent.type, ent.physical_id, ent.version, graph.graph_id, ent.name, ent.description, ent.config,
	// [[ acl.id, acl.scope, acl.permission, acl.type, acl.id_filter, acl.version ]]
	// FROM cfg_entities AS ent
	// [[ LEFT JOIN cfg_acls AS acl ON acl.entity_pk = ent.pk ]]
	// [[ WHERE (ent.network_id = $1 AND ent.key = $2 AND ent.type = $3) OR (ent.network_id ...) ... ]]
	selectBuilder := store.builder.Select(getLoadEntitiesColumns(criteria)...).
		From(fmt.Sprintf("%s AS ent", entityTable))
	if criteria.LoadPermissions {
		selectBuilder = selectBuilder.LeftJoin(fmt.Sprintf("%s AS acl ON acl.%s = ent.%s", entityAclTable, aclEntCol, entPkCol))
	}

	// The WHERE has ORs if specific IDs are provided
	if !funk.IsEmpty(filter.IDs) {
		orClause := make(sq.Or, 0, len(filter.IDs))
		funk.ForEach(filter.IDs, func(tk storage.TypeAndKey) {
			orClause = append(orClause, sq.And{
				sq.Eq{fmt.Sprintf("ent.%s", entNidCol): networkID},
				sq.Eq{fmt.Sprintf("ent.%s", entKeyCol): tk.Key},
				sq.Eq{fmt.Sprintf("ent.%s", entTypeCol): tk.Type},
			})
		})
		selectBuilder = selectBuilder.Where(orClause)
	} else {
		if filter.graphID != nil {
			selectBuilder = selectBuilder.Where(sq.Eq{fmt.Sprintf("ent.%s", entGidCol): *filter.graphID})
		} else {
			andClause := sq.And{sq.Eq{fmt.Sprintf("ent.%s", entNidCol): networkID}}
			if filter.KeyFilter != nil {
				andClause = append(andClause, sq.Eq{fmt.Sprintf("ent.%s", entKeyCol): *filter.KeyFilter})
			}
			if filter.TypeFilter != nil {
				andClause = append(andClause, sq.Eq{fmt.Sprintf("ent.%s", entTypeCol): *filter.TypeFilter})
			}
			selectBuilder = selectBuilder.Where(andClause)
		}
	}

	return selectBuilder
}

func getLoadEntitiesColumns(criteria EntityLoadCriteria) []string {
	fields := []string{
		fmt.Sprintf("ent.%s", entPkCol),
		fmt.Sprintf("ent.%s", entKeyCol),
		fmt.Sprintf("ent.%s", entTypeCol),
		fmt.Sprintf("ent.%s", entPidCol),
		fmt.Sprintf("ent.%s", entVerCol),
		fmt.Sprintf("ent.%s", entGidCol),
	}
	if criteria.LoadMetadata {
		fields = append(fields, fmt.Sprintf("ent.%s", entNameCol), fmt.Sprintf("ent.%s", entDescCol))
	}
	if criteria.LoadConfig {
		fields = append(fields, fmt.Sprintf("ent.%s", entConfCol))
	}
	if criteria.LoadPermissions {
		fields = append(
			fields,
			fmt.Sprintf("acl.%s", aclIdCol),
			fmt.Sprintf("acl.%s", aclScopeCol),
			fmt.Sprintf("acl.%s", aclPermCol),
			fmt.Sprintf("acl.%s", aclTypeCol),
			fmt.Sprintf("acl.%s", aclIdFilterCol),
			fmt.Sprintf("acl.%s", aclVerCol),
		)
	}
	return fields
}

// existingEntsByPkOut is an output parameter
func scanNextEntityRow(rows *sql.Rows, criteria EntityLoadCriteria, existingEntsByPkOut map[string]*NetworkEntity) error {
	var pk, key, entType, graphID string
	var physicalID sql.NullString
	var name, description sql.NullString

	var config []byte
	var entVersion uint64

	// Nullstrings here in case the entity doesn't have perms
	var aclid, aclscope, acltype sql.NullString
	var aclIdFilter sql.NullString
	var aclPermission, aclVersion sql.NullInt64

	// This corresponds with the order of the columns queried in the SELECT
	scanArgs := []interface{}{&pk, &key, &entType, &physicalID, &entVersion, &graphID}
	if criteria.LoadMetadata {
		scanArgs = append(scanArgs, &name, &description)
	}
	if criteria.LoadConfig {
		scanArgs = append(scanArgs, &config)
	}
	if criteria.LoadPermissions {
		scanArgs = append(scanArgs, &aclid, &aclscope, &aclPermission, &acltype, &aclIdFilter, &aclVersion)
	}

	err := rows.Scan(scanArgs...)
	if err != nil {
		return fmt.Errorf("error while scanning entity row: %s", err)
	}

	ent := NetworkEntity{
		Key:  key,
		Type: entType,

		Name:        nullStringToValue(name),
		Description: nullStringToValue(description),

		PhysicalID: nullStringToValue(physicalID),

		Config: config,

		GraphID: graphID,

		Version: entVersion,
	}
	if criteria.LoadPermissions && aclid.Valid {
		ent.Permissions = []ACL{
			{
				ID:         aclid.String,
				Scope:      deserializeACLScope(aclscope.String),
				Permission: ACLPermission(aclPermission.Int64),
				Type:       deserializeACLType(acltype.String),
				IDFilter:   deserializeACLIDFilter(aclIdFilter),
				Version:    uint64(aclVersion.Int64),
			},
		}
	}

	existingEnt, entExists := existingEntsByPkOut[pk]
	if entExists {
		if existingEnt.Permissions == nil {
			existingEnt.Permissions = []ACL{}
		}
		existingEnt.Permissions = append(existingEnt.Permissions, ent.Permissions...)
	} else {
		existingEntsByPkOut[pk] = &ent
	}
	return nil
}

func deserializeACLScope(aclScope string) ACLScope {
	if aclScope == wildcardAllString {
		return ACLScope{Wildcard: WildcardAll}
	} else {
		return ACLScope{NetworkIDs: strings.Split(aclScope, ",")}
	}
}

func deserializeACLType(aclType string) ACLType {
	if aclType == wildcardAllString {
		return ACLType{Wildcard: WildcardAll}
	} else {
		return ACLType{EntityType: aclType}
	}
}

func deserializeACLIDFilter(aclIdFilter sql.NullString) []string {
	if aclIdFilter.Valid {
		return strings.Split(aclIdFilter.String, ",")
	} else {
		return nil
	}
}

type loadedAssoc struct {
	fromPk, toPk string
}

func (store *sqlConfiguratorStorage) loadFromAssocsTable(filter EntityLoadFilter, criteria EntityLoadCriteria, entsByPk map[string]*NetworkEntity) ([]loadedAssoc, []string, error) {
	ret := []loadedAssoc{}
	allPks := map[string]struct{}{}
	if !criteria.LoadAssocsFromThis && !criteria.LoadAssocsToThis {
		return ret, []string{}, nil
	}

	entPks := funk.Keys(entsByPk).([]string)
	sort.Strings(entPks)

	// SELECT assoc.from_pk, assoc.to_pk FROM cfg_assocs AS assoc
	// WHERE assoc.from_pk IN ($1, $2, $3, ...)
	// OR assoc.to_pk IN ($4, $5, $6, ...)
	orClause := sq.Or{}
	if criteria.LoadAssocsFromThis {
		orClause = append(orClause, sq.Eq{fmt.Sprintf("assoc.%s", aFrCol): entPks})
	}
	if criteria.LoadAssocsToThis {
		orClause = append(orClause, sq.Eq{fmt.Sprintf("assoc.%s", aToCol): entPks})
	}
	// if we loaded all entities, save some network traffic and just load the
	// entire assocs table
	if filter.IsLoadAllEntities() {
		orClause = sq.Or{sq.Eq{"1": 1}}
	}

	assocRows, err := store.builder.Select(fmt.Sprintf("assoc.%s", aFrCol), fmt.Sprintf("assoc.%s", aToCol)).
		From("cfg_assocs AS assoc").
		Where(orClause).
		RunWith(store.tx).
		Query()
	if err != nil {
		return ret, []string{}, errors.Wrap(err, "error querying for associations")
	}
	defer sqorc.CloseRowsLogOnError(assocRows, "LoadEntities")

	for assocRows.Next() {
		var fromPk, toPk string
		err := assocRows.Scan(&fromPk, &toPk)
		if err != nil {
			return ret, []string{}, errors.Wrap(err, "error scanning association row")
		}

		ret = append(ret, loadedAssoc{fromPk: fromPk, toPk: toPk})
		allPks[fromPk] = struct{}{}
		allPks[toPk] = struct{}{}
	}

	allPksList := funk.Keys(allPks).([]string)
	sort.Strings(allPksList)

	return ret, allPksList, nil
}

func (store *sqlConfiguratorStorage) loadEntityTypeAndKeys(pks []string, loadedEntitiesByPk map[string]*NetworkEntity) (map[string]storage.TypeAndKey, error) {
	ret := map[string]storage.TypeAndKey{}
	pksToLoad := []interface{}{}
	for _, pk := range pks {
		if ent, exists := loadedEntitiesByPk[pk]; exists {
			ret[pk] = storage.TypeAndKey{Type: ent.Type, Key: ent.Key}
		} else {
			pksToLoad = append(pksToLoad, pk)
		}
	}
	// Early exit if we don't need to load anything from DB
	if len(pksToLoad) == 0 {
		return ret, nil
	}
	rows, err := store.builder.Select(entPkCol, entTypeCol, entKeyCol).
		From(entityTable).
		Where(sq.Eq{entPkCol: pksToLoad}).
		RunWith(store.tx).
		Query()
	defer sqorc.CloseRowsLogOnError(rows, "LoadEntities")
	if err != nil {
		return ret, errors.Wrap(err, "failed to query for entity IDs")
	}

	for rows.Next() {
		var pk, t, k string
		err := rows.Scan(&pk, &t, &k)
		if err != nil {
			return ret, errors.Wrap(err, "failed to scan entity ID")
		}
		ret[pk] = storage.TypeAndKey{Type: t, Key: k}
	}

	return ret, nil
}

// entsByPkOut is an output parameter but will also be returned
func updateEntitiesWithAssocs(entsByPkOut map[string]*NetworkEntity, assocs []loadedAssoc, entTksByPk map[string]storage.TypeAndKey, loadCriteria EntityLoadCriteria) (map[string]*NetworkEntity, []GraphEdge, error) {
	retEdges := make([]GraphEdge, 0, len(assocs))
	for _, assoc := range assocs {
		fromTk, fromTkExists := entTksByPk[assoc.fromPk]
		toTk, toTkExists := entTksByPk[assoc.toPk]

		if !fromTkExists && !toTkExists {
			return entsByPkOut, retEdges, errors.Errorf("one end of assoc from %s to %s does not exist", assoc.fromPk, assoc.toPk)
		}
		retEdges = append(retEdges, GraphEdge{From: fromTk, To: toTk})

		// We could load assocs to/from entities that weren't selected for loading
		if loadCriteria.LoadAssocsFromThis {
			fromEnt, exists := entsByPkOut[assoc.fromPk]
			if exists {
				fromEnt.Associations = append(fromEnt.Associations, toTk)
			}
		}
		if loadCriteria.LoadAssocsToThis {
			toEnt, exists := entsByPkOut[assoc.toPk]
			if exists {
				toEnt.ParentAssociations = append(toEnt.ParentAssociations, fromTk)
			}
		}
	}

	sort.Slice(retEdges, func(i, j int) bool { return retEdges[i].String() < retEdges[j].String() })
	for _, ent := range entsByPkOut {
		if loadCriteria.LoadAssocsFromThis {
			sort.Slice(ent.Associations, func(i, j int) bool { return storage.IsTKLessThan(ent.Associations[i], ent.Associations[j]) })
		}
		if loadCriteria.LoadAssocsToThis {
			sort.Slice(ent.ParentAssociations, func(i, j int) bool { return storage.IsTKLessThan(ent.ParentAssociations[i], ent.ParentAssociations[j]) })
		}
	}
	return entsByPkOut, retEdges, nil
}

func calculateEntitiesNotFound(entsByPk map[string]*NetworkEntity, requestedIDs []storage.TypeAndKey) []storage.TypeAndKey {
	if funk.IsEmpty(requestedIDs) {
		return []storage.TypeAndKey{}
	}

	foundIDsMapper := func(pk string, entity *NetworkEntity) (storage.TypeAndKey, struct{}) {
		return storage.TypeAndKey{Type: entity.Type, Key: entity.Key}, struct{}{}
	}
	foundIDsSet := funk.Map(entsByPk, foundIDsMapper).(map[storage.TypeAndKey]struct{})

	ret := []storage.TypeAndKey{}
	for _, requestedID := range requestedIDs {
		_, loaded := foundIDsSet[requestedID]
		if !loaded {
			ret = append(ret, requestedID)
		}
	}
	return ret
}
