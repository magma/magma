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
	"database/sql"
	"encoding/base64"
	"fmt"
	"sort"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/util"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type loadType int

const (
	countEntities loadType = iota
	loadEntities
	loadChildren
	loadParents
)

func (store *sqlConfiguratorStorage) countEntities(networkID string, filter EntityLoadFilter, criteria EntityLoadCriteria) (uint64, error) {
	selectBuilder, err := store.getBuilder(networkID, filter, criteria, countEntities)
	if err != nil {
		return 0, err
	}
	var count uint64
	if err = selectBuilder.RunWith(store.tx).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (store *sqlConfiguratorStorage) loadEntities(networkID string, filter EntityLoadFilter, criteria EntityLoadCriteria) (EntitiesByTK, error) {
	entsByTK := EntitiesByTK{}

	builder, err := store.getBuilder(networkID, filter, criteria, loadEntities)
	if err != nil {
		return nil, err
	}

	rows, err := builder.RunWith(store.tx).Query()
	if err != nil {
		return nil, errors.Wrap(err, "error querying for entities")
	}
	defer sqorc.CloseRowsLogOnError(rows, "loadEntities")

	for rows.Next() {
		ent, err := scanEntityRow(rows, criteria)
		if err != nil {
			return nil, err
		}
		entsByTK[ent.GetTypeAndKey()] = &ent
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "sql rows err")
	}

	return entsByTK, nil
}

func (store *sqlConfiguratorStorage) loadAssocs(networkID string, filter EntityLoadFilter, criteria EntityLoadCriteria, loadTyp loadType) (loadedAssocs, error) {
	if loadTyp != loadChildren && loadTyp != loadParents {
		return nil, errors.Errorf("wrong load type received: '%v'", loadTyp)
	}

	assocs := loadedAssocs{}

	builder, err := store.getBuilder(networkID, filter, criteria, loadTyp)
	if err != nil {
		return nil, err
	}

	rows, err := builder.RunWith(store.tx).Query()
	if err != nil {
		return nil, errors.Wrap(err, "error querying for entities")
	}
	defer sqorc.CloseRowsLogOnError(rows, "loadAssocs")

	for rows.Next() {
		a, err := scanAssocRow(rows, loadTyp)
		if err != nil {
			return nil, err
		}
		assocs = append(assocs, a)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "sql rows err")
	}

	return assocs, nil
}

func (store *sqlConfiguratorStorage) getBuilder(networkID string, filter EntityLoadFilter, criteria EntityLoadCriteria, loadTyp loadType) (sq.SelectBuilder, error) {
	// Something like:
	//
	// SELECT ent.pk, ent.key, ent.type, ent.physical_id, ent.version, ent.graph_id, ent.name, ent.description, ent.config
	// FROM cfg_entities AS ent
	// [[ JOIN (on cfg_assocs and cfg_entities to get child/parent assocs) ]]
	// [[ WHERE ent.network_id = $network_filter AND ent.key = $key_filter AND ent.type = $type_filter AND ent.key > $page_token ]]
	// ORDER BY ent.key
	// LIMIT $page_size ;

	pageSize := store.getEntityLoadPageSize(criteria)
	pageToken, err := DeserializePageToken(criteria.PageToken)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Helpers
	entCol := func(c string) string {
		return fmt.Sprintf("ent.%s", c)
	}
	addSuffix := func(b sq.SelectBuilder) sq.SelectBuilder {
		switch loadTyp {
		case countEntities:
			return b
		case loadEntities:
			return b.OrderBy(entCol(entKeyCol)).Limit(uint64(pageSize))
		case loadChildren, loadParents:
			return b.OrderBy(entCol(entKeyCol))
		default:
			glog.Errorf("Unsupported configurator SQL load type '%v'", loadTyp)
			return b
		}
	}
	isInNetwork := sq.Eq{entCol(entNidCol): networkID}

	var cols []string
	switch loadTyp {
	case countEntities:
		cols = []string{"COUNT(1)"}
	case loadEntities:
		cols = getLoadEntitiesCols(criteria)
	case loadChildren, loadParents:
		cols = getLoadAssocCols()
	default:
		return sq.SelectBuilder{}, fmt.Errorf("unsupported load type '%v'", loadTyp)
	}

	builder := store.builder.Select(cols...).From(fmt.Sprintf("%s AS ent", entityTable))

	// If we're loading assocs: query should be identical, except need to join on the assocs
	if loadTyp == loadChildren {
		builder = builder.
			Join(fmt.Sprintf("%s ON ent.%s=%s", entityAssocTable, entPkCol, aFrCol)).
			Join(fmt.Sprintf("%s AS assoc ON %s=assoc.%s", entityTable, aToCol, entPkCol))
	} else if loadTyp == loadParents {
		builder = builder.
			Join(fmt.Sprintf("%s ON ent.%s=%s", entityAssocTable, entPkCol, aToCol)).
			Join(fmt.Sprintf("%s AS assoc ON %s=assoc.%s", entityTable, aFrCol, entPkCol))
	}

	// Select only specific TKs
	if funk.NotEmpty(filter.IDs) {
		orClause := make(sq.Or, 0, len(filter.IDs))
		for _, id := range filter.IDs {
			c := sq.And{
				isInNetwork,
				sq.Eq{entCol(entKeyCol): id.Key},
				sq.Eq{entCol(entTypeCol): id.Type},
			}
			orClause = append(orClause, c)
		}
		b := builder.Where(orClause)
		return addSuffix(b), nil
	}

	// Select all with physical ID (generally expected to be unique)
	//
	// Physical ID is the only search not scoped to a network, since we
	// need to be able to look up a caller's network and ent based just on
	// its provided physical ID.
	if filter.PhysicalID != nil {
		b := builder.Where(sq.Eq{entCol(entPidCol): filter.PhysicalID.Value})
		return addSuffix(b), nil
	}

	// Select all with graph ID
	if filter.GraphID != nil {
		b := builder.Where(sq.And{
			isInNetwork,
			sq.Eq{entCol(entGidCol): filter.GraphID.Value},
		})
		return addSuffix(b), nil
	}

	// Default select all

	where := sq.And{isInNetwork}
	// Type, key filters
	if filter.KeyFilter != nil {
		where = append(where, sq.Eq{entCol(entKeyCol): filter.KeyFilter.Value})
	}
	if filter.TypeFilter != nil {
		where = append(where, sq.Eq{entCol(entTypeCol): filter.TypeFilter.Value})
	}
	// Specific page
	if criteria.PageToken != "" {
		where = append(where, sq.Gt{entCol(entKeyCol): pageToken.LastIncludedEntity})
	}

	b := builder.Where(where)
	return addSuffix(b), nil
}

func scanEntityRow(rows *sql.Rows, criteria EntityLoadCriteria) (NetworkEntity, error) {
	var nid, key, entType, graphID, pk string
	var physicalID sql.NullString
	var name, description sql.NullString

	var config []byte
	var version uint64

	// This corresponds with the order of the columns queried in the SELECT
	scanArgs := []interface{}{&nid, &pk, &key, &entType, &physicalID, &version, &graphID}
	if criteria.LoadMetadata {
		scanArgs = append(scanArgs, &name, &description)
	}
	if criteria.LoadConfig {
		scanArgs = append(scanArgs, &config)
	}

	err := rows.Scan(scanArgs...)
	if err != nil {
		return NetworkEntity{}, errors.Wrap(err, "error while scanning entity row")
	}

	ent := NetworkEntity{
		NetworkID: nid,
		Key:       key,
		Type:      entType,

		Name:        nullStringToValue(name),
		Description: nullStringToValue(description),
		PhysicalID:  nullStringToValue(physicalID),
		GraphID:     graphID,
		Pk:          pk,

		Config:  config,
		Version: version,
	}

	return ent, nil
}

func scanAssocRow(rows *sql.Rows, loadTyp loadType) (loadedAssoc, error) {
	a := loadedAssoc{}

	// This corresponds with the order of the columns queried in the SELECT
	// Always loaded as (ent, assoc), so
	// 	- For children => (from, to)
	// 	- For parents => (to, from)
	scanArgs := []interface{}{&a.fromTK.Key, &a.fromTK.Type, &a.toTK.Key, &a.toTK.Type, &a.fromPK, &a.toPK}
	if loadTyp == loadParents {
		scanArgs = []interface{}{&a.toTK.Key, &a.toTK.Type, &a.fromTK.Key, &a.fromTK.Type, &a.toPK, &a.fromPK}
	}

	err := rows.Scan(scanArgs...)
	if err != nil {
		return loadedAssoc{}, errors.Wrap(err, "error while scanning entity row")
	}

	return a, nil
}

func getLoadEntitiesCols(criteria EntityLoadCriteria) []string {
	cols := []string{
		fmt.Sprintf("ent.%s", entNidCol),
		fmt.Sprintf("ent.%s", entPkCol),
		fmt.Sprintf("ent.%s", entKeyCol),
		fmt.Sprintf("ent.%s", entTypeCol),
		fmt.Sprintf("ent.%s", entPidCol),
		fmt.Sprintf("ent.%s", entVerCol),
		fmt.Sprintf("ent.%s", entGidCol),
	}
	if criteria.LoadMetadata {
		cols = append(cols, fmt.Sprintf("ent.%s", entNameCol), fmt.Sprintf("ent.%s", entDescCol))
	}
	if criteria.LoadConfig {
		cols = append(cols, fmt.Sprintf("ent.%s", entConfCol))
	}
	return cols
}

func getLoadAssocCols() []string {
	cols := []string{
		fmt.Sprintf("ent.%s", entKeyCol),
		fmt.Sprintf("ent.%s", entTypeCol),
		fmt.Sprintf("assoc.%s", entKeyCol),
		fmt.Sprintf("assoc.%s", entTypeCol),
		fmt.Sprintf("ent.%s", entPkCol),
		fmt.Sprintf("assoc.%s", entPkCol),
	}
	return cols
}

// getEntityLoadPageSize returns the maximum number of loadEntities to return based
// on the EntityLoadCriteria specified. A page size of 0 will default to the
// maximum load size.
func (store *sqlConfiguratorStorage) getEntityLoadPageSize(loadCriteria EntityLoadCriteria) int {
	if loadCriteria.PageSize == 0 {
		return int(store.maxEntityLoadSize)
	}
	return util.MinInt(int(loadCriteria.PageSize), int(store.maxEntityLoadSize))
}

// updateEntitiesWithAssocs updates entsByTK in-place with the passed assocs.
func updateEntitiesWithAssocs(entsByTK EntitiesByTK, assocs loadedAssocs) ([]*GraphEdge, error) {
	edges := make([]*GraphEdge, 0, len(assocs))
	for _, assoc := range assocs {
		edges = append(edges, assoc.asGraphEdge())

		// Assoc may reference not-loaded ents
		fromEnt, ok := entsByTK[assoc.fromTK]
		if ok {
			fromEnt.Associations = append(fromEnt.Associations, (&EntityID{}).FromTypeAndKey(assoc.toTK))
		}
		toEnt, ok := entsByTK[assoc.toTK]
		if ok {
			toEnt.ParentAssociations = append(toEnt.ParentAssociations, (&EntityID{}).FromTypeAndKey(assoc.fromTK))
		}
	}

	sort.Slice(edges, func(i, j int) bool { return edges[i].ToString() < edges[j].ToString() })
	for _, ent := range entsByTK {
		SortIDs(ent.Associations)
		SortIDs(ent.ParentAssociations)
	}

	return edges, nil
}

func calculateIDsNotFound(entsByTK EntitiesByTK, requestedIDs []*EntityID) []*EntityID {
	if funk.IsEmpty(requestedIDs) {
		return nil
	}

	var foundTKs storage.TKs
	for tk := range entsByTK {
		foundTKs = append(foundTKs, tk)
	}
	var requestedTKs storage.TKs
	for _, id := range requestedIDs {
		requestedTKs = append(requestedTKs, id.ToTypeAndKey())
	}

	missingTKs, _ := requestedTKs.Difference(foundTKs)

	var missingIDs []*EntityID
	for _, tk := range missingTKs {
		missingIDs = append(missingIDs, (&EntityID{}).FromTypeAndKey(tk))
	}

	SortIDs(missingIDs) // for deterministic return

	return missingIDs
}

func getNextPageToken(entities []*NetworkEntity) (string, error) {
	lastEntity := entities[len(entities)-1]
	nextPageToken := &EntityPageToken{LastIncludedEntity: lastEntity.Key}
	return SerializePageToken(nextPageToken)
}

func SerializePageToken(token *EntityPageToken) (string, error) {
	marshalledToken, err := proto.Marshal(token)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(marshalledToken), nil
}

func DeserializePageToken(encodedToken string) (*EntityPageToken, error) {
	marshalledToken, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return nil, err
	}
	token := &EntityPageToken{}
	err = proto.Unmarshal(marshalledToken, token)
	if err != nil {
		return nil, err
	}
	return token, err
}

func validatePaginatedLoadParameters(filter EntityLoadFilter, criteria EntityLoadCriteria) error {
	err := fmt.Errorf("paginated loads cannot be used on multi-type queries")
	if criteria.PageSize != 0 && filter.TypeFilter == nil {
		return err
	}
	if criteria.PageToken != "" && filter.TypeFilter == nil {
		return err
	}
	return nil
}

type loadedAssoc struct {
	fromTK storage.TypeAndKey
	toTK   storage.TypeAndKey
	fromPK string
	toPK   string
}

func (l loadedAssoc) getFromID() *EntityID {
	return (&EntityID{}).FromTypeAndKey(l.fromTK)
}

func (l loadedAssoc) getToID() *EntityID {
	return (&EntityID{}).FromTypeAndKey(l.toTK)
}

func (l loadedAssoc) asGraphEdge() *GraphEdge {
	return &GraphEdge{From: (&EntityID{}).FromTypeAndKey(l.fromTK), To: (&EntityID{}).FromTypeAndKey(l.toTK)}
}

type loadedAssocs []loadedAssoc
