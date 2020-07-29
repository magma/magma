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
	"errors"
)

// ValidateGatewayConfig - Validate a DevmandGatewayConfig
func ValidateGatewayConfig(config *DevmandGatewayConfig) error {
	if config == nil {
		return errors.New("Gateway config is nil")
	}
	return nil
}

// ValidateManagedDevice - validate a ManagedDevice
func ValidateManagedDevice(config *ManagedDevice) error {
	if config == nil {
		return errors.New("Device config is nil")
	}
	return nil
}
