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

package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
)

func (m *WifiNetwork) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *NetworkWifiConfigs) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *WifiGateway) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MutableWifiGateway) ValidateModel() error {
	if err := m.Validate(strfmt.Default); err != nil {
		return err
	}

	// Custom validation only for wifi and device
	var res []error
	if err := m.Wifi.ValidateModel(); err != nil {
		res = append(res, err)
	}
	if err := m.Device.ValidateModel(); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GatewayWifiConfigs) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *WifiMesh) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MeshName) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *MeshWifiConfigs) ValidateModel() error {
	return m.Validate(strfmt.Default)
}
