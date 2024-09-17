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

package protos

import (
	"fmt"
)

func (m *RegisterOrUpdateDevicesRequest) Validate() error {
	if err := nonEmptyNetworkID(m.GetNetworkID()); err != nil {
		return err
	}
	if len(m.GetEntities()) == 0 {
		return fmt.Errorf("entities field must be non-empty")
	}
	return nil
}

func (m *GetDeviceInfoRequest) Validate() error {
	return nonEmptyNetworkIDAndDeviceIDs(m.GetNetworkID(), m.GetDeviceIDs())
}

func (m *DeleteDevicesRequest) Validate() error {
	return nonEmptyNetworkIDAndDeviceIDs(m.GetNetworkID(), m.GetDeviceIDs())
}

func nonEmptyNetworkID(networkID string) error {
	if len(networkID) == 0 {
		return fmt.Errorf("network ID must be non-empty")
	}
	return nil
}

func nonEmptyNetworkIDAndDeviceIDs(networkID string, deviceIDs []*DeviceID) error {
	if err := nonEmptyNetworkID(networkID); err != nil {
		return err
	}
	if len(deviceIDs) == 0 {
		return fmt.Errorf("device IDs field must be non-empty")
	}
	return nil
}
