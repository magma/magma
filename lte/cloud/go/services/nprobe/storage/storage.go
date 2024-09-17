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

package storage

import "magma/lte/cloud/go/services/nprobe/obsidian/models"

// NProbeStorage is the storage interface to manage nprobe service state.
type NProbeStorage interface {
	// StoreNProbeData stores current state for a given networkID and taskID
	StoreNProbeData(networkID, taskID string, data models.NetworkProbeData) error

	// GetNProbeData returns the state keyed by networkID and taskID
	GetNProbeData(networkID, taskID string) (*models.NetworkProbeData, error)

	// DeleteNProbeData deletes a state for a given networkID and taskID
	DeleteNProbeData(networkID, taskID string) error
}
