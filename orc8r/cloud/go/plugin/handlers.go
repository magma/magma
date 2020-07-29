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

import "magma/orc8r/cloud/go/obsidian"

// FlattenHandlerLists turns a variadic list of obsidian handlers into a
// single flattened list of handlers. This is typically used to merge handlers
// from different services into a single collection to return in an impl
// of OrchestratorPlugin.
func FlattenHandlerLists(handlersIn ...[]obsidian.Handler) []obsidian.Handler {
	var ret []obsidian.Handler
	for _, h := range handlersIn {
		ret = append(ret, h...)
	}
	return ret
}
