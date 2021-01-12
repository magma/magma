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

package plugin

import (
	"sync"
	"testing"
)

type testPluginRegistry struct {
	sync.Mutex
	plugins map[string]OrchestratorPlugin
}

var testPlugins = &testPluginRegistry{plugins: map[string]OrchestratorPlugin{}}

// RegisterPluginForTests registers all components of a given plugin with the
// corresponding component registries exposed by the orchestrator. This should
// only be used in test code to avoid the cost of building and loading plugins
// from disk for unit tests, thus the required but unused *testing.T
// parameter. This function will not register a plugin which has already been
// registered as identified by its GetName().
func RegisterPluginForTests(_ *testing.T, plugin OrchestratorPlugin) error {
	testPlugins.Lock()
	defer testPlugins.Unlock()
	if _, ok := testPlugins.plugins[plugin.GetName()]; !ok {
		testPlugins.plugins[plugin.GetName()] = plugin
		return registerPlugin(plugin)
	}
	// plugin has already been registered, no-op
	return nil
}
