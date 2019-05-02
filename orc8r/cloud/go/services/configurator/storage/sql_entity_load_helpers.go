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
	"sort"
	"strings"
	"text/template"

	"magma/orc8r/cloud/go/sql_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
	"github.com/thoas/go-funk"
)

func (store *sqlConfiguratorStorage) loadFromEntitiesTable(networkID string, filter EntityLoadFilter, criteria EntityLoadCriteria) (map[string]*NetworkEntity, error) {
	// Pointer values because we're modifying entities in-place with ACLs (LEFT JOIN)
	entsByPk := map[string]*NetworkEntity{}

	// If specific IDs are specified, we have to prepare a query, because
	// WHERE IN doesn't work with a composite search index
	if !funk.IsEmpty(filter.IDs) {
		return store.loadSpecificEntities(networkID, filter, criteria)
	}

	// Otherwise, we just fill in the specified filter fields and query once
	queryString, err := getLoadEntitiesQueryString(filter, criteria)
	if err != nil {
		return entsByPk, err
	}
	queryArgs := []interface{}{networkID}
	if filter.KeyFilter != nil {
		queryArgs = append(queryArgs, *filter.KeyFilter)
	}
	if filter.TypeFilter != nil {
		queryArgs = append(queryArgs, *filter.TypeFilter)
	}

	rows, err := store.tx.Query(queryString, queryArgs...)
	if err != nil {
		return entsByPk, fmt.Errorf("error querying for entities: %s", err)
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

func getLoadEntitiesQueryString(filter EntityLoadFilter, criteria EntityLoadCriteria) (string, error) {
	// SELECT ent.pk, ent.key, ent.type, ent.physical_id, ent.version, graph.graph_id, ent.name, ent.description, ent.config,
	// [[ acl.id, acl.scope, acl.permission, acl.type, acl.id_filter, acl.version ]]
	// FROM cfg_entities AS ent
	// INNER JOIN cfg_graphs as graph ON graph.entity_pk = ent.pk
	// [[ LEFT JOIN cfg_acls AS acl ON acl.entity_pk = ent.pk ]]
	// WHERE (ent.network_id, ent.key, ent.type) = ($1, $2, $3)
	queryTemplate := template.Must(template.New("ent_query").Parse(`
		SELECT {{.Fields}} FROM {{.TableName}} AS ent
		{{.ACLJoin}}
		WHERE ({{.WhereCondition}}) = {{.WhereArgList}}
	`))
	queryTemplateArgs := getEntityQueryTemplateArgs(filter, criteria)

	buf := new(bytes.Buffer)
	err := queryTemplate.Execute(buf, queryTemplateArgs)
	if err != nil {
		return "", fmt.Errorf("failed to format entity query: %s", err)
	}
	return buf.String(), nil
}

type entityQueryTemplateArgs struct {
	TableName, Fields, ACLJoin, WhereCondition, WhereArgList string
}

func getEntityQueryTemplateArgs(filter EntityLoadFilter, criteria EntityLoadCriteria) entityQueryTemplateArgs {
	ret := entityQueryTemplateArgs{
		TableName: entityTable,
	}

	fields := []string{"ent.pk", "ent.key", "ent.type", "ent.physical_id", "ent.version", "ent.graph_id"}
	if criteria.LoadMetadata {
		fields = append(fields, "ent.name", "ent.description")
	}
	if criteria.LoadConfig {
		fields = append(fields, "ent.config")
	}
	if criteria.LoadPermissions {
		fields = append(fields, "acl.id", "acl.scope", "acl.permission", "acl.type", "acl.id_filter", "acl.version")
		ret.ACLJoin = fmt.Sprintf(
			"LEFT JOIN %s AS acl ON acl.entity_pk = ent.pk",
			entityAclTable,
		)
	}
	ret.Fields = strings.Join(fields, ", ")

	whereFields := []string{"ent.network_id"}
	if !funk.IsEmpty(filter.IDs) {
		whereFields = append(whereFields, "ent.key", "ent.type")
	} else {
		if filter.KeyFilter != nil {
			whereFields = append(whereFields, "ent.key")
		}
		if filter.TypeFilter != nil {
			whereFields = append(whereFields, "ent.type")
		}
	}
	ret.WhereCondition = strings.Join(whereFields, ", ")
	ret.WhereArgList = sql_utils.GetPlaceholderArgList(1, len(whereFields))

	return ret
}

func (store *sqlConfiguratorStorage) loadSpecificEntities(networkID string, filter EntityLoadFilter, criteria EntityLoadCriteria) (map[string]*NetworkEntity, error) {
	// Pointer values because we're modifying entities in-place with ACLs (LEFT JOIN)
	entsByPk := map[string]*NetworkEntity{}
	if funk.IsEmpty(filter.IDs) {
		return entsByPk, nil
	}

	queryString, err := getLoadEntitiesQueryString(filter, criteria)
	if err != nil {
		return entsByPk, err
	}
	queryStmt, err := store.tx.Prepare(queryString)
	if err != nil {
		return entsByPk, fmt.Errorf("failed to prepare entity query: %s", err)
	}
	defer sql_utils.GetCloseStatementsDeferFunc([]*sql.Stmt{queryStmt}, "LoadEntities")()

	for _, requestedID := range filter.IDs {
		rows, err := queryStmt.Query(networkID, requestedID.Key, requestedID.Type)
		if err != nil {
			return entsByPk, fmt.Errorf("failed to query for entity (%s, %s): %s", requestedID.Type, requestedID.Key, err)
		}

		for rows.Next() {
			err := scanNextEntityRow(rows, criteria, entsByPk)
			if err != nil {
				return entsByPk, fmt.Errorf("error scanning entity (%s, %s): %s", requestedID.Type, requestedID.Key, err)
			}
		}
		if err := rows.Close(); err != nil {
			glog.Errorf("error closing *Rows in LoadEntities: %s", err)
		}
	}

	return entsByPk, nil
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
	// TODO: make the coupling nicer (construct scanArgs and template fields in the same function)
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

	queryString, err := getLoadAssocsQueryString(filter, criteria, len(entsByPk))
	if err != nil {
		return ret, []string{}, err
	}

	entPks := funk.Keys(entsByPk).([]string)
	sort.Strings(entPks)
	queryArgs := make([]interface{}, 0, len(entPks))
	funk.ConvertSlice(entPks, &queryArgs)

	// If we're loading assocs to AND from, then there are 2 WHERE IN
	// arg lists to fill
	if criteria.LoadAssocsToThis && criteria.LoadAssocsFromThis {
		queryArgs = append(queryArgs, queryArgs...)
	}

	assocRows, err := store.tx.Query(queryString, queryArgs...)
	if err != nil {
		return ret, []string{}, fmt.Errorf("error querying for associations: %s", err)
	}
	defer func() {
		if err := assocRows.Close(); err != nil {
			glog.Errorf("error closing *Rows in LoadEntities: %s", err)
		}
	}()
	for assocRows.Next() {
		var fromPk, toPk string
		err := assocRows.Scan(&fromPk, &toPk)
		if err != nil {
			return ret, []string{}, fmt.Errorf("error scanning association row: %s", err)
		}

		ret = append(ret, loadedAssoc{fromPk: fromPk, toPk: toPk})
		allPks[fromPk] = struct{}{}
		allPks[toPk] = struct{}{}
	}

	allPksList := funk.Keys(allPks).([]string)
	sort.Strings(allPksList)

	return ret, allPksList, nil
}

func getLoadAssocsQueryString(filter EntityLoadFilter, criteria EntityLoadCriteria, numLoadedEntities int) (string, error) {
	if !criteria.LoadAssocsToThis && !criteria.LoadAssocsFromThis {
		return "", nil
	}

	// SELECT assoc.from_pk, assoc.to_pk FROM cfg_assocs AS assoc
	// WHERE assoc.from_pk IN ($1, $2, $3, ...)
	// OR assoc.to_pk IN ($4, $5, $6, ...)
	queryTemplate := template.Must(template.New("ent_assoc_query").Parse(`
		SELECT assoc.from_pk, assoc.to_pk FROM {{.TableName}} as assoc
		WHERE {{.WhereClause}} = {{.WhereArgs}}
		{{if .HasOr}}
			OR {{.OrWhereClause}} = {{.OrWhereArgs}}
		{{end}}
	`))
	tmplArgs := getLoadAssocsQueryTemplateArgs(filter, criteria, numLoadedEntities)

	buf := new(bytes.Buffer)
	err := queryTemplate.Execute(buf, tmplArgs)
	if err != nil {
		return "", fmt.Errorf("failed to format assoc query: %s", err)
	}
	return buf.String(), nil
}

type assocsQueryTemplateArgs struct {
	TableName, WhereClause, WhereArgs, OrWhereClause, OrWhereArgs string
	HasOr                                                         bool
}

func getLoadAssocsQueryTemplateArgs(filter EntityLoadFilter, criteria EntityLoadCriteria, numLoadedEntities int) assocsQueryTemplateArgs {
	ret := assocsQueryTemplateArgs{TableName: entityAssocTable}

	whereClauses := getLoadAssocWhereClausesAndArgStrings(filter, criteria, numLoadedEntities)
	ret.WhereClause, ret.WhereArgs = whereClauses[0].clause, whereClauses[0].arg

	ret.HasOr = len(whereClauses) > 1
	if ret.HasOr {
		ret.OrWhereClause, ret.OrWhereArgs = whereClauses[1].clause, whereClauses[1].arg
	}

	return ret
}

type clauseAndArg struct{ clause, arg string }

func getLoadAssocWhereClausesAndArgStrings(filter EntityLoadFilter, criteria EntityLoadCriteria, numLoadedEntities int) []clauseAndArg {
	// If we load all entities, the WHERE clause should just evaluate to true,
	// so make it 1=1
	if filter.IsLoadAllEntities() {
		return []clauseAndArg{{clause: "1", arg: "1"}}
	}

	ret := []clauseAndArg{}
	argStartIdx := 1

	// The IN (...) list for the query needs a placeholder for each loaded entity
	if criteria.LoadAssocsFromThis {
		ret = append(ret, clauseAndArg{
			clause: "assoc.from_pk",
			arg:    sql_utils.GetPlaceholderArgList(argStartIdx, numLoadedEntities),
		})
		argStartIdx += numLoadedEntities
	}
	if criteria.LoadAssocsToThis {
		ret = append(ret, clauseAndArg{
			clause: "assoc.to_pk",
			arg:    sql_utils.GetPlaceholderArgList(argStartIdx, numLoadedEntities),
		})
		argStartIdx += numLoadedEntities
	}

	return ret
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

	query := fmt.Sprintf("SELECT pk, type, key FROM %s WHERE pk IN %s", entityTable, sql_utils.GetPlaceholderArgList(1, len(pksToLoad)))
	rows, err := store.tx.Query(query, pksToLoad...)
	defer sql_utils.CloseRowsLogOnError(rows, "LoadEntities")
	if err != nil {
		return ret, fmt.Errorf("failed to query for entity IDs: %s", err)
	}

	for rows.Next() {
		var pk, t, k string
		err := rows.Scan(&pk, &t, &k)
		if err != nil {
			return ret, fmt.Errorf("error scanning entity ID: %s", err)
		}
		ret[pk] = storage.TypeAndKey{Type: t, Key: k}
	}

	return ret, nil
}

// entsByPkOut is an output parameter but will also be returned
func updateEntitiesWithAssocs(entsByPkOut map[string]*NetworkEntity, assocs []loadedAssoc, entTksByPk map[string]storage.TypeAndKey, loadCriteria EntityLoadCriteria) (map[string]*NetworkEntity, error) {
	for _, assoc := range assocs {
		fromTk, fromTkExists := entTksByPk[assoc.fromPk]
		toTk, toTkExists := entTksByPk[assoc.toPk]

		if !fromTkExists && !toTkExists {
			return entsByPkOut, fmt.Errorf("one end of assoc from %s to %s does not exist", assoc.fromPk, assoc.toPk)
		}

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
	return entsByPkOut, nil
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
