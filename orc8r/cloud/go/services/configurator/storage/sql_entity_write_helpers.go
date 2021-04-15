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
	"fmt"
	"reflect"
	"sort"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func (store *sqlConfiguratorStorage) doesEntExist(networkID string, tk storage.TypeAndKey) (bool, error) {
	var count uint64
	err := store.builder.Select("COUNT(1)").
		From(entityTable).
		Where(sq.And{
			sq.Eq{entNidCol: networkID},
			sq.Eq{entTypeCol: tk.Type},
			sq.Eq{entKeyCol: tk.Key},
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

func (store *sqlConfiguratorStorage) doesPhysicalIDExist(physicalID string) (bool, error) {
	if physicalID == "" {
		return false, nil
	}

	var count uint64
	err := store.builder.Select("COUNT(1)").
		From(entityTable).
		Where(sq.Eq{entPidCol: physicalID}).
		RunWith(store.tx).
		QueryRow().Scan(&count)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "check for existence of physical ID %s", physicalID)
	}

	return count > 0, nil
}

func (store *sqlConfiguratorStorage) insertIntoEntityTable(networkID string, ent NetworkEntity) (NetworkEntity, error) {
	ent.Pk = store.idGenerator.New()
	ent.GraphID = store.idGenerator.New() // potentially-temporary graph ID

	_, err := store.builder.Insert(entityTable).
		Columns(entPkCol, entNidCol, entTypeCol, entKeyCol, entGidCol, entNameCol, entDescCol, entPidCol, entConfCol).
		Values(ent.Pk, networkID, ent.Type, ent.Key, ent.GraphID, toNullable(ent.Name), toNullable(ent.Description), toNullable(ent.PhysicalID), toNullable(ent.Config)).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return NetworkEntity{}, errors.Wrapf(err, "error creating entity %s", ent.GetTypeAndKey())
	}
	return ent, nil
}

func (store *sqlConfiguratorStorage) createEdges(networkID string, entity NetworkEntity) (EntitiesByTK, error) {
	// Load the associated entities first because we need to know PKs
	// This will also load graph ID on the entity because creating an edge can
	// involve merging previously disjoint graphs.
	if funk.IsEmpty(entity.GetGraphEdges()) {
		return nil, nil
	}

	// Get assoc pks, since we don't trust the pks provided by the input ent
	entsByTk, err := store.loadEntsFromEdges(networkID, entity)
	if err != nil {
		return nil, err
	}

	insertBuilder := store.builder.Insert(entityAssocTable).
		Columns(aFrCol, aToCol).
		OnConflict(nil, aFrCol, aToCol)
	for _, edge := range entity.GetGraphEdges() {
		fromPk := entsByTk[edge.From.ToTypeAndKey()].Pk
		toPk := entsByTk[edge.To.ToTypeAndKey()].Pk
		insertBuilder = insertBuilder.Values(fromPk, toPk)
	}
	_, err = insertBuilder.RunWith(store.tx).Exec()
	if err != nil {
		return entsByTk, errors.Wrap(err, "error creating assocs")
	}
	return entsByTk, nil
}

func (store *sqlConfiguratorStorage) loadEntsFromEdges(networkID string, targetEntity NetworkEntity) (EntitiesByTK, error) {
	loadedEntsByTk, err := store.loadEntitiesFromIDs(networkID, targetEntity.Associations)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	loadedEntsByTk[targetEntity.GetTypeAndKey()] = &targetEntity
	return loadedEntsByTk, nil
}

func (store *sqlConfiguratorStorage) loadEntitiesFromIDs(networkID string, idsToLoad []*EntityID) (EntitiesByTK, error) {
	loaded, err := store.loadEntities(networkID, EntityLoadFilter{IDs: idsToLoad}, EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}

	entsNotFound := calculateIDsNotFound(loaded, idsToLoad)
	if funk.NotEmpty(entsNotFound) {
		return nil, errors.Errorf("could not find entities matching %v", entsNotFound)
	}

	return loaded, nil
}

func (store *sqlConfiguratorStorage) mergeGraphs(createdEntity NetworkEntity, allAssociatedEntsByTk EntitiesByTK) (string, error) {
	// If we create a node or edge which bridges 2 previously disjoint graphs,
	// then we need to change the ID of one of the graphs to the joined one.

	// If we associate to no graphs, then no-op - we'll use the
	// system-generated graph ID for this single-node graph.

	// Otherwise, we'll take the lexicographically smallest graph ID to keep
	// and change the graph ID of every entity of the other graphs to this
	// target graph ID.
	adjacentGraphs := funk.Chain(createdEntity.Associations).
		Map(func(id *EntityID) string { return allAssociatedEntsByTk[id.ToTypeAndKey()].GraphID }).
		Uniq().
		Value().([]string)
	noMergeNecessary := funk.IsEmpty(adjacentGraphs) || (len(adjacentGraphs) == 1 && adjacentGraphs[0] == createdEntity.GraphID)
	if noMergeNecessary {
		return createdEntity.GraphID, nil
	}

	if !funk.ContainsString(adjacentGraphs, createdEntity.GraphID) {
		adjacentGraphs = append(adjacentGraphs, createdEntity.GraphID)
	}
	sort.Strings(adjacentGraphs)
	targetGraphID := adjacentGraphs[0]
	graphIDsToChange := adjacentGraphs[1:]

	// let squirrel cache prepared statements for us (there should only be 1)
	sc := sq.NewStmtCache(store.tx)
	defer sqorc.ClearStatementCacheLogOnError(sc, "mergeGraphs")

	for _, oldGraphID := range graphIDsToChange {
		_, err := store.builder.Update(entityTable).
			Set(entGidCol, targetGraphID).
			Where(sq.Eq{entGidCol: oldGraphID}).
			RunWith(sc).
			Exec()
		if err != nil {
			return "", errors.Wrap(err, "error updating entity graphs")
		}
	}

	return targetGraphID, nil
}

func (store *sqlConfiguratorStorage) loadEntToUpdate(networkID string, update EntityUpdateCriteria) (*NetworkEntity, error) {
	loaded, err := store.loadEntities(
		networkID,
		EntityLoadFilter{IDs: []*EntityID{update.GetID()}},
		EntityLoadCriteria{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load entity to update")
	}
	// don't error on deleting an entity which doesn't exist
	if len(loaded) != 1 && !update.DeleteEntity {
		return nil, errors.Errorf("expected to load 1 ent for update, got %d", len(loaded))
	}

	if funk.IsEmpty(loaded) {
		return nil, nil
	}

	// Return one ent
	for _, ent := range loaded {
		return ent, nil
	}
	return nil, nil // to appease compiler
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
		entOut.Name = (*update.NewName).Value
	}
	if update.NewDescription != nil {
		entOut.Description = (*update.NewDescription).Value
	}
	if update.NewPhysicalID != nil {
		entOut.PhysicalID = (*update.NewPhysicalID).Value
	}
	if update.NewConfig != nil {
		entOut.Config = (*update.NewConfig).Value
	}
	entOut.Version++

	return nil
}

// entToUpdateOut is an output parameter
func (store *sqlConfiguratorStorage) processEdgeUpdates(networkID string, update EntityUpdateCriteria, entToUpdateOut *NetworkEntity) error {
	assocsToSetSpecified := update.AssociationsToSet != nil
	if !assocsToSetSpecified && funk.IsEmpty(update.AssociationsToAdd) && funk.IsEmpty(update.AssociationsToDelete) {
		return nil
	}

	// If we want to set associations all at once, we'll first delete all
	// associations
	if assocsToSetSpecified {
		_, err := store.builder.Delete(entityAssocTable).
			Where(sq.Eq{aFrCol: entToUpdateOut.Pk}).
			RunWith(store.tx).
			Exec()
		if err != nil {
			return errors.Wrap(err, "failed to delete existing edges")
		}
	}

	// First, create edges. Because createEdges expects an ent with its pk set,
	// we'll just make the ent's Associations the edges we want to create
	// If we want to set associations, we'll create those
	entToUpdateOut.Associations = update.getEdgesToCreate()
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

	// Now delete edges
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
	if funk.IsEmpty(update.AssociationsToDelete) && update.AssociationsToSet == nil {
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
	updateBuilder := store.builder.Update(entityTable).Where(sq.Eq{entPkCol: pk})
	if update.NewName != nil {
		updateBuilder = updateBuilder.Set(entNameCol, update.NewName.Value)
	}
	if update.NewDescription != nil {
		updateBuilder = updateBuilder.Set(entDescCol, update.NewDescription.Value)
	}
	if update.NewPhysicalID != nil {
		updateBuilder = updateBuilder.Set(entPidCol, update.NewPhysicalID.Value)
	}
	if update.NewConfig != nil {
		updateBuilder = updateBuilder.Set(entConfCol, update.NewConfig.Value)
	}
	updateBuilder = updateBuilder.Set(entVerCol, sq.Expr(fmt.Sprintf("%s+1", entVerCol)))
	return updateBuilder
}

// entToUpdateOut is an output parameter
func (store *sqlConfiguratorStorage) deleteEdges(networkID string, edgesToDelete []*EntityID, entToUpdateOut *NetworkEntity) error {
	if funk.IsEmpty(edgesToDelete) {
		return nil
	}

	// Get PKs to delete
	loaded, err := store.loadEntitiesFromIDs(networkID, edgesToDelete)
	if err != nil {
		return errors.Wrap(err, "could not load entities matching associations to delete")
	}

	orClause := make(sq.Or, 0, len(edgesToDelete))
	for _, edge := range edgesToDelete {
		orClause = append(orClause, sq.And{
			sq.Eq{aFrCol: entToUpdateOut.Pk},
			sq.Eq{aToCol: loaded[edge.ToTypeAndKey()].Pk},
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
	edgesToDeleteSet := funk.Map(
		edgesToDelete,
		func(id *EntityID) (storage.TypeAndKey, bool) { return id.ToTypeAndKey(), true },
	).(map[storage.TypeAndKey]bool)
	entToUpdateOut.Associations = funk.Filter(entToUpdateOut.Associations, func(id *EntityID) bool {
		_, wasDeleted := edgesToDeleteSet[id.ToTypeAndKey()]
		return !wasDeleted
	}).([]*EntityID)

	return nil
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
