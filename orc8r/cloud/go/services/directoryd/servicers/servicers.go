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

package servicers

import (
	directoryd_protos "magma/orc8r/cloud/go/services/directoryd/protos"
	"magma/orc8r/cloud/go/services/directoryd/servicers/internal"
	"magma/orc8r/cloud/go/services/directoryd/storage"
)

func NewDirectoryLookupServicer(store storage.DirectorydStorage) (directoryd_protos.DirectoryLookupServer, error) {
	return internal.NewDirectoryLookupServicer(store)
}

func NewDirectoryUpdateServicer() directoryd_protos.GatewayDirectoryServiceServer {
	return internal.NewDirectoryUpdateServicer()
}
