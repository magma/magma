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

// File registry.go provides an mconfig builder registry by forwarding calls to
// the service registry.

package mconfig

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/registry"
)

// GetBuilders returns all registered mconfig builders.
func GetBuilders() ([]Builder, error) {
	services, err := registry.FindServices(orc8r.MconfigBuilderLabel)
	if err != nil {
		return []Builder{}, err
	}
	var builders []Builder
	for _, s := range services {
		builders = append(builders, NewRemoteBuilder(s))
	}

	return builders, nil
}
