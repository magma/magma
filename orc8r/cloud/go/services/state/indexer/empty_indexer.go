/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package indexer

import (
	"magma/orc8r/cloud/go/services/state/types"
)

type emptyIndexer struct {
	id string
}

// NewEmptyIndexer returns a do-nothing indexer that returns the passed ID
// to GetID() calls.
func NewEmptyIndexer(id string) Indexer {
	return &emptyIndexer{id: id}
}

func (r *emptyIndexer) GetID() string {
	return r.id
}

func (r *emptyIndexer) GetVersion() Version {
	return 0
}

func (r *emptyIndexer) GetTypes() []string {
	return nil
}

func (r *emptyIndexer) PrepareReindex(from, to Version, isFirstReindex bool) error {
	return nil
}

func (r *emptyIndexer) CompleteReindex(from, to Version) error {
	return nil
}

func (r *emptyIndexer) Index(networkID string, states types.SerializedStatesByID) (types.StateErrors, error) {
	return nil, nil
}
