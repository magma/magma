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

// DirectorydStorage is the persistence service interface for location records.
// All Directoryd data accesses from directoryd service must go through this interface.
type DirectorydStorage interface {
	// GetHostnameForHWID returns the hostname mapped to by hardware ID.
	GetHostnameForHWID(hwid string) (string, error)

	// MapHWIDsToHostnames maps {hwid -> hostname}.
	MapHWIDsToHostnames(hwidToHostname map[string]string) error

	// GetIMSIForSessionID returns the IMSI mapped to by session ID.
	GetIMSIForSessionID(networkID, sessionID string) (string, error)

	// MapSessionIDsToIMSIs maps {session ID -> IMSI}.
	MapSessionIDsToIMSIs(networkID string, sessionIDToIMSI map[string]string) error

	// GetHWIDForSgwCTeid returns the HwId mapped by teid
	GetHWIDForSgwCTeid(networkID, teid string) (string, error)

	// MapSgwCTeidToHWID maps {teid -> hwid}
	MapSgwCTeidToHWID(networkID string, s8TeidToHwid map[string]string) error
}
