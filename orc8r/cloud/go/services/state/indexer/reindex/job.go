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

package reindex

import (
	"fmt"

	"magma/orc8r/cloud/go/services/state/indexer"
)

// Job required to carry out a reindex job.
type Job struct {
	Idx  indexer.Indexer
	From indexer.Version
	To   indexer.Version
}

func (j *Job) String() string {
	return fmt.Sprintf("{id: %s, from: %d, to: %d}", j.Idx.GetID(), j.From, j.To)
}
