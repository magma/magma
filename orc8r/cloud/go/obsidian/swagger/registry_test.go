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

package swagger_test

import (
	"testing"

	"magma/orc8r/cloud/go/obsidian/swagger"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
)

func TestRegistry_GetSpecServicers(t *testing.T) {
	// Success with no registered servicers
	servicers, err := swagger.GetSpecServicers()
	assert.NoError(t, err)

	assert.Equal(t, 0, len(servicers))

	// Success with some registered servicers
	services := []string{"test_spec_service1", "test_spec_service2"}
	labels := map[string]string{
		orc8r.SpecServicerLabel: "true",
	}

	for _, s := range services {
		srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, s, labels, nil)
		go srv.RunTest(lis)
	}

	servicers, err = swagger.GetSpecServicers()
	assert.NoError(t, err)

	assert.Equal(t, services[0], servicers[0].GetService())
	assert.Equal(t, services[1], servicers[1].GetService())

	// Clean up registry
	registry.RemoveService(services[0])
	registry.RemoveService(services[1])
}
