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

// CtracedStorage is the persistence service interface for call traces.
// All call trace accesses from ctraced service must go through this interface.
// Call traces should be stored as .pcap files
type CtracedStorage interface {

	// StoreCallTrace stores the call trace file
	StoreCallTrace(networkID string, callTraceID string, data []byte) error

	// GetCallTrace returns the call trace file
	GetCallTrace(networkID string, callTraceID string) ([]byte, error)

	// DeleteCallTrace deletes the call trace file
	// Returns an error if the call trace does not exist
	DeleteCallTrace(networkID string, callTraceID string) error
}
