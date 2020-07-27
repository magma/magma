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

package configurator

import (
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/pkg/errors"
)

// why pointer receivers here? we can cache intermediate computation steps
// inside the EntityGraph instance. This allow mconfig building to save some
// search/scan operations.

func (eg *EntityGraph) GetEntity(entType string, key string) (NetworkEntity, error) {
	return eg.GetEntityByTK(storage.TypeAndKey{Type: entType, Key: key})
}

func (eg *EntityGraph) GetEntitiesOfType(entType string) []NetworkEntity {
	res := []NetworkEntity{}
	for _, networkEnt := range eg.Entities {
		if networkEnt.Type == entType {
			res = append(res, networkEnt)
		}
	}
	return res
}

func (eg *EntityGraph) GetEntityByTK(id storage.TypeAndKey) (NetworkEntity, error) {
	eg.cacheGraphHelpers()

	ret, found := eg.entsByTK[id]
	if !found {
		return NetworkEntity{}, merrors.ErrNotFound
	}
	return ret, nil
}

func (eg *EntityGraph) GetFirstAncestorOfType(start NetworkEntity, targetType string) (NetworkEntity, error) {
	eg.cacheGraphHelpers()

	start, found := eg.entsByTK[start.GetTypeAndKey()]
	if !found {
		return NetworkEntity{}, errors.Errorf("entity %s is not in graph", start.GetTypeAndKey())
	}

	ancestor := eg.upwardsDFSForType(start.GetTypeAndKey(), targetType, map[storage.TypeAndKey]bool{})
	if ancestor == nil {
		return NetworkEntity{}, merrors.ErrNotFound
	}
	return *ancestor, nil
}

func (eg *EntityGraph) GetAllChildrenOfType(parent NetworkEntity, targetType string) ([]NetworkEntity, error) {
	eg.cacheGraphHelpers()

	start, found := eg.entsByTK[parent.GetTypeAndKey()]
	if !found {
		return nil, errors.Errorf("entity %s is not in graph", start.GetTypeAndKey())
	}

	ret := []NetworkEntity{}
	for _, tk := range start.Associations {
		if tk.Type == targetType {
			ret = append(ret, eg.entsByTK[tk])
		}
	}
	return ret, nil
}

// backwards DFS search for a type. practically bfs would be more efficient but
// this is easier to implement
func (eg *EntityGraph) upwardsDFSForType(start storage.TypeAndKey, target string, seen map[storage.TypeAndKey]bool) *NetworkEntity {
	if _, alreadySeen := seen[start]; alreadySeen {
		return nil
	}
	if start.Type == target {
		thisEnt := eg.entsByTK[start]
		return &thisEnt
	}

	// mark start as seen, recursively search
	seen[start] = true
	for _, parent := range eg.reverseEdgesByTK[start] {
		ret := eg.upwardsDFSForType(parent, target, seen)
		if ret != nil {
			return ret
		}
	}
	return nil
}

// we may want to granularize the caching behavior if we start caching a lot
// of things
func (eg *EntityGraph) cacheGraphHelpers() {
	if eg.entsByTK != nil && eg.edgesByTK != nil && eg.reverseEdgesByTK != nil {
		return
	}

	eg.entsByTK = map[storage.TypeAndKey]NetworkEntity{}
	for _, ent := range eg.Entities {
		eg.entsByTK[ent.GetTypeAndKey()] = ent
	}

	eg.edgesByTK = map[storage.TypeAndKey][]storage.TypeAndKey{}
	eg.reverseEdgesByTK = map[storage.TypeAndKey][]storage.TypeAndKey{}
	for _, edge := range eg.Edges {
		eg.edgesByTK[edge.From] = append(eg.edgesByTK[edge.From], edge.To)
		eg.reverseEdgesByTK[edge.To] = append(eg.reverseEdgesByTK[edge.To], edge.From)
	}
}

// GetFirstParentOfType iterates through the parent associations of the entity
// and returns the TypeAndKey of first association that matches the target
// type.
// Returns ErrNotFound if no such association is found.
func (ent NetworkEntity) GetFirstParentOfType(target string) (storage.TypeAndKey, error) {
	for _, tk := range ent.ParentAssociations {
		if tk.Type == target {
			return tk, nil
		}
	}
	return storage.TypeAndKey{}, merrors.ErrNotFound
}
