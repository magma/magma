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
	"sort"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type internalEntityGraph struct {
	entsByTK EntitiesByTK
	edges    loadedAssocs
}

// loadGraphInternal will load all loadEntities and assocs for a given graph ID.
// This function will NOT fill loadEntities with associations.
func (store *sqlConfiguratorStorage) loadGraphInternal(networkID string, graphID string, criteria EntityLoadCriteria) (internalEntityGraph, error) {
	loadFilter := EntityLoadFilter{GraphID: &wrappers.StringValue{Value: graphID}}
	ents, err := store.loadEntities(networkID, loadFilter, criteria)
	if err != nil {
		return internalEntityGraph{}, errors.Wrap(err, "failed to load entities for graph")
	}
	if funk.IsEmpty(ents) {
		return internalEntityGraph{}, nil
	}

	// Just loading children is sufficient, since this will load all assocs
	// in the graph
	assocs, err := store.loadAssocs(networkID, loadFilter, criteria, loadChildren)
	if err != nil {
		return internalEntityGraph{}, errors.Wrap(err, "error loading child edges for graph")
	}

	return internalEntityGraph{entsByTK: ents, edges: assocs}, nil
}

// fixGraph will load a graph which may have been partitioned, do a connected
// component search on it, and relabel components if a partition is detected.
// entToUpdateOut is an output parameter
func (store *sqlConfiguratorStorage) fixGraph(networkID string, graphID string, entToUpdateOut *NetworkEntity) error {
	internalGraph, err := store.loadGraphInternal(networkID, graphID, EntityLoadCriteria{})
	if err != nil {
		return errors.Wrap(err, "failed to load graph of updated entity")
	}
	if funk.IsEmpty(internalGraph.entsByTK) {
		return nil
	}

	edges := map[string][]string{}
	funk.ForEach(internalGraph.edges, func(assoc loadedAssoc) { edges[assoc.fromPK] = append(edges[assoc.fromPK], assoc.toPK) })

	// Do a looped DFS over the graph to record connected components in a
	// union-find structure
	allPKs := internalGraph.entsByTK.PKs()
	seenPKs := map[string]bool{}
	uf := newUnionFind(allPKs)
	for _, pk := range allPKs {
		dfsFrom(pk, edges, seenPKs, uf)
	}

	// If the graph is fully connected, we don't have to do anything.
	// Otherwise, we need to generate new graph IDs for all components except
	// the last (largest) one in the list
	components := uf.getComponents()
	for _, component := range components[:len(components)-1] {
		newID := store.idGenerator.New()
		err := store.updateGraphID(component, newID)
		if err != nil {
			return errors.Wrap(err, "failed to fix graph")
		}
	}
	return nil
}

func (store *sqlConfiguratorStorage) updateGraphID(pksToUpdate []string, newGraphID string) error {
	sort.Strings(pksToUpdate)
	_, err := store.builder.Update(entityTable).
		Set(entGidCol, newGraphID).
		Where(sq.Eq{entPkCol: pksToUpdate}).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return errors.Wrap(err, "failed to update graph ID")
	}
	return nil
}

func dfsFrom(startPk string, edges map[string][]string, seenPKsOut map[string]bool, ufOut *unionFind) {
	if _, seen := seenPKsOut[startPk]; seen {
		return
	}

	// mark this node as seen before continuing
	// for each neighbor, union it with this node and recurse
	seenPKsOut[startPk] = true
	for _, nextPk := range edges[startPk] {
		ufOut.union(startPk, nextPk)
		dfsFrom(nextPk, edges, seenPKsOut, ufOut)
	}
}

func findRootNodes(graph internalEntityGraph) []string {
	// A root node is a node with no edges terminating at it.
	// To make computation easier, we'll make a reverse adjacency list first
	reverseEdges := map[string][]string{}
	for _, edge := range graph.edges {
		reverseEdges[edge.toPK] = append(reverseEdges[edge.toPK], edge.fromPK)
	}

	// Set of root nodes is the set difference between node IDs and reverse
	// edge keyset.
	rootNodes := funk.Filter(
		graph.entsByTK.PKs(),
		func(pk string) bool {
			_, hasEdgesTo := reverseEdges[pk]
			return !hasEdgesTo
		},
	).([]string)
	sort.Strings(rootNodes)
	return rootNodes
}
