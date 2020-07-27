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

package main

import (
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
)

// plugins must implement a main - these are expected to be empty
func main() {}

// GetOrchestratorPlugin is a function that all modules are expected to provide
// which returns an instance of the module's OrchestratorPlugin implementation
func GetOrchestratorPlugin() plugin.OrchestratorPlugin {
	return &pluginimpl.BaseOrchestratorPlugin{}
}
