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
	"strings"
)

func newUnionFind(pks []string) *unionFind {
	ret := &unionFind{
		parents: make(map[string]string, len(pks)),
		ranks:   make(map[string]int, len(pks)),
	}
	for _, pk := range pks {
		ret.parents[pk] = pk
		ret.ranks[pk] = 0
	}
	return ret
}

// unionFind is an implementation of a union-find data structure for efficient
// computation of connected components in a graph. The DS implements path
// compression on find and rank-rated union.
// Tutorial: https://www.youtube.com/watch?v=0jNmHPfA_yE
type unionFind struct {
	parents map[string]string
	ranks   map[string]int
}

// find finds the root of the given pk and compresses the path
func (uf *unionFind) find(pk string) string {
	p := uf.parents[pk]
	if p != pk {
		uf.parents[pk] = uf.find(p)
	}
	return uf.parents[pk]
}

func (uf *unionFind) union(x, y string) {
	xRoot, yRoot := uf.find(x), uf.find(y)
	if xRoot == yRoot {
		return
	}

	// To make things simpler, we invariantly make the rank of yRoot's tree
	// smaller than xRoot so we can always merge yRoot into xRoot
	xRank, yRank := uf.ranks[xRoot], uf.ranks[yRoot]
	if xRank < yRank {
		xRoot, yRoot = yRoot, xRoot
	}

	// The only time we need to update rank is if the trees were equal size,
	// since merging a smaller tree into a larger one doesn't affect the
	// larger one's rank.
	uf.parents[yRoot] = xRoot
	if xRank == yRank {
		uf.ranks[xRoot]++
	}
}

// getComponents returns the connected components in the data structure, sorted
// by length and tiebroken by string join
func (uf *unionFind) getComponents() [][]string {
	// do a find() on everything to fully compress all paths
	for pk := range uf.parents {
		uf.find(pk)
	}

	// now we can invert the parent map - each root points to a list of
	// child pk's and itself since root nodes are self-referential in parents
	inverseMap := map[string][]string{}
	for pk, parent := range uf.parents {
		inverseMap[parent] = append(inverseMap[parent], pk)
	}

	// construct the return value
	ret := make([][]string, 0, len(inverseMap))
	for _, children := range inverseMap {
		sort.Strings(children)
		ret = append(ret, children)
	}
	sort.Slice(
		ret,
		func(i, j int) bool {
			if len(ret[i]) == len(ret[j]) {
				return strings.Join(ret[i], "") < strings.Join(ret[j], "")
			}
			return len(ret[i]) < len(ret[j])
		},
	)
	return ret
}
