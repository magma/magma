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

package swagger

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/registry"
)

// GetSpecServicers returns all registered Swagger spec servicers.
func GetSpecServicers() ([]RemoteSpec, error) {
	services, err := registry.FindServices(orc8r.SwaggerSpecLabel)
	if err != nil {
		return nil, err
	}

	var servicers []RemoteSpec
	for _, s := range services {
		servicers = append(servicers, NewRemoteSpec(s))
	}
	return servicers, nil
}
