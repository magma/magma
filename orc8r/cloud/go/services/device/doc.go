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

// Package device contains the device service.
// The device service is a simple blob-storage service for tracking
// physical device.
package device

// SerdeDomain is the domain for all Serde implementations for the device
// service
const (
	SerdeDomain = "device"

	// ServiceName is the name of this service
	ServiceName = "DEVICE"

	// DBTableName is the name of the sql table used for this service
	DBTableName = "device"
)
