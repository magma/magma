/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import "magma/feg/cloud/go/protos"

// HealthBlobstore defines a storage interface for the health service. This
// interface defines create/update and read functionality while abstracting
// away any underlying storage transaction mechanics for clients.
type HealthBlobstore interface {
	GetHealth(networkID string, gatewayID string) (*protos.HealthStats, error)

	UpdateHealth(networkID string, gatewayID string, health *protos.HealthStats) error

	GetClusterState(networkID string, clusterID string) (*protos.ClusterState, error)

	UpdateClusterState(networkID string, clusterID string, logicalID string) error
}
