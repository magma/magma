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

// Package registry defines API to register and fing EAP providers
package registry

import (
	aka_provider "magma/feg/gateway/services/eap/providers/aka/provider"
	sim_provider "magma/feg/gateway/services/eap/providers/sim/provider"
)

func init() {
	Register(aka_provider.New())
	Register(sim_provider.New())
}
