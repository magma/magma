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

package mocks

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/services/state/indexer"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"

	"github.com/stretchr/testify/mock"
)

// NewMockIndexer returns a do-nothing test indexer with specified elements.
// 	- id		-- GetID return
//	- version	-- GetVersion return
//	- types		-- GetTypes return
//	- prepare	-- write PrepareReindex args to chan when called
//	- complete	-- write CompleteReindex args to chan when called
//	- index		-- write Index args to chan when called
func NewMockIndexer(
	t *testing.T,
	id string,
	version indexer.Version,
	types []string,
	prepare,
	complete,
	index chan mock.Arguments,
) (remoteIndexer indexer.Indexer, mockIndexer *Indexer) {
	mockIndexer = &Indexer{}
	mockIndexer.On("GetID").Return(id)
	mockIndexer.On("GetVersion").Return(version)
	mockIndexer.On("GetTypes").Return(types)
	mockIndexer.On("PrepareReindex", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if prepare != nil {
			prepare <- args
		}
	}).Return(nil)
	mockIndexer.On("CompleteReindex", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if complete != nil {
			complete <- args
		}
	}).Return(nil)
	mockIndexer.On("Index", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if index != nil {
			index <- args
		}
	}).Return(nil, nil)

	// Indexer servicer goroutine doesn't get canceled, but we don't care
	// since its service name gets overridden
	state_test_init.StartNewTestIndexer(t, mockIndexer)
	remoteIndexer = indexer.NewRemoteIndexer(id, version, types...)

	return remoteIndexer, mockIndexer
}

func (_m *Indexer) String() string {
	return fmt.Sprintf("{id: %s, version: %d}", _m.GetID(), _m.GetVersion())
}
